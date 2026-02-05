// clob_order_scoring.go 模块
package polymarket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// IsOrderScoring 查询单个订单是否参与奖励评分。
func (c *CLOBClient) IsOrderScoring(ctx context.Context, orderID string) (*OrderScoring, error) {
	if orderID == "" {
		return nil, ErrInvalidArgument("orderID is required")
	}
	path := EndpointIsOrderScoring
	vals := url.Values{}
	vals.Set("order_id", orderID)

	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp OrderScoring
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AreOrdersScoring 批量查询订单是否参与奖励评分。
func (c *CLOBClient) AreOrdersScoring(ctx context.Context, orderIDs []string) (OrdersScoring, error) {
	if len(orderIDs) == 0 {
		return nil, ErrInvalidArgument("orderIDs is required")
	}
	path := EndpointAreOrdersScoring
	// Node SDK: body is JSON array: ["orderId1","orderId2",...]
	body, err := json.Marshal(orderIDs)
	if err != nil {
		return nil, err
	}
	headers, err := c.l2Headers(http.MethodPost, path, string(body))
	if err != nil {
		return nil, err
	}

	var resp OrdersScoring
	if err := c.http.DoRaw(ctx, http.MethodPost, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
