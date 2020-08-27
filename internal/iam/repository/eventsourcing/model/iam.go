package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	IAMVersion = "v1"
)

type IAM struct {
	es_models.ObjectRoot
	SetUpStarted       bool         `json:"-"`
	SetUpDone          bool         `json:"-"`
	GlobalOrgID        string       `json:"globalOrgId,omitempty"`
	IAMProjectID       string       `json:"iamProjectId,omitempty"`
	Members            []*IAMMember `json:"-"`
	IDPs               []*IDPConfig `json:"-"`
	DefaultLoginPolicy *LoginPolicy `json:"-"`
}

func IAMFromModel(iam *model.IAM) *IAM {
	members := IAMMembersFromModel(iam.Members)
	idps := IDPConfigsFromModel(iam.IDPs)
	converted := &IAM{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: iam.SetUpStarted,
		SetUpDone:    iam.SetUpDone,
		GlobalOrgID:  iam.GlobalOrgID,
		IAMProjectID: iam.IAMProjectID,
		Members:      members,
		IDPs:         idps,
	}
	if iam.DefaultLoginPolicy != nil {
		converted.DefaultLoginPolicy = LoginPolicyFromModel(iam.DefaultLoginPolicy)
	}
	return converted
}

func IAMToModel(iam *IAM) *model.IAM {
	members := IAMMembersToModel(iam.Members)
	idps := IDPConfigsToModel(iam.IDPs)
	converted := &model.IAM{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: iam.SetUpStarted,
		SetUpDone:    iam.SetUpDone,
		GlobalOrgID:  iam.GlobalOrgID,
		IAMProjectID: iam.IAMProjectID,
		Members:      members,
		IDPs:         idps,
	}
	if iam.DefaultLoginPolicy != nil {
		converted.DefaultLoginPolicy = LoginPolicyToModel(iam.DefaultLoginPolicy)
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
		i.SetUpStarted = true
	case IAMSetupDone:
		i.SetUpDone = true
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
