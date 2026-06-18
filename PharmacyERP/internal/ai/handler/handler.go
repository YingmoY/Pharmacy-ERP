// Package handler 提供 AI 发票识别 HTTP 处理层。
package handler

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/ai/repository"
	"github.com/YingmoY/PharmacyERP/internal/ai/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler AI 发票 HTTP 处理器。
type Handler struct {
	svc service.AIInvoiceService
}

// New 创建 Handler 实例。
func New(svc service.AIInvoiceService) *Handler {
	return &Handler{svc: svc}
}

// ============================================================
// 请求 DTO
// ============================================================

// convertToInboundRequest 转换为入库单请求体。
type convertToInboundRequest struct {
	SupplierID int64                    `json:"supplier_id" binding:"required"`
	Items      []service.ConvertItemReq `json:"items" binding:"required,min=1"`
}

// ============================================================
// 辅助方法
// ============================================================

// getInvoiceID 从路径参数 ":id" 解析 AI 发票记录 ID。
func getInvoiceID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "无效的发票 ID")
		return 0, false
	}
	return id, true
}

// requireUserID 从 JWT 上下文获取当前用户 ID，不存在则响应 401。
func requireUserID(c *gin.Context) (int64, bool) {
	uid, ok := middleware.GetCurrentUserID(c)
	if !ok || uid <= 0 {
		core.Fail(c, ecode.ErrUnauthorized.Code, ecode.ErrUnauthorized.Msg)
		return 0, false
	}
	return uid, true
}

// handleBizErr 统一处理业务错误响应。
func handleBizErr(c *gin.Context, err error) {
	bizErr := ecode.FromError(err)
	core.Fail(c, bizErr.Code, bizErr.Msg)
}

// ============================================================
// HTTP 处理方法
// ============================================================

// RecognizeInvoice POST /ai/invoices/recognize
// 接受 multipart/form-data 上传发票文件，触发 AI 识别并同步返回识别结果。
func (h *Handler) RecognizeInvoice(c *gin.Context) {
	creatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, "file 字段必填")
		return
	}

	f, err := fh.Open()
	if err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "读取文件失败")
		return
	}
	defer f.Close()
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "读取文件失败")
		return
	}
	contentType := detectContentType(fh.Header.Get("Content-Type"), fh.Filename, fileBytes)

	fileID := uuid.New().String()
	fileName := fh.Filename

	var supplierID *int64
	if sv := c.PostForm("supplier_id"); sv != "" {
		if v, err2 := strconv.ParseInt(sv, 10, 64); err2 == nil && v > 0 {
			supplierID = &v
		}
	}
	remark := c.PostForm("remark")

	record, err := h.svc.RecognizeInvoice(c.Request.Context(), fileID, fileName, fileBytes, contentType, supplierID, remark, creatorID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, record)
}

// GetInvoice GET /ai/invoices/:id
func (h *Handler) GetInvoice(c *gin.Context) {
	id, ok := getInvoiceID(c)
	if !ok {
		return
	}
	record, err := h.svc.GetInvoice(c.Request.Context(), id)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, record)
}

// ListInvoices GET /ai/invoices
func (h *Handler) ListInvoices(c *gin.Context) {
	filter := repository.AIInvoiceListFilter{
		Status:   c.Query("status"),
		Page:     1,
		PageSize: 20,
	}
	if v, err := strconv.ParseInt(c.Query("supplier_id"), 10, 64); err == nil && v > 0 {
		filter.SupplierID = v
	}
	if v, err := strconv.Atoi(c.Query("page")); err == nil && v > 0 {
		filter.Page = v
	}
	if v, err := strconv.Atoi(c.Query("page_size")); err == nil && v > 0 {
		filter.PageSize = v
	}
	if v := c.Query("start_date"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err == nil {
			filter.StartDate = &t
		}
	}
	if v := c.Query("end_date"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err == nil {
			end := t.Add(24*time.Hour - time.Nanosecond)
			filter.EndDate = &end
		}
	}

	records, total, err := h.svc.ListInvoices(c.Request.Context(), filter)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, core.NewPageResult(total, filter.Page, filter.PageSize, records))
}

// ConvertToInbound POST /ai/invoices/:id/convert-to-inbound
// 将已识别完成的发票转换为草稿入库单。
func (h *Handler) ConvertToInbound(c *gin.Context) {
	id, ok := getInvoiceID(c)
	if !ok {
		return
	}
	operatorID, ok := requireUserID(c)
	if !ok {
		return
	}

	var req convertToInboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, err.Error())
		return
	}

	order, err := h.svc.ConvertToInbound(c.Request.Context(), id, service.ConvertToInboundReq{
		SupplierID: req.SupplierID,
		Items:      req.Items,
	}, operatorID)
	if err != nil {
		handleBizErr(c, err)
		return
	}
	core.Success(c, order)
}

// detectContentType 从 header、文件名后缀、字节内容三个来源推断真实 MIME 类型。
// 浏览器有时会将 multipart part 的 Content-Type 设为 application/octet-stream，
// 此时需要通过扩展名或文件头字节嗅探。
func detectContentType(headerCT, filename string, data []byte) string {
	// 已知精确类型直接使用
	known := map[string]bool{
		"image/jpeg": true, "image/jpg": true, "image/png": true,
		"image/webp": true, "application/pdf": true,
	}
	if known[strings.ToLower(headerCT)] {
		return strings.ToLower(headerCT)
	}

	// 从文件名后缀推断
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	}

	// 从字节头嗅探
	sniff := data
	if len(sniff) > 512 {
		sniff = sniff[:512]
	}
	detected := http.DetectContentType(sniff)
	if known[detected] {
		return detected
	}

	// PDF 字节头 %PDF 不总被 DetectContentType 识别为 application/pdf
	if len(data) >= 4 && string(data[:4]) == "%PDF" {
		return "application/pdf"
	}

	return "image/jpeg" // 兜底
}
