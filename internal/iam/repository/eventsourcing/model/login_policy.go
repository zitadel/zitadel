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
}

type IDPProvider struct {
	models.ObjectRoot
	Type        int32  `json:"idpProviderType"`
	IDPConfigID string `json:"idpConfigId"`
}

type IDPProviderID struct {
	IDPConfigID string `json:"idpConfigId"`
}

func GetIDPProvider(providers []*IDPProvider, id string) (int, *IDPProvider) {
	for i, p := range providers {
		if p.IDPConfigID == id {
			return i, p
		}
	}
	return -1, nil
}

func LoginPolicyToModel(policy *LoginPolicy) *iam_model.LoginPolicy {
	idps := IDPProvidersToModel(policy.IDPProviders)
	return &iam_model.LoginPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 iam_model.PolicyState(policy.State),
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		IDPProviders:          idps,
		ForceMFA:              policy.ForceMFA,
	}
}

func LoginPolicyFromModel(policy *iam_model.LoginPolicy) *LoginPolicy {
	idps := IDOProvidersFromModel(policy.IDPProviders)
	return &LoginPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 int32(policy.State),
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		IDPProviders:          idps,
		ForceMFA:              policy.ForceMFA,
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
