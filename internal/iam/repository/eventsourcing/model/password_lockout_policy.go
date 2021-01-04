package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type PasswordLockoutPolicy struct {
	models.ObjectRoot

	State               int32  `json:"-"`
	MaxAttempts         uint64 `json:"maxAttempts"`
	ShowLockOutFailures bool   `json:"showLockOutFailures"`
}

func PasswordLockoutPolicyFromModel(policy *iam_model.PasswordLockoutPolicy) *PasswordLockoutPolicy {
	return &PasswordLockoutPolicy{
		ObjectRoot:          policy.ObjectRoot,
		State:               int32(policy.State),
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

func PasswordLockoutPolicyToModel(policy *PasswordLockoutPolicy) *iam_model.PasswordLockoutPolicy {
	return &iam_model.PasswordLockoutPolicy{
		ObjectRoot:          policy.ObjectRoot,
		State:               iam_model.PolicyState(policy.State),
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

func (p *PasswordLockoutPolicy) Changes(changed *PasswordLockoutPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if p.MaxAttempts != changed.MaxAttempts {
		changes["maxAttempts"] = changed.MaxAttempts
	}
	if p.ShowLockOutFailures != changed.ShowLockOutFailures {
		changes["showLockOutFailures"] = changed.ShowLockOutFailures
	}
	return changes
}

func (i *IAM) appendAddPasswordLockoutPolicyEvent(event *es_models.Event) error {
	i.DefaultPasswordLockoutPolicy = new(PasswordLockoutPolicy)
	err := i.DefaultPasswordLockoutPolicy.SetData(event)
	if err != nil {
		return err
	}
	i.DefaultPasswordLockoutPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (i *IAM) appendChangePasswordLockoutPolicyEvent(event *es_models.Event) error {
	return i.DefaultPasswordLockoutPolicy.SetData(event)
}

func (p *PasswordLockoutPolicy) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}
