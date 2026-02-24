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
	if err := c.http.Do(ctx, http.MethodGet, EndpointGetMidpoint, vals, nil, nil, &resp); err != nil {
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
	if err := c.http.Do(ctx, http.MethodPost, EndpointGetMidpoints, nil, params, nil, &resp); err != nil {
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
	if err := c.http.Do(ctx, http.MethodPost, EndpointGetPrices, nil, params, nil, &resp); err != nil {
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
	if err := c.http.Do(ctx, http.MethodGet, EndpointGetSpread, vals, nil, nil, &resp); err != nil {
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
	if err := c.http.Do(ctx, http.MethodPost, EndpointGetSpreads, nil, params, nil, &resp); err != nil {
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
	if err := c.http.Do(ctx, http.MethodGet, EndpointGetLastTradePrice, vals, nil, nil, &resp); err != nil {
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
	if err := c.http.Do(ctx, http.MethodPost, EndpointGetLastTradesPrices, nil, params, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPricesByQuery 通过 query 参数批量获取市场价格（GET /prices）。
//
// tokenIDs 和 sides 分别为逗号分隔的 token ID 列表和对应的买卖方向（BUY/SELL）。
// 例如：tokenIDs = "id1,id2", sides = "BUY,SELL"
//
// API: GET /prices?token_ids=...&sides=...
// 文档: https://docs.polymarket.com/api-reference/market-data/get-market-prices-query-parameters
func (c *CLOBClient) GetPricesByQuery(ctx context.Context, tokenIDs string, sides string) (json.RawMessage, error) {
	if tokenIDs == "" {
		return nil, ErrInvalidArgument("tokenIDs is required")
	}
	if sides == "" {
		return nil, ErrInvalidArgument("sides is required")
	}
	vals := url.Values{}
	vals.Set("token_ids", tokenIDs)
	vals.Set("sides", sides)
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodGet, EndpointGetPrices, vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMidpointsByQuery 通过 query 参数批量获取 midpoint 中间价（GET /midpoints）。
//
// tokenIDs 为逗号分隔的 token ID 列表。
// 例如：tokenIDs = "id1,id2"
// 返回格式：{"token_id": "0.45", "token_id_2": "0.52"}
//
// API: GET /midpoints?token_ids=...
// 文档: https://docs.polymarket.com/api-reference/market-data/get-midpoint-prices-query-parameters
func (c *CLOBClient) GetMidpointsByQuery(ctx context.Context, tokenIDs string) (json.RawMessage, error) {
	if tokenIDs == "" {
		return nil, ErrInvalidArgument("tokenIDs is required")
	}
	vals := url.Values{}
	vals.Set("token_ids", tokenIDs)
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodGet, EndpointGetMidpoints, vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetLastTradesPricesByQuery 通过 query 参数批量获取最后成交价（GET /last-trades-prices）。
//
// tokenIDs 为逗号分隔的 token ID 列表，最多 500 个。
// 返回格式：[{"token_id": "...", "price": "0.45", "side": "BUY"}, ...]
//
// API: GET /last-trades-prices?token_ids=...
// 文档: https://docs.polymarket.com/api-reference/market-data/get-last-trade-prices-query-parameters
func (c *CLOBClient) GetLastTradesPricesByQuery(ctx context.Context, tokenIDs string) (json.RawMessage, error) {
	if tokenIDs == "" {
		return nil, ErrInvalidArgument("tokenIDs is required")
	}
	vals := url.Values{}
	vals.Set("token_ids", tokenIDs)
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodGet, EndpointGetLastTradesPrices, vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
