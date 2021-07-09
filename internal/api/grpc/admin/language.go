package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/text"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetLanguages(ctx context.Context, req *admin_pb.GetLanguagesRequest) (*admin_pb.GetLanguagesResponse, error) {
	langs, err := s.iam.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
