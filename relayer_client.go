// relayer_client.go 模块
package polymarket

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 已知的代理合约地址
var knownProxyContracts = map[common.Address]bool{
	common.HexToAddress(USDCAddress): true, // Polygon 主网 USDC
}

// 公共的 ERC20 ABI 定义
const erc20ABIString = `[
	{
		"constant": false,
		"inputs": [
			{"name": "spender", "type": "address"},
			{"name": "amount", "type": "uint256"}
		],
		"name": "approve",
		"outputs": [{"name": "", "type": "bool"}],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{"name": "owner", "type": "address"},
			{"name": "spender", "type": "address"}
		],
		"name": "allowance",
		"outputs": [{"name": "", "type": "uint256"}],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [{"name": "who", "type": "address"}],
		"name": "balanceOf",
		"outputs": [{"name": "", "type": "uint256"}],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`

// CTF Exchange redeemPositions ABI（CTF 交易所赎回位置 ABI）
const ctfRedeemABIString = `[
	{
		"constant": false,
		"inputs": [
			{"name": "collateralToken", "type": "address"},
			{"name": "parentCollectionId", "type": "bytes32"},
			{"name": "conditionId", "type": "bytes32"},
			{"name": "indexSets", "type": "uint256[]"}
		],
		"name": "redeemPositions",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

// NegRisk Adapter redeemPositions ABI（NegRisk 适配器赎回位置 ABI）
const negRiskRedeemABIString = `[
	{
		"inputs": [
			{"internalType": "bytes32", "name": "_conditionId", "type": "bytes32"},
			{"internalType": "uint256[]", "name": "_amounts", "type": "uint256[]"}
		],
		"name": "redeemPositions",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

var (
	// 解析后的 ABI
	parsedERC20ABI     abi.ABI
	parsedCTFRedeemABI abi.ABI
	parsedNegRiskABI   abi.ABI
)

func init() {
	var err error
	// 解析 ERC20 ABI
	parsedERC20ABI, err = abi.JSON(strings.NewReader(erc20ABIString))
	if err != nil {
		panic(fmt.Sprintf("解析 ERC20 ABI 失败: %v", err))
	}

	// 解析 CTF redeem ABI
	parsedCTFRedeemABI, err = abi.JSON(strings.NewReader(ctfRedeemABIString))
	if err != nil {
		panic(fmt.Sprintf("解析 CTF redeem ABI 失败: %v", err))
	}

	// 解析 NegRisk redeem ABI
	parsedNegRiskABI, err = abi.JSON(strings.NewReader(negRiskRedeemABIString))
	if err != nil {
		panic(fmt.Sprintf("解析 NegRisk redeem ABI 失败: %v", err))
	}
}

// IsProxyContract 检查是否为已知的代理合约
func (c *RelayerClient) IsProxyContract(_ context.Context, contractAddress common.Address) bool {
	return knownProxyContracts[contractAddress]
}

// RelayerConfig Relayer 配置
type RelayerConfig struct {
	RelayerURL  string       `json:"relayerUrl"`
	RPCURL      string       `json:"rpcUrl"`
	PrivateKey  string       `json:"privateKey"`
	ChainID     int64        `json:"chainId"`
	BuilderAuth *BuilderAuth `json:"builderAuth,omitempty"`
}

// RedeemRelayerRequest 赎回请求（重命名以避免冲突）
type RedeemRelayerRequest struct {
	ConditionID string `json:"conditionId"`
	// For CTF (IsNegRisk=false): Index sets to redeem
	IndexSets []*big.Int `json:"indexSets,omitempty"`
	// For NegRisk (IsNegRisk=true): Amounts to redeem [yesAmount, noAmount]
	RedeemAmounts []*big.Int `json:"redeemAmounts,omitempty"`
	// For CTF only: Collateral token address
	CollateralToken common.Address `json:"collateralToken,omitempty"`
	// For CTF only: Parent collection ID (use zero hash for default)
	ParentCollection [32]byte `json:"parentCollection,omitempty"`
	// Determines which contract to call
	IsNegRisk bool `json:"isNegRisk"`
	// Optional metadata
	Metadata string `json:"metadata,omitempty"`
}

// RelayerResponse Relayer 响应
type RelayerResponse struct {
	TransactionID   string `json:"transactionId"`
	TransactionHash string `json:"transactionHash"`
	State           string `json:"state"`
	Message         string `json:"message"`
}

// RelayerClient Relayer 客户端
type RelayerClient struct {
	config     RelayerConfig
	ethClient  *ethclient.Client
	privateKey *ecdsa.PrivateKey
	httpClient *http.Client

	// 缓存
	safeAddress common.Address
}

// NewRelayerClient 创建 Relayer 客户端
func NewRelayerClient(ctx context.Context, cfg RelayerConfig) (*RelayerClient, error) {
	// 连接到以太坊客户端
	ethClient, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("连接 RPC 失败: %w", err)
	}

	// 解析私钥
	privateKeyHex := strings.TrimPrefix(cfg.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}

	client := &RelayerClient{
		config:     cfg,
		ethClient:  ethClient,
		privateKey: privateKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// 派生 Safe 地址
	client.safeAddress = client.deriveSafeAddress()

	return client, nil
}

// GetAddress 获取 EOA 地址
func (c *RelayerClient) GetAddress() common.Address {
	return crypto.PubkeyToAddress(c.privateKey.PublicKey)
}

// GetSafeAddress 获取 Safe 地址
func (c *RelayerClient) GetSafeAddress() common.Address {
	return c.safeAddress
}

// IsSafeDeployed 检查 Safe 是否已部署
func (c *RelayerClient) IsSafeDeployed(ctx context.Context) (bool, error) {
	code, err := c.ethClient.CodeAt(ctx, c.safeAddress, nil)
	if err != nil {
		return false, fmt.Errorf("检查 Safe 代码失败: %w", err)
	}
	return len(code) > 0, nil
}

// DeploySafe 部署 Safe（如果尚未部署）
func (c *RelayerClient) DeploySafe(ctx context.Context) error {
	_, err := c.DeploySafeSubmit(ctx)
	return err
}

// DeploySafeSubmit 部署 Safe 并返回 relayer 响应（如果 Safe 已部署则返回 State=already_deployed）。
func (c *RelayerClient) DeploySafeSubmit(ctx context.Context) (*RelayerResponse, error) {
	deployed, err := c.IsSafeDeployed(ctx)
	if err != nil {
		return nil, err
	}
	if deployed {
		return &RelayerResponse{State: "already_deployed"}, nil
	}

	// 构建部署请求
	req := &safeCreateRequest{
		From:        c.GetAddress().Hex(),
		SafeAddress: c.safeAddress.Hex(),
		SaltNonce:   "0",
		Signature:   "",
		SignatureParams: safeSignatureParams{
			GasPrice:       "0",
			Operation:      "0",
			SafeTxnGas:     "0",
			BaseGas:        "0",
			GasToken:       "0x0000000000000000000000000000000000000000",
			RefundReceiver: "0x0000000000000000000000000000000000000000",
		},
	}

	// 发送部署请求
	resp, err := c.submit(ctx, "/deploy-safe", req)
	if err != nil {
		return nil, fmt.Errorf("部署 Safe 失败: %w", err)
	}

	if resp.State != "submitted" {
		return resp, fmt.Errorf("Safe 部署失败: %s", resp.Message)
	}

	return resp, nil
}

// RedeemPositions 赎回位置（统一接口）
func (c *RelayerClient) RedeemPositions(ctx context.Context, req *RedeemRelayerRequest) (*RelayerResponse, error) {
	// 确保 Safe 已部署
	if err := c.ensureSafeDeployed(ctx); err != nil {
		return nil, fmt.Errorf("确保 Safe 部署失败: %w", err)
	}

	// 构建 Safe 交易请求
	txReq, err := c.buildSafeSubmitRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("构建交易请求失败: %w", err)
	}

	// 提交交易
	resp, err := c.submit(ctx, "/submit", txReq)
	if err != nil {
		return nil, fmt.Errorf("提交交易失败: %w", err)
	}

	return resp, nil
}

// ApproveUSDC 授权 USDC 给指定地址（支持代理合约）
func (c *RelayerClient) ApproveUSDC(ctx context.Context, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	usdcAddress := common.HexToAddress(USDCAddress)

	// 检查是否为代理合约
	isProxy := c.IsProxyContract(ctx, usdcAddress)
	if isProxy {
		fmt.Printf("检测到 USDC 合约 %s: 代理合约\n", USDCAddress)
	}

	// 获取当前 nonce
	nonce, err := c.ethClient.PendingNonceAt(ctx, c.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("获取 nonce 失败: %w", err)
	}

	// 获取建议的 gas price
	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 gas 价格失败: %w", err)
	}

	// 构建 approve 调用数据
	data, err := parsedERC20ABI.Pack("approve", spender, amount)
	if err != nil {
		return nil, fmt.Errorf("构建 approve 调用数据失败: %w", err)
	}

	// 估算 gas
	msg := ethereum.CallMsg{
		From: c.GetAddress(),
		To:   &usdcAddress,
		Data: data,
	}
	gasLimit, err := c.ethClient.EstimateGas(ctx, msg)
	if err != nil {
		// 如果估算失败，使用默认值
		fmt.Printf("估算 gas 失败，使用默认值: %v\n", err)
		gasLimit = 100000 // 代理合约可能需要更多 gas
	} else {
		// 增加 20% 缓冲
		gasLimit = gasLimit * 120 / 100
		fmt.Printf("估算 Gas: %d (增加20%%缓冲后: %d)\n", gasLimit*100/120, gasLimit)
	}

	// 创建交易
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &usdcAddress, // 直接发送到合约地址（代理合约会自动转发）
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	// 签名交易
	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取链 ID 失败: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名交易失败: %w", err)
	}

	return signedTx, nil
}

