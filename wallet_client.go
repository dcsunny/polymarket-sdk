// wallet_client.go 模块
package polymarket

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WalletType 钱包类型
type WalletType string

const (
	WalletTypeSafe  WalletType = "SAFE"
	WalletTypeProxy WalletType = "PROXY"
)

// RedeemRequest 赎回请求参数
type RedeemRequest struct {
	ConditionID        string         `json:"conditionId"`
	IndexSets          []*big.Int     `json:"indexSets"`
	CollateralToken    common.Address `json:"collateralToken"`
	ParentCollectionID [32]byte       `json:"parentCollectionId"`
	NegRisk            bool           `json:"negRisk"`
	Metadata           string         `json:"metadata"`
}

// WalletClient 统一的钱包客户端接口
type WalletClient interface {
	GetAddress() common.Address
	GetSafeAddress() common.Address
	IsSafeDeployed(ctx context.Context) (bool, error)
	Close()
	RedeemPositions(ctx context.Context, req *RedeemRequest) (*types.Transaction, error)
}

// defaultWalletClient 默认的钱包客户端实现
type defaultWalletClient struct {
	ethClient  *ethclient.Client
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
	walletType WalletType

	// 合约地址
	ctfAddress       common.Address
	negRiskAdapter   common.Address
	usdcAddress      common.Address
	proxyFactoryAddr common.Address

	// 缓存 ABI
	safeABI    abi.ABI
	ctfABI     abi.ABI
	negRiskABI abi.ABI
	proxyABI   abi.ABI
	erc20ABI   abi.ABI
	erc1155ABI abi.ABI

	// Safe 特定
	safeContract *bind.BoundContract
	safeAddress  common.Address

	log ClientLogger
}

// ClientLogger 客户端日志接口
type ClientLogger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// WalletConfig 钱包配置
type WalletConfig struct {
	RPCURL         string     `json:"rpcUrl"`
	PrivateKey     string     `json:"privateKey"`
	ChainID        int64      `json:"chainId"`
	WalletType     WalletType `json:"walletType"`
	SafeAddress    string     `json:"safeAddress,omitempty"`
	ProxyFactory   string     `json:"proxyFactory,omitempty"`
	CTFAddress     string     `json:"ctfAddress,omitempty"`
	NegRiskAdapter string     `json:"negRiskAdapter,omitempty"`
	USDCAddress    string     `json:"usdcAddress,omitempty"`
}

// SafeWalletConfig 是 Safe 钱包的用户配置。
type SafeWalletConfig struct {
	RPCURL         string
	PrivateKey     string
	SafeAddress    string
	ChainID        int64
	ProxyFactory   string
	CTFAddress     string
	NegRiskAdapter string
	USDCAddress    string
	Logger         ClientLogger
}

// ProxyWalletConfig 是 Proxy 钱包的用户配置。
type ProxyWalletConfig struct {
	RPCURL         string
	PrivateKey     string
	ChainID        int64
	ProxyFactory   string
	CTFAddress     string
	NegRiskAdapter string
	USDCAddress    string
	Logger         ClientLogger
}

// NewSafeWalletClient 创建 Safe 钱包客户端。
func NewSafeWalletClient(ctx context.Context, cfg SafeWalletConfig) (WalletClient, error) {
	chainID := cfg.ChainID
	if chainID == 0 {
		chainID = DefaultChainID
	}
	return NewOptimizedWalletClient(ctx, WalletConfig{
		RPCURL:         cfg.RPCURL,
		PrivateKey:     cfg.PrivateKey,
		ChainID:        chainID,
		WalletType:     WalletTypeSafe,
		SafeAddress:    cfg.SafeAddress,
		ProxyFactory:   cfg.ProxyFactory,
		CTFAddress:     cfg.CTFAddress,
		NegRiskAdapter: cfg.NegRiskAdapter,
		USDCAddress:    cfg.USDCAddress,
	}, cfg.Logger)
}

