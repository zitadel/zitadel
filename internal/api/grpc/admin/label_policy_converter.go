package admin

import (
	"github.com/caos/zitadel/internal/domain"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func updateLabelPolicyToDomain(policy *admin_pb.UpdateLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:        policy.PrimaryColor,
		SecondaryColor:      policy.SecondaryColor,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
	}
}
