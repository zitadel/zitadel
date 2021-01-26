package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

type ApplicationWriteModel struct {
	eventstore.WriteModel

	AppID string
	State domain.AppState
	Name  string
	Type  domain.AppType
}

func NewApplicationWriteModelWithAppIDC(projectID, appID, resourceOwner string) *ApplicationWriteModel {
	return &ApplicationWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		AppID: appID,
	}
}

func NewApplicationWriteModel(projectID, resourceOwner string) *ApplicationWriteModel {
	return &ApplicationWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
	}
}
func (wm *ApplicationWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.ApplicationAddedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ApplicationChangedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ApplicationDeactivatedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ApplicationReactivatedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ApplicationRemovedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ProjectRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *ApplicationWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ApplicationAddedEvent:
			wm.Name = e.Name
			wm.State = domain.AppStateActive
		case *project.ApplicationChangedEvent:
			wm.Name = e.Name
		case *project.ApplicationDeactivatedEvent:
			if wm.State == domain.AppStateRemoved {
				continue
			}
			wm.State = domain.AppStateInactive
		case *project.ApplicationReactivatedEvent:
			if wm.State == domain.AppStateRemoved {
				continue
			}
			wm.State = domain.AppStateActive
		case *project.ApplicationRemovedEvent:
			wm.State = domain.AppStateRemoved
		case *project.ProjectRemovedEvent:
			wm.State = domain.AppStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ApplicationWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
	//EventTypes(
	//	project.ApplicationAddedType,
	//	project.ApplicationChangedType,
	//	project.ApplicationDeactivatedType,
	//	project.ApplicationReactivatedType,
	//	project.ApplicationRemovedType,
	//	project.ProjectRemovedType,
	//)
}
