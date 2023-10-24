package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/project/model"
	"github.com/zitadel/zitadel/internal/repository/project"
)

const (
	ApplicationKeyID            = "id"
	ApplicationKeyProjectID     = "project_id"
	ApplicationKeyResourceOwner = "resource_owner"
	ApplicationKeyOIDCClientID  = "oidc_client_id"
	ApplicationKeyName          = "app_name"
)

type ApplicationView struct {
	ID                     string                        `json:"appId" gorm:"column:id;primary_key"`
	ProjectID              string                        `json:"-" gorm:"column:project_id"`
	Name                   string                        `json:"name" gorm:"column:app_name"`
	CreationDate           time.Time                     `json:"-" gorm:"column:creation_date"`
	ChangeDate             time.Time                     `json:"-" gorm:"column:change_date"`
	State                  int32                         `json:"-" gorm:"column:app_state"`
	ResourceOwner          string                        `json:"-" gorm:"column:resource_owner"`
	ProjectRoleAssertion   bool                          `json:"projectRoleAssertion" gorm:"column:project_role_assertion"`
	ProjectRoleCheck       bool                          `json:"projectRoleCheck" gorm:"column:project_role_check"`
	HasProjectCheck        bool                          `json:"hasProjectCheck" gorm:"column:has_project_check"`
	PrivateLabelingSetting domain.PrivateLabelingSetting `json:"privateLabelingSetting" gorm:"column:private_labeling_setting"`

	IsOIDC                     bool                                    `json:"-" gorm:"column:is_oidc"`
	OIDCVersion                int32                                   `json:"oidcVersion" gorm:"column:oidc_version"`
	OIDCClientID               string                                  `json:"clientId" gorm:"column:oidc_client_id"`
	OIDCRedirectUris           database.TextArray[string]              `json:"redirectUris" gorm:"column:oidc_redirect_uris"`
	OIDCResponseTypes          database.Array[domain.OIDCResponseType] `json:"responseTypes" gorm:"column:oidc_response_types"`
	OIDCGrantTypes             database.Array[domain.OIDCGrantType]    `json:"grantTypes" gorm:"column:oidc_grant_types"`
	OIDCApplicationType        int32                                   `json:"applicationType" gorm:"column:oidc_application_type"`
	OIDCAuthMethodType         int32                                   `json:"authMethodType" gorm:"column:oidc_auth_method_type"`
	OIDCPostLogoutRedirectUris database.TextArray[string]              `json:"postLogoutRedirectUris" gorm:"column:oidc_post_logout_redirect_uris"`
	NoneCompliant              bool                                    `json:"-" gorm:"column:none_compliant"`
	ComplianceProblems         database.TextArray[string]              `json:"-" gorm:"column:compliance_problems"`
	DevMode                    bool                                    `json:"devMode" gorm:"column:dev_mode"`
	OriginAllowList            database.TextArray[string]              `json:"-" gorm:"column:origin_allow_list"`
	AdditionalOrigins          database.TextArray[string]              `json:"additionalOrigins" gorm:"column:additional_origins"`
	AccessTokenType            int32                                   `json:"accessTokenType" gorm:"column:access_token_type"`
	AccessTokenRoleAssertion   bool                                    `json:"accessTokenRoleAssertion" gorm:"column:access_token_role_assertion"`
	IDTokenRoleAssertion       bool                                    `json:"idTokenRoleAssertion" gorm:"column:id_token_role_assertion"`
	IDTokenUserinfoAssertion   bool                                    `json:"idTokenUserinfoAssertion" gorm:"column:id_token_userinfo_assertion"`
	ClockSkew                  time.Duration                           `json:"clockSkew" gorm:"column:clock_skew"`

	IsSAML      bool   `json:"-" gorm:"column:is_saml"`
	Metadata    []byte `json:"metadata" gorm:"column:metadata"`
	MetadataURL string `json:"metadata_url" gorm:"column:metadata_url"`

	Sequence uint64 `json:"-" gorm:"sequence"`
}

func OIDCResponseTypesToModel(oidctypes []domain.OIDCResponseType) []model.OIDCResponseType {
	result := make([]model.OIDCResponseType, len(oidctypes))
	for i, t := range oidctypes {
		result[i] = model.OIDCResponseType(t)
	}
	return result
}

func OIDCGrantTypesToModel(granttypes []domain.OIDCGrantType) []model.OIDCGrantType {
	result := make([]model.OIDCGrantType, len(granttypes))
	for i, t := range granttypes {
		result[i] = model.OIDCGrantType(t)
	}
	return result
}

