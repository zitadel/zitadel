package command

import (
	"context"
	"sort"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/detection"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) CreateDetectionRule(ctx context.Context, rule detection.Rule) (*domain.ObjectDetails, error) {
	if rule.ID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-kB1uE", "Errors.IDMissing")
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	settingsWM, rulesWM, err := c.loadDetectionPolicyWriteModels(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	effectiveRules := detectionRuleMap(c.effectiveDetectionRules(settingsWM, rulesWM))
	if _, exists := effectiveRules[rule.ID]; exists {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-S3t0j", "Errors.Instance.Detection.Rule.AlreadyExists")
	}
	targetRules := cloneDetectionRuleMap(effectiveRules)
	targetRules[rule.ID] = rule
	cmds, err := c.detectionRuleCommands(ctx, instanceID, settingsWM, rulesWM, targetRules, false)
	if err != nil {
		return nil, err
	}
	return c.pushDetectionPolicyChanges(ctx, instanceID, settingsWM, rulesWM, cmds...)
}

func (c *Commands) ChangeDetectionRule(ctx context.Context, rule detection.Rule) (*domain.ObjectDetails, error) {
	if rule.ID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-y9jSn", "Errors.IDMissing")
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	settingsWM, rulesWM, err := c.loadDetectionPolicyWriteModels(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	effectiveRules := detectionRuleMap(c.effectiveDetectionRules(settingsWM, rulesWM))
	current, exists := effectiveRules[rule.ID]
	if !exists {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-H0mYt", "Errors.NotFound")
	}
	if detectionRulesEqual(current, rule) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-xJ0sK", "Errors.NoChangesFound")
	}
	targetRules := cloneDetectionRuleMap(effectiveRules)
	targetRules[rule.ID] = rule
	cmds, err := c.detectionRuleCommands(ctx, instanceID, settingsWM, rulesWM, targetRules, false)
	if err != nil {
		return nil, err
	}
	return c.pushDetectionPolicyChanges(ctx, instanceID, settingsWM, rulesWM, cmds...)
}

func (c *Commands) RemoveDetectionRule(ctx context.Context, ruleID string) (*domain.ObjectDetails, error) {
	if ruleID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-5gN4N", "Errors.IDMissing")
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	settingsWM, rulesWM, err := c.loadDetectionPolicyWriteModels(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	effectiveRules := detectionRuleMap(c.effectiveDetectionRules(settingsWM, rulesWM))
	if _, exists := effectiveRules[ruleID]; !exists {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-4QeYq", "Errors.NotFound")
	}
	targetRules := cloneDetectionRuleMap(effectiveRules)
	delete(targetRules, ruleID)
	cmds, err := c.detectionRuleCommands(ctx, instanceID, settingsWM, rulesWM, targetRules, true)
	if err != nil {
		return nil, err
	}
	return c.pushDetectionPolicyChanges(ctx, instanceID, settingsWM, rulesWM, cmds...)
}

func (c *Commands) loadDetectionPolicyWriteModels(ctx context.Context, instanceID string) (*InstanceDetectionSettingsWriteModel, *InstanceDetectionRulesWriteModel, error) {
	settingsWM, err := c.getInstanceDetectionSettingsWriteModel(ctx, instanceID)
	if err != nil {
		return nil, nil, err
	}
	rulesWM := NewInstanceDetectionRulesWriteModel(instanceID)
	if err := c.eventstore.FilterToQueryReducer(ctx, rulesWM); err != nil {
		return nil, nil, err
	}
	return settingsWM, rulesWM, nil
}

func (c *Commands) effectiveDetectionRules(settingsWM *InstanceDetectionSettingsWriteModel, rulesWM *InstanceDetectionRulesWriteModel) []detection.Rule {
	if settingsWM != nil && settingsWM.RulesManaged {
		if rulesWM == nil {
			return nil
		}
		return rulesWM.RulesSlice()
	}
	return c.DefaultDetectionRules()
}

func (c *Commands) detectionRuleCommands(ctx context.Context, instanceID string, settingsWM *InstanceDetectionSettingsWriteModel, rulesWM *InstanceDetectionRulesWriteModel, targetRules map[string]detection.Rule, deleting bool) ([]eventstore.Command, error) {
	cfg := c.defaultDetectionConfig
	if settingsWM != nil && settingsWM.SettingsSet {
		applyDetectionSettings(&cfg, settingsWM.Settings())
	}
	cfg.Rules = detectionRulesSlice(targetRules)
	if _, err := detection.NewPolicy(cfg); err != nil {
		return nil, err
	}
	aggregate := &instance.NewAggregate(instanceID).Aggregate
	cmds := make([]eventstore.Command, 0, len(targetRules)+1)
	if settingsWM == nil || !settingsWM.RulesManaged {
		managedEvent, err := instance.NewDetectionSettingsSetEvent(ctx, aggregate, []instance.DetectionSettingsChanges{instance.ChangeDetectionSettingsRulesManaged(true)})
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, managedEvent)
		for _, rule := range detectionRulesSlice(targetRules) {
			cmds = append(cmds, instance.NewDetectionRuleAddedEvent(ctx, aggregate, rule))
		}
		return cmds, nil
	}
	currentRules := detectionRuleMap(rulesWM.RulesSlice())
	currentIDs := make([]string, 0, len(currentRules))
	for id := range currentRules {
		currentIDs = append(currentIDs, id)
	}
	sort.Strings(currentIDs)
	for _, id := range currentIDs {
		current := currentRules[id]
		target, exists := targetRules[id]
		if !exists {
			cmds = append(cmds, instance.NewDetectionRuleRemovedEvent(ctx, aggregate, id))
			continue
		}
		changes := detectionRuleChanges(current, target)
		if len(changes) > 0 {
			cmd, err := instance.NewDetectionRuleChangedEvent(ctx, aggregate, id, changes)
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, cmd)
		}
	}
	targetIDs := make([]string, 0, len(targetRules))
	for id := range targetRules {
		if _, exists := currentRules[id]; exists {
			continue
		}
		targetIDs = append(targetIDs, id)
	}
	sort.Strings(targetIDs)
	for _, id := range targetIDs {
		cmds = append(cmds, instance.NewDetectionRuleAddedEvent(ctx, aggregate, targetRules[id]))
	}
	if len(cmds) == 0 && !deleting {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-F5m0p", "Errors.NoChangesFound")
	}
	return cmds, nil
}

