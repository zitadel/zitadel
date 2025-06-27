package app

import (
	"context"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
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
		Pagination: &filter.PaginationResponse{
			TotalResult:  res.Count,
			AppliedLimit: uint64(req.GetPagination().GetLimit()),
		},
	}, nil
}

func (s *Server) GetApplicationKey(ctx context.Context, req *app.GetApplicationKeyRequest) (*app.GetApplicationKeyResponse, error) {
	queries, err := convert.GetApplicationKeyQueriesRequestToDomain(req.GetOrganizationId(), req.GetProjectId(), req.GetApplicationId())
	if err != nil {
		return nil, err
	}

	key, err := s.query.GetAuthNKeyByIDWithPermission(ctx, true, strings.TrimSpace(req.GetId()), s.checkPermission, queries...)
	if err != nil {
		return nil, err
	}

	return &app.GetApplicationKeyResponse{
		Id:             key.ID,
		Type:           app.ApplicationKeyType(key.Type),
		CreationDate:   timestamppb.New(key.CreationDate),
		ExpirationDate: timestamppb.New(key.Expiration),
	}, nil
}
