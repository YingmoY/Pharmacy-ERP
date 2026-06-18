package middleware

import (
	"runtime/debug"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, rec interface{}) {
		log.Error("panic recovered",
			zap.Any("panic", rec),
			zap.ByteString("stack", debug.Stack()),
		)
		core.Fail(c, ecode.ErrSystem.Code, "system internal error")
		c.Abort()
	})
}
