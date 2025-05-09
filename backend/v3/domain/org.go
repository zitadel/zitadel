package domain

import (
	"context"
	"time"
)

type Org struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type OrgRepository interface {
	ByID(ctx context.Context, orgID string) (*Org, error)
	Create(ctx context.Context, org *Org) error
	On(id string) OrgOperation
}

type OrgOperation interface {
	AdminRepository
	DomainRepository
	Update(ctx context.Context, org *Org) error
	Delete(ctx context.Context) error
}

type AdminRepository interface {
	AddAdmin(ctx context.Context, userID string, roles []string) error
	SetAdminRoles(ctx context.Context, userID string, roles []string) error
	RemoveAdmin(ctx context.Context, userID string) error
}

type DomainRepository interface {
	AddDomain(ctx context.Context, domain string) error
	SetDomainVerified(ctx context.Context, domain string) error
	RemoveDomain(ctx context.Context, domain string) error
}
