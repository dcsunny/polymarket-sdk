package polymarket

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dcsunny/polymarket-sdk/internal/httpx"
)

// DataClient 处理 Polymarket Data API。
type DataClient struct {
	http *httpx.Client
}

// NewDataClient 创建一个新的 DataClient。
func NewDataClient(http *httpx.Client) *DataClient {
	return &DataClient{http: http}
}

// GetPositionsQuery 过滤持仓列表。
type GetPositionsQuery struct {
	User          string
	SizeThreshold *float64
	Redeemable    *bool
	Limit         int
	Offset        int
	Market        string //就是conditionId
}

// GetPositions 返回用户的持仓列表。
func (c *DataClient) GetPositions(ctx context.Context, q GetPositionsQuery) ([]*GammaPosition, error) {
	vals := url.Values{}
	if q.User != "" {
		vals.Set("user", q.User)
	}
	if q.SizeThreshold != nil {
		vals.Set("sizeThreshold", strconv.FormatFloat(*q.SizeThreshold, 'f', -1, 64))
	}
	if q.Redeemable != nil {
		vals.Set("redeemable", strconv.FormatBool(*q.Redeemable))
	}
	if q.Limit > 0 {
		vals.Set("limit", strconv.Itoa(q.Limit))
	}
	if q.Offset > 0 {
		vals.Set("offset", strconv.Itoa(q.Offset))
	}
	if q.Market != "" {
		vals.Set("market", q.Market)
	}

	var positions []*GammaPosition
	if err := c.http.Do(ctx, http.MethodGet, EndpointDataPositions, vals, nil, nil, &positions); err != nil {
		return nil, err
	}
	return positions, nil
}