// ApproveUSDCInfinite 授权无限 USDC (uint256 最大值)
func (c *RelayerClient) ApproveUSDCInfinite(ctx context.Context, spender common.Address) (*types.Transaction, error) {
	maxUint256 := CreateInfiniteApprovalAmount()
	return c.ApproveUSDC(ctx, spender, maxUint256)
}

// ApproveUSDCForAmount 授权指定 USDC 数量（以 USDC 为单位）
func (c *RelayerClient) ApproveUSDCForAmount(ctx context.Context, spender common.Address, amount float64) (*types.Transaction, error) {
	amountWei := CreateUSDCAmount(amount)
	return c.ApproveUSDC(ctx, spender, amountWei)
}

// SendApproveTransaction 发送授权交易
func (c *RelayerClient) SendApproveTransaction(ctx context.Context, tx *types.Transaction) error {
	// 发送交易
	err := c.ethClient.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("发送授权交易失败: %w", err)
	}

	fmt.Printf("授权交易已发送: %s\n", tx.Hash().Hex())
	return nil
}

// ApproveUSDCAndSend 授权 USDC 并发送交易（便捷方法）
func (c *RelayerClient) ApproveUSDCAndSend(ctx context.Context, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	// 构建授权交易
	tx, err := c.ApproveUSDC(ctx, spender, amount)
	if err != nil {
		return nil, err
	}

	// 发送交易
	err = c.SendApproveTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// ApproveUSDCInfiniteAndSend 授权无限 USDC 并发送交易
func (c *RelayerClient) ApproveUSDCInfiniteAndSend(ctx context.Context, spender common.Address) (*types.Transaction, error) {
	// uint256 最大值: 2^256 - 1
	maxUint256 := new(big.Int)
	maxUint256.Exp(big.NewInt(2), big.NewInt(256), nil)
	maxUint256.Sub(maxUint256, big.NewInt(1))

	return c.ApproveUSDCAndSend(ctx, spender, maxUint256)
}

// GetUSDCAllowance 获取 USDC 授权额度（支持代理合约）
func (c *RelayerClient) GetUSDCAllowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	usdcAddress := common.HexToAddress(USDCAddress)

	// 构建 allowance 调用数据
	data, err := parsedERC20ABI.Pack("allowance", owner, spender)
	if err != nil {
		return nil, fmt.Errorf("构建 allowance 调用数据失败: %w", err)
	}

	// 调用合约
	result, err := c.ethClient.CallContract(ctx, ethereum.CallMsg{
		From: c.GetAddress(),
		To:   &usdcAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("调用 allowance 失败: %w", err)
	}

	// 解析结果
	var allowance *big.Int
	err = parsedERC20ABI.UnpackIntoInterface(&allowance, "allowance", result)
	if err != nil {
		return nil, fmt.Errorf("解析 allowance 结果失败: %w", err)
	}

	return allowance, nil
}

// GetUSDCBalance 获取 USDC 余额（支持代理合约）
func (c *RelayerClient) GetUSDCBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	usdcAddress := common.HexToAddress(USDCAddress)

	// 构建 balanceOf 调用数据
	data, err := parsedERC20ABI.Pack("balanceOf", address)
	if err != nil {
		return nil, fmt.Errorf("构建 balanceOf 调用数据失败: %w", err)
	}

	// 调用合约
	result, err := c.ethClient.CallContract(ctx, ethereum.CallMsg{
		From: c.GetAddress(),
		To:   &usdcAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("调用 balanceOf 失败: %w", err)
	}

	// 解析结果
	var balance *big.Int
	err = parsedERC20ABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, fmt.Errorf("解析 balanceOf 结果失败: %w", err)
	}

	return balance, nil
}

