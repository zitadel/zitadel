package model

import (
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/crypto"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	"github.com/caos/logging"
	"github.com/lib/pq"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	IDPConfigKeyIdpConfigID  = "idp_config_id"
	IDPConfigKeyAggregateID  = "aggregate_id"
	IDPConfigKeyName         = "name"
	IDPConfigKeyProviderType = "idp_provider_type"
)

type IDPConfigView struct {
	IDPConfigID     string    `json:"idpConfigId" gorm:"column:idp_config_id;primary_key"`
	AggregateID     string    `json:"-" gorm:"column:aggregate_id"`
	Name            string    `json:"name" gorm:"column:name"`
	StylingType     int32     `json:"stylingType" gorm:"column:styling_type"`
	CreationDate    time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate      time.Time `json:"-" gorm:"column:change_date"`
	IDPState        int32     `json:"-" gorm:"column:idp_state"`
	IDPProviderType int32     `json:"-" gorm:"column:idp_provider_type"`

	IsOIDC                     bool                `json:"-" gorm:"column:is_oidc"`
	OIDCClientID               string              `json:"clientId" gorm:"column:oidc_client_id"`
	OIDCClientSecret           *crypto.CryptoValue `json:"clientSecret" gorm:"column:oidc_client_secret"`
	OIDCIssuer                 string              `json:"issuer" gorm:"column:oidc_issuer"`
	OIDCScopes                 pq.StringArray      `json:"scopes" gorm:"column:oidc_scopes"`
	OIDCIDPDisplayNameMapping  int32               `json:"idpDisplayNameMapping" gorm:"column:oidc_idp_display_name_mapping"`
	OIDCUsernameMapping        int32               `json:"usernameMapping" gorm:"column:oidc_idp_username_mapping"`
	OAuthAuthorizationEndpoint string              `json:"authorizationEndpoint" gorm:"column:oauth_authorization_endpoint"`
	OAuthTokenEndpoint         string              `json:"tokenEndpoint" gorm:"column:oauth_token_endpoint"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func IDPConfigViewToModel(idp *IDPConfigView) *model.IDPConfigView {
	return &model.IDPConfigView{
		IDPConfigID:                idp.IDPConfigID,
		AggregateID:                idp.AggregateID,
		State:                      model.IDPConfigState(idp.IDPState),
		Name:                       idp.Name,
		StylingType:                model.IDPStylingType(idp.StylingType),
		Sequence:                   idp.Sequence,
		CreationDate:               idp.CreationDate,
		ChangeDate:                 idp.ChangeDate,
		IDPProviderType:            model.IDPProviderType(idp.IDPProviderType),
		IsOIDC:                     idp.IsOIDC,
		OIDCClientID:               idp.OIDCClientID,
		OIDCClientSecret:           idp.OIDCClientSecret,
		OIDCIssuer:                 idp.OIDCIssuer,
		OIDCScopes:                 idp.OIDCScopes,
		OIDCIDPDisplayNameMapping:  model.OIDCMappingField(idp.OIDCIDPDisplayNameMapping),
		OIDCUsernameMapping:        model.OIDCMappingField(idp.OIDCUsernameMapping),
		OAuthAuthorizationEndpoint: idp.OAuthAuthorizationEndpoint,
		OAuthTokenEndpoint:         idp.OAuthTokenEndpoint,
	}
}

func IdpConfigViewsToModel(idps []*IDPConfigView) []*model.IDPConfigView {
	result := make([]*model.IDPConfigView, len(idps))
	for i, idp := range idps {
		result[i] = IDPConfigViewToModel(idp)
	}
	return result
}

func (i *IDPConfigView) AppendEvent(providerType model.IDPProviderType, event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.IDPConfigAdded, org_es_model.IDPConfigAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		i.IDPProviderType = int32(providerType)
		err = i.SetData(event)
	case es_model.OIDCIDPConfigAdded, org_es_model.OIDCIDPConfigAdded:
		i.IsOIDC = true
		err = i.SetData(event)
	case es_model.OIDCIDPConfigChanged, org_es_model.OIDCIDPConfigChanged,
		es_model.IDPConfigChanged, org_es_model.IDPConfigChanged:
		err = i.SetData(event)
	case es_model.IDPConfigDeactivated, org_es_model.IDPConfigDeactivated:
		i.IDPState = int32(model.IDPConfigStateInactive)
	case es_model.IDPConfigReactivated, org_es_model.IDPConfigReactivated:
		i.IDPState = int32(model.IDPConfigStateActive)
	}
	return err
}

func (r *IDPConfigView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *IDPConfigView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Smkld").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}
