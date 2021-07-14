package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

type UserMetaDataWriteModel struct {
	MetaDataWriteModel
}

func NewUserMetaDataWriteModel(userID, resourceOwner, key string) *UserMetaDataWriteModel {
	return &UserMetaDataWriteModel{
		MetaDataWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   userID,
				ResourceOwner: resourceOwner,
			},
			Key: key,
		},
	}
}

func (wm *UserMetaDataWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MetaDataSetEvent:
			wm.MetaDataWriteModel.AppendEvents(&e.SetEvent)
		case *user.MetaDataRemovedEvent:
			wm.MetaDataWriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *UserMetaDataWriteModel) Reduce() error {
	return wm.MetaDataWriteModel.Reduce()
}

func (wm *UserMetaDataWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetaDataWriteModel.AggregateID).
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.MetaDataSetType,
			user.MetaDataRemovedType).
		Builder()
}

type UserMetaDataListWriteModel struct {
	MetaDataWriteModel
}

func NewUserMetaDataListWriteModel(userID, resourceOwner string) *UserMetaDataListWriteModel {
	return &UserMetaDataListWriteModel{
		MetaDataWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   userID,
				ResourceOwner: resourceOwner,
			},
		},
	}
}

func (wm *UserMetaDataListWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MetaDataSetEvent:
			wm.MetaDataWriteModel.AppendEvents(&e.SetEvent)
		case *user.MetaDataRemovedEvent:
			wm.MetaDataWriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *UserMetaDataListWriteModel) Reduce() error {
	return wm.MetaDataWriteModel.Reduce()
}

func (wm *UserMetaDataListWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetaDataWriteModel.AggregateID).
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.MetaDataSetType,
			user.MetaDataRemovedType).
		Builder()
}
