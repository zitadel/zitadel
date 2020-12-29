package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type OrgIAMPolicy struct {
	models.ObjectRoot

	State                 int32 `json:"-"`
	UserLoginMustBeDomain bool  `json:"userLoginMustBeDomain"`
}

func OrgIAMPolicyFromModel(policy *iam_model.OrgIAMPolicy) *OrgIAMPolicy {
	return &OrgIAMPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 int32(policy.State),
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func OrgIAMPolicyToModel(policy *OrgIAMPolicy) *iam_model.OrgIAMPolicy {
	return &iam_model.OrgIAMPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 iam_model.PolicyState(policy.State),
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func (p *OrgIAMPolicy) Changes(changed *OrgIAMPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 1)

	if p.UserLoginMustBeDomain != changed.UserLoginMustBeDomain {
		changes["userLoginMustBeDomain"] = changed.UserLoginMustBeDomain
	}
	return changes
}

func (i *IAM) appendAddOrgIAMPolicyEvent(event *es_models.Event) error {
	i.DefaultOrgIAMPolicy = new(OrgIAMPolicy)
	err := i.DefaultOrgIAMPolicy.SetData(event)
	if err != nil {
		return err
	}
	i.DefaultOrgIAMPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (i *IAM) appendChangeOrgIAMPolicyEvent(event *es_models.Event) error {
	return i.DefaultOrgIAMPolicy.SetData(event)
}

func (p *OrgIAMPolicy) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}
