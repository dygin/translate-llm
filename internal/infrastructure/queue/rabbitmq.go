package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// RabbitMQ RabbitMQ客户端
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  map[string]amqp.Queue
}

// NewRabbitMQ 创建RabbitMQ客户端
func NewRabbitMQ() (*RabbitMQ, error) {
	// 获取RabbitMQ配置
	host := g.Cfg().MustGet("rabbitmq.host").String()
	port := g.Cfg().MustGet("rabbitmq.port").Int()
	username := g.Cfg().MustGet("rabbitmq.username").String()
	password := g.Cfg().MustGet("rabbitmq.password").String()
	vhost := g.Cfg().MustGet("rabbitmq.vhost").String()

	// 连接RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("连接RabbitMQ失败: %v", err)
	}

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建通道失败: %v", err)
	}

	// 创建交换机
	err = channel.ExchangeDeclare(
		"task_exchange", // 交换机名称
		"direct",        // 交换机类型
		true,           // 持久化
		false,          // 自动删除
		false,          // 内部使用
		false,          // 不等待
		nil,            // 参数
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("创建交换机失败: %v", err)
	}

	// 创建队列
	queues := make(map[string]amqp.Queue)
	queueNames := []string{"content_generation", "translation"}

	for _, name := range queueNames {
		queue, err := channel.QueueDeclare(
			name,  // 队列名称
			true,  // 持久化
			false, // 自动删除
			false, // 独占
			false, // 不等待
			nil,   // 参数
		)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("创建队列失败: %v", err)
		}

		// 绑定队列到交换机
		err = channel.QueueBind(
			queue.Name,      // 队列名称
			name,           // 路由键
			"task_exchange", // 交换机名称
			false,          // 不等待
			nil,            // 参数
		)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("绑定队列失败: %v", err)
		}

		queues[name] = queue
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		queues:  queues,
	}, nil
}

// Close 关闭连接
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}

// Publish 发布消息
func (r *RabbitMQ) Publish(ctx context.Context, msg *Message) error {
	// 序列化消息
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 发布消息
	err = r.channel.PublishWithContext(ctx,
		"task_exchange", // 交换机名称
		msg.Type,       // 路由键
		false,         // 强制
		false,         // 立即
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息
			Timestamp:    time.Now(),
			MessageId:    msg.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("发布消息失败: %v", err)
	}

	return nil
}

// Consume 消费消息
func (r *RabbitMQ) Consume(ctx context.Context) (*Message, error) {
	// 设置预取数量
	err := r.channel.Qos(
		1,     // 预取数量
		0,     // 预取大小
		false, // 全局
	)
	if err != nil {
		return nil, fmt.Errorf("设置预取数量失败: %v", err)
	}

	// 开始消费
	msgs, err := r.channel.Consume(
		"content_generation", // 队列名称
		"",                   // 消费者标签
		false,               // 自动确认
		false,               // 独占
		false,               // 不等待
		false,               // 不阻塞
		nil,                 // 参数
	)
	if err != nil {
		return nil, fmt.Errorf("开始消费失败: %v", err)
	}

	// 等待消息
	select {
	case msg := <-msgs:
		// 解析消息
		var message Message
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			msg.Nack(false, true)
			return nil, fmt.Errorf("解析消息失败: %v", err)
		}

		// 确认消息
		msg.Ack(false)

		return &message, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
} 