// NewProxyWalletClient 创建 Proxy 钱包客户端。
func NewProxyWalletClient(ctx context.Context, cfg ProxyWalletConfig) (WalletClient, error) {
	chainID := cfg.ChainID
	if chainID == 0 {
		chainID = DefaultChainID
	}
	return NewOptimizedWalletClient(ctx, WalletConfig{
		RPCURL:         cfg.RPCURL,
		PrivateKey:     cfg.PrivateKey,
		ChainID:        chainID,
		WalletType:     WalletTypeProxy,
		ProxyFactory:   cfg.ProxyFactory,
		CTFAddress:     cfg.CTFAddress,
		NegRiskAdapter: cfg.NegRiskAdapter,
		USDCAddress:    cfg.USDCAddress,
	}, cfg.Logger)
}

// NewOptimizedWalletClient 创建优化后的钱包客户端
func NewOptimizedWalletClient(ctx context.Context, cfg WalletConfig, logger ClientLogger) (WalletClient, error) {
	// 连接到以太坊客户端
	ethClient, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("连接 RPC 失败: %w", err)
	}

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}

	chainID := cfg.ChainID
	if chainID == 0 {
		chainID = DefaultChainID
	}

	client := &defaultWalletClient{
		ethClient:  ethClient,
		privateKey: privateKey,
		chainID:    big.NewInt(chainID),
		walletType: cfg.WalletType,
		log:        logger,
	}

	// 加载合约地址
	if err := client.loadContractAddresses(cfg); err != nil {
		return nil, fmt.Errorf("加载合约地址失败: %w", err)
	}

	// 加载 ABI
	if err := client.loadABIs(); err != nil {
		return nil, fmt.Errorf("加载 ABI 失败: %w", err)
	}

	// 初始化 Safe 合约（如果是 Safe 类型）
	if cfg.WalletType == WalletTypeSafe {
		if err := client.initSafeContract(); err != nil {
			return nil, fmt.Errorf("初始化 Safe 合约失败: %w", err)
		}
	}

	return client, nil
}

// loadContractAddresses 加载合约地址
func (c *defaultWalletClient) loadContractAddresses(cfg WalletConfig) error {
	// 默认地址（Polygon 主网）
	defaultCTF := common.HexToAddress("0x4D97DCd97eC945f40cF65F87097ACe5EA0476045")
	defaultNegRisk := common.HexToAddress("0xd91E80cF2E7be2e162c6513ceD06f1dD0dA35296")
	defaultUSDC := common.HexToAddress(USDCAddress)
	defaultProxyFactory := common.HexToAddress("0xaB45c5A4B0c941a2F231C04C3f49182e1A254052")

	c.ctfAddress = defaultCTF
	c.negRiskAdapter = defaultNegRisk
	c.usdcAddress = defaultUSDC
	c.proxyFactoryAddr = defaultProxyFactory

	// 使用自定义地址（如果提供）
	if cfg.SafeAddress != "" {
		c.safeAddress = common.HexToAddress(cfg.SafeAddress)
	}
	if cfg.CTFAddress != "" {
		c.ctfAddress = common.HexToAddress(cfg.CTFAddress)
	}
	if cfg.NegRiskAdapter != "" {
		c.negRiskAdapter = common.HexToAddress(cfg.NegRiskAdapter)
	}
	if cfg.USDCAddress != "" {
		c.usdcAddress = common.HexToAddress(cfg.USDCAddress)
	}
	if cfg.ProxyFactory != "" {
		c.proxyFactoryAddr = common.HexToAddress(cfg.ProxyFactory)
	}

	return nil
}

