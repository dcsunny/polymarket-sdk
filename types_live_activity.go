// types_live_activity.go 模块
package polymarket

// MarketTradeEvent 表示市场实时活动（成交）事件。
// 对照 Node SDK：MarketTradeEvent
type MarketTradeEvent struct {
	EventType string `json:"event_type"`

	Market struct {
		ConditionID string `json:"condition_id"`
		AssetID     string `json:"asset_id"`
		Question    string `json:"question"`
		Icon        string `json:"icon"`
		Slug        string `json:"slug"`
	} `json:"market"`

	User struct {
		Address                 string `json:"address"`
		Username                string `json:"username"`
		ProfilePicture          string `json:"profile_picture"`
		OptimizedProfilePicture string `json:"optimized_profile_picture"`
		Pseudonym               string `json:"pseudonym"`
	} `json:"user"`

	Side            string `json:"side"`
	Size            string `json:"size"`
	FeeRateBps      string `json:"fee_rate_bps"`
	Price           string `json:"price"`
	Outcome         string `json:"outcome"`
	OutcomeIndex    int    `json:"outcome_index"`
	TransactionHash string `json:"transaction_hash"`
	Timestamp       string `json:"timestamp"`
}
