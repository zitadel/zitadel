package command

import (
"context"
"time"

"github.com/zitadel/zitadel/internal/api/authz"
"github.com/zitadel/zitadel/internal/eventstore"
"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceDetectionSettingsWriteModel struct {
eventstore.WriteModel
Enabled               bool
FailOpen              bool
FailureBurstThreshold int
HistoryWindow         time.Duration
ContextChangeWindow   time.Duration
MaxSignalsPerUser     int
MaxSignalsPerSession  int
RulesManaged          bool
// SettingsSet is true only when at least one settings field (other than
// RulesManaged) has been explicitly configured. When false the runtime
// defaults are used as-is; only the rules list may be overridden.
SettingsSet bool
}

func NewInstanceDetectionSettingsWriteModel(instanceID string) *InstanceDetectionSettingsWriteModel {
return &InstanceDetectionSettingsWriteModel{
WriteModel: eventstore.WriteModel{
AggregateID:   instanceID,
ResourceOwner: instanceID,
InstanceID:    instanceID,
},
}
}

func NewInstanceDetectionSettingsWriteModelFromContext(ctx context.Context) *InstanceDetectionSettingsWriteModel {
return NewInstanceDetectionSettingsWriteModel(authz.GetInstance(ctx).InstanceID())
}

func (wm *InstanceDetectionSettingsWriteModel) GetWriteModel() *eventstore.WriteModel {
return &wm.WriteModel
}

func (wm *InstanceDetectionSettingsWriteModel) Settings() DetectionSettings {
return DetectionSettings{
Enabled:               wm.Enabled,
FailOpen:              wm.FailOpen,
FailureBurstThreshold: wm.FailureBurstThreshold,
HistoryWindow:         wm.HistoryWindow,
ContextChangeWindow:   wm.ContextChangeWindow,
MaxSignalsPerUser:     wm.MaxSignalsPerUser,
MaxSignalsPerSession:  wm.MaxSignalsPerSession,
}
}

func (wm *InstanceDetectionSettingsWriteModel) Reduce() error {
for _, event := range wm.Events {
e, ok := event.(*instance.DetectionSettingsSetEvent)
if !ok {
continue
}
if e.Enabled != nil {
wm.Enabled = *e.Enabled
wm.SettingsSet = true
}
if e.FailOpen != nil {
wm.FailOpen = *e.FailOpen
wm.SettingsSet = true
}
if e.FailureBurstThreshold != nil {
wm.FailureBurstThreshold = *e.FailureBurstThreshold
wm.SettingsSet = true
}
if e.HistoryWindow != nil {
wm.HistoryWindow = *e.HistoryWindow
wm.SettingsSet = true
}
if e.ContextChangeWindow != nil {
wm.ContextChangeWindow = *e.ContextChangeWindow
wm.SettingsSet = true
}
if e.MaxSignalsPerUser != nil {
wm.MaxSignalsPerUser = *e.MaxSignalsPerUser
wm.SettingsSet = true
}
if e.MaxSignalsPerSession != nil {
wm.MaxSignalsPerSession = *e.MaxSignalsPerSession
wm.SettingsSet = true
}
if e.RulesManaged != nil {
wm.RulesManaged = *e.RulesManaged
}
}
return wm.WriteModel.Reduce()
}

func (wm *InstanceDetectionSettingsWriteModel) Query() *eventstore.SearchQueryBuilder {
return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
ResourceOwner(wm.ResourceOwner).
AddQuery().
AggregateTypes(instance.AggregateType).
AggregateIDs(wm.AggregateID).
EventTypes(instance.DetectionSettingsSetEventType).
Builder()
}
