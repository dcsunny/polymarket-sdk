// types_notifications.go 模块
package polymarket

import "encoding/json"

// DropNotificationParams 删除通知参数。
type DropNotificationParams struct {
	IDs []string `json:"ids"`
}

// Notification 通知结构（GET /notifications）。
type Notification struct {
	Type    int             `json:"type"`
	Owner   string          `json:"owner"`
	Payload json.RawMessage `json:"payload"`
}
