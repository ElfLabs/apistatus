package apistatus

import (
	"fmt"
)

// Option 选项
type Option func(s *Status)

// Apply 应用选项
func (s *Status) Apply(opts ...Option) *Status {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// NewStatus 新建
func NewStatus(opt ...Option) *Status {
	status := &Status{}
	return status.Apply(opt...)
}

func WithErrno(errno Errno) Option {
	return func(s *Status) {
		s.SetErrno(errno)
	}
}

func WithErrnoJust(errno Errno) Option {
	return func(s *Status) {
		s.Errno = errno
	}
}

func WithMessage(message string) Option {
	return func(s *Status) {
		s.Message = message
	}
}

func WithMessagef(format string, args ...interface{}) Option {
	return func(s *Status) {
		s.Message = fmt.Sprintf(format, args...)
	}
}

func WithError(err error) Option {
	return func(s *Status) {
		s.SetError(err)
	}
}

func WithErrorJust(err error) Option {
	return func(s *Status) {
		s.Err = err
	}
}

func WithData(data any) Option {
	return func(s *Status) {
		s.Data = data
	}
}
