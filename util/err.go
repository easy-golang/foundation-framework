package util

import (
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/zeromicro/go-zero/core/logx"
)

func Recover(fun func(err error)) {
	if r := recover(); r != nil {
		err := comm.WrapWithSkip(r, 2)
		logx.Error(err)
		if fun != nil {
			fun(err)
		}
	}
}
