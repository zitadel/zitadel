package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceOIDCIDPWriteModel struct {
	OIDCIDPWriteModel
}

func NewOIDCInstanceIDPWriteModel(instanceID, id string) *InstanceOIDCIDPWriteModel {
	return &InstanceOIDCIDPWriteModel{
		OIDCIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceOIDCIDPWriteModel) Reduce() error {
	return wm.OIDCIDPWriteModel.Reduce()
}

func (wm *InstanceOIDCIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OIDCIDPAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *instance.OIDCIDPChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.RemovedEvent)

			// old events
		case *instance.IDPConfigAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPConfigChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *instance.IDPOIDCConfigAddedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *instance.IDPOIDCConfigChangedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *instance.IDPConfigRemovedEvent:
			wm.OIDCIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.OIDCIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceOIDCIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.OIDCIDPAddedEventType,
			instance.OIDCIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or(). // old events
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.IDPConfigAddedEventType,
			instance.IDPConfigChangedEventType,
			instance.IDPOIDCConfigAddedEventType,
			instance.IDPOIDCConfigChangedEventType,
			instance.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Builder()
}

func (wm *InstanceOIDCIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.OIDCIDPChangedEvent, error) {

	changes, err := wm.OIDCIDPWriteModel.NewChanges(
		name,
		issuer,
		clientID,
		clientSecretString,
		secretCrypto,
		scopes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewOIDCIDPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type InstanceJWTIDPWriteModel struct {
	JWTIDPWriteModel
}

func NewJWTInstanceIDPWriteModel(instanceID, id string) *InstanceJWTIDPWriteModel {
	return &InstanceJWTIDPWriteModel{
		JWTIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceJWTIDPWriteModel) Reduce() error {
	return wm.JWTIDPWriteModel.Reduce()
}

func (wm *InstanceJWTIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.JWTIDPAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *instance.JWTIDPChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.RemovedEvent)

			// old events
		case *instance.IDPConfigAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPConfigChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *instance.IDPJWTConfigAddedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *instance.IDPJWTConfigChangedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *instance.IDPConfigRemovedEvent:
			wm.JWTIDPWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.JWTIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceJWTIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.JWTIDPAddedEventType,
			instance.JWTIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or(). // old events
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.IDPConfigAddedEventType,
			instance.IDPConfigChangedEventType,
			instance.IDPJWTConfigAddedEventType,
			instance.IDPJWTConfigChangedEventType,
			instance.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Builder()
}

func (wm *InstanceJWTIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options idp.Options,
) (*instance.JWTIDPChangedEvent, error) {

	changes, err := wm.JWTIDPWriteModel.NewChanges(
		name,
		issuer,
		jwtEndpoint,
		keysEndpoint,
		headerName,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewJWTIDPChangedEvent(ctx, aggregate, id, changes)
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
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *instance.GoogleIDPChangedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.GoogleIDPWriteModel.AppendEvents(&e.RemovedEvent)
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
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Builder()
}

func (wm *InstanceGoogleIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) (*instance.GoogleIDPChangedEvent, error) {

	changes, err := wm.GoogleIDPWriteModel.NewChanges(name, clientID, clientSecretString, secretCrypto, scopes, options)
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

type InstanceLDAPIDPWriteModel struct {
	LDAPIDPWriteModel
}

func NewLDAPInstanceIDPWriteModel(instanceID, id string) *InstanceLDAPIDPWriteModel {
	return &InstanceLDAPIDPWriteModel{
		LDAPIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceLDAPIDPWriteModel) Reduce() error {
	return wm.LDAPIDPWriteModel.Reduce()
}

func (wm *InstanceLDAPIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.LDAPIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *instance.LDAPIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.LDAPIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceLDAPIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.LDAPIDPAddedEventType,
			instance.LDAPIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		Builder()
}

func (wm *InstanceLDAPIDPWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	oldName,
	name,
	host,
	port string,
	tls bool,
	baseDN,
	userObjectClass,
	userUniqueAttribute,
	admin string,
	password string,
	secretCrypto crypto.Crypto,
	attributes idp.LDAPAttributes,
	options idp.Options,
) (*instance.LDAPIDPChangedEvent, error) {

	changes, err := wm.LDAPIDPWriteModel.NewChanges(
		name,
		host,
		port,
		tls,
		baseDN,
		userObjectClass,
		userUniqueAttribute,
		admin,
		password,
		secretCrypto,
		attributes,
		options,
	)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return nil, nil
	}
	changeEvent, err := instance.NewLDAPIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type InstanceIDPRemoveWriteModel struct {
	IDPRemoveWriteModel
}

func NewInstanceIDPRemoveWriteModel(instanceID, id string) *InstanceIDPRemoveWriteModel {
	return &InstanceIDPRemoveWriteModel{
		IDPRemoveWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			ID: id,
		},
	}
}

func (wm *InstanceIDPRemoveWriteModel) Reduce() error {
	return wm.IDPRemoveWriteModel.Reduce()
}

func (wm *InstanceIDPRemoveWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.OIDCIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.OIDCIDPAddedEvent)
		case *instance.JWTIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.JWTIDPAddedEvent)
		case *instance.GoogleIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GoogleIDPAddedEvent)
		case *instance.GoogleIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.GoogleIDPChangedEvent)
		case *instance.LDAPIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *instance.LDAPIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.RemovedEvent)
		case *instance.IDPConfigAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPConfigRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.IDPRemoveWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceIDPRemoveWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.GoogleIDPAddedEventType,
			instance.GoogleIDPChangedEventType,
			instance.LDAPIDPAddedEventType,
			instance.LDAPIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		Builder()
}
