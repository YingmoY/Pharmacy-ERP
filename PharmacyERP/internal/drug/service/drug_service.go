package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/YingmoY/PharmacyERP/internal/drug/model"
	"github.com/YingmoY/PharmacyERP/internal/drug/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateDrugRequest 创建药品请求。
type CreateDrugRequest struct {
	DrugCode         string   `json:"drug_code"`
	CommonName       string   `json:"common_name"`
	TradeName        string   `json:"trade_name"`
	Specification    string   `json:"specification"`
	DosageForm       string   `json:"dosage_form"`
	Manufacturer     string   `json:"manufacturer"`
	ApprovalNumber   string   `json:"approval_number"`
	IsPrescription   bool     `json:"is_prescription"`
	IsMedicare       bool     `json:"is_medicare"`
	Barcode          *string  `json:"barcode"`
	Unit             string   `json:"unit"`
	RetailPrice      *float64 `json:"retail_price"`
	PurchasePrice    *float64 `json:"purchase_price"`
	StorageCondition string   `json:"storage_condition"`
	Remark           string   `json:"remark"`
	OperatorID       int64    `json:"-"`
}

// UpdateDrugRequest 更新药品请求。
type UpdateDrugRequest struct {
	CommonName       *string  `json:"common_name"`
	TradeName        *string  `json:"trade_name"`
	Specification    *string  `json:"specification"`
	DosageForm       *string  `json:"dosage_form"`
	Manufacturer     *string  `json:"manufacturer"`
	ApprovalNumber   *string  `json:"approval_number"`
	IsPrescription   *bool    `json:"is_prescription"`
	IsMedicare       *bool    `json:"is_medicare"`
	Barcode          *string  `json:"barcode"`
	Unit             *string  `json:"unit"`
	RetailPrice      *float64 `json:"retail_price"`
	PurchasePrice    *float64 `json:"purchase_price"`
	StorageCondition *string  `json:"storage_condition"`
	Remark           *string  `json:"remark"`
	OperatorID       int64    `json:"-"`
}

// DrugListRequest 药品列表查询请求。
type DrugListRequest struct {
	Keyword        string `form:"keyword"`
	Status         *int8  `form:"status"`
	IsPrescription *bool  `form:"is_prescription"`
	Manufacturer   string `form:"manufacturer"`
	IsMedicare     *bool  `form:"is_medicare"`
	Page           int    `form:"page"`
	PageSize       int    `form:"page_size"`
}

