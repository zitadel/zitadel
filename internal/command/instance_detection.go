package command

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/zitadel/zitadel/internal/detection"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const detectionPolicyCacheTTL = 5 * time.Second

type DetectionSettings struct {
	Enabled               bool
	FailOpen              bool
	FailureBurstThreshold int
	HistoryWindow         time.Duration
	ContextChangeWindow   time.Duration
	MaxSignalsPerUser     int
	MaxSignalsPerSession  int
}

func (s *DetectionSettings) IsValid() error {
	if s == nil {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Q8tE1", "Errors.Risk.Invalid")
	}
	if s.FailureBurstThreshold <= 0 || s.HistoryWindow <= 0 || s.ContextChangeWindow <= 0 || s.MaxSignalsPerUser <= 0 || s.MaxSignalsPerSession <= 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-X0m4R", "Errors.Risk.Invalid")
	}
	return nil
}

func detectionSettingsFromConfig(cfg detection.Config) DetectionSettings {
	return DetectionSettings{
		Enabled:               cfg.Enabled,
		FailOpen:              cfg.FailOpen,
		FailureBurstThreshold: cfg.FailureBurstThreshold,
		HistoryWindow:         cfg.HistoryWindow,
		ContextChangeWindow:   cfg.ContextChangeWindow,
		MaxSignalsPerUser:     cfg.MaxSignalsPerUser,
		MaxSignalsPerSession:  cfg.MaxSignalsPerSession,
	}
}

func applyDetectionSettings(cfg *detection.Config, settings DetectionSettings) {
	cfg.Enabled = settings.Enabled
	cfg.FailOpen = settings.FailOpen
	cfg.FailureBurstThreshold = settings.FailureBurstThreshold
	cfg.HistoryWindow = settings.HistoryWindow
	cfg.ContextChangeWindow = settings.ContextChangeWindow
	cfg.MaxSignalsPerUser = settings.MaxSignalsPerUser
	cfg.MaxSignalsPerSession = settings.MaxSignalsPerSession
}

func cloneDetectionRules(rules []detection.Rule) []detection.Rule {
	if len(rules) == 0 {
		return nil
	}
	cloned := make([]detection.Rule, len(rules))
	copy(cloned, rules)
	return cloned
}

func detectionRuleMap(rules []detection.Rule) map[string]detection.Rule {
	if len(rules) == 0 {
		return map[string]detection.Rule{}
	}
	mapped := make(map[string]detection.Rule, len(rules))
	for _, rule := range rules {
		mapped[rule.ID] = rule
	}
	return mapped
}

func cloneDetectionRuleMap(rules map[string]detection.Rule) map[string]detection.Rule {
	cloned := make(map[string]detection.Rule, len(rules))
	for id, rule := range rules {
		cloned[id] = rule
	}
	return cloned
}

func detectionRulesSlice(rules map[string]detection.Rule) []detection.Rule {
	if len(rules) == 0 {
		return nil
	}
	ids := make([]string, 0, len(rules))
	for id := range rules {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	result := make([]detection.Rule, 0, len(ids))
	for _, id := range ids {
		result = append(result, rules[id])
	}
	return result
}

func detectionRulesEqual(a, b detection.Rule) bool {
	return a.ID == b.ID &&
		a.Description == b.Description &&
		a.Expr == b.Expr &&
		a.Action == b.Action &&
		a.Priority == b.Priority &&
		a.StopOnMatch == b.StopOnMatch &&
		a.FindingCfg == b.FindingCfg &&
		a.ContextTemplate == b.ContextTemplate &&
		a.RateLimitCfg == b.RateLimitCfg
}

func (c *Commands) DefaultDetectionSettings() DetectionSettings {
	if c == nil {
		return DetectionSettings{}
	}
	return detectionSettingsFromConfig(c.defaultDetectionConfig)
}

func (c *Commands) DefaultDetectionRules() []detection.Rule {
	if c == nil {
		return nil
	}
	return cloneDetectionRules(c.defaultDetectionConfig.Rules)
}

type cachedDetectionPolicy struct {
	policy    detection.Policy
	expiresAt time.Time
}

type instanceDetectionPolicyProvider struct {
	eventstore *eventstore.Eventstore
	defaults   detection.Config
	ttl        time.Duration

	mu    sync.RWMutex
	cache map[string]cachedDetectionPolicy
}

func newInstanceDetectionPolicyProvider(eventstore *eventstore.Eventstore, defaults detection.Config) *instanceDetectionPolicyProvider {
	return &instanceDetectionPolicyProvider{
		eventstore: eventstore,
		defaults:   defaults,
		ttl:        detectionPolicyCacheTTL,
		cache:      make(map[string]cachedDetectionPolicy),
	}
}

func (p *instanceDetectionPolicyProvider) Policy(ctx context.Context, instanceID string) (detection.Policy, error) {
	if p == nil || instanceID == "" {
		return detection.NewPolicy(p.defaults)
	}
	if policy, ok := p.cached(instanceID); ok {
		return policy, nil
	}
	policy, err := p.load(ctx, instanceID)
	if err != nil {
		return detection.Policy{}, err
	}
	p.mu.Lock()
	p.cache[instanceID] = cachedDetectionPolicy{policy: policy, expiresAt: time.Now().Add(p.ttl)}
	p.mu.Unlock()
	return policy, nil
}

func (p *instanceDetectionPolicyProvider) cached(instanceID string) (detection.Policy, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	cached, ok := p.cache[instanceID]
	if !ok || time.Now().After(cached.expiresAt) {
		return detection.Policy{}, false
	}
	return cached.policy, true
}

func (p *instanceDetectionPolicyProvider) Invalidate(instanceID string) {
	if p == nil || instanceID == "" {
		return
	}
	p.mu.Lock()
	delete(p.cache, instanceID)
	p.mu.Unlock()
}

func (p *instanceDetectionPolicyProvider) load(ctx context.Context, instanceID string) (detection.Policy, error) {
	settingsWM := NewInstanceDetectionSettingsWriteModel(instanceID)
	if err := p.eventstore.FilterToQueryReducer(ctx, settingsWM); err != nil {
		return detection.Policy{}, err
	}
	rulesWM := NewInstanceDetectionRulesWriteModel(instanceID)
	if err := p.eventstore.FilterToQueryReducer(ctx, rulesWM); err != nil {
		return detection.Policy{}, err
	}
	cfg := p.defaults
	if settingsWM.SettingsSet {
		applyDetectionSettings(&cfg, settingsWM.Settings())
	}
	if settingsWM.RulesManaged {
		cfg.Rules = rulesWM.RulesSlice()
	} else {
		cfg.Rules = cloneDetectionRules(p.defaults.Rules)
	}
	return detection.NewPolicy(cfg)
}

func (c *Commands) invalidateDetectionPolicyCache(instanceID string) {
	if c == nil || c.detectionPolicyProvider == nil {
		return
	}
	c.detectionPolicyProvider.Invalidate(instanceID)
}
