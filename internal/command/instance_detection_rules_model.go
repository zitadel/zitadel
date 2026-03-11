package command

import (
	"sort"
	"time"

	"github.com/zitadel/zitadel/internal/detection"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type DetectionRuleState struct {
	Rule         detection.Rule
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}

type InstanceDetectionRulesWriteModel struct {
	eventstore.WriteModel
	Rules map[string]*DetectionRuleState
}

func NewInstanceDetectionRulesWriteModel(instanceID string) *InstanceDetectionRulesWriteModel {
	return &InstanceDetectionRulesWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
		Rules: make(map[string]*DetectionRuleState),
	}
}

func (wm *InstanceDetectionRulesWriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

func (wm *InstanceDetectionRulesWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.DetectionRuleAddedEvent:
			wm.Rules[e.RuleID] = &DetectionRuleState{
				Rule: detection.Rule{
					ID:          e.RuleID,
					Description: e.Description,
					Expr:        e.Expr,
					Action:      e.Engine,
					Priority:    e.Priority,
					StopOnMatch: e.StopOnMatch,
					FindingCfg: detection.RuleFinding{
						Name:    e.FindingName,
						Message: e.FindingMessage,
						Block:   e.FindingBlock,
					},
					ContextTemplate: e.ContextTemplate,
					RateLimitCfg: detection.RuleRateLimit{
						KeyTemplate: e.RateLimitKey,
						Window:      e.RateLimitWindow,
						Max:         e.RateLimitMax,
					},
				},
				CreationDate: e.CreationDate(),
				ChangeDate:   e.CreationDate(),
				Sequence:     e.Sequence(),
			}
		case *instance.DetectionRuleChangedEvent:
			rule := wm.Rules[e.RuleID]
			if rule == nil {
				continue
			}
			if e.Description != nil {
				rule.Rule.Description = *e.Description
			}
			if e.Expr != nil {
				rule.Rule.Expr = *e.Expr
			}
			if e.Engine != nil {
				rule.Rule.Action = *e.Engine
			}
			if e.FindingName != nil {
				rule.Rule.FindingCfg.Name = *e.FindingName
			}
			if e.FindingMessage != nil {
				rule.Rule.FindingCfg.Message = *e.FindingMessage
			}
			if e.FindingBlock != nil {
				rule.Rule.FindingCfg.Block = *e.FindingBlock
			}
			if e.ContextTemplate != nil {
				rule.Rule.ContextTemplate = *e.ContextTemplate
			}
			if e.RateLimitKey != nil {
				rule.Rule.RateLimitCfg.KeyTemplate = *e.RateLimitKey
			}
			if e.RateLimitWindow != nil {
				rule.Rule.RateLimitCfg.Window = *e.RateLimitWindow
			}
			if e.RateLimitMax != nil {
				rule.Rule.RateLimitCfg.Max = *e.RateLimitMax
			}
			if e.Priority != nil {
				rule.Rule.Priority = *e.Priority
			}
			if e.StopOnMatch != nil {
				rule.Rule.StopOnMatch = *e.StopOnMatch
			}
			rule.ChangeDate = e.CreationDate()
			rule.Sequence = e.Sequence()
		case *instance.DetectionRuleRemovedEvent:
			delete(wm.Rules, e.RuleID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *InstanceDetectionRulesWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(instance.DetectionRuleAddedEventType, instance.DetectionRuleChangedEventType, instance.DetectionRuleRemovedEventType).
		Builder()
}

func (wm *InstanceDetectionRulesWriteModel) RulesSlice() []detection.Rule {
	if len(wm.Rules) == 0 {
		return nil
	}
	ids := make([]string, 0, len(wm.Rules))
	for id := range wm.Rules {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	rules := make([]detection.Rule, 0, len(ids))
	for _, id := range ids {
		rules = append(rules, wm.Rules[id].Rule)
	}
	return rules
}
