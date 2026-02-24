// types_ws.go 模块
package polymarket

import (
	"encoding/json"
	"strconv"
)

// FlexInt 可以从 JSON 字符串或数字反序列化为 int64。
type FlexInt int64

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexInt(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*f = FlexInt(parsed)
	return nil
}

func (f FlexInt) Int64() int64 {
	return int64(f)
}

// =====================================================
// WebSocket 频道类型
// =====================================================

// WSSChannelType 表示 WebSocket 频道类型。
type WSSChannelType string

const (
	// WSSChannelTypeUser 用户频道（需认证，接收私有订单/交易更新）
	WSSChannelTypeUser WSSChannelType = "user"
	// WSSChannelTypeMarket 市场频道（公开，接收实时订单簿/价格/市场生命周期更新）
	WSSChannelTypeMarket WSSChannelType = "market"
)

// =====================================================
// WebSocket 事件类型常量
// 对应文档：
//   - Market Channel: https://docs.polymarket.com/api-reference/wss/market
//   - User Channel:   https://docs.polymarket.com/api-reference/wss/user
// =====================================================

// WSSEventType 表示 WebSocket 消息的事件类型。
type WSSEventType = string

const (
	// ----- Market Channel 事件类型 -----

	// WSSEventBook 订单簿快照（market channel 接收）
	WSSEventBook WSSEventType = "book"
	// WSSEventPriceChange 价格变动（market channel 接收）
	WSSEventPriceChange WSSEventType = "price_change"
	// WSSEventLastTradePrice 最后成交价更新（market channel 接收）
	WSSEventLastTradePrice WSSEventType = "last_trade_price"
	// WSSEventTickSizeChange 最小价格刻度变更（market channel 接收）
	WSSEventTickSizeChange WSSEventType = "tick_size_change"
	// WSSEventBestBidAsk 最优买卖价更新（market channel 接收）
	WSSEventBestBidAsk WSSEventType = "best_bid_ask"
	// WSSEventNewMarket 新市场创建通知（market channel 接收）
	WSSEventNewMarket WSSEventType = "new_market"
	// WSSEventMarketResolved 市场结算通知（market channel 接收）
	WSSEventMarketResolved WSSEventType = "market_resolved"

	// ----- User Channel 事件类型 -----

	// WSSEventOrder 订单事件（user channel 接收）
	WSSEventOrder WSSEventType = "order"
	// WSSEventTrade 成交事件（user channel 接收）
	WSSEventTrade WSSEventType = "trade"
)

// =====================================================
// WebSocket 订阅方法常量
// =====================================================

const (
	// WSSMethodSubscribe 订阅资产（用于 WSSSubscriptionUpdate）
	WSSMethodSubscribe = "subscribe"
	// WSSMethodUnsubscribe 取消订阅资产（用于 WSSSubscriptionUpdate）
	WSSMethodUnsubscribe = "unsubscribe"
)

// =====================================================
// WebSocket 发送消息（客户端 → 服务端）
// =====================================================

// WSSAuth 用户频道认证信息。
type WSSAuth struct {
	// API 密钥
	APIKey string `json:"apiKey"`
	// API 密钥对应的密钥
	Secret string `json:"secret"`
	// API 密钥对应的口令
	Passphrase string `json:"passphrase"`
}

// WSSSubscription 初始订阅请求（连接后发送的第一条消息）。
type WSSSubscription struct {
	// 认证信息（仅 user channel 需要）
	Auth *WSSAuth `json:"auth,omitempty"`
	// 频道类型："market" 或 "user"
	Type WSSChannelType `json:"type"`
	// 要订阅的市场 ID 列表（user channel 使用）
	Markets []string `json:"markets,omitempty"`
	// 要订阅的 Asset ID（token ID）列表（market channel 使用）
	AssetsIDs []string `json:"assets_ids,omitempty"`
	// 是否在订阅时发送初始订单簿快照（默认 true）
	InitialDump *bool `json:"initial_dump,omitempty"`
	// 订阅级别（例如 level 2 为完整订单簿，默认 2）
	Level *int `json:"level,omitempty"`
}

// WSSSubscriptionUpdate 动态订阅/取消订阅请求（无需重连即可修改订阅）。
type WSSSubscriptionUpdate struct {
	// 方法："subscribe" 或 "unsubscribe"
	Method string `json:"method"`
	// 要添加或移除的 Asset ID 列表（market channel）
	AssetsIDs []string `json:"assets_ids,omitempty"`
	// 要添加或移除的市场 ID 列表（user channel）
	Markets []string `json:"markets,omitempty"`
	// 订阅级别
	Level *int `json:"level,omitempty"`
}

// =====================================================
// WebSocket 接收消息（服务端 → 客户端）
// =====================================================

// WSSMessageHandler 处理 WebSocket 原始消息的回调函数。
type WSSMessageHandler func(data json.RawMessage) error

// ----- Market Channel 消息 -----

