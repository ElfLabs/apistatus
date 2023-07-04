package apistatus

import (
	"bytes"
	"fmt"
	"strconv"
)

// IStatus 状态接口
type IStatus interface {
	GetErrno() Errno
	GetMessage() string
	GetError() error
}

// Status 状态
type Status struct {
	// Errno 错误代码
	Errno Errno `json:"errno"`
	// Message 消息
	Message string `json:"message"`
	// Err 错误
	Err error `json:"error,omitempty"`
	// Data 数据
	Data any `json:"data,omitempty"`
}

func (s *Status) GetErrno() Errno {
	return s.Errno
}

func (s *Status) GetMessage() string {
	return s.Message
}

func (s *Status) GetError() error {
	return s.Err
}

// String 实现Stringer接口
func (s *Status) String() string {
	data, err := s.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("MarshalJSON error: %s", err)
	}
	return string(data)
}

// Error 实现error接口
func (s *Status) Error() string {
	if s.Err != nil {
		return s.Err.Error()
	}
	return s.String()
}

func (s *Status) Unwrap() error {
	return s.Err
}

func (s *Status) Is(target error) bool {
	switch o := target.(type) {
	case *Status:
		return s.Errno == o.Errno
	case IStatus:
		return s.Errno == o.GetErrno()
	case Errno:
		return s.Errno == o
	}
	return false
}

func (s *Status) As(target any) bool {
	switch target.(type) {
	case *Status:
		return true
	case Status:
		return true
	}
	return false
}

func (s *Status) SetErrno(errno Errno) *Status {
	s.Errno = errno
	s.Message = errno.String()
	return s
}

func (s *Status) JustSetErrno(errno Errno) *Status {
	s.Errno = errno
	return s
}

func (s *Status) SetMessage(message string) *Status {
	s.Message = message
	return s
}

func (s *Status) SetError(err error) *Status {
	switch o := err.(type) {
	case *Status:
		s.Errno = o.GetErrno()
		s.Message = o.GetMessage()
		s.Err = o.GetError()
	case Errno:
		s.Errno = o
		s.Message = o.String()
		s.Err = err
	case IStatus:
		s.Errno = o.GetErrno()
		s.Message = o.GetMessage()
		s.Err = o.GetError()
	default:
		s.Err = err
	}
	return s
}

func (s *Status) JustSetError(err error) *Status {
	s.Err = err
	return s
}

func (s *Status) SetData(data any) *Status {
	s.Data = data
	return s
}

// MarshalJSON 实现Marshaler接口
func (s *Status) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}

	buf.WriteString(`{"errno":`)
	buf.WriteString(strconv.Itoa(int(s.Errno)))
	buf.WriteString(`,"message":"`)
	buf.WriteString(s.Message)
	buf.WriteString(`"`)

	if s.Err != nil {
		buf.WriteString(`,"error":"`)
		buf.WriteString(s.Err.Error())
		buf.WriteString(`"`)
	}

	if s.Data != nil {
		buf.WriteString(`,"data":`)
		err := WriteData(&buf, s.Data)
		if err != nil {
			return nil, err
		}
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

// NewErrnoStatus 新建错误代码状态
func NewErrnoStatus(errno Errno) *Status {
	status := &Status{}
	return status.SetErrno(errno)
}

// NewErrorStatus 新建错误状态
func NewErrorStatus(errno Errno, err error) *Status {
	if status, ok := err.(*Status); ok {
		return status.SetErrno(errno)
	}
	status := &Status{}
	return status.SetErrno(errno).SetError(err)
}

// NewMessageStatus 新建消息状态
func NewMessageStatus(errno Errno, message string) *Status {
	return &Status{
		Errno:   errno,
		Message: message,
	}
}

// NewSuccessStatus 新建成功状态
func NewSuccessStatus() *Status {
	status := Status{}
	return status.SetErrno(Success)
}

// NewDataStatus 新建数据状态
func NewDataStatus(data any) *Status {
	status := Status{}
	return status.SetError(Success).SetData(data)
}
