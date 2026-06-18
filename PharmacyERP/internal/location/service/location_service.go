package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/location/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateLocationRequest 创建货位请求。
type CreateLocationRequest struct {
	LocationCode string  `json:"location_code"`
	LocationName string  `json:"location_name"`
	Area         string  `json:"area"`
	Shelf        string  `json:"shelf"`
	Layer        string  `json:"layer"`
	Position     string  `json:"position"`
	Capacity     int     `json:"capacity"`
	Remark       *string `json:"remark"`
	OperatorID   int64   `json:"-"`
}

// UpdateLocationRequest 更新货位请求（所有字段可选）。
type UpdateLocationRequest struct {
	LocationName *string `json:"location_name"`
	Area         *string `json:"area"`
	Shelf        *string `json:"shelf"`
	Layer        *string `json:"layer"`
	Position     *string `json:"position"`
	Capacity     *int    `json:"capacity"`
	Remark       *string `json:"remark"`
	OperatorID   int64   `json:"-"`
}

// LocationListRequest 货位列表查询请求。
type LocationListRequest struct {
	Keyword  string `form:"keyword"`
	Status   *int8  `form:"status"`
	Area     string `form:"area"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// LocationDTO 货位数据传输对象。
type LocationDTO struct {
	ID           int64   `json:"id"`
	LocationCode string  `json:"location_code"`
	LocationName string  `json:"location_name"`
	Area         string  `json:"area"`
	Shelf        string  `json:"shelf"`
	Layer        string  `json:"layer"`
	Position     string  `json:"position"`
	Capacity     int     `json:"capacity"`
	Status       int8    `json:"status"`
	Remark       *string `json:"remark,omitempty"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// LocationDrugDTO 货位药品列表项 DTO（在库药品汇总）。
type LocationDrugDTO struct {
	DrugID     int64  `json:"drug_id"`
	CommonName string `json:"common_name"`
	DrugCode   string `json:"drug_code"`
	Count      int64  `json:"count"`
}

// LocationService 定义货位业务逻辑接口。
type LocationService interface {
	// CreateLocation 创建货位，校验编码唯一性与必填字段。
	CreateLocation(ctx context.Context, req CreateLocationRequest) (*LocationDTO, error)
	// GetLocation 根据 ID 获取货位。
	GetLocation(ctx context.Context, id int64) (*LocationDTO, error)
	// GetLocationByCode 根据货位编码获取货位。
	GetLocationByCode(ctx context.Context, code string) (*LocationDTO, error)
	// ListLocations 分页查询货位列表。
	ListLocations(ctx context.Context, req LocationListRequest) ([]*LocationDTO, int64, error)
	// UpdateLocation 更新货位信息。
	UpdateLocation(ctx context.Context, id int64, req UpdateLocationRequest) (*LocationDTO, error)
	// DeleteLocation 软删除货位，删前检查活跃库存。
	DeleteLocation(ctx context.Context, id int64, operatorID int64) error
	// UpdateLocationStatus 更新货位状态（0 停用 / 1 启用）。
	UpdateLocationStatus(ctx context.Context, id int64, status int8, operatorID int64) error
	// GetLocationDrugs 获取货位当前在库的药品列表（按药品聚合）。
	GetLocationDrugs(ctx context.Context, locationID int64) ([]*LocationDrugDTO, error)
	// GetAreas 获取所有不重复的区域列表。
	GetAreas(ctx context.Context) ([]string, error)
}

type locationService struct {
	db       *gorm.DB
	repo     repository.LocationRepository
	log      *zap.Logger
	mqClient *mq.Client
}

// NewLocationService 创建货位服务实例。
func NewLocationService(db *gorm.DB, repo repository.LocationRepository, log *zap.Logger, mqClient *mq.Client) LocationService {
	return &locationService{
		db:       db,
		repo:     repo,
		log:      log,
		mqClient: mqClient,
	}
}

// toDTO 将 GORM 模型转换为 DTO。
func toDTO(m *model.LocationInfo) *LocationDTO {
	return &LocationDTO{
		ID:           m.ID,
		LocationCode: m.LocationCode,
		LocationName: m.LocationName,
		Area:         m.Area,
		Shelf:        m.Shelf,
		Layer:        m.Layer,
		Position:     m.Position,
		Capacity:     m.Capacity,
		Status:       m.Status,
		Remark:       m.Remark,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// publishLogAsync 异步发送操作日志，不阻塞主流程。
func (s *locationService) publishLogAsync(ctx context.Context, action string, businessID string, operatorID int64, detail interface{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.log.Error("publishLogAsync panic recovered", zap.Any("recover", r))
			}
		}()
		if err := s.mqClient.PublishLogEvent(ctx, mq.LogEvent{
			BusinessType: "location",
			BusinessID:   businessID,
			Action:       action,
			OperatorID:   operatorID,
			Detail:       detail,
		}); err != nil {
			s.log.Warn("发布操作日志失败", zap.Error(err))
		}
	}()
}

// CreateLocation 创建货位。
func (s *locationService) CreateLocation(ctx context.Context, req CreateLocationRequest) (*LocationDTO, error) {
	if strings.TrimSpace(req.LocationCode) == "" || strings.TrimSpace(req.LocationName) == "" ||
		strings.TrimSpace(req.Area) == "" {
		return nil, ecode.ErrParamInvalid
	}

	// 检查编码唯一性
	existing, err := s.repo.FindByCode(ctx, s.db, req.LocationCode)
	if err == nil && existing != nil {
		return nil, ecode.New(40901, "location_code already exists")
	}

	loc := &model.LocationInfo{
		LocationCode: req.LocationCode,
		LocationName: req.LocationName,
		Area:         req.Area,
		Shelf:        req.Shelf,
		Layer:        req.Layer,
		Position:     req.Position,
		Capacity:     req.Capacity,
		Status:       model.LocationStatusEnabled,
		Remark:       req.Remark,
	}

	if err := s.repo.Create(ctx, s.db, loc); err != nil {
		return nil, err
	}

	s.publishLogAsync(ctx, "create", fmt.Sprintf("%d", loc.ID), req.OperatorID, loc)
	return toDTO(loc), nil
}

// GetLocation 根据 ID 获取货位。
func (s *locationService) GetLocation(ctx context.Context, id int64) (*LocationDTO, error) {
	loc, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	return toDTO(loc), nil
}

// GetLocationByCode 根据货位编码获取货位。
func (s *locationService) GetLocationByCode(ctx context.Context, code string) (*LocationDTO, error) {
	loc, err := s.repo.FindByCode(ctx, s.db, code)
	if err != nil {
		return nil, err
	}
	return toDTO(loc), nil
}

// ListLocations 分页查询货位列表。
func (s *locationService) ListLocations(ctx context.Context, req LocationListRequest) ([]*LocationDTO, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	filter := repository.LocationFilter{
		Keyword:  req.Keyword,
		Status:   req.Status,
		Area:     req.Area,
		Page:     page,
		PageSize: pageSize,
	}

	list, total, err := s.repo.List(ctx, s.db, filter)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]*LocationDTO, 0, len(list))
	for _, m := range list {
		dtos = append(dtos, toDTO(m))
	}
	return dtos, total, nil
}

