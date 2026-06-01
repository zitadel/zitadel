package command

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type ApplicationSecretWriteModel struct {
	eventstore.WriteModel

	ApplicationID string
	ClientID      string
	HashedSecret  string

	State         domain.AppState
	SecretAllowed bool
	IsAPI         bool
}

func NewApplicationSecretWriteModel(projectID, applicationID, resourceOwner string) *ApplicationSecretWriteModel {
	return &ApplicationSecretWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		ApplicationID: applicationID,
	}
}

func (wm *ApplicationSecretWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *project.ApplicationRemovedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.OIDCConfigAddedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.OIDCConfigChangedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.APIConfigAddedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.APIConfigChangedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.OIDCConfigSecretChangedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.OIDCConfigSecretHashUpdatedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.APIConfigSecretChangedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.APIConfigSecretHashUpdatedEvent:
			if e.AppID != wm.ApplicationID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ProjectRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *ApplicationSecretWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *project.ApplicationRemovedEvent:
			wm.State = domain.AppStateRemoved
		case *project.OIDCConfigAddedEvent:
			wm.appendAddOIDCEvent(e)
		case *project.OIDCConfigChangedEvent:
			wm.appendChangeOIDCEvent(e)
		case *project.APIConfigAddedEvent:
			wm.appendAddAPIEvent(e)
		case *project.APIConfigChangedEvent:
			wm.appendChangeAPIEvent(e)
		case *project.OIDCConfigSecretChangedEvent:
			wm.HashedSecret = crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)
		case *project.OIDCConfigSecretHashUpdatedEvent:
			wm.HashedSecret = e.HashedSecret
		case *project.APIConfigSecretChangedEvent:
			wm.HashedSecret = crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)
		case *project.APIConfigSecretHashUpdatedEvent:
			wm.HashedSecret = e.HashedSecret
		case *project.ProjectRemovedEvent:
			wm.State = domain.AppStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ApplicationSecretWriteModel) appendAddOIDCEvent(e *project.OIDCConfigAddedEvent) {
	wm.State = domain.AppStateActive
	wm.ClientID = e.ClientID
	wm.SecretAllowed = e.AuthMethodType == domain.OIDCAuthMethodTypeBasic || e.AuthMethodType == domain.OIDCAuthMethodTypePost
	wm.IsAPI = false
}

func (wm *ApplicationSecretWriteModel) appendChangeOIDCEvent(e *project.OIDCConfigChangedEvent) {
	if e.AuthMethodType != nil {
		wm.SecretAllowed = *e.AuthMethodType == domain.OIDCAuthMethodTypeBasic || *e.AuthMethodType == domain.OIDCAuthMethodTypePost
	}
}

func (wm *ApplicationSecretWriteModel) appendAddAPIEvent(e *project.APIConfigAddedEvent) {
	wm.State = domain.AppStateActive
	wm.ClientID = e.ClientID
	wm.SecretAllowed = e.AuthMethodType == domain.APIAuthMethodTypeBasic
	wm.IsAPI = true
}

func (wm *ApplicationSecretWriteModel) appendChangeAPIEvent(e *project.APIConfigChangedEvent) {
	if e.AuthMethodType != nil {
		wm.SecretAllowed = *e.AuthMethodType == domain.APIAuthMethodTypeBasic
	}
}

func (wm *ApplicationSecretWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			project.ApplicationRemovedType,
			project.OIDCConfigAddedType,
			project.OIDCConfigChangedType,
			project.APIConfigAddedType,
			project.APIConfigChangedType,
			project.OIDCConfigSecretChangedType,
			project.OIDCConfigSecretHashUpdatedType,
			project.APIConfigSecretChangedType,
			project.APIConfigSecretHashUpdatedType,
			project.ProjectRemovedType).
		Builder()
}
