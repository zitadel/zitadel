package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

// ToDo Michi
type LabelPolicy struct {
	models.ObjectRoot
	State          int32  `json:"-"`
	PrimaryColor   string `json:"primaryColor"`
	SecundaryColor string `json:"secundaryColor"`
}

func LabelPolicyToModel(policy *LabelPolicy) *iam_model.LabelPolicy {
	return &iam_model.LabelPolicy{
		ObjectRoot:     policy.ObjectRoot,
		State:          iam_model.PolicyState(policy.State),
		PrimaryColor:   policy.PrimaryColor,
		SecundaryColor: policy.SecundaryColor,
	}
}

func LabelPolicyFromModel(policy *iam_model.LabelPolicy) *LabelPolicy {
	return &LabelPolicy{
		ObjectRoot:     policy.ObjectRoot,
		State:          int32(policy.State),
		PrimaryColor:   policy.PrimaryColor,
		SecundaryColor: policy.SecundaryColor,
	}
}

func (p *LabelPolicy) Changes(changed *LabelPolicy) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if changed.PrimaryColor != p.PrimaryColor {
		changes["primaryColor"] = changed.PrimaryColor
	}
	if changed.SecundaryColor != p.SecundaryColor {
		changes["secundaryColor"] = changed.SecundaryColor
	}

	return changes
}

// func (i *IAM) appendAddLabelPolicyEvent(event *es_models.Event) error {
// 	i.DefaultLabelPolicy = new(LabelPolicy)
// 	err := i.DefaultLabelPolicy.SetData(event)
// 	if err != nil {
// 		return err
// 	}
// 	i.DefaultLabelPolicy.ObjectRoot.CreationDate = event.CreationDate
// 	return nil
// }

func (i *IAM) appendAddLabelPolicyEvent(event *es_models.Event) error {
	i.DefaultLabelPolicy = new(LabelPolicy)
	err := i.DefaultLabelPolicy.SetDataLabel(event)
	if err != nil {
		return err
	}
	i.DefaultLabelPolicy.ObjectRoot.CreationDate = event.CreationDate
	return nil
}

func (i *IAM) appendChangeLabelPolicyEvent(event *es_models.Event) error {
	return i.DefaultLabelPolicy.SetDataLabel(event)
}

// func (iam *IAM) appendAddIDPProviderToLabelPolicyEvent(event *es_models.Event) error {
// 	provider := new(IDPProvider)
// 	err := provider.SetData(event)
// 	if err != nil {
// 		return err
// 	}
// 	provider.ObjectRoot.CreationDate = event.CreationDate
// 	iam.DefaultLabelPolicy.IDPProviders = append(iam.DefaultLabelPolicy.IDPProviders, provider)
// 	return nil
// }

// func (iam *IAM) appendRemoveIDPProviderFromLabelPolicyEvent(event *es_models.Event) error {
// 	provider := new(IDPProvider)
// 	err := provider.SetData(event)
// 	if err != nil {
// 		return err
// 	}
// 	if i, m := GetIDPProvider(iam.DefaultLabelPolicy.IDPProviders, provider.IDPConfigID); m != nil {
// 		iam.DefaultLabelPolicy.IDPProviders[i] = iam.DefaultLabelPolicy.IDPProviders[len(iam.DefaultLabelPolicy.IDPProviders)-1]
// 		iam.DefaultLabelPolicy.IDPProviders[len(iam.DefaultLabelPolicy.IDPProviders)-1] = nil
// 		iam.DefaultLabelPolicy.IDPProviders = iam.DefaultLabelPolicy.IDPProviders[:len(iam.DefaultLabelPolicy.IDPProviders)-1]
// 	}
// 	return nil
// }

func (p *LabelPolicy) SetDataLabel(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "MODEL-ikjhf", "unable to unmarshal data")
	}
	return nil
}

func (p *IDPProvider) SetDataLabel(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "MODEL-c41Hn", "unable to unmarshal data")
	}
	return nil
}
