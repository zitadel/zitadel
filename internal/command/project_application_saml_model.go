package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

type SAMLApplicationWriteModel struct {
	eventstore.WriteModel

	AppID       string
	AppName     string
	Metadata    string
	MetadataURL string

	State domain.AppState
	saml  bool
}

func NewSAMLApplicationWriteModelWithAppID(projectID, appID, resourceOwner string) *SAMLApplicationWriteModel {
	return &SAMLApplicationWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		AppID: appID,
	}
}

func NewSAMLApplicationWriteModel(projectID, resourceOwner string) *SAMLApplicationWriteModel {
	return &SAMLApplicationWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *SAMLApplicationWriteModel) AppendEvents(events ...eventstore.EventReader) {
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
		case *project.OIDCConfigAddedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.OIDCConfigChangedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.OIDCConfigSecretChangedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ProjectRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *SAMLApplicationWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ApplicationAddedEvent:
			wm.AppName = e.Name
			wm.State = domain.AppStateActive
		case *project.ApplicationChangedEvent:
			wm.AppName = e.Name
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
		case *project.SAMLConfigAddedEvent:
			wm.appendAddSAMLEvent(e)
		case *project.SAMLConfigChangedEvent:
			wm.appendChangeSAMLEvent(e)
		case *project.ProjectRemovedEvent:
			wm.State = domain.AppStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *SAMLApplicationWriteModel) appendAddSAMLEvent(e *project.SAMLConfigAddedEvent) {
	wm.saml = true
	wm.Metadata = e.Metadata
	wm.MetadataURL = e.MetadataURL
}

func (wm *SAMLApplicationWriteModel) appendChangeSAMLEvent(e *project.SAMLConfigChangedEvent) {
	if e.Metadata != nil {
		wm.Metadata = *e.Metadata
	}
	if e.MetadataURL != nil {
		wm.Metadata = *e.MetadataURL
	}
}

func (wm *SAMLApplicationWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			project.ApplicationAddedType,
			project.ApplicationChangedType,
			project.ApplicationDeactivatedType,
			project.ApplicationReactivatedType,
			project.ApplicationRemovedType,
			project.SAMLConfigAddedType,
			project.SAMLConfigChangedType,
			project.ProjectRemovedType).
		Builder()
}

func (wm *SAMLApplicationWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	entityID string,
	metadata string,
	metadataURL string,
) (*project.SAMLConfigChangedEvent, bool, error) {
	changes := make([]project.SAMLConfigChanges, 0)
	var err error
	if wm.Metadata != metadata {
		changes = append(changes, project.ChangeMetadata(metadata))
	}
	if wm.MetadataURL != metadataURL {
		changes = append(changes, project.ChangeMetadataURL(metadataURL))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := project.NewSAMLConfigChangedEvent(ctx, aggregate, appID, entityID, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func (wm *SAMLApplicationWriteModel) IsSAML() bool {
	return wm.saml
}
