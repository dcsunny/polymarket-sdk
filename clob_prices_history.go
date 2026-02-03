// clob_prices_history.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// GetPricesHistory 获取市场价格历史（GET /prices-history）。
func (c *CLOBClient) GetPricesHistory(ctx context.Context, params PriceHistoryFilterParams) ([]MarketPrice, error) {
	vals := url.Values{}
	if params.Market != "" {
		vals.Set("market", params.Market)
	}
	if params.StartTs != nil {
		vals.Set("startTs", strconv.FormatInt(*params.StartTs, 10))
	}
	if params.EndTs != nil {
		vals.Set("endTs", strconv.FormatInt(*params.EndTs, 10))
	}
	if params.Fidelity != nil {
		vals.Set("fidelity", strconv.Itoa(*params.Fidelity))
	}
	if params.Interval != "" {
		vals.Set("interval", string(params.Interval))
	}

	var resp []MarketPrice
	if err := c.http.Do(ctx, http.MethodGet, "/prices-history", vals, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
