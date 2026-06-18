// Package service 实现告警模块的业务逻辑层。
package service

import (
	"context"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/alert/model"
	"github.com/YingmoY/PharmacyERP/internal/alert/repository"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AlertFilter 告警列表服务层过滤条件。
type AlertFilter struct {
	// Status 按状态筛选（"ACTIVE"/"RESOLVED"/"IGNORED"，空字符串表示不限）。
	Status string
	// EventType 按告警类型筛选。
	EventType string
	// Severity 按严重等级筛选。
	Severity string
	// Page 页码（从 1 开始）。
	Page int
	// PageSize 每页条数。
	PageSize int
}

// AlertInfo 告警详情 DTO。
type AlertInfo struct {
	ID          int64      `json:"id"`
	EventType   string     `json:"event_type"`
	RelatedType string     `json:"related_type"`
	RelatedID   string     `json:"related_id"`
	Description string     `json:"description"`
	// Status 状态字符串：ACTIVE/RESOLVED/IGNORED。
	Status      string     `json:"status"`
	Severity    string     `json:"severity"`
	Resolution  string     `json:"resolution"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	IgnoredAt   *time.Time `json:"ignored_at,omitempty"`
	AssignedTo  *int64     `json:"assigned_to,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// NearExpireInfo 近效期药品信息 DTO。
type NearExpireInfo struct {
	TraceCode       string    `json:"trace_code"`
	DrugID          int64     `json:"drug_id"`
	DrugName        string    `json:"drug_name"`
	BatchNumber     string    `json:"batch_number"`
	ExpireDate      time.Time `json:"expire_date"`
	LocationCode    string    `json:"location_code"`
	DaysUntilExpire int       `json:"days_until_expire"`
	Severity        string    `json:"severity"`
}

// LossCandidateInfo 盘亏候选药品信息 DTO。
type LossCandidateInfo struct {
	TraceCode    string `json:"trace_code"`
	DrugID       int64  `json:"drug_id"`
	DrugName     string `json:"drug_name"`
	BatchNumber  string `json:"batch_number"`
	LocationCode string `json:"location_code"`
	Status       string `json:"status"`
}

// AlertService 告警服务接口。
type AlertService interface {
	// ListAlerts 分页查询告警列表。
	ListAlerts(ctx context.Context, filter AlertFilter) ([]*AlertInfo, int64, error)
	// GetAlert 按 ID 查询告警详情。
	GetAlert(ctx context.Context, id int64) (*AlertInfo, error)
	// ResolveAlert 解决告警。
	ResolveAlert(ctx context.Context, id int64, resolvedBy int64, resolution string) error
	// IgnoreAlert 忽略告警。
	IgnoreAlert(ctx context.Context, id int64, ignoredBy int64, reason string) error
	// ListNearExpire 实时查询近效期药品列表（直接查询 drug_trace_inventory）。
	ListNearExpire(ctx context.Context, days, page, pageSize int) ([]*NearExpireInfo, int64, error)
	// ListLossCandidates 实时查询盘亏候选药品列表（直接查询 drug_trace_inventory）。
	ListLossCandidates(ctx context.Context, page, pageSize int) ([]*LossCandidateInfo, int64, error)
}

type alertService struct {
	repo repository.AlertRepo
	db   *gorm.DB
	log  *zap.Logger
}

// NewAlertService 创建告警服务实例。
func NewAlertService(repo repository.AlertRepo, db *gorm.DB, log *zap.Logger) AlertService {
	return &alertService{repo: repo, db: db, log: log}
}

// statusToInt16 将状态字符串转换为数据库整数值。
func statusToInt16(s string) *int16 {
	switch s {
	case "ACTIVE":
		v := model.AlertStatusActive
		return &v
	case "RESOLVED":
		v := model.AlertStatusResolved
		return &v
	case "IGNORED":
		v := model.AlertStatusIgnored
		return &v
	default:
		return nil
	}
}

// int16ToStatus 将数据库整数状态值转换为字符串。
func int16ToStatus(s int16) string {
	switch s {
	case model.AlertStatusActive:
		return "ACTIVE"
	case model.AlertStatusResolved:
		return "RESOLVED"
	case model.AlertStatusIgnored:
		return "IGNORED"
	default:
		return "UNKNOWN"
	}
}

// toAlertInfo 将数据库模型转换为 DTO。
func toAlertInfo(a *model.Alert) *AlertInfo {
	info := &AlertInfo{
		ID:          a.ID,
		EventType:   a.EventType,
		RelatedType: a.RelatedType,
		RelatedID:   a.RelatedID,
		Status:      int16ToStatus(a.Status),
		Severity:    a.Severity,
		ClosedAt:    a.ClosedAt,
		IgnoredAt:   a.IgnoredAt,
		AssignedTo:  a.AssignedTo,
		CreatedAt:   a.CreatedAt,
	}
	if a.Description != nil {
		info.Description = *a.Description
	}
	if a.Resolution != nil {
		info.Resolution = *a.Resolution
	}
	return info
}

// ListAlerts 分页查询告警列表。
func (s *alertService) ListAlerts(ctx context.Context, filter AlertFilter) ([]*AlertInfo, int64, error) {
	repoFilter := repository.AlertFilter{
		Status:    statusToInt16(filter.Status),
		EventType: filter.EventType,
		Severity:  filter.Severity,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
	}

	alerts, total, err := s.repo.List(ctx, repoFilter)
	if err != nil {
		s.log.Error("查询告警列表失败", zap.Error(err))
		return nil, 0, ecode.ErrSystem
	}

	result := make([]*AlertInfo, 0, len(alerts))
	for _, a := range alerts {
		result = append(result, toAlertInfo(a))
	}
	return result, total, nil
}

// GetAlert 按 ID 查询告警详情。
func (s *alertService) GetAlert(ctx context.Context, id int64) (*AlertInfo, error) {
	alert, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ecode.ErrNotFound
		}
		s.log.Error("查询告警详情失败", zap.Int64("id", id), zap.Error(err))
		return nil, ecode.ErrSystem
	}
	return toAlertInfo(alert), nil
}

