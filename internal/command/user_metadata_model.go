package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/metadata"
	"github.com/zitadel/zitadel/internal/repository/user"
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

func (wm *UserMetadataWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MetadataSetEvent:
			wm.MetadataWriteModel.AppendEvents(&e.SetEvent)
		case *user.MetadataRemovedEvent:
			wm.MetadataWriteModel.AppendEvents(&e.RemovedEvent)
		case *user.MetadataRemovedAllEvent:
			wm.MetadataWriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *UserMetadataWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataWriteModel.AggregateID).
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.MetadataSetType,
			user.MetadataRemovedType,
			user.MetadataRemovedAllType).
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
			metadataList: make(map[string][]byte),
		},
	}
}

func (wm *UserMetadataListWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MetadataSetEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.SetEvent)
		case *user.MetadataRemovedEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.RemovedEvent)
		case *user.MetadataRemovedAllEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.RemovedAllEvent)
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
			user.MetadataRemovedType,
			user.MetadataRemovedAllType).
		Builder()
}

type UserMetadataByOrgListWriteModel struct {
	eventstore.WriteModel
	UserMetadata map[string]map[string][]byte
}

func NewUserMetadataByOrgListWriteModel(resourceOwner string) *UserMetadataByOrgListWriteModel {
	return &UserMetadataByOrgListWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: resourceOwner,
		},
		UserMetadata: make(map[string]map[string][]byte),
	}
}

func (wm *UserMetadataByOrgListWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.MetadataSetEvent:
			wm.WriteModel.AppendEvents(&e.SetEvent)
		case *user.MetadataRemovedEvent:
			wm.WriteModel.AppendEvents(&e.RemovedEvent)
		case *user.MetadataRemovedAllEvent:
			wm.WriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *UserMetadataByOrgListWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *metadata.SetEvent:
			if val, ok := wm.UserMetadata[e.Aggregate().ID]; ok {
				val[e.Key] = e.Value
			} else {
				wm.UserMetadata[e.Aggregate().ID] = map[string][]byte{
					e.Key: e.Value,
				}
			}
		case *metadata.RemovedEvent:
			if val, ok := wm.UserMetadata[e.Aggregate().ID]; ok {
				delete(val, e.Key)
			}
		case *metadata.RemovedAllEvent:
			delete(wm.UserMetadata, e.Aggregate().ID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserMetadataByOrgListWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.MetadataSetType,
			user.MetadataRemovedType,
			user.MetadataRemovedAllType).
		Builder()
}
