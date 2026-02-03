package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/joho/godotenv"
)

// 下单
func main() {
	_ = godotenv.Load()

	cfg := pm.Config{
		PrivateKey:    os.Getenv("POLYMARKET_PRIVATE_KEY"),
		Address:       os.Getenv("POLYMARKET_ADDRESS"),
		APIKey:        os.Getenv("POLYMARKET_API_KEY"),
		APISecret:     os.Getenv("POLYMARKET_API_SECRET"),
		Passphrase:    os.Getenv("POLYMARKET_PASSPHRASE"),
		Proxy:         os.Getenv("PROXY"),
		Funder:        os.Getenv("POLYMARKET_FUNDER"),
		SignatureType: envInt("POLYMARKET_SIG_TYPE", pm.SignatureTypePolyGnosisSafe),
		ChainID:       int64(envInt("POLYMARKET_CHAIN_ID", pm.DefaultChainID)),
	}

	sdk, err := pm.New(cfg)
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		return
	}

	tokenID := os.Getenv("POLYMARKET_TOKEN_ID")
	side := strings.ToUpper(envString("POLYMARKET_SIDE", "BUY"))
	if tokenID == "" {
		fmt.Println("missing POLYMARKET_TOKEN_ID")
		return
	}

	orderArgs := &pm.OrderArgs{
		TokenID:     tokenID,
		MakerAmount: "5500000",
		TakerAmount: "10000000",
		Side:        side,
		Taker:       envString("POLYMARKET_TAKER", "0x0000000000000000000000000000000000000000"),
		FeeRateBps:  envString("POLYMARKET_FEE_RATE_BPS", "0"),
		Nonce:       envString("POLYMARKET_NONCE", "0"),
		Expiration:  envString("POLYMARKET_EXPIRATION", "0"),
	}

	order, err := sdk.CLOB.CreateOrder(orderArgs)
	if err != nil {
		fmt.Printf("create order failed: %v\n", err)
		return
	}

	resp, err := sdk.CLOB.PostOrder(context.Background(), order, pm.OrderTypeGTC)
	if err != nil {
		fmt.Printf("post order failed: %v\n", err)
		return
	}

	fmt.Printf("order response: success=%v orderId=%s status=%s\n", resp.Success, resp.OrderID, resp.Status)
}

func envString(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return def
}
