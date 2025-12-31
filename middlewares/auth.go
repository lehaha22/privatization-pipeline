package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const SecretToken = "1mCjc7eLRMjm7Nh3fCpqLn58Phr2Knf6yHnKqg94k88=" // 设置你的 Token

// AuthMiddleware 用于验证 Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization 字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// 判断 Token 格式是否为 Bearer <token>
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// 提取 Token 并与 SecretToken 校验
		token := authHeader[len("Bearer "):]
		if token != SecretToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Token 校验通过，继续执行
		c.Next()
	}
}
