package app

import (
	"context"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func (s *Server) CreateApplicationKey(ctx context.Context, req *app.CreateApplicationKeyRequest) (*app.CreateApplicationKeyResponse, error) {
	domainReq := convert.CreateAPIClientKeyRequestToDomain(req)

	appKey, err := s.command.AddApplicationKey(ctx, domainReq, "")
	if err != nil {
		return nil, err
	}

	keyDetails, err := appKey.Detail()
	if err != nil {
		return nil, err
	}

	return &app.CreateApplicationKeyResponse{
		Id:           appKey.KeyID,
		CreationDate: timestamppb.New(appKey.ChangeDate),
		KeyDetails:   keyDetails,
	}, nil
}

func (s *Server) DeleteApplicationKey(ctx context.Context, req *app.DeleteApplicationKeyRequest) (*app.DeleteApplicationKeyResponse, error) {
	deletionDetails, err := s.command.RemoveApplicationKey(ctx,
		strings.TrimSpace(req.GetProjectId()),
		strings.TrimSpace(req.GetApplicationId()),
		strings.TrimSpace(req.GetId()),
		strings.TrimSpace(req.GetOrganizationId()),
	)
	if err != nil {
		return nil, err
	}

	return &app.DeleteApplicationKeyResponse{
		DeletionDate: timestamppb.New(deletionDetails.EventDate),
	}, nil
}
