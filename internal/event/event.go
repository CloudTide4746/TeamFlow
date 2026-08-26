package event

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const TaskAssignedV1 = "task.assigned.v1"

type publisher interface {
	Publish(ctx context.Context, message Envelope) error
}
type Publisher struct {
	ch *amqp091.Channel // 此 Publisher 是这个 Channel 的唯一 owner
}

// Envelope 事件封包
// 事件封包用于封装事件的元数据和事件数据
type Envelope struct {
	EventID       string
	EventType     string
	SchemaVersion string
	OccurredAt    time.Time
	Payload       interface{}
}
