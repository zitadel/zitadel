package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ListSMTPConfigs(ctx context.Context, req *admin_pb.ListSMTPConfigsRequest) (*admin_pb.ListSMTPConfigsResponse, error) {
	queries, err := listSMTPConfigsToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchSMTPConfigs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListSMTPConfigsResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.Timestamp),
		Result:  SMTPConfigsToPb(result.Configs),
	}, nil
}
