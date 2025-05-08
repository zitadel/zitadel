package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
)

type orgDomain struct {
	*org
}

// AddDomain implements [domain.DomainRepository].
func (o *orgDomain) AddDomain(ctx context.Context, domain string) error {
	panic("unimplemented")
}

// RemoveDomain implements [domain.DomainRepository].
func (o *orgDomain) RemoveDomain(ctx context.Context, domain string) error {
	panic("unimplemented")
}

// SetDomainVerified implements [domain.DomainRepository].
func (o *orgDomain) SetDomainVerified(ctx context.Context, domain string) error {
	panic("unimplemented")
}

var _ domain.DomainRepository = (*orgDomain)(nil)
