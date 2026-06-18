package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/YingmoY/PharmacyERP/internal/supplier/model"
	"github.com/YingmoY/PharmacyERP/internal/supplier/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateSupplierRequest 创建供应商请求。
type CreateSupplierRequest struct {
	SupplierCode string  `json:"supplier_code"`
	Name         string  `json:"name"`
	ContactName  string  `json:"contact_name"`
	ContactPhone string  `json:"contact_phone"`
	LicenseNo    string  `json:"license_no"`
	Address      string  `json:"address"`
	Remark       *string `json:"remark"`
	OperatorID   int64   `json:"-"`
}

// UpdateSupplierRequest 更新供应商请求（所有字段可选）。
type UpdateSupplierRequest struct {
	Name         *string `json:"name"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *string `json:"contact_phone"`
	LicenseNo    *string `json:"license_no"`
	Address      *string `json:"address"`
	Remark       *string `json:"remark"`
	OperatorID   int64   `json:"-"`
}

// SupplierListRequest 供应商列表查询请求。
type SupplierListRequest struct {
	Keyword  string `form:"keyword"`
	Status   *int8  `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// SupplierDTO 供应商数据传输对象。
type SupplierDTO struct {
	ID           int64   `json:"id"`
	SupplierCode string  `json:"supplier_code"`
	Name         string  `json:"name"`
	ContactName  string  `json:"contact_name"`
	ContactPhone string  `json:"contact_phone"`
	LicenseNo    string  `json:"license_no"`
	Address      string  `json:"address"`
	Status       int8    `json:"status"`
	Remark       *string `json:"remark,omitempty"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// SupplierService 定义供应商业务逻辑接口。
type SupplierService interface {
	// CreateSupplier 创建供应商，校验编码唯一性与必填字段。
	CreateSupplier(ctx context.Context, req CreateSupplierRequest) (*SupplierDTO, error)
	// GetSupplier 根据 ID 获取供应商。
	GetSupplier(ctx context.Context, id int64) (*SupplierDTO, error)
	// ListSuppliers 分页查询供应商列表。
	ListSuppliers(ctx context.Context, req SupplierListRequest) ([]*SupplierDTO, int64, error)
	// UpdateSupplier 更新供应商信息（不允许修改 supplier_code）。
	UpdateSupplier(ctx context.Context, id int64, req UpdateSupplierRequest) (*SupplierDTO, error)
	// DeleteSupplier 软删除供应商，删前检查关联入库单。
	DeleteSupplier(ctx context.Context, id int64, operatorID int64) error
	// UpdateSupplierStatus 更新供应商状态（0 停用 / 1 启用）。
	UpdateSupplierStatus(ctx context.Context, id int64, status int8, operatorID int64) error
}

type supplierService struct {
	db       *gorm.DB
	repo     repository.SupplierRepository
	log      *zap.Logger
	mqClient *mq.Client
}

// NewSupplierService 创建供应商服务实例。
func NewSupplierService(db *gorm.DB, repo repository.SupplierRepository, log *zap.Logger, mqClient *mq.Client) SupplierService {
	return &supplierService{
		db:       db,
		repo:     repo,
		log:      log,
		mqClient: mqClient,
	}
}

// toDTO 将 GORM 模型转换为 DTO。
func toDTO(m *model.Supplier) *SupplierDTO {
	return &SupplierDTO{
		ID:           m.ID,
		SupplierCode: m.SupplierCode,
		Name:         m.Name,
		ContactName:  m.ContactName,
		ContactPhone: m.ContactPhone,
		LicenseNo:    m.LicenseNo,
		Address:      m.Address,
		Status:       m.Status,
		Remark:       m.Remark,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// publishLogAsync 异步发送操作日志，不阻塞主流程。
func (s *supplierService) publishLogAsync(ctx context.Context, action string, businessID string, operatorID int64, detail interface{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("publishLogAsync panic recovered", zap.Any("recover", r))
			}
		}()
		if err := s.mqClient.PublishLogEvent(ctx, mq.LogEvent{
			BusinessType: "supplier",
			BusinessID:   businessID,
			Action:       action,
			OperatorID:   operatorID,
			Detail:       detail,
		}); err != nil {
			s.log.Warn("发布操作日志失败", zap.Error(err))
		}
	}()
}

// CreateSupplier 创建供应商。
func (s *supplierService) CreateSupplier(ctx context.Context, req CreateSupplierRequest) (*SupplierDTO, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, ecode.ErrParamInvalid
	}
	if strings.TrimSpace(req.SupplierCode) == "" {
		req.SupplierCode = fmt.Sprintf("SUP%s", time.Now().Format("20060102150405"))
	}

	// 检查编码唯一性
	existing, err := s.repo.FindByCode(ctx, s.db, req.SupplierCode)
	if err == nil && existing != nil {
		return nil, ecode.New(40901, "supplier_code already exists")
	}

	supplier := &model.Supplier{
		SupplierCode: req.SupplierCode,
		Name:         req.Name,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		LicenseNo:    req.LicenseNo,
		Address:      req.Address,
		Status:       model.SupplierStatusEnabled,
		Remark:       req.Remark,
	}

	if err := s.repo.Create(ctx, s.db, supplier); err != nil {
		return nil, err
	}

	s.publishLogAsync(ctx, "create", fmt.Sprintf("%d", supplier.ID), req.OperatorID, supplier)
	return toDTO(supplier), nil
}

// GetSupplier 根据 ID 获取供应商。
func (s *supplierService) GetSupplier(ctx context.Context, id int64) (*SupplierDTO, error) {
	supplier, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	return toDTO(supplier), nil
}

// ListSuppliers 分页查询供应商列表。
func (s *supplierService) ListSuppliers(ctx context.Context, req SupplierListRequest) ([]*SupplierDTO, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	filter := repository.SupplierFilter{
		Keyword:  req.Keyword,
		Status:   req.Status,
		Page:     page,
		PageSize: pageSize,
	}

	list, total, err := s.repo.List(ctx, s.db, filter)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]*SupplierDTO, 0, len(list))
	for _, m := range list {
		dtos = append(dtos, toDTO(m))
	}
	return dtos, total, nil
}

// UpdateSupplier 更新供应商信息（不允许修改 supplier_code）。
func (s *supplierService) UpdateSupplier(ctx context.Context, id int64, req UpdateSupplierRequest) (*SupplierDTO, error) {
	if _, err := s.repo.FindByID(ctx, s.db, id); err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, ecode.ErrParamInvalid
		}
		updates["name"] = *req.Name
	}
	if req.ContactName != nil {
		updates["contact_name"] = *req.ContactName
	}
	if req.ContactPhone != nil {
		updates["contact_phone"] = *req.ContactPhone
	}
	if req.LicenseNo != nil {
		updates["license_no"] = *req.LicenseNo
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if len(updates) > 0 {
		if err := s.repo.Update(ctx, s.db, id, updates); err != nil {
			return nil, err
		}
		s.publishLogAsync(ctx, "update", fmt.Sprintf("%d", id), req.OperatorID, updates)
	}

	return s.GetSupplier(ctx, id)
}

// DeleteSupplier 软删除供应商，删前检查关联入库单。
func (s *supplierService) DeleteSupplier(ctx context.Context, id int64, operatorID int64) error {
	if _, err := s.repo.FindByID(ctx, s.db, id); err != nil {
		return err
	}

	hasOrders, err := s.repo.HasInboundOrders(ctx, s.db, id)
	if err != nil {
		return err
	}
	if hasOrders {
		return ecode.New(40901, "supplier has inbound orders, cannot delete")
	}

	if err := s.repo.Delete(ctx, s.db, id); err != nil {
		return err
	}

	s.publishLogAsync(ctx, "delete", fmt.Sprintf("%d", id), operatorID, map[string]interface{}{"id": id})
	return nil
}

// UpdateSupplierStatus 更新供应商状态。
func (s *supplierService) UpdateSupplierStatus(ctx context.Context, id int64, status int8, operatorID int64) error {
	if status != model.SupplierStatusEnabled && status != model.SupplierStatusDisabled {
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