// WSSBookMessage 订单簿快照消息。
//
// 事件类型：book
// 频道：market channel
// 在订阅时发送完整的订单簿快照，后续通过 price_change 增量更新。
type WSSBookMessage struct {
	// 事件类型，固定为 "book"
	EventType string `json:"event_type"`
	// 市场条件 ID
	Market string `json:"market"`
	// 资产 ID（token ID）
	AssetID string `json:"asset_id"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
	// 订单簿哈希，用于一致性校验
	Hash string `json:"hash"`
	// 买单列表（按价格排序）
	Bids []WSSOrderSummary `json:"bids"`
	// 卖单列表（按价格排序）
	Asks []WSSOrderSummary `json:"asks"`
}

// WSSOrderSummary 订单簿中的单个价格层级。
type WSSOrderSummary struct {
	// 价格
	Price string `json:"price"`
	// 该价格上的总数量
	Size string `json:"size"`
}

// WSSPriceChangeMessage 价格变动消息（订单簿增量更新）。
//
// 事件类型：price_change
// 频道：market channel
// 当订单簿上的买卖价格/数量发生变化时推送。
type WSSPriceChangeMessage struct {
	// 事件类型，固定为 "price_change"
	EventType string `json:"event_type"`
	// 市场条件 ID
	Market string `json:"market"`
	// 资产 ID
	AssetID string `json:"asset_id"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
	// 价格层级变化列表
	Changes []PriceLevelChange `json:"changes"`
}

// PriceLevelChange 单个价格层级的变化。
type PriceLevelChange struct {
	// 方向："BUY" 或 "SELL"
	Side string `json:"side"`
	// 变化的价格层级
	Price string `json:"price"`
	// 该价格上的新总数量（为 "0" 表示该层级被移除）
	Size string `json:"size"`
	// 变化后的最优买价
	BestBid string `json:"best_bid"`
	// 变化后的最优卖价
	BestAsk string `json:"best_ask"`
	// 变化后的订单簿哈希
	Hash string `json:"hash"`
}

// WSSLastTradePriceMessage 最后成交价消息。
//
// 事件类型：last_trade_price
// 频道：market channel
// 当某个资产发生新的成交时推送。
type WSSLastTradePriceMessage struct {
	// 事件类型，固定为 "last_trade_price"
	EventType string `json:"event_type"`
	// 市场条件 ID
	Market string `json:"market"`
	// 资产 ID
	AssetID string `json:"asset_id"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
	// 最后成交价
	Price string `json:"price"`
	// 最后成交数量
	Size string `json:"size"`
	// 成交方向："BUY" 或 "SELL"
	Side string `json:"side"`
	// 手续费率
	FeeRate string `json:"fee_rate"`
}

// WSSTickSizeChangeMessage 最小价格刻度变更消息。
//
// 事件类型：tick_size_change
// 频道：market channel
// 当某个资产的最小价格刻度（tick size）发生变化时推送。
type WSSTickSizeChangeMessage struct {
	// 事件类型，固定为 "tick_size_change"
	EventType string `json:"event_type"`
	// 市场条件 ID
	Market string `json:"market"`
	// 资产 ID
	AssetID string `json:"asset_id"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
	// 变更前的最小价格刻度
	OldTickSize string `json:"old_tick_size"`
	// 变更后的最小价格刻度
	NewTickSize string `json:"new_tick_size"`
}

// WSSBestBidAskMessage 最优买卖价消息。
//
// 事件类型：best_bid_ask
// 频道：market channel
// 当某个资产的最优买/卖价发生变化时推送。
type WSSBestBidAskMessage struct {
	// 事件类型，固定为 "best_bid_ask"
	EventType string `json:"event_type"`
	// 市场条件 ID
	Market string `json:"market"`
	// 资产 ID
	AssetID string `json:"asset_id"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
	// 最优买价
	BestBid string `json:"best_bid"`
	// 最优卖价
	BestAsk string `json:"best_ask"`
	// 买卖价差
	Spread string `json:"spread"`
}

// WSSNewMarketMessage 新市场创建通知。
//
// 事件类型：new_market
// 频道：market channel
// 当平台上创建新市场时推送。
type WSSNewMarketMessage struct {
	// 事件类型，固定为 "new_market"
	EventType string `json:"event_type"`
	// 市场 ID
	ID string `json:"id"`
	// 市场条件地址
	Market string `json:"market"`
	// 市场问题
	Question string `json:"question"`
	// 市场 slug（URL 友好标识）
	Slug string `json:"slug"`
	// 资产 ID 列表
	AssetsIDs []string `json:"assets_ids"`
	// 结果选项列表
	Outcomes []string `json:"outcomes"`
	// 标签列表
	Tags []string `json:"tags"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
}

