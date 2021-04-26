package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	LabelPolicyAddedEventType   = "policy.label.added"
	LabelPolicyChangedEventType = "policy.label.changed"
	LabelPolicyRemovedEventType = "policy.label.removed"
)

type LabelPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor        string `json:"primaryColor,omitempty"`
	SecondaryColor      string `json:"secondaryColor,omitempty"`
	HideLoginNameSuffix bool   `json:"hideLoginNameSuffix,omitempty"`
	LogoDarkThemeID     string
}

func (e *LabelPolicyAddedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *LabelPolicyAddedEvent) Assets() []*eventstore.Asset {
	return nil
}

func NewLabelPolicyAddedEvent(
	base *eventstore.BaseEvent,
	primaryColor,
	secondaryColor string,
	hideLoginNameSuffix bool,
) *LabelPolicyAddedEvent {

	return &LabelPolicyAddedEvent{
		BaseEvent:           *base,
		PrimaryColor:        primaryColor,
		SecondaryColor:      secondaryColor,
		HideLoginNameSuffix: hideLoginNameSuffix,
	}
}

func LabelPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LabelPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type LabelPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor        *string `json:"primaryColor,omitempty"`
	SecondaryColor      *string `json:"secondaryColor,omitempty"`
	HideLoginNameSuffix *bool   `json:"hideLoginNameSuffix,omitempty"`
}

func (e *LabelPolicyChangedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *LabelPolicyChangedEvent) Assets() []*eventstore.Asset {
	return nil
}

func NewLabelPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []LabelPolicyChanges,
) (*LabelPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-Asfd3", "Errors.NoChangesFound")
	}
	changeEvent := &LabelPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type LabelPolicyChanges func(*LabelPolicyChangedEvent)

func ChangePrimaryColor(primaryColor string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.PrimaryColor = &primaryColor
	}
}

func ChangeSecondaryColor(secondaryColor string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.SecondaryColor = &secondaryColor
	}
}

func ChangeHideLoginNameSuffix(hideLoginNameSuffix bool) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.HideLoginNameSuffix = &hideLoginNameSuffix
	}
}

func LabelPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LabelPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-qhfFb", "unable to unmarshal label policy")
	}

	return e, nil
}

type LabelPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LabelPolicyRemovedEvent) Data() interface{} {
	return nil
}

func (e *LabelPolicyRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *LabelPolicyRemovedEvent) Assets() []*eventstore.Asset {
	return nil
}

func NewLabelPolicyRemovedEvent(base *eventstore.BaseEvent) *LabelPolicyRemovedEvent {
	return &LabelPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func LabelPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &LabelPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
