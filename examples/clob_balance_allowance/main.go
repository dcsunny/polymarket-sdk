package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/joho/godotenv"
)

// 获取账户余额与授权（CLOB /balance-allowance）
// - COLLATERAL：USDC 余额/授权
// - CONDITIONAL：指定 token_id 的持仓/授权（如果提供 POLYMARKET_TOKEN_ID）
func main() {
	_ = godotenv.Load()

	cfg := pm.Config{
		Address:       os.Getenv("POLYMARKET_ADDRESS"),
		APIKey:        os.Getenv("POLYMARKET_API_KEY"),
		APISecret:     os.Getenv("POLYMARKET_API_SECRET"),
		Passphrase:    os.Getenv("POLYMARKET_PASSPHRASE"),
		Proxy:         os.Getenv("PROXY"),
		SignatureType: envInt("POLYMARKET_SIG_TYPE", pm.SignatureTypeEOA),
		ChainID:       int64(envInt("POLYMARKET_CHAIN_ID", pm.DefaultChainID)),
	}

	sdk, err := pm.New(cfg)
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		return
	}

	ctx := context.Background()

	// 1) USDC（Collateral）余额/授权
	coll, err := sdk.CLOB.GetBalanceAllowance(ctx, &pm.BalanceAllowanceParams{
		AssetType: pm.AssetTypeCollateral,
	})
	if err != nil {
		fmt.Printf("get collateral balance failed: %v\n", err)
		return
	}
	fmt.Printf("collateral balance=%s (raw=%s)\n", formatUnits(coll.Balance, 6), coll.Balance)
	fmt.Printf("collateral allowance=%s (raw=%s)\n", formatUnits(coll.Allowance, 6), coll.Allowance)

	// 2) Conditional Token（可选）
	tokenID := os.Getenv("POLYMARKET_TOKEN_ID")
	if tokenID == "" {
		return
	}

	pos, err := sdk.CLOB.GetBalanceAllowance(ctx, &pm.BalanceAllowanceParams{
		AssetType: pm.AssetTypeConditional,
		TokenID:   tokenID,
	})
	if err != nil {
		fmt.Printf("get conditional balance failed: %v\n", err)
		return
	}
	fmt.Printf("conditional token_id=%s balance=%s (raw=%s)\n", tokenID, formatUnits(pos.Balance, 6), pos.Balance)
	fmt.Printf("conditional token_id=%s allowance=%s (raw=%s)\n", tokenID, formatUnits(pos.Allowance, 6), pos.Allowance)
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return def
}

// formatUnits 将整数（base unit）字符串按 decimals 格式化为可读小数。
func formatUnits(s string, decimals int) string {
	s = strings.TrimSpace(s)
	if s == "" || decimals <= 0 {
		return s
	}

	n := new(big.Int)
	if _, ok := n.SetString(s, 10); !ok {
		return s
	}

	base := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	intPart := new(big.Int).Quo(n, base)
	fracPart := new(big.Int).Mod(n, base)

	frac := fmt.Sprintf("%0*s", decimals, fracPart.String())
	frac = strings.TrimRight(frac, "0")
	if frac == "" {
		return intPart.String()
	}
	return intPart.String() + "." + frac
}
