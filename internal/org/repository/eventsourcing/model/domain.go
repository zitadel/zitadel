package model

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/org/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type OrgDomain struct {
	es_models.ObjectRoot `json:"-"`

	Domain         string              `json:"domain"`
	Verified       bool                `json:"-"`
	Primary        bool                `json:"-"`
	ValidationType int32               `json:"validationType"`
	ValidationCode *crypto.CryptoValue `json:"validationCode"`
}

func GetDomain(domains []*OrgDomain, domain string) (int, *OrgDomain) {
	for i, d := range domains {
		if d.Domain == domain {
			return i, d
		}
	}
	return -1, nil
}

func (o *Org) appendAddDomainEvent(event eventstore.Event) error {
	domain := new(OrgDomain)
	err := domain.SetData(event)
	if err != nil {
		return err
	}
	domain.ObjectRoot.CreationDate = event.CreatedAt()
	o.Domains = append(o.Domains, domain)
	return nil
}

func (o *Org) appendRemoveDomainEvent(event eventstore.Event) error {
	domain := new(OrgDomain)
	err := domain.SetData(event)
	if err != nil {
		return err
	}
	if i, r := GetDomain(o.Domains, domain.Domain); r != nil {
		o.Domains[i] = o.Domains[len(o.Domains)-1]
		o.Domains[len(o.Domains)-1] = nil
		o.Domains = o.Domains[:len(o.Domains)-1]
	}
	return nil
}

func (o *Org) appendVerifyDomainEvent(event eventstore.Event) error {
	domain := new(OrgDomain)
	err := domain.SetData(event)
	if err != nil {
		return err
	}
	if i, d := GetDomain(o.Domains, domain.Domain); d != nil {
		d.Verified = true
		o.Domains[i] = d
	}
	return nil
}

func (o *Org) appendPrimaryDomainEvent(event eventstore.Event) error {
	domain := new(OrgDomain)
	err := domain.SetData(event)
	if err != nil {
		return err
	}
	for _, d := range o.Domains {
		d.Primary = false
		if d.Domain == domain.Domain {
			d.Primary = true
		}
	}
	return nil
}

func (o *Org) appendVerificationDomainEvent(event eventstore.Event) error {
	domain := new(OrgDomain)
	err := domain.SetData(event)
	if err != nil {
		return err
	}
	for _, d := range o.Domains {
		if d.Domain == domain.Domain {
			d.ValidationType = domain.ValidationType
			d.ValidationCode = domain.ValidationCode
		}
	}
	return nil
}

func (m *OrgDomain) SetData(event eventstore.Event) error {
	err := event.Unmarshal(m)
	if err != nil {
		return zerrors.ThrowInternal(err, "EVENT-Hz7Mb", "unable to unmarshal data")
	}
	return nil
}

func OrgDomainsFromModel(domains []*model.OrgDomain) []*OrgDomain {
	convertedDomainss := make([]*OrgDomain, len(domains))
	for i, m := range domains {
		convertedDomainss[i] = OrgDomainFromModel(m)
	}
	return convertedDomainss
}

func OrgDomainFromModel(domain *model.OrgDomain) *OrgDomain {
	return &OrgDomain{
		ObjectRoot:     domain.ObjectRoot,
		Domain:         domain.Domain,
		Verified:       domain.Verified,
		Primary:        domain.Primary,
		ValidationType: int32(domain.ValidationType),
		ValidationCode: domain.ValidationCode,
	}
}

func OrgDomainsToModel(domains []*OrgDomain) []*model.OrgDomain {
	convertedDomains := make([]*model.OrgDomain, len(domains))
	for i, m := range domains {
		convertedDomains[i] = OrgDomainToModel(m)
	}
	return convertedDomains
}

func OrgDomainToModel(domain *OrgDomain) *model.OrgDomain {
	return &model.OrgDomain{
		ObjectRoot:     domain.ObjectRoot,
		Domain:         domain.Domain,
		Primary:        domain.Primary,
		Verified:       domain.Verified,
		ValidationType: model.OrgDomainValidationType(domain.ValidationType),
		ValidationCode: domain.ValidationCode,
	}
}
