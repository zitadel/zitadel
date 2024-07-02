package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Org struct {
	models.ObjectRoot

	State OrgState
	Name  string

	PrimaryDomain string
	Domains       []*OrgDomain
}

func (o *Org) IsValid() bool {
	if o == nil {
		return false
	}
	o.Name = strings.TrimSpace(o.Name)
	return o.Name != ""
}

func (o *Org) AddIAMDomain(iamDomain string) {
	orgDomain, _ := NewIAMDomainName(o.Name, iamDomain)
	o.Domains = append(o.Domains, &OrgDomain{Domain: orgDomain, Verified: true, Primary: true})
}

type OrgState int32

const (
	OrgStateUnspecified OrgState = iota
	OrgStateActive
	OrgStateInactive
	OrgStateRemoved

	orgStateMax
)

func (s OrgState) Valid() bool {
	return s > OrgStateUnspecified && s < orgStateMax
}
