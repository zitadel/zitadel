package management

import (
	"github.com/caos/zitadel/internal/domain"
	mgmt "github.com/caos/zitadel/pkg/grpc/management"
)

func AddLockoutPolicyToDomain(p *mgmt.AddCustomLockoutPolicyRequest) *domain.LockoutPolicy {
	return &domain.LockoutPolicy{
		MaxPasswordAttempts: uint64(p.MaxPasswordAttempts),
		ShowLockOutFailures: p.ShowLockoutFailure,
	}
}

func UpdateLockoutPolicyToDomain(p *mgmt.UpdateCustomLockoutPolicyRequest) *domain.LockoutPolicy {
	return &domain.LockoutPolicy{
		MaxPasswordAttempts: uint64(p.MaxPasswordAttempts),
		ShowLockOutFailures: p.ShowLockoutFailure,
	}
}
