package model

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
)

const (
	OrgVersion = "v1"
)

type Org struct {
	es_models.ObjectRoot `json:"-"`

	Name  string `json:"name,omitempty"`
	State int32  `json:"-"`

	Domains                  []*OrgDomain                           `json:"-"`
	Members                  []*OrgMember                           `json:"-"`
	OrgIAMPolicy             *iam_es_model.OrgIAMPolicy             `json:"-"`
	LabelPolicy              *iam_es_model.LabelPolicy              `json:"-"`
	MailTemplate             *iam_es_model.MailTemplate             `json:"-"`
	MailTexts                []*iam_es_model.MailText               `json:"-"`
	IDPs                     []*iam_es_model.IDPConfig              `json:"-"`
	LoginPolicy              *iam_es_model.LoginPolicy              `json:"-"`
	PasswordComplexityPolicy *iam_es_model.PasswordComplexityPolicy `json:"-"`
	PasswordAgePolicy        *iam_es_model.PasswordAgePolicy        `json:"-"`
	PasswordLockoutPolicy    *iam_es_model.PasswordLockoutPolicy    `json:"-"`
}

func OrgFromModel(org *org_model.Org) *Org {
	members := OrgMembersFromModel(org.Members)
	domains := OrgDomainsFromModel(org.Domains)
	idps := iam_es_model.IDPConfigsFromModel(org.IDPs)
	mailTexts := iam_es_model.MailTextsFromModel(org.MailTexts)
	converted := &Org{
		ObjectRoot: org.ObjectRoot,
		Name:       org.Name,
		State:      int32(org.State),
		Domains:    domains,
		MailTexts:  mailTexts,
		Members:    members,
		IDPs:       idps,
	}
	if org.OrgIamPolicy != nil {
		converted.OrgIAMPolicy = iam_es_model.OrgIAMPolicyFromModel(org.OrgIamPolicy)
	}
	if org.LoginPolicy != nil {
		converted.LoginPolicy = iam_es_model.LoginPolicyFromModel(org.LoginPolicy)
	}
	if org.LabelPolicy != nil {
		converted.LabelPolicy = iam_es_model.LabelPolicyFromModel(org.LabelPolicy)
	}
	if org.MailTemplate != nil {
		converted.MailTemplate = iam_es_model.MailTemplateFromModel(org.MailTemplate)
	}
	if org.PasswordComplexityPolicy != nil {
		converted.PasswordComplexityPolicy = iam_es_model.PasswordComplexityPolicyFromModel(org.PasswordComplexityPolicy)
	}
	if org.PasswordAgePolicy != nil {
		converted.PasswordAgePolicy = iam_es_model.PasswordAgePolicyFromModel(org.PasswordAgePolicy)
	}
	if org.PasswordLockoutPolicy != nil {
		converted.PasswordLockoutPolicy = iam_es_model.PasswordLockoutPolicyFromModel(org.PasswordLockoutPolicy)
	}
	return converted
}

func OrgToModel(org *Org) *org_model.Org {
	converted := &org_model.Org{
		ObjectRoot: org.ObjectRoot,
		Name:       org.Name,
		State:      org_model.OrgState(org.State),
		Domains:    OrgDomainsToModel(org.Domains),
		Members:    OrgMembersToModel(org.Members),
		MailTexts:  iam_es_model.MailTextsToModel(org.MailTexts),
		IDPs:       iam_es_model.IDPConfigsToModel(org.IDPs),
	}
	if org.OrgIAMPolicy != nil {
		converted.OrgIamPolicy = iam_es_model.OrgIAMPolicyToModel(org.OrgIAMPolicy)
	}
	if org.LoginPolicy != nil {
		converted.LoginPolicy = iam_es_model.LoginPolicyToModel(org.LoginPolicy)
	}
	if org.LabelPolicy != nil {
		converted.LabelPolicy = iam_es_model.LabelPolicyToModel(org.LabelPolicy)
	}
	if org.MailTemplate != nil {
		converted.MailTemplate = iam_es_model.MailTemplateToModel(org.MailTemplate)
	}
	if org.PasswordComplexityPolicy != nil {
		converted.PasswordComplexityPolicy = iam_es_model.PasswordComplexityPolicyToModel(org.PasswordComplexityPolicy)
	}
	if org.PasswordAgePolicy != nil {
		converted.PasswordAgePolicy = iam_es_model.PasswordAgePolicyToModel(org.PasswordAgePolicy)
	}
	if org.PasswordLockoutPolicy != nil {
		converted.PasswordLockoutPolicy = iam_es_model.PasswordLockoutPolicyToModel(org.PasswordLockoutPolicy)
	}
	return converted
}

