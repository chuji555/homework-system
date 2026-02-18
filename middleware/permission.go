package middleware

import (
	"github.com/chuji555/homework-system/pkg/errcode"
	"github.com/chuji555/homework-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware：只允许admin（老登）访问
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取角色（AuthMiddleware之后才能用）
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			response.Error(c, errcode.PermissionDenied)
			c.Abort()
			return
		}
		c.Next()
	}
}

// StudentMiddleware：只允许student（小登）访问
func StudentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "student" {
			response.Error(c, errcode.PermissionDenied)
			c.Abort()
			return
		}
		c.Next()
	}
}
