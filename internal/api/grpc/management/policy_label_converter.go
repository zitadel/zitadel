package management

import (
	"github.com/caos/zitadel/internal/domain"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func addLabelPolicyToDomain(p *mgmt_pb.AddCustomLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:        p.PrimaryColor,
		SecondaryColor:      p.SecondaryColor,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
	}
}

func updateLabelPolicyToDomain(p *mgmt_pb.UpdateCustomLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:        p.PrimaryColor,
		SecondaryColor:      p.SecondaryColor,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
	}
}
