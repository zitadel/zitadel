package admin

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func UpdatePasswordLockoutPolicyToDomain(p *admin.UpdatePasswordLockoutPolicyRequest) *domain.PasswordLockoutPolicy {
	return &domain.PasswordLockoutPolicy{
		MaxAttempts:         uint64(p.MaxAttempts),
		ShowLockOutFailures: p.ShowLockoutFailure,
	}
}