// Close 关闭客户端
func (c *RelayerClient) Close() {
	c.ethClient.Close()
}

// ensureSafeDeployed 确保 Safe 已部署
func (c *RelayerClient) ensureSafeDeployed(ctx context.Context) error {
	deployed, err := c.IsSafeDeployed(ctx)
	if err != nil {
		return err
	}
	if !deployed {
		return fmt.Errorf("Safe 尚未部署，请先调用 DeploySafe")
	}
	return nil
}

// deriveSafeAddress 派生 Safe 地址
func (c *RelayerClient) deriveSafeAddress() common.Address {
	// 对齐 Python 实现：
	// salt = keccak256(abi.encode(["address"], [owner]))
	// safe = keccak256(0xff + factory + salt + SAFE_INIT_CODE_HASH)[12:]
	//
	// 参考：
	// /Users/dachang/Workspace/python/py-builder-relayer-client/py_builder_relayer_client/builder/derive.py
	const safeInitCodeHashHex = "0x2bce2127ff07fb632d16c8347c4ebf501f4841168bed00d9e6ef715ddb6fcecf"
	const safeFactoryHex = "0xaacFeEa03eb1561C4e67d661e40682Bd20E3541b"

	owner := c.GetAddress()
	factory := common.HexToAddress(safeFactoryHex)

	// abi.encode(address) -> 32 bytes left padded
	var encoded [32]byte
	copy(encoded[12:], owner.Bytes())
	salt := crypto.Keccak256Hash(encoded[:]) // 32 bytes

	initCodeHash := common.HexToHash(safeInitCodeHashHex).Bytes()
	buf := make([]byte, 0, 1+20+32+32)
	buf = append(buf, 0xff)
	buf = append(buf, factory.Bytes()...)
	buf = append(buf, salt.Bytes()...)
	buf = append(buf, initCodeHash...)
	h := crypto.Keccak256Hash(buf).Bytes()
	return common.BytesToAddress(h[12:])
}

