// clob.go 模块
package polymarket

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dcsunny/polymarket-sdk/internal/auth"
	"github.com/dcsunny/polymarket-sdk/internal/httpx"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// CLOBClient 处理 Polymarket CLOB API。
type CLOBClient struct {
	http *httpx.Client
	cfg  Config

	address    string
	privateKey string
	apiKey     string
	apiSecret  string
	passphrase string

	orderBuilder *OrderBuilder
	sigType      int
	funder       string
	chainID      int64

	builderAuth *BuilderAuth

	tickSizeCache map[string]string
	negRiskCache  map[string]bool
	feeRateCache  map[string]int
}

func NewCLOBClient(http *httpx.Client, cfg Config) *CLOBClient {
	if cfg.Address == "" && cfg.PrivateKey != "" {
		if addr, err := PrivateKeyToAddress(cfg.PrivateKey); err == nil {
			cfg.Address = addr
		}
	}

	chainID := cfg.ChainID
	if chainID == 0 {
		chainID = DefaultChainID
	}

	client := &CLOBClient{
		http:          http,
		cfg:           cfg,
		address:       cfg.Address,
		privateKey:    cfg.PrivateKey,
		apiKey:        cfg.APIKey,
		apiSecret:     cfg.APISecret,
		passphrase:    cfg.Passphrase,
		sigType:       cfg.SignatureType,
		funder:        cfg.Funder,
		chainID:       chainID,
		tickSizeCache: make(map[string]string),
		negRiskCache:  make(map[string]bool),
		feeRateCache:  make(map[string]int),
	}

	if cfg.BuilderAPIKey != "" && cfg.BuilderAPISecret != "" && cfg.BuilderPassphrase != "" {
		client.builderAuth = NewBuilderAuth(cfg.BuilderAPIKey, cfg.BuilderAPISecret, cfg.BuilderPassphrase)
	}

	if client.privateKey != "" && client.address != "" {
		client.orderBuilder = NewOrderBuilder(client, client.privateKey, client.address, client.sigType, client.funder)
	}
	return client
}

// APIKey 返回当前的 API 密钥。
func (c *CLOBClient) APIKey() string {
	return c.apiKey
}

// ChainID 返回订单所使用的链 ID。
func (c *CLOBClient) ChainID() int64 {
	return c.chainID
}

// SetAPICredentials 更新客户端上的 L2 凭证。
func (c *CLOBClient) SetAPICredentials(key, secret, passphrase string) {
	c.apiKey = key
	c.apiSecret = secret
	c.passphrase = passphrase
}

