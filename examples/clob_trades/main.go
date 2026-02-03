package main

import (
	"context"
	"fmt"
	"os"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/joho/godotenv"
)

//查询某个market的成交记录，“已经成交的记录”

func main() {
	_ = godotenv.Load()

	marketID := os.Getenv("POLYMARKET_MARKET_ID")
	if marketID == "" {
		fmt.Println("set POLYMARKET_MARKET_ID")
		return
	}

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

	trades, err := sdk.CLOB.GetTradesByMarket(context.Background(), marketID)
	if err != nil {
		fmt.Printf("get trades failed: %v\n", err)
		return
	}

	fmt.Printf("trades: %d\n", len(trades))
	if len(trades) > 0 {
		fmt.Printf("latest trade: id=%s price=%s size=%s\n", trades[0].ID, trades[0].Price, trades[0].Size)
	}
}
