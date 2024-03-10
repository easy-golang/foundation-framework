package rabbitmq

import (
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/wangliujing/foundation-framework/orm"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/clause"
	"time"
)

type transactionalMsgDao struct {
	connection orm.Connection
}

func NewTransactionalMsgDao(connection orm.Connection) transactionalMsgDao {
	if connection == nil {
		logx.Must(comm.New("connection must not nil"))
	}
	return transactionalMsgDao{
		connection: connection,
	}
}
func (t transactionalMsgDao) ListTrigerMsg() (result []TransactionalMsgPO, err error) {
	result = make([]TransactionalMsgPO, 0)
	db := t.connection.(*orm.GormConnection).DB
	err = db.Where("trigger_time <= ?", time.Now()).Limit(100).Find(&result).Error
	return
}

func (t transactionalMsgDao) Save(transactionalMsgPO *TransactionalMsgPO) (*string, error) {
	db := t.connection.(*orm.GormConnection).DB
	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(transactionalMsgPO).Error
	return &transactionalMsgPO.Id, err
}

func (t transactionalMsgDao) DeletedById(id string) error {
	db := t.connection.(*orm.GormConnection).DB
	return db.Where("id = ?", id).Delete(&TransactionalMsgPO{}).Error
}

func (t transactionalMsgDao) BatchPhysicsDeletedWithEndTime(end time.Time) error {
	db := t.connection.(*orm.GormConnection).DB
	return db.Unscoped().Where("created_at < ?", end).Delete(&TransactionalMsgPO{}).Error
}
