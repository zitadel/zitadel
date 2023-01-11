package management

import (
	"github.com/zitadel/zitadel/internal/domain"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func AddNotificationPolicyToDomain(req *mgmt_pb.AddCustomNotificationPolicyRequest) *domain.NotificationPolicy {
	return &domain.NotificationPolicy{
		PasswordChange: req.PasswordChange,
	}
}

func UpdateNotificationPolicyToDomain(req *mgmt_pb.UpdateCustomNotificationPolicyRequest) *domain.NotificationPolicy {
	return &domain.NotificationPolicy{
		PasswordChange: req.PasswordChange,
	}
}
