package rabbitmq

import (
	"encoding/json"
	"github.com/easy-golang/foundation-framework/err/biz"
	"github.com/easy-golang/foundation-framework/err/comm"
	"github.com/easy-golang/foundation-framework/orm"
	"github.com/easy-golang/foundation-framework/redis"
	"github.com/easy-golang/foundation-framework/util"
	"github.com/easy-golang/foundation-framework/util/collection"
	"github.com/easy-golang/foundation-framework/util/ustring"
	"github.com/panjf2000/ants"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"gorm.io/gorm"
	"sync"
	"sync/atomic"
	"time"
)

type Conf struct {
	Url                 string
	SendReliableMessage bool   `json:",default=false"` // 使用可靠消息
	MaxRetryAttempts    int    `json:",default=3"`     // 错误消息重试次数(消息监听才需配置)
	DeadExchangeName    string `json:",optional"`      // 死信交换机(消息监听才需配置)
	DeadQueueName       string `json:",optional"`      // 死信队列(消息监听才需配置)
}

const (
	OPEN  = 1
	CLOSE = 0
)

type Client struct {
	*Conn
	conf                     Conf
	QueueNameListenerMapping map[string]Listener
	connection               orm.Connection
	redisClient              *redis.Client
	state                    int32
}

func NewClient(conf Conf, connection orm.Connection, redisClient *redis.Client) *Client {
	conn := newConn(conf.Url)
	client := &Client{
		Conn:        conn,
		conf:        conf,
		connection:  connection,
		redisClient: redisClient,
		state:       OPEN,
	}
	proc.AddShutdownListener(func() {
		err := client.Close()
		if err != nil {
			logx.Error(comm.WrapNew("关闭rabbit客户端失败", err))
		}
	})
	// 可靠消息
	if conf.SendReliableMessage {
		if connection == nil || redisClient == nil {
			logx.Must(comm.New("dbConnection or redisClient must not be nil"))
		}
		client.initTransactionalMsg()
	}
	return client
}

func (c *Client) Close() error {
	atomic.StoreInt32(&c.state, CLOSE)
	return c.Conn.Close()
}

func (c *Client) initTransactionalMsg() {
	gormConn := c.connection.(*orm.GormConnection)
	gormConn.DB.AutoMigrate(TransactionalMsgPO{})
	go c.deletedTransactionalMsg()
	go c.rabbitMqMessageRetrySend()
}

func (c *Client) initDeadMsg(connection orm.Connection) {
	gormConn := connection.(*orm.GormConnection)
	gormConn.DB.AutoMigrate(DeadMsgPO{})
	go c.rabbitMqDeadMessageRetrySend()
}

// 定时删除数据，每天执行一次
func (c *Client) deletedTransactionalMsg() {

	if atomic.LoadInt32(&c.state) == CLOSE {
		logx.Info("stop deletedTransactionalMsg...")
		return
	}
	logx.Info("begin deletedTransactionalMsg...")
	defer func() {
		go c.deletedTransactionalMsg()
	}()
	defer util.Recover(nil)
	time.Sleep(time.Hour * 24)
	c.redisClient.NewRedisLock("deletedTransactionalMsg").TryLockFunc(func() error {
		err := NewTransactionalMsgDao(c.connection).BatchPhysicsDeletedWithEndTime(time.Now().Add(-7 * 24 * time.Hour))
		if err == nil {
			logx.Info("deletedTransactionalMsg数据删除成功")
			return nil
		}
		logx.Error("deletedTransactionalMsg数据删除失败")
		return err
	})
}

func (c *Client) RegisterListener(listeners ...Listener) {
	if listeners == nil || len(listeners) == 0 {
		return
	}

	// 使用死信
	useDead := len(c.conf.DeadQueueName) != 0 && len(c.conf.DeadExchangeName) != 0
	if useDead {
		if c.connection == nil || c.redisClient == nil {
			logx.Must(comm.New("dbConnection or redisClient must not be nil"))
		}
		c.initDeadMsg(c.connection)

		// 声明死信交换机
		exchangeDeclareParam := ExchangeDeclareParam{
			Name:       c.conf.DeadExchangeName,
			Kind:       "direct",
			Durable:    true,
			AutoDelete: false,
			Internal:   true,
			NoWait:     false,
			Args:       nil,
		}
		c.DeclareExchange(exchangeDeclareParam)
		// 加死信监听
		listeners = append(listeners, deadMsgListener{DeadQueueName: c.conf.DeadQueueName,
			DeadExchangeName: c.conf.DeadExchangeName, client: c})
	}

	queueNameListenerMapping := make(map[string]Listener)
	c.QueueNameListenerMapping = queueNameListenerMapping
	// 声明监听队列队列
	for _, listener := range listeners {
		queueDeclareParam := listener.QueueDeclareParam()
		queueNameListenerMapping[queueDeclareParam.QueueName] = listener
		var args amqp.Table

		if useDead && listener.QueueDeclareParam().QueueName != c.conf.DeadQueueName {
			// 声明业务队列属性与死信交换机绑定
			args = amqp.Table{
				"x-dead-letter-exchange":    c.conf.DeadExchangeName,
				"x-dead-letter-routing-key": deadRoutingKey,
			}

		}
		go c.doStartListener(listener, args)
	}

}

