package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"teamflow/internal/model"
	"teamflow/internal/ws"
	"teamflow/pkg/utils"

	"gorm.io/gorm"
)

type NotificationService interface {
	OnTaskStatusChange(task *model.Task, oldStatus, newStatus model.TaskStatus) error
	OnTaskAssigned(task *model.Task, assigneeID uint) error
}
type notificationService struct {
	db  *gorm.DB
	hub *ws.Hub
}

//为什么此处无需声明 直接在NotificationService接口中定义即可呢?
//Go 刻意去掉了 implements 关键字，让接口的实现关系变成隐式的、解耦的。好处是：
//你的 NoopNotificationService 完全不需要 import 接口所在的包，也能"无意中"实现它
//后续你写 EmailNotifier，也只需要把两个方法写出来，不用写任何声明，它就自动成了 NotificationService 的实现者
//这也就是注释里说的 "依赖倒置原则" —— 接口定义方和实现方完全解耦
//所以你的困惑很正常（从 Java/C# 转过来都会有这个疑问），但这正是 Go 接口设计的精妙之处。

func (n *notificationService) OnTaskStatusChange(task *model.Task, oldStatus, newStatus model.TaskStatus) error {
	return nil
}

func (n *notificationService) OnTaskAssigned(task *model.Task, assigneeID uint) error {
	//校验assignee是否在线
	if !n.hub.IsConnected(utils.FormatUintToString(assigneeID)) {
		return nil
	}
	//对assigneeID发送通知
	data, err := json.Marshal(ws.WebSocketMessage{
		Type: "notification",
		Payload: map[string]interface{}{
			"type":    "task_assigned",
			"title":   "任务分配通知",
			"content": fmt.Sprintf("您已被分配了任务 %d", task.ID),
		},
	})
	if err != nil {
		return err
	}
	if err := n.hub.SendToUser(utils.FormatUintToString(assigneeID), data); err != nil {
		return err
	}

	return nil
}

//### 后续实现思路
//
//| 实现类型                | 说明                             |
//| ------------------- | ------------------------------ |
//| `EmailNotifier`     | 通过 SMTP 发送邮件，调用 `email.Send()` |
//| `WebSocketNotifier` | 通过 WebSocket Hub 推送实时消息        |
//| `CompositeNotifier` | 组合多个 Notifier，同时发送多种通知         |
//
//后续在 `wire.go` 或初始化代码中替换 `NoopNotifier` 即可，**Service 层代码无需修改**（依赖倒置原则）。

func NewNotificationService(db *gorm.DB, hub *ws.Hub) *notificationService {
	return &notificationService{db: db, hub: hub}
}

// CreateNotification 创建并持久化一条通知
func (s *notificationService) CreateNotification(
	userID, senderID uint,
	notifyType model.NotificationType,
	title, content string,
	resourceID uint,
) (*model.Notification, error) {
	n := &model.Notification{
		UserID:   userID,
		SenderID: senderID,
		Type:     string(notifyType),
		Title:    title,
		Content:  content,
		RefID:    &resourceID,
		RefType:  string(notifyType),
	}
	if err := s.db.Create(n).Error; err != nil {
		return nil, err
	}
	return n, nil
}

// GetNotifications 分页获取指定用户的通知列表
func (s *notificationService) GetNotifications(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := s.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error
	return notifications, total, err
}

// MarkAsRead 将指定通知标记为已读
func (s *notificationService) MarkAsRead(notificationID, userID uint) error {
	result := s.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)
	if result.RowsAffected == 0 {
		return errors.New("通知不存在或无权操作")
	}
	return result.Error
}

// MarkAllAsRead 将用户所有通知标记为已读
func (s *notificationService) MarkAllAsRead(userID uint) error {
	return s.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
