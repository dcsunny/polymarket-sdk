// types_pagination.go 模块
package polymarket

import "encoding/json"

const (
	// InitialCursor 初始分页游标（与 Node SDK 对齐）。
	InitialCursor = "MA=="
)

// PaginationPayload 通用分页响应结构（CLOB 服务）。
type PaginationPayload struct {
	Limit      int               `json:"limit"`
	Count      int               `json:"count"`
	NextCursor string            `json:"next_cursor"`
	Data       []json.RawMessage `json:"data"`
}
