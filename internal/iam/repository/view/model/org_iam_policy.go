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
	OrgIAMPolicyKeyAggregateID = "aggregate_id"
)

type OrgIAMPolicyView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:org_iam_policy_state"`

	UserLoginMustBeDomain bool `json:"userLoginMustBeDomain" gorm:"column:user_login_must_be_domain"`
	Default               bool `json:"-" gorm:"-"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func OrgIAMViewFromModel(policy *model.OrgIAMPolicyView) *OrgIAMPolicyView {
	return &OrgIAMPolicyView{
		AggregateID:           policy.AggregateID,
		Sequence:              policy.Sequence,
		CreationDate:          policy.CreationDate,
		ChangeDate:            policy.ChangeDate,
		Default:               policy.Default,
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func OrgIAMViewToModel(policy *OrgIAMPolicyView) *model.OrgIAMPolicyView {
	return &model.OrgIAMPolicyView{
		AggregateID:           policy.AggregateID,
		Sequence:              policy.Sequence,
		CreationDate:          policy.CreationDate,
		ChangeDate:            policy.ChangeDate,
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		Default:               policy.Default,
	}
}

func (i *OrgIAMPolicyView) AppendEvent(event *models.Event) (err error) {
	i.Sequence = event.Sequence
	i.ChangeDate = event.CreationDate
	switch event.Type {
	case es_model.OrgIAMPolicyAdded, org_es_model.OrgIAMPolicyAdded:
		i.setRootData(event)
		i.CreationDate = event.CreationDate
		err = i.SetData(event)
	case es_model.OrgIAMPolicyChanged, org_es_model.OrgIAMPolicyChanged:
		err = i.SetData(event)
	}
	return err
}

func (r *OrgIAMPolicyView) setRootData(event *models.Event) {
	r.AggregateID = event.AggregateID
}

func (r *OrgIAMPolicyView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, r); err != nil {
		logging.Log("EVEN-Dmi9g").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Hs8uf", "Could not unmarshal data")
	}
	return nil
}
