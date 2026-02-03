// errors_extra.go 模块
package polymarket

// InvalidArgumentError 表示无效的用户输入。
type InvalidArgumentError struct {
	Message string
}

func (e *InvalidArgumentError) Error() string {
	return e.Message
}

// ErrInvalidArgument 返回 InvalidArgumentError。
func ErrInvalidArgument(msg string) error {
	return &InvalidArgumentError{Message: msg}
}
