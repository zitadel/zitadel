package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

type UserMetadataWriteModel struct {
	MetadataWriteModel
}

func NewUserMetadataWriteModel(userID, resourceOwner, key string) *UserMetadataWriteModel {
	return &UserMetadataWriteModel{
		MetadataWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   userID,
				ResourceOwner: resourceOwner,
			},
			Key: key,
		},
	}
}

func (wm *UserMetadataWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MetadataSetEvent:
			wm.MetadataWriteModel.AppendEvents(&e.SetEvent)
		case *user.MetadataRemovedEvent:
			wm.MetadataWriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *UserMetadataWriteModel) Reduce() error {
	return wm.MetadataWriteModel.Reduce()
}

func (wm *UserMetadataWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataWriteModel.AggregateID).
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.MetadataSetType,
			user.MetadataRemovedType).
		Builder()
}

type UserMetadataListWriteModel struct {
	MetadataListWriteModel
}

func NewUserMetadataListWriteModel(userID, resourceOwner string) *UserMetadataListWriteModel {
	return &UserMetadataListWriteModel{
		MetadataListWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   userID,
				ResourceOwner: resourceOwner,
			},
			metaDataList: make(map[string][]byte),
		},
	}
}

func (wm *UserMetadataListWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MetadataSetEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.SetEvent)
		case *user.MetadataRemovedEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *UserMetadataListWriteModel) Reduce() error {
	return wm.MetadataListWriteModel.Reduce()
}

func (wm *UserMetadataListWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataListWriteModel.AggregateID).
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.MetadataSetType,
			user.MetadataRemovedType).
		Builder()
}
