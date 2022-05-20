package policy

import (
	timestamp_pb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
)

func ModelLoginPolicyToPb(policy *query.LoginPolicy) *policy_pb.LoginPolicy {
	return &policy_pb.LoginPolicy{
		IsDefault:              policy.IsDefault,
		AllowUsernamePassword:  policy.AllowUsernamePassword,
		AllowRegister:          policy.AllowRegister,
		AllowExternalIdp:       policy.AllowExternalIDPs,
		ForceMfa:               policy.ForceMFA,
		PasswordlessType:       ModelPasswordlessTypeToPb(policy.PasswordlessType),
		HidePasswordReset:      policy.HidePasswordReset,
		IgnoreUnknownUsernames: policy.IgnoreUnknownUsernames,
		DefaultRedirectUri:     policy.DefaultRedirectURI,
		Details: &object.ObjectDetails{
			Sequence:      policy.Sequence,
			CreationDate:  timestamp_pb.New(policy.CreationDate),
			ChangeDate:    timestamp_pb.New(policy.ChangeDate),
			ResourceOwner: policy.OrgID,
		},
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

func ModelPasswordlessTypeToPb(passwordlessType domain.PasswordlessType) policy_pb.PasswordlessType {
	switch passwordlessType {
	case domain.PasswordlessTypeAllowed:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_ALLOWED
	case domain.PasswordlessTypeNotAllowed:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED
	default:
		return policy_pb.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED
	}
}
