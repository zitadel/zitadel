package model

import (
	"encoding/json"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/lib/pq"
	"time"

	es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	LoginPolicyKeyAggregateID = "aggregate_id"
	LoginPolicyKeyDefault     = "default_policy"
)

type LoginPolicyView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:login_policy_state"`

	AllowRegister         bool          `json:"allowRegister" gorm:"column:allow_register"`
	AllowUsernamePassword bool          `json:"allowUsernamePassword" gorm:"column:allow_username_password"`
	AllowExternalIDP      bool          `json:"allowExternalIdp" gorm:"column:allow_external_idp"`
	ForceMFA              bool          `json:"forceMFA" gorm:"column:force_mfa"`
	PasswordlessType      int32         `json:"passwordlessType" gorm:"column:passwordless_type"`
	SecondFactors         pq.Int64Array `json:"-" gorm:"column:second_factors"`
	MultiFactors          pq.Int64Array `json:"-" gorm:"column:multi_factors"`
	Default               bool          `json:"-" gorm:"column:default_policy"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func LoginPolicyViewFromModel(policy *model.LoginPolicyView) *LoginPolicyView {
	return &LoginPolicyView{
		AggregateID:           policy.AggregateID,
		Sequence:              policy.Sequence,
		CreationDate:          policy.CreationDate,
		ChangeDate:            policy.ChangeDate,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIDP:      policy.AllowExternalIDP,
		AllowUsernamePassword: policy.AllowUsernamePassword,
		ForceMFA:              policy.ForceMFA,
		PasswordlessType:      int32(policy.PasswordlessType),
		SecondFactors:         secondFactorsFromModel(policy.SecondFactors),
		MultiFactors:          multiFactorsFromModel(policy.MultiFactors),
		Default:               policy.Default,
	}
}

func secondFactorsFromModel(mfas []model.SecondFactorType) []int64 {
	convertedMFAs := make([]int64, len(mfas))
	for i, m := range mfas {
		convertedMFAs[i] = int64(m)
	}
	return convertedMFAs
}

func multiFactorsFromModel(mfas []model.MultiFactorType) []int64 {
	convertedMFAs := make([]int64, len(mfas))
	for i, m := range mfas {
		convertedMFAs[i] = int64(m)
	}
	return convertedMFAs
}

func LoginPolicyViewToModel(policy *LoginPolicyView) *model.LoginPolicyView {
	return &model.LoginPolicyView{
		AggregateID:           policy.AggregateID,
		Sequence:              policy.Sequence,
		CreationDate:          policy.CreationDate,
		ChangeDate:            policy.ChangeDate,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIDP:      policy.AllowExternalIDP,
		AllowUsernamePassword: policy.AllowUsernamePassword,
		ForceMFA:              policy.ForceMFA,
		PasswordlessType:      model.PasswordlessType(policy.PasswordlessType),
		SecondFactors:         secondFactorsToModel(policy.SecondFactors),
		MultiFactors:          multiFactorsToToModel(policy.MultiFactors),
		Default:               policy.Default,
	}
}

func secondFactorsToModel(mfas []int64) []model.SecondFactorType {
	convertedMFAs := make([]model.SecondFactorType, len(mfas))
	for i, m := range mfas {
		convertedMFAs[i] = model.SecondFactorType(m)
	}
	return convertedMFAs
}

func multiFactorsToToModel(mfas []int64) []model.MultiFactorType {
	convertedMFAs := make([]model.MultiFactorType, len(mfas))
	for i, m := range mfas {
		convertedMFAs[i] = model.MultiFactorType(m)
	}
	return convertedMFAs
}

func (p *LoginPolicyView) AppendEvent(event *models.Event) (err error) {
	p.Sequence = event.Sequence
	p.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.LoginPolicyAdded:
		p.setRootData(event)
		p.CreationDate = event.CreationDate
		p.Default = true
		err = p.SetData(event)
	case org_es_model.LoginPolicyAdded:
		p.setRootData(event)
		p.CreationDate = event.CreationDate
		err = p.SetData(event)
		p.Default = false
	case es_model.LoginPolicyChanged, org_es_model.LoginPolicyChanged:
		err = p.SetData(event)
	case es_model.LoginPolicySecondFactorAdded, org_es_model.LoginPolicySecondFactorAdded:
		mfa := new(es_model.MFA)
		err := mfa.SetData(event)
		if err != nil {
			return err
		}
		if !existsMFA(p.SecondFactors, int64(mfa.MFAType)) {
			p.SecondFactors = append(p.SecondFactors, int64(mfa.MFAType))
		}

	case es_model.LoginPolicySecondFactorRemoved, org_es_model.LoginPolicySecondFactorRemoved:
		err = p.removeSecondFactor(event)
	case es_model.LoginPolicyMultiFactorAdded, org_es_model.LoginPolicyMultiFactorAdded:
		mfa := new(es_model.MFA)
		err := mfa.SetData(event)
		if err != nil {
			return err
		}
		if !existsMFA(p.MultiFactors, int64(mfa.MFAType)) {
			p.MultiFactors = append(p.MultiFactors, int64(mfa.MFAType))
		}
	case es_model.LoginPolicyMultiFactorRemoved, org_es_model.LoginPolicyMultiFactorRemoved:
		err = p.removeMultiFactor(event)
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

func (p *LoginPolicyView) removeSecondFactor(event *models.Event) error {
	mfa := new(es_model.MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	for i := len(p.SecondFactors) - 1; i >= 0; i-- {
		if p.SecondFactors[i] == int64(mfa.MFAType) {
			copy(p.SecondFactors[i:], p.SecondFactors[i+1:])
			p.SecondFactors[len(p.SecondFactors)-1] = 0
			p.SecondFactors = p.SecondFactors[:len(p.SecondFactors)-1]
			return nil
		}
	}
	return nil
}

func (p *LoginPolicyView) removeMultiFactor(event *models.Event) error {
	mfa := new(es_model.MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	for i := len(p.MultiFactors) - 1; i >= 0; i-- {
		if p.MultiFactors[i] == int64(mfa.MFAType) {
			copy(p.MultiFactors[i:], p.MultiFactors[i+1:])
			p.MultiFactors[len(p.MultiFactors)-1] = 0
			p.MultiFactors = p.MultiFactors[:len(p.MultiFactors)-1]
			return nil
		}
	}
	return nil
}

func existsMFA(mfas []int64, mfaType int64) bool {
	for _, m := range mfas {
		if m == mfaType {
			return true
		}
	}
	return false
}
