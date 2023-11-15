package admin

import (
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/pkg/grpc/admin"
)

func UpdateLockoutPolicyToDomain(p *admin.UpdateLockoutPolicyRequest) *domain.LockoutPolicy {
	return &domain.LockoutPolicy{
		MaxPasswordAttempts: uint64(p.MaxPasswordAttempts),
	}
}
