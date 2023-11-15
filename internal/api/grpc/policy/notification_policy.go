package policy

import (
	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	"github.com/zitadel/zitadel/v2/internal/query"
	policy_pb "github.com/zitadel/zitadel/v2/pkg/grpc/policy"
)

func ModelNotificationPolicyToPb(policy *query.NotificationPolicy) *policy_pb.NotificationPolicy {
	return &policy_pb.NotificationPolicy{
		IsDefault:      policy.IsDefault,
		PasswordChange: policy.PasswordChange,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}
