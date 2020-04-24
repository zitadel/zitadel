package eventsourcing

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
)

type PasswordComplexityPolicy struct {
	models.ObjectRoot

	Description  string
	State        int32
	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool
}

func (p *PasswordComplexityPolicy) ComplexityChanges(changed *PasswordComplexityPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Description != "" && p.Description != changed.Description {
		changes["Description"] = changed.Description
	}
	if p.MinLength != changed.MinLength {
		changes["MinLength"] = changed.MinLength
	}
	if p.HasLowercase != changed.HasLowercase {
		changes["HasLowercase"] = changed.HasLowercase
	}
	if p.HasUppercase != changed.HasUppercase {
		changes["HasUppercase"] = changed.HasUppercase
	}
	if p.HasNumber != changed.HasNumber {
		changes["HasNumber"] = changed.HasNumber
	}
	if p.HasSymbol != changed.HasSymbol {
		changes["HasSymbol"] = changed.HasSymbol
	}
	return changes
}

func PasswordComplexityPolicyFromModel(policy *model.PasswordComplexityPolicy) *PasswordComplexityPolicy {
	return &PasswordComplexityPolicy{
		ObjectRoot: models.ObjectRoot{
			ID:           policy.ObjectRoot.ID,
			ChangeDate:   policy.ChangeDate,
			CreationDate: policy.CreationDate,
			Sequence:     policy.Sequence,
		},
		Description:  policy.Description,
		State:        policy.State,
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

func PasswordComplexityPolicyToModel(policy *PasswordComplexityPolicy) *model.PasswordComplexityPolicy {
	return &model.PasswordComplexityPolicy{
		ObjectRoot: models.ObjectRoot{
			ID:           policy.ID,
			ChangeDate:   policy.ChangeDate,
			CreationDate: policy.CreationDate,
			Sequence:     policy.Sequence,
		},
		Description:  policy.Description,
		State:        policy.State,
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
