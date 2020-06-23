package model

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgIamPolicy struct {
	models.ObjectRoot

	Description           string `json:"description,omitempty"`
	State                 int32  `json:"-"`
	UserLoginMustBeDomain bool   `json:"userLoginMustBeDomain"`
}

func OrgIamPolicyToModel(policy *OrgIamPolicy) *org_model.OrgIamPolicy {
	return &org_model.OrgIamPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 org_model.PolicyState(policy.State),
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func OrgIamPolicyFromModel(policy *org_model.OrgIamPolicy) *OrgIamPolicy {
	return &OrgIamPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 int32(policy.State),
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func (o *Org) appendAddOrgIamPolicyEvent(event *es_models.Event) error {
	o.OrgIamPolicy = new(OrgIamPolicy)
	err := o.OrgIamPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.OrgIamPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeOrgIamPolicyEvent(event *es_models.Event) error {
	return o.OrgIamPolicy.SetData(event)
}

func (o *Org) appendRemoveOrgIamPolicyEvent() {
	o.OrgIamPolicy = nil
}

func (p *OrgIamPolicy) Changes(changed *OrgIamPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if changed.Description != p.Description {
		changes["description"] = changed.Description
	}
	if changed.UserLoginMustBeDomain != p.UserLoginMustBeDomain {
		changes["userLoginMustBeDomain"] = changed.UserLoginMustBeDomain
	}

	return changes
}

func (p *OrgIamPolicy) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}
