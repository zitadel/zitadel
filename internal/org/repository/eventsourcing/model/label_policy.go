package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

// ToDo Michi
func (o *Org) appendAddLabelPolicyEvent(event *es_models.Event) error {
	o.LabelPolicy = new(iam_es_model.LabelPolicy)
	err := o.LabelPolicy.SetDataLabel(event)
	if err != nil {
		return err
	}
	o.LabelPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (o *Org) appendChangeLabelPolicyEvent(event *es_models.Event) error {
	return o.LabelPolicy.SetDataLabel(event)
}

func (o *Org) appendRemoveLabelPolicyEvent(event *es_models.Event) {
	o.LabelPolicy = nil
}

// func (o *Org) appendAddIdpProviderToLabelPolicyEvent(event *es_models.Event) error {
// 	provider := &iam_es_model.IDPProvider{}
// 	err := provider.SetData(event)
// 	if err != nil {
// 		return err
// 	}
// 	provider.ObjectRoot.CreationDate = event.CreationDate
// 	o.LabelPolicy.IDPProviders = append(o.LabelPolicy.IDPProviders, provider)
// 	return nil
// }

// func (o *Org) appendRemoveIdpProviderFromLabelPolicyEvent(event *es_models.Event) error {
// 	provider := &iam_es_model.IDPProvider{}
// 	err := provider.SetData(event)
// 	if err != nil {
// 		return err
// 	}
// 	if i, m := iam_es_model.GetIDPProvider(o.LabelPolicy.IDPProviders, provider.IDPConfigID); m != nil {
// 		o.LabelPolicy.IDPProviders[i] = o.LabelPolicy.IDPProviders[len(o.LabelPolicy.IDPProviders)-1]
// 		o.LabelPolicy.IDPProviders[len(o.LabelPolicy.IDPProviders)-1] = nil
// 		o.LabelPolicy.IDPProviders = o.LabelPolicy.IDPProviders[:len(o.LabelPolicy.IDPProviders)-1]
// 	}
// 	return nil
// }
