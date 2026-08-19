package middleware

import (
	"strconv"
	"teamflow/internal/service"
	"teamflow/pkg/apperr"

	"github.com/gin-gonic/gin"
)

func TeamPermission(permSvc service.PermissionService, reqTeamRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextKeyUserID)
		if !exists {
			unauthorized := apperr.ErrUnauthorized
			_ = c.Error(unauthorized)
			return
		}
		teamIDStr := c.Param("teamID")
		teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
		if err != nil {
			unauthorized := apperr.ErrUnauthorized
			_ = c.Error(unauthorized)
			return
		}
		actualRole, ok, err := permSvc.CheckTeamPermisson(userID, teamID, reqTeamRole)
		if err != nil {
			unauthorized := apperr.ErrServerError
			_ = c.Error(unauthorized)
			return
		}
		if !ok {
			unauthorized := apperr.ErrUnauthorized
			_ = c.Error(unauthorized)
			return
		}
		c.Set(ContextKeyUserID, actualRole)
		c.Next()
	}

}