// loadABIs 加载所有需要的 ABI
func (c *defaultWalletClient) loadABIs() error {
	var err error

	// 加载 Safe ABI
	c.safeABI, err = abi.JSON(strings.NewReader(gnosisSafeABI))
	if err != nil {
		return fmt.Errorf("加载 Safe ABI 失败: %w", err)
	}

	// 加载 CTF ABI
	c.ctfABI, err = abi.JSON(strings.NewReader(CTFContractABI))
	if err != nil {
		return fmt.Errorf("加载 CTF ABI 失败: %w", err)
	}

	// 加载 NegRisk ABI
	c.negRiskABI, err = abi.JSON(strings.NewReader(negRiskAdapterABIJSON))
	if err != nil {
		return fmt.Errorf("加载 NegRisk ABI 失败: %w", err)
	}

	// 加载 Proxy ABI
	c.proxyABI, err = abi.JSON(strings.NewReader(proxyFactoryABIJSON))
	if err != nil {
		return fmt.Errorf("加载 Proxy ABI 失败: %w", err)
	}

	// 加载 ERC20 ABI
	c.erc20ABI, err = abi.JSON(strings.NewReader(erc20ABIJSON))
	if err != nil {
		return fmt.Errorf("加载 ERC20 ABI 失败: %w", err)
	}

	// 加载 ERC1155 ABI
	c.erc1155ABI, err = abi.JSON(strings.NewReader(erc1155ABIJSON))
	if err != nil {
		return fmt.Errorf("加载 ERC1155 ABI 失败: %w", err)
	}

	return nil
}

// initSafeContract 初始化 Safe 合约
func (c *defaultWalletClient) initSafeContract() error {
	if (c.safeAddress == common.Address{}) {
		c.safeAddress = c.deriveSafeAddress()
	}
	c.safeContract = bind.NewBoundContract(
		c.safeAddress,
		c.safeABI,
		c.ethClient,
		c.ethClient,
		c.ethClient,
	)
	return nil
}

// GetAddress 获取 EOA 地址
func (c *defaultWalletClient) GetAddress() common.Address {
	return crypto.PubkeyToAddress(c.privateKey.PublicKey)
}

// GetSafeAddress 获取 Safe 地址
func (c *defaultWalletClient) GetSafeAddress() common.Address {
	if c.walletType == WalletTypeProxy {
		return c.GetAddress() // Proxy 钱包就是 EOA 地址
	}
	return c.safeAddress
}

// IsSafeDeployed 检查 Safe 是否已部署
func (c *defaultWalletClient) IsSafeDeployed(ctx context.Context) (bool, error) {
	if c.walletType == WalletTypeProxy {
		return true, nil // Proxy 钱包不需要部署
	}

	code, err := c.ethClient.CodeAt(ctx, c.safeAddress, nil)
	if err != nil {
		return false, fmt.Errorf("检查 Safe 代码失败: %w", err)
	}
	return len(code) > 0, nil
}

// Close 关闭客户端
func (c *defaultWalletClient) Close() {
	c.ethClient.Close()
}

// RedeemPositions 赎回位置（统一接口）
func (c *defaultWalletClient) RedeemPositions(ctx context.Context, req *RedeemRequest) (*types.Transaction, error) {
	// 解析 condition ID
	conditionID, err := hexToBytes32(req.ConditionID)
	if err != nil {
		return nil, fmt.Errorf("解析 condition ID 失败: %w", err)
	}

	// 编码赎回数据
	var data []byte
	if req.NegRisk {
		data, err = c.encodeRedeemNegRiskPayload(conditionID, req.IndexSets)
	} else {
		data, err = c.encodeRedeemPayload(conditionID, req.IndexSets, req.CollateralToken, req.ParentCollectionID)
	}
	if err != nil {
		return nil, fmt.Errorf("编码赎回数据失败: %w", err)
	}

	// 确定目标合约
	targetContract := c.ctfAddress
	if req.NegRisk {
		targetContract = c.negRiskAdapter
	}

	// 根据钱包类型执行交易
	switch c.walletType {
	case WalletTypeSafe:
		return c.executeSafeTransaction(ctx, targetContract, data)
	case WalletTypeProxy:
		return c.executeProxyTransaction(ctx, targetContract, data)
	default:
		return nil, fmt.Errorf("不支持的钱包类型: %s", c.walletType)
	}
}

