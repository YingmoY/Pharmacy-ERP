package service

import (
	"context"
	"errors"
	"testing"

	"github.com/YingmoY/PharmacyERP/internal/inventory/model"
	"github.com/YingmoY/PharmacyERP/internal/pkg/ecode"
	"go.uber.org/zap"
)

// TestNormalizeAndValidateTraceCodes 验证追溯码清洗与去重规则。
//
// 覆盖点：
// 1) 正常输入可被原样返回；
// 2) 空字符串会触发参数错误；
// 3) 同一请求内重复追溯码会触发重复扫描错误。
func TestNormalizeAndValidateTraceCodes(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		wantLen   int
		wantError error
	}{
		{
			name:    "ok",
			input:   []string{"T001", "T002"},
			wantLen: 2,
		},
		{
			name:      "contains empty code",
			input:     []string{"T001", ""},
			wantError: ecode.ErrParamInvalid,
		},
		{
			name:      "duplicate code",
			input:     []string{"T001", "T001"},
			wantError: ecode.ErrDuplicateScan,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAndValidateTraceCodes(tt.input)
			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Fatalf("want error %v, got %v", tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("want len %d, got %d", tt.wantLen, len(got))
			}
		})
	}
}

// TestInboundService_StatusTransitions 验证状态流转规则函数。
//
// 这里不依赖数据库，聚焦纯规则判断，确保核心状态机约束稳定。
func TestInboundService_StatusTransitions(t *testing.T) {
	svc := &inboundService{}

	if !svc.IsValidOrderStatusTransition(model.InboundOrderStatusDraft, model.InboundOrderStatusPendingConfirm) {
		t.Fatalf("order transition DRAFT -> PENDING_CONFIRM should be valid")
	}
	if svc.IsValidOrderStatusTransition(model.InboundOrderStatusCompleted, model.InboundOrderStatusDraft) {
		t.Fatalf("order transition COMPLETED -> DRAFT should be invalid")
	}

	if !svc.IsValidTraceStatusTransition(model.TraceInventoryStatusPending, model.TraceInventoryStatusInStock) {
		t.Fatalf("trace transition PENDING -> IN_STOCK should be valid")
	}
	if svc.IsValidTraceStatusTransition(model.TraceInventoryStatusSold, model.TraceInventoryStatusInStock) {
		t.Fatalf("trace transition SOLD -> IN_STOCK should be invalid")
	}
}

// TestInboundService_RequestValidation 验证服务层入口参数校验。
//
// 说明：
// - 这些用例都在进入数据库事务之前返回，因此无需构造 DB 依赖；
// - 目标是兜住最常见的非法输入，避免无效请求进入事务逻辑。
func TestInboundService_RequestValidation(t *testing.T) {
	svc := &inboundService{
		logger: zap.NewNop(),
	}
	ctx := context.Background()

	// ConfirmInboundTraceCodes：必填参数缺失。
	if err := svc.ConfirmInboundTraceCodes(ctx, ConfirmInboundTraceCodesRequest{}); !errors.Is(err, ecode.ErrParamInvalid) {
		t.Fatalf("ConfirmInboundTraceCodes want ErrParamInvalid, got %v", err)
	}

	// ConfirmInboundTraceCodes：请求内重复追溯码。
	err := svc.ConfirmInboundTraceCodes(ctx, ConfirmInboundTraceCodesRequest{
		OrderID:    1,
		DetailID:   1,
		OperatorID: 1,
		TraceCodes: []string{"X001", "X001"},
	})
	if !errors.Is(err, ecode.ErrDuplicateScan) {
		t.Fatalf("ConfirmInboundTraceCodes duplicate code want ErrDuplicateScan, got %v", err)
	}

	// PutawayTraceCodes：必填参数缺失。
	if err := svc.PutawayTraceCodes(ctx, PutawayTraceCodesRequest{}); !errors.Is(err, ecode.ErrParamInvalid) {
		t.Fatalf("PutawayTraceCodes want ErrParamInvalid, got %v", err)
	}

	// PutawayTraceCodes：请求内重复追溯码。
	err = svc.PutawayTraceCodes(ctx, PutawayTraceCodesRequest{
		LocationID: 10,
		TraceCodes: []string{"Y001", "Y001"},
	})
	if !errors.Is(err, ecode.ErrDuplicateScan) {
		t.Fatalf("PutawayTraceCodes duplicate code want ErrDuplicateScan, got %v", err)
	}
}
