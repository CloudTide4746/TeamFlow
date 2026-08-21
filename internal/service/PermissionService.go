package service

import (
	"errors"
	"teamflow/internal/model"
	"teamflow/pkg/apperr"
	"teamflow/storage"

	"gorm.io/gorm"
)

type PermissionService struct {
}
type ProjectPermissionService interface {
	CheckTeamPermisson(userID uint, teamID uint, reqTeamRole string) (string, bool, error)
	CheckProjectPermisson(userID uint, projectID uint, reqProjectRole string) (string, bool, error)
}

func NewProjectPermissionService() *PermissionService {
	return &PermissionService{}
}

func (s *PermissionService) CheckTeamPermisson(userID uint, teamID uint, reqTeamRole string) (string, bool, error) {
	var teamMember model.TeamMember
	if err := storage.DB.Where("user_id = ? AND team_id = ?", userID, teamID).First(&teamMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			unauthorized := apperr.ErrServerError
			return "", false, unauthorized
		}
		return "", false, apperr.ErrServerError
	}
	return teamMember.Role, true, nil
}

func (s *PermissionService) CheckProjectPermisson(userID uint, projectID uint, reqProjectRole string) (string, bool, error) {
	var projectMember model.ProjectMember
	if err := storage.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&projectMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			unauthorized := apperr.ErrServerError
			return "", false, unauthorized
		}
		return "", false, apperr.ErrServerError
	}
	return projectMember.Role, true, nil
}
