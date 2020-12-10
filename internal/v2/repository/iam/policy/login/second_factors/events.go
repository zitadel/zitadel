package second_factors

import (
	"context"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/second_factors"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

var (
	iamEventPrefix                         = eventstore.EventType("iam.")
	LoginPolicySecondfAddedEventType       = iamEventPrefix + second_factors.LoginPolicySecondFactorAddedEventType
	LoginPolicySecondFactoremovedEventType = iamEventPrefix + second_factors.LoginPolicySecondFactorRemovedEventType
)

type LoginPolicySecondFactorAddedEvent struct {
	second_factors.SecondFactorAddedEvent
}

func NewLoginPolicySecondFactorAddedEvent(
	ctx context.Context,
	mfaType second_factors.SecondFactorType,
) *LoginPolicySecondFactorAddedEvent {
	return &LoginPolicySecondFactorAddedEvent{
		SecondFactorAddedEvent: *second_factors.NewSecondFactorAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicySecondfAddedEventType),
			mfaType),
	}
}

//
//func SecondFactorAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
//	e, err := second_factors.SecondFactorAddedEventEventMapper(event)
//	if err != nil {
//		return nil, err
//	}
//
//	return &LoginPolicySecondFactorAddedEvent{
//		SecondFactorAddedEvent: *e.(*second_factors.SecondFactorAddedEvent),
//	}, nil
//}
//
//type LoginPolicySecondFactorRemovedEvent struct {
//	idpprovider.SecondFactorRemovedEvent
//}
//
//func NewLoginPolicySecondFactorRemovedEvent(
//	ctx context.Context,
//	idpConfigID string,
//) *LoginPolicySecondFactorRemovedEvent {
//
//	return &LoginPolicySecondFactorRemovedEvent{
//		SecondFactorRemovedEvent: *idpprovider.NewSecondFactorRemovedEvent(
//			eventstore.NewBaseEventForPush(ctx, login.LoginPolicySecondFactorRemovedEventType),
//			idpConfigID),
//	}
//}
//
//func SecondFactorRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
//	e, err := idpprovider.SecondFactorRemovedEventEventMapper(event)
//	if err != nil {
//		return nil, err
//	}
//
//	return &LoginPolicySecondFactorRemovedEvent{
//		SecondFactorRemovedEvent: *e.(*idpprovider.SecondFactorRemovedEvent),
//	}, nil
//}
