package file

import (
	"github.com/YingmoY/PharmacyERP/internal/file/handler"
	"github.com/YingmoY/PharmacyERP/internal/file/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Module 是 file 业务模块的 HTTP 适配层。
type Module struct {
	db      *gorm.DB
	log     *zap.Logger
	handler *handler.Handler
}

// NewModule 构建文件模块，注入依赖。
// baseDir 为文件存储根目录，传空串使用默认值 "./uploads"。
func NewModule(db *gorm.DB, log *zap.Logger, baseDir string) *Module {
	svc := service.NewFileService(db, log, baseDir)
	h := handler.NewHandler(svc, log)
	return &Module{
		db:      db,
		log:     log,
		handler: h,
	}
}

// RegisterRoutes 注册文件模块路由。
// JWT 鉴权由上层路由组统一应用，此处不再重复添加。
func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	g := group.Group("/files")

	g.POST("/upload", m.handler.UploadFile)
	g.GET("/:file_id", m.handler.GetFileInfo)
	g.GET("/:file_id/download", m.handler.DownloadFile)
}
