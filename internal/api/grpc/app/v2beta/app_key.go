package app

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func (s *Server) CreateApplicationKey(ctx context.Context, req *connect.Request[app.CreateApplicationKeyRequest]) (*connect.Response[app.CreateApplicationKeyResponse], error) {
	domainReq := convert.CreateAPIClientKeyRequestToDomain(req.Msg)

	appKey, err := s.command.AddApplicationKey(ctx, domainReq, "")
	if err != nil {
		return nil, err
	}

	keyDetails, err := appKey.Detail()
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.CreateApplicationKeyResponse{
		Id:           appKey.KeyID,
		CreationDate: timestamppb.New(appKey.ChangeDate),
		KeyDetails:   keyDetails,
	}), nil
}

func (s *Server) DeleteApplicationKey(ctx context.Context, req *connect.Request[app.DeleteApplicationKeyRequest]) (*connect.Response[app.DeleteApplicationKeyResponse], error) {
	deletionDetails, err := s.command.RemoveApplicationKey(ctx,
		strings.TrimSpace(req.Msg.GetProjectId()),
		strings.TrimSpace(req.Msg.GetApplicationId()),
		strings.TrimSpace(req.Msg.GetId()),
		strings.TrimSpace(req.Msg.GetOrganizationId()),
	)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.DeleteApplicationKeyResponse{
		DeletionDate: timestamppb.New(deletionDetails.EventDate),
	}), nil
}
