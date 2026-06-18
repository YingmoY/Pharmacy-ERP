package task

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	inventoryModel "github.com/YingmoY/PharmacyERP/internal/inventory/model"
	inventoryRepo "github.com/YingmoY/PharmacyERP/internal/inventory/repository"
	inventoryService "github.com/YingmoY/PharmacyERP/internal/inventory/service"
	locationModel "github.com/YingmoY/PharmacyERP/internal/location/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/core"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"github.com/YingmoY/PharmacyERP/internal/pkg/middleware"
	"github.com/YingmoY/PharmacyERP/internal/pkg/mq"
	"github.com/YingmoY/PharmacyERP/internal/shelving/service"
	"github.com/YingmoY/PharmacyERP/internal/task/handler"
	taskModel "github.com/YingmoY/PharmacyERP/internal/task/model"
	"github.com/YingmoY/PharmacyERP/internal/task/repository"
	taskService "github.com/YingmoY/PharmacyERP/internal/task/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// inboundScanAdapter 将 InboundService 适配为 ScanTaskService 要求的 InboundScanProvider 接口。
type inboundScanAdapter struct {
	db      *gorm.DB
	inbound inventoryService.InboundService
}

func (a *inboundScanAdapter) ShelveByLocationCode(ctx context.Context, traceCode, locationCode string, operatorID int64) error {
	// 根据货位编码查找货位 ID。
	var loc locationModel.LocationInfo
	err := a.db.WithContext(ctx).
		Where("location_code = ? AND status = ?", strings.TrimSpace(locationCode), locationModel.LocationStatusEnabled).
		First(&loc).Error
	if err != nil {
		return ecode.ErrLocationNotFound
	}
	return a.inbound.PutawayTraceCodes(ctx, inventoryService.PutawayTraceCodesRequest{
		LocationID: loc.ID,
		TraceCodes: []string{traceCode},
	})
}

// Module 是 task 业务模块的 HTTP 适配层。
type Module struct {
	db               *gorm.DB
	inventoryTaskHdl *handler.InventoryTaskHandler
	scanTaskHdl      *handler.ScanTaskHandler
	inbound          inventoryService.InboundService
	mqClient         *mq.Client
	jwtSecret        string
}

func NewModule(db *gorm.DB, logger *zap.Logger, mqClient *mq.Client, jwtSecret string) *Module {
	taskRepo := repository.NewInventoryTask()
	orderRepo := inventoryRepo.NewInboundOrder()
	traceRepo := inventoryRepo.NewTraceInventory()
	inboundSvc := inventoryService.NewInboundService(db, orderRepo, traceRepo, logger)

	// 创建上架服务（用于 SHELVING 类型扫码任务路由）。
	shelvingSvc := service.NewShelvingService(db, logger)

	// 创建盘点任务服务。
	inventoryTaskSvc := taskService.NewInventoryTaskService(db, taskRepo, logger)

	// 创建扫码任务服务，注入入库和上架业务提供方。
	scanTaskSvc := taskService.NewScanTaskService(
		db,
		&inboundScanAdapter{db: db, inbound: inboundSvc},
		shelvingSvc,
		logger,
	)

	return &Module{
		db:               db,
		inventoryTaskHdl: handler.NewInventoryTaskHandler(inventoryTaskSvc),
		scanTaskHdl:      handler.NewScanTaskHandler(scanTaskSvc),
		inbound:          inboundSvc,
		mqClient:         mqClient,
		jwtSecret:        jwtSecret,
	}
}

