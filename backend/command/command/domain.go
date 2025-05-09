package command

import (
	"slices"

	"github.com/zitadel/zitadel/backend/command/receiver"
)

type SetPrimaryDomain struct {
	Domains []*receiver.Domain

	Domain string
}

func (s *SetPrimaryDomain) Execute() error {
	for domain := range slices.Values(s.Domains) {
		domain.IsPrimary = domain.Name == s.Domain
	}
	return nil
}

type RemoveDomain struct {
	Domains []*receiver.Domain

	Domain string
}

func (r *RemoveDomain) Execute() error {
	r.Domains = slices.DeleteFunc(r.Domains, func(domain *receiver.Domain) bool {
		return domain.Name == r.Domain
	})
	return nil
}
