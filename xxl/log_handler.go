package xxl

import (
	"encoding/json"
	"fmt"
	"github.com/easy-golang/foundation-framework/orm"
	"github.com/easy-golang/foundation-framework/redis"
	"github.com/easy-golang/foundation-framework/util"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"net/http"
	"strings"
	"time"
)

/*
*
用来日志查询，显示到xxl-job-admin后台
*/

type Log interface {
	Info(param *RunReq, format string, a ...any)
	Error(param *RunReq, format string, a ...any)
	Handler(req *LogReq) *LogRes
}

type LogPO struct {
	Id        int64     `gorm:"type:bigint"`
	LogId     int64     `gorm:"index:xxl_log_id,priority:1,type:bigint"`
	Level     int8      `gorm:"type:tinyint"`
	Text      string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index:xxl_log_id,priority:2""`
}

func (LogPO) TableName() string {
	return "xxl_log"
}

type mysqlLog struct {
	conn        *orm.GormConnection
	close       bool
	redisClient *redis.Client
	conf        LogConf
}

func NewMysqlLog(conn *orm.GormConnection, redisClient *redis.Client, conf LogConf) *mysqlLog {
	// 自动创建表
	conn.DB.AutoMigrate(LogPO{})
	m := &mysqlLog{conn: conn, redisClient: redisClient, conf: conf}
	proc.AddShutdownListener(func() {
		m.close = true
	})
	go m.clearLog()
	return m
}

func (m *mysqlLog) clearLog() {
	if m.close {
		logx.Info("stop clearLog...")
		return
	}
	logx.Info("begin clearLog...")
	defer func() {
		go m.clearLog()
	}()
	defer util.Recover(nil)
	time.Sleep(time.Hour * 24)
	m.redisClient.NewRedisLock("clearLog").TryLockFunc(func() error {
		return m.conn.Where("created_at < ?", time.Now().Add(-24*time.Duration(m.conf.KeepDays)*time.Hour)).
			Delete(new(LogPO)).Error
	})
}

func (m *mysqlLog) Info(param *RunReq, format string, a ...interface{}) {
	logx.Infof(format, a...)
	go func() {
		log := &LogPO{
			LogId: param.LogID,
			Level: 0,
			Text:  fmt.Sprintf(format, a...),
		}
		m.conn.Save(log)
	}()
}

func (m *mysqlLog) Error(param *RunReq, format string, a ...any) {
	logx.Errorf(format, a...)
	go func() {
		log := &LogPO{
			LogId: param.LogID,
			Level: 1,
			Text:  fmt.Sprintf(format, a...),
		}
		m.conn.Save(log)
	}()
}

func (m *mysqlLog) Handler(req *LogReq) *LogRes {
	// 根据日志ID查询
	result := make([]LogPO, 0)
	err := m.conn.Where("log_id = ?", req.LogID).Order("created_at").Find(&result).Error
	if err != nil {
		logx.Errorf("find log error：%+v", err)
		return &LogRes{Code: 500, Msg: "", Content: LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  fmt.Sprintf("find log error：%+v", err),
			IsEnd:       true,
		}}
	}
	var builder strings.Builder
	for _, log := range result {
		builder.WriteString(log.CreatedAt.Format(util.NORMAL))
		if log.Level == 0 {
			builder.WriteString(" INFO  => ")
		} else {
			builder.WriteString(" ERROR => ")
		}
		builder.WriteString(log.Text)
		builder.WriteString("\n")
	}

	return &LogRes{Code: 200, Msg: "", Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  builder.String(),
		IsEnd:       true,
	}}
}

// 请求错误
func reqErrLogHandler(w http.ResponseWriter, req *LogReq, err error) {
	res := &LogRes{Code: 500, Msg: err.Error(), Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   0,
		LogContent:  err.Error(),
		IsEnd:       true,
	}}
	str, _ := json.Marshal(res)
	_, _ = w.Write(str)
}
