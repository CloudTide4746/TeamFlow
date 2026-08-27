package publisher

import (
	"context"
	"encoding/json"
	"teamflow/internal/event"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp091.Channel // 此 Publisher 是这个 Channel 的唯一 owner
}

func NewPublisher(ch *amqp091.Channel) *Publisher {
	return &Publisher{ch: ch}
}

func (p *Publisher) Publish(ctx context.Context, e event.Envelope) error {
	//序列化body
	body, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return p.ch.PublishWithContext(ctx, event.EventsExchange, e.EventType, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		MessageId:    e.EventID,
		Type:         e.EventType,
		Timestamp:    e.OccurredAt,
		Body:         body,
	})
}
