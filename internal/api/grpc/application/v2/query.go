package app

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/application/v2/convert"
	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func (s *Server) GetApplication(ctx context.Context, req *connect.Request[application.GetApplicationRequest]) (*connect.Response[application.GetApplicationResponse], error) {
	res, err := s.query.AppByIDWithPermission(ctx, req.Msg.GetApplicationId(), false, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.GetApplicationResponse{
		Application: convert.AppToPb(res),
	}), nil
}

func (s *Server) ListApplications(ctx context.Context, req *connect.Request[application.ListApplicationsRequest]) (*connect.Response[application.ListApplicationsResponse], error) {
	queries, err := convert.ListApplicationsRequestToModel(s.systemDefaults, req.Msg)
	if err != nil {
		return nil, err
	}

	res, err := s.query.SearchApps(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.ListApplicationsResponse{
		Applications: convert.AppsToPb(res.Apps),
		Pagination:   filter.QueryToPaginationPb(queries.SearchRequest, res.SearchResponse),
	}), nil
}

func (s *Server) GetApplicationKey(ctx context.Context, req *connect.Request[application.GetApplicationKeyRequest]) (*connect.Response[application.GetApplicationKeyResponse], error) {
	key, err := s.query.GetAuthNKeyByIDWithPermission(ctx, true, strings.TrimSpace(req.Msg.GetKeyId()), s.checkPermission)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.GetApplicationKeyResponse{
		KeyId:          key.ID,
		CreationDate:   timestamppb.New(key.CreationDate),
		ExpirationDate: timestamppb.New(key.Expiration),
	}), nil
}

func (s *Server) ListApplicationKeys(ctx context.Context, req *connect.Request[application.ListApplicationKeysRequest]) (*connect.Response[application.ListApplicationKeysResponse], error) {
	queries, err := convert.ListApplicationKeysRequestToDomain(s.systemDefaults, req.Msg)
	if err != nil {
		return nil, err
	}

	res, err := s.query.SearchAuthNKeys(ctx, queries, query.JoinFilterUnspecified, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.ListApplicationKeysResponse{
		Keys:       convert.ApplicationKeysToPb(res.AuthNKeys),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, res.SearchResponse),
	}), nil
}
