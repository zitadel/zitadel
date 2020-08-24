package model

import (
	"encoding/json"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
)

const (
	IamVersion = "v1"
)

type Iam struct {
	es_models.ObjectRoot
	SetUpStarted       bool         `json:"-"`
	SetUpDone          bool         `json:"-"`
	GlobalOrgID        string       `json:"globalOrgId,omitempty"`
	IamProjectID       string       `json:"iamProjectId,omitempty"`
	Members            []*IamMember `json:"-"`
	IDPs               []*IDPConfig `json:"-"`
	DefaultLoginPolicy *LoginPolicy `json:"-"`
}

func IamFromModel(iam *model.Iam) *Iam {
	members := IamMembersFromModel(iam.Members)
	idps := IDPConfigsFromModel(iam.IDPs)
	converted := &Iam{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: iam.SetUpStarted,
		SetUpDone:    iam.SetUpDone,
		GlobalOrgID:  iam.GlobalOrgID,
		IamProjectID: iam.IamProjectID,
		Members:      members,
		IDPs:         idps,
	}
	if iam.DefaultLoginPolicy != nil {
		converted.DefaultLoginPolicy = LoginPolicyFromModel(iam.DefaultLoginPolicy)
	}
	return converted
}

func IamToModel(iam *Iam) *model.Iam {
	members := IamMembersToModel(iam.Members)
	idps := IDPConfigsToModel(iam.IDPs)
	converted := &model.Iam{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: iam.SetUpStarted,
		SetUpDone:    iam.SetUpDone,
		GlobalOrgID:  iam.GlobalOrgID,
		IamProjectID: iam.IamProjectID,
		Members:      members,
		IDPs:         idps,
	}
	if iam.DefaultLoginPolicy != nil {
		converted.DefaultLoginPolicy = LoginPolicyToModel(iam.DefaultLoginPolicy)
	}
	return converted
}

func (i *Iam) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := i.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (i *Iam) AppendEvent(event *es_models.Event) (err error) {
	i.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case IamSetupStarted:
		i.SetUpStarted = true
	case IamSetupDone:
		i.SetUpDone = true
	case IamProjectSet,
		GlobalOrgSet:
		err = i.SetData(event)
	case IamMemberAdded:
		err = i.appendAddMemberEvent(event)
	case IamMemberChanged:
		err = i.appendChangeMemberEvent(event)
	case IamMemberRemoved:
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

func (i *Iam) SetData(event *es_models.Event) error {
	i.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, i); err != nil {
		logging.Log("EVEN-9sie4").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-slwi3", "could not unmarshal event")
	}
	return nil
}
