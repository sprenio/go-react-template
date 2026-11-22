package ratelimit

import (
	"context"
	"fmt"
	"time"

	"backend/internal/contexthelper"
)

type Service struct {
	limiter Limiter
	cfg     *Config
}

func NewService(limiter Limiter, cfg *Config) *Service {
	effectiveConfig := cfg
	if effectiveConfig == nil {
		defaultCfg := DefaultConfig()
		effectiveConfig = &defaultCfg
	}
	return &Service{
		limiter: limiter,
		cfg:     effectiveConfig,
	}
}

// CheckGlobal verifies global rate limits based on whether user is logged in or anonymous.
func (s *Service) CheckGlobal(ctx context.Context) (*Result, error) {
	if !s.cfg.Enabled {
		return &Result{Allowed: true}, nil
	}

	userId, isAuthorized := contexthelper.GetUserId(ctx)
	if isAuthorized && userId > 0 {
		keyPrefix := fmt.Sprintf("user:%d:global", userId)
		return s.limiter.Check(ctx, keyPrefix, s.cfg.AuthorizedGlobal)
	}

	clientIP := contexthelper.GetClientIp(ctx)
	if clientIP == "" {
		clientIP = "unknown"
	}
	keyPrefix := fmt.Sprintf("ip:%s:global", clientIP)
	return s.limiter.Check(ctx, keyPrefix, s.cfg.AnonymousGlobal)
}

// CheckEndpoint verifies endpoint-specific rate limits for the given endpoint identifier.
// Limits are tracked per User ID (if authenticated) or per Client IP (if unauthenticated).
func (s *Service) CheckEndpoint(ctx context.Context, endpointIdentifier string) (*Result, error) {
	if !s.cfg.Enabled {
		return &Result{Allowed: true}, nil
	}

	rules, hasRules := s.cfg.Endpoints[endpointIdentifier]
	if !hasRules || len(rules) == 0 {
		return &Result{Allowed: true}, nil
	}

	userId, isAuthorized := contexthelper.GetUserId(ctx)
	if isAuthorized && userId > 0 {
		keyPrefix := fmt.Sprintf("user:%d:endpoint:%s", userId, endpointIdentifier)
		return s.limiter.Check(ctx, keyPrefix, rules)
	}

	clientIP := contexthelper.GetClientIp(ctx)
	if clientIP == "" {
		clientIP = "unknown"
	}
	keyPrefix := fmt.Sprintf("ip:%s:endpoint:%s", clientIP, endpointIdentifier)
	return s.limiter.Check(ctx, keyPrefix, rules)
}

// CheckCustom allows checking any custom business logic rate limits.
// If customRules is provided, those rules are evaluated. If empty, rules configured under customIdentifier in Config are used.
func (s *Service) CheckCustom(ctx context.Context, customIdentifier string, subjectKey string, customRules []LimitRule) (*Result, error) {
	if !s.cfg.Enabled {
		return &Result{Allowed: true}, nil
	}

	rules := customRules
	if len(rules) == 0 {
		if configuredRules, exists := s.cfg.CustomLimits[customIdentifier]; exists {
			rules = configuredRules
		}
	}

	if len(rules) == 0 {
		return &Result{Allowed: true}, nil
	}

	keyPrefix := fmt.Sprintf("custom:%s:%s", customIdentifier, subjectKey)
	return s.limiter.Check(ctx, keyPrefix, rules)
}

// CheckCustomSingle is a convenience helper for single rate/window checks in business code.
func (s *Service) CheckCustomSingle(ctx context.Context, customIdentifier string, subjectKey string, maxRate int64, windowDuration time.Duration) (*Result, error) {
	rule := LimitRule{
		Rate:   maxRate,
		Window: windowDuration,
	}
	return s.CheckCustom(ctx, customIdentifier, subjectKey, []LimitRule{rule})
}
