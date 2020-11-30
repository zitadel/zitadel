package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicyAddedEventType   = iamEventTypePrefix + policy.LoginPolicyAddedEventType
	LoginPolicyChangedEventType = iamEventTypePrefix + policy.LoginPolicyChangedEventType

	LoginPolicyIDPProviderAddedEventType   = iamEventTypePrefix + policy.LoginPolicyIDPProviderAddedEventType
	LoginPolicyIDPProviderRemovedEventType = iamEventTypePrefix + policy.LoginPolicyIDPProviderRemovedEventType
)

type LoginPolicyReadModel struct{ policy.LoginPolicyReadModel }

func (rm *LoginPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *LoginPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.LoginPolicyChangedEvent)
		case *policy.LoginPolicyAddedEvent, *policy.LoginPolicyChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type LoginPolicyAddedEvent struct {
	policy.LoginPolicyAddedEvent
}

func LoginPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LoginPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyAddedEvent{LoginPolicyAddedEvent: *e.(*policy.LoginPolicyAddedEvent)}, nil
}

type LoginPolicyChangedEvent struct {
	policy.LoginPolicyChangedEvent
}

func LoginPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LoginPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyChangedEvent{LoginPolicyChangedEvent: *e.(*policy.LoginPolicyChangedEvent)}, nil
}

type LoginPolicyIDPProviderWriteModel struct {
	eventstore.WriteModel
	policy.IDPProviderWriteModel

	idpConfigID string
	iamID       string

	IsRemoved bool
}

func NewLoginPolicyIDPProviderWriteModel(iamID, idpConfigID string) *LoginPolicyIDPProviderWriteModel {
	return &LoginPolicyIDPProviderWriteModel{
		iamID:       iamID,
		idpConfigID: idpConfigID,
	}
}

func (wm *LoginPolicyIDPProviderWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyIDPProviderAddedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.IDPProviderWriteModel.AppendEvents(&e.IDPProviderAddedEvent)
		}
	}
}

func (wm *LoginPolicyIDPProviderWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *LoginPolicyIDPProviderAddedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.IsRemoved = false
		case *LoginPolicyIDPProviderRemovedEvent:
			if e.IDPConfigID != wm.idpConfigID {
				continue
			}
			wm.IsRemoved = true
		}
	}
	if err := wm.IDPProviderWriteModel.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *LoginPolicyIDPProviderWriteModel) Query() *eventstore.SearchQueryFactory {
	return eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}

type LoginPolicyIDPProviderAddedEvent struct {
	policy.IDPProviderAddedEvent
}

func NewLoginPolicyIDPProviderAddedEvent(
	ctx context.Context,
	idpConfigID string,
	idpProviderType provider.Type,
) *LoginPolicyIDPProviderAddedEvent {

	return &LoginPolicyIDPProviderAddedEvent{
		IDPProviderAddedEvent: *policy.NewIDPProviderAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyIDPProviderAddedEventType),
			idpConfigID,
			provider.TypeSystem),
	}
}

func IDPProviderAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.IDPProviderAddedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyIDPProviderAddedEvent{
		IDPProviderAddedEvent: *e.(*policy.IDPProviderAddedEvent),
	}, nil
}

type LoginPolicyIDPProviderRemovedEvent struct {
	policy.IDPProviderRemovedEvent
}

func NewLoginPolicyIDPProviderRemovedEvent(
	ctx context.Context,
	idpConfigID string,
) *LoginPolicyIDPProviderRemovedEvent {

	return &LoginPolicyIDPProviderRemovedEvent{
		IDPProviderRemovedEvent: *policy.NewIDPProviderRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyIDPProviderRemovedEventType),
			idpConfigID),
	}
}

func IDPProviderRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.IDPProviderRemovedEventEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyIDPProviderRemovedEvent{
		IDPProviderRemovedEvent: *e.(*policy.IDPProviderRemovedEvent),
	}, nil
}
