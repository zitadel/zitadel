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

func (o *Org) appendAddSecondFactorToLoginPolicyEvent(event *es_models.Event) error {
	mfa := &iam_es_model.MFA{}
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	o.LoginPolicy.SecondFactors = append(o.LoginPolicy.SecondFactors, mfa.MFAType)
	return nil
}

func (o *Org) appendRemoveSecondFactorFromLoginPolicyEvent(event *es_models.Event) error {
	mfa := &iam_es_model.MFA{}
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	if i, m := iam_es_model.GetMFA(o.LoginPolicy.SecondFactors, mfa.MFAType); m != 0 {
		o.LoginPolicy.SecondFactors[i] = o.LoginPolicy.SecondFactors[len(o.LoginPolicy.SecondFactors)-1]
		o.LoginPolicy.SecondFactors[len(o.LoginPolicy.SecondFactors)-1] = 0
		o.LoginPolicy.SecondFactors = o.LoginPolicy.SecondFactors[:len(o.LoginPolicy.SecondFactors)-1]
		return nil
	}
	return nil
}

func (o *Org) appendAddMultiFactorToLoginPolicyEvent(event *es_models.Event) error {
	mfa := &iam_es_model.MFA{}
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	o.LoginPolicy.MultiFactors = append(o.LoginPolicy.MultiFactors, mfa.MFAType)
	return nil
}

func (o *Org) appendRemoveMultiFactorFromLoginPolicyEvent(event *es_models.Event) error {
	mfa := &iam_es_model.MFA{}
	err := mfa.SetData(event)
	if err != nil {
		return err
	}
	if i, m := iam_es_model.GetMFA(o.LoginPolicy.MultiFactors, mfa.MFAType); m != 0 {
		o.LoginPolicy.MultiFactors[i] = o.LoginPolicy.MultiFactors[len(o.LoginPolicy.MultiFactors)-1]
		o.LoginPolicy.MultiFactors[len(o.LoginPolicy.MultiFactors)-1] = 0
		o.LoginPolicy.MultiFactors = o.LoginPolicy.MultiFactors[:len(o.LoginPolicy.MultiFactors)-1]
		return nil
	}
	return nil
}
