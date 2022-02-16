package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/text"
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	"golang.org/x/text/language"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *admin_pb.GetSupportedLanguagesRequest) (*admin_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.query.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
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

func (s *Server) GetDefaultLanguage(ctx context.Context, req *admin_pb.GetDefaultLanguageRequest) (*admin_pb.GetDefaultLanguageResponse, error) {
	iam, err := s.query.IAMByID(ctx, domain.IAMID)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultLanguageResponse{Language: iam.DefaultLanguage.String()}, nil
}
