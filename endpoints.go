// endpoints.go 模块
package polymarket

// CLOB API 路径常量（对齐 Node clob-client：`src/endpoints.ts`）。
//
// 说明：
// - 这里仅包含「path」部分（以 `/` 开头），不包含 host/baseURL。
// - 部分端点在 Node 侧用不同常量名指向同一个 path（例如 create/delete），此处保留同样的可读性写法。
const (
	// Server Time（服务器时间）
	EndpointTime = "/time"

	// API Key endpoints（API Key 相关）
	EndpointCreateAPIKey = "/auth/api-key"
	EndpointGetAPIKeys   = "/auth/api-keys"
	EndpointDeleteAPIKey = "/auth/api-key"
	EndpointDeriveAPIKey = "/auth/derive-api-key"
	EndpointClosedOnly   = "/auth/ban-status/closed-only"

	// Readonly API Key endpoints（只读 API Key）
	EndpointCreateReadonlyAPIKey   = "/auth/readonly-api-key"
	EndpointGetReadonlyAPIKeys     = "/auth/readonly-api-keys"
	EndpointDeleteReadonlyAPIKey   = "/auth/readonly-api-key"
	EndpointValidateReadonlyAPIKey = "/auth/validate-readonly-api-key"

	// Builder API Key endpoints（Builder API Key）
	EndpointCreateBuilderAPIKey = "/auth/builder-api-key"
	EndpointGetBuilderAPIKeys   = "/auth/builder-api-key"
	EndpointRevokeBuilderAPIKey = "/auth/builder-api-key"

	// Markets（市场）
	EndpointGetSamplingSimplifiedMarkets = "/sampling-simplified-markets"
	EndpointGetSamplingMarkets           = "/sampling-markets"
	EndpointGetSimplifiedMarkets         = "/simplified-markets"
	EndpointGetMarkets                   = "/markets"
	EndpointGetMarketPrefix              = "/markets/"
	EndpointGetOrderBook                 = "/book"
	EndpointGetOrderBooks                = "/books"
	EndpointGetMidpoint                  = "/midpoint"
	EndpointGetMidpoints                 = "/midpoints"
	EndpointGetPrice                     = "/price"
	EndpointGetPrices                    = "/prices"
	EndpointGetSpread                    = "/spread"
	EndpointGetSpreads                   = "/spreads"
	EndpointGetLastTradePrice            = "/last-trade-price"
	EndpointGetLastTradesPrices          = "/last-trades-prices"
	EndpointGetTickSize                  = "/tick-size"
	EndpointGetNegRisk                   = "/neg-risk"
	EndpointGetFeeRate                   = "/fee-rate"

	// Order endpoints（订单）
	EndpointPostOrder          = "/order"
	EndpointPostOrders         = "/orders"
	EndpointCancelOrder        = "/order"
	EndpointCancelOrders       = "/orders"
	EndpointGetOrderPrefix     = "/data/order/"
	EndpointCancelAll          = "/cancel-all"
	EndpointCancelMarketOrders = "/cancel-market-orders"
	EndpointGetOpenOrders      = "/data/orders"
	EndpointGetTrades          = "/data/trades"
	EndpointIsOrderScoring     = "/order-scoring"
	EndpointAreOrdersScoring   = "/orders-scoring"

	// Price history（历史价格）
	EndpointGetPricesHistory = "/prices-history"

	// Notifications（通知）
	EndpointGetNotifications  = "/notifications"
	EndpointDropNotifications = "/notifications"

	// Balance（余额与授权）
	EndpointGetBalanceAllowance    = "/balance-allowance"
	EndpointUpdateBalanceAllowance = "/balance-allowance/update"

	// Live activity（成交活动流）
	EndpointGetMarketTradesEventsPrefix = "/live-activity/events/"

	// Rewards（奖励）
	EndpointGetEarningsForUserForDay      = "/rewards/user"
	EndpointGetTotalEarningsForUserForDay = "/rewards/user/total"
	EndpointGetLiquidityRewardPercentages = "/rewards/user/percentages"
	EndpointGetRewardsMarketsCurrent      = "/rewards/markets/current"
	EndpointGetRewardsMarketsPrefix       = "/rewards/markets/"
	EndpointGetRewardsEarningsPercentages = "/rewards/user/markets"

	// Builder endpoints（Builder）
	EndpointGetBuilderTrades = "/builder/trades"

	// Heartbeats（心跳）
	EndpointPostHeartbeat = "/v1/heartbeats"

	// RFQ
	EndpointCreateRfqRequest      = "/rfq/request"
	EndpointCancelRfqRequest      = "/rfq/request"
	EndpointGetRfqRequests        = "/rfq/data/requests"
	EndpointCreateRfqQuote        = "/rfq/quote"
	EndpointCancelRfqQuote        = "/rfq/quote"
	EndpointRfqRequestsAccept     = "/rfq/request/accept"
	EndpointRfqQuoteApprove       = "/rfq/quote/approve"
	EndpointGetRfqRequesterQuotes = "/rfq/data/requester/quotes"
	EndpointGetRfqQuoterQuotes    = "/rfq/data/quoter/quotes"
	EndpointGetRfqBestQuote       = "/rfq/data/best-quote"
	EndpointRfqConfig             = "/rfq/config"
)
