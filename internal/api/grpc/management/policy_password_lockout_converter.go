package management

import (
	"github.com/caos/zitadel/internal/domain"
	mgmt "github.com/caos/zitadel/pkg/grpc/management"
)

func AddPasswordLockoutPolicyToDomain(p *mgmt.AddCustomPasswordLockoutPolicyRequest) *domain.PasswordLockoutPolicy {
	return &domain.PasswordLockoutPolicy{
		MaxAttempts:         uint64(p.MaxAttempts),
		ShowLockOutFailures: p.ShowLockoutFailure,
	}
}

func UpdatePasswordLockoutPolicyToDomain(p *mgmt.UpdateCustomPasswordLockoutPolicyRequest) *domain.PasswordLockoutPolicy {
	return &domain.PasswordLockoutPolicy{
		MaxAttempts:         uint64(p.MaxAttempts),
		ShowLockOutFailures: p.ShowLockoutFailure,
	}
}
