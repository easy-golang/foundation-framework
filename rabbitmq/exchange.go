package rabbitmq

import "github.com/streadway/amqp"

/*type ExchangeDeclare interface {
	ExchangeParam() ExchangeParam
}*/

type ExchangeDeclareParam struct {
	Name       string     // 交换器名
	Kind       string     // exchange type：一般用fanout、direct、topic
	Durable    bool       // 是否持久化
	AutoDelete bool       // 是否自动删除（自动删除的前提是至少有一个队列或者交换器与这和交换器绑定，之后所有与这个交换器绑定的队列或者交换器都与此解绑）
	Internal   bool       // 设置是否内置的。true表示是内置的交换器，客户端程序无法直接发送消息到这个交换器中，只能通过交换器路由到交换器这种方式
	NoWait     bool       // 是否阻塞
	Args       amqp.Table // 额外属性
}
