package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
)

type UserV3WriteModel struct {
	eventstore.WriteModel

	SchemaType     string
	SchemaRevision uint64

	Email           string
	IsEmailVerified bool
	Phone           string
	IsPhoneVerified bool

	Data json.RawMessage

	State domain.UserState
}

func NewUserV3WriteModel(resourceOwner, userID string) *UserV3WriteModel {
	return &UserV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *UserV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *schemauser.CreatedEvent:
			wm.SchemaType = e.SchemaType
			wm.SchemaRevision = 1
			wm.Data = e.Data
			wm.Email = e.Email
			wm.Phone = e.Phone

			wm.State = domain.UserStateActive
		case *schemauser.UpdatedEvent:
			if e.SchemaType != nil {
				wm.SchemaType = *e.SchemaType
			}
			if e.SchemaRevision != nil {
				wm.SchemaRevision = *e.SchemaRevision
			}
			if len(e.Data) > 0 {
				wm.Data = e.Data
			}
			if e.Email != nil {
				wm.Email = *e.Email
			}
			if e.Phone != nil {
				wm.Phone = *e.Phone
			}
		case *schemauser.DeletedEvent:
			wm.State = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(schemauser.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			schemauser.CreatedType,
			schemauser.UpdatedType,
			schemauser.DeletedType,
		)

	if wm.SchemaType != "" {
		query = query.EventData(map[string]interface{}{"schemaType": wm.SchemaType})
	}

	return query.Builder()
}

func (wm *UserV3WriteModel) NewUpdatedEvent(
	ctx context.Context,
	agg *eventstore.Aggregate,
	schemaType *string,
	schemaRevision *uint64,
	data json.RawMessage,
	email *string,
	phone *string,
) *schemauser.UpdatedEvent {
	changes := make([]schemauser.Changes, 0)
	if schemaType != nil && wm.SchemaType != *schemaType {
		changes = append(changes, schemauser.ChangeSchemaType(wm.SchemaType, *schemaType))
	}
	if schemaRevision != nil && wm.SchemaRevision != *schemaRevision {
		changes = append(changes, schemauser.ChangeSchemaRevision(wm.SchemaRevision, *schemaRevision))
	}
	if !bytes.Equal(wm.Data, data) {
		changes = append(changes, schemauser.ChangeData(data))
	}
	if email != nil && wm.Email != *email {
		changes = append(changes, schemauser.ChangeEmail(*email))
	}
	if phone != nil && wm.Phone != *phone {
		changes = append(changes, schemauser.ChangePhone(*phone))
	}
	if len(changes) == 0 {
		return nil
	}
	return schemauser.NewUpdatedEvent(ctx, agg, changes)
}

func UserV3AggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          schemauser.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       schemauser.AggregateVersion,
	}
}

func (wm *UserV3WriteModel) Exists() bool {
	return wm.State != domain.UserStateDeleted && wm.State != domain.UserStateUnspecified
}
