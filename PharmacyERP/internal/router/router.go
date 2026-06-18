package router

import (
	"net/http"

	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/gin-gonic/gin"
)

// ModuleRouter is the shared contract each business module router should implement.
type ModuleRouter interface {
	RegisterRoutes(group *gin.RouterGroup)
}

func RegisterAPIRoutes(engine *gin.Engine, modules ...ModuleRouter) {
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, core.Response{
			Code:      http.StatusNotFound,
			Message:   "route not found",
			Data:      nil,
			RequestID: core.GetRequestID(c),
		})
	})

	apiV1 := engine.Group("/api/v1")
	apiV1.GET("/health", func(c *gin.Context) {
		core.Success(c, gin.H{"status": "ok"})
	})

	engine.GET("/healthz", func(c *gin.Context) {
		core.Success(c, gin.H{"status": "ok"})
	})

	for _, m := range modules {
		if m == nil {
			continue
		}
		m.RegisterRoutes(apiV1)
	}
}
