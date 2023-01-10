package admin

import (
	"github.com/zitadel/zitadel/internal/domain"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func UpdateNotificationPolicyToDomain(req *admin_pb.UpdateNotificationPolicyRequest) *domain.NotificationPolicy {
	return &domain.NotificationPolicy{
		PasswordChange: req.PasswordChange,
	}
}
