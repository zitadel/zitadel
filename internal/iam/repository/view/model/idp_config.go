package model

import (
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	org_repo "github.com/caos/zitadel/internal/repository/org"

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
	Sequence        uint64    `json:"-" gorm:"column:sequence"`

	*IDPConfigOIDCView
	*IDPConfigAuthConnectorView
}

type IDPConfigOIDCView struct {
	OIDCClientID              string              `json:"clientId" gorm:"column:oidc_client_id"`
	OIDCClientSecret          *crypto.CryptoValue `json:"clientSecret" gorm:"column:oidc_client_secret"`
	OIDCIssuer                string              `json:"issuer" gorm:"column:oidc_issuer"`
	OIDCScopes                pq.StringArray      `json:"scopes" gorm:"column:oidc_scopes"`
	OIDCIDPDisplayNameMapping int32               `json:"idpDisplayNameMapping" gorm:"column:oidc_idp_display_name_mapping"`
	OIDCUsernameMapping       int32               `json:"usernameMapping" gorm:"column:oidc_idp_username_mapping"`
}

func (v *IDPConfigOIDCView) IsZero() bool {
	return v == nil || v.OIDCIssuer == ""
}

type IDPConfigAuthConnectorView struct {
	AuthConnectorBaseURL     string `json:"baseUrl" gorm:"column:auth_connector_base_url"`
	AuthConnectorProviderID  string `json:"providerId" gorm:"column:auth_connector_provider_id"`
	AuthConnectorMachineID   string `json:"machineId" gorm:"column:auth_connector_machine_id"`
	AuthConnectorMachineName string `json:"-" gorm:"column:auth_connector_machine_name"`
}

func (v *IDPConfigAuthConnectorView) IsZero() bool {
	return v == nil || v.AuthConnectorBaseURL == ""
}

func IDPConfigViewToModel(idp *IDPConfigView) *model.IDPConfigView {
	idpView := &model.IDPConfigView{
		IDPConfigID:     idp.IDPConfigID,
		AggregateID:     idp.AggregateID,
		State:           model.IDPConfigState(idp.IDPState),
		Name:            idp.Name,
		StylingType:     model.IDPStylingType(idp.StylingType),
		Sequence:        idp.Sequence,
		CreationDate:    idp.CreationDate,
		ChangeDate:      idp.ChangeDate,
		IDPProviderType: model.IDPProviderType(idp.IDPProviderType),
	}
	if !idp.IDPConfigOIDCView.IsZero() {
		idpView.IDPConfigOIDCView = &model.IDPConfigOIDCView{
			OIDCClientID:              idp.OIDCClientID,
			OIDCClientSecret:          idp.OIDCClientSecret,
			OIDCIssuer:                idp.OIDCIssuer,
			OIDCScopes:                idp.OIDCScopes,
			OIDCIDPDisplayNameMapping: model.OIDCMappingField(idp.OIDCIDPDisplayNameMapping),
			OIDCUsernameMapping:       model.OIDCMappingField(idp.OIDCUsernameMapping),
		}
	}
	if !idp.IDPConfigAuthConnectorView.IsZero() {
		idpView.IDPConfigAuthConnectorView = &model.IDPConfigAuthConnectorView{
			AuthConnectorBaseURL:     idp.AuthConnectorBaseURL,
			AuthConnectorProviderID:  idp.AuthConnectorProviderID,
			AuthConnectorMachineID:   idp.AuthConnectorMachineID,
			AuthConnectorMachineName: idp.AuthConnectorMachineName,
		}
	}
	return idpView
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
	case es_model.OIDCIDPConfigAdded, org_es_model.OIDCIDPConfigAdded,
		models.EventType(iam_repo.IDPAuthConnectorConfigAddedEventType),
		models.EventType(org_repo.IDPAuthConnectorConfigAddedEventType),
		es_model.OIDCIDPConfigChanged, org_es_model.OIDCIDPConfigChanged,
		es_model.IDPConfigChanged, org_es_model.IDPConfigChanged,
		models.EventType(iam_repo.IDPAuthConnectorConfigChangedEventType),
		models.EventType(org_repo.IDPAuthConnectorConfigChangedEventType):
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
