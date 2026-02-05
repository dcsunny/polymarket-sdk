package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/joho/godotenv"
)

// 获取 CLOB 市场列表示例
// 演示如何使用 GetMarkets 获取 CLOB 服务下的市场数据

func main() {
	_ = godotenv.Load()

	cfg := pm.Config{
		APIKey:     os.Getenv("POLYMARKET_API_KEY"),
		APISecret:  os.Getenv("POLYMARKET_API_SECRET"),
		Passphrase: os.Getenv("POLYMARKET_PASSPHRASE"),
		Address:    os.Getenv("POLYMARKET_ADDRESS"),
		Proxy:      os.Getenv("PROXY"),
	}

	sdk, err := pm.New(cfg)
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		return
	}

	// 第一次请求，使用初始游标
	nextCursor := pm.InitialCursor
	pageNum := 1
	maxPages := 3 // 最多获取 3 页数据，避免输出太多

	for {
		if pageNum > maxPages {
			fmt.Printf("\n已达到最大页数限制 (%d 页)，停止获取\n", maxPages)
			break
		}

		fmt.Printf("\n========== 第 %d 页 ==========\n", pageNum)

		// 获取市场列表
		resp, err := sdk.CLOB.GetMarkets(context.Background(), nextCursor)
		if err != nil {
			fmt.Printf("get markets failed: %v\n", err)
			return
		}

		// 显示分页信息
		fmt.Printf("Limit: %d\n", resp.Limit)
		fmt.Printf("Count: %d\n", resp.Count)
		fmt.Printf("NextCursor: %s\n", resp.NextCursor)
		fmt.Printf("Data length: %d\n", len(resp.Data))

		// 解析并显示市场数据
		for i, rawMarket := range resp.Data {
			var market map[string]interface{}
			if err := json.Unmarshal(rawMarket, &market); err != nil {
				fmt.Printf("  [%d] 解析失败: %v\n", i, err)
				continue
			}

			// 显示关键信息
			fmt.Printf("\n  [%d] 市场信息:\n", i+1)
			if conditionID, ok := market["condition_id"].(string); ok {
				fmt.Printf("      Condition ID: %s\n", conditionID)
			}
			if question, ok := market["question"].(string); ok {
				fmt.Printf("      Question: %s\n", question)
			}
			if ticker, ok := market["ticker"].(string); ok {
				fmt.Printf("      Ticker: %s\n", ticker)
			}
			if description, ok := market["description"].(string); ok && len(description) > 0 {
				fmt.Printf("      Description: %s\n", description)
			}
			if active, ok := market["active"].(bool); ok {
				fmt.Printf("      Active: %v\n", active)
			}
			if closed, ok := market["closed"].(bool); ok {
				fmt.Printf("      Closed: %v\n", closed)
			}
			if orders, ok := market["orders"].(string); ok {
				fmt.Printf("      Orders: %s\n", orders)
			}

			// 只显示前 3 个市场的详细信息
			if i >= 2 {
				fmt.Printf("\n  ... (还有 %d 个市场)\n", len(resp.Data)-3)
				break
			}
		}

		// 检查是否还有下一页
		if resp.NextCursor == "" || resp.NextCursor == pm.EndCursor || resp.NextCursor == nextCursor {
			fmt.Println("\n没有更多数据了")
			break
		}

		nextCursor = resp.NextCursor
		pageNum++
	}

	// 演示如何获取单个市场（通过 condition_id）
	fmt.Println("\n========== 获取单个市场示例 ==========")
	if len(os.Args) > 1 {
		conditionID := os.Args[1]
		rawMarket, err := sdk.CLOB.GetMarket(context.Background(), conditionID)
		if err != nil {
			fmt.Printf("get market failed: %v\n", err)
			return
		}

		var market map[string]interface{}
		if err := json.Unmarshal(rawMarket, &market); err != nil {
			fmt.Printf("parse market failed: %v\n", err)
			return
		}

		// 格式化输出
		prettyJSON, _ := json.MarshalIndent(market, "", "  ")
		fmt.Printf("\n市场详情:\n%s\n", string(prettyJSON))
	} else {
		fmt.Println("\n提示: 可以传入 condition_id 参数获取单个市场详情")
		fmt.Println("例如: go run main.go 0x123456...")
	}
}
