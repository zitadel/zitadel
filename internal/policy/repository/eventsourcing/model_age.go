package eventsourcing

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
)

const (
	policyAgeVersion = "v1"
)

type PasswordAgePolicy struct {
	models.ObjectRoot

	Description    string `json:"description,omitempty"`
	State          int32  `json:"-"`
	MaxAgeDays     uint64 `json:"maxAgeDays"`
	ExpireWarnDays uint64 `json:"expireWarnDays"`
}

func (p *PasswordAgePolicy) AgeChanges(changed *PasswordAgePolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Description != "" && p.Description != changed.Description {
		changes["description"] = changed.Description
	}
	if p.MaxAgeDays != changed.MaxAgeDays {
		changes["maxAgeDays"] = changed.MaxAgeDays
	}
	if p.ExpireWarnDays != changed.ExpireWarnDays {
		changes["expireWarnDays"] = changed.ExpireWarnDays
	}
	return changes
}

func PasswordAgePolicyFromModel(policy *model.PasswordAgePolicy) *PasswordAgePolicy {
	return &PasswordAgePolicy{
		ObjectRoot:     policy.ObjectRoot,
		Description:    policy.Description,
		State:          int32(policy.State),
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}

func PasswordAgePolicyToModel(policy *PasswordAgePolicy) *model.PasswordAgePolicy {
	return &model.PasswordAgePolicy{
		ObjectRoot:     policy.ObjectRoot,
		Description:    policy.Description,
		State:          model.PolicyState(policy.State),
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}

func (p *PasswordAgePolicy) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *PasswordAgePolicy) AppendEvent(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case model.PasswordAgePolicyAdded, model.PasswordAgePolicyChanged:
		if err := json.Unmarshal(event.Data, p); err != nil {
			logging.Log("EVEN-idl93").WithError(err).Error("could not unmarshal event data")
			return err
		}
		return nil
	}
	return nil
}
