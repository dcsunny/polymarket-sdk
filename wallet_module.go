// wallet_module.go 模块
package polymarket

import "context"

// WalletModule 提供使用 SDK 配置默认值的便捷构造函数。
type WalletModule struct {
	cfg Config
}

func NewWalletModule(cfg Config) *WalletModule {
	return &WalletModule{cfg: cfg}
}

// Safe 使用 SDK 配置默认值创建 Safe 钱包客户端。
func (w *WalletModule) Safe(ctx context.Context, cfg SafeWalletConfig) (WalletClient, error) {
	if cfg.RPCURL == "" {
		cfg.RPCURL = w.cfg.RPCURL
	}
	if cfg.PrivateKey == "" {
		cfg.PrivateKey = w.cfg.PrivateKey
	}
	if cfg.ChainID == 0 {
		cfg.ChainID = w.cfg.ChainID
	}
	return NewSafeWalletClient(ctx, cfg)
}

// Proxy 使用 SDK 配置默认值创建 Proxy 钱包客户端。
func (w *WalletModule) Proxy(ctx context.Context, cfg ProxyWalletConfig) (WalletClient, error) {
	if cfg.RPCURL == "" {
		cfg.RPCURL = w.cfg.RPCURL
	}
	if cfg.PrivateKey == "" {
		cfg.PrivateKey = w.cfg.PrivateKey
	}
	if cfg.ChainID == 0 {
		cfg.ChainID = w.cfg.ChainID
	}
	return NewProxyWalletClient(ctx, cfg)
}

// Relayer 使用 SDK 配置默认值创建 Relayer 客户端。
func (w *WalletModule) Relayer(ctx context.Context, cfg RelayerConfig) (*RelayerClient, error) {
	if cfg.RPCURL == "" {
		cfg.RPCURL = w.cfg.RPCURL
	}
	if cfg.PrivateKey == "" {
		cfg.PrivateKey = w.cfg.PrivateKey
	}
	if cfg.ChainID == 0 {
		cfg.ChainID = w.cfg.ChainID
	}
	if cfg.RelayerURL == "" {
		cfg.RelayerURL = w.cfg.RelayerURL
	}
	return NewRelayerClient(ctx, cfg)
}
