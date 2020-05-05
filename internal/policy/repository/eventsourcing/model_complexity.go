package eventsourcing

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
)

const (
	policyComplexityVersion = "v1"
)

type PasswordComplexityPolicy struct {
	models.ObjectRoot

	Description  string `json:"description,omitempty"`
	State        int32  `json:"-"`
	MinLength    uint64 `json:"minLength"`
	HasLowercase bool   `json:"hasLowercase"`
	HasUppercase bool   `json:"hasUppercase"`
	HasNumber    bool   `json:"hasNumber"`
	HasSymbol    bool   `json:"hasSymbol"`
}

func (p *PasswordComplexityPolicy) ComplexityChanges(changed *PasswordComplexityPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Description != "" && p.Description != changed.Description {
		changes["description"] = changed.Description
	}
	if p.MinLength != changed.MinLength {
		changes["minLength"] = changed.MinLength
	}
	if p.HasLowercase != changed.HasLowercase {
		changes["hasLowercase"] = changed.HasLowercase
	}
	if p.HasUppercase != changed.HasUppercase {
		changes["hasUppercase"] = changed.HasUppercase
	}
	if p.HasNumber != changed.HasNumber {
		changes["hasNumber"] = changed.HasNumber
	}
	if p.HasSymbol != changed.HasSymbol {
		changes["hasSymbol"] = changed.HasSymbol
	}
	return changes
}

func PasswordComplexityPolicyFromModel(policy *model.PasswordComplexityPolicy) *PasswordComplexityPolicy {
	return &PasswordComplexityPolicy{
		ObjectRoot:   policy.ObjectRoot,
		Description:  policy.Description,
		State:        int32(policy.State),
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

func PasswordComplexityPolicyToModel(policy *PasswordComplexityPolicy) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot:   policy.ObjectRoot,
		Description:  policy.Description,
		State:        model.PolicyState(policy.State),
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

func (p *PasswordComplexityPolicy) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := p.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (p *PasswordComplexityPolicy) AppendEvent(event *es_models.Event) error {
	p.ObjectRoot.AppendEvent(event)

	switch event.Type {
	case model.PasswordComplexityPolicyAdded, model.PasswordComplexityPolicyChanged:
		if err := json.Unmarshal(event.Data, p); err != nil {
			logging.Log("EVEN-idl93").WithError(err).Error("could not unmarshal event data")
			return err
		}
		return nil
	}
	return nil
}
