# Polymarket Go SDK

面向 Polymarket 的 Go SDK，提供 REST、CLOB、WSS、RTDS、钱包与 Relayer 等能力，强调“简单、清晰、可用”。

## 功能概览

- REST：事件与市场查询
- CLOB：订单管理、订单簿、交易、价格与评分
- WSS：市场与用户频道订阅
- RTDS：实时行情订阅
- 钱包：Safe / Proxy / Relayer

## 快速开始

```go
import pm "github.com/dcsunny/polymarket-sdk"

sdk, _ := pm.New(pm.Config{
    Address:    "0x...",
    APIKey:     "...",
    APISecret:  "...",
    Passphrase: "...",
})

orders, _ := sdk.CLOB.GetActiveOrders(ctx, &pm.GetActiveOrdersRequest{})
```

## 订单管理接口

- `GetOrder`：获取单个订单
- `GetActiveOrders` / `GetActiveOrdersPage`：获取活跃订单（支持分页/自动拉取）
- `PostOrder` / `PostOrderWithOptions`：下单（支持 `deferExec` / `postOnly`）
- `PostOrders` / `PostOrdersSigned`：批量下单
- `CancelOrder` / `CancelOrders` / `CancelAllOrders` / `CancelMarketOrders`
- `GetOrderBook` / `GetOrderBooks`
- `IsOrderScoring` / `AreOrdersScoring`
- `GetMarketTradesEvents`：市场成交活动流（公开接口）
- Rewards：`GetCurrentRewards` / `GetRawRewardsForMarket` / `GetEarningsForUserForDay` 等
- RFQ：`CreateRfqRequest` / `CreateRfqQuote` / `GetRfqRequests` 等

接口实现对照：[Get Order](https://docs.polymarket.com/developers/CLOB/orders/get-order)、
[Get Active Orders](https://docs.polymarket.com/developers/CLOB/orders/get-active-order)、
[Check Order Reward Scoring](https://docs.polymarket.com/developers/CLOB/orders/check-scoring)

## 示例

示例代码位于 `examples/`，并已独立为子模块：

- `examples/wss_orderbook_by_event`：输入 `EVENT_SLUG`，实时拉取订单簿
- `examples/clob_place_order`：下单示例
- `examples/clob_trades`：交易数据示例
- `examples/clob_balance_allowance`：余额与授权示例
- `examples/clob_open_orders`：账户活跃订单示例
- `examples/clob_order_history`：通过 trades 反推历史订单示例
- `examples/relayer_deploy_safe`：部署 Safe 示例
- `examples/relayer_redeem_tokens`：赎回代币（redeem positions）示例

示例使用 `github.com/joho/godotenv` 读取 `.env`。

## 配置说明

`Config` 常用字段：

- `BaseURL` / `CLOBBaseURL` / `WSSMarketURL` / `WSSUserURL` / `RTDSURL`
- `Address` / `PrivateKey` / `APIKey` / `APISecret` / `Passphrase`
- `SignatureType` / `Funder` / `ChainID`
- （可选）builder flow：`BuilderAPIKey` / `BuilderAPISecret` / `BuilderPassphrase`

## 目录结构

- `client.go`：SDK 聚合入口
- `rest.go`：REST 客户端
- `clob*.go`：CLOB 相关能力
- `wss.go` / `rtds.go`：实时与订阅
- `wallet_client.go` / `relayer_client.go`：钱包与 relayer

## 说明

- 下单测试建议参考 `clob_order_creation_test.go` 中的流程与参数。
- 真实交易请确保账户余额与授权充足。
