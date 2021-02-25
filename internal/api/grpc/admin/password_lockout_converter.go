package admin

import (
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func UpdateDefaultPasswordLockoutPolicyToDomain(p *admin.UpdateDefaultPasswordLockoutPolicyRequest) *domain.PasswordLockoutPolicy {
	return &domain.PasswordLockoutPolicy{
		MaxAttempts:         uint64(p.MaxAttempts),
		ShowLockOutFailures: p.ShowLockoutFailure,
	}
}
