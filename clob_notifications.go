// clob_notifications.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetNotifications 获取通知列表（L2 认证）。
func (c *CLOBClient) GetNotifications(ctx context.Context) ([]Notification, error) {
	path := EndpointGetNotifications
	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	vals := url.Values{}
	vals.Set("signature_type", strconv.Itoa(c.sigType))

	var resp []Notification
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DropNotifications 删除通知（L2 认证）。
// ids 为空时表示删除全部（与 Node SDK 行为一致）。
func (c *CLOBClient) DropNotifications(ctx context.Context, ids []string) error {
	path := EndpointDropNotifications
	headers, err := c.l2Headers(http.MethodDelete, path, "")
	if err != nil {
		return err
	}

	vals := url.Values{}
	if len(ids) > 0 {
		vals.Set("ids", strings.Join(ids, ","))
	}

	return c.http.Do(ctx, http.MethodDelete, path, vals, nil, headers, nil)
}
