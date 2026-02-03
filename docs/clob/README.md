# CLOB

本 SDK 的 `CLOBClient` 覆盖 Polymarket CLOB REST API（下单/撤单/订单簿/成交等）。

## 初始化

```go
sdk, _ := pm.New(pm.Config{
    Address:    "0x...",
    APIKey:     "...",
    APISecret:  "...",
    Passphrase: "...",
})
```

## 订单管理

- `CreateOrder`：构建并签名订单（需要 `PrivateKey`/`Address`）
- `PostOrder` / `PostOrderWithOptions`：提交单个订单（L2 认证）
- `PostOrders` / `PostOrdersSigned`：批量提交订单（L2 认证，最多 15）
- `GetOrder`：获取单个订单（L2 认证）
- `GetActiveOrders` / `GetActiveOrdersPage`：获取活跃订单（L2 认证）
- `CancelOrder` / `CancelOrders` / `CancelAllOrders` / `CancelMarketOrders`：撤单（L2 认证）
- `IsOrderScoring` / `AreOrdersScoring`：奖励评分状态（L2 认证）

## 订单簿与成交

- `GetOrderBook` / `GetOrderBooks`：订单簿快照
- `GetTrades` / `GetTradesPage`：成交列表（L2 认证）

## 市场数据

- `GetMidpoint` / `GetMidpoints`
- `GetPrices`
- `GetSpread` / `GetSpreads`
- `GetLastTradePrice` / `GetLastTradesPrices`
- `GetTickSize` / `GetNegRisk` / `GetFeeRateBps`

## 通知、余额、心跳

- `GetNotifications` / `DropNotifications`（L2 认证）
- `GetBalanceAllowance` / `UpdateBalanceAllowance`（L2 认证）
- `PostHeartbeat`（L2 认证）

## Live Activity

- `GetMarketTradesEvents`：市场成交活动流（公开接口）

## Rewards

- `GetCurrentRewards` / `GetRawRewardsForMarket`（公开接口）
- `GetEarningsForUserForDay` / `GetTotalEarningsForUserForDay` / `GetUserEarningsAndMarketsConfig` / `GetRewardPercentages`（L2 认证）

## RFQ

- `CreateRfqRequest` / `CancelRfqRequest` / `GetRfqRequests`（L2 认证）
- `CreateRfqQuote` / `CancelRfqQuote` / `GetRfqRequesterQuotes` / `GetRfqQuoterQuotes` / `GetRfqBestQuote`（L2 认证）
- `GetRfqConfig`（L2 认证）
- `AcceptRfqQuote` / `ApproveRfqOrder`（L2 认证，payload 透传）

