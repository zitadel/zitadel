package model

import (
	"encoding/json"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/lib/pq"
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

	AllowRegister         bool          `json:"allowRegister" gorm:"column:allow_register"`
	AllowUsernamePassword bool          `json:"allowUsernamePassword" gorm:"column:allow_username_password"`
	AllowExternalIDP      bool          `json:"allowExternalIdp" gorm:"column:allow_external_idp"`
	ForceMFA              bool          `json:"forceMFA" gorm:"column:force_mfa"`
	SoftwareMFAs          pq.Int64Array `json:"-" gorm:"column:software_mfas"`
	HardwareMFAs          pq.Int64Array `json:"-" gorm:"column:hardware_mfas"`
	Default               bool          `json:"-" gorm:"-"`

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
		SoftwareMFAs:          softwareMFAsFromModel(policy.SoftwareMFAs),
		HardwareMFAs:          hardwareMFAsFromModel(policy.HardwareMFAs),
		Default:               policy.Default,
	}
}

func softwareMFAsFromModel(mfas []model.SoftwareMFAType) []int64 {
	convertedMFAs := make([]int64, len(mfas))
	for i, m := range mfas {
		convertedMFAs[i] = int64(m)
	}
	return convertedMFAs
}

func hardwareMFAsFromModel(mfas []model.HardwareMFAType) []int64 {
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
		SoftwareMFAs:          softwareMFAsToModel(policy.SoftwareMFAs),
		HardwareMFAs:          hardwareMFAsToToModel(policy.HardwareMFAs),
		Default:               policy.Default,
	}
}

func softwareMFAsToModel(mfas []int64) []model.SoftwareMFAType {
	convertedMFAs := make([]model.SoftwareMFAType, len(mfas))
	for i, m := range mfas {
		convertedMFAs[i] = model.SoftwareMFAType(m)
	}
	return convertedMFAs
}

func hardwareMFAsToToModel(mfas []int64) []model.HardwareMFAType {
	convertedMFAs := make([]model.HardwareMFAType, len(mfas))
	for i, m := range mfas {
		convertedMFAs[i] = model.HardwareMFAType(m)
	}
	return convertedMFAs
}

func (p *LoginPolicyView) AppendEvent(event *models.Event) (err error) {
	p.Sequence = event.Sequence
	p.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.LoginPolicyAdded, org_es_model.LoginPolicyAdded:
		p.setRootData(event)
		p.CreationDate = event.CreationDate
		err = p.SetData(event)
	case es_model.LoginPolicyChanged, org_es_model.LoginPolicyChanged:
		err = p.SetData(event)
	case es_model.LoginPolicySoftwareMFAAdded, org_es_model.LoginPolicySoftwareMFAAdded:
		mfa := new(es_model.MFA)
		err := mfa.SetData(event)
		if err != nil {
			return err
		}
		p.SoftwareMFAs = append(p.SoftwareMFAs, int64(mfa.MfaType))
	case es_model.LoginPolicySoftwareMFARemoved, org_es_model.LoginPolicySoftwareMFARemoved:
		err = p.removeSoftwareMFA(event)
	case es_model.LoginPolicyHardwareMFAAdded, org_es_model.LoginPolicyHardwareMFAAdded:
		mfa := new(es_model.MFA)
		err := mfa.SetData(event)
		if err != nil {
			return err
		}
		p.HardwareMFAs = append(p.HardwareMFAs, int64(mfa.MfaType))
	case es_model.LoginPolicyHardwareMFARemoved, org_es_model.LoginPolicyHardwareMFARemoved:
		err = p.removeHardwareMFA(event)
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

func (p *LoginPolicyView) removeSoftwareMFA(event *models.Event) error {
	mfa := new(es_model.MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	for i := len(p.SoftwareMFAs) - 1; i >= 0; i-- {
		if p.SoftwareMFAs[i] == int64(mfa.MfaType) {
			copy(p.SoftwareMFAs[i:], p.SoftwareMFAs[i+1:])
			p.SoftwareMFAs[len(p.SoftwareMFAs)-1] = 0
			p.SoftwareMFAs = p.SoftwareMFAs[:len(p.SoftwareMFAs)-1]
			return nil
		}
	}
	return nil
}

func (p *LoginPolicyView) removeHardwareMFA(event *models.Event) error {
	mfa := new(es_model.MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	for i := len(p.HardwareMFAs) - 1; i >= 0; i-- {
		if p.HardwareMFAs[i] == int64(mfa.MfaType) {
			copy(p.HardwareMFAs[i:], p.HardwareMFAs[i+1:])
			p.HardwareMFAs[len(p.HardwareMFAs)-1] = 0
			p.HardwareMFAs = p.HardwareMFAs[:len(p.HardwareMFAs)-1]
			return nil
		}
	}
	return nil
}
