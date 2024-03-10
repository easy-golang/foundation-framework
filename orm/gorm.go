package orm

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type GormConfig struct {
	Dsn              string
	MaxIdleConns     int
	MaxOpenConns     int
	ConnMaxIdleTime  time.Duration
	ConnMaxLifetime  time.Duration
	OpenLogicDeleted bool
	LogLevel         int
}

type GormConnection struct {
	*gorm.DB
}

func NewGormConnection(config GormConfig) *GormConnection {
	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{
		// 开启软删除，使用 timestamp 记录删除时间
		DisableForeignKeyConstraintWhenMigrating: config.OpenLogicDeleted,
		Logger:                                   logger.Default.LogMode(logger.LogLevel(config.LogLevel)),
	})
	proc.AddShutdownListener(func() {
		db, err := db.DB()
		if err == nil {
			err = db.Close()
			if err == nil {
				return
			}
		}
		logx.Error("gorm数据库连接关闭失败")
	})
	if err == nil {
		sqlDB, err := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
		if err == nil {
			err := sqlDB.Ping()
			if err == nil {
				return &GormConnection{
					DB: db,
				}
			}
		}
	}
	logx.Must(err)
	return nil
}

func (g *GormConnection) Begin(opts ...*sql.TxOptions) Connection {
	return &GormConnection{
		DB: g.DB.Begin(opts...),
	}
}

func (g *GormConnection) Rollback() {
	g.DB.Rollback()
}

func (g *GormConnection) Commit() {
	g.DB.Commit()
}

func (g *GormConnection) SavePoint(name string) {
	g.DB.SavePoint(name)
}

func (g *GormConnection) RollbackTo(name string) {
	g.DB.RollbackTo(name)
}

type defaultTransaction struct {
	connection Connection
	opts       []*sql.TxOptions
}

func NewTransaction(connection Connection, opts ...*sql.TxOptions) Transaction {
	return &defaultTransaction{
		connection: connection,
		opts:       opts,
	}
}

func (g *defaultTransaction) Do(fun func(conn Connection) error) error {
	connection := g.connection.Begin(g.opts...)
	err := fun(connection)
	if err != nil {
		connection.Rollback()
		return err
	}
	connection.Commit()
	return nil
}

func (g *defaultTransaction) DoWithListener(fun func(conn Connection) error, commit CommitListener,
	rollback RollbackListener) error {
	connection := g.connection.Begin(g.opts...)
	err := fun(connection)
	if err != nil {
		connection.Rollback()
		if rollback != nil {
			rollback(err)
		}
		return err
	}

	connection.Commit()
	if commit != nil {
		commit()
	}
	return nil
}
