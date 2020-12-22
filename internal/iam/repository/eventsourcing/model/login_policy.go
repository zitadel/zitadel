package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type LoginPolicy struct {
	models.ObjectRoot
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
	models.ObjectRoot
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

func LoginPolicyFromModel(policy *iam_model.LoginPolicy) *LoginPolicy {
	idps := IDOProvidersFromModel(policy.IDPProviders)
	secondFactors := SecondFactorsFromModel(policy.SecondFactors)
	multiFactors := MultiFactorsFromModel(policy.MultiFactors)
	return &LoginPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 int32(policy.State),
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		IDPProviders:          idps,
		ForceMFA:              policy.ForceMFA,
		SecondFactors:         secondFactors,
		MultiFactors:          multiFactors,
		PasswordlessType:      int32(policy.PasswordlessType),
	}
}

func IDPProvidersToModel(members []*IDPProvider) []*iam_model.IDPProvider {
	convertedProviders := make([]*iam_model.IDPProvider, len(members))
	for i, m := range members {
		convertedProviders[i] = IDPProviderToModel(m)
	}
	return convertedProviders
}

func IDOProvidersFromModel(members []*iam_model.IDPProvider) []*IDPProvider {
	convertedProviders := make([]*IDPProvider, len(members))
	for i, m := range members {
		convertedProviders[i] = IDPProviderFromModel(m)
	}
	return convertedProviders
}

func IDPProviderToModel(provider *IDPProvider) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		ObjectRoot:  provider.ObjectRoot,
		Type:        iam_model.IDPProviderType(provider.Type),
		IdpConfigID: provider.IDPConfigID,
	}
}

func IDPProviderFromModel(provider *iam_model.IDPProvider) *IDPProvider {
	return &IDPProvider{
		ObjectRoot:  provider.ObjectRoot,
		Type:        int32(provider.Type),
		IDPConfigID: provider.IdpConfigID,
	}
}

func SecondFactorsFromModel(mfas []iam_model.SecondFactorType) []int32 {
	convertedMFAs := make([]int32, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = int32(mfa)
	}
	return convertedMFAs
}

func SecondFactorFromModel(mfa iam_model.SecondFactorType) *MFA {
	return &MFA{MFAType: int32(mfa)}
}

func SecondFactorsToModel(mfas []int32) []iam_model.SecondFactorType {
	convertedMFAs := make([]iam_model.SecondFactorType, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = iam_model.SecondFactorType(mfa)
	}
	return convertedMFAs
}

func MultiFactorsFromModel(mfas []iam_model.MultiFactorType) []int32 {
	convertedMFAs := make([]int32, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = int32(mfa)
	}
	return convertedMFAs
}

func MultiFactorFromModel(mfa iam_model.MultiFactorType) *MFA {
	return &MFA{MFAType: int32(mfa)}
}

func MultiFactorsToModel(mfas []int32) []iam_model.MultiFactorType {
	convertedMFAs := make([]iam_model.MultiFactorType, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = iam_model.MultiFactorType(mfa)
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

func (i *IAM) appendAddLoginPolicyEvent(event *es_models.Event) error {
	i.DefaultLoginPolicy = new(LoginPolicy)
	err := i.DefaultLoginPolicy.SetData(event)
	if err != nil {
		return err
	}
	i.DefaultLoginPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (i *IAM) appendChangeLoginPolicyEvent(event *es_models.Event) error {
	return i.DefaultLoginPolicy.SetData(event)
}

func (iam *IAM) appendAddIDPProviderToLoginPolicyEvent(event *es_models.Event) error {
	provider := new(IDPProvider)
	err := provider.SetData(event)
	if err != nil {
		return err
	}
	provider.ObjectRoot.CreationDate = event.CreationDate
	iam.DefaultLoginPolicy.IDPProviders = append(iam.DefaultLoginPolicy.IDPProviders, provider)
	return nil
}

func (iam *IAM) appendRemoveIDPProviderFromLoginPolicyEvent(event *es_models.Event) error {
	provider := new(IDPProvider)
	err := provider.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetIDPProvider(iam.DefaultLoginPolicy.IDPProviders, provider.IDPConfigID); m != nil {
		iam.DefaultLoginPolicy.IDPProviders[i] = iam.DefaultLoginPolicy.IDPProviders[len(iam.DefaultLoginPolicy.IDPProviders)-1]
		iam.DefaultLoginPolicy.IDPProviders[len(iam.DefaultLoginPolicy.IDPProviders)-1] = nil
		iam.DefaultLoginPolicy.IDPProviders = iam.DefaultLoginPolicy.IDPProviders[:len(iam.DefaultLoginPolicy.IDPProviders)-1]
		return nil
	}
	return nil
}

func (iam *IAM) appendAddSecondFactorToLoginPolicyEvent(event *es_models.Event) error {
	mfa := new(MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	iam.DefaultLoginPolicy.SecondFactors = append(iam.DefaultLoginPolicy.SecondFactors, mfa.MFAType)
	return nil
}

func (iam *IAM) appendRemoveSecondFactorFromLoginPolicyEvent(event *es_models.Event) error {
	mfa := new(MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetMFA(iam.DefaultLoginPolicy.SecondFactors, mfa.MFAType); m != 0 {
		iam.DefaultLoginPolicy.SecondFactors[i] = iam.DefaultLoginPolicy.SecondFactors[len(iam.DefaultLoginPolicy.SecondFactors)-1]
		iam.DefaultLoginPolicy.SecondFactors[len(iam.DefaultLoginPolicy.SecondFactors)-1] = 0
		iam.DefaultLoginPolicy.SecondFactors = iam.DefaultLoginPolicy.SecondFactors[:len(iam.DefaultLoginPolicy.SecondFactors)-1]
		return nil
	}
	return nil
}

func (iam *IAM) appendAddMultiFactorToLoginPolicyEvent(event *es_models.Event) error {
	mfa := new(MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	iam.DefaultLoginPolicy.MultiFactors = append(iam.DefaultLoginPolicy.MultiFactors, mfa.MFAType)
	return nil
}

func (iam *IAM) appendRemoveMultiFactorFromLoginPolicyEvent(event *es_models.Event) error {
	mfa := new(MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetMFA(iam.DefaultLoginPolicy.MultiFactors, mfa.MFAType); m != 0 {
		iam.DefaultLoginPolicy.MultiFactors[i] = iam.DefaultLoginPolicy.MultiFactors[len(iam.DefaultLoginPolicy.MultiFactors)-1]
		iam.DefaultLoginPolicy.MultiFactors[len(iam.DefaultLoginPolicy.MultiFactors)-1] = 0
		iam.DefaultLoginPolicy.MultiFactors = iam.DefaultLoginPolicy.MultiFactors[:len(iam.DefaultLoginPolicy.MultiFactors)-1]
		return nil
	}
	return nil
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
