package middleware

import (
	"github.com/chuji555/homework-system/pkg/errcode"
	"github.com/chuji555/homework-system/pkg/jwt"
	"github.com/chuji555/homework-system/pkg/response"

	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errcode.AuthError)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, errcode.AuthError)
			c.Abort()
			return
		}
		// 解析Token
		claims, errCode := jwt.ParseAccessToken(parts[1])
		if errCode != errcode.Success {
			response.Error(c, errCode)
			c.Abort()
			return
		}
		// 将用户信息存入上下文（后续接口可直接获取）
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("department", claims.Department)
		c.Next()
	}
}
