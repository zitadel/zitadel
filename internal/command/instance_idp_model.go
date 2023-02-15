package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

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
		case *instance.LDAPIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *instance.LDAPIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *instance.IDPRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.RemovedEvent)
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
			instance.LDAPIDPAddedEventType,
			instance.LDAPIDPChangedEventType,
			instance.IDPRemovedEventType,
		).
		Builder()
}
