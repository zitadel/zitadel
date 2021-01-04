package model

import (
	"encoding/json"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	IAMVersion = "v1"
)

type Step int

const (
	Step1     = Step(model.Step1)
	Step2     = Step(model.Step2)
	StepCount = Step(model.StepCount)
)

type IAM struct {
	es_models.ObjectRoot
	SetUpStarted                    Step                      `json:"-"`
	SetUpDone                       Step                      `json:"-"`
	GlobalOrgID                     string                    `json:"globalOrgId,omitempty"`
	IAMProjectID                    string                    `json:"iamProjectId,omitempty"`
	Members                         []*IAMMember              `json:"-"`
	IDPs                            []*IDPConfig              `json:"-"`
	DefaultLoginPolicy              *LoginPolicy              `json:"-"`
	DefaultLabelPolicy              *LabelPolicy              `json:"-"`
	DefaultOrgIAMPolicy             *OrgIAMPolicy             `json:"-"`
	DefaultPasswordComplexityPolicy *PasswordComplexityPolicy `json:"-"`
	DefaultPasswordAgePolicy        *PasswordAgePolicy        `json:"-"`
	DefaultPasswordLockoutPolicy    *PasswordLockoutPolicy    `json:"-"`
}

func IAMFromModel(iam *model.IAM) *IAM {
	members := IAMMembersFromModel(iam.Members)
	idps := IDPConfigsFromModel(iam.IDPs)
	converted := &IAM{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: Step(iam.SetUpStarted),
		SetUpDone:    Step(iam.SetUpDone),
		GlobalOrgID:  iam.GlobalOrgID,
		IAMProjectID: iam.IAMProjectID,
		Members:      members,
		IDPs:         idps,
	}
	if iam.DefaultLoginPolicy != nil {
		converted.DefaultLoginPolicy = LoginPolicyFromModel(iam.DefaultLoginPolicy)
	}
	if iam.DefaultLabelPolicy != nil {
		converted.DefaultLabelPolicy = LabelPolicyFromModel(iam.DefaultLabelPolicy)
	}
	if iam.DefaultPasswordComplexityPolicy != nil {
		converted.DefaultPasswordComplexityPolicy = PasswordComplexityPolicyFromModel(iam.DefaultPasswordComplexityPolicy)
	}
	if iam.DefaultPasswordAgePolicy != nil {
		converted.DefaultPasswordAgePolicy = PasswordAgePolicyFromModel(iam.DefaultPasswordAgePolicy)
	}
	if iam.DefaultPasswordLockoutPolicy != nil {
		converted.DefaultPasswordLockoutPolicy = PasswordLockoutPolicyFromModel(iam.DefaultPasswordLockoutPolicy)
	}
	if iam.DefaultOrgIAMPolicy != nil {
		converted.DefaultOrgIAMPolicy = OrgIAMPolicyFromModel(iam.DefaultOrgIAMPolicy)
	}
	return converted
}

func IAMToModel(iam *IAM) *model.IAM {
	members := IAMMembersToModel(iam.Members)
	idps := IDPConfigsToModel(iam.IDPs)
	converted := &model.IAM{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: domain.Step(iam.SetUpStarted),
		SetUpDone:    domain.Step(iam.SetUpDone),
		GlobalOrgID:  iam.GlobalOrgID,
		IAMProjectID: iam.IAMProjectID,
		Members:      members,
		IDPs:         idps,
	}
	if iam.DefaultLoginPolicy != nil {
		converted.DefaultLoginPolicy = LoginPolicyToModel(iam.DefaultLoginPolicy)
	}
	if iam.DefaultLabelPolicy != nil {
		converted.DefaultLabelPolicy = LabelPolicyToModel(iam.DefaultLabelPolicy)
	}
	if iam.DefaultPasswordComplexityPolicy != nil {
		converted.DefaultPasswordComplexityPolicy = PasswordComplexityPolicyToModel(iam.DefaultPasswordComplexityPolicy)
	}
	if iam.DefaultPasswordAgePolicy != nil {
		converted.DefaultPasswordAgePolicy = PasswordAgePolicyToModel(iam.DefaultPasswordAgePolicy)
	}
	if iam.DefaultPasswordLockoutPolicy != nil {
		converted.DefaultPasswordLockoutPolicy = PasswordLockoutPolicyToModel(iam.DefaultPasswordLockoutPolicy)
	}
	if iam.DefaultOrgIAMPolicy != nil {
		converted.DefaultOrgIAMPolicy = OrgIAMPolicyToModel(iam.DefaultOrgIAMPolicy)
	}
	return converted
}

