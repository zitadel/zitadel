package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type PasswordAgePolicy struct {
	es_models.ObjectRoot

	State          int32  `json:"-"`
	MaxAgeDays     uint64 `json:"maxAgeDays"`
	ExpireWarnDays uint64 `json:"expireWarnDays"`
}

func PasswordAgePolicyToModel(policy *PasswordAgePolicy) *iam_model.PasswordAgePolicy {
	return &iam_model.PasswordAgePolicy{
		ObjectRoot:     policy.ObjectRoot,
		State:          iam_model.PolicyState(policy.State),
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}

func (p *PasswordAgePolicy) Changes(changed *PasswordAgePolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 1)

	if p.MaxAgeDays != changed.MaxAgeDays {
		changes["maxAgeDays"] = changed.MaxAgeDays
	}
	if p.ExpireWarnDays != changed.ExpireWarnDays {
		changes["expireWarnDays"] = changed.ExpireWarnDays
	}
	return changes
}

func (p *PasswordAgePolicy) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}
