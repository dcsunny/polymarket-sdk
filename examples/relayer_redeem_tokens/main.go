package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

// 赎回代币（Redeem Positions）
//
// 说明：
// - 该示例使用 Polymarket Relayer（`RelayerClient`）提交 redeem 交易。
// - redeem 通常需要 Safe 已部署；若未部署，可设置 AUTO_DEPLOY_SAFE=true 自动部署（需要 Builder API 认证）。
//
// 环境变量（必填）：
// - POLYMARKET_RPC_URL
// - POLYMARKET_PRIVATE_KEY
// - POLYMARKET_CONDITION_ID（或 POLYMARKET_MARKET_ID）
//
// 可选：
// - POLYMARKET_RELAYER_URL（默认 https://relayer-v2.polymarket.com）
// - POLYMARKET_CHAIN_ID（默认 137）
// - AUTO_DEPLOY_SAFE（true/false，默认 false）
//
// Builder Auth（Safe 部署/某些 relayer 操作需要）：
// - POLYMARKET_BUILDER_API_KEY
// - POLYMARKET_BUILDER_API_SECRET
// - POLYMARKET_BUILDER_PASSPHRASE
//
// 赎回模式：
// - IS_NEGRISK=true：NegRisk 赎回，使用 YES_AMOUNT/NO_AMOUNT（整数，base unit）
// - IS_NEGRISK=false：CTF 赎回，使用 INDEX_SETS（例如 "1,2"），以及 COLLATERAL_TOKEN（默认 Polygon USDC）
func main() {
	_ = godotenv.Load()

	rpcURL := os.Getenv("POLYMARKET_RPC_URL")
	if rpcURL == "" {
		fmt.Println("missing POLYMARKET_RPC_URL")
		return
	}
	privateKey := os.Getenv("POLYMARKET_PRIVATE_KEY")
	if privateKey == "" {
		fmt.Println("missing POLYMARKET_PRIVATE_KEY")
		return
	}

	conditionID := envString("POLYMARKET_CONDITION_ID", "")
	if conditionID == "" {
		conditionID = envString("POLYMARKET_MARKET_ID", "")
	}
	if conditionID == "" {
		fmt.Println("missing POLYMARKET_CONDITION_ID (or POLYMARKET_MARKET_ID)")
		return
	}

	relayerURL := strings.TrimRight(envString("POLYMARKET_RELAYER_URL", "https://relayer-v2.polymarket.com"), "/")
	chainID := int64(envInt("POLYMARKET_CHAIN_ID", 137))

	// Builder auth：可选（但 AUTO_DEPLOY_SAFE=true 时基本必需）
	var builderAuth *pm.BuilderAuth
	if k := os.Getenv("POLYMARKET_BUILDER_API_KEY"); k != "" {
		builderAuth = pm.NewBuilderAuth(
			k,
			os.Getenv("POLYMARKET_BUILDER_API_SECRET"),
			os.Getenv("POLYMARKET_BUILDER_PASSPHRASE"),
		)
	}

	ctx := context.Background()
	client, err := pm.NewRelayerClient(ctx, pm.RelayerConfig{
		RelayerURL:  relayerURL,
		RPCURL:      rpcURL,
		PrivateKey:  strings.TrimPrefix(privateKey, "0x"),
		ChainID:     chainID,
		BuilderAuth: builderAuth,
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
	if !deployed {
		if envBool("AUTO_DEPLOY_SAFE", false) {
			fmt.Println("Safe not deployed, deploying...")
			if err := client.DeploySafe(ctx); err != nil {
				fmt.Printf("deploy safe failed: %v\n", err)
				return
			}
			fmt.Println("Safe deployed.")
		} else {
			fmt.Println("Safe not deployed. Set AUTO_DEPLOY_SAFE=true to deploy automatically (requires builder auth).")
			return
		}
	}

	isNegRisk := envBool("IS_NEGRISK", true)

	// 1) NegRisk redeem
	if isNegRisk {
		yesAmount := mustBigInt(envString("YES_AMOUNT", "0"))
		noAmount := mustBigInt(envString("NO_AMOUNT", "0"))

		resp, err := client.RedeemNegRiskPositions(ctx, conditionID, yesAmount, noAmount)
		if err != nil {
			fmt.Printf("redeem negrisk failed: %v\n", err)
			return
		}
		printRelayerResp(resp)
		return
	}

	// 2) CTF redeem
	indexSetsStr := envString("INDEX_SETS", "1,2")
	indexSets, err := parseBigIntList(indexSetsStr)
	if err != nil || len(indexSets) == 0 {
		fmt.Printf("invalid INDEX_SETS: %q\n", indexSetsStr)
		return
	}
	collateral := common.HexToAddress(envString("COLLATERAL_TOKEN", pm.USDCAddress))

	resp, err := client.RedeemCTFPositions(ctx, conditionID, indexSets, collateral)
	if err != nil {
		fmt.Printf("redeem ctf failed: %v\n", err)
		return
	}
	printRelayerResp(resp)
}

func printRelayerResp(resp *pm.RelayerResponse) {
	fmt.Printf("transaction id: %s\n", resp.TransactionID)
	fmt.Printf("state: %s\n", resp.State)
	fmt.Printf("tx hash: %s\n", resp.TransactionHash)
	if resp.TransactionHash != "" {
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

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		v = strings.ToLower(strings.TrimSpace(v))
		return v == "1" || v == "true" || v == "yes" || v == "y"
	}
	return def
}

func mustBigInt(s string) *big.Int {
	n := new(big.Int)
	s = strings.TrimSpace(s)
	if s == "" {
		return big.NewInt(0)
	}
	if _, ok := n.SetString(s, 10); !ok {
		panic(fmt.Sprintf("invalid int: %q", s))
	}
	return n
}

func parseBigIntList(csv string) ([]*big.Int, error) {
	parts := strings.Split(csv, ",")
	out := make([]*big.Int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n := new(big.Int)
		if _, ok := n.SetString(p, 10); !ok {
			return nil, fmt.Errorf("invalid int: %q", p)
		}
		out = append(out, n)
	}
	return out, nil
}
