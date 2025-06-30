package app

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func (s *Server) GetApplication(ctx context.Context, req *app.GetApplicationRequest) (*app.GetApplicationResponse, error) {
	res, err := s.query.AppByIDWithPermission(ctx, req.GetId(), false, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return &app.GetApplicationResponse{
		App: convert.AppToPb(res),
	}, nil
}

func (s *Server) ListApplications(ctx context.Context, req *app.ListApplicationsRequest) (*app.ListApplicationsResponse, error) {
	queries, err := convert.ListApplicationsRequestToModel(s.systemDefaults, req)
	if err != nil {
		return nil, err
	}

	res, err := s.query.SearchApps(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}

	return &app.ListApplicationsResponse{
		Applications: convert.AppsToPb(res.Apps),
		Pagination:   filter.QueryToPaginationPb(queries.SearchRequest, res.SearchResponse),
	}, nil
}
