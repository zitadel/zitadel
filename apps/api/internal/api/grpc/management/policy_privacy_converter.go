package management

import (
	"github.com/zitadel/zitadel/internal/domain"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func AddPrivacyPolicyToDomain(req *mgmt_pb.AddCustomPrivacyPolicyRequest) *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		TOSLink:        req.TosLink,
		PrivacyLink:    req.PrivacyLink,
		HelpLink:       req.HelpLink,
		SupportEmail:   domain.EmailAddress(req.SupportEmail),
		DocsLink:       req.DocsLink,
		CustomLink:     req.CustomLink,
		CustomLinkText: req.CustomLinkText,
	}
}

func UpdatePrivacyPolicyToDomain(req *mgmt_pb.UpdateCustomPrivacyPolicyRequest) *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		TOSLink:        req.TosLink,
		PrivacyLink:    req.PrivacyLink,
		HelpLink:       req.HelpLink,
		SupportEmail:   domain.EmailAddress(req.SupportEmail),
		DocsLink:       req.DocsLink,
		CustomLink:     req.CustomLink,
		CustomLinkText: req.CustomLinkText,
	}
}
