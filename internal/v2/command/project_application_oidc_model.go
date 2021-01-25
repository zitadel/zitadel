package command

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
	"time"
)

type ApplicationOIDCConfigWriteModel struct {
	eventstore.WriteModel

	AppID                    string
	AppName                  string
	ClientID                 string
	ClientSecret             *crypto.CryptoValue
	ClientSecretString       string
	RedirectUris             []string
	ResponseTypes            []domain.OIDCResponseType
	GrantTypes               []domain.OIDCGrantType
	ApplicationType          domain.OIDCApplicationType
	AuthMethodType           domain.OIDCAuthMethodType
	PostLogoutRedirectUris   []string
	OIDCVersion              domain.OIDCVersion
	Compliance               *domain.Compliance
	DevMode                  bool
	AccessTokenType          domain.OIDCTokenType
	AccessTokenRoleAssertion bool
	IDTokenRoleAssertion     bool
	IDTokenUserinfoAssertion bool
	ClockSkew                time.Duration
	State                    domain.AppState
}

func NewApplicationOIDCConfigWriteModelWithAppIDC(projectID, appID, resourceOwner string) *ApplicationOIDCConfigWriteModel {
	return &ApplicationOIDCConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		AppID: appID,
	}
}

func NewApplicationOIDCConfigWriteModel(projectID, resourceOwner string) *ApplicationOIDCConfigWriteModel {
	return &ApplicationOIDCConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
	}
}
func (wm *ApplicationOIDCConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
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

func (wm *ApplicationOIDCConfigWriteModel) Reduce() error {
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
		case *project.OIDCConfigAddedEvent:

		case *project.OIDCConfigChangedEvent:

		case *project.OIDCConfigSecretChangedEvent:

		case *project.ProjectRemovedEvent:
			wm.State = domain.AppStateRemoved
		}
	}
	return nil
}

func (wm *ApplicationOIDCConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, project.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			project.ApplicationAddedType,
			project.ApplicationChangedType,
			project.ApplicationDeactivatedType,
			project.ApplicationReactivatedType,
			project.ApplicationRemovedType,
			project.OIDCConfigAddedType,
			project.OIDCConfigChangedType,
			project.OIDCConfigSecretChangedType,
			project.ProjectRemovedType,
		)
}
