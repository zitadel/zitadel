package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type LockoutPolicy struct {
	es_models.ObjectRoot

	State               int32  `json:"-"`
	MaxPasswordAttempts uint64 `json:"maxPasswordAttempts"`
	ShowLockOutFailures bool   `json:"showLockOutFailures"`
}

func LockoutPolicyToModel(policy *LockoutPolicy) *iam_model.LockoutPolicy {
	return &iam_model.LockoutPolicy{
		ObjectRoot:          policy.ObjectRoot,
		State:               iam_model.PolicyState(policy.State),
		MaxPasswordAttempts: policy.MaxPasswordAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

func (p *LockoutPolicy) Changes(changed *LockoutPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if p.MaxPasswordAttempts != changed.MaxPasswordAttempts {
		changes["maxAttempts"] = changed.MaxPasswordAttempts
	}
	if p.ShowLockOutFailures != changed.ShowLockOutFailures {
		changes["showLockOutFailures"] = changed.ShowLockOutFailures
	}
	return changes
}

func (i *IAM) appendAddLockoutPolicyEvent(event *es_models.Event) error {
	i.DefaultLockoutPolicy = new(LockoutPolicy)
	err := i.DefaultLockoutPolicy.SetData(event)
	if err != nil {
		return err
	}
	i.DefaultLockoutPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (i *IAM) appendChangeLockoutPolicyEvent(event *es_models.Event) error {
	return i.DefaultLockoutPolicy.SetData(event)
}

func (p *LockoutPolicy) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}
