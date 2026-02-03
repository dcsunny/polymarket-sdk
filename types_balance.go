// types_balance.go 模块
package polymarket

// AssetType 资产类型。
type AssetType string

const (
	AssetTypeCollateral  AssetType = "COLLATERAL"
	AssetTypeConditional AssetType = "CONDITIONAL"
)

// BalanceAllowanceParams 余额与授权查询参数。
type BalanceAllowanceParams struct {
	AssetType AssetType
	TokenID   string
}

// BalanceAllowanceResponse 余额与授权响应。
type BalanceAllowanceResponse struct {
	Balance   string `json:"balance"`
	Allowance string `json:"allowance"`
}
