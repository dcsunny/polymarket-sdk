// types_heartbeat.go 模块
package polymarket

// HeartbeatResponse 心跳响应。
type HeartbeatResponse struct {
	HeartbeatID string `json:"heartbeat_id"`
	Error       string `json:"error,omitempty"`
}
