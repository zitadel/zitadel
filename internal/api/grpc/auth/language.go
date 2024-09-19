package auth

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/i18n"
	auth_pb "github.com/zitadel/zitadel/v2/pkg/grpc/auth"
)

func (s *Server) GetSupportedLanguages(context.Context, *auth_pb.GetSupportedLanguagesRequest) (*auth_pb.GetSupportedLanguagesResponse, error) {
	return &auth_pb.GetSupportedLanguagesResponse{Languages: domain.LanguagesToStrings(i18n.SupportedLanguages())}, nil
}
