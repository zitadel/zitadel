package admin

import (
	"github.com/caos/zitadel/internal/v2/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func updateDefaultLabelPolicyToDomain(policy *admin_pb.UpdateDefaultLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecondaryColor: policy.SecondaryColor,
	}
}
