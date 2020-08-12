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
	IdpProviders          []*IdpProvider `json:"-"`
}

type IdpProvider struct {
	models.ObjectRoot
	Type        int32  `json:"idpProviderType"`
	IdpConfigID string `json:"idpConfigId"`
}

type IdpProviderID struct {
	IdpConfigID string `json:"idpConfigId"`
}

func GetIdpProvider(providers []*IdpProvider, id string) (int, *IdpProvider) {
	for i, p := range providers {
		if p.IdpConfigID == id {
			return i, p
		}
	}
	return -1, nil
}

func LoginPolicyToModel(policy *LoginPolicy) *iam_model.LoginPolicy {
	idps := IdpProvidersToModel(policy.IdpProviders)
	return &iam_model.LoginPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 iam_model.PolicyState(policy.State),
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		IdpProviders:          idps,
	}
}

func LoginPolicyFromModel(policy *iam_model.LoginPolicy) *LoginPolicy {
	idps := IdpProvidersFromModel(policy.IdpProviders)
	return &LoginPolicy{
		ObjectRoot:            policy.ObjectRoot,
		State:                 int32(policy.State),
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowExternalIdp,
		IdpProviders:          idps,
	}
}

func IdpProvidersToModel(members []*IdpProvider) []*iam_model.IdpProvider {
	convertedProviders := make([]*iam_model.IdpProvider, len(members))
	for i, m := range members {
		convertedProviders[i] = IdpProviderToModel(m)
	}
	return convertedProviders
}

func IdpProvidersFromModel(members []*iam_model.IdpProvider) []*IdpProvider {
	convertedProviders := make([]*IdpProvider, len(members))
	for i, m := range members {
		convertedProviders[i] = IdpProviderFromModel(m)
	}
	return convertedProviders
}

func IdpProviderToModel(provider *IdpProvider) *iam_model.IdpProvider {
	return &iam_model.IdpProvider{
		ObjectRoot:  provider.ObjectRoot,
		Type:        iam_model.IdpProviderType(provider.Type),
		IdpConfigID: provider.IdpConfigID,
	}
}

func IdpProviderFromModel(provider *iam_model.IdpProvider) *IdpProvider {
	return &IdpProvider{
		ObjectRoot:  provider.ObjectRoot,
		Type:        int32(provider.Type),
		IdpConfigID: provider.IdpConfigID,
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

	return changes
}

func (i *Iam) appendAddLoginPolicyEvent(event *es_models.Event) error {
	i.DefaultLoginPolicy = new(LoginPolicy)
	err := i.DefaultLoginPolicy.SetData(event)
	if err != nil {
		return err
	}
	i.DefaultLoginPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (i *Iam) appendChangeLoginPolicyEvent(event *es_models.Event) error {
	return i.DefaultLoginPolicy.SetData(event)
}

func (iam *Iam) appendAddIdpProviderToLoginPolicyEvent(event *es_models.Event) error {
	provider := &IdpProvider{}
	err := provider.SetData(event)
	if err != nil {
		return err
	}
	provider.ObjectRoot.CreationDate = event.CreationDate
	iam.DefaultLoginPolicy.IdpProviders = append(iam.DefaultLoginPolicy.IdpProviders, provider)
	return nil
}

func (iam *Iam) appendRemoveIdpProviderFromLoginPolicyEvent(event *es_models.Event) error {
	provider := &IdpProvider{}
	err := provider.SetData(event)
	if err != nil {
		return err
	}
	if i, m := GetIdpProvider(iam.DefaultLoginPolicy.IdpProviders, provider.IdpConfigID); m != nil {
		iam.DefaultLoginPolicy.IdpProviders[i] = iam.DefaultLoginPolicy.IdpProviders[len(iam.DefaultLoginPolicy.IdpProviders)-1]
		iam.DefaultLoginPolicy.IdpProviders[len(iam.DefaultLoginPolicy.IdpProviders)-1] = nil
		iam.DefaultLoginPolicy.IdpProviders = iam.DefaultLoginPolicy.IdpProviders[:len(iam.DefaultLoginPolicy.IdpProviders)-1]
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

func (p *IdpProvider) SetData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-ldos9", "unable to unmarshal data")
	}
	return nil
}
