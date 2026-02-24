// types_bridge.go 模块
package polymarket

// =====================================================
// Bridge API 类型定义
// 对应文档：https://docs.polymarket.com/api-reference/bridge
// =====================================================

// ----- 获取支持的资产 (GET /supported-assets) -----

// BridgeToken 代币基本信息。
type BridgeToken struct {
	// 代币名称，例如 "USD Coin"
	Name string `json:"name"`
	// 代币符号，例如 "USDC"
	Symbol string `json:"symbol"`
	// 代币合约地址
	Address string `json:"address"`
	// 代币精度
	Decimals int `json:"decimals"`
}

// SupportedAsset 支持的单个资产，包含链和代币信息。
type SupportedAsset struct {
	// 链 ID，例如 "1"（Ethereum）、"137"（Polygon）
	ChainID string `json:"chainId"`
	// 链名称，例如 "Ethereum"
	ChainName string `json:"chainName"`
	// 代币详细信息
	Token BridgeToken `json:"token"`
	// 最小结算金额（USD）
	MinCheckoutUSD float64 `json:"minCheckoutUsd"`
}

// SupportedAssetsResponse 获取支持资产的响应。
type SupportedAssetsResponse struct {
	// 支持的资产列表，包含各链上的代币及其最低充值/提现额度
	SupportedAssets []SupportedAsset `json:"supportedAssets"`
}

// ----- 创建充值地址 (POST /deposit) -----

// DepositRequest 创建充值地址的请求参数。
type DepositRequest struct {
	// Polymarket 钱包地址（Polygon 上），充值的资金将以 USDC.e 形式入账到该地址
	Address string `json:"address"`
}

// BridgeAddresses 多链充值/提现地址。
type BridgeAddresses struct {
	// EVM 链充值/提现地址（以太坊、Polygon、Base 等）
	EVM string `json:"evm"`
	// Solana 链充值/提现地址
	SVM string `json:"svm"`
	// 比特币充值/提现地址
	BTC string `json:"btc"`
}

// DepositResponse 创建充值地址的响应。
type DepositResponse struct {
	// 各区块链网络的充值地址
	Address BridgeAddresses `json:"address"`
	// 附加说明信息
	Note string `json:"note"`
}

// ----- 获取报价 (POST /quote) -----

// QuoteRequest 获取跨链/跨币报价的请求参数。
type QuoteRequest struct {
	// 发送代币数量（最小单位），例如 "10000000"（即 10 USDC，精度为 6）
	FromAmountBaseUnit string `json:"fromAmountBaseUnit"`
	// 来源链 ID，例如 "137"（Polygon）
	FromChainID string `json:"fromChainId"`
	// 来源代币合约地址
	FromTokenAddress string `json:"fromTokenAddress"`
	// 接收方钱包地址
	RecipientAddress string `json:"recipientAddress"`
	// 目标链 ID，例如 "137"（Polygon）
	ToChainID string `json:"toChainId"`
	// 目标代币合约地址
	ToTokenAddress string `json:"toTokenAddress"`
}

// EstFeeBreakdown 预估费用明细。
type EstFeeBreakdown struct {
	// 应用手续费标签
	AppFeeLabel string `json:"appFeeLabel"`
	// 应用手续费百分比
	AppFeePercent float64 `json:"appFeePercent"`
	// 应用手续费（USD）
	AppFeeUSD float64 `json:"appFeeUsd"`
	// 填充成本百分比
	FillCostPercent float64 `json:"fillCostPercent"`
	// 填充成本（USD）
	FillCostUSD float64 `json:"fillCostUsd"`
	// Gas 费用（USD）
	GasUSD float64 `json:"gasUsd"`
	// 最大滑点
	MaxSlippage float64 `json:"maxSlippage"`
	// 最低接收数量
	MinReceived float64 `json:"minReceived"`
	// 兑换影响
	SwapImpact float64 `json:"swapImpact"`
	// 兑换影响（USD）
	SwapImpactUSD float64 `json:"swapImpactUsd"`
	// 总影响
	TotalImpact float64 `json:"totalImpact"`
	// 总影响（USD）
	TotalImpactUSD float64 `json:"totalImpactUsd"`
}

// QuoteResponse 获取报价的响应。
type QuoteResponse struct {
	// 预估完成结算所需时间（毫秒）
	EstCheckoutTimeMs int64 `json:"estCheckoutTimeMs"`
	// 预估费用明细
	EstFeeBreakdown EstFeeBreakdown `json:"estFeeBreakdown"`
	// 预估输入金额（USD）
	EstInputUSD float64 `json:"estInputUsd"`
	// 预估输出金额（USD）
	EstOutputUSD float64 `json:"estOutputUsd"`
	// 预估收到的代币数量（最小单位）
	EstToTokenBaseUnit string `json:"estToTokenBaseUnit"`
	// 报价的唯一标识
	QuoteID string `json:"quoteId"`
}

// ----- 获取交易状态 (GET /status/{address}) -----

// BridgeTransaction 单笔跨链桥交易信息。
type BridgeTransaction struct {
	// 来源链 ID
	FromChainID string `json:"fromChainId"`
	// 来源代币合约地址
	FromTokenAddress string `json:"fromTokenAddress"`
	// 发送金额（最小单位）
	FromAmountBaseUnit string `json:"fromAmountBaseUnit"`
	// 目标链 ID
	ToChainID string `json:"toChainId"`
	// 目标代币合约地址
	ToTokenAddress string `json:"toTokenAddress"`
	// 交易哈希（完成后才有值）
	TxHash string `json:"txHash,omitempty"`
	// 交易创建时间（毫秒级时间戳）
	CreatedTimeMs int64 `json:"createdTimeMs,omitempty"`
	// 交易状态：DEPOSIT_DETECTED / PROCESSING / COMPLETED
	Status string `json:"status"`
}

// TransactionStatusResponse 交易状态查询的响应。
type TransactionStatusResponse struct {
	// 交易列表
	Transactions []BridgeTransaction `json:"transactions"`
}

// ----- 创建提现地址 (POST /withdraw) -----

// WithdrawRequest 创建提现地址的请求参数。
type WithdrawRequest struct {
	// Polymarket 钱包地址（Polygon 上的来源地址）
	Address string `json:"address"`
	// 目标链 ID，例如 "1"（Ethereum）、"8453"（Base）、"1151111081099710"（Solana）
	ToChainID string `json:"toChainId"`
	// 目标代币合约地址
	ToTokenAddress string `json:"toTokenAddress"`
	// 目标接收钱包地址
	RecipientAddr string `json:"recipientAddr"`
}

// WithdrawResponse 创建提现地址的响应。
type WithdrawResponse struct {
	// 各区块链网络的提现地址
	Address BridgeAddresses `json:"address"`
	// 附加说明信息
	Note string `json:"note"`
}
