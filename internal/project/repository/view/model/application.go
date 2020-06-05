package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
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
	ApplicationKeyName          = "app_name"
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
	OIDCPostLogoutRedirectUris pq.StringArray `json:"postLogoutRedirectUris" gorm:"column:oidc_post_logout_redirect_uris"`

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
	result := make([]int64, len(oidctypes))
	for i, t := range oidctypes {
		result[i] = int64(t)
	}
	return result
}

func OIDCGrantTypesFromModel(granttypes []model.OIDCGrantType) []int64 {
	result := make([]int64, len(granttypes))
	for i, t := range granttypes {
		result[i] = int64(t)
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
	result := make([]model.OIDCResponseType, len(oidctypes))
	for i, t := range oidctypes {
		result[i] = model.OIDCResponseType(t)
	}
	return result
}

func OIDCGrantTypesToModel(granttypes []int64) []model.OIDCGrantType {
	result := make([]model.OIDCGrantType, len(granttypes))
	for i, t := range granttypes {
		result[i] = model.OIDCGrantType(t)
	}
	return result
}

func ApplicationViewsToModel(roles []*ApplicationView) []*model.ApplicationView {
	result := make([]*model.ApplicationView, len(roles))
	for i, r := range roles {
		result[i] = ApplicationViewToModel(r)
	}
	return result
}

func (a *ApplicationView) AppendEvent(event *models.Event) (err error) {
	a.Sequence = event.Sequence
	a.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.ApplicationAdded:
		a.setRootData(event)
		a.CreationDate = event.CreationDate
		err = a.SetData(event)
	case es_model.OIDCConfigAdded:
		a.IsOIDC = true
		err = a.SetData(event)
	case es_model.OIDCConfigChanged,
		es_model.ApplicationChanged:
		err = a.SetData(event)
	case es_model.ApplicationDeactivated:
		a.State = int32(model.APPSTATE_INACTIVE)
	case es_model.ApplicationReactivated:
		a.State = int32(model.APPSTATE_ACTIVE)
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
