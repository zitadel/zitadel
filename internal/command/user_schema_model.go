package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
)

type UserSchemaWriteModel struct {
	eventstore.WriteModel

	Type                   string
	Schema                 map[string]any
	PossibleAuthenticators []domain.AuthenticatorType
	State                  domain.UserSchemaState
}

func NewUserSchemaWriteModel(schemaID, resourceOwner string) *UserSchemaWriteModel {
	return &UserSchemaWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   schemaID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *UserSchemaWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *schema.CreatedEvent:
			wm.Type = e.SchemaType
			wm.Schema = e.Schema
			wm.PossibleAuthenticators = e.PossibleAuthenticators
			wm.State = domain.UserSchemaStateActive
			//case *user.PersonalAccessTokenRemovedEvent:
			//	wm.State = domain.PersonalAccessTokenStateRemoved
			//case *user.UserRemovedEvent:
			//	wm.State = domain.PersonalAccessTokenStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserSchemaWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(schema.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			schema.CreatedType,
		).
		Builder()
}

//
//func (wm *UserSchemaWriteModel) Exists() bool {
//	return wm.State != domain.PersonalAccessTokenStateUnspecified && wm.State != domain.PersonalAccessTokenStateRemoved
//}