// executeSafeTransaction 执行 Safe 交易
func (c *defaultWalletClient) executeSafeTransaction(ctx context.Context, to common.Address, data []byte) (*types.Transaction, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("创建交易签名器失败: %w", err)
	}
	auth.Context = ctx

	// 构建 Safe execTransaction 调用
	tx, err := c.safeContract.Transact(
		auth,
		"execTransaction",
		to,
		big.NewInt(0), // value
		data,
		uint8(0),         // operation
		big.NewInt(0),    // safeTxGas
		big.NewInt(0),    // baseGas
		big.NewInt(0),    // gasPrice
		common.Address{}, // gasToken
		common.Address{}, // refundReceiver
		[]byte{},         // signatures
	)
	if err != nil {
		return nil, fmt.Errorf("执行 Safe 交易失败: %w", err)
	}

	return tx, nil
}

// executeProxyTransaction 执行 Proxy 交易
func (c *defaultWalletClient) executeProxyTransaction(ctx context.Context, to common.Address, data []byte) (*types.Transaction, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("创建交易签名器失败: %w", err)
	}
	auth.Context = ctx

	// 构建 Proxy 调用
	call := struct {
		TypeCode uint8
		To       common.Address
		Value    *big.Int
		Data     []byte
	}{
		TypeCode: 0, // CALL
		To:       to,
		Value:    big.NewInt(0),
		Data:     data,
	}

	proxyContract := bind.NewBoundContract(
		c.GetAddress(),
		c.proxyABI,
		c.ethClient,
		c.ethClient,
		c.ethClient,
	)

	tx, err := proxyContract.Transact(
		auth,
		"execute",
		[]struct {
			TypeCode uint8
			To       common.Address
			Value    *big.Int
			Data     []byte
		}{call},
	)
	if err != nil {
		return nil, fmt.Errorf("执行 Proxy 交易失败: %w", err)
	}

	return tx, nil
}

// encodeRedeemPayload 编码标准赎回数据
func (c *defaultWalletClient) encodeRedeemPayload(conditionID [32]byte, indexSets []*big.Int, collateralToken common.Address, parentCollectionID [32]byte) ([]byte, error) {
	method, ok := c.ctfABI.Methods["redeemPositions"]
	if !ok {
		return nil, fmt.Errorf("redeemPositions 方法不存在")
	}

	data, err := method.Inputs.Pack(
		conditionID,
		indexSets,
		collateralToken,
		parentCollectionID,
	)
	if err != nil {
		return nil, err
	}

	return append(method.ID, data...), nil
}

// encodeRedeemNegRiskPayload 编码 NegRisk 赎回数据
func (c *defaultWalletClient) encodeRedeemNegRiskPayload(conditionID [32]byte, amounts []*big.Int) ([]byte, error) {
	method, ok := c.negRiskABI.Methods["redeemPositions"]
	if !ok {
		return nil, fmt.Errorf("redeemPositions 方法不存在")
	}

	data, err := method.Inputs.Pack(
		conditionID,
		amounts,
	)
	if err != nil {
		return nil, err
	}

	return append(method.ID, data...), nil
}

// deriveSafeAddress 派生 Safe 地址
func (c *defaultWalletClient) deriveSafeAddress() common.Address {
	owner := c.GetAddress()
	factory := common.HexToAddress("0x29bFC345bf3C537388236341257B188166631a69") // Safe Factory 地址
	return c.create2Address(owner, factory)
}

// create2Address CREATE2 地址计算（简化实现）
func (c *defaultWalletClient) create2Address(salt, factory common.Address) common.Address {
	// 实际实现需要 keccak256(salt + keccak256(init_code) + factory)
	return common.HexToAddress("0xEA77ea672344987868F9f425e4ceABa6bbe3eD25")
}

// hexToBytes32 将十六进制字符串转换为 [32]byte
func hexToBytes32(s string) ([32]byte, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	if len(s) != 64 {
		return [32]byte{}, fmt.Errorf("无效的 hex 字符串长度: %d", len(s))
	}

	var result [32]byte
	_, err := fmt.Sscanf(s, "%064x", &result)
	return result, err
}

