package management

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/grpc/text"
	mgmt_pb "github.com/zitadel/zitadel/v2/pkg/grpc/management"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *mgmt_pb.GetSupportedLanguagesRequest) (*mgmt_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.query.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
