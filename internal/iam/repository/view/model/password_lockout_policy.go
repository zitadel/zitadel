package model

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/domain"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"time"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	LockoutKeyAggregateID = "aggregate_id"
)

type LockoutPolicyView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:lockout_policy_state"`

	MaxPasswordAttempts uint64 `json:"maxPasswordAttempts" gorm:"column:max_password_attempts"`
	ShowLockOutFailures bool   `json:"showLockOutFailures" gorm:"column:show_lockout_failures"`
	Default             bool   `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func LockoutViewToModel(policy *LockoutPolicyView) *model.LockoutPolicyView {
	return &model.LockoutPolicyView{
		AggregateID:         policy.AggregateID,
		Sequence:            policy.Sequence,
		CreationDate:        policy.CreationDate,
		ChangeDate:          policy.ChangeDate,
		MaxPasswordAttempts: policy.MaxPasswordAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
		Default:             policy.Default,
	}
}

func (p *LockoutPolicyView) ToDomain() *domain.LockoutPolicy {
	return &domain.LockoutPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  p.AggregateID,
			CreationDate: p.CreationDate,
			ChangeDate:   p.ChangeDate,
			Sequence:     p.Sequence,
		},
		MaxPasswordAttempts: p.MaxPasswordAttempts,
		ShowLockOutFailures: p.ShowLockOutFailures,
		Default:             p.Default,
	}
}

func (i *LockoutPolicyView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.LockoutPolicyAdded, org_es_model.LockoutPolicyAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	case es_model.LockoutPolicyChanged, org_es_model.LockoutPolicyChanged:
		err = i.SetData(event)
	}
	return err
}

func (r *LockoutPolicyView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *LockoutPolicyView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-gHls0").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
