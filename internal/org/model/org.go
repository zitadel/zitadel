package model

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
)

type Org struct {
	es_models.ObjectRoot

	State   OrgState
	Name    string
	Domains []*OrgDomain

	DomainPolicy *iam_model.DomainPolicy
}

type OrgState int32

const (
	OrgStateActive OrgState = iota
	OrgStateInactive
)

func (o *Org) IsActive() bool {
	return o.State == OrgStateActive
}

func (o *Org) GetDomain(domain *OrgDomain) (int, *OrgDomain) {
	for i, d := range o.Domains {
		if d.Domain == domain.Domain {
			return i, d
		}
	}
	return -1, nil
}

func (o *Org) GetPrimaryDomain() (string, error) {
	for _, d := range o.Domains {
		if d.Primary {
			return d.Domain, nil
		}
	}
	return "", errors.ThrowInternalf(nil, "ORG-Dertg", "no primary domain found for org: %s (instanceID: %s)", o.AggregateID, o.InstanceID)
}

func (o *Org) AddIAMDomain(iamDomain string) {
	orgDomain, _ := domain.NewIAMDomainName(o.Name, iamDomain)
	o.Domains = append(o.Domains, &OrgDomain{Domain: orgDomain, Verified: true, Primary: true})
}
