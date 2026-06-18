// Package service 实现 AI 发票识别相关业务逻辑。
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"time"

	aiModel "github.com/YingmoY/PharmacyERP/internal/ai/model"
	"github.com/YingmoY/PharmacyERP/internal/ai/repository"
	inboundModel "github.com/YingmoY/PharmacyERP/internal/inbound/model"
	inboundRepo "github.com/YingmoY/PharmacyERP/internal/inbound/repository"
	inboundService "github.com/YingmoY/PharmacyERP/internal/inbound/service"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// =====================================================
// DTO 定义
// =====================================================

// ConvertItemReq 转换为入库单时的明细行请求。
type ConvertItemReq struct {
	DrugID      int64   `json:"drug_id"`
	BatchNumber string  `json:"batch_number"`
	ExpireDate  string  `json:"expire_date"` // YYYY-MM-DD
	PlannedQty  int32   `json:"planned_qty"`
	UnitPrice   float64 `json:"unit_price"`
	Remark      string  `json:"remark"`
}

// ConvertToInboundReq 转换为入库单请求参数。
type ConvertToInboundReq struct {
	SupplierID int64            `json:"supplier_id"`
	Items      []ConvertItemReq `json:"items"`
}

// AIInvoiceListFilter 透传仓储层查询条件。
type AIInvoiceListFilter = repository.AIInvoiceListFilter

// =====================================================
// Python AI 服务响应结构体
// =====================================================

// aiInvoiceResult 对应 Python ai-openapi.yaml InvoiceRecognizeResult。
type aiInvoiceResult struct {
	RecognizedSupplierName *string         `json:"recognized_supplier_name"`
	MatchedSupplierID      *int64          `json:"matched_supplier_id"`
	InvoiceNo              *string         `json:"invoice_no"`
	InvoiceDate            *string         `json:"invoice_date"` // "YYYY-MM-DD"
	TotalAmount            *string         `json:"total_amount"`
	Confidence             float64         `json:"confidence"`
	Items                  json.RawMessage `json:"items"`
}

// aiInvoiceResponseData 对应 Python InvoiceRecognizeResponse。
type aiInvoiceResponseData struct {
	RequestID    string           `json:"request_id"`
	Status       string           `json:"status"`
	Result       *aiInvoiceResult `json:"result"`
	ErrorCode    *string          `json:"error_code"`
	ErrorMessage *string          `json:"error_message"`
}

// aiServiceEnvelope Python AI 服务 ok() 包装层。
type aiServiceEnvelope struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    *aiInvoiceResponseData `json:"data"`
}

// =====================================================
// Service 接口
// =====================================================

// AIInvoiceService AI 发票业务接口。
type AIInvoiceService interface {
	// RecognizeInvoice 触发 AI 识别流程，返回识别记录。
	RecognizeInvoice(ctx context.Context, fileID, fileName string, fileBytes []byte, contentType string, supplierID *int64, remark string, creatorID int64) (*aiModel.AIInvoiceRecord, error)
	// GetInvoice 查询 AI 发票识别记录详情。
	GetInvoice(ctx context.Context, id int64) (*aiModel.AIInvoiceRecord, error)
	// ListInvoices 分页查询 AI 发票记录列表。
	ListInvoices(ctx context.Context, filter AIInvoiceListFilter) ([]*aiModel.AIInvoiceRecord, int64, error)
	// ConvertToInbound 将已识别完成的发票记录转换为草稿入库单（含明细）。
	ConvertToInbound(ctx context.Context, invoiceID int64, req ConvertToInboundReq, operatorID int64) (*inboundModel.InboundOrder, error)
}

// =====================================================
// Service 实现
// =====================================================

// aiInvoiceService 是 AIInvoiceService 的默认实现。
type aiInvoiceService struct {
	db           *gorm.DB
	invoiceRepo  repository.AIInvoiceRepo
	inboundSvc   inboundService.InboundService
	logger       *zap.Logger
	aiServiceURL string
	httpClient   *http.Client
}

