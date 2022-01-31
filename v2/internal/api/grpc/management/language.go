package management

import (
	"context"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
	"github.com/caos/zitadel/v2/internal/api/grpc/text"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *mgmt_pb.GetSupportedLanguagesRequest) (*mgmt_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.query.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
