package admin

import (
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func updateLoginPolicyToDomain(p *admin_pb.UpdateLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword: p.AllowUsernamePassword,
		AllowRegister:         p.AllowRegister,
		AllowExternalIDP:      p.AllowExternalIdp,
		ForceMFA:              p.ForceMfa,
		PasswordlessType:      policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
	}
}