func (c *CLOBClient) l2Headers(method, path, body string) (map[string]string, error) {
	if c.address == "" || c.apiKey == "" || c.apiSecret == "" || c.passphrase == "" {
		return nil, errors.New("missing L2 credentials")
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := auth.L2Signature(c.apiSecret, timestamp, method, path, body)
	return map[string]string{
		"POLY_ADDRESS":    c.address,
		"POLY_SIGNATURE":  signature,
		"POLY_TIMESTAMP":  timestamp,
		"POLY_API_KEY":    c.apiKey,
		"POLY_PASSPHRASE": c.passphrase,
	}, nil
}

func (c *CLOBClient) l1Headers(nonce int) (map[string]string, error) {
	if c.address == "" || c.privateKey == "" {
		return nil, errors.New("missing address or private key")
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig, err := auth.ClobAuthSignature(c.privateKey, c.address, timestamp, nonce, c.cfg.ChainID)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"POLY_ADDRESS":   c.address,
		"POLY_SIGNATURE": sig,
		"POLY_TIMESTAMP": timestamp,
		"POLY_NONCE":     strconv.Itoa(nonce),
	}, nil
}

// CreateAPIKey 创建新的 API 密钥（L1 认证）。
func (c *CLOBClient) CreateAPIKey(ctx context.Context, nonce int) (*APICredentials, error) {
	headers, err := c.l1Headers(nonce)
	if err != nil {
		return nil, err
	}

	var creds APICredentials
	if err := c.http.Do(ctx, http.MethodPost, "/auth/api-key", nil, nil, headers, &creds); err != nil {
		return nil, err
	}
	c.SetAPICredentials(creds.APIKey, creds.Secret, creds.Passphrase)
	return &creds, nil
}

// DeriveAPIKey 推导现有的 API 密钥（L1 认证）。
func (c *CLOBClient) DeriveAPIKey(ctx context.Context, nonce int) (*APICredentials, error) {
	headers, err := c.l1Headers(nonce)
	if err != nil {
		return nil, err
	}

	var creds APICredentials
	if err := c.http.Do(ctx, http.MethodGet, "/auth/derive-api-key", nil, nil, headers, &creds); err != nil {
		return nil, err
	}
	c.SetAPICredentials(creds.APIKey, creds.Secret, creds.Passphrase)
	return &creds, nil
}

// CreateOrDeriveAPIKey 尝试先推导，然后创建。
func (c *CLOBClient) CreateOrDeriveAPIKey(ctx context.Context, nonce int) (*APICredentials, error) {
	creds, err := c.DeriveAPIKey(ctx, nonce)
	if err == nil {
		return creds, nil
	}
	return c.CreateAPIKey(ctx, nonce)
}

// GetAPIKeys 列出 API 密钥（L2 认证）。
func (c *CLOBClient) GetAPIKeys(ctx context.Context) (*APIKeysResponse, error) {
	headers, err := c.l2Headers(http.MethodGet, "/auth/api-keys", "")
	if err != nil {
		return nil, err
	}

	var result APIKeysResponse
	if err := c.http.Do(ctx, http.MethodGet, "/auth/api-keys", nil, nil, headers, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteAPIKey 删除当前的 API 密钥（L2 认证）。
func (c *CLOBClient) DeleteAPIKey(ctx context.Context) error {
	headers, err := c.l2Headers(http.MethodDelete, "/auth/api-key", "")
	if err != nil {
		return err
	}
	return c.http.Do(ctx, http.MethodDelete, "/auth/api-key", nil, nil, headers, nil)
}

// GetActiveOrders 返回活跃订单，自动分页直到结束。
func (c *CLOBClient) GetActiveOrders(ctx context.Context, req *GetActiveOrdersRequest) ([]*OpenOrder, error) {
	nextCursor := InitialCursor
	var all []*OpenOrder

	for nextCursor != EndCursor {
		resp, err := c.GetActiveOrdersPage(ctx, req, nextCursor)
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Data...)
		nextCursor = resp.NextCursor
	}

	return all, nil
}

// GetOpenOrders 是 GetActiveOrders 的别名（与官方文档/Node SDK 命名对齐）。
func (c *CLOBClient) GetOpenOrders(ctx context.Context, req *GetActiveOrdersRequest) ([]*OpenOrder, error) {
	return c.GetActiveOrders(ctx, req)
}

// GetOpenOrdersPage 是 GetActiveOrdersPage 的别名（与官方文档/Node SDK 命名对齐）。
func (c *CLOBClient) GetOpenOrdersPage(ctx context.Context, req *GetActiveOrdersRequest, nextCursor string) (*GetActiveOrdersResponse, error) {
	return c.GetActiveOrdersPage(ctx, req, nextCursor)
}

// GetActiveOrdersPage 获取活跃订单单页（GET /data/orders）。
// nextCursor 为空时默认使用 InitialCursor。
func (c *CLOBClient) GetActiveOrdersPage(ctx context.Context, req *GetActiveOrdersRequest, nextCursor string) (*GetActiveOrdersResponse, error) {
	path := "/data/orders"
	vals := url.Values{}
	if req != nil {
		if req.ID != "" {
			vals.Set("id", req.ID)
		}
		if req.Market != "" {
			vals.Set("market", req.Market)
		}
		if req.AssetID != "" {
			vals.Set("asset_id", req.AssetID)
		}
	}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	vals.Set("next_cursor", nextCursor)

	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp GetActiveOrdersResponse
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOrder 根据哈希获取单个订单。
func (c *CLOBClient) GetOrder(ctx context.Context, orderHash string) (*OpenOrder, error) {
	if orderHash == "" {
		return nil, ErrInvalidArgument("orderHash is required")
	}
	path := "/data/order/" + url.PathEscape(orderHash)
	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}
	var raw []byte
	if err := c.http.Do(ctx, http.MethodGet, path, nil, nil, headers, &raw); err != nil {
		return nil, err
	}

	// 兼容不同返回格式：直接返回 order 或 {order: {...}}
	var wrapped struct {
		Order OpenOrder `json:"order"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Order.ID != "" {
		return &wrapped.Order, nil
	}

	var order OpenOrder
	if err := json.Unmarshal(raw, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

// Price 返回代币侧的价格。
func (c *CLOBClient) Price(ctx context.Context, tokenID string, side PriceSide) (string, error) {
	if tokenID == "" {
		return "", ErrInvalidArgument("tokenID is required")
	}
	if side != PriceSideBuy && side != PriceSideSell {
		return "", ErrInvalidArgument("side must be BUY or SELL")
	}
	vals := url.Values{}
	vals.Set("token_id", tokenID)
	vals.Set("side", side.String())

	var resp PriceResponse
	if err := c.http.Do(ctx, http.MethodGet, "/price", vals, nil, nil, &resp); err != nil {
		return "", err
	}
	return resp.Price, nil
}

// BuyPrice 返回代币的买入价格。
func (c *CLOBClient) BuyPrice(ctx context.Context, tokenID string) (string, error) {
	return c.Price(ctx, tokenID, PriceSideBuy)
}

// SellPrice 返回代币的卖出价格。
func (c *CLOBClient) SellPrice(ctx context.Context, tokenID string) (string, error) {
	return c.Price(ctx, tokenID, PriceSideSell)
}

// PrivateKeyToAddress 将私钥转换为地址。
func PrivateKeyToAddress(privateKeyHex string) (string, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}
	pub := privateKey.PublicKey
	return crypto.PubkeyToAddress(pub).Hex(), nil
}

// ParsePrivateKey 验证并解析十六进制私钥。
func ParsePrivateKey(privateKeyHex string) error {
	_, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return err
	}
	return nil
}

// ValidateAddress 检查十六进制地址是否有效。
func ValidateAddress(address string) bool {
	return common.IsHexAddress(address)
}

// GenerateNonce 返回基于时间戳的随机数。
func GenerateNonce() int {
	return int(time.Now().UnixNano() % 1000000)
}

// HexSignature 规范化签名十六进制字符串。
func HexSignature(sig []byte) string {
	return "0x" + hex.EncodeToString(sig)
}
