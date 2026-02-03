# Wallet / Relayer

钱包模块用于链上操作（Safe / Proxy）以及 Relayer 提交流程的封装。

入口：

```go
sdk, _ := pm.New(pm.Config{
    RPCURL: "...",
    PrivateKey: "0x...",
})

safe, _ := sdk.Wallet.Safe(ctx, pm.SafeWalletConfig{ /* ... */ })
proxy, _ := sdk.Wallet.Proxy(ctx, pm.ProxyWalletConfig{ /* ... */ })
relayer, _ := sdk.Wallet.Relayer(ctx, pm.RelayerConfig{ /* ... */ })
```

具体能力请参考：

- `wallet_client.go`
- `relayer_client.go`

