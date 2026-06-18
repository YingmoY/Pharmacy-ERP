// Package service 实现仪表盘模块的业务逻辑层。
package service

import (
	"context"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OverviewData 仪表盘概览数据。
type OverviewData struct {
	// TodaySalesAmount 今日销售总金额（COMPLETED/PARTIALLY_REFUNDED/REFUNDED 状态，按 paid_at 统计）。
	TodaySalesAmount float64 `json:"today_sales_amount"`
	// TodayOrderCount 今日订单数。
	TodayOrderCount int64 `json:"today_order_count"`
	// TodayInboundCount 今日完成入库单数（按 completed_at 统计）。
	TodayInboundCount int64 `json:"today_inbound_count"`
	// TotalInStock 当前在库追溯码数量。
	TotalInStock int64 `json:"total_in_stock"`
	// NearExpireCount 近效期药品数量（IN_STOCK 且 30 天内到期）。
	NearExpireCount int64 `json:"near_expire_count"`
	// LossCandidateCount 盘亏候选数量。
	LossCandidateCount int64 `json:"loss_candidate_count"`
	// PendingShelving 待上架数量（PENDING 状态）。
	PendingShelvingCount int64 `json:"pending_shelving_count"`
	// ActiveAlertCount 激活告警数量（audit_event status=0）。
	ActiveAlertCount int64 `json:"active_alert_count"`
}

// DailyStats 每日销售统计。
type DailyStats struct {
	// Date 日期，格式 YYYY-MM-DD。
	Date         string  `json:"date"`
	OrderCount   int64   `json:"order_count"`
	SalesAmount  float64 `json:"sales_amount"`
	RefundAmount float64 `json:"refund_amount"`
}

// DrugSalesStats 药品销售排行统计。
type DrugSalesStats struct {
	DrugID      int64   `json:"drug_id"`
	DrugName    string  `json:"drug_name"`
	DrugCode    string  `json:"drug_code"`
	SalesCount  int64   `json:"sales_count"`
	SalesAmount float64 `json:"sales_amount"`
}

// InboundStats 入库统计。
type InboundStats struct {
	TotalOrders     int64   `json:"total_orders"`
	CompletedOrders int64   `json:"completed_orders"`
	TotalAmount     float64 `json:"total_amount"`
	TodayCompleted  int64   `json:"today_completed"`
}

// InventoryStats 库存状态统计。
type InventoryStats struct {
	InStock       int64 `json:"in_stock"`
	Pending       int64 `json:"pending"`
	Sold          int64 `json:"sold"`
	Lost          int64 `json:"lost"`
	Misplaced     int64 `json:"misplaced"`
	LossCandidate int64 `json:"loss_candidate"`
	NearExpire7   int64 `json:"near_expire_7"`
	NearExpire15  int64 `json:"near_expire_15"`
	NearExpire30  int64 `json:"near_expire_30"`
}

// DashboardService 仪表盘服务接口。
type DashboardService interface {
	// GetOverview 获取仪表盘概览数据。
	GetOverview(ctx context.Context) (*OverviewData, error)
	// GetSalesTrend 获取最近 N 天销售趋势数据。
	GetSalesTrend(ctx context.Context, days int) ([]*DailyStats, error)
	// GetTopDrugs 获取最近 N 天销售排行前 limit 名的药品。
	GetTopDrugs(ctx context.Context, limit int, days int) ([]*DrugSalesStats, error)
	// GetInboundStats 获取最近 N 天入库统计。
	GetInboundStats(ctx context.Context, days int) (*InboundStats, error)
	// GetInventoryStats 获取当前库存状态统计。
	GetInventoryStats(ctx context.Context) (*InventoryStats, error)
}

type dashboardService struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewDashboardService 创建仪表盘服务实例。
func NewDashboardService(db *gorm.DB, log *zap.Logger) DashboardService {
	return &dashboardService{db: db, log: log}
}

// GetOverview 获取仪表盘概览数据。
func (s *dashboardService) GetOverview(ctx context.Context) (*OverviewData, error) {
	data := &OverviewData{}

	// 今日销售金额和订单数（按上海时间 paid_at 统计，已完成/部分退款/已退款订单）
	type salesRow struct {
		TotalAmount float64
		OrderCount  int64
	}
	var salesResult salesRow
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(actual_amount), 0) AS total_amount,
			COUNT(*) AS order_count
		FROM sales_order
		WHERE deleted_at IS NULL
		  AND status IN ('FINISHED', 'COMPLETED', 'PARTIALLY_REFUNDED', 'REFUNDED')
		  AND date_trunc('day', paid_at AT TIME ZONE 'Asia/Shanghai') = date_trunc('day', NOW() AT TIME ZONE 'Asia/Shanghai')
	`).Scan(&salesResult).Error; err != nil {
		s.log.Error("查询今日销售数据失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}
	data.TodaySalesAmount = salesResult.TotalAmount
	data.TodayOrderCount = salesResult.OrderCount

	// 今日完成入库单数
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM inbound_order
		WHERE deleted_at IS NULL
		  AND status = 'COMPLETED'
		  AND date_trunc('day', completed_at AT TIME ZONE 'Asia/Shanghai') = date_trunc('day', NOW() AT TIME ZONE 'Asia/Shanghai')
	`).Scan(&data.TodayInboundCount).Error; err != nil {
		s.log.Error("查询今日入库数量失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	// 在库数量
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM drug_trace_inventory WHERE deleted_at IS NULL AND status = 'IN_STOCK'
	`).Scan(&data.TotalInStock).Error; err != nil {
		s.log.Error("查询在库数量失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	// 近效期数量（30 天内）
	threshold := time.Now().AddDate(0, 0, 30)
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM drug_trace_inventory
		WHERE deleted_at IS NULL AND status = 'IN_STOCK' AND expire_date <= ?
	`, threshold).Scan(&data.NearExpireCount).Error; err != nil {
		s.log.Error("查询近效期数量失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	// 盘亏候选数量
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM drug_trace_inventory WHERE deleted_at IS NULL AND status = 'LOSS_CANDIDATE'
	`).Scan(&data.LossCandidateCount).Error; err != nil {
		s.log.Error("查询盘亏候选数量失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	// 待上架数量
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM drug_trace_inventory WHERE deleted_at IS NULL AND status = 'PENDING'
	`).Scan(&data.PendingShelvingCount).Error; err != nil {
		s.log.Error("查询待上架数量失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	// 激活告警数量
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM audit_event WHERE deleted_at IS NULL AND status = 0
	`).Scan(&data.ActiveAlertCount).Error; err != nil {
		s.log.Error("查询激活告警数量失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	return data, nil
}

// GetSalesTrend 获取最近 N 天销售趋势数据（按上海时区分组）。
func (s *dashboardService) GetSalesTrend(ctx context.Context, days int) ([]*DailyStats, error) {
	if days <= 0 {
		days = 7
	}

	type row struct {
		Date         string  `gorm:"column:date"`
		OrderCount   int64   `gorm:"column:order_count"`
		SalesAmount  float64 `gorm:"column:sales_amount"`
		RefundAmount float64 `gorm:"column:refund_amount"`
	}

	var rows []row
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			to_char(date_trunc('day', paid_at AT TIME ZONE 'Asia/Shanghai'), 'YYYY-MM-DD') AS date,
			COUNT(*) AS order_count,
			COALESCE(SUM(actual_amount), 0) AS sales_amount,
			COALESCE(SUM(refund_amount), 0) AS refund_amount
		FROM sales_order
		WHERE deleted_at IS NULL
		  AND status IN ('FINISHED', 'COMPLETED', 'PARTIALLY_REFUNDED', 'REFUNDED')
		  AND paid_at >= NOW() - INTERVAL '1 day' * ?
		GROUP BY date_trunc('day', paid_at AT TIME ZONE 'Asia/Shanghai')
		ORDER BY date ASC
	`, days).Scan(&rows).Error; err != nil {
		s.log.Error("查询销售趋势失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	result := make([]*DailyStats, 0, len(rows))
	for _, r := range rows {
		result = append(result, &DailyStats{
			Date:         r.Date,
			OrderCount:   r.OrderCount,
			SalesAmount:  r.SalesAmount,
			RefundAmount: r.RefundAmount,
		})
	}
	return result, nil
}

// GetTopDrugs 获取最近 N 天销售排行前 limit 名药品。
func (s *dashboardService) GetTopDrugs(ctx context.Context, limit int, days int) ([]*DrugSalesStats, error) {
	if limit <= 0 {
		limit = 10
	}
	if days <= 0 {
		days = 30
	}

	type row struct {
		DrugID      int64   `gorm:"column:drug_id"`
		DrugName    string  `gorm:"column:drug_name"`
		DrugCode    string  `gorm:"column:drug_code"`
		SalesCount  int64   `gorm:"column:sales_count"`
		SalesAmount float64 `gorm:"column:sales_amount"`
	}

	var rows []row
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			soi.drug_id,
			di.common_name AS drug_name,
			di.drug_code,
			SUM(soi.quantity) AS sales_count,
			COALESCE(SUM(soi.subtotal_amount), 0) AS sales_amount
		FROM sales_order_item soi
		JOIN sales_order so ON so.id = soi.order_id AND so.deleted_at IS NULL
		JOIN drug_info di ON di.id = soi.drug_id AND di.deleted_at IS NULL
		WHERE soi.deleted_at IS NULL
		  AND so.status IN ('FINISHED', 'COMPLETED', 'PARTIALLY_REFUNDED', 'REFUNDED')
		  AND so.paid_at >= NOW() - INTERVAL '1 day' * ?
		GROUP BY soi.drug_id, di.common_name, di.drug_code
		ORDER BY sales_count DESC
		LIMIT ?
	`, days, limit).Scan(&rows).Error; err != nil {
		s.log.Error("查询药品销售排行失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	result := make([]*DrugSalesStats, 0, len(rows))
	for _, r := range rows {
		result = append(result, &DrugSalesStats{
			DrugID:      r.DrugID,
			DrugName:    r.DrugName,
			DrugCode:    r.DrugCode,
			SalesCount:  r.SalesCount,
			SalesAmount: r.SalesAmount,
		})
	}
	return result, nil
}

// GetInboundStats 获取最近 N 天入库统计。
func (s *dashboardService) GetInboundStats(ctx context.Context, days int) (*InboundStats, error) {
	if days <= 0 {
		days = 30
	}

	stats := &InboundStats{}

	// 查询周期内所有入库单数量和金额
	type row struct {
		TotalOrders     int64   `gorm:"column:total_orders"`
		CompletedOrders int64   `gorm:"column:completed_orders"`
		TotalAmount     float64 `gorm:"column:total_amount"`
	}
	var r row
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) AS total_orders,
			SUM(CASE WHEN status = 'COMPLETED' THEN 1 ELSE 0 END) AS completed_orders,
			COALESCE(SUM(CASE WHEN status = 'COMPLETED' THEN total_amount ELSE 0 END), 0) AS total_amount
		FROM inbound_order
		WHERE deleted_at IS NULL
		  AND created_at >= NOW() - INTERVAL '1 day' * ?
	`, days).Scan(&r).Error; err != nil {
		s.log.Error("查询入库统计失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}
	stats.TotalOrders = r.TotalOrders
	stats.CompletedOrders = r.CompletedOrders
	stats.TotalAmount = r.TotalAmount

	// 今日完成入库单数
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM inbound_order
		WHERE deleted_at IS NULL
		  AND status = 'COMPLETED'
		  AND date_trunc('day', completed_at AT TIME ZONE 'Asia/Shanghai') = date_trunc('day', NOW() AT TIME ZONE 'Asia/Shanghai')
	`).Scan(&stats.TodayCompleted).Error; err != nil {
		s.log.Error("查询今日完成入库数量失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	return stats, nil
}

// GetInventoryStats 获取当前库存状态统计。
func (s *dashboardService) GetInventoryStats(ctx context.Context) (*InventoryStats, error) {
	stats := &InventoryStats{}
	now := time.Now()

	// 按状态分组统计各状态数量
	type statusRow struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:cnt"`
	}
	var statusRows []statusRow
	if err := s.db.WithContext(ctx).Raw(`
		SELECT status, COUNT(*) AS cnt
		FROM drug_trace_inventory
		WHERE deleted_at IS NULL
		GROUP BY status
	`).Scan(&statusRows).Error; err != nil {
		s.log.Error("查询库存状态统计失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	for _, r := range statusRows {
		switch r.Status {
		case "IN_STOCK":
			stats.InStock = r.Count
		case "PENDING":
			stats.Pending = r.Count
		case "SOLD":
			stats.Sold = r.Count
		case "LOST":
			stats.Lost = r.Count
		case "MISPLACED":
			stats.Misplaced = r.Count
		case "LOSS_CANDIDATE":
			stats.LossCandidate = r.Count
		}
	}

	// 近效期统计（仅 IN_STOCK）
	thresholds := []struct {
		days  int
		field *int64
	}{
		{7, &stats.NearExpire7},
		{15, &stats.NearExpire15},
		{30, &stats.NearExpire30},
	}

	for _, t := range thresholds {
		threshold := now.AddDate(0, 0, t.days)
		if err := s.db.WithContext(ctx).Raw(`
			SELECT COUNT(*) FROM drug_trace_inventory
			WHERE deleted_at IS NULL AND status = 'IN_STOCK' AND expire_date <= ?
		`, threshold).Scan(t.field).Error; err != nil {
			s.log.Error("查询近效期统计失败", zap.Int("days", t.days), zap.Error(err))
			return nil, ecode.ErrSystem
		}
	}

	return stats, nil
}
