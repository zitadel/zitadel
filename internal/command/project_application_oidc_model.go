package command

import (
	"context"
	"reflect"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type OIDCApplicationWriteModel struct {
	eventstore.WriteModel

	AppID                    string
	AppName                  string
	ClientID                 string
	HashedSecret             string
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
	AdditionalOrigins        []string
	SkipNativeAppSuccessPage bool
	BackChannelLogoutURI     string
	oidc                     bool
}

func NewOIDCApplicationWriteModelWithAppID(projectID, appID, resourceOwner string) *OIDCApplicationWriteModel {
	return &OIDCApplicationWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
		AppID: appID,
	}
}

func NewOIDCApplicationWriteModel(projectID, resourceOwner string) *OIDCApplicationWriteModel {
	return &OIDCApplicationWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   projectID,
			ResourceOwner: resourceOwner,
		},
	}
}
func (wm *OIDCApplicationWriteModel) AppendEvents(events ...eventstore.Event) {
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
		case *project.OIDCConfigSecretHashUpdatedEvent:
			if e.AppID != wm.AppID {
				continue
			}
			wm.WriteModel.AppendEvents(e)
		case *project.ProjectRemovedEvent:
			wm.WriteModel.AppendEvents(e)
		}
	}
}

func (wm *OIDCApplicationWriteModel) Reduce() error {
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
			wm.appendAddOIDCEvent(e)
		case *project.OIDCConfigChangedEvent:
			wm.appendChangeOIDCEvent(e)
		case *project.OIDCConfigSecretChangedEvent:
			wm.HashedSecret = crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)
		case *project.OIDCConfigSecretHashUpdatedEvent:
			wm.HashedSecret = e.HashedSecret
		case *project.ProjectRemovedEvent:
			wm.State = domain.AppStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OIDCApplicationWriteModel) appendAddOIDCEvent(e *project.OIDCConfigAddedEvent) {
	wm.oidc = true
	wm.ClientID = e.ClientID
	wm.HashedSecret = crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)
	wm.RedirectUris = e.RedirectUris
	wm.ResponseTypes = e.ResponseTypes
	wm.GrantTypes = e.GrantTypes
	wm.ApplicationType = e.ApplicationType
	wm.AuthMethodType = e.AuthMethodType
	wm.PostLogoutRedirectUris = e.PostLogoutRedirectUris
	wm.OIDCVersion = e.Version
	wm.DevMode = e.DevMode
	wm.AccessTokenType = e.AccessTokenType
	wm.AccessTokenRoleAssertion = e.AccessTokenRoleAssertion
	wm.IDTokenRoleAssertion = e.IDTokenRoleAssertion
	wm.IDTokenUserinfoAssertion = e.IDTokenUserinfoAssertion
	wm.ClockSkew = e.ClockSkew
	wm.AdditionalOrigins = e.AdditionalOrigins
	wm.SkipNativeAppSuccessPage = e.SkipNativeAppSuccessPage
	wm.BackChannelLogoutURI = e.BackChannelLogoutURI
}

func (wm *OIDCApplicationWriteModel) appendChangeOIDCEvent(e *project.OIDCConfigChangedEvent) {
	if e.RedirectUris != nil {
		wm.RedirectUris = *e.RedirectUris
	}
	if e.ResponseTypes != nil {
		wm.ResponseTypes = *e.ResponseTypes
	}
	if e.GrantTypes != nil {
		wm.GrantTypes = *e.GrantTypes
	}
	if e.ApplicationType != nil {
		wm.ApplicationType = *e.ApplicationType
	}
	if e.AuthMethodType != nil {
		wm.AuthMethodType = *e.AuthMethodType
	}
	if e.PostLogoutRedirectUris != nil {
		wm.PostLogoutRedirectUris = *e.PostLogoutRedirectUris
	}
	if e.Version != nil {
		wm.OIDCVersion = *e.Version
	}
	if e.DevMode != nil {
		wm.DevMode = *e.DevMode
	}
	if e.AccessTokenType != nil {
		wm.AccessTokenType = *e.AccessTokenType
	}
	if e.AccessTokenRoleAssertion != nil {
		wm.AccessTokenRoleAssertion = *e.AccessTokenRoleAssertion
	}
	if e.IDTokenRoleAssertion != nil {
		wm.IDTokenRoleAssertion = *e.IDTokenRoleAssertion
	}
	if e.IDTokenUserinfoAssertion != nil {
		wm.IDTokenUserinfoAssertion = *e.IDTokenUserinfoAssertion
	}
	if e.ClockSkew != nil {
		wm.ClockSkew = *e.ClockSkew
	}
	if e.AdditionalOrigins != nil {
		wm.AdditionalOrigins = *e.AdditionalOrigins
	}
	if e.SkipNativeAppSuccessPage != nil {
		wm.SkipNativeAppSuccessPage = *e.SkipNativeAppSuccessPage
	}
	if e.BackChannelLogoutURI != nil {
		wm.BackChannelLogoutURI = *e.BackChannelLogoutURI
	}
}