// DrugDTO 药品数据传输对象，用于 HTTP 响应。
type DrugDTO struct {
	ID               int64    `json:"id"`
	DrugCode         string   `json:"drug_code"`
	CommonName       string   `json:"common_name"`
	TradeName        string   `json:"trade_name"`
	Specification    string   `json:"specification"`
	DosageForm       string   `json:"dosage_form"`
	Manufacturer     string   `json:"manufacturer"`
	ApprovalNumber   string   `json:"approval_number"`
	IsPrescription   bool     `json:"is_prescription"`
	IsMedicare       bool     `json:"is_medicare"`
	Status           int8     `json:"status"`
	Barcode          *string  `json:"barcode,omitempty"`
	Unit             string   `json:"unit"`
	RetailPrice      *float64 `json:"retail_price,omitempty"`
	PurchasePrice    *float64 `json:"purchase_price,omitempty"`
	StorageCondition string   `json:"storage_condition"`
	Remark           string   `json:"remark"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	InStockCount     int64    `json:"in_stock_count"`
}

// InventorySummaryDTO 药品库存汇总数据传输对象。
type InventorySummaryDTO struct {
	DrugID             int64 `json:"drug_id"`
	InStockCount       int64 `json:"in_stock_count"`
	AvailableCount     int64 `json:"available_count"`
	PendingCount       int64 `json:"pending_count"`
	SoldCount          int64 `json:"sold_count"`
	LostCount          int64 `json:"lost_count"`
	MisplacedCount     int64 `json:"misplaced_count"`
	LossCandidateCount int64 `json:"loss_candidate_count"`
}

// DrugSaleInfoDTO 药品销售信息数据传输对象。
type DrugSaleInfoDTO struct {
	ID             int64    `json:"id"`
	DrugCode       string   `json:"drug_code"`
	CommonName     string   `json:"common_name"`
	TradeName      string   `json:"trade_name"`
	Specification  string   `json:"specification"`
	Manufacturer   string   `json:"manufacturer"`
	Unit           string   `json:"unit"`
	RetailPrice    *float64 `json:"retail_price,omitempty"`
	IsPrescription bool     `json:"is_prescription"`
	InStockCount   int64    `json:"in_stock_count"`
}

// DrugService 定义药品业务逻辑接口。
type DrugService interface {
	// CreateDrug 创建新药品，校验编码唯一性与必填字段。
	CreateDrug(ctx context.Context, req CreateDrugRequest) (*DrugDTO, error)
	// GetDrug 根据 ID 获取药品。
	GetDrug(ctx context.Context, id int64) (*DrugDTO, error)
	// GetDrugByCode 根据药品编码获取药品。
	GetDrugByCode(ctx context.Context, code string) (*DrugDTO, error)
	// ListDrugs 分页查询药品列表。
	ListDrugs(ctx context.Context, req DrugListRequest) ([]*DrugDTO, int64, error)
	// UpdateDrug 更新药品信息。
	UpdateDrug(ctx context.Context, id int64, req UpdateDrugRequest) (*DrugDTO, error)
	// DeleteDrug 软删除药品，删前检查关联引用。
	DeleteDrug(ctx context.Context, id int64, operatorID int64) error
	// UpdateDrugStatus 更新药品状态（0 停用 / 1 启用）。
	UpdateDrugStatus(ctx context.Context, id int64, status int8, operatorID int64) error
	// GetInventorySummary 获取药品库存各状态汇总。
	GetInventorySummary(ctx context.Context, drugID int64) (*InventorySummaryDTO, error)
	// GetDrugSaleInfo 获取药品销售关键信息。
	GetDrugSaleInfo(ctx context.Context, id int64) (*DrugSaleInfoDTO, error)
	// SearchDrugs 药品快速搜索（用于自动补全）。
	SearchDrugs(ctx context.Context, q string, limit int) ([]*DrugDTO, error)
}

type drugService struct {
	db       *gorm.DB
	repo     repository.DrugRepository
	log      *zap.Logger
	mqClient *mq.Client
}

// NewDrugService 创建药品服务实例。
func NewDrugService(db *gorm.DB, repo repository.DrugRepository, log *zap.Logger, mqClient *mq.Client) DrugService {
	return &drugService{
		db:       db,
		repo:     repo,
		log:      log,
		mqClient: mqClient,
	}
}

// toDTO 将 GORM 模型转换为 DTO。
func toDTO(m *model.DrugInfo) *DrugDTO {
	return &DrugDTO{
		ID:               m.ID,
		DrugCode:         m.DrugCode,
		CommonName:       m.CommonName,
		TradeName:        m.TradeName,
		Specification:    m.Specification,
		DosageForm:       m.DosageForm,
		Manufacturer:     m.Manufacturer,
		ApprovalNumber:   m.ApprovalNumber,
		IsPrescription:   m.IsPrescription,
		IsMedicare:       m.IsMedicare,
		Status:           m.Status,
		Barcode:          m.Barcode,
		Unit:             m.Unit,
		RetailPrice:      m.RetailPrice,
		PurchasePrice:    m.PurchasePrice,
		StorageCondition: m.StorageCondition,
		Remark:           m.Remark,
		CreatedAt:        m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		InStockCount:     m.InStockCount,
	}
}

// publishLogAsync 异步发送操作日志，不阻塞主流程。
func (s *drugService) publishLogAsync(ctx context.Context, action string, businessID string, operatorID int64, detail interface{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("publishLogAsync panic recovered", zap.Any("recover", r))
			}
		}()
		if err := s.mqClient.PublishLogEvent(ctx, mq.LogEvent{
			BusinessType: "drug",
			BusinessID:   businessID,
			Action:       action,
			OperatorID:   operatorID,
			Detail:       detail,
		}); err != nil {
			s.log.Warn("发布操作日志失败", zap.Error(err))
		}
	}()
}

// CreateDrug 创建新药品。
func (s *drugService) CreateDrug(ctx context.Context, req CreateDrugRequest) (*DrugDTO, error) {
	// 必填字段校验
	if strings.TrimSpace(req.DrugCode) == "" || strings.TrimSpace(req.CommonName) == "" ||
		strings.TrimSpace(req.Specification) == "" || strings.TrimSpace(req.Manufacturer) == "" {
		return nil, ecode.ErrParamInvalid
	}

	// 检查编码唯一性
	existing, err := s.repo.FindByCode(ctx, s.db, req.DrugCode)
	if err == nil && existing != nil {
		return nil, ecode.New(40901, "drug_code already exists")
	}

	drug := &model.DrugInfo{
		DrugCode:         req.DrugCode,
		CommonName:       req.CommonName,
		TradeName:        req.TradeName,
		Specification:    req.Specification,
		DosageForm:       req.DosageForm,
		Manufacturer:     req.Manufacturer,
		ApprovalNumber:   req.ApprovalNumber,
		IsPrescription:   req.IsPrescription,
		IsMedicare:       req.IsMedicare,
		Status:           model.DrugStatusEnabled,
		Barcode:          req.Barcode,
		Unit:             req.Unit,
		RetailPrice:      req.RetailPrice,
		PurchasePrice:    req.PurchasePrice,
		StorageCondition: req.StorageCondition,
		Remark:           req.Remark,
	}

	if err := s.repo.Create(ctx, s.db, drug); err != nil {
		return nil, err
	}

	s.publishLogAsync(ctx, "create", fmt.Sprintf("%d", drug.ID), req.OperatorID, drug)
	return toDTO(drug), nil
}

// GetDrug 根据 ID 获取药品。
func (s *drugService) GetDrug(ctx context.Context, id int64) (*DrugDTO, error) {
	drug, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	return toDTO(drug), nil
}

// GetDrugByCode 根据药品编码获取药品。
func (s *drugService) GetDrugByCode(ctx context.Context, code string) (*DrugDTO, error) {
	drug, err := s.repo.FindByCode(ctx, s.db, code)
	if err != nil {
		return nil, err
	}
	return toDTO(drug), nil
}

// ListDrugs 分页查询药品列表。
func (s *drugService) ListDrugs(ctx context.Context, req DrugListRequest) ([]*DrugDTO, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	filter := repository.DrugFilter{
		Keyword:        req.Keyword,
		Status:         req.Status,
		IsPrescription: req.IsPrescription,
		Manufacturer:   req.Manufacturer,
		IsMedicare:     req.IsMedicare,
		Page:           page,
		PageSize:       pageSize,
	}

	list, total, err := s.repo.List(ctx, s.db, filter)
	if err != nil {
		return nil, 0, err
	}

	// Batch-enrich in-stock counts from drug_trace_inventory.
	if len(list) > 0 {
		ids := make([]int64, 0, len(list))
		for _, m := range list {
			ids = append(ids, m.ID)
		}
		type stockRow struct {
			DrugID int64 `gorm:"column:drug_id"`
			Cnt    int64 `gorm:"column:cnt"`
		}
		var rows []stockRow
		if err2 := s.db.WithContext(ctx).
			Table("drug_trace_inventory").
			Select("drug_id, COUNT(*) AS cnt").
			Where("drug_id IN ? AND status = 'IN_STOCK'", ids).
			Group("drug_id").
			Scan(&rows).Error; err2 == nil {
			stockMap := make(map[int64]int64, len(rows))
			for _, r := range rows {
				stockMap[r.DrugID] = r.Cnt
			}
			for _, m := range list {
				m.InStockCount = stockMap[m.ID]
			}
		}
	}

	dtos := make([]*DrugDTO, 0, len(list))
	for _, m := range list {
		dtos = append(dtos, toDTO(m))
	}
	return dtos, total, nil
}

// UpdateDrug 更新药品信息。
func (s *drugService) UpdateDrug(ctx context.Context, id int64, req UpdateDrugRequest) (*DrugDTO, error) {
	// 确认药品存在
	if _, err := s.repo.FindByID(ctx, s.db, id); err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.CommonName != nil {
		updates["common_name"] = *req.CommonName
	}
	if req.TradeName != nil {
		updates["trade_name"] = *req.TradeName
	}
	if req.Specification != nil {
		updates["specification"] = *req.Specification
	}
	if req.DosageForm != nil {
		updates["dosage_form"] = *req.DosageForm
	}
	if req.Manufacturer != nil {
		updates["manufacturer"] = *req.Manufacturer
	}
	if req.ApprovalNumber != nil {
		updates["approval_number"] = *req.ApprovalNumber
	}
	if req.IsPrescription != nil {
		updates["is_prescription"] = *req.IsPrescription
	}
	if req.IsMedicare != nil {
		updates["is_medicare"] = *req.IsMedicare
	}
	if req.Barcode != nil {
		updates["barcode"] = *req.Barcode
	}
	if req.Unit != nil {
		updates["unit"] = *req.Unit
	}
	if req.RetailPrice != nil {
		updates["retail_price"] = *req.RetailPrice
	}
	if req.PurchasePrice != nil {
		updates["purchase_price"] = *req.PurchasePrice
	}
	if req.StorageCondition != nil {
		updates["storage_condition"] = *req.StorageCondition
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if len(updates) == 0 {
		return s.GetDrug(ctx, id)
	}

	if err := s.repo.Update(ctx, s.db, id, updates); err != nil {
		return nil, err
	}

	s.publishLogAsync(ctx, "update", fmt.Sprintf("%d", id), req.OperatorID, updates)

	return s.GetDrug(ctx, id)
}

// DeleteDrug 软删除药品，删前检查关联引用。
func (s *drugService) DeleteDrug(ctx context.Context, id int64, operatorID int64) error {
	// 确认药品存在
	if _, err := s.repo.FindByID(ctx, s.db, id); err != nil {
		return err
	}

	// 检查库存引用
	hasInv, err := s.repo.HasInventory(ctx, s.db, id)
	if err != nil {
		return err
	}
	if hasInv {
		return ecode.New(40901, "drug has active inventory records, cannot delete")
	}

	// 检查销售引用
	hasSales, err := s.repo.HasSalesItems(ctx, s.db, id)
	if err != nil {
		return err
	}
	if hasSales {
		return ecode.New(40901, "drug has sales order items, cannot delete")
	}

	// 检查入库引用
	hasInbound, err := s.repo.HasInboundDetails(ctx, s.db, id)
	if err != nil {
		return err
	}
	if hasInbound {
		return ecode.New(40901, "drug has inbound order details, cannot delete")
	}

	if err := s.repo.Delete(ctx, s.db, id); err != nil {
		return err
	}

	s.publishLogAsync(ctx, "delete", fmt.Sprintf("%d", id), operatorID, map[string]interface{}{"id": id})
	return nil
}

// UpdateDrugStatus 更新药品状态。
func (s *drugService) UpdateDrugStatus(ctx context.Context, id int64, status int8, operatorID int64) error {
	if status != model.DrugStatusEnabled && status != model.DrugStatusDisabled {
		return ecode.ErrParamInvalid
	}

	if _, err := s.repo.FindByID(ctx, s.db, id); err != nil {
		return err
	}

	if err := s.repo.UpdateStatus(ctx, s.db, id, status); err != nil {
		return err
	}

	s.publishLogAsync(ctx, "update_status", fmt.Sprintf("%d", id), operatorID,
		map[string]interface{}{"id": id, "status": status})
	return nil
}

// GetInventorySummary 查询药品库存各状态数量汇总。
func (s *drugService) GetInventorySummary(ctx context.Context, drugID int64) (*InventorySummaryDTO, error) {
	// 确认药品存在
	if _, err := s.repo.FindByID(ctx, s.db, drugID); err != nil {
		return nil, err
	}

	type statusCount struct {
		Status string
		Count  int64
	}

	var rows []statusCount
	if err := s.db.WithContext(ctx).
		Table("drug_trace_inventory").
		Select("status, COUNT(*) as count").
		Where("drug_id = ? AND deleted_at IS NULL", drugID).
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	summary := &InventorySummaryDTO{DrugID: drugID}
	for _, row := range rows {
		switch row.Status {
		case "IN_STOCK":
			summary.InStockCount = row.Count
		case "PENDING":
			summary.PendingCount = row.Count
		case "SOLD":
			summary.SoldCount = row.Count
		case "LOST":
			summary.LostCount = row.Count
		case "MISPLACED":
			summary.MisplacedCount = row.Count
		case "LOSS_CANDIDATE":
			summary.LossCandidateCount = row.Count
		}
	}

	// available = in_stock（当前版本未实现预留逻辑，与 in_stock 一致）
	summary.AvailableCount = summary.InStockCount

	return summary, nil
}

// GetDrugSaleInfo 获取药品销售信息。
func (s *drugService) GetDrugSaleInfo(ctx context.Context, id int64) (*DrugSaleInfoDTO, error) {
	drug, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	var inStockCount int64
	if err := s.db.WithContext(ctx).
		Table("drug_trace_inventory").
		Where("drug_id = ? AND status = ? AND deleted_at IS NULL", id, "IN_STOCK").
		Count(&inStockCount).Error; err != nil {
		return nil, err
	}
	return &DrugSaleInfoDTO{
		ID:             drug.ID,
		DrugCode:       drug.DrugCode,
		CommonName:     drug.CommonName,
		TradeName:      drug.TradeName,
		Specification:  drug.Specification,
		Manufacturer:   drug.Manufacturer,
		Unit:           drug.Unit,
		RetailPrice:    drug.RetailPrice,
		IsPrescription: drug.IsPrescription,
		InStockCount:   inStockCount,
	}, nil
}

// SearchDrugs 药品快速搜索，用于销售前台自动补全。
func (s *drugService) SearchDrugs(ctx context.Context, q string, limit int) ([]*DrugDTO, error) {
	if strings.TrimSpace(q) == "" {
		return []*DrugDTO{}, nil
	}
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	filter := repository.DrugFilter{
		Keyword:  q,
		Page:     1,
		PageSize: limit,
	}

	// 只搜索启用状态的药品
	enabled := model.DrugStatusEnabled
	filter.Status = &enabled

	list, _, err := s.repo.List(ctx, s.db, filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]*DrugDTO, 0, len(list))
	for _, m := range list {
		dtos = append(dtos, toDTO(m))
	}
	return dtos, nil
}
