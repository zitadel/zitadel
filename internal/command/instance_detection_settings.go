package command

import (
"context"

"github.com/zitadel/zitadel/internal/api/authz"
"github.com/zitadel/zitadel/internal/detection"
"github.com/zitadel/zitadel/internal/domain"
"github.com/zitadel/zitadel/internal/repository/instance"
)

func (c *Commands) SetDetectionSettings(ctx context.Context, settings *DetectionSettings) (*domain.ObjectDetails, error) {
if err := settings.IsValid(); err != nil {
return nil, err
}
instanceID := authz.GetInstance(ctx).InstanceID()
settingsWM, rulesWM, err := c.loadDetectionPolicyWriteModels(ctx, instanceID)
if err != nil {
return nil, err
}
current := c.DefaultDetectionSettings()
if settingsWM.SettingsSet {
current = settingsWM.Settings()
}
cfg := c.defaultDetectionConfig
applyDetectionSettings(&cfg, *settings)
cfg.Rules = c.effectiveDetectionRules(settingsWM, rulesWM)
if _, err := detection.NewPolicy(cfg); err != nil {
return nil, err
}
aggregate := &instance.NewAggregate(instanceID).Aggregate
changes := make([]instance.DetectionSettingsChanges, 0, 7)
if current.Enabled != settings.Enabled {
changes = append(changes, instance.ChangeDetectionSettingsEnabled(settings.Enabled))
}
if current.FailOpen != settings.FailOpen {
changes = append(changes, instance.ChangeDetectionSettingsFailOpen(settings.FailOpen))
}
if current.FailureBurstThreshold != settings.FailureBurstThreshold {
changes = append(changes, instance.ChangeDetectionSettingsFailureBurstThreshold(settings.FailureBurstThreshold))
}
if current.HistoryWindow != settings.HistoryWindow {
changes = append(changes, instance.ChangeDetectionSettingsHistoryWindow(settings.HistoryWindow))
}
if current.ContextChangeWindow != settings.ContextChangeWindow {
changes = append(changes, instance.ChangeDetectionSettingsContextChangeWindow(settings.ContextChangeWindow))
}
if current.MaxSignalsPerUser != settings.MaxSignalsPerUser {
changes = append(changes, instance.ChangeDetectionSettingsMaxSignalsPerUser(settings.MaxSignalsPerUser))
}
if current.MaxSignalsPerSession != settings.MaxSignalsPerSession {
changes = append(changes, instance.ChangeDetectionSettingsMaxSignalsPerSession(settings.MaxSignalsPerSession))
}
cmd, err := instance.NewDetectionSettingsSetEvent(ctx, aggregate, changes)
if err != nil {
return nil, err
}
result, err := c.pushAppendAndReduceDetails(ctx, settingsWM, cmd)
if err == nil {
c.invalidateDetectionPolicyCache(instanceID)
}
return result, err
}

// GetEffectiveDetectionSettings returns the effective detection settings for the current instance.
// When settings have never been explicitly saved (only rule seeding has occurred), the system
// defaults are returned — consistent with how the runtime policy provider behaves.
func (c *Commands) GetEffectiveDetectionSettings(ctx context.Context) (DetectionSettings, error) {
instanceID := authz.GetInstance(ctx).InstanceID()
wm, err := c.getInstanceDetectionSettingsWriteModel(ctx, instanceID)
if err != nil {
return DetectionSettings{}, err
}
if wm.SettingsSet {
return wm.Settings(), nil
}
return c.DefaultDetectionSettings(), nil
}

func (c *Commands) getInstanceDetectionSettingsWriteModel(ctx context.Context, instanceID string) (*InstanceDetectionSettingsWriteModel, error) {
wm := NewInstanceDetectionSettingsWriteModel(instanceID)
if err := c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
return nil, err
}
return wm, nil
}
