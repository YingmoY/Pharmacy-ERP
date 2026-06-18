package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const contextKeyRequestID = "requestID"

// Response 是所有接口统一的响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id,omitempty"`
}

// GetRequestID 从 gin 上下文中获取请求 ID
func GetRequestID(c *gin.Context) string {
	if v, exists := c.Get(contextKeyRequestID); exists {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// Success 返回 HTTP 200 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   "success",
		Data:      data,
		RequestID: GetRequestID(c),
	})
}

// Fail 返回业务错误响应，HTTP 状态码始终为 200（保持向后兼容）
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:      code,
		Message:   msg,
		Data:      nil,
		RequestID: GetRequestID(c),
	})
	c.Abort()
}

// FailWithStatus 返回指定 HTTP 状态码的错误响应
func FailWithStatus(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code:      code,
		Message:   msg,
		Data:      nil,
		RequestID: GetRequestID(c),
	})
	c.Abort()
}
