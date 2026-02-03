// types_rest.go 模块
package polymarket

import (
	"encoding/json"
	"strconv"
	"time"
)

// Event represents a Polymarket event.
type Event struct {
	ID     int64  `json:"-"`
	IDRaw  string `json:"id"`
	Ticker string `json:"ticker"`
	Slug   string `json:"slug"`
	Title  string `json:"title"`

	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`

	Image          string                 `json:"image"`
	Icon           string                 `json:"icon"`
	FeaturedImage  string                 `json:"featuredImage"`
	ImageOptimized map[string]interface{} `json:"imageOptimized"`
	IconOptimized  map[string]interface{} `json:"iconOptimized"`

	StartDate         *time.Time `json:"startDate"`
	EndDate           *time.Time `json:"endDate"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	CreationDate      *time.Time `json:"creationDate"`
	ClosedTime        *time.Time `json:"closedTime"`
	FinishedTimestamp *time.Time `json:"finishedTimestamp"`

	Active                bool `json:"active"`
	Closed                bool `json:"closed"`
	Archived              bool `json:"archived"`
	New                   bool `json:"new"`
	Featured              bool `json:"featured"`
	Restricted            bool `json:"restricted"`
	Live                  bool `json:"live"`
	Ended                 bool `json:"ended"`
	AutomaticallyResolved bool `json:"automaticallyResolved"`
	AutomaticallyActive   bool `json:"automaticallyActive"`
	PendingDeployment     bool `json:"pendingDeployment"`
	Deploying             bool `json:"deploying"`
	CYOM                  bool `json:"cyom"`

	Liquidity     float64 `json:"liquidity"`
	Volume        float64 `json:"volume"`
	OpenInterest  float64 `json:"openInterest"`
	LiquidityAmm  float64 `json:"liquidityAmm"`
	LiquidityClob float64 `json:"liquidityClob"`
	Volume24hr    float64 `json:"volume24hr"`
	Volume1wk     float64 `json:"volume1wk"`
	Volume1mo     float64 `json:"volume1mo"`
	Volume1yr     float64 `json:"volume1yr"`

	NegRisk          bool   `json:"negRisk"`
	NegRiskMarketID  string `json:"negRiskMarketID"`
	NegRiskFeeBips   int    `json:"negRiskFeeBips"`
	NegRiskAugmented bool   `json:"negRiskAugmented"`

	CommentCount int           `json:"commentCount"`
	TweetCount   int           `json:"tweetCount"`
	Competitive  FloatOrString `json:"competitive"`
	Score        FloatOrString `json:"score"`
	Recurrence   string        `json:"recurrence"`

	Markets       []interface{} `json:"markets"`
	Series        []interface{} `json:"series"`
	Categories    []interface{} `json:"categories"`
	Collections   []interface{} `json:"collections"`
	Tags          []interface{} `json:"tags"`
	EventCreators []interface{} `json:"eventCreators"`
	Chats         []interface{} `json:"chats"`
	Templates     []interface{} `json:"templates"`
	SubEvents     []interface{} `json:"subEvents"`
}

// FloatOrString handles fields that can be string or number.
type FloatOrString float64

func (f *FloatOrString) Float64() float64 {
	return float64(*f)
}

func (f *FloatOrString) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		if value == "" {
			*f = 0
			return nil
		}
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			*f = 0
			return nil
		}
		*f = FloatOrString(parsed)
	case float64:
		*f = FloatOrString(value)
	case int:
		*f = FloatOrString(value)
	case int64:
		*f = FloatOrString(value)
	default:
		*f = 0
	}
	return nil
}

func (f FloatOrString) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}

