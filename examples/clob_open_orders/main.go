package main

import (
	"context"
	"fmt"
	"os"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/joho/godotenv"
)

// 获取账户在某个“持仓/市场”（market 或 asset_id）上的活跃订单（open orders） 看“还挂着的单”
//
// 可选环境变量：
// - POLYMARKET_MARKET_ID：condition_id
// - POLYMARKET_ASSET_ID：clob 的 asset_id（token_id）
func main() {
	_ = godotenv.Load()

	cfg := pm.Config{
		Address:    os.Getenv("POLYMARKET_ADDRESS"),
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

	req := &pm.GetActiveOrdersRequest{
		Market:  os.Getenv("POLYMARKET_MARKET_ID"),
		AssetID: os.Getenv("POLYMARKET_ASSET_ID"),
	}

	orders, err := sdk.CLOB.GetOpenOrders(context.Background(), req)
	if err != nil {
		fmt.Printf("get open orders failed: %v\n", err)
		return
	}

	fmt.Printf("open orders: %d\n", len(orders))
	for i, o := range orders {
		fmt.Printf(
			"#%d id=%s status=%s side=%s price=%s size=%s market=%s asset_id=%s\n",
			i, o.ID, o.Status, o.Side, o.Price, o.Size, o.Market, o.AssetID,
		)
	}
}
