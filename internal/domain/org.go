package domain

import (
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Org struct {
	models.ObjectRoot

	State OrgState
	Name  string

	PrimaryDomain string
	Domains       []*OrgDomain
}

func (o *Org) IsValid() bool {
	return o != nil && o.Name != ""
}

func (o *Org) AddIAMDomain(iamDomain string) {
	o.Domains = append(o.Domains, &OrgDomain{Domain: NewIAMDomainName(o.Name, iamDomain), Verified: true, Primary: true})
}

type OrgState int32

const (
	OrgStateUnspecified OrgState = iota
	OrgStateActive
	OrgStateInactive
	OrgStateRemoved
)