// Market represents a Polymarket market.
type Market struct {
	ID          string `json:"id"`
	ConditionID string `json:"conditionId"`
	Slug        string `json:"slug"`
	Question    string `json:"question"`

	Description string `json:"description"`
	Category    string `json:"category"`

	Image          string                 `json:"image"`
	Icon           string                 `json:"icon"`
	ImageOptimized map[string]interface{} `json:"imageOptimized"`
	IconOptimized  map[string]interface{} `json:"iconOptimized"`

	StartDate                *time.Time `json:"startDate"`
	EndDate                  *time.Time `json:"endDate"`
	EndDateIso               string     `json:"endDateIso"`
	EventStartTime           *time.Time `json:"eventStartTime"`
	CreatedAt                *time.Time `json:"createdAt"`
	UpdatedAt                *time.Time `json:"updatedAt"`
	AcceptingOrdersTimestamp *time.Time `json:"acceptingOrdersTimestamp"`

	Active                       bool `json:"active"`
	Closed                       bool `json:"closed"`
	Archived                     bool `json:"archived"`
	New                          bool `json:"new"`
	Featured                     bool `json:"featured"`
	Restricted                   bool `json:"restricted"`
	Approved                     bool `json:"approved"`
	AcceptingOrders              bool `json:"acceptingOrders"`
	AutomaticallyActive          bool `json:"automaticallyActive"`
	ClearBookOnStart             bool `json:"clearBookOnStart"`
	ManualActivation             bool `json:"manualActivation"`
	PendingDeployment            bool `json:"pendingDeployment"`
	Deploying                    bool `json:"deploying"`
	EnableOrderBook              bool `json:"enableOrderBook"`
	FeesEnabled                  bool `json:"feesEnabled"`
	Funded                       bool `json:"funded"`
	HasReviewedDates             bool `json:"hasReviewedDates"`
	HoldingRewardsEnabled        bool `json:"holdingRewardsEnabled"`
	Ready                        bool `json:"ready"`
	RfqEnabled                   bool `json:"rfqEnabled"`
	ShowGmpOutcome               bool `json:"showGmpOutcome"`
	ShowGmpSeries                bool `json:"showGmpSeries"`
	PagerDutyNotificationEnabled bool `json:"pagerDutyNotificationEnabled"`

	Liquidity      FloatOrString `json:"liquidity"`
	LiquidityNum   FloatOrString `json:"liquidityNum"`
	Volume         FloatOrString `json:"volume"`
	VolumeNum      FloatOrString `json:"volumeNum"`
	LiquidityAmm   float64       `json:"liquidityAmm"`
	LiquidityClob  float64       `json:"liquidityClob"`
	Volume24hr     float64       `json:"volume24hr"`
	Volume24hrAmm  float64       `json:"volume24hrAmm"`
	Volume24hrClob float64       `json:"volume24hrClob"`
	Volume1wk      float64       `json:"volume1wk"`
	Volume1wkAmm   float64       `json:"volume1wkAmm"`
	Volume1wkClob  float64       `json:"volume1wkClob"`
	Volume1mo      float64       `json:"volume1mo"`
	Volume1moAmm   float64       `json:"volume1moAmm"`
	Volume1moClob  float64       `json:"volume1moClob"`
	Volume1yr      float64       `json:"volume1yr"`
	Volume1yrAmm   float64       `json:"volume1yrAmm"`
	Volume1yrClob  float64       `json:"volume1yrClob"`
	VolumeAmm      float64       `json:"volumeAmm"`
	VolumeClob     float64       `json:"volumeClob"`

	BestAsk             float64 `json:"bestAsk"`
	BestBid             float64 `json:"bestBid"`
	LastTradePrice      float64 `json:"lastTradePrice"`
	OneDayPriceChange   float64 `json:"oneDayPriceChange"`
	OneHourPriceChange  float64 `json:"oneHourPriceChange"`
	OneMonthPriceChange float64 `json:"oneMonthPriceChange"`
	OneWeekPriceChange  float64 `json:"oneWeekPriceChange"`
	OneYearPriceChange  float64 `json:"oneYearPriceChange"`
	OutcomePrices       string  `json:"outcomePrices"`

	Competitive           interface{} `json:"competitive"`
	GroupItemThreshold    interface{} `json:"groupItemThreshold"`
	OrderMinSize          float64     `json:"orderMinSize"`
	OrderPriceMinTickSize float64     `json:"orderPriceMinTickSize"`
	RewardsMaxSpread      float64     `json:"rewardsMaxSpread"`
	RewardsMinSize        float64     `json:"rewardsMinSize"`
	Spread                float64     `json:"spread"`

	NegRisk      bool `json:"negRisk"`
	NegRiskOther bool `json:"negRiskOther"`

	Outcomes           string `json:"outcomes"`
	ClobTokenIds       string `json:"clobTokenIds"`
	QuestionID         string `json:"questionID"`
	ResolutionSource   string `json:"resolutionSource"`
	MarketMakerAddress string `json:"marketMakerAddress"`

	Events                []MarketEvent `json:"events"`
	Categories            []interface{} `json:"categories"`
	Tags                  []interface{} `json:"tags"`
	UmaResolutionStatuses interface{}   `json:"umaResolutionStatuses"`
}

// MarketEvent represents event data attached to a market.
type MarketEvent struct {
	Id                string    `json:"id"`
	Ticker            string    `json:"ticker"`
	Slug              string    `json:"slug"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"startDate"`
	CreationDate      time.Time `json:"creationDate"`
	EndDate           time.Time `json:"endDate"`
	Image             string    `json:"image"`
	Icon              string    `json:"icon"`
	Active            bool      `json:"active"`
	Closed            bool      `json:"closed"`
	Archived          bool      `json:"archived"`
	Featured          bool      `json:"featured"`
	Restricted        bool      `json:"restricted"`
	Liquidity         float64   `json:"liquidity"`
	Volume            float64   `json:"volume"`
	OpenInterest      float64   `json:"openInterest"`
	SortBy            string    `json:"sortBy"`
	Category          string    `json:"category"`
	PublishedAt       string    `json:"published_at"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Competitive       float64   `json:"competitive"`
	Volume24Hr        float64   `json:"volume24hr"`
	Volume1Wk         float64   `json:"volume1wk"`
	Volume1Mo         float64   `json:"volume1mo"`
	Volume1Yr         float64   `json:"volume1yr"`
	LiquidityAmm      float64   `json:"liquidityAmm"`
	LiquidityClob     float64   `json:"liquidityClob"`
	CommentCount      int       `json:"commentCount"`
	Cyom              bool      `json:"cyom"`
	ClosedTime        time.Time `json:"closedTime"`
	ShowAllOutcomes   bool      `json:"showAllOutcomes"`
	ShowMarketImages  bool      `json:"showMarketImages"`
	EnableNegRisk     bool      `json:"enableNegRisk"`
	NegRiskAugmented  bool      `json:"negRiskAugmented"`
	PendingDeployment bool      `json:"pendingDeployment"`
	Deploying         bool      `json:"deploying"`
}

// UnmarshalJSON handles Event ID that can be string or number.
func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		ID interface{} `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.ID.(type) {
	case string:
		e.IDRaw = v
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			e.ID = id
		}
	case float64:
		e.ID = int64(v)
		e.IDRaw = strconv.FormatInt(int64(v), 10)
	case int64:
		e.ID = v
		e.IDRaw = strconv.FormatInt(v, 10)
	}

	return nil
}
