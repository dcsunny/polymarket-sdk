// clob_heartbeat.go 模块
package polymarket

import (
	"context"
	"encoding/json"
	"net/http"
)

// PostHeartbeat 发送心跳（L2 认证，POST /v1/heartbeats）。
// 注意：启动 heartbeat 后，如果 10 秒内不继续发送，可能会触发订单自动取消（参考官方说明）。
func (c *CLOBClient) PostHeartbeat(ctx context.Context, heartbeatID *string) (*HeartbeatResponse, error) {
	path := "/v1/heartbeats"

	bodyObj := struct {
		HeartbeatID *string `json:"heartbeat_id"`
	}{
		HeartbeatID: heartbeatID,
	}
	body, err := json.Marshal(bodyObj)
	if err != nil {
		return nil, err
	}

	headers, err := c.l2Headers(http.MethodPost, path, string(body))
	if err != nil {
		return nil, err
	}

	var resp HeartbeatResponse
	if err := c.http.DoRaw(ctx, http.MethodPost, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
