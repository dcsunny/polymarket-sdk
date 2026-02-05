// types_clob.go 模块
package polymarket

// APICredentials holds API key info.
type APICredentials struct {
	APIKey     string `json:"apiKey"`
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}

// APIKeysResponse is returned by GetAPIKeys.
type APIKeysResponse struct {
	APIKeys []string `json:"apiKeys"`
}

// GetActiveOrdersRequest filters active orders.
type GetActiveOrdersRequest struct {
	ID      string
	Market  string
	AssetID string
}

// GetActiveOrdersResponse represents paginated response.
type GetActiveOrdersResponse struct {
	Count      int          `json:"count"`
	Data       []*OpenOrder `json:"data"`
	Limit      int          `json:"limit"`
	NextCursor string       `json:"next_cursor"`
}

const (
	// EndCursor 表示分页结束（与 Node SDK 对齐）。
	EndCursor = "LTE="
)

// OpenOrder represents an active order.
type OpenOrder struct {
	ID              string      `json:"id"`
	Market          string      `json:"market"`
	AssetID         string      `json:"asset_id"`
	Price           string      `json:"price"`
	Size            string      `json:"size"`
	Side            string      `json:"side"`
	OrderType       string      `json:"type"`
	Status          string      `json:"status"`
	Owner           string      `json:"owner"`
	MakerAddress    string      `json:"maker_address"`
	FilledSize      string      `json:"filled_size"`
	RemainingSize   string      `json:"remaining_size"`
	OriginalSize    string      `json:"original_size"`
	SizeMatched     string      `json:"size_matched"`
	Outcome         string      `json:"outcome"`
	CreatedAt       interface{} `json:"created_at"`
	UpdatedAt       interface{} `json:"updated_at"`
	AssociateTrades []string    `json:"associate_trades"`
}

// MakerOrder represents a maker order in a trade.
type MakerOrder struct {
	OrderID         string `json:"order_id"`
	Price           string `json:"price"`
	Size            string `json:"size"`
	MatchedSize     string `json:"matched_size"`
	Outcome         string `json:"outcome"`
	OwnerAddress    string `json:"owner_address"`
	FeeRateBps      string `json:"fee_rate_bps"`
	AssetID         string `json:"asset_id"`
	MarketID        string `json:"market_id"`
	TransactionHash string `json:"transaction_hash"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
}

// PriceSide indicates BUY or SELL for price.
type PriceSide string

const (
	PriceSideBuy  PriceSide = "BUY"
	PriceSideSell PriceSide = "SELL"
)

func (s PriceSide) String() string {
	return string(s)
}

// PriceResponse represents price response.
type PriceResponse struct {
	Price string `json:"price"`
}

// CancelMarketOrdersRequest cancels orders by market or asset.
type CancelMarketOrdersRequest struct {
	Market  string `json:"market,omitempty"`
	AssetID string `json:"asset_id,omitempty"`
}

// OrderBookSummary 订单簿快照。
type OrderBookSummary struct {
	Market         string         `json:"market"`
	AssetID        string         `json:"asset_id"`
	Timestamp      string         `json:"timestamp"`
	Bids           []OrderSummary `json:"bids"`
	Asks           []OrderSummary `json:"asks"`
	MinOrderSize   string         `json:"min_order_size"`
	TickSize       string         `json:"tick_size"`
	NegRisk        bool           `json:"neg_risk"`
	LastTradePrice string         `json:"last_trade_price"`
	Hash           string         `json:"hash"`
}

// OrderSummary 订单簿档位。
type OrderSummary struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// BookParams 批量订单簿参数。
type BookParams struct {
	TokenID string `json:"token_id"`
	Side    string `json:"side,omitempty"`
}

// OrderScoringParams 单个订单评分参数。
type OrderScoringParams struct {
	OrderID string `json:"order_id"`
}

// OrderScoring 单个订单评分结果。
type OrderScoring struct {
	Scoring bool `json:"scoring"`
}

// OrdersScoringParams 批量订单评分参数。
type OrdersScoringParams struct {
	OrderIDs []string `json:"orderIds"`
}

// OrdersScoring 批量订单评分结果。
type OrdersScoring map[string]bool
