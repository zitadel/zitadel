package model

import (
	"encoding/json"
	"time"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	IdpProviderKeyAggregateID = "aggregate_id"
)

type IdpProviderView struct {
	AggregateID string `json:"-" gorm:"column:aggregate_id;primary_key"`
	IdpConfigID string `json:"idpConfigID" gorm:"column:idp_config_id;primary_key"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`

	Name            string `json:"-" gorm:"column:name"`
	IdpConfigType   int32  `json:"-" gorm:"column:idp_config_type"`
	IdpProviderType int32  `json:"idpProviderType" gorm:"column:idp_provider_type"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func IdpProviderViewFromModel(policy *model.IdpProviderView) *IdpProviderView {
	return &IdpProviderView{
		AggregateID:     policy.AggregateID,
		Sequence:        policy.Sequence,
		CreationDate:    policy.CreationDate,
		ChangeDate:      policy.ChangeDate,
		Name:            policy.Name,
		IdpConfigType:   int32(policy.IdpConfigType),
		IdpProviderType: int32(policy.IdpProviderType),
	}
}

func IdpProviderViewToModel(policy *IdpProviderView) *model.IdpProviderView {
	return &model.IdpProviderView{
		AggregateID:     policy.AggregateID,
		Sequence:        policy.Sequence,
		CreationDate:    policy.CreationDate,
		ChangeDate:      policy.ChangeDate,
		Name:            policy.Name,
		IdpConfigType:   model.IdpConfigType(policy.IdpConfigType),
		IdpProviderType: model.IdpProviderType(policy.IdpProviderType),
	}
}

func IdpProviderViewsToModel(providers []*IdpProviderView) []*model.IdpProviderView {
	result := make([]*model.IdpProviderView, len(providers))
	for i, r := range providers {
		result[i] = IdpProviderViewToModel(r)
	}
	return result
}

func (i *IdpProviderView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.LoginPolicyIdpProviderAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	}
	return err
}

func (r *IdpProviderView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *IdpProviderView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Lso0d").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