// buildSafeSubmitRequest 构建 Safe 提交请求
func (c *RelayerClient) buildSafeSubmitRequest(ctx context.Context, req *RedeemRelayerRequest) (*safeSubmitRequest, error) {
	// 编码交易数据
	data, err := c.encodeRedeemData(req)
	if err != nil {
		return nil, err
	}

	// 派生目标合约地址
	targetContract := c.getTargetContract(req.IsNegRisk)

	// 获取 nonce (需要在签名前获取)
	nonce, err := c.fetchNonce(ctx)
	if err != nil {
		return nil, err
	}

	// 构建 EIP-712 签名
	signature, err := c.signSafeHash(ctx, targetContract, data, nonce)
	if err != nil {
		return nil, err
	}

	return &safeSubmitRequest{
		From:        c.GetAddress().Hex(),
		To:          targetContract.Hex(),
		ProxyWallet: c.safeAddress.Hex(),
		Data:        hexutil.Encode(data),
		Nonce:       fmt.Sprintf("%d", nonce),
		Signature:   hexutil.Encode(signature),
		SignatureParams: safeSignatureParams{
			GasPrice:       "0",
			Operation:      "0",
			SafeTxnGas:     "0",
			BaseGas:        "0",
			GasToken:       "0x0000000000000000000000000000000000000000",
			RefundReceiver: "0x0000000000000000000000000000000000000000",
		},
		Type:     "SAFE",
		Metadata: "",
	}, nil
}

