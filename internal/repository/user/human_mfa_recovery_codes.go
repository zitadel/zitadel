package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	recoveryCodeEventPrefix             = mfaEventPrefix + "recoverycode."
	HumanRecoveryCodesAddedType         = recoveryCodeEventPrefix + "added"
	HumanRecoveryCodesRemovedType       = recoveryCodeEventPrefix + "removed"
	HumanRecoveryCodeCheckSucceededType = recoveryCodeEventPrefix + "check.succeeded"
	HumanRecoveryCodeCheckFailedType    = recoveryCodeEventPrefix + "check.failed"
)

type HumanRecoveryCodesAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
	Codes []string `json:"codes,omitempty"`
}

func (e *HumanRecoveryCodesAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanRecoveryCodesAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanRecoveryCodesAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanRecoveryCodesAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	hashedCodes []string,
	optionalAuthRequest *AuthRequestInfo,
) *HumanRecoveryCodesAddedEvent {
	return &HumanRecoveryCodesAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanRecoveryCodesAddedType,
		),
		Codes:           hashedCodes,
		AuthRequestInfo: optionalAuthRequest,
	}
}

type HumanRecoveryCodesRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanRecoveryCodesRemovedEvent) Payload() interface{} {
	return nil
}

func (e *HumanRecoveryCodesRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanRecoveryCodesRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanRecoveryCodeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	authRequest *AuthRequestInfo,
) *HumanRecoveryCodesRemovedEvent {
	return &HumanRecoveryCodesRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanRecoveryCodesRemovedType,
		),
		AuthRequestInfo: authRequest,
	}
}

type HumanRecoveryCodeCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
	CodeIndex int `json:"codeIndex,omitempty"`
}

func (e *HumanRecoveryCodeCheckSucceededEvent) Payload() interface{} {
	return e
}

func (e *HumanRecoveryCodeCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanRecoveryCodeCheckSucceededEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanRecoveryCodeCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	codeIndex int,
	info *AuthRequestInfo,
) *HumanRecoveryCodeCheckSucceededEvent {
	return &HumanRecoveryCodeCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanRecoveryCodeCheckSucceededType,
		),
		AuthRequestInfo: info,
		CodeIndex:       codeIndex,
	}
}

type HumanRecoveryCodeCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanRecoveryCodeCheckFailedEvent) Payload() interface{} {
	return e
}

func (e *HumanRecoveryCodeCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanRecoveryCodeCheckFailedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanRecoveryCodeCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanRecoveryCodeCheckFailedEvent {
	return &HumanRecoveryCodeCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanRecoveryCodeCheckFailedType,
		),
		AuthRequestInfo: info,
	}
}
