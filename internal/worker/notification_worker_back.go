package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"teamflow/internal/event"
	"teamflow/internal/model"
	"teamflow/internal/publisher"
	"teamflow/internal/service"

	"github.com/rabbitmq/amqp091-go"
)

type notificationWorker struct {
	notifier  service.NotificationService
	publisher publisher.Publisher
}

func NewNotificationWorker(notifier service.NotificationService, publishers ...publisher.Publisher) *notificationWorker {
	worker := &notificationWorker{notifier: notifier}
	if len(publishers) > 0 {
		worker.publisher = publishers[0]
	}
	return worker
}

// Handle keeps the original delivery-based entry point for callers that still
// consume raw task-assignment payloads directly.
func (w *notificationWorker) Handle(delivery amqp091.Delivery) error {
	var payload event.TaskAssignedPayload
	if err := json.Unmarshal(delivery.Body, &payload); err != nil {
		return fmt.Errorf("%w: decode task assignment payload: %v", event.ErrPermanent, err)
	}
	return w.handlePayload(payload)
}

// HandleEnvelope is the event-registry entry point used by the MQ consumer.
func (w *notificationWorker) HandleEnvelope(envelope event.Envelope) error {
	if envelope.EventType != event.TaskAssignedV1 {
		return fmt.Errorf("%w: unexpected event type: %s", event.ErrPermanent, envelope.EventType)
	}
	var payload event.TaskAssignedPayload
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return fmt.Errorf("%w: decode task assignment payload: %v", event.ErrPermanent, err)
	}
	return w.handlePayload(payload)
}

func (w *notificationWorker) handlePayload(payload event.TaskAssignedPayload) error {

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

	if err := w.notifier.PushIfOnline(notification); err != nil {
		log.Printf("notification realtime push failed: %v", err)
	}
	return nil
}

// Retry 处理失败的投递：按重试次数投递到对应的延迟重试队列，超过上限则进入停车场。
// 重试副本发布成功后才 ACK 原消息，避免副本未入队就丢消息。
func (w *notificationWorker) Retry(delivery amqp091.Delivery) error {
	attempt := retryCount(delivery.Headers) + 1
	routingKey, retry := retryRoute(attempt)
	if !retry {
		// 重试次数用尽：requeue=false，由主队列的 DLX 送入停车场
		return delivery.Nack(false, false)
	}

	// 复制原 headers 并记录新一次重试次数
	headers := make(amqp091.Table, len(delivery.Headers)+1)
	for k, v := range delivery.Headers {
		headers[k] = v
	}
	headers["x-retry-count"] = int32(attempt)

	if err := w.publisher.Publish(context.Background(), event.RetryExchange, routingKey, delivery.Body, headers); err != nil {
		// 副本未入队：不 ACK，原消息保持未确认，连接恢复后由 Broker 重新投递
		return err
	}
	return delivery.Ack(false)
}

// retryCount 读取已执行的失败重试次数
func retryCount(headers amqp091.Table) int32 {
	if v, ok := headers["x-retry-count"]; ok {
		switch n := v.(type) {
		case int32:
			return n
		case int64:
			return int32(n)
		}
	}
	return 0
}

// retryRoute 按重试次数选择延迟档位：1→5s，2→30s，3→5m，超过不再重试
func retryRoute(attempt int32) (string, bool) {
	switch attempt {
	case 1:
		return event.Retry5sBinding, true
	case 2:
		return event.Retry30sBinding, true
	case 3:
		return event.Retry5mBinding, true
	default:
		return "", false
	}
}
