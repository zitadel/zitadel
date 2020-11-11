package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	LabelPolicyAddedEventType   = "policy.label.added"
	LabelPolicyChangedEventType = "policy.label.changed"
	LabelPolicyRemovedEventType = "policy.label.removed"
)

type LabelPolicyAggregate struct {
	eventstore.Aggregate

	PrimaryColor   string
	SecondaryColor string
}

type LabelPolicyReadModel struct {
	eventstore.ReadModel

	PrimaryColor   string
	SecondaryColor string
}

type LabelPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor   string `json:"primaryColor"`
	SecondaryColor string `json:"secondaryColor"`
}

func (e *LabelPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *LabelPolicyAddedEvent) Data() interface{} {
	return e
}

func NewLabelPolicyAddedEvent(
	ctx context.Context,
	primaryColor,
	secondaryColor string,
) *LabelPolicyAddedEvent {

	return &LabelPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			LabelPolicyAddedEventType,
		),
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
	}
}

type LabelPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PrimaryColor   string `json:"primaryColor,omitempty"`
	SecondaryColor string `json:"secondaryColor,omitempty"`
}

func (e *LabelPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *LabelPolicyChangedEvent) Data() interface{} {
	return e
}

func NewLabelPolicyChangedEvent(
	ctx context.Context,
	current,
	changed *LabelPolicyAggregate,
) *LabelPolicyChangedEvent {

	e := &LabelPolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			LabelPolicyChangedEventType,
		),
	}
	if current.PrimaryColor != changed.PrimaryColor {
		e.PrimaryColor = changed.PrimaryColor
	}
	if current.SecondaryColor != changed.SecondaryColor {
		e.SecondaryColor = changed.SecondaryColor
	}

	return e
}

type LabelPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LabelPolicyRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *LabelPolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewLabelPolicyRemovedEvent(ctx context.Context) *LabelPolicyRemovedEvent {
	return &LabelPolicyRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			LabelPolicyRemovedEventType,
		),
	}
}