// NewAIInvoiceService 创建 AI 发票服务实例。
func NewAIInvoiceService(
	db *gorm.DB,
	invoiceRepo repository.AIInvoiceRepo,
	_ inboundRepo.InboundRepo,
	inboundSvc inboundService.InboundService,
	logger *zap.Logger,
	aiServiceBaseURL string,
	aiServiceTimeout time.Duration,
) AIInvoiceService {
	if aiServiceTimeout <= 0 {
		aiServiceTimeout = 600 * time.Second
	}
	return &aiInvoiceService{
		db:           db,
		invoiceRepo:  invoiceRepo,
		inboundSvc:   inboundSvc,
		logger:       logger,
		aiServiceURL: aiServiceBaseURL,
		httpClient:   &http.Client{Timeout: aiServiceTimeout},
	}
}

// =====================================================
// RecognizeInvoice
// =====================================================

// RecognizeInvoice 创建识别记录并调用 Python AI 服务执行识别。
func (s *aiInvoiceService) RecognizeInvoice(
	ctx context.Context,
	fileID, fileName string,
	fileBytes []byte,
	contentType string,
	supplierID *int64,
	remark string,
	creatorID int64,
) (*aiModel.AIInvoiceRecord, error) {
	if fileID == "" {
		return nil, ecode.ErrParamInvalid
	}
	if creatorID <= 0 {
		return nil, ecode.ErrParamInvalid
	}

	// 1. 确保 file_info 中存在该 file_id（满足外键约束）。
	if err := s.db.WithContext(ctx).Exec(
		`INSERT INTO file_info (file_id, original_name, storage_name, storage_path, content_type, business_type, uploader_id, status)
		 VALUES (?, ?, ?, ?, ?, 'ai_invoice', ?, 1)
		 ON CONFLICT (file_id) DO NOTHING`,
		fileID, fileName, fileName, "/uploads/"+fileID, contentType, creatorID,
	).Error; err != nil {
		return nil, fmt.Errorf("写入文件记录失败: %w", err)
	}

	// 2. 创建 PENDING 识别记录。
	record := &aiModel.AIInvoiceRecord{
		FileID:            fileID,
		FileName:          fileName,
		Status:            aiModel.AIInvoiceStatusPending,
		CreatorID:         creatorID,
		MatchedSupplierID: supplierID,
	}
	if remark != "" {
		record.Remark = &remark
	}
	if err := s.invoiceRepo.Create(ctx, record); err != nil {
		return nil, err
	}

	// 3. 后台异步调用 Python AI 服务，立即返回 PROCESSING 状态给前端。
	// AI 识别耗时 2-5 分钟，同步等待会导致请求超时，改为后台 goroutine。
	record.Status = aiModel.AIInvoiceStatusProcessing
	go func() {
		bgCtx := context.Background()
		if err := s.callAIService(bgCtx, record, fileBytes, contentType, fileName); err != nil {
			s.logger.Warn("AI 发票识别失败", zap.Int64("invoice_id", record.ID), zap.Error(err))
			errMsg := err.Error()
			_ = s.invoiceRepo.UpdateStatus(bgCtx, record.ID, map[string]interface{}{
				"status":        aiModel.AIInvoiceStatusFailed,
				"error_message": errMsg,
			})
		}
	}()

	return record, nil
}

