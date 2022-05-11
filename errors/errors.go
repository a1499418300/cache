package errors

var (
	ErrMemorySize      = NewError(10001, "内存单位有误，请检查")
	ErrGobEncodeFailed = NewError(10002, "gob编码失败")
)

// Error 错误对象封装
type Error struct {
	Code int32  `json:"errcode"`
	Msg  string `json:"errmsg"`
}

// Error 继承error
func (e Error) Error() string {
	return e.Msg
}

// NewError 新建error
func NewError(code int32, msg string) Error {
	return Error{
		Code: code,
		Msg:  msg,
	}
}
