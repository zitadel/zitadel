package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OrgMetadataWriteModel struct {
	MetadataWriteModel
}

func NewOrgMetadataWriteModel(orgID, key string) *OrgMetadataWriteModel {
	return &OrgMetadataWriteModel{
		MetadataWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			Key: key,
		},
	}
}

func (wm *OrgMetadataWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.MetadataSetEvent:
			wm.MetadataWriteModel.AppendEvents(&e.SetEvent)
		case *org.MetadataRemovedEvent:
			wm.MetadataWriteModel.AppendEvents(&e.RemovedEvent)
		case *org.MetadataRemovedAllEvent:
			wm.MetadataWriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *OrgMetadataWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataWriteModel.AggregateID).
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.MetadataSetType,
			org.MetadataRemovedType,
			org.MetadataRemovedAllType).
		Builder()
}

type OrgMetadataListWriteModel struct {
	MetadataListWriteModel
}

func NewOrgMetadataListWriteModel(orgID string) *OrgMetadataListWriteModel {
	return &OrgMetadataListWriteModel{
		MetadataListWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			metadataList: make(map[string][]byte),
		},
	}
}

func (wm *OrgMetadataListWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.MetadataSetEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.SetEvent)
		case *org.MetadataRemovedEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.RemovedEvent)
		case *org.MetadataRemovedAllEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.RemovedAllEvent)
		}
	}
}

func (wm *OrgMetadataListWriteModel) Reduce() error {
	return wm.MetadataListWriteModel.Reduce()
}

func (wm *OrgMetadataListWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataListWriteModel.AggregateID).
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.MetadataSetType,
			org.MetadataRemovedType,
			org.MetadataRemovedAllType).
		Builder()
}
