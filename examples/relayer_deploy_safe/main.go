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

// 部署 Safe（relayer /deploy-safe）
//
// 注意：该接口通常需要 Builder API 认证（POLYMARKET_BUILDER_*）。
//
// 环境变量（必填）：
// - POLYMARKET_RPC_URL
// - POLYMARKET_PRIVATE_KEY
//
// Builder Auth（通常必填）：
// - POLYMARKET_BUILDER_API_KEY
// - POLYMARKET_BUILDER_API_SECRET
// - POLYMARKET_BUILDER_PASSPHRASE
//
// 可选：
// - POLYMARKET_RELAYER_URL（默认 https://relayer-v2.polymarket.com）
// - POLYMARKET_CHAIN_ID（默认 137）
func main() {
	_ = godotenv.Load()

	rpcURL := os.Getenv("POLYMARKET_RPC_URL")
	privateKey := os.Getenv("POLYMARKET_PRIVATE_KEY")
	if rpcURL == "" || privateKey == "" {
		fmt.Println("missing POLYMARKET_RPC_URL or POLYMARKET_PRIVATE_KEY")
		return
	}

	relayerURL := strings.TrimRight(envString("POLYMARKET_RELAYER_URL", "https://relayer-v2.polymarket.com"), "/")
	chainID := int64(envInt("POLYMARKET_CHAIN_ID", 137))

	// Builder auth（部署 Safe 通常需要）
	builderKey := os.Getenv("POLYMARKET_BUILDER_API_KEY")
	builderSecret := os.Getenv("POLYMARKET_BUILDER_API_SECRET")
	builderPassphrase := os.Getenv("POLYMARKET_BUILDER_PASSPHRASE")
	if builderKey == "" || builderSecret == "" || builderPassphrase == "" {
		fmt.Println("missing builder auth env vars (POLYMARKET_BUILDER_API_KEY/SECRET/PASSPHRASE)")
		return
	}

	ctx := context.Background()
	client, err := pm.NewRelayerClient(ctx, pm.RelayerConfig{
		RelayerURL: relayerURL,
		RPCURL:     rpcURL,
		PrivateKey: privateKey, // 支持带/不带 0x
		ChainID:    chainID,
		BuilderAuth: pm.NewBuilderAuth(
			builderKey,
			builderSecret,
			builderPassphrase,
		),
	})
	if err != nil {
		fmt.Printf("init relayer client failed: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Printf("EOA address: %s\n", client.GetAddress().Hex())
	fmt.Printf("Safe address: %s\n", client.GetSafeAddress().Hex())

	deployed, err := client.IsSafeDeployed(ctx)
	if err != nil {
		fmt.Printf("check safe deployed failed: %v\n", err)
		return
	}
	if deployed {
		fmt.Println("Safe already deployed.")
		return
	}

	resp, err := client.DeploySafeSubmit(ctx)
	if err != nil {
		fmt.Printf("deploy safe failed: %v\n", err)
		if resp != nil && resp.Message != "" {
			fmt.Printf("relayer message: %s\n", resp.Message)
		}
		return
	}

	if resp != nil && resp.State == "already_deployed" {
		fmt.Println("Safe already deployed.")
		return
	}

	fmt.Printf("state: %s\n", resp.State)
	if resp.TransactionID != "" {
		fmt.Printf("transaction id: %s\n", resp.TransactionID)
	}
	if resp.TransactionHash != "" {
		fmt.Printf("tx hash: %s\n", resp.TransactionHash)
		fmt.Printf("polygonscan: https://polygonscan.com/tx/%s\n", resp.TransactionHash)
	}
	if resp.Message != "" {
		fmt.Printf("message: %s\n", resp.Message)
	}
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
