// clob_public.go 模块
package polymarket

import (
	"context"
	"encoding/json"
	"net/http"
)

// OK 用于探测 CLOB 服务是否可用（GET /）。
func (c *CLOBClient) OK(ctx context.Context) ([]byte, error) {
	var raw []byte
	if err := c.http.Do(ctx, http.MethodGet, "/", nil, nil, nil, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// GetServerTime 获取服务器时间（GET /time）。
// 返回值为 unix 秒时间戳（与 Node SDK 行为对齐）。
func (c *CLOBClient) GetServerTime(ctx context.Context) (int64, error) {
	var raw []byte
	if err := c.http.Do(ctx, http.MethodGet, EndpointTime, nil, nil, nil, &raw); err != nil {
		return 0, err
	}

	// 兼容不同返回格式：直接是数字 or {"time":123}
	var t int64
	if err := json.Unmarshal(raw, &t); err == nil {
		return t, nil
	}
	var wrapped struct {
		Time int64 `json:"time"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Time != 0 {
		return wrapped.Time, nil
	}
	return 0, json.Unmarshal(raw, &t)
}
