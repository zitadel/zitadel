package model

import (
	"encoding/json"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"time"

	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	PasswordAgeKeyAggregateID = "aggregate_id"
)

type PasswordAgePolicyView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:age_policy_state"`

	MaxAgeDays     uint64 `json:"maxAgeDays" gorm:"column:max_age_days"`
	ExpireWarnDays uint64 `json:"expireWarnDays" gorm:"column:expire_warn_days"`
	Default        bool   `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func PasswordAgeViewFromModel(policy *model.PasswordAgePolicyView) *PasswordAgePolicyView {
	return &PasswordAgePolicyView{
		AggregateID:    policy.AggregateID,
		Sequence:       policy.Sequence,
		CreationDate:   policy.CreationDate,
		ChangeDate:     policy.ChangeDate,
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		Default:        policy.Default,
	}
}

func PasswordAgeViewToModel(policy *PasswordAgePolicyView) *model.PasswordAgePolicyView {
	return &model.PasswordAgePolicyView{
		AggregateID:    policy.AggregateID,
		Sequence:       policy.Sequence,
		CreationDate:   policy.CreationDate,
		ChangeDate:     policy.ChangeDate,
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		Default:        policy.Default,
	}
}

func (i *PasswordAgePolicyView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case iam_es_model.PasswordAgePolicyAdded, org_es_model.PasswordAgePolicyAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	case iam_es_model.PasswordAgePolicyChanged, org_es_model.PasswordAgePolicyChanged:
		err = i.SetData(event)
	}
	return err
}

func (r *PasswordAgePolicyView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *PasswordAgePolicyView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-gH9os").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
