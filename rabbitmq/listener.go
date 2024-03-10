package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

const deadRoutingKey = "dead"

type Listener interface {
	CoroutineNum() int
	ConsumerParam() ConsumerParam
	AckParam() AckParam
	QueueDeclareParam() QueueDeclareParam
	OnDelivery(delivery *amqp.Delivery) error
}

type ConsumerParam struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

type AckParam struct {
	Multipart bool
	Requeue   bool
}

type QueueDeclareParam struct {
	QueueName    string
	Durable      bool
	AutoDelete   bool
	Exclusive    bool
	NoWait       bool
	RoutingKey   string
	ExchangeName string
}

type deadMsgListener struct {
	DeadQueueName    string
	DeadExchangeName string
	client           *Client
}

func (l deadMsgListener) CoroutineNum() int {
	return 1
}

func (l deadMsgListener) ConsumerParam() ConsumerParam {
	return ConsumerParam{
		Consumer:  "dead_queue_consumer",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}
}

func (l deadMsgListener) AckParam() AckParam {
	return AckParam{Multipart: false, Requeue: true}
}

func (l deadMsgListener) QueueDeclareParam() QueueDeclareParam {
	return QueueDeclareParam{
		QueueName:    l.DeadQueueName,
		Durable:      true,
		AutoDelete:   false,
		Exclusive:    false,
		NoWait:       false,
		RoutingKey:   deadRoutingKey,
		ExchangeName: l.DeadExchangeName,
	}
}

func (l deadMsgListener) OnDelivery(delivery *amqp.Delivery) error {
	// 设置这个属性保证死信消息消费不成功一直重试，防止死信消息丢失
	delivery.Redelivered = false
	logx.Info("收到死信消息：", string(delivery.Body))
	data, err := json.Marshal(NewDeliveryMirror(delivery))
	if err != nil {
		logx.Errorf("死信消息消费失败： %s \n %+v", string(delivery.Body), err)
		return err
	}
	err = NewDeadMsgDao(l.client.connection).Save(&DeadMsgPO{MessageId: delivery.MessageId,
		QueueName:       delivery.Headers["x-first-death-queue"].(string),
		Message:         string(delivery.Body),
		DeliveryMessage: string(data),
	})
	if err != nil {
		logx.Errorf("死信消息消费失败： %s \n %+v", string(delivery.Body), err)
		return err
	}
	return nil
}

type DeliveryMirror struct {
	Headers         amqp.Table
	ContentType     string
	ContentEncoding string
	DeliveryMode    uint8
	Priority        uint8
	CorrelationId   string
	ReplyTo         string
	Expiration      string
	MessageId       string
	Timestamp       time.Time
	Type            string
	UserId          string
	AppId           string

	ConsumerTag string

	MessageCount uint32

	DeliveryTag uint64
	Redelivered bool
	Exchange    string
	RoutingKey  string

	Body []byte
}

func NewDeliveryMirror(delivery *amqp.Delivery) *DeliveryMirror {
	return &DeliveryMirror{
		//Headers:         delivery.Headers,
		ContentType:     delivery.ContentType,
		ContentEncoding: delivery.ContentEncoding,
		DeliveryMode:    delivery.DeliveryMode,
		Priority:        delivery.Priority,
		CorrelationId:   delivery.CorrelationId,
		ReplyTo:         delivery.ReplyTo,
		Expiration:      delivery.Expiration,
		MessageId:       delivery.MessageId,
		Timestamp:       delivery.Timestamp,
		Type:            delivery.Type,
		UserId:          delivery.UserId,
		AppId:           delivery.AppId,
		ConsumerTag:     delivery.ConsumerTag,
		MessageCount:    delivery.MessageCount,
		DeliveryTag:     delivery.DeliveryTag,
		Redelivered:     delivery.Redelivered,
		Exchange:        delivery.Exchange,
		RoutingKey:      delivery.RoutingKey,
		Body:            delivery.Body,
	}
}

func (delivery *DeliveryMirror) ConvertorToDelivery() *amqp.Delivery {
	return &amqp.Delivery{
		//Headers:         delivery.Headers,
		ContentType:     delivery.ContentType,
		ContentEncoding: delivery.ContentEncoding,
		DeliveryMode:    delivery.DeliveryMode,
		Priority:        delivery.Priority,
		CorrelationId:   delivery.CorrelationId,
		ReplyTo:         delivery.ReplyTo,
		Expiration:      delivery.Expiration,
		MessageId:       delivery.MessageId,
		Timestamp:       delivery.Timestamp,
		Type:            delivery.Type,
		UserId:          delivery.UserId,
		AppId:           delivery.AppId,
		ConsumerTag:     delivery.ConsumerTag,
		MessageCount:    delivery.MessageCount,
		DeliveryTag:     delivery.DeliveryTag,
		Redelivered:     delivery.Redelivered,
		Exchange:        delivery.Exchange,
		RoutingKey:      delivery.RoutingKey,
		Body:            delivery.Body,
	}
}
