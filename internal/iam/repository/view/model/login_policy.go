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
	LoginPolicyKeyAggregateID = "aggregate_id"
)

type LoginPolicyView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:login_policy_state"`

	AllowRegister         bool `json:"allowRegister" gorm:"column:allow_register"`
	AllowUsernamePassword bool `json:"allowUsernamePassword" gorm:"column:allow_username_password"`
	AllowExternalIdp      bool `json:"allowExternalIdp" gorm:"column:allow_external_idp"`
	Default               bool `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func LoginPolicyViewFromModel(policy *model.LoginPolicyView) *LoginPolicyView {
	return &LoginPolicyView{
		AggregateID:           policy.AggregateID,
		Sequence:              policy.Sequence,
		CreationDate:          policy.CreationDate,
		ChangeDate:            policy.ChangeDate,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowUsernamePassword: policy.AllowUsernamePassword,
		Default:               policy.Default,
	}
}

func LoginPolicyViewToModel(policy *LoginPolicyView) *model.LoginPolicyView {
	return &model.LoginPolicyView{
		AggregateID:           policy.AggregateID,
		Sequence:              policy.Sequence,
		CreationDate:          policy.CreationDate,
		ChangeDate:            policy.ChangeDate,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowUsernamePassword: policy.AllowUsernamePassword,
		Default:               policy.Default,
	}
}

func (i *LoginPolicyView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.LoginPolicyAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	case es_model.LoginPolicyChanged:
		err = i.SetData(event)
	}
	return err
}

func (r *LoginPolicyView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *LoginPolicyView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Kn7ds").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
