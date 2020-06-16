package model

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type OrgDomain struct {
	es_models.ObjectRoot
	Domain   string
	Primary  bool
	Verified bool
}

func NewOrgDomain(orgID, domain string) *OrgDomain {
	return &OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: orgID}, Domain: domain}
}

func (domain *OrgDomain) IsValid() bool {
	return domain.AggregateID != "" && domain.Domain != ""
}
