package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ContextKeyRequestID = "requestID"

// RequestID 为每个请求生成唯一的 X-Request-ID 并写入上下文和响应头
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先使用客户端传入的 X-Request-ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 写入 gin 上下文
		c.Set(ContextKeyRequestID, requestID)

		// 在响应头中返回
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
