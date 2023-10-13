package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ListSMTPProviders(ctx context.Context, req *admin_pb.ListSMSProvidersRequest) (*admin_pb.ListSMTPProvidersResponse, error) {
	queries, err := listSMSConfigsToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchSMSConfigs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListSMTPProvidersResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.Timestamp),
		Result:  SMSConfigsToPb(result.Configs),
	}, nil
}
