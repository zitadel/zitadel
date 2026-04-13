package app

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func (s *Server) GetApplication(ctx context.Context, req *connect.Request[app.GetApplicationRequest]) (*connect.Response[app.GetApplicationResponse], error) {
	res, err := s.query.AppByIDWithPermission(ctx, req.Msg.GetId(), false, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.GetApplicationResponse{
		App: convert.AppToPb(res),
	}), nil
}

func (s *Server) ListApplications(ctx context.Context, req *connect.Request[app.ListApplicationsRequest]) (*connect.Response[app.ListApplicationsResponse], error) {
	queries, err := convert.ListApplicationsRequestToModel(s.systemDefaults, req.Msg)
	if err != nil {
		return nil, err
	}

	res, err := s.query.SearchApps(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.ListApplicationsResponse{
		Applications: convert.AppsToPb(res.Apps),
		Pagination:   filter.QueryToPaginationPb(queries.SearchRequest, res.SearchResponse),
	}), nil
}

func (s *Server) GetApplicationKey(ctx context.Context, req *connect.Request[app.GetApplicationKeyRequest]) (*connect.Response[app.GetApplicationKeyResponse], error) {
	queries, err := convert.GetApplicationKeyQueriesRequestToDomain(req.Msg.GetOrganizationId(), req.Msg.GetProjectId(), req.Msg.GetApplicationId())
	if err != nil {
		return nil, err
	}

	key, err := s.query.GetAuthNKeyByIDWithPermission(ctx, true, strings.TrimSpace(req.Msg.GetId()), s.checkPermission, queries...)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.GetApplicationKeyResponse{
		Id:             key.ID,
		CreationDate:   timestamppb.New(key.CreationDate),
		ExpirationDate: timestamppb.New(key.Expiration),
	}), nil
}

func (s *Server) ListApplicationKeys(ctx context.Context, req *connect.Request[app.ListApplicationKeysRequest]) (*connect.Response[app.ListApplicationKeysResponse], error) {
	queries, err := convert.ListApplicationKeysRequestToDomain(s.systemDefaults, req.Msg)
	if err != nil {
		return nil, err
	}

	res, err := s.query.SearchAuthNKeys(ctx, queries, query.JoinFilterUnspecified, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.ListApplicationKeysResponse{
		Keys:       convert.ApplicationKeysToPb(res.AuthNKeys),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, res.SearchResponse),
	}), nil
}
