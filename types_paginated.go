// types_paginated.go 模块
package polymarket

// PaginatedResponse 通用分页响应结构（CLOB 服务常见格式）。
// data 为强类型切片；next_cursor 为 base64 游标。
type PaginatedResponse[T any] struct {
	Limit      int    `json:"limit"`
	Count      int    `json:"count"`
	NextCursor string `json:"next_cursor"`
	Data       []T    `json:"data"`

	// 某些接口会返回 total_count（例如 RFQ）。
	TotalCount *int `json:"total_count,omitempty"`
}
