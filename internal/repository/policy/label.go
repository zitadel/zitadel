package policy

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	LabelPolicyAddedEventType     = "policy.label.added"
	LabelPolicyChangedEventType   = "policy.label.changed"
	LabelPolicyActivatedEventType = "policy.label.activated"
	LabelPolicyRemovedEventType   = "policy.label.removed"
)

type LabelPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor        string `json:"primaryColor,omitempty"`
	SecondaryColor      string `json:"secondaryColor,omitempty"`
	WarnColor           string `json:"warnColor,omitempty"`
	PrimaryColorDark    string `json:"primaryColorDark,omitempty"`
	SecondaryColorDark  string `json:"secondaryColorDark,omitempty"`
	WarnColorDark       string `json:"warnColorDark,omitempty"`
	HideLoginNameSuffix bool   `json:"hideLoginNameSuffix,omitempty"`
	ErrorMsgPopup       bool   `json:"errorMsgPopup,omitempty"`
	DisableWatermark    bool   `json:"disableMsgPopup,omitempty"`
}

func (e *LabelPolicyAddedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyAddedEvent(
	base *eventstore.BaseEvent,
	primaryColor,
	secondaryColor,
	warnColor,
	primaryColorDark,
	secondaryColorDark,
	warnColorDark string,
	hideLoginNameSuffix,
	errorMsgPopup,
	disableWatermark bool,
) *LabelPolicyAddedEvent {

	return &LabelPolicyAddedEvent{
		BaseEvent:           *base,
		PrimaryColor:        primaryColor,
		SecondaryColor:      secondaryColor,
		WarnColor:           warnColor,
		PrimaryColorDark:    primaryColorDark,
		SecondaryColorDark:  secondaryColorDark,
		WarnColorDark:       warnColorDark,
		HideLoginNameSuffix: hideLoginNameSuffix,
		ErrorMsgPopup:       errorMsgPopup,
		DisableWatermark:    disableWatermark,
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
	WarnColor           *string `json:"warnColor,omitempty"`
	PrimaryColorDark    *string `json:"primaryColorDark,omitempty"`
	SecondaryColorDark  *string `json:"secondaryColorDark,omitempty"`
	WarnColorDark       *string `json:"warnColorDark,omitempty"`
	HideLoginNameSuffix *bool   `json:"hideLoginNameSuffix,omitempty"`
	ErrorMsgPopup       *bool   `json:"errorMsgPopup,omitempty"`
	DisableWatermark    *bool   `json:"disableMsgPopup,omitempty"`
}

func (e *LabelPolicyChangedEvent) Data() interface{} {
	return e
}

func (e *LabelPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func ChangeWarnColor(warnColor string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.WarnColor = &warnColor
	}
}

func ChangePrimaryColorDark(primaryColorDark string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.PrimaryColorDark = &primaryColorDark
	}
}

func ChangeSecondaryColorDark(secondaryColorDark string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.SecondaryColorDark = &secondaryColorDark
	}
}

func ChangeWarnColorDark(warnColorDark string) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.WarnColorDark = &warnColorDark
	}
}

func ChangeHideLoginNameSuffix(hideLoginNameSuffix bool) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.HideLoginNameSuffix = &hideLoginNameSuffix
	}
}

func ChangeErrorMsgPopup(errMsgPopup bool) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.ErrorMsgPopup = &errMsgPopup
	}
}

func ChangeDisableWatermark(disableWatermark bool) func(*LabelPolicyChangedEvent) {
	return func(e *LabelPolicyChangedEvent) {
		e.DisableWatermark = &disableWatermark
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

type LabelPolicyActivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LabelPolicyActivatedEvent) Data() interface{} {
	return nil
}

func (e *LabelPolicyActivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLabelPolicyActivatedEvent(base *eventstore.BaseEvent) *LabelPolicyActivatedEvent {
	return &LabelPolicyActivatedEvent{
		BaseEvent: *base,
	}
}

func LabelPolicyActivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &LabelPolicyActivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
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
