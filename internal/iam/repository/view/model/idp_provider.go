package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	IDPProviderKeyAggregateID = "aggregate_id"
	IDPProviderKeyIdpConfigID = "idp_config_id"
	IDPProviderKeyState       = "idp_state"
	IDPProviderKeyInstanceID  = "instance_id"
)

type IDPProviderView struct {
	AggregateID string `json:"-" gorm:"column:aggregate_id;primary_key"`
	IDPConfigID string `json:"idpConfigID" gorm:"column:idp_config_id;primary_key"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`

	Name            string `json:"-" gorm:"column:name"`
	StylingType     int32  `json:"-" gorm:"column:styling_type"`
	IDPConfigType   int32  `json:"-" gorm:"column:idp_config_type"`
	IDPProviderType int32  `json:"idpProviderType" gorm:"column:idp_provider_type"`
	IDPState        int32  `json:"-" gorm:"column:idp_state"`

	Sequence   uint64 `json:"-" gorm:"column:sequence"`
	InstanceID string `json:"instanceID" gorm:"column:instance_id;primary_key"`
}

func IDPProviderViewToModel(provider *IDPProviderView) *model.IDPProviderView {
	return &model.IDPProviderView{
		AggregateID:     provider.AggregateID,
		Sequence:        provider.Sequence,
		CreationDate:    provider.CreationDate,
		ChangeDate:      provider.ChangeDate,
		Name:            provider.Name,
		StylingType:     model.IDPStylingType(provider.StylingType),
		IDPConfigID:     provider.IDPConfigID,
		IDPConfigType:   model.IdpConfigType(provider.IDPConfigType),
		IDPProviderType: model.IDPProviderType(provider.IDPProviderType),
		IDPState:        model.IDPConfigState(provider.IDPState),
	}
}

func IDPProviderViewsToModel(providers []*IDPProviderView) []*model.IDPProviderView {
	result := make([]*model.IDPProviderView, len(providers))
	for i, r := range providers {
		result[i] = IDPProviderViewToModel(r)
	}
	return result
}

func (i *IDPProviderView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch eventstore.EventType(event.Type) {
	case instance.LoginPolicyIDPProviderAddedEventType,
		org.LoginPolicyIDPProviderAddedEventType:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	}
	return err
}

func (r *IDPProviderView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
	r.InstanceID = event.InstanceID
}

func (r *IDPProviderView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.New().WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
