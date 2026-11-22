package ratelimit

import (
	"backend/config"
	"time"
)

// LimitRule defines rate and window duration (e.g. Rate: 10, Window: 1*time.Second).
type LimitRule struct {
	Rate   int64         `mapstructure:"rate" yaml:"rate"`
	Window time.Duration `mapstructure:"window" yaml:"window"`
}

// Config holds rate limit settings for anonymous, authorized, endpoints, and custom scopes.
type Config struct {
	Enabled          bool                   `mapstructure:"enabled" yaml:"enabled"`
	AnonymousGlobal  []LimitRule            `mapstructure:"anonymous_global" yaml:"anonymous_global"`
	AuthorizedGlobal []LimitRule            `mapstructure:"authorized_global" yaml:"authorized_global"`
	Endpoints        map[string][]LimitRule `mapstructure:"endpoints" yaml:"endpoints"`
	CustomLimits     map[string][]LimitRule `mapstructure:"custom_limits" yaml:"custom_limits"`
}

// FromAppConfig converts app config.RateLimitConfig to internal ratelimit.Config.
func FromAppConfig(cfg config.RateLimitConfig) Config {
	mapRules := func(ruleConfigs []config.RateLimitRuleConfig) []LimitRule {
		rules := make([]LimitRule, len(ruleConfigs))
		for index, ruleConfig := range ruleConfigs {
			rules[index] = LimitRule{
				Rate:   ruleConfig.Rate,
				Window: ruleConfig.Window,
			}
		}
		return rules
	}

	mapRuleMap := func(endpointMap map[string][]config.RateLimitRuleConfig) map[string][]LimitRule {
		result := make(map[string][]LimitRule, len(endpointMap))
		for key, ruleConfigs := range endpointMap {
			result[key] = mapRules(ruleConfigs)
		}
		return result
	}

	converted := Config{
		Enabled:          cfg.Enabled,
		AnonymousGlobal:  mapRules(cfg.AnonymousGlobal),
		AuthorizedGlobal: mapRules(cfg.AuthorizedGlobal),
		Endpoints:        mapRuleMap(cfg.Endpoints),
		CustomLimits:     mapRuleMap(cfg.CustomLimits),
	}

	// Fallback to defaults if global lists are not configured
	defaults := DefaultConfig()
	if len(converted.AnonymousGlobal) == 0 {
		converted.AnonymousGlobal = defaults.AnonymousGlobal
	}
	if len(converted.AuthorizedGlobal) == 0 {
		converted.AuthorizedGlobal = defaults.AuthorizedGlobal
	}
	if len(converted.Endpoints) == 0 {
		converted.Endpoints = defaults.Endpoints
	}

	return converted
}

// DefaultConfig returns default rate limit configuration matching system specifications.
func DefaultConfig() Config {
	return Config{
		Enabled: true,
		AnonymousGlobal: []LimitRule{
			{Rate: 10, Window: time.Second},
			{Rate: 100, Window: time.Minute},
			{Rate: 800, Window: 10 * time.Minute},
		},
		AuthorizedGlobal: []LimitRule{
			{Rate: 15, Window: time.Second},
			{Rate: 150, Window: time.Minute},
			{Rate: 1200, Window: 10 * time.Minute},
		},
		Endpoints: map[string][]LimitRule{
			"register": {
				{Rate: 1, Window: time.Second},
				{Rate: 3, Window: time.Minute},
				{Rate: 10, Window: 10 * time.Minute},
			},
			"login": {
				{Rate: 1, Window: time.Second},
				{Rate: 5, Window: time.Minute},
				{Rate: 20, Window: 10 * time.Minute},
			},
		},
		CustomLimits: map[string][]LimitRule{},
	}
}
