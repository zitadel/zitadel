package admin

import (
	"context"

	instance_grpc "github.com/zitadel/zitadel/internal/api/grpc/instance"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetInstanceDomains(ctx context.Context, req *admin_pb.ListInstanceDomainsRequest) (*admin_pb.ListInstanceDomainsResponse, error) {
	queries, err := ListInstanceDomainsRequestToModel(req)
	if err != nil {
		return nil, err
	}

	domains, err := s.query.SearchInstanceDomains(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListInstanceDomainsResponse{
		Result: instance_grpc.DomainsToPb(domains.Domains),
		Details: object.ToListDetails(
			domains.Count,
			domains.Sequence,
			domains.Timestamp,
		),
	}, nil
}
