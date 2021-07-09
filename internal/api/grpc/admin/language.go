package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/text"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *admin_pb.GetSupportedLanguagesRequest) (*admin_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.iam.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
