// types_builder_trades.go 模块
package polymarket

// TradeParams 交易查询参数（用于 /data/trades /builder/trades 等）。
type TradeParams struct {
	ID        string
	MakerAddr string
	Market    string
	AssetID   string
	Before    *int64
	After     *int64
}

// BuilderTrade builder/trades 返回的成交记录。
type BuilderTrade struct {
	ID              string  `json:"id"`
	TradeType       string  `json:"tradeType"`
	TakerOrderHash  string  `json:"takerOrderHash"`
	Builder         string  `json:"builder"`
	Market          string  `json:"market"`
	AssetID         string  `json:"assetId"`
	Side            string  `json:"side"`
	Size            string  `json:"size"`
	SizeUsdc        string  `json:"sizeUsdc"`
	Price           string  `json:"price"`
	Status          string  `json:"status"`
	Outcome         string  `json:"outcome"`
	OutcomeIndex    int     `json:"outcomeIndex"`
	Owner           string  `json:"owner"`
	Maker           string  `json:"maker"`
	TransactionHash string  `json:"transactionHash"`
	MatchTime       string  `json:"matchTime"`
	BucketIndex     int     `json:"bucketIndex"`
	Fee             string  `json:"fee"`
	FeeUsdc         string  `json:"feeUsdc"`
	ErrMsg          *string `json:"err_msg,omitempty"`
	CreatedAt       *string `json:"createdAt,omitempty"`
	UpdatedAt       *string `json:"updatedAt,omitempty"`
}

// BuilderTradesPage builder/trades 分页响应。
type BuilderTradesPage struct {
	Data       []BuilderTrade `json:"data"`
	NextCursor string         `json:"next_cursor"`
	Limit      int            `json:"limit"`
	Count      int            `json:"count"`
}