func (a *ApplicationView) AppendEventIfMyApp(event *models.Event) (err error) {
	view := new(ApplicationView)
	switch event.Type() {
	case project.ApplicationAddedType:
		err = view.SetData(event)
		if err != nil {
			return err
		}
	case project.ApplicationChangedType,
		project.OIDCConfigAddedType,
		project.OIDCConfigChangedType,
		project.APIConfigAddedType,
		project.APIConfigChangedType,
		project.ApplicationDeactivatedType,
		project.ApplicationReactivatedType,
		project.SAMLConfigAddedType,
		project.SAMLConfigChangedType:
		err = view.SetData(event)
		if err != nil {
			return err
		}
	case project.ApplicationRemovedType:
		err = view.SetData(event)
		if err != nil {
			return err
		}
	case project.ProjectChangedType:
		return a.AppendEvent(event)
	case project.ProjectRemovedType:
		return a.AppendEvent(event)
	default:
		return nil
	}
	if view.ID == a.ID {
		return a.AppendEvent(event)
	}
	return nil
}

func (a *ApplicationView) AppendEvent(event *models.Event) (err error) {
	a.Sequence = event.Seq
	a.ChangeDate = event.CreationDate
	switch event.Type() {
	case project.ApplicationAddedType:
		a.setRootData(event)
		a.CreationDate = event.CreationDate
		a.ResourceOwner = event.ResourceOwner
		err = a.SetData(event)
	case project.OIDCConfigAddedType:
		a.IsOIDC = true
		err = a.SetData(event)
		if err != nil {
			return err
		}
		a.setCompliance()
		return a.setOriginAllowList()
	case project.SAMLConfigAddedType:
		a.IsSAML = true
		return a.SetData(event)
	case project.APIConfigAddedType:
		a.IsOIDC = false
		return a.SetData(event)
	case project.ApplicationChangedType:
		return a.SetData(event)
	case project.OIDCConfigChangedType:
		err = a.SetData(event)
		if err != nil {
			return err
		}
		a.setCompliance()
		return a.setOriginAllowList()
	case project.SAMLConfigChangedType:
		return a.SetData(event)
	case project.APIConfigChangedType:
		return a.SetData(event)
	case project.ProjectChangedType:
		return a.setProjectChanges(event)
	case project.ApplicationDeactivatedType:
		a.State = int32(model.AppStateInactive)
	case project.ApplicationReactivatedType:
		a.State = int32(model.AppStateActive)
	case project.ApplicationRemovedType, project.ProjectRemovedType:
		a.State = int32(model.AppStateRemoved)
	}
	return err
}

func (a *ApplicationView) setRootData(event *models.Event) {
	a.ProjectID = event.AggregateID
}

func (a *ApplicationView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-lo9ds").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-8suie", "Could not unmarshal data")
	}
	return nil
}

func (a *ApplicationView) setOriginAllowList() error {
	allowList := make(database.TextArray[string], 0)
	for _, redirect := range a.OIDCRedirectUris {
		origin, err := http_util.GetOriginFromURLString(redirect)
		if err != nil {
			return err
		}
		if !http_util.IsOriginAllowed(allowList, origin) {
			allowList = append(allowList, origin)
		}
	}
	for _, origin := range a.AdditionalOrigins {
		if !http_util.IsOriginAllowed(allowList, origin) {
			allowList = append(allowList, origin)
		}
	}
	a.OriginAllowList = allowList
	return nil
}

func (a *ApplicationView) setCompliance() {
	compliance := model.GetOIDCCompliance(model.OIDCVersion(a.OIDCVersion), model.OIDCApplicationType(a.OIDCApplicationType), OIDCGrantTypesToModel(a.OIDCGrantTypes), OIDCResponseTypesToModel(a.OIDCResponseTypes), model.OIDCAuthMethodType(a.OIDCAuthMethodType), a.OIDCRedirectUris)
	a.NoneCompliant = compliance.NoneCompliant
	a.ComplianceProblems = compliance.Problems
}

func (a *ApplicationView) setProjectChanges(event *models.Event) error {
	changes := struct {
		ProjectRoleAssertion   *bool                          `json:"projectRoleAssertion,omitempty"`
		ProjectRoleCheck       *bool                          `json:"projectRoleCheck,omitempty"`
		HasProjectCheck        *bool                          `json:"hasProjectCheck,omitempty"`
		PrivateLabelingSetting *domain.PrivateLabelingSetting `json:"privateLabelingSetting,omitempty"`
	}{}
	if err := json.Unmarshal(event.Data, &changes); err != nil {
		logging.Log("EVEN-DFbfg").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Bw221", "Could not unmarshal data")
	}
	if changes.ProjectRoleAssertion != nil {
		a.ProjectRoleAssertion = *changes.ProjectRoleAssertion
	}
	if changes.ProjectRoleCheck != nil {
		a.ProjectRoleCheck = *changes.ProjectRoleCheck
	}
	if changes.HasProjectCheck != nil {
		a.HasProjectCheck = *changes.HasProjectCheck
	}
	if changes.PrivateLabelingSetting != nil {
		a.PrivateLabelingSetting = *changes.PrivateLabelingSetting
	}
	return nil
}
