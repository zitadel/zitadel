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
	IDPProviders          []*IDPProvider `json:"-"`
	SoftwareMFAs          []int32        `json:"-"`
	HardwareMFAs          []int32        `json:"-"`
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
	MfaType int32 `json:"mfaType"`
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
	softwareMFAs := SoftwareMFAsToModel(policy.SoftwareMFAs)
	hardwareMFAs := HardwareMFAsToModel(policy.HardwareMFAs)
	return &iam_model.LoginPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 iam_model.PolicyState(policy.State),
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		IDPProviders:          idps,
		ForceMFA:              policy.ForceMFA,
		SoftwareMFAs:          softwareMFAs,
		HardwareMFAs:          hardwareMFAs,
	}
}

func LoginPolicyFromModel(policy *iam_model.LoginPolicy) *LoginPolicy {
	idps := IDOProvidersFromModel(policy.IDPProviders)
	softwareMFAs := SoftwareMFAsFromModel(policy.SoftwareMFAs)
	hardwareMFAs := HardwareMFAsFromModel(policy.HardwareMFAs)
	return &LoginPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 int32(policy.State),
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		IDPProviders:          idps,
		ForceMFA:              policy.ForceMFA,
		SoftwareMFAs:          softwareMFAs,
		HardwareMFAs:          hardwareMFAs,
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

func SoftwareMFAsFromModel(mfas []iam_model.SoftwareMFAType) []int32 {
	convertedMFAs := make([]int32, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = int32(mfa)
	}
	return convertedMFAs
}

func SoftwareMFAsToModel(mfas []int32) []iam_model.SoftwareMFAType {
	convertedMFAs := make([]iam_model.SoftwareMFAType, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = iam_model.SoftwareMFAType(mfa)
	}
	return convertedMFAs
}

func HardwareMFAsFromModel(mfas []iam_model.HardwareMFAType) []int32 {
	convertedMFAs := make([]int32, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = int32(mfa)
	}
	return convertedMFAs
}

func HardwareMFAsToModel(mfas []int32) []iam_model.HardwareMFAType {
	convertedMFAs := make([]iam_model.HardwareMFAType, len(mfas))
	for i, mfa := range mfas {
		convertedMFAs[i] = iam_model.HardwareMFAType(mfa)
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

func (iam *IAM) appendAddSoftwareMFAToLoginPolicyEvent(event *es_models.Event) error {
	mfa := new(MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	iam.DefaultLoginPolicy.SoftwareMFAs = append(iam.DefaultLoginPolicy.SoftwareMFAs, mfa.MfaType)
	return nil
}

func (iam *IAM) appendRemoveSoftwareMfaFromLoginPolicyEvent(event *es_models.Event) error {
	mfa := new(MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetMFA(iam.DefaultLoginPolicy.SoftwareMFAs, mfa.MfaType); m != 0 {
		iam.DefaultLoginPolicy.SoftwareMFAs[i] = iam.DefaultLoginPolicy.SoftwareMFAs[len(iam.DefaultLoginPolicy.SoftwareMFAs)-1]
		iam.DefaultLoginPolicy.SoftwareMFAs[len(iam.DefaultLoginPolicy.SoftwareMFAs)-1] = 0
		iam.DefaultLoginPolicy.SoftwareMFAs = iam.DefaultLoginPolicy.SoftwareMFAs[:len(iam.DefaultLoginPolicy.SoftwareMFAs)-1]
		return nil
	}
	return nil
}

func (iam *IAM) appendAddHardwareMFAToLoginPolicyEvent(event *es_models.Event) error {
	mfa := new(MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	iam.DefaultLoginPolicy.HardwareMFAs = append(iam.DefaultLoginPolicy.HardwareMFAs, mfa.MfaType)
	return nil
}

func (iam *IAM) appendRemoveHardwareMfaFromLoginPolicyEvent(event *es_models.Event) error {
	mfa := new(MFA)
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetMFA(iam.DefaultLoginPolicy.HardwareMFAs, mfa.MfaType); m != 0 {
		iam.DefaultLoginPolicy.HardwareMFAs[i] = iam.DefaultLoginPolicy.HardwareMFAs[len(iam.DefaultLoginPolicy.HardwareMFAs)-1]
		iam.DefaultLoginPolicy.HardwareMFAs[len(iam.DefaultLoginPolicy.HardwareMFAs)-1] = 0
		iam.DefaultLoginPolicy.HardwareMFAs = iam.DefaultLoginPolicy.HardwareMFAs[:len(iam.DefaultLoginPolicy.HardwareMFAs)-1]
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
