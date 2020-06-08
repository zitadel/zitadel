package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"strings"
)

type Org struct {
	es_models.ObjectRoot

	State   OrgState
	Name    string
	Domains []*OrgDomain

	Members      []*OrgMember
	OrgIamPolicy *OrgIamPolicy
}

type OrgState int32

const (
	ORGSTATE_ACTIVE OrgState = iota
	ORGSTATE_INACTIVE
)

func NewOrg(id string) *Org {
	return &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: id}, State: ORGSTATE_ACTIVE}
}

func (o *Org) IsActive() bool {
	return o.State == ORGSTATE_ACTIVE
}

func (o *Org) IsValid() bool {
	return o.Name != ""
}

func (o *Org) ContainsDomain(domain *OrgDomain) bool {
	for _, d := range o.Domains {
		if d.Domain == domain.Domain {
			return true
		}
	}
	return false
}

func (o *Org) ContainsMember(userID string) bool {
	for _, member := range o.Members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}

func (o *Org) nameForDomain(iamDomain string) string {
	return strings.ToLower(strings.ReplaceAll(o.Name, " ", "-") + "." + iamDomain)
}

func (o *Org) AddIAMDomain(iamDomain string) {
	o.Domains = append(o.Domains, &OrgDomain{Domain: o.nameForDomain(iamDomain), Verified: true, Primary: true})
}
