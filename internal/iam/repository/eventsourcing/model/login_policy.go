package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type LoginPolicy struct {
	es_models.ObjectRoot
	State                 int32          `json:"-"`
	AllowUsernamePassword bool           `json:"allowUsernamePassword"`
	AllowRegister         bool           `json:"allowRegister"`
	AllowExternalIdp      bool           `json:"allowExternalIdp"`
	ForceMFA              bool           `json:"forceMFA"`
	PasswordlessType      int32          `json:"passwordlessType"`
	IDPProviders          []*IDPProvider `json:"-"`
	SecondFactors         []int32        `json:"-"`
	MultiFactors          []int32        `json:"-"`
}

type IDPProvider struct {
	es_models.ObjectRoot
	Type        int32  `json:"idpProviderType"`
	IDPConfigID string `json:"idpConfigId"`
}

type IDPProviderID struct {
	IDPConfigID string `json:"idpConfigId"`
}

type MFA struct {
	MFAType int32 `json:"mfaType"`
}

func GetIDPProvider(providers []*IDPProvider, id string) (int, *IDPProvider) {
	for i, p := range providers {
		if p.IDPConfigID == id {
			return i, p
		}
	}
	return -1, nil
}

func GetMFA(mfas []int32, mfaType int32) (int, int32) {
	for i, m := range mfas {
		if m == mfaType {
			return i, m
		}
	}
	return -1, 0
}
func LoginPolicyToModel(policy *LoginPolicy) *iam_model.LoginPolicy {
	idps := IDPProvidersToModel(policy.IDPProviders)
	secondFactors := SecondFactorsToModel(policy.SecondFactors)
	multiFactors := MultiFactorsToModel(policy.MultiFactors)
	return &iam_model.LoginPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 iam_model.PolicyState(policy.State),
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		IDPProviders:          idps,
		ForceMFA:              policy.ForceMFA,
		SecondFactors:         secondFactors,
		MultiFactors:          multiFactors,
		PasswordlessType:      iam_model.PasswordlessType(policy.PasswordlessType),
	}
}

func IDPProvidersToModel(members []*IDPProvider) []*iam_model.IDPProvider {
	convertedProviders := make([]*iam_model.IDPProvider, len(members))
	for i, m := range members {
		convertedProviders[i] = IDPProviderToModel(m)
	}
	return convertedProviders
}

func IDPProviderToModel(provider *IDPProvider) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		ObjectRoot:  provider.ObjectRoot,
		Type:        iam_model.IDPProviderType(provider.Type),
		IDPConfigID: provider.IDPConfigID,
	}
}

func SecondFactorsToModel(mfas []int32) []domain.SecondFactorType {
	convertedMFAs := make([]domain.SecondFactorType, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = domain.SecondFactorType(mfa)
	}
	return convertedMFAs
}

func MultiFactorsToModel(mfas []int32) []domain.MultiFactorType {
	convertedMFAs := make([]domain.MultiFactorType, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = domain.MultiFactorType(mfa)
	}
	return convertedMFAs
}

func (p *LoginPolicy) Changes(changed *LoginPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if changed.AllowUsernamePassword != p.AllowUsernamePassword {
		changes["allowUsernamePassword"] = changed.AllowUsernamePassword
	}
	if changed.AllowRegister != p.AllowRegister {
		changes["allowRegister"] = changed.AllowRegister
	}
	if changed.AllowExternalIdp != p.AllowExternalIdp {
		changes["allowExternalIdp"] = changed.AllowExternalIdp
	}
	if changed.ForceMFA != p.ForceMFA {
		changes["forceMFA"] = changed.ForceMFA
	}
	if changed.PasswordlessType != p.PasswordlessType {
		changes["passwordlessType"] = changed.PasswordlessType
	}
	return changes
}

func (p *LoginPolicy) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-7JS9d", "unable to unmarshal data")
	}
	return nil
}

func (p *IDPProvider) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-ldos9", "unable to unmarshal data")
	}
	return nil
}

func (m *MFA) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, m)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-4G9os", "unable to unmarshal data")
	}
	return nil
}