func (c *Client) doStartListener(listener Listener, args amqp.Table) {
	c.DeclareQueue(args, listener.QueueDeclareParam())
	channel := c.GetChannelUntilSucess()
	defer channel.Close()
	consumerParam := listener.ConsumerParam()
	msgChanl, err := channel.Consume(
		listener.QueueDeclareParam().QueueName,
		consumerParam.Consumer,
		consumerParam.AutoAck,
		consumerParam.Exclusive,
		consumerParam.NoLocal,
		consumerParam.NoWait,
		consumerParam.Args,
	)
	if err != nil {
		logx.Errorf("启动监听器错误：%v %+v", consumerParam, err)
		time.Sleep(time.Second)
		logx.Infof("重新启动监听器：%v", consumerParam)
		go c.doStartListener(listener, args)
		return
	}
	coroutineNum := listener.CoroutineNum()
	if coroutineNum == 0 {
		coroutineNum = 1
	}
	pool, err := ants.NewPool(coroutineNum)
	defer pool.Release()
	if err != nil {
		logx.Must(comm.Wrap(err))
	}
	ackParam := listener.AckParam()
	// 从通道中接收数据
	for {
		delivery, ok := <-msgChanl
		if !ok {
			logx.Errorf("Channel is closed: %v", consumerParam)
			go c.doStartListener(listener, args)
			return
		}
		pool.Submit(func() {
			defer util.Recover(func(err error) {
				if !consumerParam.AutoAck {
					delivery.Nack(ackParam.Multipart, false)
				}
			})
			err := listener.OnDelivery(&delivery)
			if err == nil {
				if !consumerParam.AutoAck {
					delivery.Ack(ackParam.Multipart)
				}
				return
			}
			logx.Errorf("consumer error: %s %s", string(delivery.Body), err)
			if consumerParam.AutoAck {
				return
			}
			if ackParam.Requeue && !delivery.Redelivered {
				delivery.Nack(ackParam.Multipart, true)
			} else {
				delivery.Nack(ackParam.Multipart, false)
			}
		})
	}
}

/*func (c *Client) GetQueueDeclareParam(queueName string) *QueueDeclareParam {
	if c.queueDeclareParamMapping == nil {
		return nil
	}
	return c.queueDeclareParamMapping[queueName]
}*/

func (c *Client) DeclareExchange(exchangeParams ...ExchangeDeclareParam) {
	if exchangeParams == nil || len(exchangeParams) == 0 {
		return
	}
	channel := c.GetChannelUntilSucess()
	defer channel.Close()
	for _, param := range exchangeParams {
		// 声明交换器
		err := channel.ExchangeDeclare(
			param.Name,       //交换器名
			param.Kind,       //exchange type：一般用fanout、direct、topic
			param.Durable,    // 是否持久化
			param.AutoDelete, //是否自动删除（自动删除的前提是至少有一个队列或者交换器与这和交换器绑定，之后所有与这个交换器绑定的队列或者交换器都与此解绑）
			param.Internal,   //设置是否内置的。true表示是内置的交换器，客户端程序无法直接发送消息到这个交换器中，只能通过交换器路由到交换器这种方式
			param.NoWait,     // 是否阻塞
			param.Args,       // 额外属性
		)
		if err != nil {
			logx.Must(comm.Wrap(err))
		}
	}
}

func (c *Client) DeclareQueue(args amqp.Table, queueParams ...QueueDeclareParam) {
	if queueParams == nil || len(queueParams) == 0 {
		return
	}
	channel := c.GetChannelUntilSucess()
	defer channel.Close()
	for _, param := range queueParams {
		_, err := channel.QueueDeclare(
			param.QueueName,
			param.Durable,
			param.AutoDelete,
			param.Exclusive,
			param.NoWait,
			args,
		)
		if err != nil {
			logx.Must(comm.Wrap(err))
		}
		err = channel.QueueBind(
			param.QueueName,    // 绑定的队列名称
			param.RoutingKey,   // bindkey 用于消息路由分发的key
			param.ExchangeName, // 绑定的exchange名
			false,              // 是否阻塞
			nil,
		)
		if err != nil {
			logx.Must(comm.Wrap(err))
		}
	}
}

