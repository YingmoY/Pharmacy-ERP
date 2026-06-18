package handler

import (
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/YingmoY/PharmacyERP/internal/file/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// maxUploadSize 单次上传文件大小上限：50 MB。
	maxUploadSize = 50 << 20
)

// Handler 文件模块 HTTP 处理器。
type Handler struct {
	svc service.FileService
	log *zap.Logger
}

// NewHandler 创建文件处理器实例。
func NewHandler(svc service.FileService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// UploadFile 处理文件上传请求。
// POST /files/upload
// Content-Type: multipart/form-data
// Form fields: file（必填），business_type（可选），business_id（可选），remark（可选）
func (h *Handler) UploadFile(c *gin.Context) {
	// 限制最大上传大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	if err := c.Request.ParseMultipartForm(maxUploadSize); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "file too large or invalid multipart form")
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "file field is required")
		return
	}
	defer file.Close()

	var uploaderID int64
	if uid, ok := middleware.GetCurrentUserID(c); ok {
		uploaderID = uid
	}

	businessType := c.PostForm("business_type")
	businessID := c.PostForm("business_id")
	remark := c.PostForm("remark")

	var remarkPtr *string
	if remark != "" {
		remarkPtr = &remark
	}

	req := service.UploadFileRequest{
		File:         file,
		FileHeader:   fileHeader,
		BusinessType: businessType,
		BusinessID:   businessID,
		UploaderID:   uploaderID,
		Remark:       remarkPtr,
	}

	dto, err := h.svc.UploadFile(c.Request.Context(), req)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// GetFileInfo 获取文件元数据。
// GET /files/:file_id
func (h *Handler) GetFileInfo(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "file_id is required")
		return
	}

	dto, err := h.svc.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	core.Success(c, dto)
}

// DownloadFile 流式下载文件。
// GET /files/:file_id/download
func (h *Handler) DownloadFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "file_id is required")
		return
	}

	filePath, originalName, err := h.svc.GetFilePath(c.Request.Context(), fileID)
	if err != nil {
		bizErr := ecode.FromError(err)
		core.Fail(c, bizErr.Code, bizErr.Msg)
		return
	}

	// 设置下载响应头，兼容中文文件名（RFC 5987 编码）
	encodedName := url.QueryEscape(originalName)
	c.Header("Content-Disposition",
		`attachment; filename="`+filepath.Base(originalName)+`"; filename*=UTF-8''`+encodedName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.File(filePath)
}
