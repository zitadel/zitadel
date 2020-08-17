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
	provider := &iam_es_model.IdpProvider{}
	err := provider.SetData(event)
	if err != nil {
		return err
	}
	provider.ObjectRoot.CreationDate = event.CreationDate
	o.LoginPolicy.IdpProviders = append(o.LoginPolicy.IdpProviders, provider)
	return nil
}

func (o *Org) appendRemoveIdpProviderFromLoginPolicyEvent(event *es_models.Event) error {
	provider := &iam_es_model.IdpProvider{}
	err := provider.SetData(event)
	if err != nil {
		return err
	}
	if i, m := iam_es_model.GetIdpProvider(o.LoginPolicy.IdpProviders, provider.IdpConfigID); m != nil {
		o.LoginPolicy.IdpProviders[i] = o.LoginPolicy.IdpProviders[len(o.LoginPolicy.IdpProviders)-1]
		o.LoginPolicy.IdpProviders[len(o.LoginPolicy.IdpProviders)-1] = nil
		o.LoginPolicy.IdpProviders = o.LoginPolicy.IdpProviders[:len(o.LoginPolicy.IdpProviders)-1]
	}
	return nil
}
