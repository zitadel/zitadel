package admin

import (
	"context"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/i18n"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *admin_pb.GetSupportedLanguagesRequest) (*admin_pb.GetSupportedLanguagesResponse, error) {
	return &admin_pb.GetSupportedLanguagesResponse{Languages: domain.LanguagesToStrings(i18n.SupportedLanguages())}, nil
}

func (s *Server) GetAllowedLanguages(ctx context.Context, _ *admin_pb.GetAllowedLanguagesRequest) (*admin_pb.GetAllowedLanguagesResponse, error) {
	restrictions, err := s.query.GetInstanceRestrictions(ctx)
	if err != nil {
		return nil, err
	}
	allowed := restrictions.AllowedLanguages
	if len(allowed) == 0 {
		allowed = i18n.SupportedLanguages()
	}
	return &admin_pb.GetAllowedLanguagesResponse{Languages: domain.LanguagesToStrings(allowed)}, nil
}

func (s *Server) SetDefaultLanguage(ctx context.Context, req *admin_pb.SetDefaultLanguageRequest) (*admin_pb.SetDefaultLanguageResponse, error) {
	lang, err := language.Parse(req.Language)
	if err != nil {
		return nil, caos_errors.ThrowInvalidArgument(err, "API-39nnf", "Errors.Language.Parse")
	}
	details, err := s.command.SetDefaultLanguage(ctx, lang)
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetDefaultLanguageResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) GetDefaultLanguage(ctx context.Context, _ *admin_pb.GetDefaultLanguageRequest) (*admin_pb.GetDefaultLanguageResponse, error) {
	return &admin_pb.GetDefaultLanguageResponse{Language: authz.GetInstance(ctx).DefaultLanguage().String()}, nil
}
