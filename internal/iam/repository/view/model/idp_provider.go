package model

import (
	"encoding/json"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"time"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	IDPProviderKeyAggregateID = "aggregate_id"
	IDPProviderKeyIdpConfigID = "idp_config_id"
)

type IDPProviderView struct {
	AggregateID string `json:"-" gorm:"column:aggregate_id;primary_key"`
	IDPConfigID string `json:"idpConfigID" gorm:"column:idp_config_id;primary_key"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`

	Name            string `json:"-" gorm:"column:name"`
	IDPConfigType   int32  `json:"-" gorm:"column:idp_config_type"`
	IDPProviderType int32  `json:"idpProviderType" gorm:"column:idp_provider_type"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func IDPProviderViewFromModel(policy *model.IDPProviderView) *IDPProviderView {
	return &IDPProviderView{
		AggregateID:     policy.AggregateID,
		Sequence:        policy.Sequence,
		CreationDate:    policy.CreationDate,
		ChangeDate:      policy.ChangeDate,
		Name:            policy.Name,
		IDPConfigType:   int32(policy.IDPConfigType),
		IDPProviderType: int32(policy.IDPProviderType),
	}
}

func IDPProviderViewToModel(policy *IDPProviderView) *model.IDPProviderView {
	return &model.IDPProviderView{
		AggregateID:     policy.AggregateID,
		Sequence:        policy.Sequence,
		CreationDate:    policy.CreationDate,
		ChangeDate:      policy.ChangeDate,
		Name:            policy.Name,
		IDPConfigType:   model.IdpConfigType(policy.IDPConfigType),
		IDPProviderType: model.IDPProviderType(policy.IDPProviderType),
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
	switch event.Type {
	case es_model.LoginPolicyIDPProviderAdded, org_es_model.LoginPolicyIDPProviderAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	}
	return err
}

func (r *IDPProviderView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *IDPProviderView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Lso0d").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
