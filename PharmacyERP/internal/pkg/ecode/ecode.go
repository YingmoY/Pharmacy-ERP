package ecode

import "fmt"

// Error 是统一的业务错误类型
type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

func New(code int, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

// 通用系统错误码 1xxxx
var (
	ErrSystem       = New(10000, "system internal error")
	ErrParamInvalid = New(10001, "request params invalid")
	ErrUnauthorized = New(10002, "unauthorized")
	ErrPermission   = New(10003, "permission denied")
	ErrNotFound     = New(10004, "resource not found")
	ErrConflict     = New(10005, "resource conflict")
	ErrForbidden    = New(10006, "forbidden")
)

// 追溯/库存相关错误码 2xxxx
var (
	ErrTraceCodeNotFound = New(20001, "trace code not found")
	ErrStatusInvalid     = New(20002, "status transition invalid")
	ErrDuplicateScan     = New(20003, "trace code duplicated in task")
	ErrLocationNotFound  = New(20004, "location not found or disabled")
	ErrTraceStatusLocked = New(20005, "trace code status not allowed for operation")
	// ErrQtyExceeded 确认数量超过计划数量。
	ErrQtyExceeded = New(20006, "confirmed qty exceeds planned qty")
	// ErrNotFullyConfirmed 入库单存在未完全确认的明细。
	ErrNotFullyConfirmed = New(20007, "not all details fully confirmed")
	// ErrHasInStockTrace 入库单存在已上架追溯码，无法取消。
	ErrHasInStockTrace = New(20008, "order has in-stock trace codes, cannot cancel")
	// ErrInvoiceNotCompleted 发票未识别完成，不可转换。
	ErrInvoiceNotCompleted = New(20009, "invoice recognition not completed")
	// ErrInvoiceAlreadyConverted 发票已转换为入库单。
	ErrInvoiceAlreadyConverted = New(20010, "invoice already converted to inbound order")
	// ErrNoDetailForSubmit 入库单无明细，无法提交。
	ErrNoDetailForSubmit = New(20011, "inbound order has no details, cannot submit")
	// ErrSupplierInactive 供应商不存在或已停用。
	ErrSupplierInactive = New(10007, "supplier not found or inactive")
	// ErrDrugInactive 药品不存在或已停用。
	ErrDrugInactive = New(10008, "drug not found or inactive")
)

// 医保/处方相关错误码 3xxxx
var (
	ErrMedicareTrialFail = New(30001, "medicare trial settle failed")
	ErrPrescriptionAudit = New(30002, "prescription order requires pharmacist approval")
)

// 用户相关错误码 4xxxx
var (
	ErrUserNotFound       = New(40001, "user not found")
	ErrUserDisabled       = New(40002, "user is disabled")
	ErrPasswordWrong      = New(40003, "password is incorrect")
	ErrDuplicateUsername  = New(40004, "username already exists")
)

// 角色相关错误码 5xxxx
var (
	ErrRoleNotFound = New(50001, "role not found")
	ErrBuiltInRole  = New(50002, "built-in role cannot be modified or deleted")
	ErrRoleInUse    = New(50003, "role is assigned to users and cannot be deleted")
)

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return ErrSystem
}
