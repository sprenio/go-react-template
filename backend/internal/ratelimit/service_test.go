package ratelimit

import (
	"context"
	"testing"
	"time"
)

type mockLimiter struct {
	lastPrefix string
	lastRules  []LimitRule
	result     *Result
}

func (m *mockLimiter) Check(ctx context.Context, keyPrefix string, rules []LimitRule) (*Result, error) {
	m.lastPrefix = keyPrefix
	m.lastRules = rules
	return m.result, nil
}

func TestServiceCheckGlobalAnonymous(t *testing.T) {
	mock := &mockLimiter{
		result: &Result{Allowed: true, Remaining: 9},
	}
	cfg := DefaultConfig()
	service := NewService(mock, &cfg)

	ctx := context.Background()
	result, err := service.CheckGlobal(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !result.Allowed {
		t.Errorf("Expected allowed, got false")
	}
	if mock.lastPrefix != "ip:unknown:global" {
		t.Errorf("Expected prefix ip:unknown:global, got %s", mock.lastPrefix)
	}
	if len(mock.lastRules) != 3 {
		t.Errorf("Expected 3 rules, got %d", len(mock.lastRules))
	}
}

func TestServiceCheckCustom(t *testing.T) {
	mock := &mockLimiter{
		result: &Result{Allowed: false, Remaining: 0, RetryAfter: 5 * time.Second},
	}
	cfg := DefaultConfig()
	cfg.CustomLimits["sms_verification"] = []LimitRule{
		{Rate: 3, Window: time.Hour},
	}
	service := NewService(mock, &cfg)

	ctx := context.Background()
	result, err := service.CheckCustom(ctx, "sms_verification", "user_123", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Allowed {
		t.Errorf("Expected rejected, got allowed")
	}
	if mock.lastPrefix != "custom:sms_verification:user_123" {
		t.Errorf("Expected prefix custom:sms_verification:user_123, got %s", mock.lastPrefix)
	}
	if len(mock.lastRules) != 1 || mock.lastRules[0].Rate != 3 {
		t.Errorf("Expected custom configured rule to be passed")
	}
}