type Conn struct {
	*amqp.Connection
	lock *sync.Mutex
	//connCloseChan chan *amqp.Error
	config amqp.Config
	url    string
}

func newConn(url string) *Conn {
	connConfig := amqp.Config{
		Heartbeat: 10 * time.Second, // 设置心跳时间间隔
	}
	connection, err := amqp.DialConfig(url, connConfig)
	if err != nil {
		logx.Must(comm.Wrap(err))
	}
	conn := &Conn{
		Connection: connection,
		lock:       new(sync.Mutex),
		//connCloseChan: connCloseChan,
		url:    url,
		config: connConfig,
	}
	startNotifyCloseListener(conn)
	return conn
}

func startNotifyCloseListener(conn *Conn) {
	connCloseChan := make(chan *amqp.Error, 1)
	go func() {
		for {
			_, ok := <-connCloseChan
			if !ok {
				logx.Info("connCloseChan closed")
				conn.reConnect()
				break
			}
		}
	}()
	conn.Connection.NotifyClose(connCloseChan)
}

func (c *Conn) reConnect() {
	defer util.Recover(nil)
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.IsClosed() {
		return
	}
	// 连接配置
	Connection, err := amqp.DialConfig(c.url, c.config)

	if err == nil {
		logx.Info("重连成功")
		// 防止连接未关闭
		c.Connection.Close()
		c.Connection = Connection
		startNotifyCloseListener(c)
		// 再次校验
		if !c.IsClosed() {
			return
		}
	}
	logx.Info("重连失败：", err)
	// 一秒后重试
	time.Sleep(time.Second)
	// 这个地方不用异步会导致死锁
	go c.reConnect()
}

func (c *Conn) Close() error {
	logx.Info("关闭rabbitmq...")
	defer util.Recover(nil)
	err := c.Connection.Close()
	if err != nil {
		return comm.WrapNew("rabbitmq关闭失败", err)
	}
	return nil
}

func (c *Conn) GetChannel() (*amqp.Channel, error) {
	if c.Connection == nil || c.Connection.IsClosed() {
		logx.Info("连接关闭，重新初始化连接...")
		c.reConnect()
	}
	return c.Connection.Channel()
}

func (c *Conn) GetChannelUntilSucess() *amqp.Channel {
	for {
		if c.Connection == nil || c.Connection.IsClosed() {
			logx.Info("连接关闭，重新初始化连接...")
			c.reConnect()
		}
		channel, err := c.Connection.Channel()
		if err == nil {
			return channel
		}
		// 防止不停的重试导致cpu飙高
		time.Sleep(time.Second)
		logx.Error("rabbitMq获取会话异常,重新获取：", err)
	}
}

func (c *Client) rabbitMqMessageRetrySend() {
	if atomic.LoadInt32(&c.state) == CLOSE {
		logx.Info("stop rabbitMqMessageRetrySend...")
		return
	}
	logx.Info("begin rabbitMqMessageRetrySend...")
	defer func() {
		go c.rabbitMqMessageRetrySend()
	}()
	defer util.Recover(nil)
	time.Sleep(time.Minute)
	lock := c.redisClient.NewRedisLock("rabbitMqMessageRetrySend")
	lock.TryLockFunc(func() error {
		for {
			result, err := NewTransactionalMsgDao(c.connection).ListTrigerMsg()
			if collection.IsEmptySlice[TransactionalMsgPO](result) {
				break
			}
			if err != nil {
				logx.Error("定时任务出错：", err)
				return err
			}
			for _, value := range result {
				err := c.doRabbitMqMessageRetrySend(value)
				if err != nil {
					logx.Infof("定时任务发送消息失败：%+v", err)
				}
			}
			time.Sleep(time.Second * 10)
		}
		return nil
	})
}

