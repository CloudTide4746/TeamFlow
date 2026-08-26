package event

import (
	"github.com/rabbitmq/amqp091-go"
)

const (
	EventsExchange      = "teamflow.events.v1"
	NotificationQueue   = "teamflow.notification.websocket.q"
	TaskAssignedRouting = "task.assigned.v1"
)

func DelareTopology(ch *amqp091.Channel) error {
	// 声明exchange
	err := ch.ExchangeDeclare(EventsExchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}
	// 声明queue
	_, err = ch.QueueDeclare(NotificationQueue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return ch.QueueBind(NotificationQueue, TaskAssignedRouting, EventsExchange, false, nil)
}
