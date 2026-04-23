package group

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata/v2"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	filter_pb "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func (s *Server) SetGroupMetadata(ctx context.Context, req *connect.Request[group_v2.SetGroupMetadataRequest]) (*connect.Response[group_v2.SetGroupMetadataResponse], error) {
	// resourceOwner is empty here — the command loads the group and resolves the
	// owning org from its write model. Tenant isolation is enforced via the
	// permission check inside the command.
	details, err := s.command.BulkSetGroupMetadata(ctx, req.Msg.GetId(), "", bulkSetGroupMetadataToDomain(req.Msg)...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.SetGroupMetadataResponse{
		SetDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) ListGroupMetadata(ctx context.Context, req *connect.Request[group_v2.ListGroupMetadataRequest]) (*connect.Response[group_v2.ListGroupMetadataResponse], error) {
	queries, err := listGroupMetadataToQuery(s.systemDefaults, req.Msg)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchGroupMetadata(ctx, true, req.Msg.GetId(), queries, false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.ListGroupMetadataResponse{
		Metadata: metadata.GroupMetadataListToPb(res.Metadata),
		Pagination: &filter_pb.PaginationResponse{
			TotalResult:  res.Count,
			AppliedLimit: uint64(req.Msg.GetPagination().GetLimit()),
		},
	}), nil
}

func (s *Server) DeleteGroupMetadata(ctx context.Context, req *connect.Request[group_v2.DeleteGroupMetadataRequest]) (*connect.Response[group_v2.DeleteGroupMetadataResponse], error) {
	details, err := s.command.BulkRemoveGroupMetadata(ctx, req.Msg.GetId(), "", req.Msg.GetKeys()...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.DeleteGroupMetadataResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

func bulkSetGroupMetadataToDomain(req *group_v2.SetGroupMetadataRequest) []*domain.Metadata {
	out := make([]*domain.Metadata, len(req.Metadata))
	for i, m := range req.Metadata {
		out[i] = &domain.Metadata{
			Key:   m.GetKey(),
			Value: m.GetValue(),
		}
	}
	return out
}

func listGroupMetadataToQuery(defaults systemdefaults.SystemDefaults, req *group_v2.ListGroupMetadataRequest) (*query.GroupMetadataSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(defaults, req.GetPagination())
	if err != nil {
		return nil, err
	}
	queries, err := metadata.GroupMetadataQueriesToQuery(req.GetFilters())
	if err != nil {
		return nil, err
	}
	return &query.GroupMetadataSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}