// encodeRedeemData 编码赎回数据
func (c *RelayerClient) encodeRedeemData(req *RedeemRelayerRequest) ([]byte, error) {
	if req.IsNegRisk {
		return c.encodeNegRiskRedeemData(req)
	}
	return c.encodeCTFRedeemData(req)
}

// encodeCTFRedeemData 编码 CTF 赎回数据
func (c *RelayerClient) encodeCTFRedeemData(req *RedeemRelayerRequest) ([]byte, error) {
	// 验证必需参数
	if req.CollateralToken == (common.Address{}) {
		return nil, fmt.Errorf("CTF 赎回需要提供 CollateralToken")
	}
	if len(req.IndexSets) == 0 {
		return nil, fmt.Errorf("CTF 赎回需要提供 IndexSets")
	}

	// 解析 condition ID
	conditionID, err := hex.DecodeString(strings.TrimPrefix(req.ConditionID, "0x"))
	if err != nil {
		return nil, fmt.Errorf("解析 condition ID 失败: %w", err)
	}
	if len(conditionID) != 32 {
		return nil, fmt.Errorf("condition ID 长度必须为 32 字节")
	}

	// 转换为 32 字节数组
	var conditionIDBytes32 [32]byte
	copy(conditionIDBytes32[:], conditionID)

	// 使用解析后的 ABI 编码
	data, err := parsedCTFRedeemABI.Pack(
		"redeemPositions",
		req.CollateralToken,
		req.ParentCollection, // 使用零 hash 表示默认 parent collection
		conditionIDBytes32,
		req.IndexSets,
	)
	if err != nil {
		return nil, fmt.Errorf("编码 CTF 赎回数据失败: %w", err)
	}

	return data, nil
}

// encodeNegRiskRedeemData 编码 NegRisk 赎回数据
func (c *RelayerClient) encodeNegRiskRedeemData(req *RedeemRelayerRequest) ([]byte, error) {
	// 验证必需参数
	if len(req.RedeemAmounts) != 2 {
		return nil, fmt.Errorf("NegRisk 赎回需要提供 2 个金额: [yesAmount, noAmount]")
	}

	// 解析 condition ID
	conditionID, err := hex.DecodeString(strings.TrimPrefix(req.ConditionID, "0x"))
	if err != nil {
		return nil, fmt.Errorf("解析 condition ID 失败: %w", err)
	}
	if len(conditionID) != 32 {
		return nil, fmt.Errorf("condition ID 长度必须为 32 字节")
	}

	// 转换为 32 字节数组
	var conditionIDBytes32 [32]byte
	copy(conditionIDBytes32[:], conditionID)

	// 使用解析后的 ABI 编码
	data, err := parsedNegRiskABI.Pack(
		"redeemPositions",
		conditionIDBytes32,
		req.RedeemAmounts, // [yesAmount, noAmount]
	)
	if err != nil {
		return nil, fmt.Errorf("编码 NegRisk 赎回数据失败: %w", err)
	}

	return data, nil
}

