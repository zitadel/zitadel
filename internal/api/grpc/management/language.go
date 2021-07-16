package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/text"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *mgmt_pb.GetSupportedLanguagesRequest) (*mgmt_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.org.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