// ResolveAlert 将告警标记为已解决。
func (s *alertService) ResolveAlert(ctx context.Context, id int64, resolvedBy int64, resolution string) error {
	// 先查询是否存在
	alert, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ecode.ErrNotFound
		}
		return ecode.ErrSystem
	}
	// 已解决或已忽略的告警不允许重复处理
	if alert.Status != model.AlertStatusActive {
		return ecode.New(10005, "告警已处理，无法重复操作")
	}

	now := time.Now()
	if err := s.repo.UpdateResolve(ctx, id, resolvedBy, resolution, now); err != nil {
		s.log.Error("解决告警失败", zap.Int64("id", id), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}

// IgnoreAlert 将告警标记为已忽略。
func (s *alertService) IgnoreAlert(ctx context.Context, id int64, ignoredBy int64, reason string) error {
	// 先查询是否存在
	alert, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ecode.ErrNotFound
		}
		return ecode.ErrSystem
	}
	// 已处理的告警不允许重复操作
	if alert.Status != model.AlertStatusActive {
		return ecode.New(10005, "告警已处理，无法重复操作")
	}

	now := time.Now()
	if err := s.repo.UpdateIgnore(ctx, id, ignoredBy, now, reason); err != nil {
		s.log.Error("忽略告警失败", zap.Int64("id", id), zap.Error(err))
		return ecode.ErrSystem
	}
	return nil
}

// nearExpireRow 用于接收近效期查询结果。
type nearExpireRow struct {
	TraceCode    string    `gorm:"column:trace_code"`
	DrugID       int64     `gorm:"column:drug_id"`
	DrugName     string    `gorm:"column:drug_name"`
	BatchNumber  string    `gorm:"column:batch_number"`
	ExpireDate   time.Time `gorm:"column:expire_date"`
	LocationCode string    `gorm:"column:location_code"`
}

