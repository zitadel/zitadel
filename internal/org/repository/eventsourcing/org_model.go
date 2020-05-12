package eventsourcing

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
)

const (
	orgVersion = "v1"
)

type Org struct {
	es_models.ObjectRoot `json:"-"`

	Name   string `json:"name,omitempty"`
	Domain string `json:"domain,omitempty"`
	State  int32  `json:"-"`

	Members []*OrgMember `json:"-"`
}

func OrgFromModel(org *org_model.Org) *Org {
	members := OrgMembersFromModel(org.Members)

	return &Org{
		ObjectRoot: org.ObjectRoot,
		Domain:     org.Domain,
		Name:       org.Name,
		State:      int32(org.State),
		Members:    members,
	}
}

func OrgToModel(org *Org) *org_model.Org {
	return &org_model.Org{
		ObjectRoot: org.ObjectRoot,
		Domain:     org.Domain,
		Name:       org.Name,
		State:      org_model.OrgState(org.State),
		Members:    OrgMembersToModel(org.Members),
	}
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

func (o *Org) AppendEvent(event *es_models.Event) error {
	switch event.Type {
	case org_model.OrgAdded:
		*o = Org{}
		fallthrough
	case org_model.OrgChanged:
		err := json.Unmarshal(event.Data, o)
		if err != nil {
			return errors.ThrowInternal(err, "EVENT-BpbQZ", "unable to unmarshal event")
		}
	case org_model.OrgDeactivated:
		o.State = int32(org_model.ORGSTATE_INACTIVE)
	case org_model.OrgReactivated:
		o.State = int32(org_model.ORGSTATE_ACTIVE)
	case org_model.OrgMemberAdded:
		member, err := OrgMemberFromEvent(nil, event)
		if err != nil {
			return err
		}
		member.CreationDate = event.CreationDate

		o.setMember(member)
	case org_model.OrgMemberChanged:
		member, err := OrgMemberFromEvent(nil, event)
		if err != nil {
			return err
		}
		existingMember := o.getMember(member.UserID)
		member.CreationDate = existingMember.CreationDate

		o.setMember(member)
	case org_model.OrgMemberRemoved:
		member, err := OrgMemberFromEvent(nil, event)
		if err != nil {
			return err
		}
		o.removeMember(member.UserID)
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
	if changed.Domain != "" && changed.Domain != o.Domain {
		changes["domain"] = changed.Domain
	}

	return changes
}
