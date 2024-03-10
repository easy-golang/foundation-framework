package handler

import (
	"github.com/wangliujing/foundation-framework/err"
	"github.com/wangliujing/foundation-framework/err/biz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type standardResult struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Code    int32  `json:"code"`
}

func OK(data any) standardResult {
	return standardResult{
		State:   true,
		Message: "操作成功",
		Data:    data,
		Code:    200,
	}
}

func Fail() standardResult {
	return standardResult{
		State:   false,
		Message: "操作失败",
		Code:    500,
	}
}

func CustomFail(code int32, message string) standardResult {
	return standardResult{
		State:   false,
		Message: message,
		Code:    code,
	}
}

func ExceptionFail(e error) standardResult {
	fromError, ok := status.FromError(e)
	if ok {
		if fromError.Code() != codes.Unknown {
			return standardResult{
				State:   false,
				Message: fromError.Message(),
				Code:    int32(fromError.Code()),
			}
		}
	}
	bizError, ok := biz.IsBizErrorOrCause(e)
	if ok {
		return standardResult{
			State:   false,
			Message: bizError.Message,
			Code:    bizError.Code,
		}
	}
	return standardResult{
		State:   false,
		Message: err.ServerError,
		Code:    500,
	}

}
