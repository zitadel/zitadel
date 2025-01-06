package service

import (
	"context"

	"github.com/zitadel/zitadel/backend/internal/port"
)

type InstanceDomainRepository interface {
	// CreateInstanceDomain creates a new instance domain
	CreateInstanceDomain(ctx context.Context, executor port.Executor, instanceID string, domain *Domain) error
	// CreateInstanceDomains creates multiple instance domains
	CreateInstanceDomains(ctx context.Context, executor port.Executor, instanceID string, domains []*Domain) error
}
