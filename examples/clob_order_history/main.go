package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/joho/godotenv"
)

// 获取“历史订单”的一种实用做法：
// 1) 先拉取账户的历史成交（trades）
// 2) 从 trades 中提取 order hash（taker_order_id / maker_orders[].order_id）
// 3) 再用 GetOrder 拉取订单详情
//
// 注意：
// - 这种方式只能覆盖“发生过成交”的订单；纯撤单且无成交的订单不会出现在 trades 里。
// - 为避免拉取过多数据，默认只拉取前几页 trades，并限制最多查询 N 个 order。
//
// 环境变量：
// - POLYMARKET_ADDRESS（必填）
// - POLYMARKET_API_KEY / POLYMARKET_API_SECRET / POLYMARKET_PASSPHRASE（必填）
// - POLYMARKET_MARKET_ID（可选：只看某个市场）
// - MAX_PAGES（可选：默认 2）
// - MAX_ORDERS（可选：默认 20）
func main() {
	_ = godotenv.Load()

	address := os.Getenv("POLYMARKET_ADDRESS")
	if address == "" {
		fmt.Println("missing POLYMARKET_ADDRESS")
		return
	}

	cfg := pm.Config{
		Address:    address,
		APIKey:     os.Getenv("POLYMARKET_API_KEY"),
		APISecret:  os.Getenv("POLYMARKET_API_SECRET"),
		Passphrase: os.Getenv("POLYMARKET_PASSPHRASE"),
		Proxy:      os.Getenv("PROXY"),
	}

	sdk, err := pm.New(cfg)
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		return
	}

	marketID := os.Getenv("POLYMARKET_MARKET_ID")
	maxPages := envInt("MAX_PAGES", 2)
	maxOrders := envInt("MAX_ORDERS", 20)

	ctx := context.Background()

	// trades：分别以 maker / taker 拉取，再去重
	trades := make(map[string]*pm.Trade)
	fetch := func(req *pm.GetTradesRequest) error {
		next := pm.InitialCursor
		for page := 0; page < maxPages && next != pm.EndCursor; page++ {
			resp, err := sdk.CLOB.GetTradesPage(ctx, req, next)
			if err != nil {
				return err
			}
			for _, t := range resp.Data {
				if t != nil && t.ID != "" {
					trades[t.ID] = t
				}
			}
			next = resp.NextCursor
		}
		return nil
	}

	if err := fetch(&pm.GetTradesRequest{Market: marketID, Maker: address}); err != nil {
		fmt.Printf("get trades (maker) failed: %v\n", err)
		return
	}
	if err := fetch(&pm.GetTradesRequest{Market: marketID, Taker: address}); err != nil {
		fmt.Printf("get trades (taker) failed: %v\n", err)
		return
	}

	fmt.Printf("trades loaded: %d\n", len(trades))

	// 提取 order hashes
	orderIDs := make(map[string]struct{})
	for _, t := range trades {
		if t == nil {
			continue
		}
		if t.TakerOrderID != "" {
			orderIDs[t.TakerOrderID] = struct{}{}
		}
		for _, mo := range t.MakerOrders {
			if mo.OrderID != "" {
				orderIDs[mo.OrderID] = struct{}{}
			}
		}
	}

	fmt.Printf("order ids (from trades): %d\n", len(orderIDs))

	// 拉取订单详情（限制数量）
	n := 0
	for id := range orderIDs {
		if n >= maxOrders {
			break
		}
		o, err := sdk.CLOB.GetOrder(ctx, id)
		if err != nil {
			fmt.Printf("get order failed: id=%s err=%v\n", id, err)
			continue
		}
		fmt.Printf("order id=%s status=%s side=%s price=%s size=%s market=%s asset_id=%s\n", o.ID, o.Status, o.Side, o.Price, o.Size, o.Market, o.AssetID)
		n++
	}
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return def
}
