package model

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"time"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/lib/pq"
)

const (
	IdpConfigKeyIdpConfigID  = "idp_config_id"
	IdpConfigKeyAggregateID  = "aggregate_id"
	IdpConfigKeyName         = "name"
	IdpConfigKeyProviderType = "idp_provider_type"
)

type IdpConfigView struct {
	IdpConfigID     string    `json:"idpConfigId" gorm:"column:idp_config_id;primary_key"`
	AggregateID     string    `json:"-" gorm:"column:aggregate_id"`
	Name            string    `json:"name" gorm:"column:name"`
	LogoSrc         string    `json:"logoSrc" gorm:"column:logo_src"`
	CreationDate    time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate      time.Time `json:"-" gorm:"column:change_date"`
	IdpState        int32     `json:"-" gorm:"column:idp_state"`
	IdpProviderType int32     `json:"-" gorm:"column:idp_provider_type"`

	IsOidc           bool                `json:"-" gorm:"column:is_oidc"`
	OidcClientID     string              `json:"clientId" gorm:"column:oidc_client_id"`
	OidcClientSecret *crypto.CryptoValue `json:"clientSecret" gorm:"column:oidc_client_secret"`
	OidcIssuer       string              `json:"issuer" gorm:"column:oidc_issuer"`
	OidcScopes       pq.StringArray      `json:"scopes" gorm:"column:oidc_scopes"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func IdpConfigViewFromModel(idp *model.IdpConfigView) *IdpConfigView {
	return &IdpConfigView{
		IdpConfigID:      idp.IdpConfigID,
		AggregateID:      idp.AggregateID,
		Name:             idp.Name,
		LogoSrc:          idp.LogoSrc,
		Sequence:         idp.Sequence,
		CreationDate:     idp.CreationDate,
		ChangeDate:       idp.ChangeDate,
		IdpProviderType:  int32(idp.IdpProviderType),
		IsOidc:           idp.IsOidc,
		OidcClientID:     idp.OidcClientID,
		OidcClientSecret: idp.OidcClientSecret,
		OidcIssuer:       idp.OidcIssuer,
		OidcScopes:       idp.OidcScopes,
	}
}

func IdpConfigViewToModel(idp *IdpConfigView) *model.IdpConfigView {
	return &model.IdpConfigView{
		IdpConfigID:      idp.IdpConfigID,
		AggregateID:      idp.AggregateID,
		Name:             idp.Name,
		LogoSrc:          idp.LogoSrc,
		Sequence:         idp.Sequence,
		CreationDate:     idp.CreationDate,
		ChangeDate:       idp.ChangeDate,
		IdpProviderType:  model.IdpProviderType(idp.IdpProviderType),
		IsOidc:           idp.IsOidc,
		OidcClientID:     idp.OidcClientID,
		OidcClientSecret: idp.OidcClientSecret,
		OidcIssuer:       idp.OidcIssuer,
		OidcScopes:       idp.OidcScopes,
	}
}

func IdpConfigViewsToModel(idps []*IdpConfigView) []*model.IdpConfigView {
	result := make([]*model.IdpConfigView, len(idps))
	for i, idp := range idps {
		result[i] = IdpConfigViewToModel(idp)
	}
	return result
}

func (i *IdpConfigView) AppendEvent(providerType model.IdpProviderType, event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.IdpConfigAdded, org_es_model.IdpConfigAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		i.IdpProviderType = int32(providerType)
		err = i.SetData(event)
	case es_model.OidcIdpConfigAdded, org_es_model.OidcIdpConfigAdded:
		i.IsOidc = true
		err = i.SetData(event)
	case es_model.OidcIdpConfigChanged, org_es_model.OidcIdpConfigChanged,
		es_model.IdpConfigChanged, org_es_model.IdpConfigChanged:
		err = i.SetData(event)
	case es_model.IdpConfigDeactivated, org_es_model.IdpConfigDeactivated:
		i.IdpState = int32(model.IdpConfigStateInactive)
	case es_model.IdpConfigReactivated, org_es_model.IdpConfigReactivated:
		i.IdpState = int32(model.IdpConfigStateActive)
	}
	return err
}

func (r *IdpConfigView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *IdpConfigView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Smkld").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-lub6s", "Could not unmarshal data")
	}
	return nil
}
