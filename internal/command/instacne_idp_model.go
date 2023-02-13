package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceOAuthIDPWriteModel struct {
	OAuthIDPWriteModel
}

func NewOAuthInstanceIDPWriteModel(instanceID, id string) *InstanceOAuthIDPWriteModel {
	return &InstanceOAuthIDPWriteModel{
		OAuthIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceOAuthIDPWriteModel) Reduce() error {
	return wm.OAuthIDPWriteModel.Reduce()
}

func (wm *InstanceOAuthIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OAuthIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPAddedEvent)
		case *instance.OAuthIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.OAuthIDPWriteModel.AppendEvents(&e.OAuthIDPChangedEvent)
		}
	}
}

func (wm *InstanceOAuthIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.OAuthIDPAddedEventType,
			instance.OAuthIDPChangedEventType,
		).
		Builder()
}

func (wm *InstanceOAuthIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options idp.Options,
) (*instance.OAuthIDPChangedEvent, error) {

	changes, err := wm.OAuthIDPWriteModel.NewChanges(
		name,
		clientID,
		clientSecretString,
		secretCrypto,
		authorizationEndpoint,
		tokenEndpoint,
		userEndpoint,
		scopes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewOAuthIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type InstanceGoogleIDPWriteModel struct {
	GoogleIDPWriteModel
}

func NewGoogleInstanceIDPWriteModel(instanceID, id string) *InstanceGoogleIDPWriteModel {
	return &InstanceGoogleIDPWriteModel{
		GoogleIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceGoogleIDPWriteModel) Reduce() error {
	return wm.GoogleIDPWriteModel.Reduce()
}

func (wm *InstanceGoogleIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.GoogleIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *instance.GoogleIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		}
	}
}

func (wm *InstanceGoogleIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.GoogleIDPAddedEventType,
			instance.GoogleIDPChangedEventType,
		).
		Builder()
}

func (wm *InstanceGoogleIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GoogleIDPChangedEvent, error) {

	changes, err := wm.GoogleIDPWriteModel.NewChanges(clientID, clientSecretString, secretCrypto, scopes, options)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewGoogleIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}
