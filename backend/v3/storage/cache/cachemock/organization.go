package cachemock

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/cache"
)

type OrganizationCacheMock struct {
	c map[string]*domain.Organization
}

func NewOrganizationCacheMock() *OrganizationCacheMock {
	return &OrganizationCacheMock{
		c: make(map[string]*domain.Organization),
	}
}

// Delete implements cache.Cache.
func (o *OrganizationCacheMock) Delete(_ context.Context, _ domain.OrgCacheIndex, keys ...string) error {
	for _, key := range keys {
		delete(o.c, key)
	}
	return nil
}

// Get implements cache.Cache.
func (o *OrganizationCacheMock) Get(_ context.Context, _ domain.OrgCacheIndex, key string) (*domain.Organization, bool) {
	res, ok := o.c[key]
	return res, ok
}

// Invalidate implements cache.Cache.
func (o *OrganizationCacheMock) Invalidate(ctx context.Context, i domain.OrgCacheIndex, key ...string) error {
	return o.Delete(ctx, i, key...)
}

// Set implements cache.Cache.
func (o *OrganizationCacheMock) Set(_ context.Context, organization *domain.Organization) {
	o.c[organization.ID] = organization
}

// Truncate implements cache.Cache.
func (o *OrganizationCacheMock) Truncate(_ context.Context) error {
	o.c = make(map[string]*domain.Organization)
	return nil
}

var _ cache.Cache[domain.OrgCacheIndex, string, *domain.Organization] = (*OrganizationCacheMock)(nil)
