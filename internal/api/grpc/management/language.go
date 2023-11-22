package management

import (
	"context"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetSupportedLanguages(context.Context, *mgmt_pb.GetSupportedLanguagesRequest) (*mgmt_pb.GetSupportedLanguagesResponse, error) {
	return &mgmt_pb.GetSupportedLanguagesResponse{Languages: domain.LanguagesToStrings(i18n.SupportedLanguages())}, nil
}

func (s *Server) GetAllowedLanguages(ctx context.Context, _ *mgmt_pb.GetAllowedLanguagesRequest) (*mgmt_pb.GetAllowedLanguagesResponse, error) {
	restrictions, err := s.query.GetInstanceRestrictions(ctx)
	if err != nil {
		return nil, err
	}
	allowed := restrictions.AllowedLanguages
	if len(allowed) == 0 {
		allowed = i18n.SupportedLanguages()
	}
	return &mgmt_pb.GetAllowedLanguagesResponse{Languages: domain.LanguagesToStrings(allowed)}, nil
}