// callAIService 向 Python AI 服务 POST multipart 文件，解析并持久化识别结果。
func (s *aiInvoiceService) callAIService(
	ctx context.Context,
	record *aiModel.AIInvoiceRecord,
	fileBytes []byte,
	contentType string,
	fileName string,
) error {
	// 标记为处理中。
	if err := s.invoiceRepo.UpdateStatus(ctx, record.ID, map[string]interface{}{
		"status": aiModel.AIInvoiceStatusProcessing,
	}); err != nil {
		return fmt.Errorf("更新识别状态失败: %w", err)
	}
	record.Status = aiModel.AIInvoiceStatusProcessing

	// 构建 multipart/form-data 请求体。
	// 注意：必须手动设置文件 part 的 Content-Type，否则 Python AI 服务会收到
	// application/octet-stream 并以 415 拒绝请求。
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
	partHeader.Set("Content-Type", contentType)
	part, err := mw.CreatePart(partHeader)
	if err != nil {
		return fmt.Errorf("构建 multipart 字段失败: %w", err)
	}
	if _, err = part.Write(fileBytes); err != nil {
		return fmt.Errorf("写入文件内容失败: %w", err)
	}
	_ = mw.WriteField("erp_file_id", record.FileID)
	_ = mw.WriteField("match_master_data", "true")
	mw.Close()

	aiURL := s.aiServiceURL + "/ai/api/v1/invoices/recognize"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, aiURL, &body)
	if err != nil {
		return fmt.Errorf("构建 AI 请求失败: %w", err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("调用 AI 服务失败: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取 AI 服务响应失败: %w", err)
	}

	var envelope aiServiceEnvelope
	if err := json.Unmarshal(respBytes, &envelope); err != nil {
		return fmt.Errorf("解析 AI 响应 JSON 失败: %w", err)
	}

	if envelope.Code != 200 || envelope.Data == nil {
		return fmt.Errorf("AI 服务返回错误 [%d]: %s", envelope.Code, envelope.Message)
	}

	data := envelope.Data
	if data.Status == "FAILED" {
		errMsg := "AI 识别失败"
		if data.ErrorMessage != nil && *data.ErrorMessage != "" {
			errMsg = *data.ErrorMessage
		}
		return fmt.Errorf("%s", errMsg)
	}
	if data.Result == nil {
		return fmt.Errorf("AI 服务返回空识别结果")
	}

	result := data.Result

	// 将完整 result 作为 result_json 存储，raw_response 存原始响应。
	resultBytes, _ := json.Marshal(result)

	updates := map[string]interface{}{
		"status":            aiModel.AIInvoiceStatusCompleted,
		"result_json":       datatypes.JSON(resultBytes),
		"raw_response_json": datatypes.JSON(respBytes),
	}
	if result.RecognizedSupplierName != nil {
		updates["recognized_supplier_name"] = *result.RecognizedSupplierName
	}
	if result.InvoiceNo != nil {
		updates["invoice_no"] = *result.InvoiceNo
	}
	if result.InvoiceDate != nil {
		if d, err2 := time.Parse("2006-01-02", *result.InvoiceDate); err2 == nil {
			updates["invoice_date"] = d
		}
	}
	// 只在 AI 推荐且原始请求未指定时写入 matched_supplier_id。
	if result.MatchedSupplierID != nil && *result.MatchedSupplierID > 0 && record.MatchedSupplierID == nil {
		updates["matched_supplier_id"] = *result.MatchedSupplierID
	}

	if err := s.invoiceRepo.UpdateStatus(ctx, record.ID, updates); err != nil {
		return fmt.Errorf("写入识别结果失败: %w", err)
	}

	// 同步更新内存中的 record 字段。
	record.Status = aiModel.AIInvoiceStatusCompleted
	record.RecognizedSupplierName = result.RecognizedSupplierName
	record.InvoiceNo = result.InvoiceNo
	if result.InvoiceDate != nil {
		if d, err2 := time.Parse("2006-01-02", *result.InvoiceDate); err2 == nil {
			record.InvoiceDate = &d
		}
	}
	if result.MatchedSupplierID != nil && *result.MatchedSupplierID > 0 && record.MatchedSupplierID == nil {
		record.MatchedSupplierID = result.MatchedSupplierID
	}
	record.ResultJSON = datatypes.JSON(resultBytes)
	record.RawResponseJSON = datatypes.JSON(respBytes)

	s.logger.Info("AI 发票识别完成",
		zap.Int64("invoice_id", record.ID),
		zap.String("ai_request_id", data.RequestID),
		zap.Float64("confidence", result.Confidence),
	)
	return nil
}

// =====================================================
// GetInvoice / ListInvoices
// =====================================================

// GetInvoice 查询 AI 发票记录详情。
func (s *aiInvoiceService) GetInvoice(ctx context.Context, id int64) (*aiModel.AIInvoiceRecord, error) {
	record, err := s.invoiceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}
	s.enrichInvoiceFields(ctx, []*aiModel.AIInvoiceRecord{record})
	return record, nil
}

// ListInvoices 分页查询 AI 发票记录列表。
func (s *aiInvoiceService) ListInvoices(ctx context.Context, filter AIInvoiceListFilter) ([]*aiModel.AIInvoiceRecord, int64, error) {
	records, total, err := s.invoiceRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	s.enrichInvoiceFields(ctx, records)
	return records, total, nil
}

