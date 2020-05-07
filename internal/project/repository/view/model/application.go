package model

import (
	"encoding/json"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/lib/pq"
	"time"
)

const (
	ApplicationKeyID            = "id"
	ApplicationKeyProjectID     = "project_id"
	ApplicationKeyResourceOwner = "resource_owner"
	ApplicationKeyOIDCClientID  = "oidc_client_id"
	ApplicationKeyName          = "name"
)

type ApplicationView struct {
	ID           string    `json:"appId" gorm:"column:id;primary_key"`
	ProjectID    string    `json:"-" gorm:"column:project_id"`
	Name         string    `json:"name" gorm:"column:app_name"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:app_state"`

	IsOIDC                     bool           `json:"-" gorm:"column:is_oidc"`
	OIDCClientID               string         `json:"clientId" gorm:"column:oidc_client_id"`
	OIDCRedirectUris           pq.StringArray `json:"redirectUris" gorm:"column:oidc_redirect_uris"`
	OIDCResponseTypes          pq.Int64Array  `json:"responseTypes" gorm:"column:oidc_response_types"`
	OIDCGrantTypes             pq.Int64Array  `json:"grantTypes" gorm:"column:oidc_grant_types"`
	OIDCApplicationType        int32          `json:"applicationType" gorm:"column:oidc_application_type"`
	OIDCAuthMethodType         int32          `json:"authMethodType" gorm:"column:oidc_auth_method_type"`
	OIDCPostLogoutRedirectUris pq.StringArray `json:"redirectUris" gorm:"column:oidc_post_logout_redirect_uris"`

	Sequence uint64 `json:"-" gorm:"sequence"`
}

func ApplicationViewFromModel(app *model.ApplicationView) *ApplicationView {
	return &ApplicationView{
		ID:           app.ID,
		ProjectID:    app.ProjectID,
		Name:         app.Name,
		State:        int32(app.State),
		Sequence:     app.Sequence,
		CreationDate: app.CreationDate,
		ChangeDate:   app.ChangeDate,

		IsOIDC:                     app.IsOIDC,
		OIDCClientID:               app.OIDCClientID,
		OIDCRedirectUris:           app.OIDCRedirectUris,
		OIDCResponseTypes:          OIDCResponseTypesFromModel(app.OIDCResponseTypes),
		OIDCGrantTypes:             OIDCGrantTypesFromModel(app.OIDCGrantTypes),
		OIDCApplicationType:        int32(app.OIDCApplicationType),
		OIDCAuthMethodType:         int32(app.OIDCAuthMethodType),
		OIDCPostLogoutRedirectUris: app.OIDCPostLogoutRedirectUris,
	}
}

func OIDCResponseTypesFromModel(oidctypes []model.OIDCResponseType) []int64 {
	result := make([]int64, 0)
	for _, t := range oidctypes {
		result = append(result, int64(t))
	}
	return result
}

func OIDCGrantTypesFromModel(granttypes []model.OIDCGrantType) []int64 {
	result := make([]int64, 0)
	for _, t := range granttypes {
		result = append(result, int64(t))
	}
	return result
}

func ApplicationViewToModel(app *ApplicationView) *model.ApplicationView {
	return &model.ApplicationView{
		ID:           app.ID,
		ProjectID:    app.ProjectID,
		Name:         app.Name,
		State:        model.AppState(app.State),
		Sequence:     app.Sequence,
		CreationDate: app.CreationDate,
		ChangeDate:   app.ChangeDate,

		IsOIDC:                     app.IsOIDC,
		OIDCClientID:               app.OIDCClientID,
		OIDCRedirectUris:           app.OIDCRedirectUris,
		OIDCResponseTypes:          OIDCResponseTypesToModel(app.OIDCResponseTypes),
		OIDCGrantTypes:             OIDCGrantTypesToModel(app.OIDCGrantTypes),
		OIDCApplicationType:        model.OIDCApplicationType(app.OIDCApplicationType),
		OIDCAuthMethodType:         model.OIDCAuthMethodType(app.OIDCAuthMethodType),
		OIDCPostLogoutRedirectUris: app.OIDCPostLogoutRedirectUris,
	}
}

func OIDCResponseTypesToModel(oidctypes []int64) []model.OIDCResponseType {
	result := make([]model.OIDCResponseType, 0)
	for _, t := range oidctypes {
		result = append(result, model.OIDCResponseType(t))
	}
	return result
}

func OIDCGrantTypesToModel(granttypes []int64) []model.OIDCGrantType {
	result := make([]model.OIDCGrantType, 0)
	for _, t := range granttypes {
		result = append(result, model.OIDCGrantType(t))
	}
	return result
}

func ApplicationViewsToModel(roles []*ApplicationView) []*model.ApplicationView {
	result := make([]*model.ApplicationView, 0)
	for _, r := range roles {
		result = append(result, ApplicationViewToModel(r))
	}
	return result
}

func (a *ApplicationView) AppendEvent(event *models.Event) error {
	a.Sequence = event.Sequence
	switch event.Type {
	case es_model.ApplicationAdded:
		a.setRootData(event)
		a.SetData(event)
		a.CreationDate = event.CreationDate
	case es_model.OIDCConfigAdded:
		a.IsOIDC = true
		a.SetData(event)
	case es_model.OIDCConfigChanged,
		es_model.ApplicationChanged:
		a.SetData(event)
	case es_model.ApplicationDeactivated:
		a.State = int32(model.APPSTATE_INACTIVE)
	case es_model.ApplicationReactivated:
		a.State = int32(model.APPSTATE_ACTIVE)
	}
	return nil
}

func (a *ApplicationView) setRootData(event *models.Event) {
	a.ProjectID = event.AggregateID
	a.ChangeDate = event.CreationDate
}

func (a *ApplicationView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, a); err != nil {
		logging.Log("EVEN-lo9ds").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}
