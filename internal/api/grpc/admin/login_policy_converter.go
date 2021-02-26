package admin

import (
	"github.com/caos/logging"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func updateDefaultLoginPolicyToDomain(p *admin_pb.UpdateDefaultLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword: p.AllowUsernamePassword,
		AllowRegister:         p.AllowRegister,
		AllowExternalIDP:      p.AllowExternalIdp,
		ForceMFA:              p.ForceMfa,
		PasswordlessType:      policy_grpc.PasswordlessTypeToDomain(p.PasswordlessType),
	}
}
