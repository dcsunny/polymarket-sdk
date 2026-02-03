// client.go 模块
package polymarket

import (
	"errors"

	"github.com/dcsunny/polymarket-sdk/internal/httpx"
)

// SDK 是 Polymarket SDK 的主入口。
type SDK struct {
	cfg Config

	REST   *RESTClient
	CLOB   *CLOBClient
	WSS    *WSSClient
	RTDS   *RTDSClient
	Wallet *WalletModule
}

// New 创建一个应用了默认配置的新 SDK 客户端。
func New(cfg Config) (*SDK, error) {
	cfg = cfg.withDefaults()

	if cfg.BaseURL == "" || cfg.CLOBBaseURL == "" {
		return nil, errors.New("base urls are required")
	}

	restHTTP, err := httpx.New(cfg.BaseURL, cfg.Timeout, cfg.Proxy, cfg.UserAgent, cfg.Debug)
	if err != nil {
		return nil, err
	}
	clobHTTP, err := httpx.New(cfg.CLOBBaseURL, cfg.Timeout, cfg.Proxy, cfg.UserAgent, cfg.Debug)
	if err != nil {
		return nil, err
	}

	sdk := &SDK{cfg: cfg}
	sdk.REST = NewRESTClient(restHTTP)
	sdk.CLOB = NewCLOBClient(clobHTTP, cfg)
	sdk.WSS = NewWSSClient(cfg)
	sdk.RTDS = NewRTDSClient(cfg)
	sdk.Wallet = NewWalletModule(cfg)

	return sdk, nil
}

// Config 返回 SDK 配置的副本。
func (s *SDK) Config() Config {
	return s.cfg
}