// UpdateLocation 更新货位信息。
func (s *locationService) UpdateLocation(ctx context.Context, id int64, req UpdateLocationRequest) (*LocationDTO, error) {
	if _, err := s.repo.FindByID(ctx, s.db, id); err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}

	if req.LocationName != nil {
		if strings.TrimSpace(*req.LocationName) == "" {
			return nil, ecode.ErrParamInvalid
		}
		updates["location_name"] = *req.LocationName
	}
	if req.Area != nil {
		if strings.TrimSpace(*req.Area) == "" {
			return nil, ecode.ErrParamInvalid
		}
		updates["area"] = *req.Area
	}
	if req.Shelf != nil {
		updates["shelf"] = *req.Shelf
	}
	if req.Layer != nil {
		updates["layer"] = *req.Layer
	}
	if req.Position != nil {
		updates["position"] = *req.Position
	}
	if req.Capacity != nil {
		updates["capacity"] = *req.Capacity
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

	return s.GetLocation(ctx, id)
}

// DeleteLocation 软删除货位，删前检查活跃库存。
func (s *locationService) DeleteLocation(ctx context.Context, id int64, operatorID int64) error {
	if _, err := s.repo.FindByID(ctx, s.db, id); err != nil {
		return err
	}

	hasInv, err := s.repo.HasActiveInventory(ctx, s.db, id)
	if err != nil {
		return err
	}
	if hasInv {
		return ecode.New(40901, "location has active inventory (IN_STOCK/MISPLACED/LOSS_CANDIDATE), cannot delete")
	}

	if err := s.repo.Delete(ctx, s.db, id); err != nil {
		return err
	}

	s.publishLogAsync(ctx, "delete", fmt.Sprintf("%d", id), operatorID, map[string]interface{}{"id": id})
	return nil
}

// UpdateLocationStatus 更新货位状态。
func (s *locationService) UpdateLocationStatus(ctx context.Context, id int64, status int8, operatorID int64) error {
	if status != model.LocationStatusEnabled && status != model.LocationStatusDisabled {
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

// GetLocationDrugs 获取货位当前在库的药品汇总列表。
// 通过 drug_trace_inventory 关联 drug_info 进行聚合统计。
func (s *locationService) GetLocationDrugs(ctx context.Context, locationID int64) ([]*LocationDrugDTO, error) {
	// 确认货位存在
	if _, err := s.repo.FindByID(ctx, s.db, locationID); err != nil {
		return nil, err
	}

	type row struct {
		DrugID     int64  `gorm:"column:drug_id"`
		CommonName string `gorm:"column:common_name"`
		DrugCode   string `gorm:"column:drug_code"`
		Count      int64  `gorm:"column:count"`
	}

	var rows []row
	if err := s.db.WithContext(ctx).
		Table("drug_trace_inventory ti").
		Select("ti.drug_id, di.common_name, di.drug_code, COUNT(*) as count").
		Joins("LEFT JOIN drug_info di ON di.id = ti.drug_id AND di.deleted_at IS NULL").
		Where("ti.location_id = ? AND ti.status = ? AND ti.deleted_at IS NULL",
			locationID, "IN_STOCK").
		Group("ti.drug_id, di.common_name, di.drug_code").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	dtos := make([]*LocationDrugDTO, 0, len(rows))
	for _, r := range rows {
		dtos = append(dtos, &LocationDrugDTO{
			DrugID:     r.DrugID,
			CommonName: r.CommonName,
			DrugCode:   r.DrugCode,
			Count:      r.Count,
		})
	}
	return dtos, nil
}

// GetAreas 获取所有不重复的区域列表。
func (s *locationService) GetAreas(ctx context.Context) ([]string, error) {
	return s.repo.GetAreas(ctx, s.db)
}
