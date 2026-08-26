package worker

import (
	"encoding/json"
	"fmt"
	"teamflow/internal/event"
	"teamflow/internal/model"
	"teamflow/internal/service"

	"github.com/rabbitmq/amqp091-go"
)

type notificationWorker struct {
	notifier service.NotificationService
}

func NewNotificationWorker(notifier service.NotificationService) *notificationWorker {
	return &notificationWorker{notifier: notifier}
}
func (w *notificationWorker) Handle(delivery amqp091.Delivery) error {
	var envelope event.Envelope
	if err := json.Unmarshal(delivery.Body, &envelope); err != nil {
		return fmt.Errorf("decode envelope: %w", err)
	}
	if envelope.EventType != event.TaskAssignedV1 {
		return fmt.Errorf("unexpected event type: %s", envelope.EventType)
	}

	payloadBytes, _ := json.Marshal(envelope.Payload)
	var payload event.TaskAssignedPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return fmt.Errorf("decode task assignment payload: %w", err)
	}

	notification, err := w.notifier.CreateNotification(
		payload.AssigneeID,
		payload.OperatorID,
		model.NotifyTaskAssigned,
		"任务分配通知",
		fmt.Sprintf("您已被分配任务 %d", payload.TaskID),
		payload.TaskID,
	)
	if err != nil {
		return err
	}
	return w.notifier.PushIfOnline(notification)
}
