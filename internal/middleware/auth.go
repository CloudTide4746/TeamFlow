package middleware

import (
	"strings"
	"teamflow/pkg/jwt"
	"teamflow/pkg/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		// 验证token
		parts := strings.SplitN(token, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		// 验证token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func GetCurrentUsername(c *gin.Context) (string, bool) {
	userID, ok := c.Get("username")
	if !ok {
		return "", false
	}
	return userID.(string), true
}