func (m *Module) RegisterRoutes(group *gin.RouterGroup) {
	auth := middleware.JWTAuth(m.jwtSecret)

	// ─── 盘点任务 API ─────────────────────────────────────────────────
	invTaskGroup := group.Group("/inventory-tasks", auth)
	invTaskGroup.GET("", m.inventoryTaskHdl.ListInventoryTasks)
	invTaskGroup.POST("", m.inventoryTaskHdl.CreateInventoryTask)
	invTaskGroup.GET("/:id", m.inventoryTaskHdl.GetInventoryTask)
	invTaskGroup.POST("/:id/start", m.inventoryTaskHdl.StartInventoryTask)
	invTaskGroup.POST("/:id/scan", m.inventoryTaskHdl.SubmitInventoryScan)
	invTaskGroup.POST("/:id/complete", m.inventoryTaskHdl.CompleteInventoryTask)
	invTaskGroup.POST("/:id/cancel", m.inventoryTaskHdl.CancelInventoryTask)
	invTaskGroup.POST("/:id/assign", m.inventoryTaskHdl.AssignInventoryTask)
	invTaskGroup.GET("/:id/details", m.inventoryTaskHdl.GetInventoryTaskDetails)
	invTaskGroup.GET("/:id/summary", m.inventoryTaskHdl.GetInventoryTaskSummary)
	invTaskGroup.GET("/:id/loss-candidates", m.inventoryTaskHdl.GetLossCandidates)
	invTaskGroup.POST("/:id/loss-candidates/:trace_code/confirm", m.inventoryTaskHdl.ConfirmLoss)
	invTaskGroup.POST("/:id/loss-candidates/:trace_code/reject", m.inventoryTaskHdl.RejectLoss)
	invTaskGroup.GET("/:id/misplaced", m.inventoryTaskHdl.GetMisplaced)
	invTaskGroup.POST("/:id/misplaced/:trace_code/relocate", m.inventoryTaskHdl.RelocateMisplaced)

	// ─── 扫码任务 API ─────────────────────────────────────────────────
	scanTaskGroup := group.Group("/scan-tasks", auth)
	scanTaskGroup.GET("", m.scanTaskHdl.ListScanTasks)
	scanTaskGroup.POST("", m.scanTaskHdl.CreateScanTask)
	scanTaskGroup.GET("/:id", m.scanTaskHdl.GetScanTask)
	scanTaskGroup.POST("/:id/start", m.scanTaskHdl.StartScanTask)
	scanTaskGroup.POST("/:id/submit", m.scanTaskHdl.SubmitScan)
	scanTaskGroup.POST("/:id/complete", m.scanTaskHdl.CompleteScanTask)
	scanTaskGroup.POST("/:id/cancel", m.scanTaskHdl.CancelScanTask)
	scanTaskGroup.GET("/:id/details", m.scanTaskHdl.GetScanTaskDetails)

	// ─── 移动端兼容路由（向后兼容保留）──────────────────────────────────
	mobileGroup := group.Group("/mobile")
	mobileGroup.POST("/auth/login", m.login)
	mobileGroup.POST("/auth/register", m.register)
	mobileGroup.GET("/trace/detail", m.traceDetail)
	mobileGroup.GET("/trace/timeline", m.traceTimeline)
	mobileGroup.POST("/scan/inbound/submit", m.submitInboundScanCompat)
	mobileGroup.POST("/scan/inventory/submit", m.submitInventoryScanCompat)
}

// ──────────────────────────────────────────────────────────────────────────────
// 以下为向后兼容的实现，保留原有代码结构
// ──────────────────────────────────────────────────────────────────────────────

type sysUser struct {
	ID           int64  `gorm:"column:id"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password_hash"`
	RealName     string `gorm:"column:real_name"`
	Role         string `gorm:"column:role"`
	Status       int16  `gorm:"column:status"`
}

func (sysUser) TableName() string { return "sys_user" }

