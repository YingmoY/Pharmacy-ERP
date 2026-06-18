// Package service 实现报表模块的业务逻辑层。
package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/report/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SalesFilter 销售报表查询过滤条件。
type SalesFilter struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	CashierID     int64  `json:"cashier_id"`
	PaymentMethod string `json:"payment_method"`
	DrugID        int64  `json:"drug_id"`
	Status        string `json:"status"`
}

// InboundFilter 入库报表查询过滤条件。
type InboundFilter struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	SupplierID int64  `json:"supplier_id"`
	DrugID     int64  `json:"drug_id"`
	Status     string `json:"status"`
}

// InventoryFilter 库存报表查询过滤条件。
type InventoryFilter struct {
	DrugID      int64  `json:"drug_id"`
	LocationID  int64  `json:"location_id"`
	Status      string `json:"status"`
	BatchNumber string `json:"batch_number"`
}

// TraceLogFilter 追溯日志报表查询过滤条件。
type TraceLogFilter struct {
	TraceCode  string `json:"trace_code"`
	DrugID     int64  `json:"drug_id"`
	ActionType string `json:"action_type"`
	OperatorID int64  `json:"operator_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

// SalesReportItem 销售报表明细行。
type SalesReportItem struct {
	OrderNo      string  `json:"order_no"`
	Date         string  `json:"date"`
	CashierName  string  `json:"cashier_name"`
	DrugName     string  `json:"drug_name"`
	TraceCode    string  `json:"trace_code"`
	Price        float64 `json:"price"`
	RefundStatus string  `json:"refund_status"`
}

// SalesReport 销售报表。
type SalesReport struct {
	TotalOrders int64              `json:"total_orders"`
	TotalAmount float64            `json:"total_amount"`
	TotalRefund float64            `json:"total_refund"`
	NetAmount   float64            `json:"net_amount"`
	Items       []*SalesReportItem `json:"items"`
}

// InboundReportItem 入库报表明细行。
type InboundReportItem struct {
	OrderNo      string  `json:"order_no"`
	Date         string  `json:"date"`
	SupplierName string  `json:"supplier_name"`
	DrugName     string  `json:"drug_name"`
	Batch        string  `json:"batch"`
	Qty          int32   `json:"qty"`
	UnitPrice    float64 `json:"unit_price"`
	Amount       float64 `json:"amount"`
}

// InboundReport 入库报表。
type InboundReport struct {
	TotalOrders int64                `json:"total_orders"`
	TotalAmount float64              `json:"total_amount"`
	TotalQty    int64                `json:"total_qty"`
	Items       []*InboundReportItem `json:"items"`
}

// InventoryReportSummary 库存报表汇总。
type InventoryReportSummary struct {
	InStock       int64 `json:"in_stock"`
	Pending       int64 `json:"pending"`
	Sold          int64 `json:"sold"`
	Lost          int64 `json:"lost"`
	Misplaced     int64 `json:"misplaced"`
	LossCandidate int64 `json:"loss_candidate"`
}

// InventoryReportItem 库存报表明细行。
type InventoryReportItem struct {
	TraceCode   string    `json:"trace_code"`
	DrugName    string    `json:"drug_name"`
	BatchNumber string    `json:"batch_number"`
	ExpireDate  time.Time `json:"expire_date"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
}

// InventoryReport 库存报表。
type InventoryReport struct {
	Summary *InventoryReportSummary `json:"summary"`
	Items   []*InventoryReportItem  `json:"items"`
}

