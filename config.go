// config.go 模块
package polymarket

import "time"

const (
	DefaultBaseURL      = "https://gamma-api.polymarket.com"
	DefaultCLOBBaseURL  = "https://clob.polymarket.com"
	DefaultWSSMarketURL = "wss://ws-subscriptions-clob.polymarket.com/ws/market"
	DefaultWSSUserURL   = "wss://ws-subscriptions-clob.polymarket.com/ws/user"
	DefaultRTDSURL      = "wss://ws-live-data.polymarket.com"
	DefaultRelayerURL   = "https://relay-v2.polymarket.com/"
	DefaultTimeout      = 30 * time.Second
	DefaultChainID      = ChainIDPolygon
)

const (
	SignatureTypeEOA            = 0
	SignatureTypePolyProxy      = 1
	SignatureTypePolyGnosisSafe = 2
)

// Config 定义 SDK 配置。
type Config struct {
	BaseURL      string
	CLOBBaseURL  string
	WSSMarketURL string
	WSSUserURL   string
	RTDSURL      string
	RelayerURL   string

	Timeout   time.Duration
	Proxy     string
	Debug     bool
	ChainID   int64
	UserAgent string

	// Auth
	Address    string
	PrivateKey string
	APIKey     string
	APISecret  string
	Passphrase string

	SignatureType int
	Funder        string

	// Builder（可选）：用于 builder flow（下单时注入 builder headers）
	BuilderAPIKey     string
	BuilderAPISecret  string
	BuilderPassphrase string

	// Wallet
	RPCURL      string
	BuilderAuth string
}

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.CLOBBaseURL == "" {
		c.CLOBBaseURL = DefaultCLOBBaseURL
	}
	if c.WSSMarketURL == "" {
		c.WSSMarketURL = DefaultWSSMarketURL
	}
	if c.WSSUserURL == "" {
		c.WSSUserURL = DefaultWSSUserURL
	}
	if c.RTDSURL == "" {
		c.RTDSURL = DefaultRTDSURL
	}
	if c.RelayerURL == "" {
		c.RelayerURL = DefaultRelayerURL
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}
	if c.ChainID == 0 {
		c.ChainID = DefaultChainID
	}
	if c.SignatureType == 0 {
		c.SignatureType = SignatureTypeEOA
	}
	return c
}
