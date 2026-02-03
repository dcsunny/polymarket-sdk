// api_error.go 模块
package errors

import "fmt"

// APIError 表示非 2xx 响应。
type APIError struct {
	Status    int
	Code      string
	Message   string
	RequestID string
	Body      string
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Code != "" && e.Message != "" {
		return fmt.Sprintf("api error: status=%d code=%s msg=%s", e.Status, e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("api error: status=%d msg=%s", e.Status, e.Message)
	}
	return fmt.Sprintf("api error: status=%d", e.Status)
}