// getTargetContract 获取目标合约地址
func (c *RelayerClient) getTargetContract(isNegRisk bool) common.Address {
	chainID := c.config.ChainID
	if chainID == 0 {
		chainID = DefaultChainID
	}
	contractCfg, err := GetContractConfig(chainID)
	if err != nil {
		// 未识别链时回退到 Polygon 主网配置，保持兼容。
		contractCfg = PolygonContractConfig
	}

	if isNegRisk {
		return common.HexToAddress(contractCfg.NegRiskAdapter)
	}
	return common.HexToAddress(contractCfg.ConditionalTokens)
}

// fetchNonce 获取 Safe nonce
func (c *RelayerClient) fetchNonce(ctx context.Context) (uint64, error) {
	// 使用地址和类型参数请求 nonce，保持与 JS 客户端一致
	u, err := url.Parse(c.config.RelayerURL)
	if err != nil {
		return 0, fmt.Errorf("解析 relayer URL 失败: %w", err)
	}

	u.Path = strings.TrimRight(u.Path, "/") + "/nonce"
	q := u.Query()
	q.Set("address", c.GetAddress().Hex())
	q.Set("type", "SAFE")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var nonceResp struct {
		Nonce string `json:"nonce"`
	}
	if err := json.Unmarshal(body, &nonceResp); err != nil {
		return 0, err
	}

	if nonceResp.Nonce == "" {
		return 0, fmt.Errorf("nonce 为空: %s", string(body))
	}

	nonceValue, ok := new(big.Int).SetString(nonceResp.Nonce, 10)
	if !ok {
		return 0, fmt.Errorf("解析 nonce 失败: %s", nonceResp.Nonce)
	}

	return nonceValue.Uint64(), nil
}

// signSafeHash 签名 Safe 哈希 (EIP-712)
func (c *RelayerClient) signSafeHash(ctx context.Context, to common.Address, data []byte, nonce uint64) ([]byte, error) {
	// 获取 chain ID
	chainID, err := c.ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取链 ID 失败: %w", err)
	}

	// 1. 构建类型哈希
	safeTxTypeHash := crypto.Keccak256Hash([]byte("SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)"))

	// 2. 构建各个字段的哈希值
	dataHash := crypto.Keccak256Hash(data)

	// 3. 将各字段编码为 32 字节数组 (遵循 ABI 编码)
	toBytes32 := common.LeftPadBytes(to.Bytes(), 32)
	valueBytes32 := common.LeftPadBytes(common.Big0.Bytes(), 32) // value = 0
	dataBytes32 := dataHash.Bytes()
	operationBytes32 := common.LeftPadBytes(big.NewInt(0).Bytes(), 32) // Call = 0
	safeTxGasBytes32 := common.LeftPadBytes(common.Big0.Bytes(), 32)   // safeTxGas = 0
	baseGasBytes32 := common.LeftPadBytes(common.Big0.Bytes(), 32)     // baseGas = 0
	gasPriceBytes32 := common.LeftPadBytes(common.Big0.Bytes(), 32)    // gasPrice = 0
	gasTokenBytes32 := common.LeftPadBytes(common.HexToAddress("0x0000000000000000000000000000000000000").Bytes(), 32)
	refundReceiverBytes32 := common.LeftPadBytes(common.HexToAddress("0x0000000000000000000000000000000000000").Bytes(), 32)
	nonceBytes32 := common.LeftPadBytes(new(big.Int).SetUint64(nonce).Bytes(), 32)

	// 4. 构建结构数据哈希
	safeTxData := make([]byte, 0, 32*10)
	safeTxData = append(safeTxData, safeTxTypeHash.Bytes()...)
	safeTxData = append(safeTxData, toBytes32...)
	safeTxData = append(safeTxData, valueBytes32...)
	safeTxData = append(safeTxData, dataBytes32...)
	safeTxData = append(safeTxData, operationBytes32...)
	safeTxData = append(safeTxData, safeTxGasBytes32...)
	safeTxData = append(safeTxData, baseGasBytes32...)
	safeTxData = append(safeTxData, gasPriceBytes32...)
	safeTxData = append(safeTxData, gasTokenBytes32...)
	safeTxData = append(safeTxData, refundReceiverBytes32...)
	safeTxData = append(safeTxData, nonceBytes32...)

	safeTxHash := crypto.Keccak256Hash(safeTxData)

	// 5. 构建域分隔符 (EIP-712 domain separator)
	domainSeparatorTypeHash := crypto.Keccak256Hash([]byte("EIP712Domain(uint256 chainId,address verifyingContract)"))
	chainIDBytes32 := common.LeftPadBytes(chainID.Bytes(), 32)
	verifyingContractBytes32 := common.LeftPadBytes(c.safeAddress.Bytes(), 32)

	domainSeparatorData := make([]byte, 0, 32+32+32)
	domainSeparatorData = append(domainSeparatorData, domainSeparatorTypeHash.Bytes()...)
	domainSeparatorData = append(domainSeparatorData, chainIDBytes32...)
	domainSeparatorData = append(domainSeparatorData, verifyingContractBytes32...)

	domainSeparator := crypto.Keccak256Hash(domainSeparatorData)

	// 6. 构建最终哈希 (0x1901 || domainSeparator || safeTxHash)
	finalHashData := make([]byte, 0, 2+32+32)
	finalHashData = append(finalHashData, 0x19, 0x01) // EIP-712 prefix
	finalHashData = append(finalHashData, domainSeparator.Bytes()...)
	finalHashData = append(finalHashData, safeTxHash.Bytes()...)

	finalHash := crypto.Keccak256Hash(finalHashData)

	// 7. 签名最终哈希（先用 EIP-191 prefix 进行二次哈希，再签名）
	prefixed := addEthereumMessagePrefix(finalHash.Bytes())
	signature, err := crypto.Sign(prefixed.Bytes(), c.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名 EIP-712 哈希失败: %w", err)
	}

	// 8. 调整 v 参数以符合 Relayer 的 Safe 格式 (31/32)
	if signature[64] < 27 {
		signature[64] += 27
	}
	if signature[64] == 27 {
		signature[64] = 31
	} else if signature[64] == 28 {
		signature[64] = 32
	}

	return signature, nil
}

