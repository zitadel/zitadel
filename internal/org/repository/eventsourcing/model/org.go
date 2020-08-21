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

	Domains      []*OrgDomain              `json:"-"`
	Members      []*OrgMember              `json:"-"`
	OrgIamPolicy *OrgIamPolicy             `json:"-"`
	LoginPolicy  *iam_es_model.LoginPolicy `json:"-"`
	IDPs         []*iam_es_model.IdpConfig `json:"-"`
}

func OrgFromModel(org *org_model.Org) *Org {
	members := OrgMembersFromModel(org.Members)
	domains := OrgDomainsFromModel(org.Domains)
	idps := iam_es_model.IdpConfigsFromModel(org.IDPs)
	converted := &Org{
		ObjectRoot: org.ObjectRoot,
		Name:       org.Name,
		State:      int32(org.State),
		Domains:    domains,
		Members:    members,
		IDPs:       idps,
	}
	if org.OrgIamPolicy != nil {
		converted.OrgIamPolicy = OrgIamPolicyFromModel(org.OrgIamPolicy)
	}
	if org.LoginPolicy != nil {
		converted.LoginPolicy = iam_es_model.LoginPolicyFromModel(org.LoginPolicy)
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
		IDPs:       iam_es_model.IdpConfigsToModel(org.IDPs),
	}
	if org.OrgIamPolicy != nil {
		converted.OrgIamPolicy = OrgIamPolicyToModel(org.OrgIamPolicy)
	}
	if org.LoginPolicy != nil {
		converted.LoginPolicy = iam_es_model.LoginPolicyToModel(org.LoginPolicy)
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
	case OrgIamPolicyAdded:
		err = o.appendAddOrgIamPolicyEvent(event)
	case OrgIamPolicyChanged:
		err = o.appendChangeOrgIamPolicyEvent(event)
	case OrgIamPolicyRemoved:
		o.appendRemoveOrgIamPolicyEvent()
	case IdpConfigAdded:
		err = o.appendAddIdpConfigEvent(event)
	case IdpConfigChanged:
		err = o.appendChangeIdpConfigEvent(event)
	case IdpConfigRemoved:
		err = o.appendRemoveIdpConfigEvent(event)
	case IdpConfigDeactivated:
		err = o.appendIdpConfigStateEvent(event, model.IdpConfigStateInactive)
	case IdpConfigReactivated:
		err = o.appendIdpConfigStateEvent(event, model.IdpConfigStateActive)
	case OidcIdpConfigAdded:
		err = o.appendAddOidcIdpConfigEvent(event)
	case OidcIdpConfigChanged:
		err = o.appendChangeOidcIdpConfigEvent(event)
	case LoginPolicyAdded:
		err = o.appendAddLoginPolicyEvent(event)
	case LoginPolicyChanged:
		err = o.appendChangeLoginPolicyEvent(event)
	case LoginPolicyRemoved:
		o.appendRemoveLoginPolicyEvent(event)
	case LoginPolicyIdpProviderAdded:
		err = o.appendAddIdpProviderToLoginPolicyEvent(event)
	case LoginPolicyIdpProviderRemoved:
		err = o.appendRemoveIdpProviderFromLoginPolicyEvent(event)
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
