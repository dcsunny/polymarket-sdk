// clob_markets.go 模块
package polymarket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// GetSamplingSimplifiedMarkets 获取 sampling simplified markets（分页）。
func (c *CLOBClient) GetSamplingSimplifiedMarkets(ctx context.Context, nextCursor string) (*PaginationPayload, error) {
	return c.getMarketList(ctx, "/sampling-simplified-markets", nextCursor)
}

// GetSamplingMarkets 获取 sampling markets（分页）。
func (c *CLOBClient) GetSamplingMarkets(ctx context.Context, nextCursor string) (*PaginationPayload, error) {
	return c.getMarketList(ctx, "/sampling-markets", nextCursor)
}

// GetSimplifiedMarkets 获取 simplified markets（分页）。
func (c *CLOBClient) GetSimplifiedMarkets(ctx context.Context, nextCursor string) (*PaginationPayload, error) {
	return c.getMarketList(ctx, "/simplified-markets", nextCursor)
}

// GetMarkets 获取 markets（分页，CLOB 服务下的 markets，不是 gamma-api）。
func (c *CLOBClient) GetMarkets(ctx context.Context, nextCursor string) (*PaginationPayload, error) {
	return c.getMarketList(ctx, "/markets", nextCursor)
}

// GetMarket 获取单个 market（GET /markets/{conditionId}）。
// 返回原始 JSON，调用方可自行解析为结构体。
func (c *CLOBClient) GetMarket(ctx context.Context, conditionID string) (json.RawMessage, error) {
	if conditionID == "" {
		return nil, ErrInvalidArgument("conditionID is required")
	}
	path := "/markets/" + url.PathEscape(conditionID)
	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodGet, path, nil, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CLOBClient) getMarketList(ctx context.Context, path string, nextCursor string) (*PaginationPayload, error) {
	vals := url.Values{}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	vals.Set("next_cursor", nextCursor)

	var resp PaginationPayload
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
