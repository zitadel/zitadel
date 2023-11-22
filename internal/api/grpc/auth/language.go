package auth

import (
	"context"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"

	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) GetSupportedLanguages(context.Context, *auth_pb.GetSupportedLanguagesRequest) (*auth_pb.GetSupportedLanguagesResponse, error) {
	return &auth_pb.GetSupportedLanguagesResponse{Languages: domain.LanguagesToStrings(i18n.SupportedLanguages())}, nil
}

func (s *Server) GetAllowedLanguages(ctx context.Context, _ *auth_pb.GetAllowedLanguagesRequest) (*auth_pb.GetAllowedLanguagesResponse, error) {
	restrictions, err := s.query.GetInstanceRestrictions(ctx)
	if err != nil {
		return nil, err
	}
	allowed := restrictions.AllowedLanguages
	if len(allowed) == 0 {
		allowed = i18n.SupportedLanguages()
	}
	return &auth_pb.GetAllowedLanguagesResponse{Languages: domain.LanguagesToStrings(allowed)}, nil
}