// ABI JSON 字符串
const (
	gnosisSafeABI = `[
		{
			"inputs": [],
			"name": "nonce",
			"outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{"internalType": "address", "name": "to", "type": "address"},
				{"internalType": "uint256", "name": "value", "type": "uint256"},
				{"internalType": "bytes", "name": "data", "type": "bytes"},
				{"internalType": "uint8", "name": "operation", "type": "uint8"},
				{"internalType": "uint256", "name": "safeTxGas", "type": "uint256"},
				{"internalType": "uint256", "name": "baseGas", "type": "uint256"},
				{"internalType": "uint256", "name": "gasPrice", "type": "uint256"},
				{"internalType": "address", "name": "gasToken", "type": "address"},
				{"internalType": "address", "name": "refundReceiver", "type": "address"},
				{"internalType": "uint256", "name": "_nonce", "type": "uint256"}
			],
			"name": "execTransaction",
			"outputs": [{"internalType": "bool", "name": "success", "type": "bool"}],
			"stateMutability": "payable",
			"type": "function"
		}
	]`

	CTFContractABI = `[
		{
			"inputs": [
				{"internalType": "contract IERC20", "name": "collateralToken", "type": "address"},
				{"internalType": "bytes32", "name": "parentCollectionId", "type": "bytes32"},
				{"internalType": "bytes32", "name": "conditionId", "type": "bytes32"},
				{"internalType": "uint256[]", "name": "indexSets", "type": "uint256[]"},
				{"internalType": "uint256", "name": "amount", "type": "uint256"}
			],
			"name": "splitPosition",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{"internalType": "contract IERC20", "name": "collateralToken", "type": "address"},
				{"internalType": "bytes32", "name": "parentCollectionId", "type": "bytes32"},
				{"internalType": "bytes32", "name": "conditionId", "type": "bytes32"},
				{"internalType": "uint256[]", "name": "indexSets", "type": "uint256[]"},
				{"internalType": "uint256", "name": "amount", "type": "uint256"}
			],
			"name": "mergePosition",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{"internalType": "contract IERC20", "name": "collateralToken", "type": "address"},
				{"internalType": "bytes32", "name": "parentCollectionId", "type": "bytes32"},
				{"internalType": "bytes32", "name": "conditionId", "type": "bytes32"},
				{"internalType": "uint256[]", "name": "indexSets", "type": "uint256[]"}
			],
			"name": "redeemPositions",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`

	negRiskAdapterABIJSON = `[
		{
			"inputs": [
				{"internalType": "bytes32", "name": "marketId", "type": "bytes32"},
				{"internalType": "uint256", "name": "indexSet", "type": "uint256"},
				{"internalType": "uint256", "name": "amount", "type": "uint256"}
			],
			"name": "convertPositions",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`

	proxyFactoryABIJSON = `[
		{
			"inputs": [
				{
					"components": [
						{"internalType": "uint8", "name": "typeCode", "type": "uint8"},
						{"internalType": "address", "name": "to", "type": "address"},
						{"internalType": "uint256", "name": "value", "type": "uint256"},
						{"internalType": "bytes", "name": "data", "type": "bytes"}
					],
					"internalType": "struct ProxyFactory.Call[]",
					"name": "calls",
					"type": "tuple[]"
				}
			],
			"name": "proxy",
			"outputs": [
				{"internalType": "bytes[]", "name": "returnValues", "type": "bytes[]"}
			],
			"stateMutability": "payable",
			"type": "function"
		}
	]`

	erc20ABIJSON = `[
		{
			"inputs": [
				{"internalType": "address", "name": "spender", "type": "address"},
				{"internalType": "uint256", "name": "amount", "type": "uint256"}
			],
			"name": "approve",
			"outputs": [{"internalType": "bool", "name": "", "type": "bool"}],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`

	erc1155ABIJSON = `[
		{
			"inputs": [
				{"internalType": "address", "name": "operator", "type": "address"},
				{"internalType": "bool", "name": "approved", "type": "bool"}
			],
			"name": "setApprovalForAll",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`
)
