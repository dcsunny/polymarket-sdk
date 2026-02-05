// constants.go 模块
package polymarket

import "fmt"

// 合约地址相关常量（Polygon 主网）。
const (
	CTFContractAddress        = "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045"
	UMACTFAdapterAddress      = "0x6A9D222616C90FcA5754cd1333cFD9b7fb6a4F74"
	USDCAddress               = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"
	CTFExchangeAddress        = "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E"
	NegRiskCTFExchangeAddress = "0xC5d563A36AE78145C45a50134d48A1215220f80a"
	NegRiskAdapterAddress     = "0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296"
)

// 合约地址相关常量（Polygon Amoy 测试网）。
const (
	AmoyCTFContractAddress     = "0x69308FB512518e39F9b16112fA8d994F4e2Bf8bB"
	AmoyCollateralAddress      = "0x9c4e1703476e875070ee25b56a58b008cfb8fa78"
	AmoyCTFExchangeAddress     = "0xdFE02Eb6733538f8Ea35D585af8DE5958AD99E40"
	AmoyNegRiskExchangeAddress = NegRiskCTFExchangeAddress
	AmoyNegRiskAdapterAddress  = NegRiskAdapterAddress
)

// 链 ID 常量（与 Node SDK 对齐）。
const (
	ChainIDPolygon int64 = 137
	ChainIDAmoy    int64 = 80002
)

// 代币精度常量（与 Node SDK 对齐）。
const (
	CollateralTokenDecimals  = 6
	ConditionalTokenDecimals = 6
)

// CredsCreationWarning 创建 API 凭证时的提示文案（对齐 Node SDK）。
const CredsCreationWarning = `Your credentials CANNOT be recovered after they've been created.
Be sure to store them safely!`

// ContractConfig 合约配置（对齐 Node SDK 的 ContractConfig）。
type ContractConfig struct {
	Exchange          string
	NegRiskAdapter    string
	NegRiskExchange   string
	Collateral        string
	ConditionalTokens string
}

var (
	// PolygonContractConfig Polygon 主网合约配置。
	PolygonContractConfig = ContractConfig{
		Exchange:          CTFExchangeAddress,
		NegRiskAdapter:    NegRiskAdapterAddress,
		NegRiskExchange:   NegRiskCTFExchangeAddress,
		Collateral:        USDCAddress,
		ConditionalTokens: CTFContractAddress,
	}
	// AmoyContractConfig Polygon Amoy 测试网合约配置。
	AmoyContractConfig = ContractConfig{
		Exchange:          AmoyCTFExchangeAddress,
		NegRiskAdapter:    AmoyNegRiskAdapterAddress,
		NegRiskExchange:   AmoyNegRiskExchangeAddress,
		Collateral:        AmoyCollateralAddress,
		ConditionalTokens: AmoyCTFContractAddress,
	}
)

// GetContractConfig 根据链 ID 获取合约配置。
func GetContractConfig(chainID int64) (ContractConfig, error) {
	switch chainID {
	case ChainIDPolygon:
		return PolygonContractConfig, nil
	case ChainIDAmoy:
		return AmoyContractConfig, nil
	default:
		return ContractConfig{}, fmt.Errorf("unsupported chain id: %d", chainID)
	}
}

// UMAAdapterABI is a simplified ABI for resolution status.
const UMAAdapterABI = `[
	{
		"inputs": [
			{"internalType": "bytes32", "name": "questionId", "type": "bytes32"}
		],
		"name": "ready",
		"outputs": [{"internalType": "bool", "name": "", "type": "bool"}],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "bytes32", "name": "questionId", "type": "bytes32"}
		],
		"name": "payouts",
		"outputs": [{"internalType": "uint256[]", "name": "", "type": "uint256[]"}],
		"stateMutability": "view",
		"type": "function"
	}
]`

// Side 常量（与 Node SDK 对齐）。
const (
	SideBuy  = "BUY"
	SideSell = "SELL"
)