// submit 提交请求到 relayer
func (c *RelayerClient) submit(ctx context.Context, path string, req interface{}) (*RelayerResponse, error) {
	// 编码请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("编码请求体失败: %w", err)
	}

	// 创建 HTTP 请求
	url := c.config.RelayerURL + path
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置头部
	httpReq.Header.Set("Content-Type", "application/json")

	// 添加 Builder Auth（如果提供）
	if c.config.BuilderAuth != nil && c.config.BuilderAuth.APIKey != "" {
		auth := NewBuilderAuth(c.config.BuilderAuth.APIKey, c.config.BuilderAuth.Secret, c.config.BuilderAuth.Passphrase)
		headers, err := auth.Headers("POST", path, reqBody)
		if err != nil {
			return nil, fmt.Errorf("获取认证头部失败: %w", err)
		}
		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}
	}

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var result RelayerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查 HTTP 状态
	if resp.StatusCode >= 400 {
		return &result, fmt.Errorf("relayer 错误: %d %s", resp.StatusCode, string(body))
	}

	return &result, nil
}

// 请求和响应结构体

type safeCreateRequest struct {
	From            string              `json:"from"`
	SafeAddress     string              `json:"safe"`
	SaltNonce       string              `json:"saltNonce"`
	Signature       string              `json:"signature"`
	SignatureParams safeSignatureParams `json:"signatureParams"`
	Type            string              `json:"type"`
	Metadata        string              `json:"metadata"`
}

type safeSubmitRequest struct {
	From            string              `json:"from"`
	To              string              `json:"to"`
	ProxyWallet     string              `json:"proxyWallet"`
	Data            string              `json:"data"`
	Nonce           string              `json:"nonce"`
	Signature       string              `json:"signature"`
	SignatureParams safeSignatureParams `json:"signatureParams"`
	Type            string              `json:"type"`
	Metadata        string              `json:"metadata"`
}

