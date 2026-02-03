// clob_live_activity.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
)

// GetMarketTradesEvents 获取某个市场（condition_id）的实时成交活动（GET /live-activity/events/{condition_id}）。
func (c *CLOBClient) GetMarketTradesEvents(ctx context.Context, conditionID string) ([]MarketTradeEvent, error) {
	if conditionID == "" {
		return nil, ErrInvalidArgument("conditionID is required")
	}
	path := "/live-activity/events/" + url.PathEscape(conditionID)

	var resp []MarketTradeEvent
	if err := c.http.Do(ctx, http.MethodGet, path, nil, nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
