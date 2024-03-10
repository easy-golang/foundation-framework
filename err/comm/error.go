package comm

import (
	"fmt"
	"github.com/wangliujing/foundation-framework/err"
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
		commonError := Error{Code: 500, Message: ie.GetMessage(), StackTrace: err.GetStackTrace(0)}
		commonError.Cause = ie
		return commonError
	}
	e, flag := obj.(error)
	if flag {
		commonError := Error{Code: 500, Message: e.Error(), StackTrace: err.GetStackTrace(0)}
		commonError.Cause = e
		return commonError
	}

	commonError := Error{Code: 500, Message: fmt.Sprintf("%v", obj), StackTrace: err.GetStackTrace(0)}
	return commonError
}

func WrapNew(message string, e error) Error {
	commonError := Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(0)}
	commonError.Cause = e
	return commonError
}

func WrapNewCode(message string, code int32, e error) Error {
	commonError := Error{Code: code, Message: message, StackTrace: err.GetStackTrace(0)}
	commonError.Cause = e
	return commonError
}

func WrapNewf(e error, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	commonError := Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(0)}
	commonError.Cause = e
	return commonError
}

func WrapNewCodef(e error, code int32, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	commonError := Error{Code: code, Message: message, StackTrace: err.GetStackTrace(0)}
	commonError.Cause = e
	return commonError
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
		commonError := Error{Code: 500, Message: e.Error(), StackTrace: err.GetStackTrace(skip)}
		commonError.Cause = e
		return commonError
	}
	commonError := Error{Code: 500, Message: fmt.Sprintf("%v", obj), StackTrace: err.GetStackTrace(skip)}
	return commonError
}

func WrapNewWithSkip(message string, e error, skip int) Error {
	commonError := Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(skip)}
	commonError.Cause = e
	return commonError
}

func WrapNewCodeWithSkip(message string, code int32, e error, skip int) Error {
	commonError := Error{Code: code, Message: message, StackTrace: err.GetStackTrace(skip)}
	commonError.Cause = e
	return commonError
}

func WrapNewWithSkipf(e error, skip int, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	commonError := Error{Code: 500, Message: message, StackTrace: err.GetStackTrace(skip)}
	commonError.Cause = e
	return commonError
}

func WrapNewCodeWithSkipf(e error, code int32, skip int, format string, args ...any) Error {
	message := fmt.Sprintf(format, args...)
	commonError := Error{Code: code, Message: message, StackTrace: err.GetStackTrace(skip)}
	commonError.Cause = e
	return commonError
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

func IsCommonErrorOrCause(e error) (*Error, bool) {
	commonError, ok := e.(Error)
	if ok {
		return &commonError, true
	} else {
		ic, ok := e.(err.ICause)
		if ok {
			return IsCommonErrorOrCause(ic.GetCause())
		}
	}
	return nil, false
}
