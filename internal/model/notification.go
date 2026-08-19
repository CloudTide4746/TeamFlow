package model

import "gorm.io/gorm"

// Notification 通知模型
type Notification struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	Type    string `gorm:"size:50;not null" json:"type"`    // task_assigned/comment/mention 等
	Title   string `gorm:"size:200;not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	IsRead  bool   `gorm:"default:false" json:"is_read"`
	RefID   *uint  `json:"ref_id"`      // 关联资源 ID（如 TaskID）
	RefType string `gorm:"size:50" json:"ref_type"` // 关联资源类型（如 "task"）

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
