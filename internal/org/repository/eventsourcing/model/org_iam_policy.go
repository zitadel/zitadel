package model

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgIAMPolicy struct {
	models.ObjectRoot

	Description           string `json:"description,omitempty"`
	State                 int32  `json:"-"`
	UserLoginMustBeDomain bool   `json:"userLoginMustBeDomain"`
}

func OrgIAMPolicyToModel(policy *OrgIAMPolicy) *org_model.OrgIAMPolicy {
	return &org_model.OrgIAMPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 org_model.PolicyState(policy.State),
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func OrgIAMPolicyFromModel(policy *org_model.OrgIAMPolicy) *OrgIAMPolicy {
	return &OrgIAMPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 int32(policy.State),
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func (o *Org) appendAddOrgIAMPolicyEvent(event *es_models.Event) error {
	o.OrgIamPolicy = new(OrgIAMPolicy)
	err := o.OrgIamPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.OrgIamPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeOrgIAMPolicyEvent(event *es_models.Event) error {
	return o.OrgIamPolicy.SetData(event)
}

func (o *Org) appendRemoveOrgIAMPolicyEvent() {
	o.OrgIamPolicy = nil
}

func (p *OrgIAMPolicy) Changes(changed *OrgIAMPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if changed.Description != p.Description {
		changes["description"] = changed.Description
	}
	if changed.UserLoginMustBeDomain != p.UserLoginMustBeDomain {
		changes["userLoginMustBeDomain"] = changed.UserLoginMustBeDomain
	}

	return changes
}

func (p *OrgIAMPolicy) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}
