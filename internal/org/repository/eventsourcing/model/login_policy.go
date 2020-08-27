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
	}
	return nil
}
