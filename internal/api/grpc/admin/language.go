package admin

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/text"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *admin_pb.GetSupportedLanguagesRequest) (*admin_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.query.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}

func (s *Server) SetDefaultLanguage(ctx context.Context, req *admin_pb.SetDefaultLanguageRequest) (*admin_pb.SetDefaultLanguageResponse, error) {
	_, err := language.Parse(req.Language)
	if err != nil {
		return nil, caos_errors.ThrowInvalidArgument(err, "API-39nnf", "Errors.Language.Parse")
	}
	//TODO: Will be added by silvan
	//details, err := s.command.SetDefaultLanguage(ctx, lang)
	//if err != nil {
	//	return nil, err
	//}
	//return &admin_pb.SetDefaultLanguageResponse{
	//	Details: object.DomainToChangeDetailsPb(details),
	//}, nil
	return nil, nil
}

func (s *Server) GetDefaultLanguage(ctx context.Context, _ *admin_pb.GetDefaultLanguageRequest) (*admin_pb.GetDefaultLanguageResponse, error) {
	return &admin_pb.GetDefaultLanguageResponse{Language: authz.GetInstance(ctx).DefaultLanguage().String()}, nil
}