func detectionRuleChanges(current, target detection.Rule) []instance.DetectionRuleChanges {
	changes := make([]instance.DetectionRuleChanges, 0, 11)
	if current.Description != target.Description {
		changes = append(changes, instance.ChangeDetectionRuleDescription(target.Description))
	}
	if current.Expr != target.Expr {
		changes = append(changes, instance.ChangeDetectionRuleExpr(target.Expr))
	}
	if current.Action != target.Action {
		changes = append(changes, instance.ChangeDetectionRuleEngine(target.Action))
	}
	if current.Priority != target.Priority {
		changes = append(changes, instance.ChangeDetectionRulePriority(target.Priority))
	}
	if current.StopOnMatch != target.StopOnMatch {
		changes = append(changes, instance.ChangeDetectionRuleStopOnMatch(target.StopOnMatch))
	}
	if current.FindingCfg.Name != target.FindingCfg.Name {
		changes = append(changes, instance.ChangeDetectionRuleFindingName(target.FindingCfg.Name))
	}
	if current.FindingCfg.Message != target.FindingCfg.Message {
		changes = append(changes, instance.ChangeDetectionRuleFindingMessage(target.FindingCfg.Message))
	}
	if current.FindingCfg.Block != target.FindingCfg.Block {
		changes = append(changes, instance.ChangeDetectionRuleFindingBlock(target.FindingCfg.Block))
	}
	if current.ContextTemplate != target.ContextTemplate {
		changes = append(changes, instance.ChangeDetectionRuleContextTemplate(target.ContextTemplate))
	}
	if current.RateLimitCfg.KeyTemplate != target.RateLimitCfg.KeyTemplate {
		changes = append(changes, instance.ChangeDetectionRuleRateLimitKey(target.RateLimitCfg.KeyTemplate))
	}
	if current.RateLimitCfg.Window != target.RateLimitCfg.Window {
		changes = append(changes, instance.ChangeDetectionRuleRateLimitWindow(target.RateLimitCfg.Window))
	}
	if current.RateLimitCfg.Max != target.RateLimitCfg.Max {
		changes = append(changes, instance.ChangeDetectionRuleRateLimitMax(target.RateLimitCfg.Max))
	}
	return changes
}

func (c *Commands) pushDetectionPolicyChanges(ctx context.Context, instanceID string, settingsWM *InstanceDetectionSettingsWriteModel, rulesWM *InstanceDetectionRulesWriteModel, cmds ...eventstore.Command) (*domain.ObjectDetails, error) {
	if len(cmds) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-c0Mmf", "Errors.NoChangesFound")
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	if err := AppendAndReduce(settingsWM, events...); err != nil {
		return nil, err
	}
	if err := AppendAndReduce(rulesWM, events...); err != nil {
		return nil, err
	}
	c.invalidateDetectionPolicyCache(instanceID)
	lastEvent := events[len(events)-1]
	return &domain.ObjectDetails{
		Sequence:      lastEvent.Sequence(),
		EventDate:     lastEvent.CreatedAt(),
		ResourceOwner: lastEvent.Aggregate().InstanceID,
	}, nil
}