// enrichInvoiceFields 批量填充 supplier_name、drug_count、total_amount、confidence。
func (s *aiInvoiceService) enrichInvoiceFields(ctx context.Context, records []*aiModel.AIInvoiceRecord) {
	if len(records) == 0 {
		return
	}
	supplierIDs := make([]int64, 0)
	for _, r := range records {
		if r.MatchedSupplierID != nil && *r.MatchedSupplierID > 0 {
			supplierIDs = append(supplierIDs, *r.MatchedSupplierID)
		}
	}
	supplierNameMap := make(map[int64]string)
	if len(supplierIDs) > 0 {
		type row struct {
			ID   int64
			Name string
		}
		var rows []row
		s.db.WithContext(ctx).Table("supplier").Select("id, name").Where("id IN ?", supplierIDs).Scan(&rows)
		for _, r := range rows {
			supplierNameMap[r.ID] = r.Name
		}
	}

	for _, r := range records {
		if r.MatchedSupplierID != nil {
			if name, ok := supplierNameMap[*r.MatchedSupplierID]; ok && name != "" {
				r.SupplierName = name
			}
		}
		if r.SupplierName == "" && r.RecognizedSupplierName != nil {
			r.SupplierName = *r.RecognizedSupplierName
		}

		// 从 result_json 提取统计字段。
		// Python AI 服务返回 total_amount 为字符串（如 "1024.50"）。
		if len(r.ResultJSON) > 0 {
			var result struct {
				TotalAmount interface{}   `json:"total_amount"`
				Confidence  float64       `json:"confidence"`
				Items       []interface{} `json:"items"`
			}
			if err := json.Unmarshal(r.ResultJSON, &result); err == nil {
				switch v := result.TotalAmount.(type) {
				case float64:
					r.TotalAmount = v
				case string:
					r.TotalAmount, _ = strconv.ParseFloat(v, 64)
				}
				r.Confidence = result.Confidence
				r.DrugCount = len(result.Items)
			}
		}
	}
}

// =====================================================
// ConvertToInbound
// =====================================================

// ConvertToInbound 将 COMPLETED 状态的 AI 发票转换为草稿入库单。
func (s *aiInvoiceService) ConvertToInbound(ctx context.Context, invoiceID int64, req ConvertToInboundReq, operatorID int64) (*inboundModel.InboundOrder, error) {
	if req.SupplierID <= 0 {
		return nil, ecode.ErrParamInvalid
	}
	if len(req.Items) == 0 {
		return nil, ecode.New(20011, "转换入库单须至少有一条明细")
	}

	record, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrNotFound
		}
		return nil, err
	}

	if record.Status != aiModel.AIInvoiceStatusCompleted {
		return nil, ecode.New(20009, "只有识别完成的发票才能转换为入库单")
	}
	if record.InboundOrderID != nil {
		return nil, ecode.New(20010, "该发票已转换为入库单，不可重复转换")
	}

	details := make([]inboundService.CreateDetailReq, 0, len(req.Items))
	for _, item := range req.Items {
		details = append(details, inboundService.CreateDetailReq{
			DrugID:      item.DrugID,
			BatchNumber: item.BatchNumber,
			ExpireDate:  item.ExpireDate,
			PlannedQty:  item.PlannedQty,
			UnitPrice:   item.UnitPrice,
			Remark:      item.Remark,
		})
	}

	invoiceNo := ""
	if record.InvoiceNo != nil {
		invoiceNo = *record.InvoiceNo
	}

	var createdOrder *inboundModel.InboundOrder
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order, err := s.inboundSvc.CreateOrderWithDetails(ctx, tx, inboundService.CreateOrderReq{
			SupplierID: req.SupplierID,
			InvoiceNo:  invoiceNo,
			Details:    details,
		}, operatorID)
		if err != nil {
			return err
		}
		createdOrder = order

		return tx.WithContext(ctx).
			Model(&aiModel.AIInvoiceRecord{}).
			Where("id = ?", invoiceID).
			Updates(map[string]interface{}{
				"inbound_order_id": order.ID,
				"converted_at":     time.Now(),
			}).Error
	})
	if err != nil {
		return nil, err
	}

	return s.inboundSvc.GetOrder(ctx, createdOrder.ID)
}