func (m *Module) register(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		DisplayName string `json:"displayName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}
	username := strings.TrimSpace(req.Username)
	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" {
		displayName = username
	}
	if username == "" || len(req.Password) < 6 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "username/password invalid")
		return
	}

	var existed int64
	if err := m.db.Model(&sysUser{}).Where("username = ?", username).Count(&existed).Error; err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "database error")
		return
	}
	if existed > 0 {
		core.Fail(c, 4009, "username already exists")
		return
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "password hash failed")
		return
	}

	user := sysUser{
		Username:     username,
		PasswordHash: string(hashBytes),
		RealName:     displayName,
		Role:         "WAREHOUSE",
		Status:       1,
	}
	if err := m.db.Create(&user).Error; err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "create user failed")
		return
	}

	core.Success(c, gin.H{
		"access_token":  "mock-jwt-token-for-" + user.Username,
		"refresh_token": "mock-refresh-token-for-" + user.Username,
		"expires_in":    86400,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
			"role":      user.Role,
			"status":    user.Status,
		},
		"userId":      fmt.Sprintf("%d", user.ID),
		"username":    user.Username,
		"displayName": user.RealName,
		"token":       "mock-jwt-token-for-" + user.Username,
	})
}

func (m *Module) login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	username := strings.TrimSpace(req.Username)
	var user sysUser
	if err := m.db.Where("username = ? AND status = 1", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			core.Fail(c, 4010, "invalid username or password")
			return
		}
		core.Fail(c, ecode.ErrSystem.Code, "database error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		core.Fail(c, 4010, "invalid username or password")
		return
	}

	core.Success(c, gin.H{
		"userId":      fmt.Sprintf("%d", user.ID),
		"username":    user.Username,
		"displayName": user.RealName,
		"token":       "mock-jwt-token-for-" + user.Username,
	})
}

type traceDetailRow struct {
	TraceCode       string `gorm:"column:trace_code"`
	DrugName        string `gorm:"column:drug_name"`
	Spec            string `gorm:"column:spec"`
	BatchNo         string `gorm:"column:batch_no"`
	ExpireDate      string `gorm:"column:expire_date"`
	CurrentLocation string `gorm:"column:current_location"`
	CurrentStatus   string `gorm:"column:current_status"`
}

type inboundScanCompatRequest struct {
	LocationCode string   `json:"locationCode" binding:"required"`
	TraceCodes   []string `json:"traceCodes" binding:"required,min=1"`
}

type inventoryScanCompatRequest struct {
	LocationCode       string   `json:"locationCode" binding:"required"`
	NormalCodes        []string `json:"normalCodes"`
	LossCandidateCodes []string `json:"lossCandidateCodes"`
	WrongShelfCodes    []string `json:"wrongShelfCodes"`
}

func (m *Module) submitInboundScanCompat(c *gin.Context) {
	var req inboundScanCompatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	var location locationModel.LocationInfo
	if err := m.db.Where("location_code = ? AND status = 1", strings.TrimSpace(req.LocationCode)).First(&location).Error; err != nil {
		core.Fail(c, ecode.ErrLocationNotFound.Code, ecode.ErrLocationNotFound.Msg)
		return
	}

	operatorID := m.currentOperatorID(c)
	scanTask, err := m.createScanTask(c, "INBOUND", location.ID, operatorID)
	if err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "create scan task failed")
		return
	}

	results := make([]gin.H, 0, len(req.TraceCodes))
	accepted := 0
	seen := make(map[string]struct{}, len(req.TraceCodes))
	for _, raw := range req.TraceCodes {
		code := strings.TrimSpace(raw)
		if code == "" {
			continue
		}

		if _, ok := seen[code]; ok {
			msg := "duplicate trace code in current scan task"
			_ = m.createScanTaskDetail(c, scanTask.ID, code, req.LocationCode, "DUPLICATE", msg)
			results = append(results, gin.H{
				"traceCode": code,
				"status":    "rejected",
				"message":   msg,
			})
			continue
		}
		seen[code] = struct{}{}

		oldStatus, _ := m.getTraceStatus(c, code)
		err := m.inbound.PutawayTraceCodes(c.Request.Context(), inventoryService.PutawayTraceCodesRequest{
			LocationID: location.ID,
			TraceCodes: []string{code},
		})
		if err != nil {
			bizErr := ecode.FromError(err)
			_ = m.createScanTaskDetail(c, scanTask.ID, code, req.LocationCode, "STATUS_ERROR", bizErr.Msg)
			results = append(results, gin.H{
				"traceCode": code,
				"status":    "rejected",
				"message":   bizErr.Msg,
			})
			continue
		}
		accepted++
		_ = m.createScanTaskDetail(c, scanTask.ID, code, req.LocationCode, "SUCCESS", "")
		_ = m.createDrugTraceLog(c, code, "SHELVING", oldStatus, inventoryModel.TraceInventoryStatusInStock, &location.ID, operatorID, scanTask.TaskNo, "mobile inbound scan accepted")
		results = append(results, gin.H{
			"traceCode": code,
			"status":    "accepted",
			"message":   "accepted",
		})
	}
	_ = m.completeScanTask(c, scanTask)

	m.publishTaskLog(c, "mobile_inbound_submit", scanTask.TaskNo, gin.H{
		"location_code": req.LocationCode,
		"total":         len(req.TraceCodes),
		"accepted":      accepted,
	})

	core.Success(c, gin.H{
		"total":    len(req.TraceCodes),
		"accepted": accepted,
		"rejected": len(req.TraceCodes) - accepted,
		"items":    results,
	})
}

func (m *Module) submitInventoryScanCompat(c *gin.Context) {
	var req inventoryScanCompatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		core.Fail(c, ecode.ErrParamInvalid.Code, ecode.ErrParamInvalid.Msg)
		return
	}

	var location locationModel.LocationInfo
	if err := m.db.Where("location_code = ? AND status = 1", strings.TrimSpace(req.LocationCode)).First(&location).Error; err != nil {
		core.Fail(c, ecode.ErrLocationNotFound.Code, ecode.ErrLocationNotFound.Msg)
		return
	}

	normalCodes := trimCodes(req.NormalCodes)
	lossCodes := trimCodes(req.LossCandidateCodes)
	wrongCodes := trimCodes(req.WrongShelfCodes)
	total := len(normalCodes) + len(lossCodes) + len(wrongCodes)
	if total == 0 {
		core.Fail(c, ecode.ErrParamInvalid.Code, "at least one trace code is required")
		return
	}

	operatorID := m.currentOperatorID(c)
	scanTask, err := m.createScanTask(c, "INVENTORY", location.ID, operatorID)
	if err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "create scan task failed")
		return
	}

	results := make([]gin.H, 0, total)
	accepted := 0
	seen := make(map[string]struct{}, total)
	accepted += m.applyInventoryScanCodes(c, scanTask, location, normalCodes, "normal", inventoryModel.TraceInventoryStatusInStock, operatorID, seen, &results)
	accepted += m.applyInventoryScanCodes(c, scanTask, location, lossCodes, "lossCandidate", inventoryModel.TraceInventoryStatusLossCandidate, operatorID, seen, &results)
	accepted += m.applyInventoryScanCodes(c, scanTask, location, wrongCodes, "wrongShelf", inventoryModel.TraceInventoryStatusMisplaced, operatorID, seen, &results)
	_ = m.completeScanTask(c, scanTask)

	m.publishTaskLog(c, "mobile_inventory_submit", scanTask.TaskNo, gin.H{
		"location_code": req.LocationCode,
		"total":         total,
		"accepted":      accepted,
	})

	core.Success(c, gin.H{
		"total":    total,
		"accepted": accepted,
		"rejected": total - accepted,
		"items":    results,
	})
}

func trimCodes(codes []string) []string {
	result := make([]string, 0, len(codes))
	for _, code := range codes {
		c := strings.TrimSpace(code)
		if c == "" {
			continue
		}
		result = append(result, c)
	}
	return result
}

func (m *Module) currentOperatorID(c *gin.Context) int64 {
	if operatorID, ok := middleware.GetCurrentUserID(c); ok && operatorID > 0 {
		return operatorID
	}
	if header := strings.TrimSpace(c.GetHeader("X-User-ID")); header != "" {
		if operatorID, err := strconv.ParseInt(header, 10, 64); err == nil && operatorID > 0 {
			return operatorID
		}
	}
	return 1
}

func (m *Module) createScanTask(c *gin.Context, taskType string, relatedID int64, operatorID int64) (*taskModel.ScanTask, error) {
	now := time.Now()
	task := &taskModel.ScanTask{
		TaskNo:     fmt.Sprintf("SCAN-%s-%d", taskType, now.UnixNano()),
		TaskType:   taskType,
		RelatedID:  relatedID,
		OperatorID: operatorID,
		Status:     "PROCESSING",
		StartTime:  &now,
	}
	if err := m.db.WithContext(c.Request.Context()).Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (m *Module) completeScanTask(c *gin.Context, task *taskModel.ScanTask) error {
	now := time.Now()
	return m.db.WithContext(c.Request.Context()).
		Model(task).
		Updates(map[string]interface{}{
			"status":   "COMPLETED",
			"end_time": now,
		}).Error
}

func (m *Module) createScanTaskDetail(c *gin.Context, taskID int64, traceCode, locationCode, result, msg string) error {
	var location *string
	if trimmed := strings.TrimSpace(locationCode); trimmed != "" {
		location = &trimmed
	}
	var errorMsg *string
	if strings.TrimSpace(msg) != "" {
		errorMsg = &msg
	}
	detail := taskModel.ScanTaskDetail{
		TaskID:       taskID,
		TraceCode:    traceCode,
		LocationCode: location,
		ScanResult:   result,
		ErrorMsg:     errorMsg,
		ScanTime:     time.Now(),
	}
	return m.db.WithContext(c.Request.Context()).Create(&detail).Error
}

func (m *Module) getTraceStatus(c *gin.Context, traceCode string) (string, error) {
	var trace inventoryModel.TraceInventory
	err := m.db.WithContext(c.Request.Context()).
		Select("status").
		Where("trace_code = ?", traceCode).
		First(&trace).Error
	if err != nil {
		return "", err
	}
	return trace.Status, nil
}

func (m *Module) createDrugTraceLog(c *gin.Context, traceCode, actionType, fromStatus, toStatus string, locationID *int64, operatorID int64, relatedNo, remark string) error {
	var from *string
	if strings.TrimSpace(fromStatus) != "" {
		from = &fromStatus
	}
	var to *string
	if strings.TrimSpace(toStatus) != "" {
		to = &toStatus
	}
	log := inventoryModel.DrugTraceLog{
		TraceCode:    traceCode,
		ActionType:   actionType,
		FromStatus:   from,
		ToStatus:     to,
		ToLocationID: locationID,
		OperatorID:   operatorID,
		RelatedNo:    &relatedNo,
		Remark:       &remark,
	}
	return m.db.WithContext(c.Request.Context()).Create(&log).Error
}

func (m *Module) applyInventoryScanCodes(c *gin.Context, task *taskModel.ScanTask, location locationModel.LocationInfo, codes []string, itemStatus string, targetStatus string, operatorID int64, seen map[string]struct{}, results *[]gin.H) int {
	accepted := 0
	for _, code := range codes {
		if _, ok := seen[code]; ok {
			msg := "duplicate trace code in current scan task"
			_ = m.createScanTaskDetail(c, task.ID, code, location.LocationCode, "DUPLICATE", msg)
			*results = append(*results, gin.H{"traceCode": code, "status": "rejected", "message": msg})
			continue
		}
		seen[code] = struct{}{}

		var trace inventoryModel.TraceInventory
		err := m.db.WithContext(c.Request.Context()).
			Where("trace_code = ?", code).
			First(&trace).Error
		if err != nil {
			bizErr := ecode.ErrTraceCodeNotFound
			_ = m.createScanTaskDetail(c, task.ID, code, location.LocationCode, "INVALID", bizErr.Msg)
			*results = append(*results, gin.H{"traceCode": code, "status": "rejected", "message": bizErr.Msg})
			continue
		}

		query := m.db.WithContext(c.Request.Context()).
			Model(&inventoryModel.TraceInventory{}).
			Where("trace_code = ?", code)
		if targetStatus == inventoryModel.TraceInventoryStatusInStock || targetStatus == inventoryModel.TraceInventoryStatusLossCandidate {
			query = query.Where("location_id = ?", location.ID)
		}
		update := map[string]interface{}{"status": targetStatus}
		if targetStatus == inventoryModel.TraceInventoryStatusMisplaced {
			update["location_id"] = location.ID
		}
		tx := query.Updates(update)
		if tx.Error != nil || tx.RowsAffected == 0 {
			msg := "trace code does not belong to current location"
			_ = m.createScanTaskDetail(c, task.ID, code, location.LocationCode, "STATUS_ERROR", msg)
			*results = append(*results, gin.H{"traceCode": code, "status": "rejected", "message": msg})
			continue
		}

		accepted++
		_ = m.createScanTaskDetail(c, task.ID, code, location.LocationCode, "SUCCESS", "")
		_ = m.createDrugTraceLog(c, code, inventoryActionType(targetStatus), trace.Status, targetStatus, &location.ID, operatorID, task.TaskNo, "mobile inventory scan accepted")
		*results = append(*results, gin.H{"traceCode": code, "status": itemStatus, "message": "accepted"})
	}
	return accepted
}

func inventoryActionType(targetStatus string) string {
	switch targetStatus {
	case inventoryModel.TraceInventoryStatusMisplaced:
		return "INVENTORY_MISPLACED"
	case inventoryModel.TraceInventoryStatusLossCandidate:
		return "INVENTORY_LOSS"
	default:
		return "INVENTORY_NORMAL"
	}
}

func (m *Module) traceDetail(c *gin.Context) {
	traceCode := strings.TrimSpace(c.Query("traceCode"))
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "traceCode is required")
		return
	}
	m.writeTraceDetail(c, traceCode)
}

func (m *Module) traceDetailByPath(c *gin.Context) {
	traceCode := strings.TrimSpace(c.Param("trace_code"))
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}
	m.writeTraceDetail(c, traceCode)
}

func (m *Module) writeTraceDetail(c *gin.Context, traceCode string) {
	var row traceDetailRow
	err := m.db.Raw(`
		SELECT
			t.trace_code AS trace_code,
			d.common_name AS drug_name,
			d.specification AS spec,
			t.batch_number AS batch_no,
			TO_CHAR(t.expire_date, 'YYYY-MM-DD') AS expire_date,
			COALESCE(l.location_code, '') AS current_location,
			t.status AS current_status
		FROM drug_trace_inventory t
		LEFT JOIN drug_info d ON d.id = t.drug_id
		LEFT JOIN location_info l ON l.id = t.location_id
		WHERE t.trace_code = ? AND t.deleted_at IS NULL
		LIMIT 1
	`, traceCode).Scan(&row).Error
	if err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "database error")
		return
	}
	if row.TraceCode == "" {
		core.Fail(c, ecode.ErrTraceCodeNotFound.Code, ecode.ErrTraceCodeNotFound.Msg)
		return
	}

	core.Success(c, gin.H{
		"trace_code":      row.TraceCode,
		"drug_name":       row.DrugName,
		"spec":            row.Spec,
		"batch_number":    row.BatchNo,
		"expire_date":     row.ExpireDate,
		"location_code":   row.CurrentLocation,
		"status":          row.CurrentStatus,
		"traceCode":       row.TraceCode,
		"drugName":        row.DrugName,
		"batchNo":         row.BatchNo,
		"expireDate":      row.ExpireDate,
		"currentLocation": row.CurrentLocation,
		"currentStatus":   row.CurrentStatus,
	})
}

func (m *Module) traceTimeline(c *gin.Context) {
	traceCode := strings.TrimSpace(c.Query("traceCode"))
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "traceCode is required")
		return
	}

	var logs []struct {
		Time         string `gorm:"column:time"`
		Event        string `gorm:"column:event"`
		OperatorName string `gorm:"column:operator_name"`
		Remark       string `gorm:"column:remark"`
	}

	err := m.db.Raw(`
		SELECT
			TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS') AS time,
			action_type AS event,
			CAST(operator_id AS TEXT) AS operator_name,
			COALESCE(remark, '') AS remark
		FROM drug_trace_log
		WHERE trace_code = ?
		ORDER BY created_at DESC
	`, traceCode).Scan(&logs).Error
	if err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "database error")
		return
	}

	if len(logs) == 0 {
		core.Success(c, []gin.H{})
		return
	}

	items := make([]gin.H, 0, len(logs))
	for _, it := range logs {
		items = append(items, gin.H{
			"time":         it.Time,
			"event":        it.Event,
			"operatorName": it.OperatorName,
			"remark":       it.Remark,
		})
	}
	core.Success(c, items)
}

func (m *Module) traceLogsByPath(c *gin.Context) {
	traceCode := strings.TrimSpace(c.Param("trace_code"))
	if traceCode == "" {
		core.Fail(c, ecode.ErrParamInvalid.Code, "trace_code is required")
		return
	}

	var logs []inventoryModel.DrugTraceLog
	if err := m.db.WithContext(c.Request.Context()).
		Where("trace_code = ?", traceCode).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		core.Fail(c, ecode.ErrSystem.Code, "database error")
		return
	}
	core.Success(c, gin.H{
		"total":     len(logs),
		"page":      1,
		"page_size": len(logs),
		"list":      logs,
	})
}

func (m *Module) publishTaskLog(c *gin.Context, action, businessID string, detail interface{}) {
	if m.mqClient == nil {
		return
	}
	operatorID, _ := middleware.GetCurrentUserID(c)
	_ = m.mqClient.PublishLogEvent(c.Request.Context(), mq.LogEvent{
		BusinessType: "task",
		BusinessID:   businessID,
		Action:       action,
		OperatorID:   operatorID,
		Detail:       detail,
	})
}
