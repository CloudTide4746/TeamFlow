package service

import (
	"log"
	"teamflow/internal/model"
)

type NotificationService interface {
	OnTaskStatusChange(task *model.Task, oldStatus, newStatus model.TaskStatus) error
	OnTaskAssigned(task *model.Task, assigneeID uint) error
}
type NoopNotificationService struct{}

var _ NotificationService = (*NoopNotificationService)(nil)

//为什么此处无需声明 直接在NotificationService接口中定义即可呢?
//Go 刻意去掉了 implements 关键字，让接口的实现关系变成隐式的、解耦的。好处是：
//你的 NoopNotificationService 完全不需要 import 接口所在的包，也能"无意中"实现它
//后续你写 EmailNotifier，也只需要把两个方法写出来，不用写任何声明，它就自动成了 NotificationService 的实现者
//这也就是注释里说的 "依赖倒置原则" —— 接口定义方和实现方完全解耦
//所以你的困惑很正常（从 Java/C# 转过来都会有这个疑问），但这正是 Go 接口设计的精妙之处。

func (n *NoopNotificationService) OnTaskStatusChange(task *model.Task, oldStatus, newStatus model.TaskStatus) error {
	return nil
}

func (n *NoopNotificationService) OnTaskAssigned(task *model.Task, assigneeID uint) error {
	log.Printf("[Notification] 任务 %d 已被分配给 %d", task.ID, assigneeID)
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
