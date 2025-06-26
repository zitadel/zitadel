package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type SAMLApplicationWriteModel struct {
	eventstore.WriteModel

	AppID        string
	AppName      string
	EntityID     string
	Metadata     []byte
	MetadataURL  string
	LoginVersion domain.LoginVersion
	LoginBaseURI string

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

func (wm *SAMLApplicationWriteModel) AppendEvents(events ...eventstore.Event) {
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
		case *project.SAMLConfigAddedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.SAMLConfigChangedEvent:
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
	wm.EntityID = e.EntityID
	wm.LoginVersion = e.LoginVersion
	wm.LoginBaseURI = e.LoginBaseURI
}

func (wm *SAMLApplicationWriteModel) appendChangeSAMLEvent(e *project.SAMLConfigChangedEvent) {
	wm.saml = true
	if e.Metadata != nil {
		wm.Metadata = e.Metadata
	}
	if e.MetadataURL != nil {
		wm.MetadataURL = *e.MetadataURL
	}
	if e.EntityID != "" {
		wm.EntityID = e.EntityID
	}
	if e.LoginVersion != nil {
		wm.LoginVersion = *e.LoginVersion
	}
	if e.LoginBaseURI != nil {
		wm.LoginBaseURI = *e.LoginBaseURI
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
	metadata []byte,
	metadataURL *string,
	loginVersion *domain.LoginVersion,
	loginBaseURI *string,
) (*project.SAMLConfigChangedEvent, bool, error) {
	changes := make([]project.SAMLConfigChanges, 0)
	var err error
	if metadata != nil && !slices.Equal(wm.Metadata, metadata) {
		changes = append(changes, project.ChangeMetadata(metadata))
	}
	if metadataURL != nil && wm.MetadataURL != *metadataURL {
		changes = append(changes, project.ChangeMetadataURL(*metadataURL))
	}
	if wm.EntityID != entityID {
		changes = append(changes, project.ChangeEntityID(entityID))
	}
	if loginVersion != nil && wm.LoginVersion != *loginVersion {
		changes = append(changes, project.ChangeSAMLLoginVersion(*loginVersion))
	}
	if loginBaseURI != nil && wm.LoginBaseURI != *loginBaseURI {
		changes = append(changes, project.ChangeSAMLLoginBaseURI(*loginBaseURI))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := project.NewSAMLConfigChangedEvent(ctx, aggregate, appID, wm.EntityID, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func (wm *SAMLApplicationWriteModel) IsSAML() bool {
	return wm.saml
}

type AppIDToEntityID struct {
	AppID    string
	EntityID string
}

type SAMLEntityIDsWriteModel struct {
	eventstore.WriteModel

	EntityIDs []*AppIDToEntityID
}

func NewSAMLEntityIDsWriteModel(projectID, resourceOwner string) *SAMLEntityIDsWriteModel {
	return &SAMLEntityIDsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		EntityIDs: []*AppIDToEntityID{},
	}
}

func (wm *SAMLEntityIDsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			project.ApplicationRemovedType,
			project.SAMLConfigAddedType,
			project.SAMLConfigChangedType).
		Builder()
}

func (wm *SAMLEntityIDsWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.ApplicationRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		case *project.SAMLConfigAddedEvent:
			wm.WriteModel.AppendEvents(e)
		case *project.SAMLConfigChangedEvent:
			if e.EntityID != "" {
				wm.WriteModel.AppendEvents(e)
			}
		}
	}
}

func (wm *SAMLEntityIDsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ApplicationRemovedEvent:
			wm.EntityIDs = removeAppIDFromEntityIDs(wm.EntityIDs, e.AppID)
		case *project.SAMLConfigAddedEvent:
			wm.EntityIDs = append(wm.EntityIDs, &AppIDToEntityID{AppID: e.AppID, EntityID: e.EntityID})
		case *project.SAMLConfigChangedEvent:
			for _, item := range wm.EntityIDs {
				if e.AppID == item.AppID && e.EntityID != "" {
					item.EntityID = e.EntityID
				}
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func removeAppIDFromEntityIDs(items []*AppIDToEntityID, appID string) []*AppIDToEntityID {
	for i := len(items) - 1; i >= 0; i-- {
		if items[i].AppID == appID {
			items[i] = items[len(items)-1]
			items[len(items)-1] = nil
			items = items[:len(items)-1]
		}
	}
	return items
}
