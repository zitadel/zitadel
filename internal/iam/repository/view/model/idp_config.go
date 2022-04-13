package model

import (
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"

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
	IDPConfigKeyInstanceID   = "instance_id"
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
	AutoRegister    bool      `json:"autoRegister" gorm:"column:auto_register"`

	IsOIDC                     bool                `json:"-" gorm:"column:is_oidc"`
	OIDCClientID               string              `json:"clientId" gorm:"column:oidc_client_id"`
	OIDCClientSecret           *crypto.CryptoValue `json:"clientSecret" gorm:"column:oidc_client_secret"`
	OIDCIssuer                 string              `json:"issuer" gorm:"column:oidc_issuer"`
	OIDCScopes                 pq.StringArray      `json:"scopes" gorm:"column:oidc_scopes"`
	OIDCIDPDisplayNameMapping  int32               `json:"idpDisplayNameMapping" gorm:"column:oidc_idp_display_name_mapping"`
	OIDCUsernameMapping        int32               `json:"usernameMapping" gorm:"column:oidc_idp_username_mapping"`
	OAuthAuthorizationEndpoint string              `json:"authorizationEndpoint" gorm:"column:oauth_authorization_endpoint"`
	OAuthTokenEndpoint         string              `json:"tokenEndpoint" gorm:"column:oauth_token_endpoint"`
	JWTEndpoint                string              `json:"jwtEndpoint" gorm:"jwt_endpoint"`
	JWTKeysEndpoint            string              `json:"keysEndpoint" gorm:"jwt_keys_endpoint"`
	JWTHeaderName              string              `json:"headerName" gorm:"jwt_header_name"`

	Sequence   uint64 `json:"-" gorm:"column:sequence"`
	InstanceID string `json:"instanceID" gorm:"column:instance_id;primary_key"`
}

func IDPConfigViewToModel(idp *IDPConfigView) *model.IDPConfigView {
	view := &model.IDPConfigView{
		IDPConfigID:                idp.IDPConfigID,
		AggregateID:                idp.AggregateID,
		State:                      model.IDPConfigState(idp.IDPState),
		Name:                       idp.Name,
		StylingType:                model.IDPStylingType(idp.StylingType),
		AutoRegister:               idp.AutoRegister,
		Sequence:                   idp.Sequence,
		CreationDate:               idp.CreationDate,
		ChangeDate:                 idp.ChangeDate,
		IDPProviderType:            model.IDPProviderType(idp.IDPProviderType),
		IsOIDC:                     idp.IsOIDC,
		OIDCClientID:               idp.OIDCClientID,
		OIDCClientSecret:           idp.OIDCClientSecret,
		OIDCScopes:                 idp.OIDCScopes,
		OIDCIDPDisplayNameMapping:  model.OIDCMappingField(idp.OIDCIDPDisplayNameMapping),
		OIDCUsernameMapping:        model.OIDCMappingField(idp.OIDCUsernameMapping),
		OAuthAuthorizationEndpoint: idp.OAuthAuthorizationEndpoint,
		OAuthTokenEndpoint:         idp.OAuthTokenEndpoint,
	}
	if idp.IsOIDC {
		view.OIDCIssuer = idp.OIDCIssuer
		return view
	}
	view.JWTEndpoint = idp.JWTEndpoint
	view.JWTIssuer = idp.OIDCIssuer
	view.JWTKeysEndpoint = idp.JWTKeysEndpoint
	view.JWTHeaderName = idp.JWTHeaderName
	return view
}

func (i *IDPConfigView) AppendEvent(providerType model.IDPProviderType, event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch eventstore.EventType(event.Type) {
	case instance.IDPConfigAddedEventType, org.IDPConfigAddedEventType:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		i.IDPProviderType = int32(providerType)
		err = i.SetData(event)
	case instance.IDPOIDCConfigAddedEventType, org.IDPOIDCConfigAddedEventType:
		i.IsOIDC = true
		err = i.SetData(event)
	case instance.IDPOIDCConfigChangedEventType, org.IDPOIDCConfigChangedEventType,
		instance.IDPConfigChangedEventType, org.IDPConfigChangedEventType,
		org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType,
		org.IDPJWTConfigChangedEventType, instance.IDPJWTConfigChangedEventType:
		err = i.SetData(event)
	case instance.IDPConfigDeactivatedEventType, org.IDPConfigDeactivatedEventType:
		i.IDPState = int32(model.IDPConfigStateInactive)
	case instance.IDPConfigReactivatedEventType, org.IDPConfigReactivatedEventType:
		i.IDPState = int32(model.IDPConfigStateActive)
	}
	return err
}

func (r *IDPConfigView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
	r.InstanceID = event.InstanceID
}

func (r *IDPConfigView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.New().WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}
