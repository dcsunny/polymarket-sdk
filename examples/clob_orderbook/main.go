package main

import (
	"context"
	"fmt"
	"os"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/joho/godotenv"
)

// 获取单个 token 的订单簿快照（买单/卖单挂单情况）
//
// 环境变量：
// - POLYMARKET_TOKEN_ID：clob 的 token_id（asset_id）
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

	tokenID := os.Getenv("POLYMARKET_TOKEN_ID")
	if tokenID == "" {
		fmt.Println("POLYMARKET_TOKEN_ID is required")
		return
	}

	book, err := sdk.CLOB.GetOrderBook(context.Background(), tokenID)
	if err != nil {
		fmt.Printf("get orderbook failed: %v\n", err)
		return
	}

	fmt.Printf("Market: %s\n", book.Market)
	fmt.Printf("Asset ID: %s\n", book.AssetID)
	fmt.Printf("Last Trade Price: %s\n", book.LastTradePrice)
	fmt.Printf("Min Order Size: %s\n", book.MinOrderSize)
	fmt.Printf("Tick Size: %s\n", book.TickSize)
	fmt.Printf("Timestamp: %s\n", book.Timestamp)
	fmt.Printf("Neg Risk: %v\n", book.NegRisk)
	fmt.Println()

	// 打印买单（价格从高到低）
	fmt.Printf("Bids (Buy Orders): %d\n", len(book.Bids))
	for i, bid := range book.Bids {
		fmt.Printf("  Bid #%d: Price=%s, Size=%s\n", i+1, bid.Price, bid.Size)
	}
	fmt.Println()

	// 打印卖单（价格从低到高）
	fmt.Printf("Asks (Sell Orders): %d\n", len(book.Asks))
	for i, ask := range book.Asks {
		fmt.Printf("  Ask #%d: Price=%s, Size=%s\n", i+1, ask.Price, ask.Size)
	}
}