// TraceLogItem 追溯日志报表明细。
type TraceLogItem struct {
	ID           int64     `json:"id"`
	TraceCode    string    `json:"trace_code"`
	DrugName     string    `json:"drug_name"`
	ActionType   string    `json:"action_type"`
	FromStatus   string    `json:"from_status"`
	ToStatus     string    `json:"to_status"`
	OperatorID   int64     `json:"operator_id"`
	OperatorName string    `json:"operator_name"`
	RelatedNo    string    `json:"related_no"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
}

// ReportService 报表服务接口。
type ReportService interface {
	// GetSalesReport 查询销售报表数据。
	GetSalesReport(ctx context.Context, filter SalesFilter) (*SalesReport, error)
	// ExportSalesReport 创建销售报表导出任务。
	ExportSalesReport(ctx context.Context, filter SalesFilter, userID int64) (*model.ReportExportTask, error)
	// GetInboundReport 查询入库报表数据。
	GetInboundReport(ctx context.Context, filter InboundFilter) (*InboundReport, error)
	// ExportInboundReport 创建入库报表导出任务。
	ExportInboundReport(ctx context.Context, filter InboundFilter, userID int64) (*model.ReportExportTask, error)
	// GetInventoryReport 查询库存报表数据。
	GetInventoryReport(ctx context.Context, filter InventoryFilter) (*InventoryReport, error)
	// ExportInventoryReport 创建库存报表导出任务。
	ExportInventoryReport(ctx context.Context, filter InventoryFilter, userID int64) (*model.ReportExportTask, error)
	// GetTraceLogReport 查询追溯日志报表数据。
	GetTraceLogReport(ctx context.Context, filter TraceLogFilter) ([]*TraceLogItem, int64, error)
	// ExportTraceLogReport 创建追溯日志报表导出任务。
	ExportTraceLogReport(ctx context.Context, filter TraceLogFilter, userID int64) (*model.ReportExportTask, error)
	// GetExportTask 按 task_id 查询导出任务状态。
	GetExportTask(ctx context.Context, taskID string) (*model.ReportExportTask, error)
}

type reportService struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewReportService 创建报表服务实例。
func NewReportService(db *gorm.DB, log *zap.Logger) ReportService {
	return &reportService{db: db, log: log}
}

// createExportTask 创建导出任务并立即标记为成功（实际文件生成超出当前版本范围）。
func (s *reportService) createExportTask(ctx context.Context, reportType string, queryParams interface{}, userID int64) (*model.ReportExportTask, error) {
	paramsBytes, err := json.Marshal(queryParams)
	if err != nil {
		return nil, ecode.ErrSystem
	}

	now := time.Now()
	msg := "导出任务已创建（文件生成功能待实现）"
	task := &model.ReportExportTask{
		TaskID:       uuid.New().String(),
		ReportType:   reportType,
		ExportFormat: "xlsx",
		QueryParams:  datatypes.JSON(paramsBytes),
		Status:       model.ExportTaskStatusSuccess,
		Message:      &msg,
		RequestedBy:  userID,
		StartedAt:    &now,
		FinishedAt:   &now,
	}

	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		s.log.Error("创建导出任务失败", zap.String("report_type", reportType), zap.Error(err))
		return nil, ecode.ErrSystem
	}
	return task, nil
}

// GetSalesReport 查询销售报表数据。
func (s *reportService) GetSalesReport(ctx context.Context, filter SalesFilter) (*SalesReport, error) {
	query := s.db.WithContext(ctx).
		Table("sales_order so").
		Joins("JOIN sales_order_item soi ON soi.order_id = so.id AND soi.deleted_at IS NULL").
		Joins("JOIN drug_info di ON di.id = soi.drug_id AND di.deleted_at IS NULL").
		Joins("LEFT JOIN sys_user u ON u.id = so.cashier_id AND u.deleted_at IS NULL").
		Where("so.deleted_at IS NULL")

	if filter.StartDate != "" {
		query = query.Where("so.created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("so.created_at < ?", filter.EndDate)
	}
	if filter.CashierID > 0 {
		query = query.Where("so.cashier_id = ?", filter.CashierID)
	}
	if filter.PaymentMethod != "" {
		query = query.Where("so.payment_method = ?", filter.PaymentMethod)
	}
	if filter.DrugID > 0 {
		query = query.Where("soi.drug_id = ?", filter.DrugID)
	}
	if filter.Status != "" {
		query = query.Where("so.status = ?", filter.Status)
	}

	// 汇总统计（单独构建简单聚合查询，避免 join 放大）
	type summaryRow struct {
		TotalOrders int64   `gorm:"column:total_orders"`
		TotalAmount float64 `gorm:"column:total_amount"`
		TotalRefund float64 `gorm:"column:total_refund"`
	}
	var summary summaryRow
	sumQ := s.db.WithContext(ctx).
		Table("sales_order so").
		Where("so.deleted_at IS NULL")
	if filter.StartDate != "" {
		sumQ = sumQ.Where("so.created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		sumQ = sumQ.Where("so.created_at < ?", filter.EndDate)
	}
	if filter.CashierID > 0 {
		sumQ = sumQ.Where("so.cashier_id = ?", filter.CashierID)
	}
	if filter.PaymentMethod != "" {
		sumQ = sumQ.Where("so.payment_method = ?", filter.PaymentMethod)
	}
	if filter.Status != "" {
		sumQ = sumQ.Where("so.status = ?", filter.Status)
	}
	if err := sumQ.
		Select("COUNT(*) AS total_orders, COALESCE(SUM(actual_amount),0) AS total_amount, COALESCE(SUM(refund_amount),0) AS total_refund").
		Scan(&summary).Error; err != nil {
		s.log.Error("查询销售报表汇总失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	// 明细列表
	type itemRow struct {
		OrderNo      string  `gorm:"column:order_no"`
		Date         string  `gorm:"column:date"`
		CashierName  string  `gorm:"column:cashier_name"`
		DrugName     string  `gorm:"column:drug_name"`
		TraceCode    string  `gorm:"column:trace_code"`
		Price        float64 `gorm:"column:price"`
		RefundStatus string  `gorm:"column:refund_status"`
	}
	var itemRows []itemRow
	if err := query.
		Select(`so.order_no,
			to_char(so.created_at AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') AS date,
			COALESCE(u.real_name, u.username, '') AS cashier_name,
			di.common_name AS drug_name,
			soi.trace_code,
			soi.price,
			COALESCE(soi.refund_status, '') AS refund_status`).
		Order("so.created_at DESC").
		Limit(1000).
		Scan(&itemRows).Error; err != nil {
		s.log.Error("查询销售报表明细失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	items := make([]*SalesReportItem, 0, len(itemRows))
	for _, r := range itemRows {
		items = append(items, &SalesReportItem{
			OrderNo:      r.OrderNo,
			Date:         r.Date,
			CashierName:  r.CashierName,
			DrugName:     r.DrugName,
			TraceCode:    r.TraceCode,
			Price:        r.Price,
			RefundStatus: r.RefundStatus,
		})
	}

	return &SalesReport{
		TotalOrders: summary.TotalOrders,
		TotalAmount: summary.TotalAmount,
		TotalRefund: summary.TotalRefund,
		NetAmount:   summary.TotalAmount - summary.TotalRefund,
		Items:       items,
	}, nil
}


// ExportSalesReport 创建销售报表导出任务。
func (s *reportService) ExportSalesReport(ctx context.Context, filter SalesFilter, userID int64) (*model.ReportExportTask, error) {
	return s.createExportTask(ctx, "SALES", filter, userID)
}

// GetInboundReport 查询入库报表数据。
func (s *reportService) GetInboundReport(ctx context.Context, filter InboundFilter) (*InboundReport, error) {
	query := s.db.WithContext(ctx).
		Table("inbound_order io").
		Joins("JOIN inbound_order_detail iod ON iod.order_id = io.id AND iod.deleted_at IS NULL").
		Joins("JOIN drug_info di ON di.id = iod.drug_id AND di.deleted_at IS NULL").
		Joins("LEFT JOIN supplier s ON s.id = io.supplier_id AND s.deleted_at IS NULL").
		Where("io.deleted_at IS NULL")

	if filter.StartDate != "" {
		query = query.Where("io.created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("io.created_at < ?", filter.EndDate)
	}
	if filter.SupplierID > 0 {
		query = query.Where("io.supplier_id = ?", filter.SupplierID)
	}
	if filter.DrugID > 0 {
		query = query.Where("iod.drug_id = ?", filter.DrugID)
	}
	if filter.Status != "" {
		query = query.Where("io.status = ?", filter.Status)
	}

	type itemRow struct {
		OrderNo      string  `gorm:"column:order_no"`
		Date         string  `gorm:"column:date"`
		SupplierName string  `gorm:"column:supplier_name"`
		DrugName     string  `gorm:"column:drug_name"`
		Batch        string  `gorm:"column:batch"`
		Qty          int32   `gorm:"column:qty"`
		UnitPrice    float64 `gorm:"column:unit_price"`
		Amount       float64 `gorm:"column:amount"`
	}
	var itemRows []itemRow
	if err := query.
		Select(`io.order_no,
			to_char(io.created_at AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') AS date,
			COALESCE(s.name, '') AS supplier_name,
			di.common_name AS drug_name,
			iod.batch_number AS batch,
			iod.confirmed_qty AS qty,
			iod.unit_price,
			iod.amount`).
		Order("io.created_at DESC").
		Limit(1000).
		Scan(&itemRows).Error; err != nil {
		s.log.Error("查询入库报表明细失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	var totalOrders int64
	var totalAmount float64
	var totalQty int64
	items := make([]*InboundReportItem, 0, len(itemRows))
	orderSet := make(map[string]bool)
	for _, r := range itemRows {
		if !orderSet[r.OrderNo] {
			orderSet[r.OrderNo] = true
			totalOrders++
		}
		totalAmount += r.Amount
		totalQty += int64(r.Qty)
		items = append(items, &InboundReportItem{
			OrderNo:      r.OrderNo,
			Date:         r.Date,
			SupplierName: r.SupplierName,
			DrugName:     r.DrugName,
			Batch:        r.Batch,
			Qty:          r.Qty,
			UnitPrice:    r.UnitPrice,
			Amount:       r.Amount,
		})
	}

	return &InboundReport{
		TotalOrders: totalOrders,
		TotalAmount: totalAmount,
		TotalQty:    totalQty,
		Items:       items,
	}, nil
}

// ExportInboundReport 创建入库报表导出任务。
func (s *reportService) ExportInboundReport(ctx context.Context, filter InboundFilter, userID int64) (*model.ReportExportTask, error) {
	return s.createExportTask(ctx, "INBOUND", filter, userID)
}

// GetInventoryReport 查询库存报表数据。
func (s *reportService) GetInventoryReport(ctx context.Context, filter InventoryFilter) (*InventoryReport, error) {
	query := s.db.WithContext(ctx).
		Table("drug_trace_inventory dti").
		Joins("JOIN drug_info di ON di.id = dti.drug_id AND di.deleted_at IS NULL").
		Joins("LEFT JOIN location_info li ON li.id = dti.location_id AND li.deleted_at IS NULL").
		Where("dti.deleted_at IS NULL")

	if filter.DrugID > 0 {
		query = query.Where("dti.drug_id = ?", filter.DrugID)
	}
	if filter.LocationID > 0 {
		query = query.Where("dti.location_id = ?", filter.LocationID)
	}
	if filter.Status != "" {
		query = query.Where("dti.status = ?", filter.Status)
	}
	if filter.BatchNumber != "" {
		query = query.Where("dti.batch_number = ?", filter.BatchNumber)
	}

	// 汇总统计（按状态分组）
	type statusRow struct {
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:cnt"`
	}
	var statusRows []statusRow
	if err := s.db.WithContext(ctx).
		Table("drug_trace_inventory").
		Select("status, COUNT(*) AS cnt").
		Where("deleted_at IS NULL").
		Group("status").
		Scan(&statusRows).Error; err != nil {
		s.log.Error("查询库存汇总失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	summary := &InventoryReportSummary{}
	for _, r := range statusRows {
		switch r.Status {
		case "IN_STOCK":
			summary.InStock = r.Count
		case "PENDING":
			summary.Pending = r.Count
		case "SOLD":
			summary.Sold = r.Count
		case "LOST":
			summary.Lost = r.Count
		case "MISPLACED":
			summary.Misplaced = r.Count
		case "LOSS_CANDIDATE":
			summary.LossCandidate = r.Count
		}
	}

	// 明细列表
	type itemRow struct {
		TraceCode   string    `gorm:"column:trace_code"`
		DrugName    string    `gorm:"column:drug_name"`
		BatchNumber string    `gorm:"column:batch_number"`
		ExpireDate  time.Time `gorm:"column:expire_date"`
		Location    string    `gorm:"column:location"`
		Status      string    `gorm:"column:status"`
	}
	var itemRows []itemRow
	if err := query.
		Select(`dti.trace_code,
			di.common_name AS drug_name,
			dti.batch_number,
			dti.expire_date,
			COALESCE(li.location_code, '') AS location,
			dti.status`).
		Order("dti.expire_date ASC, dti.status ASC").
		Limit(2000).
		Scan(&itemRows).Error; err != nil {
		s.log.Error("查询库存报表明细失败", zap.Error(err))
		return nil, ecode.ErrSystem
	}

	items := make([]*InventoryReportItem, 0, len(itemRows))
	for _, r := range itemRows {
		items = append(items, &InventoryReportItem{
			TraceCode:   r.TraceCode,
			DrugName:    r.DrugName,
			BatchNumber: r.BatchNumber,
			ExpireDate:  r.ExpireDate,
			Location:    r.Location,
			Status:      r.Status,
		})
	}

	return &InventoryReport{
		Summary: summary,
		Items:   items,
	}, nil
}

// ExportInventoryReport 创建库存报表导出任务。
func (s *reportService) ExportInventoryReport(ctx context.Context, filter InventoryFilter, userID int64) (*model.ReportExportTask, error) {
	return s.createExportTask(ctx, "INVENTORY", filter, userID)
}

// GetTraceLogReport 查询追溯日志报表数据。
func (s *reportService) GetTraceLogReport(ctx context.Context, filter TraceLogFilter) ([]*TraceLogItem, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	query := s.db.WithContext(ctx).
		Table("drug_trace_log dtl").
		Joins("LEFT JOIN sys_user u ON u.id = dtl.operator_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN drug_trace_inventory dti ON dti.trace_code = dtl.trace_code AND dti.deleted_at IS NULL").
		Joins("LEFT JOIN drug_info di ON di.id = COALESCE(dtl.drug_id, dti.drug_id) AND di.deleted_at IS NULL").
		Where("dtl.deleted_at IS NULL")

	if filter.TraceCode != "" {
		query = query.Where("dtl.trace_code = ?", filter.TraceCode)
	}
	if filter.DrugID > 0 {
		query = query.Where("dtl.drug_id = ?", filter.DrugID)
	}
	if filter.ActionType != "" {
		query = query.Where("dtl.action_type = ?", filter.ActionType)
	}
	if filter.OperatorID > 0 {
		query = query.Where("dtl.operator_id = ?", filter.OperatorID)
	}
	if filter.StartDate != "" {
		query = query.Where("dtl.created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("dtl.created_at < ?", filter.EndDate)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		s.log.Error("统计追溯日志数量失败", zap.Error(err))
		return nil, 0, ecode.ErrSystem
	}

	if total == 0 {
		return []*TraceLogItem{}, 0, nil
	}

	type itemRow struct {
		ID           int64     `gorm:"column:id"`
		TraceCode    string    `gorm:"column:trace_code"`
		DrugName     string    `gorm:"column:drug_name"`
		ActionType   string    `gorm:"column:action_type"`
		FromStatus   string    `gorm:"column:from_status"`
		ToStatus     string    `gorm:"column:to_status"`
		OperatorID   int64     `gorm:"column:operator_id"`
		OperatorName string    `gorm:"column:operator_name"`
		RelatedNo    string    `gorm:"column:related_no"`
		Remark       string    `gorm:"column:remark"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}
	var rows []itemRow
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.
		Select(`dtl.id,
			dtl.trace_code,
			COALESCE(di.common_name, '') AS drug_name,
			dtl.action_type,
			COALESCE(dtl.from_status, '') AS from_status,
			COALESCE(dtl.to_status, '') AS to_status,
			dtl.operator_id,
			COALESCE(u.real_name, u.username, '') AS operator_name,
			COALESCE(dtl.related_no, '') AS related_no,
			COALESCE(dtl.remark, '') AS remark,
			dtl.created_at`).
		Order("dtl.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Scan(&rows).Error; err != nil {
		s.log.Error("查询追溯日志报表失败", zap.Error(err))
		return nil, 0, ecode.ErrSystem
	}

	items := make([]*TraceLogItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, &TraceLogItem{
			ID:           r.ID,
			TraceCode:    r.TraceCode,
			DrugName:     r.DrugName,
			ActionType:   r.ActionType,
			FromStatus:   r.FromStatus,
			ToStatus:     r.ToStatus,
			OperatorID:   r.OperatorID,
			OperatorName: r.OperatorName,
			RelatedNo:    r.RelatedNo,
			Remark:       r.Remark,
			CreatedAt:    r.CreatedAt,
		})
	}
	return items, total, nil
}

// ExportTraceLogReport 创建追溯日志报表导出任务。
func (s *reportService) ExportTraceLogReport(ctx context.Context, filter TraceLogFilter, userID int64) (*model.ReportExportTask, error) {
	return s.createExportTask(ctx, "TRACE_LOG", filter, userID)
}

// GetExportTask 按 task_id 查询导出任务状态。
func (s *reportService) GetExportTask(ctx context.Context, taskID string) (*model.ReportExportTask, error) {
	var task model.ReportExportTask
	if err := s.db.WithContext(ctx).Where("task_id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ecode.ErrNotFound
		}
		s.log.Error("查询导出任务失败", zap.String("task_id", taskID), zap.Error(err))
		return nil, ecode.ErrSystem
	}
	return &task, nil
}
