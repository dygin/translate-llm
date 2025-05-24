package queue

import "context"

// Message 消息结构
type Message struct {
	ID      string          // 消息ID
	Type    string          // 消息类型
	Data    []byte          // 消息数据
	Retries int             // 重试次数
}

// Queue 队列接口
type Queue interface {
	// Publish 发布消息
	Publish(ctx context.Context, msg *Message) error

	// Consume 消费消息
	Consume(ctx context.Context) (*Message, error)

	// Close 关闭连接
	Close() error
} 