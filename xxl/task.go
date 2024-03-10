package xxl

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"runtime/debug"
)

// TaskFunc 任务执行函数
type TaskFunc func(cxt context.Context, param *RunReq, log Log) error

// Task 任务
type Task struct {
	Id        int64
	Name      string
	Ext       context.Context
	Param     *RunReq
	fn        TaskFunc
	Cancel    context.CancelFunc
	StartTime int64
	EndTime   int64
	log       Log
}

// Run 运行任务
func (t *Task) Run(callback func(code int64, msg string)) {
	defer func(cancel func()) {
		if err := recover(); err != nil {
			logx.Info(t.Info()+" panic: %v", err)
			debug.PrintStack() //堆栈跟踪
			callback(500, "task panic:"+fmt.Sprintf("%v", err))
			cancel()
		}
	}(t.Cancel)
	err := t.fn(t.Ext, t.Param, t.log)
	if err != nil {
		callback(500, err.Error())
		return
	}
	callback(200, "执行成功")
}

// Info 任务信息
func (t *Task) Info() string {
	return "任务ID[" + Int64ToStr(t.Id) + "]任务名称[" + t.Name + "]参数：" + t.Param.ExecutorParams
}
