package instance

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/detection"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	detectionRulePrefix           = "policy.detection.rule."
	UniqueDetectionRule           = "instance_detection_rule"
	DetectionRuleAddedEventType   = instanceEventTypePrefix + detectionRulePrefix + "added"
	DetectionRuleChangedEventType = instanceEventTypePrefix + detectionRulePrefix + "changed"
	DetectionRuleRemovedEventType = instanceEventTypePrefix + detectionRulePrefix + "removed"
)

func NewAddDetectionRuleUniqueConstraint(ruleID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueDetectionRule,
		ruleID,
		"Errors.Instance.Detection.Rule.AlreadyExists",
	)
}

func NewRemoveDetectionRuleUniqueConstraint(ruleID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(UniqueDetectionRule, ruleID)
}

type DetectionRuleAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	RuleID          string               `json:"rule_id"`
	Description     string               `json:"description,omitempty"`
	Expr            string               `json:"expr"`
	Engine          detection.ActionType `json:"engine"`
	Priority        int                  `json:"priority,omitempty"`
	StopOnMatch     bool                 `json:"stop_on_match,omitempty"`
	FindingName     string               `json:"finding_name,omitempty"`
	FindingMessage  string               `json:"finding_message,omitempty"`
	FindingBlock    bool                 `json:"finding_block,omitempty"`
	ContextTemplate string               `json:"context_template,omitempty"`
	RateLimitKey    string               `json:"rate_limit_key,omitempty"`
	RateLimitWindow time.Duration        `json:"rate_limit_window,omitempty"`
	RateLimitMax    int                  `json:"rate_limit_max,omitempty"`
}

func (e *DetectionRuleAddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func NewDetectionRuleAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, rule detection.Rule) *DetectionRuleAddedEvent {
	return &DetectionRuleAddedEvent{
		BaseEvent:       eventstore.NewBaseEventForPush(ctx, aggregate, DetectionRuleAddedEventType),
		RuleID:          rule.ID,
		Description:     rule.Description,
		Expr:            rule.Expr,
		Engine:          rule.Action,
		Priority:        rule.Priority,
		StopOnMatch:     rule.StopOnMatch,
		FindingName:     rule.FindingCfg.Name,
		FindingMessage:  rule.FindingCfg.Message,
		FindingBlock:    rule.FindingCfg.Block,
		ContextTemplate: rule.ContextTemplate,
		RateLimitKey:    rule.RateLimitCfg.KeyTemplate,
		RateLimitWindow: rule.RateLimitCfg.Window,
		RateLimitMax:    rule.RateLimitCfg.Max,
	}
}

func (e *DetectionRuleAddedEvent) Payload() any {
	return e
}

func (e *DetectionRuleAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddDetectionRuleUniqueConstraint(e.RuleID)}
}

type DetectionRuleChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	RuleID          string                `json:"rule_id"`
	Description     *string               `json:"description,omitempty"`
	Expr            *string               `json:"expr,omitempty"`
	Engine          *detection.ActionType `json:"engine,omitempty"`
	Priority        *int                  `json:"priority,omitempty"`
	StopOnMatch     *bool                 `json:"stop_on_match,omitempty"`
	FindingName     *string               `json:"finding_name,omitempty"`
	FindingMessage  *string               `json:"finding_message,omitempty"`
	FindingBlock    *bool                 `json:"finding_block,omitempty"`
	ContextTemplate *string               `json:"context_template,omitempty"`
	RateLimitKey    *string               `json:"rate_limit_key,omitempty"`
	RateLimitWindow *time.Duration        `json:"rate_limit_window,omitempty"`
	RateLimitMax    *int                  `json:"rate_limit_max,omitempty"`
}

func (e *DetectionRuleChangedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

type DetectionRuleChanges func(event *DetectionRuleChangedEvent)

func NewDetectionRuleChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, ruleID string, changes []DetectionRuleChanges) (*DetectionRuleChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-6sd4k", "Errors.NoChangesFound")
	}
	event := &DetectionRuleChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(ctx, aggregate, DetectionRuleChangedEventType),
		RuleID:    ruleID,
	}
	for _, change := range changes {
		change(event)
	}
	return event, nil
}

func ChangeDetectionRuleDescription(description string) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.Description = &description
	}
}

func ChangeDetectionRuleExpr(expr string) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.Expr = &expr
	}
}

func ChangeDetectionRuleEngine(engine detection.ActionType) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.Engine = &engine
	}
}

func ChangeDetectionRuleFindingName(name string) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.FindingName = &name
	}
}

func ChangeDetectionRuleFindingMessage(message string) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.FindingMessage = &message
	}
}

func ChangeDetectionRuleFindingBlock(block bool) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.FindingBlock = &block
	}
}

func ChangeDetectionRuleContextTemplate(tmpl string) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.ContextTemplate = &tmpl
	}
}

func ChangeDetectionRuleRateLimitKey(key string) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.RateLimitKey = &key
	}
}

func ChangeDetectionRuleRateLimitWindow(window time.Duration) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.RateLimitWindow = &window
	}
}

func ChangeDetectionRuleRateLimitMax(max int) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.RateLimitMax = &max
	}
}

func ChangeDetectionRulePriority(priority int) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.Priority = &priority
	}
}

func ChangeDetectionRuleStopOnMatch(stopOnMatch bool) DetectionRuleChanges {
	return func(e *DetectionRuleChangedEvent) {
		e.StopOnMatch = &stopOnMatch
	}
}

func (e *DetectionRuleChangedEvent) Payload() any {
	return e
}

func (e *DetectionRuleChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type DetectionRuleRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	RuleID string `json:"rule_id"`
}

func (e *DetectionRuleRemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func NewDetectionRuleRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, ruleID string) *DetectionRuleRemovedEvent {
	return &DetectionRuleRemovedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(ctx, aggregate, DetectionRuleRemovedEventType),
		RuleID:    ruleID,
	}
}

func (e *DetectionRuleRemovedEvent) Payload() any {
	return e
}

func (e *DetectionRuleRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveDetectionRuleUniqueConstraint(e.RuleID)}
}
