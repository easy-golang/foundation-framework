package biz

import (
	"context"
	"fmt"
	"github.com/wangliujing/foundation-framework/err"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(message string) Error {
	return Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(0)}
}

func NewCode(code int32, message string) Error {
	return Error{Code: code, Message: message, StackTrace: err.GetStackTrace(0)}
}

func Newf(format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	return Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(0)}
}

func NewCodef(code int32, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	return Error{Code: code, Message: message, StackTrace: err.GetStackTrace(0)}
}

func Wrap(obj any) Error {
	ie, flag := obj.(err.IError)
	if flag {
		standardError := Error{Code: 500, Message: ie.GetMessage(), StackTrace: err.GetStackTrace(0)}
		standardError.Cause = ie
		return standardError
	}
	e, flag := obj.(error)
	if flag {
		standardError := Error{Code: 500, Message: e.Error(), StackTrace: err.GetStackTrace(0)}
		standardError.Cause = e
		return standardError
	}

	bizError := Error{Code: 500, Message: fmt.Sprintf("%v", obj), StackTrace: err.GetStackTrace(0)}
	return bizError
}

func WrapNew(message string, e error) Error {
	standardError := Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(0)}
	standardError.Cause = e
	return standardError
}

func WrapNewCode(message string, code int32, e error) Error {
	standardError := Error{Code: code, Message: message, StackTrace: err.GetStackTrace(0)}
	standardError.Cause = e
	return standardError
}

func WrapNewf(e error, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	standardError := Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(0)}
	standardError.Cause = e
	return standardError
}

func WrapNewCodef(e error, code int32, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	standardError := Error{Code: code, Message: message, StackTrace: err.GetStackTrace(0)}
	standardError.Cause = e
	return standardError
}

func NewWithSkip(message string, skip int) Error {
	return Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(skip)}
}

func NewCodeWithSkip(code int32, message string, skip int) Error {
	return Error{Code: code, Message: message, StackTrace: err.GetStackTrace(skip)}
}

func NewWithSkipf(skip int, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	return Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(skip)}
}

func NewCodeWithSkipf(code int32, skip int, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	return Error{Code: code, Message: message, StackTrace: err.GetStackTrace(skip)}
}

func WrapWithSkip(obj any, skip int) Error {
	e, flag := obj.(error)
	if flag {
		standardError := Error{Code: 500, Message: e.Error(), StackTrace: err.GetStackTrace(skip)}
		standardError.Cause = e
		return standardError
	}
	standardError := Error{Code: 500, Message: fmt.Sprintf("%v", obj), StackTrace: err.GetStackTrace(skip)}
	return standardError
}

func WrapNewWithSkip(message string, e error, skip int) Error {
	standardError := Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(skip)}
	standardError.Cause = e
	return standardError
}

func WrapNewCodeWithSkip(message string, code int32, e error, skip int) Error {
	standardError := Error{Code: code, Message: message, StackTrace: err.GetStackTrace(skip)}
	standardError.Cause = e
	return standardError
}

func WrapNewWithSkipf(e error, skip int, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	standardError := Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(skip)}
	standardError.Cause = e
	return standardError
}

func WrapNewCodeWithSkipf(e error, code int32, skip int, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	standardError := Error{Code: code, Message: message, StackTrace: err.GetStackTrace(skip)}
	standardError.Cause = e
	return standardError
}

type Error struct {
	State      bool
	Code       int32
	Message    string
	Cause      error
	StackTrace string
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) String() string {
	return err.String(e)
}

func (e Error) GetMessage() string {
	return e.Message
}

func (e Error) GetStackTrace() string {
	return e.StackTrace
}
func (e Error) GetCause() error {
	return e.Cause
}

func IsBizErrorOrCause(e error) (*Error, bool) {
	bizError, ok := e.(Error)
	if ok {
		return &bizError, true
	} else {
		ic, ok := e.(err.ICause)
		if ok {
			return IsBizErrorOrCause(ic.GetCause())
		}
	}
	return nil, false
}

func ErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	bizErr, b := IsBizErrorOrCause(err)
	if b {
		err = status.Error(codes.Code(bizErr.Code), bizErr.Error())
	}
	return resp, err
}
