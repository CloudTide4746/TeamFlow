package event

import (
	"time"

	"github.com/google/uuid"
)

type TaskAssignedPayload struct {
	TaskID     uint `json:"task_id"`
	ProjectID  uint `json:"project_id"`
	AssigneeID uint `json:"assignee_id"`
	OperatorID uint `json:"operator_id"`
}

// NewTaskAssigned 创建任务分配事件 用于创建任务分配事件的封包
func NewTaskAssigned(taskID, projectID, assigneeID, operatorID uint) Envelope {
	return Envelope{
		EventID:       uuid.NewString(),
		EventType:     TaskAssignedV1,
		SchemaVersion: "1",
		OccurredAt:    time.Now().UTC(),
		Payload: TaskAssignedPayload{
			TaskID: taskID, ProjectID: projectID,
			AssigneeID: assigneeID, OperatorID: operatorID,
		},
	}
}
