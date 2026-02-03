// clob_builder_trades.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// GetBuilderTradesPage 获取 builder/trades 分页数据（builder auth）。
func (c *CLOBClient) GetBuilderTradesPage(ctx context.Context, params *TradeParams, nextCursor string) (*BuilderTradesPage, error) {
	if c.builderAuth == nil {
		return nil, ErrInvalidArgument("builder auth is not configured")
	}

	path := "/builder/trades"
	headers, err := c.builderAuth.Headers(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	vals := url.Values{}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	vals.Set("next_cursor", nextCursor)

	if params != nil {
		if params.ID != "" {
			vals.Set("id", params.ID)
		}
		if params.MakerAddr != "" {
			vals.Set("maker_address", params.MakerAddr)
		}
		if params.Market != "" {
			vals.Set("market", params.Market)
		}
		if params.AssetID != "" {
			vals.Set("asset_id", params.AssetID)
		}
		if params.Before != nil {
			vals.Set("before", strconv.FormatInt(*params.Before, 10))
		}
		if params.After != nil {
			vals.Set("after", strconv.FormatInt(*params.After, 10))
		}
	}

	var resp BuilderTradesPage
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
