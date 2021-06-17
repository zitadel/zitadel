package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type LabelPolicy struct {
	es_models.ObjectRoot
	State               int32  `json:"-"`
	PrimaryColor        string `json:"primaryColor"`
	BackgroundColor     string `json:"backgroundColor"`
	FontColor           string `json:"fontColor"`
	WarnColor           string `json:"warnColor"`
	PrimaryColorDark    string `json:"primaryColorDark"`
	BackgroundColorDark string `json:"backgroundColorDark"`
	FontColorDark       string `json:"fontColorDark"`
	WarnColorDark       string `json:"warnColorDark"`
	HideLoginNameSuffix bool   `json:"hideLoginNameSuffix"`
}

func LabelPolicyToModel(policy *LabelPolicy) *iam_model.LabelPolicy {
	return &iam_model.LabelPolicy{
		ObjectRoot:          policy.ObjectRoot,
		State:               iam_model.PolicyState(policy.State),
		PrimaryColor:        policy.PrimaryColor,
		BackgroundColor:     policy.BackgroundColor,
		WarnColor:           policy.WarnColor,
		FontColor:           policy.FontColor,
		PrimaryColorDark:    policy.PrimaryColorDark,
		BackgroundColorDark: policy.BackgroundColorDark,
		WarnColorDark:       policy.WarnColorDark,
		FontColorDark:       policy.FontColorDark,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
	}
}

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

func (p *LabelPolicy) SetDataLabel(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, p)
	if err != nil {
		return errors.ThrowInternal(err, "MODEL-Gdgwq", "unable to unmarshal data")
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
