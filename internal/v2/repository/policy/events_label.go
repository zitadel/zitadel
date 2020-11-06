package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	LabelPolicyAddedEventType = "policy.label.added"
)

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
	service string,
	primaryColor,
	secondaryColor string,
) *LabelPolicyAddedEvent {

	return &LabelPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			LabelPolicyAddedEventType,
		),
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
	}
}
