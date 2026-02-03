// constants.go 模块
package polymarket

// Contract addresses and related constants.
const (
	CTFContractAddress        = "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045"
	UMACTFAdapterAddress      = "0x6A9D222616C90FcA5754cd1333cFD9b7fb6a4F74"
	USDCAddress               = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"
	CTFExchangeAddress        = "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E"
	NegRiskCTFExchangeAddress = "0xC5d563A36AE78145C45a50134d48A1215220f80a"
)

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

// Side constants.
const (
	SideBuy  = "BUY"
	SideSell = "SELL"
)