func (wm *OIDCApplicationWriteModel) Query() *eventstore.SearchQueryBuilder {
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
			project.OIDCConfigAddedType,
			project.OIDCConfigChangedType,
			project.OIDCConfigSecretChangedType,
			project.OIDCConfigSecretHashUpdatedType,
			project.ProjectRemovedType,
		).Builder()
}

func (wm *OIDCApplicationWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	redirectURIS,
	postLogoutRedirectURIs []string,
	responseTypes []domain.OIDCResponseType,
	grantTypes []domain.OIDCGrantType,
	appType domain.OIDCApplicationType,
	authMethodType domain.OIDCAuthMethodType,
	oidcVersion domain.OIDCVersion,
	accessTokenType domain.OIDCTokenType,
	devMode,
	accessTokenRoleAssertion,
	idTokenRoleAssertion,
	idTokenUserinfoAssertion bool,
	clockSkew time.Duration,
	additionalOrigins []string,
	skipNativeAppSuccessPage bool,
	backChannelLogoutURI string,
) (*project.OIDCConfigChangedEvent, bool, error) {
	changes := make([]project.OIDCConfigChanges, 0)
	var err error

	if !reflect.DeepEqual(wm.RedirectUris, redirectURIS) {
		changes = append(changes, project.ChangeRedirectURIs(redirectURIS))
	}
	if !reflect.DeepEqual(wm.ResponseTypes, responseTypes) {
		changes = append(changes, project.ChangeResponseTypes(responseTypes))
	}
	if !reflect.DeepEqual(wm.GrantTypes, grantTypes) {
		changes = append(changes, project.ChangeGrantTypes(grantTypes))
	}
	if wm.ApplicationType != appType {
		changes = append(changes, project.ChangeApplicationType(appType))
	}
	if wm.AuthMethodType != authMethodType {
		changes = append(changes, project.ChangeAuthMethodType(authMethodType))
	}
	if !reflect.DeepEqual(wm.PostLogoutRedirectUris, postLogoutRedirectURIs) {
		changes = append(changes, project.ChangePostLogoutRedirectURIs(postLogoutRedirectURIs))
	}
	if wm.OIDCVersion != oidcVersion {
		changes = append(changes, project.ChangeVersion(oidcVersion))
	}
	if wm.DevMode != devMode {
		changes = append(changes, project.ChangeDevMode(devMode))
	}
	if wm.AccessTokenType != accessTokenType {
		changes = append(changes, project.ChangeAccessTokenType(accessTokenType))
	}
	if wm.AccessTokenRoleAssertion != accessTokenRoleAssertion {
		changes = append(changes, project.ChangeAccessTokenRoleAssertion(accessTokenRoleAssertion))
	}
	if wm.IDTokenRoleAssertion != idTokenRoleAssertion {
		changes = append(changes, project.ChangeIDTokenRoleAssertion(idTokenRoleAssertion))
	}
	if wm.IDTokenUserinfoAssertion != idTokenUserinfoAssertion {
		changes = append(changes, project.ChangeIDTokenUserinfoAssertion(idTokenUserinfoAssertion))
	}
	if wm.ClockSkew != clockSkew {
		changes = append(changes, project.ChangeClockSkew(clockSkew))
	}
	if !reflect.DeepEqual(wm.AdditionalOrigins, additionalOrigins) {
		changes = append(changes, project.ChangeAdditionalOrigins(additionalOrigins))
	}
	if wm.SkipNativeAppSuccessPage != skipNativeAppSuccessPage {
		changes = append(changes, project.ChangeSkipNativeAppSuccessPage(skipNativeAppSuccessPage))
	}
	if wm.BackChannelLogoutURI != backChannelLogoutURI {
		changes = append(changes, project.ChangeBackChannelLogoutURI(backChannelLogoutURI))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := project.NewOIDCConfigChangedEvent(ctx, aggregate, appID, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func (wm *OIDCApplicationWriteModel) IsOIDC() bool {
	return wm.oidc
}
