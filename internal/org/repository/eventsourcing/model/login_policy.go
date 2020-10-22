package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func (o *Org) appendAddLoginPolicyEvent(event *es_models.Event) error {
	o.LoginPolicy = new(iam_es_model.LoginPolicy)
	err := o.LoginPolicy.SetData(event)
	if err != nil {
		return err
	}
	o.LoginPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeLoginPolicyEvent(event *es_models.Event) error {
	return o.LoginPolicy.SetData(event)
}

func (o *Org) appendRemoveLoginPolicyEvent(event *es_models.Event) {
	o.LoginPolicy = nil
}

func (o *Org) appendAddIdpProviderToLoginPolicyEvent(event *es_models.Event) error {
	provider := &iam_es_model.IDPProvider{}
	err := provider.SetData(event)
	if err != nil {
		return err
	}
	provider.ObjectRoot.CreationDate = event.CreationDate
	o.LoginPolicy.IDPProviders = append(o.LoginPolicy.IDPProviders, provider)
	return nil
}

func (o *Org) appendRemoveIdpProviderFromLoginPolicyEvent(event *es_models.Event) error {
	provider := &iam_es_model.IDPProvider{}
	err := provider.SetData(event)
	if err != nil {
		return err
	}
	if i, m := iam_es_model.GetIDPProvider(o.LoginPolicy.IDPProviders, provider.IDPConfigID); m != nil {
		o.LoginPolicy.IDPProviders[i] = o.LoginPolicy.IDPProviders[len(o.LoginPolicy.IDPProviders)-1]
		o.LoginPolicy.IDPProviders[len(o.LoginPolicy.IDPProviders)-1] = nil
		o.LoginPolicy.IDPProviders = o.LoginPolicy.IDPProviders[:len(o.LoginPolicy.IDPProviders)-1]
		return nil
	}
	return nil
}

func (o *Org) appendAddSoftwareMFAToLoginPolicyEvent(event *es_models.Event) error {
	mfa := &iam_es_model.MFA{}
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	o.LoginPolicy.SoftwareMFAs = append(o.LoginPolicy.SoftwareMFAs, mfa.MfaType)
	return nil
}

func (o *Org) appendRemoveSoftwareMFAFromLoginPolicyEvent(event *es_models.Event) error {
	mfa := &iam_es_model.MFA{}
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	if i, m := iam_es_model.GetMFA(o.LoginPolicy.SoftwareMFAs, mfa.MfaType); m != 0 {
		o.LoginPolicy.SoftwareMFAs[i] = o.LoginPolicy.SoftwareMFAs[len(o.LoginPolicy.SoftwareMFAs)-1]
		o.LoginPolicy.SoftwareMFAs[len(o.LoginPolicy.SoftwareMFAs)-1] = 0
		o.LoginPolicy.SoftwareMFAs = o.LoginPolicy.SoftwareMFAs[:len(o.LoginPolicy.SoftwareMFAs)-1]
		return nil
	}
	return nil
}

func (o *Org) appendAddHardwareMFAToLoginPolicyEvent(event *es_models.Event) error {
	mfa := &iam_es_model.MFA{}
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	o.LoginPolicy.HardwareMFAs = append(o.LoginPolicy.HardwareMFAs, mfa.MfaType)
	return nil
}

func (o *Org) appendRemoveHardwareMFAFromLoginPolicyEvent(event *es_models.Event) error {
	mfa := &iam_es_model.MFA{}
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	if i, m := iam_es_model.GetMFA(o.LoginPolicy.HardwareMFAs, mfa.MfaType); m != 0 {
		o.LoginPolicy.HardwareMFAs[i] = o.LoginPolicy.HardwareMFAs[len(o.LoginPolicy.HardwareMFAs)-1]
		o.LoginPolicy.HardwareMFAs[len(o.LoginPolicy.HardwareMFAs)-1] = 0
		o.LoginPolicy.HardwareMFAs = o.LoginPolicy.HardwareMFAs[:len(o.LoginPolicy.HardwareMFAs)-1]
		return nil
	}
	return nil
}
