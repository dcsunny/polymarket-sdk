// clob_orderbook.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
)

// GetOrderBook 获取单个 token 的订单簿快照。
func (c *CLOBClient) GetOrderBook(ctx context.Context, tokenID string) (*OrderBookSummary, error) {
	if tokenID == "" {
		return nil, ErrInvalidArgument("tokenID is required")
	}
	vals := url.Values{}
	vals.Set("token_id", tokenID)
	var resp OrderBookSummary
	if err := c.http.Do(ctx, http.MethodGet, EndpointGetOrderBook, vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOrderBooks 批量获取多个 token 的订单簿。
func (c *CLOBClient) GetOrderBooks(ctx context.Context, params []BookParams) ([]*OrderBookSummary, error) {
	if len(params) == 0 {
		return nil, ErrInvalidArgument("params is required")
	}
	var resp []*OrderBookSummary
	if err := c.http.Do(ctx, http.MethodPost, EndpointGetOrderBooks, nil, params, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
