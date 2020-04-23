package eventsourcing

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
)

type PasswordLockoutPolicy struct {
	models.ObjectRoot

	Description         string
	State               int32
	MaxAttempts         uint64
	ShowLockOutFailures bool
}

func (p *PasswordLockoutPolicy) LockoutChanges(changed *PasswordLockoutPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Description != "" && p.Description != changed.Description {
		changes["description"] = changed.Description
	}
	// todo
	return changes
}

func PasswordLockoutPolicyFromModel(policy *model.PasswordLockoutPolicy) *PasswordLockoutPolicy {
	return &PasswordLockoutPolicy{
		ObjectRoot:          policy.ObjectRoot,
		Description:         policy.Description,
		State:               policy.State,
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

func PasswordLockoutPolicyToModel(policy *PasswordLockoutPolicy) *model.PasswordLockoutPolicy {
	return &model.PasswordLockoutPolicy{
		ObjectRoot:          policy.ObjectRoot,
		Description:         policy.Description,
		State:               policy.State,
		MaxAttempts:         policy.MaxAttempts,
		ShowLockOutFailures: policy.ShowLockOutFailures,
	}
}

func (p *PasswordLockoutPolicy) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *PasswordLockoutPolicy) AppendEvent(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case model.PasswordLockoutPolicyAdded, model.PasswordLockoutPolicyChanged:
		if err := json.Unmarshal(event.Data, p); err != nil {
			logging.Log("EVEN-idl93").WithError(err).Error("could not unmarshal event data")
			return err
		}
		return nil
	}
	return nil
}
