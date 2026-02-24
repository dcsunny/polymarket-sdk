// bridge.go 模块
package polymarket

import (
	"context"
	"net/http"

	"github.com/dcsunny/polymarket-sdk/internal/httpx"
)

// BridgeClient 处理 Polymarket Bridge（跨链桥）API。
// 对应文档：https://docs.polymarket.com/api-reference/bridge
//
// Bridge API 提供跨链资产转移功能，支持在不同区块链网络之间充值和提现。
// Base URL: https://bridge.polymarket.com
type BridgeClient struct {
	http *httpx.Client
}

// NewBridgeClient 创建新的 Bridge 客户端。
func NewBridgeClient(http *httpx.Client) *BridgeClient {
	return &BridgeClient{http: http}
}

// GetSupportedAssets 获取支持的资产列表。
//
// 返回所有支持充值和提现的链及代币信息，包括各链上的最低结算金额。
//
// API: GET /supported-assets
// 文档: https://docs.polymarket.com/api-reference/bridge/get-supported-assets
func (c *BridgeClient) GetSupportedAssets(ctx context.Context) (*SupportedAssetsResponse, error) {
	var resp SupportedAssetsResponse
	if err := c.http.Do(ctx, http.MethodGet, EndpointBridgeSupportedAssets, nil, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateDepositAddresses 创建充值地址。
//
// 根据指定的 Polymarket 钱包地址，生成多个区块链网络（EVM、Solana、Bitcoin）的充值地址。
// 用户可以向这些地址发送资金，资金将以 USDC.e 形式入账到指定的 Polymarket 钱包。
//
// API: POST /deposit
// 文档: https://docs.polymarket.com/api-reference/bridge/create-deposit-addresses
func (c *BridgeClient) CreateDepositAddresses(ctx context.Context, req DepositRequest) (*DepositResponse, error) {
	if req.Address == "" {
		return nil, ErrInvalidArgument("address is required")
	}

	var resp DepositResponse
	if err := c.http.Do(ctx, http.MethodPost, EndpointBridgeDeposit, nil, req, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetQuote 获取跨链/跨币兑换报价。
//
// 根据来源链、代币和金额，获取到目标链和代币的兑换报价，
// 包括预估到账金额、费用明细、预计完成时间等信息。
//
// API: POST /quote
// 文档: https://docs.polymarket.com/api-reference/bridge/get-a-quote
func (c *BridgeClient) GetQuote(ctx context.Context, req QuoteRequest) (*QuoteResponse, error) {
	if req.FromAmountBaseUnit == "" {
		return nil, ErrInvalidArgument("fromAmountBaseUnit is required")
	}
	if req.FromChainID == "" {
		return nil, ErrInvalidArgument("fromChainId is required")
	}
	if req.FromTokenAddress == "" {
		return nil, ErrInvalidArgument("fromTokenAddress is required")
	}
	if req.RecipientAddress == "" {
		return nil, ErrInvalidArgument("recipientAddress is required")
	}
	if req.ToChainID == "" {
		return nil, ErrInvalidArgument("toChainId is required")
	}
	if req.ToTokenAddress == "" {
		return nil, ErrInvalidArgument("toTokenAddress is required")
	}

	var resp QuoteResponse
	if err := c.http.Do(ctx, http.MethodPost, EndpointBridgeQuote, nil, req, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTransactionStatus 获取跨链桥交易状态。
//
// 根据充值/提现地址（从 /deposit 或 /withdraw 返回的 EVM、SVM 或 BTC 地址）
// 查询该地址关联的所有交易及其状态。
// 交易状态可能为：DEPOSIT_DETECTED（已检测到充值）、PROCESSING（处理中）、COMPLETED（已完成）。
//
// API: GET /status/{address}
// 文档: https://docs.polymarket.com/api-reference/bridge/get-transaction-status
func (c *BridgeClient) GetTransactionStatus(ctx context.Context, address string) (*TransactionStatusResponse, error) {
	if address == "" {
		return nil, ErrInvalidArgument("address is required")
	}

	var resp TransactionStatusResponse
	path := EndpointBridgeStatus + address
	if err := c.http.Do(ctx, http.MethodGet, path, nil, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateWithdrawalAddresses 创建提现地址。
//
// 从 Polymarket 钱包提现到指定的目标链和代币。
// 生成多个区块链网络的提现地址，用户需要通过 Polymarket 钱包向这些地址发送资金来完成提现。
//
// API: POST /withdraw
// 文档: https://docs.polymarket.com/api-reference/bridge/create-withdrawal-addresses
func (c *BridgeClient) CreateWithdrawalAddresses(ctx context.Context, req WithdrawRequest) (*WithdrawResponse, error) {
	if req.Address == "" {
		return nil, ErrInvalidArgument("address is required")
	}
	if req.ToChainID == "" {
		return nil, ErrInvalidArgument("toChainId is required")
	}
	if req.ToTokenAddress == "" {
		return nil, ErrInvalidArgument("toTokenAddress is required")
	}
	if req.RecipientAddr == "" {
		return nil, ErrInvalidArgument("recipientAddr is required")
	}

	var resp WithdrawResponse
	if err := c.http.Do(ctx, http.MethodPost, EndpointBridgeWithdraw, nil, req, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
