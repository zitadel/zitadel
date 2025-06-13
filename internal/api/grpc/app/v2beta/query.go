package app

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
)

func (s *Server) GetApplication(ctx context.Context, req *app.GetApplicationRequest) (*app.GetApplicationResponse, error) {
	res, err := s.query.AppByID(ctx, req.GetApplicationId(), false)
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

	res, err := s.query.SearchApps(ctx, queries, false)
	if err != nil {
		return nil, err
	}

	return &app.ListApplicationsResponse{
		App: convert.AppsToPb(res.Apps),
		Pagination: &filter.PaginationResponse{
			TotalResult:  res.Count,
			AppliedLimit: uint64(req.GetPagination().GetLimit()),
		},
	}, nil
}
