// clob_market_data.go 模块
package polymarket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// GetMidpoint 获取单个 token 的 midpoint（GET /midpoint）。
func (c *CLOBClient) GetMidpoint(ctx context.Context, tokenID string) (json.RawMessage, error) {
	if tokenID == "" {
		return nil, ErrInvalidArgument("tokenID is required")
	}
	vals := url.Values{}
	vals.Set("token_id", tokenID)
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodGet, "/midpoint", vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMidpoints 批量获取 midpoint（POST /midpoints）。
func (c *CLOBClient) GetMidpoints(ctx context.Context, params []BookParams) (json.RawMessage, error) {
	if len(params) == 0 {
		return nil, ErrInvalidArgument("params is required")
	}
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodPost, "/midpoints", nil, params, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPrices 批量获取价格（POST /prices）。
func (c *CLOBClient) GetPrices(ctx context.Context, params []BookParams) (json.RawMessage, error) {
	if len(params) == 0 {
		return nil, ErrInvalidArgument("params is required")
	}
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodPost, "/prices", nil, params, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetSpread 获取单个 token 的 spread（GET /spread）。
func (c *CLOBClient) GetSpread(ctx context.Context, tokenID string) (json.RawMessage, error) {
	if tokenID == "" {
		return nil, ErrInvalidArgument("tokenID is required")
	}
	vals := url.Values{}
	vals.Set("token_id", tokenID)
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodGet, "/spread", vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetSpreads 批量获取 spread（POST /spreads）。
func (c *CLOBClient) GetSpreads(ctx context.Context, params []BookParams) (json.RawMessage, error) {
	if len(params) == 0 {
		return nil, ErrInvalidArgument("params is required")
	}
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodPost, "/spreads", nil, params, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetLastTradePrice 获取最后成交价（GET /last-trade-price）。
func (c *CLOBClient) GetLastTradePrice(ctx context.Context, tokenID string) (json.RawMessage, error) {
	if tokenID == "" {
		return nil, ErrInvalidArgument("tokenID is required")
	}
	vals := url.Values{}
	vals.Set("token_id", tokenID)
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodGet, "/last-trade-price", vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetLastTradesPrices 批量获取最后成交价（POST /last-trades-prices）。
func (c *CLOBClient) GetLastTradesPrices(ctx context.Context, params []BookParams) (json.RawMessage, error) {
	if len(params) == 0 {
		return nil, ErrInvalidArgument("params is required")
	}
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodPost, "/last-trades-prices", nil, params, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
