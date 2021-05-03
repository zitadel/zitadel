package management

import (
	"github.com/caos/zitadel/internal/domain"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func addLabelPolicyToDomain(p *mgmt_pb.AddCustomLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:        p.PrimaryColor,
		SecondaryColor:      p.SecondaryColor,
		WarnColor:           p.WarnColor,
		PrimaryColorDark:    p.PrimaryColorDark,
		SecondaryColorDark:  p.SecondaryColorDark,
		WarnColorDark:       p.WarnColorDark,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
		ErrorMsgPopup:       p.ErrorMsgPopup,
		DisableWatermark:    p.DisableWatermark,
	}
}

func updateLabelPolicyToDomain(p *mgmt_pb.UpdateCustomLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:        p.PrimaryColor,
		SecondaryColor:      p.SecondaryColor,
		WarnColor:           p.WarnColor,
		PrimaryColorDark:    p.PrimaryColorDark,
		SecondaryColorDark:  p.SecondaryColorDark,
		WarnColorDark:       p.WarnColorDark,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
		ErrorMsgPopup:       p.ErrorMsgPopup,
		DisableWatermark:    p.DisableWatermark,
	}
}