// ListNearExpire 实时查询近效期药品（`days` 天内到期，状态为 IN_STOCK）。
func (s *alertService) ListNearExpire(ctx context.Context, days, page, pageSize int) ([]*NearExpireInfo, int64, error) {
	if days <= 0 {
		days = 30
	}
	now := time.Now()
	threshold := now.AddDate(0, 0, days)

	// 统计总数
	var total int64
	if err := s.db.WithContext(ctx).
		Table("drug_trace_inventory dti").
		Where("dti.status = 'IN_STOCK'").
		Where("dti.expire_date <= ?", threshold).
		Where("dti.deleted_at IS NULL").
		Count(&total).Error; err != nil {
		s.log.Error("统计近效期药品数量失败", zap.Error(err))
		return nil, 0, ecode.ErrSystem
	}

	if total == 0 {
		return []*NearExpireInfo{}, 0, nil
	}

	offset := (page - 1) * pageSize
	var rows []nearExpireRow

	// 联查药品名称和货位编码
	if err := s.db.WithContext(ctx).
		Table("drug_trace_inventory dti").
		Select("dti.trace_code, dti.drug_id, di.common_name AS drug_name, dti.batch_number, dti.expire_date, COALESCE(li.location_code, '') AS location_code").
		Joins("LEFT JOIN drug_info di ON di.id = dti.drug_id AND di.deleted_at IS NULL").
		Joins("LEFT JOIN location_info li ON li.id = dti.location_id AND li.deleted_at IS NULL").
		Where("dti.status = 'IN_STOCK'").
		Where("dti.expire_date <= ?", threshold).
		Where("dti.deleted_at IS NULL").
		Order("dti.expire_date ASC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		s.log.Error("查询近效期药品失败", zap.Error(err))
		return nil, 0, ecode.ErrSystem
	}

	result := make([]*NearExpireInfo, 0, len(rows))
	for _, row := range rows {
		// 计算距离到期天数
		daysUntil := int(row.ExpireDate.Sub(now).Hours() / 24)

		// 根据天数判断严重等级
		severity := model.AlertSeverityLow
		if daysUntil <= 7 {
			severity = model.AlertSeverityHigh
		} else if daysUntil <= 15 {
			severity = model.AlertSeverityMedium
		}

		result = append(result, &NearExpireInfo{
			TraceCode:       row.TraceCode,
			DrugID:          row.DrugID,
			DrugName:        row.DrugName,
			BatchNumber:     row.BatchNumber,
			ExpireDate:      row.ExpireDate,
			LocationCode:    row.LocationCode,
			DaysUntilExpire: daysUntil,
			Severity:        severity,
		})
	}
	return result, total, nil
}

// lossCandidateRow 用于接收盘亏候选查询结果。
type lossCandidateRow struct {
	TraceCode    string `gorm:"column:trace_code"`
	DrugID       int64  `gorm:"column:drug_id"`
	DrugName     string `gorm:"column:drug_name"`
	BatchNumber  string `gorm:"column:batch_number"`
	LocationCode string `gorm:"column:location_code"`
	Status       string `gorm:"column:status"`
}

// ListLossCandidates 实时查询盘亏候选药品列表。
func (s *alertService) ListLossCandidates(ctx context.Context, page, pageSize int) ([]*LossCandidateInfo, int64, error) {
	// 统计总数
	var total int64
	if err := s.db.WithContext(ctx).
		Table("drug_trace_inventory dti").
		Where("dti.status = 'LOSS_CANDIDATE'").
		Where("dti.deleted_at IS NULL").
		Count(&total).Error; err != nil {
		s.log.Error("统计盘亏候选数量失败", zap.Error(err))
		return nil, 0, ecode.ErrSystem
	}

	if total == 0 {
		return []*LossCandidateInfo{}, 0, nil
	}

	offset := (page - 1) * pageSize
	var rows []lossCandidateRow

	// 联查药品名称和货位编码
	if err := s.db.WithContext(ctx).
		Table("drug_trace_inventory dti").
		Select("dti.trace_code, dti.drug_id, di.common_name AS drug_name, dti.batch_number, COALESCE(li.location_code, '') AS location_code, dti.status").
		Joins("LEFT JOIN drug_info di ON di.id = dti.drug_id AND di.deleted_at IS NULL").
		Joins("LEFT JOIN location_info li ON li.id = dti.location_id AND li.deleted_at IS NULL").
		Where("dti.status = 'LOSS_CANDIDATE'").
		Where("dti.deleted_at IS NULL").
		Order("dti.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		s.log.Error("查询盘亏候选药品失败", zap.Error(err))
		return nil, 0, ecode.ErrSystem
	}

	result := make([]*LossCandidateInfo, 0, len(rows))
	for _, row := range rows {
		result = append(result, &LossCandidateInfo{
			TraceCode:    row.TraceCode,
			DrugID:       row.DrugID,
			DrugName:     row.DrugName,
			BatchNumber:  row.BatchNumber,
			LocationCode: row.LocationCode,
			Status:       row.Status,
		})
	}
	return result, total, nil
}
