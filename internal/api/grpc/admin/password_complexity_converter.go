package admin

import (
	"github.com/caos/zitadel/internal/v2/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func UpdateDefaultPasswordComplexityPolicyToDomain(req *admin_pb.UpdateDefaultPasswordComplexityPolicyRequest) *domain.PasswordComplexityPolicy {
	return &domain.PasswordComplexityPolicy{
		MinLength:    uint64(req.MinLength),
		HasLowercase: req.HasLowercase,
		HasUppercase: req.HasUppercase,
		HasNumber:    req.HasNumber,
		HasSymbol:    req.HasSymbol,
	}
}
