// types_rfq.go 模块
package polymarket

// CancelRfqRequestParams 取消 RFQ 请求的参数。
type CancelRfqRequestParams struct {
	RequestID string `json:"requestId"`
}

// CreateRfqRequestPayload 创建 RFQ 请求的请求体（Go 侧简化版，不包含 userType）。
type CreateRfqRequestPayload struct {
	AssetIn   string `json:"assetIn"`
	AssetOut  string `json:"assetOut"`
	AmountIn  string `json:"amountIn"`
	AmountOut string `json:"amountOut"`
}

// CreateRfqQuotePayload 创建 RFQ Quote 的请求体（Go 侧简化版，不包含 userType）。
type CreateRfqQuotePayload struct {
	RequestID string `json:"requestId"`
	AssetIn   string `json:"assetIn"`
	AssetOut  string `json:"assetOut"`
	AmountIn  string `json:"amountIn"`
	AmountOut string `json:"amountOut"`
}

// CancelRfqQuoteParams 取消 RFQ Quote 的参数。
type CancelRfqQuoteParams struct {
	QuoteID string `json:"quoteId"`
}

// AcceptQuoteParams 接受 RFQ Quote 的最小参数（仅用于调用方组织 payload）。
type AcceptQuoteParams struct {
	RequestID  string `json:"requestId"`
	QuoteID    string `json:"quoteId"`
	Expiration int64  `json:"expiration"`
}

// ApproveOrderParams 审批 RFQ Quote 的最小参数（仅用于调用方组织 payload）。
type ApproveOrderParams struct {
	RequestID  string `json:"requestId"`
	QuoteID    string `json:"quoteId"`
	Expiration int64  `json:"expiration"`
}

// RfqListState RFQ 列表状态筛选。
type RfqListState string

const (
	RfqListStateActive   RfqListState = "active"
	RfqListStateInactive RfqListState = "inactive"
)

// RfqSortDir RFQ 列表排序方向。
type RfqSortDir string

const (
	RfqSortDirAsc  RfqSortDir = "asc"
	RfqSortDirDesc RfqSortDir = "desc"
)

// RfqRequestsSortBy RFQ requests 的排序字段。
type RfqRequestsSortBy string

const (
	RfqRequestsSortByPrice   RfqRequestsSortBy = "price"
	RfqRequestsSortByExpiry  RfqRequestsSortBy = "expiry"
	RfqRequestsSortBySize    RfqRequestsSortBy = "size"
	RfqRequestsSortByCreated RfqRequestsSortBy = "created"
)

// RfqQuotesSortBy RFQ quotes 的排序字段。
type RfqQuotesSortBy string

const (
	RfqQuotesSortByPrice   RfqQuotesSortBy = "price"
	RfqQuotesSortByExpiry  RfqQuotesSortBy = "expiry"
	RfqQuotesSortByCreated RfqQuotesSortBy = "created"
)

// GetRfqRequestsParams 获取 RFQ requests 的查询参数（支持重复参数）。
type GetRfqRequestsParams struct {
	Offset string
	Limit  int
	State  RfqListState

	RequestIDs []string
	Markets    []string

	SizeMin     *float64
	SizeMax     *float64
	SizeUsdcMin *float64
	SizeUsdcMax *float64

	PriceMin *float64
	PriceMax *float64

	SortBy  RfqRequestsSortBy
	SortDir RfqSortDir
}

// GetRfqQuotesParams 获取 RFQ quotes 的查询参数（支持重复参数）。
type GetRfqQuotesParams struct {
	Offset string
	Limit  int
	State  RfqListState

	QuoteIDs   []string
	RequestIDs []string
	Markets    []string

	SizeMin     *float64
	SizeMax     *float64
	SizeUsdcMin *float64
	SizeUsdcMax *float64

	PriceMin *float64
	PriceMax *float64

	SortBy  RfqQuotesSortBy
	SortDir RfqSortDir
}

// GetRfqBestQuoteParams 获取最佳 quote 的查询参数。
type GetRfqBestQuoteParams struct {
	RequestID string
}

// RfqRequest RFQ request 结构（对照 Node SDK：RfqRequest）。
type RfqRequest struct {
	RequestID       string  `json:"requestId"`
	UserAddress     string  `json:"userAddress"`
	ProxyAddress    string  `json:"proxyAddress"`
	Token           string  `json:"token"`
	Complement      string  `json:"complement"`
	Condition       string  `json:"condition"`
	Side            string  `json:"side"`
	SizeIn          string  `json:"sizeIn"`
	SizeOut         string  `json:"sizeOut"`
	Price           float64 `json:"price"`
	AcceptedQuoteID string  `json:"acceptedQuoteId"`
	State           string  `json:"state"`
	Expiry          string  `json:"expiry"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// RfqMatchType RFQ match type（对照 Node SDK：RfqMatchType）。
type RfqMatchType string

const (
	RfqMatchTypeComplementary RfqMatchType = "COMPLEMENTARY"
	RfqMatchTypeMerge         RfqMatchType = "MERGE"
	RfqMatchTypeMint          RfqMatchType = "MINT"
)

// RfqQuote RFQ quote 结构（对照 Node SDK：RfqQuote）。
type RfqQuote struct {
	QuoteID      string       `json:"quoteId"`
	RequestID    string       `json:"requestId"`
	UserAddress  string       `json:"userAddress"`
	ProxyAddress string       `json:"proxyAddress"`
	Complement   string       `json:"complement"`
	Condition    string       `json:"condition"`
	Token        string       `json:"token"`
	Side         string       `json:"side"`
	SizeIn       string       `json:"sizeIn"`
	SizeOut      string       `json:"sizeOut"`
	Price        float64      `json:"price"`
	State        string       `json:"state"`
	Expiry       string       `json:"expiry"`
	MatchType    RfqMatchType `json:"matchType"`
	CreatedAt    string       `json:"createdAt"`
	UpdatedAt    string       `json:"updatedAt"`
}

// RfqRequestResponse 创建 RFQ request 的响应。
type RfqRequestResponse struct {
	RequestID string `json:"requestId"`
	Error     string `json:"error,omitempty"`
}

// RfqQuoteResponse 创建 RFQ quote 的响应。
type RfqQuoteResponse struct {
	QuoteID string `json:"quoteId"`
	Error   string `json:"error,omitempty"`
}
