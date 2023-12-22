package model

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DomainPolicy struct {
	es_models.ObjectRoot

	State                 int32 `json:"-"`
	UserLoginMustBeDomain bool  `json:"userLoginMustBeDomain"`
}

func DomainPolicyToModel(policy *DomainPolicy) *iam_model.DomainPolicy {
	return &iam_model.DomainPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 iam_model.PolicyState(policy.State),
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
	}
}

func (p *DomainPolicy) Changes(changed *DomainPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 1)

	if p.UserLoginMustBeDomain != changed.UserLoginMustBeDomain {
		changes["userLoginMustBeDomain"] = changed.UserLoginMustBeDomain
	}
	return changes
}

func (p *DomainPolicy) SetData(event eventstore.Event) error {
	err := event.Unmarshal(p)
	if err != nil {
		return zerrors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}
