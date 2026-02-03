// types_rewards.go 模块
package polymarket

// UserEarning 单日用户收益明细（对照 Node SDK：UserEarning）。
type UserEarning struct {
	Date         string  `json:"date"`
	ConditionID  string  `json:"condition_id"`
	AssetAddress string  `json:"asset_address"`
	MakerAddress string  `json:"maker_address"`
	Earnings     float64 `json:"earnings"`
	AssetRate    float64 `json:"asset_rate"`
}

// TotalUserEarning 单日用户总收益（对照 Node SDK：TotalUserEarning）。
type TotalUserEarning struct {
	Date         string  `json:"date"`
	AssetAddress string  `json:"asset_address"`
	MakerAddress string  `json:"maker_address"`
	Earnings     float64 `json:"earnings"`
	AssetRate    float64 `json:"asset_rate"`
}

// RewardsPercentages 用户在各市场的奖励占比（对照 Node SDK：RewardsPercentages）。
type RewardsPercentages map[string]float64

// RewardToken 奖励市场中的 token 信息。
type RewardToken struct {
	TokenID string  `json:"token_id"`
	Outcome string  `json:"outcome"`
	Price   float64 `json:"price"`
}

// RewardsConfig 单个奖励配置项（对照 Node SDK：RewardsConfig）。
type RewardsConfig struct {
	AssetAddress string  `json:"asset_address"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	RatePerDay   float64 `json:"rate_per_day"`
	TotalRewards float64 `json:"total_rewards"`
}

// MarketReward 市场奖励信息（对照 Node SDK：MarketReward）。
type MarketReward struct {
	ConditionID      string          `json:"condition_id"`
	Question         string          `json:"question"`
	MarketSlug       string          `json:"market_slug"`
	EventSlug        string          `json:"event_slug"`
	Image            string          `json:"image"`
	RewardsMaxSpread float64         `json:"rewards_max_spread"`
	RewardsMinSize   float64         `json:"rewards_min_size"`
	Tokens           []RewardToken   `json:"tokens"`
	RewardsConfig    []RewardsConfig `json:"rewards_config"`
}

// Earning 收益项（对照 Node SDK：Earning）。
type Earning struct {
	AssetAddress string  `json:"asset_address"`
	Earnings     float64 `json:"earnings"`
	AssetRate    float64 `json:"asset_rate"`
}

// UserRewardsEarning 用户在奖励市场的收益与配置（对照 Node SDK：UserRewardsEarning）。
type UserRewardsEarning struct {
	ConditionID           string          `json:"condition_id"`
	Question              string          `json:"question"`
	MarketSlug            string          `json:"market_slug"`
	EventSlug             string          `json:"event_slug"`
	Image                 string          `json:"image"`
	RewardsMaxSpread      float64         `json:"rewards_max_spread"`
	RewardsMinSize        float64         `json:"rewards_min_size"`
	MarketCompetitiveness float64         `json:"market_competitiveness"`
	Tokens                []RewardToken   `json:"tokens"`
	RewardsConfig         []RewardsConfig `json:"rewards_config"`
	MakerAddress          string          `json:"maker_address"`
	EarningPercentage     float64         `json:"earning_percentage"`
	Earnings              []Earning       `json:"earnings"`
}
