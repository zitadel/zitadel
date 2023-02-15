package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgLDAPIDPWriteModel struct {
	LDAPIDPWriteModel
}

func NewLDAPOrgIDPWriteModel(orgID, id string) *OrgLDAPIDPWriteModel {
	return &OrgLDAPIDPWriteModel{
		LDAPIDPWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgLDAPIDPWriteModel) Reduce() error {
	return wm.LDAPIDPWriteModel.Reduce()
}

func (wm *OrgLDAPIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LDAPIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *org.LDAPIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *org.IDPRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.LDAPIDPWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.LDAPIDPWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgLDAPIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.LDAPIDPAddedEventType,
			org.LDAPIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}

func (wm *OrgLDAPIDPWriteModel) NewChangedEvent(
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
) (*org.LDAPIDPChangedEvent, error) {

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
	changeEvent, err := org.NewLDAPIDPChangedEvent(ctx, aggregate, id, oldName, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

type OrgIDPRemoveWriteModel struct {
	IDPRemoveWriteModel
}

func NewOrgIDPRemoveWriteModel(orgID, id string) *OrgIDPRemoveWriteModel {
	return &OrgIDPRemoveWriteModel{
		IDPRemoveWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ID: id,
		},
	}
}

func (wm *OrgIDPRemoveWriteModel) Reduce() error {
	return wm.IDPRemoveWriteModel.Reduce()
}

func (wm *OrgIDPRemoveWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LDAPIDPAddedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPAddedEvent)
		case *org.LDAPIDPChangedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.LDAPIDPChangedEvent)
		case *org.IDPRemovedEvent:
			wm.IDPRemoveWriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.IDPRemoveWriteModel.AppendEvents(e)
		}
	}
}

func (wm *OrgIDPRemoveWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.LDAPIDPAddedEventType,
			org.LDAPIDPChangedEventType,
			org.IDPRemovedEventType,
		).
		Builder()
}