type safeSignatureParams struct {
	GasPrice       string `json:"gasPrice"`
	Operation      string `json:"operation"`
	SafeTxnGas     string `json:"safeTxnGas"`
	BaseGas        string `json:"baseGas"`
	GasToken       string `json:"gasToken"`
	RefundReceiver string `json:"refundReceiver"`
}

// CreateInfiniteApprovalAmount 创建无限授权金额 (2^256 - 1)
func CreateInfiniteApprovalAmount() *big.Int {
	result := new(big.Int)
	result.Exp(big.NewInt(2), big.NewInt(256), nil)
	result.Sub(result, big.NewInt(1))
	return result
}

// CreateUSDCAmount 创建 USDC 授权金额（USDC 有 6 位小数）
func CreateUSDCAmount(usdcAmount float64) *big.Int {
	amountWei := new(big.Int)
	amountFloat := big.NewFloat(usdcAmount)
	multiplier := big.NewFloat(1000000) // 6位小数
	amountFloat.Mul(amountFloat, multiplier)
	amountFloat.Int(amountWei)
	return amountWei
}

// PrintApprovalInfo 打印授权信息
func PrintApprovalInfo(usdcContractAddress string) {
	fmt.Println("=== USDC 授权信息 ===")
	fmt.Printf("USDC 合约地址: %s\n", usdcContractAddress)
	fmt.Printf("CTF Exchange 合约地址: %s\n", CTFExchangeAddress)
	fmt.Println("USDC 小数位数: 6")
	fmt.Println("1 USDC = 1,000,000 wei")

	maxUint256 := CreateInfiniteApprovalAmount()
	fmt.Printf("uint256 最大值: %s wei\n", maxUint256.String())

	usdcAmount := new(big.Float).SetInt(maxUint256)
	usdcAmount.Quo(usdcAmount, big.NewFloat(1000000))
	fmt.Printf("相当于: %s USDC (几乎无限)\n", usdcAmount.String())
}

func addEthereumMessagePrefix(message []byte) common.Hash {
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	prefixed := append([]byte(prefix), message...)
	return crypto.Keccak256Hash(prefixed)
}

// NewCTFRedeemRequest 创建 CTF 赎回请求
func NewCTFRedeemRequest(conditionID string, indexSets []*big.Int, collateralToken common.Address) *RedeemRelayerRequest {
	return &RedeemRelayerRequest{
		ConditionID:      conditionID,
		IndexSets:        indexSets,
		CollateralToken:  collateralToken,
		ParentCollection: [32]byte{}, // 使用零 hash
		IsNegRisk:        false,
	}
}

// NewNegRiskRedeemRequest 创建 NegRisk 赎回请求
func NewNegRiskRedeemRequest(conditionID string, redeemAmounts []*big.Int) *RedeemRelayerRequest {
	if len(redeemAmounts) != 2 {
		panic("NegRisk 赎回需要 2 个金额: [yesAmount, noAmount]")
	}
	return &RedeemRelayerRequest{
		ConditionID:   conditionID,
		RedeemAmounts: redeemAmounts, // [yesAmount, noAmount]
		IsNegRisk:     true,
	}
}

// RedeemCTFPositions 赎回 CTF 位置（便捷方法）
func (c *RelayerClient) RedeemCTFPositions(ctx context.Context, conditionID string, indexSets []*big.Int, collateralToken common.Address) (*RelayerResponse, error) {
	req := NewCTFRedeemRequest(conditionID, indexSets, collateralToken)
	return c.RedeemPositions(ctx, req)
}

// RedeemNegRiskPositions 赎回 NegRisk 位置（便捷方法）
func (c *RelayerClient) RedeemNegRiskPositions(ctx context.Context, conditionID string, yesAmount, noAmount *big.Int) (*RelayerResponse, error) {
	req := NewNegRiskRedeemRequest(conditionID, []*big.Int{yesAmount, noAmount})
	return c.RedeemPositions(ctx, req)
}