func OrgFromEvents(org *Org, events ...*es_models.Event) (*Org, error) {
	if org == nil {
		org = new(Org)
	}

	return org, org.AppendEvents(events...)
}

func (o *Org) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		err := o.AppendEvent(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Org) AppendEvent(event *es_models.Event) (err error) {
	switch event.Type {
	case OrgAdded:
		*o = Org{}
		err = o.setData(event)
		if err != nil {
			return err
		}
	case OrgChanged:
		err = o.setData(event)
		if err != nil {
			return err
		}
	case OrgDeactivated:
		o.State = int32(org_model.OrgStateInactive)
	case OrgReactivated:
		o.State = int32(org_model.OrgStateActive)
	case OrgMemberAdded:
		member, err := OrgMemberFromEvent(nil, event)
		if err != nil {
			return err
		}
		member.CreationDate = event.CreationDate

		o.setMember(member)
	case OrgMemberChanged:
		member, err := OrgMemberFromEvent(nil, event)
		if err != nil {
			return err
		}
		existingMember := o.getMember(member.UserID)
		member.CreationDate = existingMember.CreationDate

		o.setMember(member)
	case OrgMemberRemoved:
		member, err := OrgMemberFromEvent(nil, event)
		if err != nil {
			return err
		}
		o.removeMember(member.UserID)
	case OrgDomainAdded:
		err = o.appendAddDomainEvent(event)
	case OrgDomainVerificationAdded:
		err = o.appendVerificationDomainEvent(event)
	case OrgDomainVerified:
		err = o.appendVerifyDomainEvent(event)
	case OrgDomainPrimarySet:
		err = o.appendPrimaryDomainEvent(event)
	case OrgDomainRemoved:
		err = o.appendRemoveDomainEvent(event)
	case OrgIAMPolicyAdded:
		err = o.appendAddOrgIAMPolicyEvent(event)
	case OrgIAMPolicyChanged:
		err = o.appendChangeOrgIAMPolicyEvent(event)
	case OrgIAMPolicyRemoved:
		o.appendRemoveOrgIAMPolicyEvent()
	case IDPConfigAdded:
		err = o.appendAddIDPConfigEvent(event)
	case IDPConfigChanged:
		err = o.appendChangeIDPConfigEvent(event)
	case IDPConfigRemoved:
		err = o.appendRemoveIDPConfigEvent(event)
	case IDPConfigDeactivated:
		err = o.appendIDPConfigStateEvent(event, model.IDPConfigStateInactive)
	case IDPConfigReactivated:
		err = o.appendIDPConfigStateEvent(event, model.IDPConfigStateActive)
	case OIDCIDPConfigAdded:
		err = o.appendAddOIDCIDPConfigEvent(event)
	case OIDCIDPConfigChanged:
		err = o.appendChangeOIDCIDPConfigEvent(event)
	case LabelPolicyAdded:
		err = o.appendAddLabelPolicyEvent(event)
	case LabelPolicyChanged:
		err = o.appendChangeLabelPolicyEvent(event)
	case LabelPolicyRemoved:
		o.appendRemoveLabelPolicyEvent(event)
	case LoginPolicyAdded:
		err = o.appendAddLoginPolicyEvent(event)
	case LoginPolicyChanged:
		err = o.appendChangeLoginPolicyEvent(event)
	case LoginPolicyRemoved:
		o.appendRemoveLoginPolicyEvent(event)
	case LoginPolicyIDPProviderAdded:
		err = o.appendAddIdpProviderToLoginPolicyEvent(event)
	case LoginPolicyIDPProviderRemoved:
		err = o.appendRemoveIdpProviderFromLoginPolicyEvent(event)
	case MailTemplateAdded:
		err = o.appendAddMailTemplateEvent(event)
	case MailTemplateChanged:
		err = o.appendChangeMailTemplateEvent(event)
	case MailTemplateRemoved:
		o.appendRemoveMailTemplateEvent(event)
	case MailTextAdded:
		err = o.appendAddMailTextEvent(event)
	case MailTextChanged:
		err = o.appendChangeMailTextEvent(event)
	case MailTextRemoved:
		o.appendRemoveMailTextEvent(event)
	case LoginPolicySecondFactorAdded:
		err = o.appendAddSecondFactorToLoginPolicyEvent(event)
	case LoginPolicySecondFactorRemoved:
		err = o.appendRemoveSecondFactorFromLoginPolicyEvent(event)
	case LoginPolicyMultiFactorAdded:
		err = o.appendAddMultiFactorToLoginPolicyEvent(event)
	case LoginPolicyMultiFactorRemoved:
		err = o.appendRemoveMultiFactorFromLoginPolicyEvent(event)
	case PasswordComplexityPolicyAdded:
		err = o.appendAddPasswordComplexityPolicyEvent(event)
	case PasswordComplexityPolicyChanged:
		err = o.appendChangePasswordComplexityPolicyEvent(event)
	case PasswordComplexityPolicyRemoved:
		o.appendRemovePasswordComplexityPolicyEvent(event)
	case PasswordAgePolicyAdded:
		err = o.appendAddPasswordAgePolicyEvent(event)
	case PasswordAgePolicyChanged:
		err = o.appendChangePasswordAgePolicyEvent(event)
	case PasswordAgePolicyRemoved:
		o.appendRemovePasswordAgePolicyEvent(event)
	case PasswordLockoutPolicyAdded:
		err = o.appendAddPasswordLockoutPolicyEvent(event)
	case PasswordLockoutPolicyChanged:
		err = o.appendChangePasswordLockoutPolicyEvent(event)
	case PasswordLockoutPolicyRemoved:
		o.appendRemovePasswordLockoutPolicyEvent(event)
	}
	if err != nil {
		return err
	}
	o.ObjectRoot.AppendEvent(event)
	return nil
}

func (o *Org) setData(event *es_models.Event) error {
	err := json.Unmarshal(event.Data, o)
	if err != nil {
		return errors.ThrowInternal(err, "EVENT-BpbQZ", "unable to unmarshal event")
	}
	return nil
}

func (o *Org) getMember(userID string) *OrgMember {
	for _, member := range o.Members {
		if member.UserID == userID {
			return member
		}
	}
	return nil
}

func (o *Org) setMember(member *OrgMember) {
	for i, existingMember := range o.Members {
		if existingMember.UserID == member.UserID {
			o.Members[i] = member
			return
		}
	}
	o.Members = append(o.Members, member)
}

func (o *Org) removeMember(userID string) {
	for i := len(o.Members) - 1; i >= 0; i-- {
		if o.Members[i].UserID == userID {
			copy(o.Members[i:], o.Members[i+1:])
			o.Members[len(o.Members)-1] = nil
			o.Members = o.Members[:len(o.Members)-1]
		}
	}
}

func (o *Org) Changes(changed *Org) map[string]interface{} {
	changes := make(map[string]interface{}, 2)

	if changed.Name != "" && changed.Name != o.Name {
		changes["name"] = changed.Name
	}

	return changes
}
