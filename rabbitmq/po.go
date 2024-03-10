package rabbitmq

import (
	"gorm.io/gorm"
	"time"
)

type TransactionalMsgPO struct {
	Id string `gorm:"type:char(36)"`

	Exchange string `gorm:"type:varchar(100)"`

	RoutingKey string `gorm:"type:varchar(100)"`

	Message string

	TriggerTime time.Time `gorm:"index:rabbit_trigger_time"`

	DelayedTime int64

	CreatedAt time.Time `gorm:"index:rabbit_create_at"`

	DeletedAt gorm.DeletedAt
}

func (TransactionalMsgPO) TableName() string {
	return "transactional_msg"
}

type DeadMsgPO struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	MessageId string `gorm:"type:char(36)"`

	QueueName string `gorm:"type:varchar(100)"`

	Message string `gorm:"type:json"`

	DeliveryMessage string `gorm:"type:json"`

	RetryCount int

	CreatedAt time.Time

	DeletedAt *gorm.DeletedAt
}

func (DeadMsgPO) TableName() string {
	return "dead_msg"
}
