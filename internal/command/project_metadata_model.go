package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type ProjectMetadataWriteModel struct {
	MetadataWriteModel
}

func NewProjectMetadataWriteModel(projectID, resourceOwner, key string) *ProjectMetadataWriteModel {
	return &ProjectMetadataWriteModel{
		MetadataWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   projectID,
				ResourceOwner: resourceOwner,
			},
			Key: key,
		},
	}
}

func (wm *ProjectMetadataWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.MetadataSetEvent:
			wm.MetadataWriteModel.AppendEvents(&e.SetEvent)
		case *project.MetadataRemovedEvent:
			wm.MetadataWriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *ProjectMetadataWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataWriteModel.AggregateID).
		AggregateTypes(project.AggregateType).
		EventTypes(
			project.MetadataSetType,
			project.MetadataRemovedType).
		Builder()
}

type ProjectMetadataListWriteModel struct {
	MetadataListWriteModel
}

func NewProjectMetadataListWriteModel(projectID, resourceOwner string) *ProjectMetadataListWriteModel {
	return &ProjectMetadataListWriteModel{
		MetadataListWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   projectID,
				ResourceOwner: resourceOwner,
			},
			metadataList: make(map[string][]byte),
		},
	}
}

func (wm *ProjectMetadataListWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.MetadataSetEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.SetEvent)
		case *project.MetadataRemovedEvent:
			wm.MetadataListWriteModel.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *ProjectMetadataListWriteModel) Reduce() error {
	return wm.MetadataListWriteModel.Reduce()
}

func (wm *ProjectMetadataListWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.MetadataListWriteModel.AggregateID).
		AggregateTypes(project.AggregateType).
		EventTypes(
			project.MetadataSetType,
			project.MetadataRemovedType).
		Builder()
}
