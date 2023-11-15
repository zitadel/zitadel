package admin

import (
	"github.com/zitadel/zitadel/v2/internal/domain"
	admin_pb "github.com/zitadel/zitadel/v2/pkg/grpc/admin"
)

func UpdatePasswordAgePolicyToDomain(policy *admin_pb.UpdatePasswordAgePolicyRequest) *domain.PasswordAgePolicy {
	return &domain.PasswordAgePolicy{
		MaxAgeDays:     uint64(policy.MaxAgeDays),
		ExpireWarnDays: uint64(policy.ExpireWarnDays),
	}
}
