package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetLanguages(ctx context.Context, req *mgmt_pb.GetLanguagesRequest) (*mgmt_pb.GetLanguagesResponse, error) {
	langs, err := s.org.Languages(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetLanguagesResponse{Languages: langs}, nil
}
