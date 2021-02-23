package model

import (
	"encoding/json"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"time"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	PasswordLockoutKeyAggregateID = "aggregate_id"
)

type PasswordLockoutPolicyView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:lockout_policy_state"`

	MaxAttempts         uint64 `json:"maxAttempts" gorm:"column:max_attempts"`
	ShowLockOutFailures bool   `json:"showLockOutFailures" gorm:"column:show_lockout_failures"`
	Default             bool   `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func PasswordLockoutViewFromModel(policy *model.PasswordLockoutPolicyView) *PasswordLockoutPolicyView {
	return &PasswordLockoutPolicyView{
		AggregateID:         policy.AggregateID,
		Sequence:            policy.Sequence,
		CreationDate:        policy.CreationDate,
		ChangeDate:          policy.ChangeDate,
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
		Default:             policy.Default,
	}
}

func PasswordLockoutViewToModel(policy *PasswordLockoutPolicyView) *model.PasswordLockoutPolicyView {
	return &model.PasswordLockoutPolicyView{
		AggregateID:         policy.AggregateID,
		Sequence:            policy.Sequence,
		CreationDate:        policy.CreationDate,
		ChangeDate:          policy.ChangeDate,
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
		Default:             policy.Default,
	}
}

func (i *PasswordLockoutPolicyView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.PasswordLockoutPolicyAdded, org_es_model.PasswordLockoutPolicyAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	case es_model.PasswordLockoutPolicyChanged, org_es_model.PasswordLockoutPolicyChanged:
		err = i.SetData(event)
	}
	return err
}

func (r *PasswordLockoutPolicyView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *PasswordLockoutPolicyView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-gHls0").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
