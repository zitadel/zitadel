package command

import (
	"bytes"
	"context"
	"encoding/json"

	"golang.org/x/exp/slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
)

type UserSchemaWriteModel struct {
	eventstore.WriteModel

	SchemaType             string
	Schema                 json.RawMessage
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
			wm.SchemaType = e.SchemaType
			wm.Schema = e.Schema
			wm.PossibleAuthenticators = e.PossibleAuthenticators
			wm.State = domain.UserSchemaStateActive
		case *schema.UpdatedEvent:
			if e.SchemaType != nil {
				wm.SchemaType = *e.SchemaType
			}
			if len(e.Schema) > 0 {
				wm.Schema = e.Schema
			}
			if len(e.PossibleAuthenticators) > 0 {
				wm.PossibleAuthenticators = e.PossibleAuthenticators
			}
		case *schema.DeactivatedEvent:
			wm.State = domain.UserSchemaStateInactive
		case *schema.ReactivatedEvent:
			wm.State = domain.UserSchemaStateActive
		case *schema.DeletedEvent:
			wm.State = domain.UserSchemaStateDeleted
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
			schema.UpdatedType,
			schema.DeactivatedType,
			schema.ReactivatedType,
			schema.DeletedType,
		).
		Builder()
}
func (wm *UserSchemaWriteModel) NewUpdatedEvent(
	ctx context.Context,
	agg *eventstore.Aggregate,
	schemaType *string,
	userSchema json.RawMessage,
	possibleAuthenticators []domain.AuthenticatorType,
) *schema.UpdatedEvent {
	changes := make([]schema.Changes, 0)
	if schemaType != nil && wm.SchemaType != *schemaType {
		changes = append(changes, schema.ChangeSchemaType(wm.SchemaType, *schemaType))
	}
	if !bytes.Equal(wm.Schema, userSchema) {
		changes = append(changes, schema.ChangeSchema(userSchema))
	}
	if len(possibleAuthenticators) > 0 && slices.Compare(wm.PossibleAuthenticators, possibleAuthenticators) != 0 {
		changes = append(changes, schema.ChangePossibleAuthenticators(possibleAuthenticators))
	}
	if len(changes) == 0 {
		return nil
	}
	return schema.NewUpdatedEvent(ctx, agg, changes)
}

func UserSchemaAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          schema.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       schema.AggregateVersion,
	}
}

func (wm *UserSchemaWriteModel) Exists() bool {
	return wm.State != domain.UserSchemaStateUnspecified && wm.State != domain.UserSchemaStateDeleted
}
