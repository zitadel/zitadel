package user

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) ListUserMetadata(ctx context.Context, req *connect.Request[user.ListUserMetadataRequest]) (*connect.Response[user.ListUserMetadataResponse], error) {
	metadataQueries, err := s.listUserMetadataRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchUserMetadata(ctx, true, req.Msg.UserId, metadataQueries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.ListUserMetadataResponse{
		Metadata:   metadata.UserMetadataListToPb(res.Metadata),
		Pagination: filter.QueryToPaginationPb(metadataQueries.SearchRequest, res.SearchResponse),
	}), nil
}

func (s *Server) listUserMetadataRequestToModel(req *user.ListUserMetadataRequest) (*query.UserMetadataSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := metadata.UserMetadataFiltersToQuery(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.UserMetadataSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: query.UserMetadataCreationDateCol,
		},
		Queries: queries,
	}, nil
}

func (s *Server) SetUserMetadata(ctx context.Context, req *connect.Request[user.SetUserMetadataRequest]) (*connect.Response[user.SetUserMetadataResponse], error) {
	result, err := s.command.BulkSetUserMetadata(ctx, req.Msg.UserId, "", setUserMetadataToDomain(req.Msg)...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.SetUserMetadataResponse{
		SetDate: timestamppb.New(result.EventDate),
	}), nil
}

func setUserMetadataToDomain(req *user.SetUserMetadataRequest) []*domain.Metadata {
	metadata := make([]*domain.Metadata, len(req.Metadata))
	for i, data := range req.Metadata {
		metadata[i] = &domain.Metadata{
			Key:   data.Key,
			Value: data.Value,
		}
	}
	return metadata
}

func (s *Server) DeleteUserMetadata(ctx context.Context, req *connect.Request[user.DeleteUserMetadataRequest]) (*connect.Response[user.DeleteUserMetadataResponse], error) {
	result, err := s.command.BulkRemoveUserMetadata(ctx, req.Msg.UserId, "", req.Msg.Keys...)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.DeleteUserMetadataResponse{
		DeletionDate: timestamppb.New(result.EventDate),
	}), nil
}
