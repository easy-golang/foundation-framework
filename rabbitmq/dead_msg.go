package rabbitmq

import (
	"github.com/easy-golang/foundation-framework/err/comm"
	"github.com/easy-golang/foundation-framework/orm"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type deadMsgDao struct {
	connection orm.Connection
}

func NewDeadMsgDao(connection orm.Connection) deadMsgDao {
	if connection == nil {
		logx.Must(comm.New("connection must not nil"))
	}
	return deadMsgDao{
		connection: connection,
	}
}

func (d deadMsgDao) GetByIdUnscoped(id int64) (*DeadMsgPO, error) {
	subscription := new(DeadMsgPO)
	db := d.connection.(*orm.GormConnection).DB
	err := db.Unscoped().Where("id = ?", id).First(subscription).Error
	if err != nil {
		if err == logger.ErrRecordNotFound {
			err = nil
			subscription = nil
		}
	}
	return subscription, err
}

func (d deadMsgDao) Save(po *DeadMsgPO) error {
	db := d.connection.(*orm.GormConnection).DB
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(po).Error
}

func (d deadMsgDao) DeletedById(id int64) error {
	db := d.connection.(*orm.GormConnection).DB
	return db.Where("id = ?", id).Delete(&DeadMsgPO{}).Error
}

func (d deadMsgDao) DeletedByIdUnscoped(id int64) error {
	db := d.connection.(*orm.GormConnection).DB
	return db.Unscoped().Where("id = ?", id).Delete(&DeadMsgPO{}).Error
}

func (d deadMsgDao) ListTrigerMsg() (result []DeadMsgPO, err error) {
	result = make([]DeadMsgPO, 0)
	db := d.connection.(*orm.GormConnection).DB
	err = db.Limit(100).Find(&result).Error
	return
}