func (i *IAM) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := i.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (i *IAM) AppendEvent(event *es_models.Event) (err error) {
	i.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case IAMSetupStarted:
		if len(event.Data) == 0 {
			i.SetUpStarted = Step(model.Step1)
			return
		}
		step := new(struct{ Step Step })
		err = json.Unmarshal(event.Data, step)
		if err != nil {
			return err
		}
		i.SetUpStarted = step.Step
	case IAMSetupDone:
		if len(event.Data) == 0 {
			i.SetUpDone = Step(model.Step1)
			return
		}
		step := new(struct{ Step Step })
		err = json.Unmarshal(event.Data, step)
		if err != nil {
			return err
		}
		i.SetUpDone = step.Step
	case IAMProjectSet,
		GlobalOrgSet:
		err = i.SetData(event)
	case IAMMemberAdded:
		err = i.appendAddMemberEvent(event)
	case IAMMemberChanged:
		err = i.appendChangeMemberEvent(event)
	case IAMMemberRemoved:
		err = i.appendRemoveMemberEvent(event)
	case IDPConfigAdded:
		return i.appendAddIDPConfigEvent(event)
	case IDPConfigChanged:
		return i.appendChangeIDPConfigEvent(event)
	case IDPConfigRemoved:
		return i.appendRemoveIDPConfigEvent(event)
	case IDPConfigDeactivated:
		return i.appendIDPConfigStateEvent(event, model.IDPConfigStateInactive)
	case IDPConfigReactivated:
		return i.appendIDPConfigStateEvent(event, model.IDPConfigStateActive)
	case OIDCIDPConfigAdded:
		return i.appendAddOIDCIDPConfigEvent(event)
	case OIDCIDPConfigChanged:
		return i.appendChangeOIDCIDPConfigEvent(event)
	case LoginPolicyAdded:
		return i.appendAddLoginPolicyEvent(event)
	case LoginPolicyChanged:
		return i.appendChangeLoginPolicyEvent(event)
	case LoginPolicyIDPProviderAdded:
		return i.appendAddIDPProviderToLoginPolicyEvent(event)
	case LoginPolicyIDPProviderRemoved:
		return i.appendRemoveIDPProviderFromLoginPolicyEvent(event)
	case LoginPolicySecondFactorAdded:
		return i.appendAddSecondFactorToLoginPolicyEvent(event)
	case LoginPolicySecondFactorRemoved:
		return i.appendRemoveSecondFactorFromLoginPolicyEvent(event)
	case LoginPolicyMultiFactorAdded:
		return i.appendAddMultiFactorToLoginPolicyEvent(event)
	case LoginPolicyMultiFactorRemoved:
		return i.appendRemoveMultiFactorFromLoginPolicyEvent(event)
	case LabelPolicyAdded:
		return i.appendAddLabelPolicyEvent(event)
	case LabelPolicyChanged:
		return i.appendChangeLabelPolicyEvent(event)
	case PasswordComplexityPolicyAdded:
		return i.appendAddPasswordComplexityPolicyEvent(event)
	case PasswordComplexityPolicyChanged:
		return i.appendChangePasswordComplexityPolicyEvent(event)
	case PasswordAgePolicyAdded:
		return i.appendAddPasswordAgePolicyEvent(event)
	case PasswordAgePolicyChanged:
		return i.appendChangePasswordAgePolicyEvent(event)
	case PasswordLockoutPolicyAdded:
		return i.appendAddPasswordLockoutPolicyEvent(event)
	case PasswordLockoutPolicyChanged:
		return i.appendChangePasswordLockoutPolicyEvent(event)
	case OrgIAMPolicyAdded:
		return i.appendAddOrgIAMPolicyEvent(event)
	case OrgIAMPolicyChanged:
		return i.appendChangeOrgIAMPolicyEvent(event)
	}

	return err
}

func (i *IAM) SetData(event *es_models.Event) error {
	i.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, i); err != nil {
		logging.Log("EVEN-9sie4").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-slwi3", "could not unmarshal event")
	}
	return nil
}
