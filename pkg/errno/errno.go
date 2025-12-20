package errno

import "fmt"

type Errno struct {
	Code    int
	Message string
}

// 实现error接口，让Errno可以当作error返回
func (e Errno) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
func New(code int, msg string) Errno {
	return Errno{
		Code:    code,
		Message: msg,
	}
}
func Decode(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	switch types := err.(type) {
	case Errno:
		return types.Code, types.Message
	default:
		return InternalServerError.Code, err.Error()
	}
}
