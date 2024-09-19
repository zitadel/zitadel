package admin

import (
	"github.com/zitadel/zitadel/v2/internal/domain"
	admin_pb "github.com/zitadel/zitadel/v2/pkg/grpc/admin"
)

func UpdatePasswordComplexityPolicyToDomain(req *admin_pb.UpdatePasswordComplexityPolicyRequest) *domain.PasswordComplexityPolicy {
	return &domain.PasswordComplexityPolicy{
		MinLength:    uint64(req.MinLength),
		HasLowercase: req.HasLowercase,
		HasUppercase: req.HasUppercase,
		HasNumber:    req.HasNumber,
		HasSymbol:    req.HasSymbol,
	}
}
