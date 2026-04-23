package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/metadata"
)

type GroupMetadataWriteModel struct {
	MetadataWriteModel
}

func NewGroupMetadataWriteModel(groupID, resourceOwner, key string) *GroupMetadataWriteModel {
	return &GroupMetadataWriteModel{
		MetadataWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   groupID,
				ResourceOwner: resourceOwner,
			},
			Key: key,
		},
	}
}

func (wm *GroupMetadataWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *group.MetadataSetEvent:
			wm.MetadataWriteModel.AppendEvents(&e.SetEvent)
		case *group.MetadataRemovedEvent:
			wm.MetadataWriteModel.AppendEvents(&e.RemovedEvent)
		case *group.MetadataRemovedAllEvent:
			wm.MetadataWriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *GroupMetadataWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataWriteModel.AggregateID).
		AggregateTypes(group.AggregateType).
		EventTypes(
			group.MetadataSetType,
			group.MetadataRemovedType,
			group.MetadataRemovedAllType).
		Builder()
}

type GroupMetadataListWriteModel struct {
	MetadataListWriteModel
}

func NewGroupMetadataListWriteModel(groupID, resourceOwner string) *GroupMetadataListWriteModel {
	return &GroupMetadataListWriteModel{
		MetadataListWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   groupID,
				ResourceOwner: resourceOwner,
			},
			metadataList: make(map[string][]byte),
		},
	}
}

func (wm *GroupMetadataListWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *group.MetadataSetEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.SetEvent)
		case *group.MetadataRemovedEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.RemovedEvent)
		case *group.MetadataRemovedAllEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *GroupMetadataListWriteModel) Reduce() error {
	return wm.MetadataListWriteModel.Reduce()
}

func (wm *GroupMetadataListWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataListWriteModel.AggregateID).
		AggregateTypes(group.AggregateType).
		EventTypes(
			group.MetadataSetType,
			group.MetadataRemovedType,
			group.MetadataRemovedAllType).
		Builder()
}

type GroupMetadataByOrgListWriteModel struct {
	eventstore.WriteModel
	GroupMetadata map[string]map[string][]byte
}

func NewGroupMetadataByOrgListWriteModel(resourceOwner string) *GroupMetadataByOrgListWriteModel {
	return &GroupMetadataByOrgListWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: resourceOwner,
		},
		GroupMetadata: make(map[string]map[string][]byte),
	}
}

func (wm *GroupMetadataByOrgListWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *group.MetadataSetEvent:
			wm.WriteModel.AppendEvents(&e.SetEvent)
		case *group.MetadataRemovedEvent:
			wm.WriteModel.AppendEvents(&e.RemovedEvent)
		case *group.MetadataRemovedAllEvent:
			wm.WriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *GroupMetadataByOrgListWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *metadata.SetEvent:
			if val, ok := wm.GroupMetadata[e.Aggregate().ID]; ok {
				val[e.Key] = e.Value
			} else {
				wm.GroupMetadata[e.Aggregate().ID] = map[string][]byte{
					e.Key: e.Value,
				}
			}
		case *metadata.RemovedEvent:
			if val, ok := wm.GroupMetadata[e.Aggregate().ID]; ok {
				delete(val, e.Key)
			}
		case *metadata.RemovedAllEvent:
			delete(wm.GroupMetadata, e.Aggregate().ID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GroupMetadataByOrgListWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(group.AggregateType).
		EventTypes(
			group.MetadataSetType,
			group.MetadataRemovedType,
			group.MetadataRemovedAllType).
		Builder()
}