func (c *Client) doRabbitMqMessageRetrySend(value TransactionalMsgPO) error {
	var message any
	if ustring.IsJson(value.Message) {
		message = make(map[string]any)
		err := json.Unmarshal([]byte(value.Message), &message)
		if err != nil {
			return err
		}
	} else {
		message = value.Message
	}
	if value.DelayedTime > 0 {
		err := c.NotifyPublishDelay(value.Exchange, value.RoutingKey, message, value.Id, value.DelayedTime, func(success bool, messageId string) error {
			if success {
				return NewTransactionalMsgDao(c.connection).DeletedById(messageId)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	err := c.NotifyPublish(value.Exchange, value.RoutingKey, message, value.Id, func(success bool, messageId string) error {
		if success {
			return NewTransactionalMsgDao(c.connection).DeletedById(messageId)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) rabbitMqDeadMessageRetrySend() {
	if atomic.LoadInt32(&c.state) == CLOSE {
		logx.Info("stop rabbitMqDeadMessageRetrySend...")
		return
	}
	logx.Info("begin rabbitMqDeadMessageRetrySend...")
	defer func() {
		go c.rabbitMqDeadMessageRetrySend()
	}()
	defer util.Recover(nil)
	time.Sleep(time.Minute)
	lock := c.redisClient.NewRedisLock("rabbitMqDeadMessageRetrySend")
	lock.TryLockFunc(func() error {
		for {
			result, err := NewDeadMsgDao(c.connection).ListTrigerMsg()
			if collection.IsEmptySlice[DeadMsgPO](result) {
				break
			}
			if err != nil {
				logx.Errorf("定时任务出错：%+v", err)
			}
			for _, value := range result {
				err := c.doRabbitMqDeadMessageRetrySend(value)
				if err != nil {
					logx.Errorf("重试消息失败：%+v", err)
				}
			}
			time.Sleep(time.Second * 10)
		}
		return nil
	})
}

func (c *Client) doRabbitMqDeadMessageRetrySend(value DeadMsgPO) error {
	// 获取执行监听器
	listener := c.QueueNameListenerMapping[value.QueueName]
	if listener == nil {
		return biz.Newf("队列[%s]不存在监听器", value.QueueName)
	}
	deliveryMirror := new(DeliveryMirror)
	err := json.Unmarshal([]byte(value.DeliveryMessage), deliveryMirror)
	if err == nil {
		err = listener.OnDelivery(deliveryMirror.ConvertorToDelivery())
	}
	if err != nil {
		value.RetryCount++
		if value.RetryCount >= c.conf.MaxRetryAttempts {
			// 更新次数，并逻辑删除
			deletedAt := new(gorm.DeletedAt)
			deletedAt.Scan(time.Now())
			value.DeletedAt = deletedAt
		}
		return NewDeadMsgDao(c.connection).Save(&value)
	}
	return NewDeadMsgDao(c.connection).DeletedByIdUnscoped(value.Id)
}

//const DeadLetter = "deadLetter"

func (c *Client) SendMessageInTransaction(exchange, key string, fn func(connection orm.Connection) error, messages ...any) error {
	return c.SendDelayMessageInTransaction(exchange, key, 0, fn, messages...)
}

func (c *Client) SendDelayMessageInTransaction(exchange, key string, delay int64, fn func(connection orm.Connection) error, messages ...any) error {
	messageMap := make(map[string]any)
	return orm.NewTransaction(c.connection).DoWithListener(func(conn orm.Connection) error {
		// 执行业务逻辑
		if fn != nil {
			err := fn(conn)
			if err != nil {
				return err
			}
		}

		for _, message := range messages {
			marshal, err := json.Marshal(message)
			if err != nil {
				return err
			}
			// 保存事物消息
			id, err := NewTransactionalMsgDao(conn).Save(&TransactionalMsgPO{Id: ustring.GetUUID(),
				Exchange: exchange, RoutingKey: key, Message: string(marshal), TriggerTime: time.Now().Add(time.Minute)})
			if err != nil {
				return err
			}
			messageMap[*id] = message
		}
		return nil
	}, func() {
		go func() {
			// 事物执行成功过后发送消息
			for id, message := range messageMap {
				c.NotifyPublishDelay(exchange, key, message, id, delay, func(success bool, messageId string) error {
					if success {
						return NewTransactionalMsgDao(c.connection).DeletedById(messageId)
					}
					return nil
				})
			}
		}()
	}, nil)
}

func (c *Client) NotifyPublish(exchange, key string, message any, messageId string, callBack func(success bool, messageId string) error) error {
	if body, flag := message.(string); flag {
		msg := amqp.Publishing{Body: []byte(body), MessageId: messageId, DeliveryMode: amqp.Persistent,
			ContentType: "application/json"}
		return c.DoNotifyPublish(exchange, key, msg, callBack)
	}
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{Body: body, MessageId: messageId, DeliveryMode: amqp.Persistent,
		ContentType: "application/json"}
	return c.DoNotifyPublish(exchange, key, msg, callBack)
}

func (c *Client) NotifyPublishDelay(exchange, key string, message any, messageId string, delay int64,
	callBack func(success bool, messageId string) error) error {
	if delay <= 0 {
		return c.NotifyPublish(exchange, key, message, messageId, callBack)
	}
	headers := amqp.Table{}
	headers["x-delay"] = delay

	if body, flag := message.(string); flag {
		msg := amqp.Publishing{Body: []byte(body), Headers: headers, MessageId: messageId, DeliveryMode: amqp.Persistent,
			ContentType: "application/json"}
		return c.DoNotifyPublish(exchange, key, msg, callBack)
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{Body: body, Headers: headers, MessageId: messageId, DeliveryMode: amqp.Persistent,
		ContentType: "application/json"}
	return c.DoNotifyPublish(exchange, key, msg, callBack)
}

func (c *Client) DoNotifyPublish(exchange, key string, msg amqp.Publishing, callBack func(success bool, messageId string) error) error {
	if !c.conf.SendReliableMessage {
		return comm.New("SendReliableMessage must be true")
	}
	// 开启发布者确认模式 注意这个地方不能直接关闭channel 否者确认模式无效，必须等确认完再关闭
	channel, err := c.GetChannel()
	if err != nil {
		return err
	}
	if err := channel.Confirm(false); err != nil {
		channel.Close()
		return err
	}
	go confirmOne(exchange, key, msg, channel, msg.MessageId, callBack)
	err = channel.Publish(
		exchange, // 交换器名
		key,      // routing key
		true,     // 是否返回消息(匹配队列)，如果为true, 会根据binding规则匹配queue，如未匹配queue，则把发送的消息返回给发送者
		false,    // 是否返回消息(匹配消费者)，如果为true, 消息发送到queue后发现没有绑定消费者，则把发送的消息返回给发送者
		msg,
	)
	if err != nil {
		logx.Error("rabbitMq发送消息失败")
		return err
	}
	return nil
}

func (c *Client) Publish(exchange, key string, message any) error {
	if body, flag := message.(string); flag {
		return c.DoPublish(exchange, key, amqp.Publishing{Body: []byte(body), DeliveryMode: amqp.Persistent, ContentType: "application/json"})
	}
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return c.DoPublish(exchange, key, amqp.Publishing{Body: body, DeliveryMode: amqp.Persistent, ContentType: "application/json"})
}

// PublishDelay delay 单位毫秒
func (c *Client) PublishDelay(exchange, key string, message any, delay int64) error {
	if delay <= 0 {
		return c.Publish(exchange, key, message)
	}
	headers := amqp.Table{}
	headers["x-delay"] = delay

	if body, flag := message.(string); flag {
		return c.DoPublish(exchange, key, amqp.Publishing{Body: []byte(body), DeliveryMode: amqp.Persistent,
			Headers: headers, ContentType: "application/json"})
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return c.DoPublish(exchange, key, amqp.Publishing{Body: body, DeliveryMode: amqp.Persistent,
		Headers: headers, ContentType: "application/json"})
}

func (c *Client) DoPublish(exchange, key string, msg amqp.Publishing) error {
	errChan := make(chan error, 0)
	defer close(errChan)
	err := c.DoNotifyPublish(exchange, key, msg, func(success bool, messageId string) error {
		if success {
			errChan <- nil
		} else {
			errChan <- biz.New("message send failed")
		}
		return nil
	})

	if err != nil {
		return err
	}
	err = <-errChan
	return err
}

func confirmOne(exchange, key string, msg amqp.Publishing, channel *amqp.Channel, messageId string,
	listener func(success bool, messageId string) error) {
	defer channel.Close()
	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	returns := channel.NotifyReturn(make(chan amqp.Return, 1))
	select {
	case confirm := <-confirms:
		if confirm.Ack {
			logx.Debug("发送消息成功：[", exchange, "][", key, "]\n", string(msg.Body))
			if listener != nil {
				err := listener(true, messageId)
				if err != nil {
					logx.Errorf("回调方法出错：%+v", err)
				}
			}
			return
		}
		logx.Debug("发送消息失败：[", exchange, "][", key, "]\n", string(msg.Body))
		if listener != nil {
			err := listener(false, messageId)
			if err != nil {
				logx.Errorf("回调方法出错：%+v", err)
			}
		}
	case _ = <-returns:
		logx.Error("消息无法路由被退回：[", exchange, "][", key, "]\n", string(msg.Body))
		if listener != nil {
			err := listener(false, messageId)
			if err != nil {
				logx.Errorf("回调方法出错：%+v", err)
			}
		}
	}
}