// WSSMarketResolvedMessage 市场结算通知。
//
// 事件类型：market_resolved
// 频道：market channel
// 当某个市场结算完成时推送，包含胜出资产的信息。
type WSSMarketResolvedMessage struct {
	// 事件类型，固定为 "market_resolved"
	EventType string `json:"event_type"`
	// 市场 ID
	ID string `json:"id"`
	// 市场条件地址
	Market string `json:"market"`
	// 该市场所有资产 ID 列表
	AssetsIDs []string `json:"assets_ids"`
	// 胜出的资产 ID
	WinnerAssetID string `json:"winner_asset_id"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
}

// ----- User Channel 消息 -----

// WSSOrderEvent 订单事件消息。
//
// 事件类型：order
// 频道：user channel（需认证）
// 当用户的订单状态发生变化时推送（下单、成交、取消等）。
type WSSOrderEvent struct {
	// 事件类型，固定为 "order"
	EventType string `json:"event_type"`
	// 订单 ID
	ID string `json:"id"`
	// 市场条件 ID
	Market string `json:"market"`
	// 资产 ID
	AssetID string `json:"asset_id"`
	// 订单所有者钱包地址
	OrderOwner string `json:"order_owner"`
	// 订单价格
	Price string `json:"price"`
	// 方向："BUY" 或 "SELL"
	Side string `json:"side"`
	// 原始下单数量
	OriginalSize string `json:"original_size"`
	// 已成交数量
	SizeMatched string `json:"size_matched"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
	// 订单变更类型："PLACEMENT"（下单）、"CANCELLATION"（取消）、"MATCH"（成交）等
	Type string `json:"type"`
	// 结果标签："YES" 或 "NO"
	Outcome string `json:"outcome,omitempty"`
	// 订单状态："OPEN"、"MATCHED"、"CANCELED" 等
	Status string `json:"status,omitempty"`
	// 关联的交易 ID 列表
	AssociatedTrades []string `json:"associated_trades,omitempty"`
}

// WSSTradeEvent 成交事件消息。
//
// 事件类型：trade
// 频道：user channel（需认证）
// 当用户的订单发生成交时推送。
type WSSTradeEvent struct {
	// 事件类型，固定为 "trade"
	EventType string `json:"event_type"`
	// 成交 ID
	ID string `json:"id"`
	// 市场条件 ID
	Market string `json:"market"`
	// 资产 ID
	AssetID string `json:"asset_id"`
	// 成交所有者钱包地址
	Owner string `json:"owner"`
	// 成交价格
	Price string `json:"price"`
	// 方向："BUY" 或 "SELL"
	Side string `json:"side"`
	// 成交数量
	Size string `json:"size"`
	// 成交状态："MATCHED"
	Status string `json:"status"`
	// 时间戳
	Timestamp FlexInt `json:"timestamp"`
	// taker 订单 ID
	TakerOrderID string `json:"taker_order_id"`
	// 交易方身份："TAKER" 或 "MAKER"
	TraderSide string `json:"trader_side,omitempty"`
	// 链上交易哈希
	TransactionHash string `json:"transaction_hash,omitempty"`
	// Maker 订单明细列表
	MakerOrders []MakerOrder `json:"maker_orders"`
	// 原始数据（供调用方自行解析）
	RawData json.RawMessage `json:"-"`
}

// ----- Sports Channel 消息 -----

// WSSSportsResultMessage 体育赛事实时结果消息。
//
// 频道：sports channel（wss://ws-subscriptions-clob.polymarket.com/ws/sports）
// 注意：Sports channel 使用 "ping"/"pong" 心跳（每 5 秒），且消息没有 event_type 字段。
type WSSSportsResultMessage struct {
	// 体育赛事标识（slug）
	Slug string `json:"slug"`
	// 比赛是否正在进行中
	Live bool `json:"live"`
	// 比赛是否已结束
	Ended bool `json:"ended"`
	// 比分信息
	Score WSSSportsScore `json:"score"`
	// 当前阶段信息（如 "Q1"、"2nd Half"）
	Period string `json:"period"`
	// 当前阶段已用时间
	Elapsed string `json:"elapsed"`
	// 数据最后更新时间（ISO 8601）
	LastUpdate string `json:"last_update"`
	// 比赛结束时间（ISO 8601，仅比赛结束后有值）
	FinishedTimestamp string `json:"finished_timestamp,omitempty"`
	// 当前控球/回合方（部分赛事适用）
	Turn string `json:"turn,omitempty"`
}

// WSSSportsScore 体育赛事比分。
type WSSSportsScore struct {
	// 主队得分
	Home string `json:"home"`
	// 客队得分
	Away string `json:"away"`
}

// =====================================================
// WebSocket 订单变更类型常量（User Channel order 事件的 type 字段）
// =====================================================

const (
	// WSSOrderTypePlacement 下单
	WSSOrderTypePlacement = "PLACEMENT"
	// WSSOrderTypeCancellation 取消
	WSSOrderTypeCancellation = "CANCELLATION"
	// WSSOrderTypeMatch 成交
	WSSOrderTypeMatch = "MATCH"
)

// =====================================================
// WebSocket 订单状态常量（User Channel order 事件的 status 字段）
// =====================================================

const (
	// WSSOrderStatusOpen 订单挂单中
	WSSOrderStatusOpen = "OPEN"
	// WSSOrderStatusMatched 订单已完全成交
	WSSOrderStatusMatched = "MATCHED"
	// WSSOrderStatusCanceled 订单已取消
	WSSOrderStatusCanceled = "CANCELED"
)
