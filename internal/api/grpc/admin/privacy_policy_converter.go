package admin

import (
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func UpdatePrivacyPolicyToDomain(req *admin_pb.UpdatePrivacyPolicyRequest) *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		TOSLink:     req.TosLink,
		PrivacyLink: req.PrivacyLink,
	}
}
