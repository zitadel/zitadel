package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

type ApplicationWriteModel struct {
	eventstore.WriteModel

	AppID      string
	State      domain.AppState
	Name       string
	Type       domain.AppType
	OIDCConfig *domain.OIDCConfig
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
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *project.ApplicationAddedEvent:
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
			//case *project.ApplicationChangedEvent:
			//	wm.Name = e.Name
		}
	}
	return nil
}

func (wm *ApplicationWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.AggregateID).ResourceOwner(wm.ResourceOwner)
}
