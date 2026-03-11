package instance

import (
"context"
"time"

"github.com/zitadel/zitadel/internal/eventstore"
"github.com/zitadel/zitadel/internal/zerrors"
)

const (
detectionSettingsPrefix       = "policy.detection."
DetectionSettingsSetEventType = instanceEventTypePrefix + detectionSettingsPrefix + "set"
)

type DetectionSettingsSetEvent struct {
*eventstore.BaseEvent `json:"-"`

Enabled               *bool          `json:"enabled,omitempty"`
FailOpen              *bool          `json:"fail_open,omitempty"`
FailureBurstThreshold *int           `json:"failure_burst_threshold,omitempty"`
HistoryWindow         *time.Duration `json:"history_window,omitempty"`
ContextChangeWindow   *time.Duration `json:"context_change_window,omitempty"`
MaxSignalsPerUser     *int           `json:"max_signals_per_user,omitempty"`
MaxSignalsPerSession  *int           `json:"max_signals_per_session,omitempty"`
RulesManaged          *bool          `json:"rules_managed,omitempty"`
}

func (e *DetectionSettingsSetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
e.BaseEvent = b
}

func NewDetectionSettingsSetEvent(ctx context.Context, aggregate *eventstore.Aggregate, changes []DetectionSettingsChanges) (*DetectionSettingsSetEvent, error) {
if len(changes) == 0 {
return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-xW1Ck", "Errors.NoChangesFound")
}
event := &DetectionSettingsSetEvent{
BaseEvent: eventstore.NewBaseEventForPush(ctx, aggregate, DetectionSettingsSetEventType),
}
for _, change := range changes {
change(event)
}
return event, nil
}

type DetectionSettingsChanges func(event *DetectionSettingsSetEvent)

func ChangeDetectionSettingsEnabled(enabled bool) DetectionSettingsChanges {
return func(e *DetectionSettingsSetEvent) {
e.Enabled = &enabled
}
}

func ChangeDetectionSettingsFailOpen(failOpen bool) DetectionSettingsChanges {
return func(e *DetectionSettingsSetEvent) {
e.FailOpen = &failOpen
}
}

func ChangeDetectionSettingsFailureBurstThreshold(threshold int) DetectionSettingsChanges {
return func(e *DetectionSettingsSetEvent) {
e.FailureBurstThreshold = &threshold
}
}

func ChangeDetectionSettingsHistoryWindow(window time.Duration) DetectionSettingsChanges {
return func(e *DetectionSettingsSetEvent) {
e.HistoryWindow = &window
}
}

func ChangeDetectionSettingsContextChangeWindow(window time.Duration) DetectionSettingsChanges {
return func(e *DetectionSettingsSetEvent) {
e.ContextChangeWindow = &window
}
}

func ChangeDetectionSettingsMaxSignalsPerUser(max int) DetectionSettingsChanges {
return func(e *DetectionSettingsSetEvent) {
e.MaxSignalsPerUser = &max
}
}

func ChangeDetectionSettingsMaxSignalsPerSession(max int) DetectionSettingsChanges {
return func(e *DetectionSettingsSetEvent) {
e.MaxSignalsPerSession = &max
}
}

func ChangeDetectionSettingsRulesManaged(managed bool) DetectionSettingsChanges {
return func(e *DetectionSettingsSetEvent) {
e.RulesManaged = &managed
}
}

func (e *DetectionSettingsSetEvent) Payload() any {
return e
}

func (e *DetectionSettingsSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
return nil
}
