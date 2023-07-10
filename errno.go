package apistatus

import (
	"log"
	"net/http"
	"strconv"
)

var (
	WarningLogFunc func(format string, v ...any) = log.Printf
)

// Errno 错误代码
type Errno int

// String 实现Stringer接口
func (no Errno) String() string {
	msg, ok := errnoMap[no]
	if ok {
		return msg
	}
	if msg = http.StatusText(int(no)); msg != "" {
		return msg
	}
	return "errno: " + strconv.Itoa(int(no))
}

// Error 实现error接口
func (no Errno) Error() string {
	return no.String()
}

// StatusCode HTTP状态码
func (no Errno) StatusCode() int {
	code, ok := errnoStatusCode[no]
	if ok {
		return code
	}
	if msg := http.StatusText(int(no)); msg != "" {
		return int(no)
	}
	return http.StatusInternalServerError
}

func (no Errno) GetErrno() Errno {
	return no
}

func (no Errno) GetMessage() string {
	return no.String()
}

func (no Errno) GetError() error {
	return no
}

func (no Errno) Is(target error) bool {
	if v, ok := target.(Errno); ok {
		return v == no
	}
	if i, ok := target.(interface {
		GetErrno() int
	}); ok {
		return i.GetErrno() == int(no)
	}
	return false
}

func (no Errno) As(target any) bool {
	switch target.(type) {
	case Errno:
		return true
	case *Errno:
		return true
	}
	return false
}

// RegisterErrno 注册错误代码
func RegisterErrno(m map[Errno]string) {
	for errno, message := range m {
		if msg, ok := errnoMap[errno]; ok && WarningLogFunc != nil {
			WarningLogFunc("errno: %d Message will replace: %s to %s!\n", errno, msg, message)
		}
		errnoMap[errno] = message
	}
}

// GetAllErrno 获取所有的错误代码
func GetAllErrno() map[Errno]string {
	m := make(map[Errno]string, len(errnoMap))
	for k, v := range errnoMap {
		m[k] = v
	}
	return m
}

// RegisteredErrno 错误代码是否已经注册
func RegisteredErrno(errno Errno) bool {
	_, ok := errnoMap[errno]
	return ok
}

// RegisterHTTPStatusCode 注册错误代码-HTTP状态码映射
func RegisterHTTPStatusCode(m map[Errno]int) {
	for errno, statusCode := range m {
		if code, ok := errnoStatusCode[errno]; ok && WarningLogFunc != nil {
			WarningLogFunc("errno: %d StatusCode will replace: %d to %d!\n", errno, code, statusCode)
		}
		errnoStatusCode[errno] = statusCode
	}
}

// GetAllHTTPStatusCode 获取所有的错误代码-HTTP状态码映射
func GetAllHTTPStatusCode() map[Errno]int {
	m := make(map[Errno]int, len(errnoStatusCode))
	for k, v := range errnoStatusCode {
		m[k] = v
	}
	return m
}

// RegisteredHTTPStatusCode 已经注册的HTTP状态码
func RegisteredHTTPStatusCode(errno Errno) bool {
	_, ok := errnoStatusCode[errno]
	return ok
}

const (
	Success             Errno = iota // 成功
	Failure                          // 失败
	Error                            // 错误
	Exception                        // 异常
	ErrUnauthorized                  // 认证失败
	ErrForbidden                     // 禁止
	ErrInvalidParameter              // 参数错误
	ErrValidate                      // 验证错误
	ErrSerialization                 // 序列化错误
	ErrDeserialization               // 反序列化错误
	ErrNotFound                      // 未找到
	ErrAlreadyExists                 // 已存在
	MaxErrno
)

const (
	// ErrnoReserve 保留的错误代码范围
	ErrnoReserve Errno = 1000 // 自定义的错误代码建议从此开始
)

// errnoMap 错误代码-消息映射表
var errnoMap = map[Errno]string{
	Success:             "Success",
	Failure:             "Failure",
	Error:               "Error",
	Exception:           "Exception",
	ErrUnauthorized:     "Unauthorized",
	ErrForbidden:        "Forbidden",
	ErrInvalidParameter: "Invalid Parameter",
	ErrValidate:         "Validate Error",
	ErrSerialization:    "Serialization Error",
	ErrDeserialization:  "Deserialization Error",
	ErrNotFound:         "Not Found",
	ErrAlreadyExists:    "Already Exists",
}

// errnoStatusCode 错误代码-HTTP状态码映射表
var errnoStatusCode = map[Errno]int{
	Success:             http.StatusOK,
	Failure:             http.StatusInternalServerError,
	Error:               http.StatusInternalServerError,
	Exception:           http.StatusInternalServerError,
	ErrUnauthorized:     http.StatusUnauthorized,
	ErrForbidden:        http.StatusForbidden,
	ErrInvalidParameter: http.StatusBadRequest,
	ErrValidate:         http.StatusBadRequest,
	ErrNotFound:         http.StatusNotFound,
	ErrAlreadyExists:    http.StatusConflict,
}
