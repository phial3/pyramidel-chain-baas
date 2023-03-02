package mq

import (
	"fmt"
	"github.com/hxx258456/pyramidel-chain-baas/internal/localconfig"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	Mqurl string
}

// NewRabbitMQ 实例创建
func NewRabbitMQ() *RabbitMQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", localconfig.Defaultconfig.AMQP.User, localconfig.Defaultconfig.AMQP.Password, localconfig.Defaultconfig.AMQP.Host, localconfig.Defaultconfig.AMQP.Port, localconfig.Defaultconfig.AMQP.Vhost)
	return &RabbitMQ{QueueName: localconfig.Defaultconfig.AMQP.Queue, Exchange: "", Key: "", Mqurl: url}
}

func (mq *RabbitMQ) Destroy() {
	mq.channel.Close()
	mq.conn.Close()
}

func (mq *RabbitMQ) Connect() (err error) {
	mq.conn, err = amqp.Dial(mq.Mqurl)
	if err != nil {
		return err
	}
	mq.channel, err = mq.conn.Channel()
	if err != nil {
		return err
	}
	return nil
}

// PublishSimple 直接模式队列生产
func (r *RabbitMQ) Publish(message []byte) (err error) {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	queue, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		true,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	r.queue = &queue
	if err != nil {
		return err
	}
	//调用channel 发送消息到队列中
	if err := r.channel.Publish(
		r.Exchange,
		r.queue.Name,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: localconfig.Defaultconfig.AMQP.ContentType,
			Body:        message,
		}); err != nil {
		return err
	}
	return nil
}
