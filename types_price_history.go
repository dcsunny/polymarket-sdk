// types_price_history.go 模块
package polymarket

// MarketPrice 历史价格点（与 Node SDK 对齐：t=timestamp, p=price）。
type MarketPrice struct {
	T int64   `json:"t"`
	P float64 `json:"p"`
}

// PriceHistoryInterval 历史价格区间。
type PriceHistoryInterval string

const (
	PriceHistoryIntervalMax      PriceHistoryInterval = "max"
	PriceHistoryIntervalOneWeek  PriceHistoryInterval = "1w"
	PriceHistoryIntervalOneDay   PriceHistoryInterval = "1d"
	PriceHistoryIntervalSixHours PriceHistoryInterval = "6h"
	PriceHistoryIntervalOneHour  PriceHistoryInterval = "1h"
)

// PriceHistoryFilterParams 历史价格查询参数（GET /prices-history）。
type PriceHistoryFilterParams struct {
	Market   string
	StartTs  *int64
	EndTs    *int64
	Fidelity *int
	Interval PriceHistoryInterval
}
