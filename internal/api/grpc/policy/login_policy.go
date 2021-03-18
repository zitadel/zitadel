package policy

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func ModelLoginPolicyToPb(policy *model.LoginPolicyView) *policy_pb.LoginPolicy {
	return &policy_pb.LoginPolicy{
		IsDefault:             policy.Default,
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowRegister:         policy.AllowRegister,
		AllowExternalIdp:      policy.AllowRegister,
		ForceMfa:              policy.ForceMFA,
		PasswordlessType:      ModelPasswordlessTypeToPb(policy.PasswordlessType),
	}
}

func PasswordlessTypeToDomain(passwordlessType policy_pb.PasswordlessType) domain.PasswordlessType {
	switch passwordlessType {
	case policy_pb.PasswordlessType_PASSWORDLESS_TYPE_ALLOWED:
		return domain.PasswordlessTypeAllowed
	case policy_pb.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED:
		return domain.PasswordlessTypeNotAllowed
	default:
		return -1
	}
}

func ModelPasswordlessTypeToPb(passwordlessType model.PasswordlessType) policy_pb.PasswordlessType {
	switch passwordlessType {
	case model.PasswordlessTypeAllowed:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_ALLOWED
	case model.PasswordlessTypeNotAllowed:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED
	default:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED
	}
}
