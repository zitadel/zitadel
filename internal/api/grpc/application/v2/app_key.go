package app

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/application/v2/convert"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func (s *Server) CreateApplicationKey(ctx context.Context, req *connect.Request[application.CreateApplicationKeyRequest]) (*connect.Response[application.CreateApplicationKeyResponse], error) {
	domainReq := convert.CreateAPIClientKeyRequestToDomain(req.Msg)

	appKey, err := s.command.AddApplicationKey(ctx, domainReq, "")
	if err != nil {
		return nil, err
	}

	keyDetails, err := appKey.Detail()
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.CreateApplicationKeyResponse{
		KeyId:        appKey.KeyID,
		CreationDate: timestamppb.New(appKey.ChangeDate),
		KeyDetails:   keyDetails,
	}), nil
}

func (s *Server) DeleteApplicationKey(ctx context.Context, req *connect.Request[application.DeleteApplicationKeyRequest]) (*connect.Response[application.DeleteApplicationKeyResponse], error) {
	deletionDetails, err := s.command.RemoveApplicationKey(ctx,
		strings.TrimSpace(req.Msg.GetProjectId()),
		strings.TrimSpace(req.Msg.GetApplicationId()),
		strings.TrimSpace(req.Msg.GetKeyId()),
		"",
	)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.DeleteApplicationKeyResponse{
		DeletionDate: timestamppb.New(deletionDetails.EventDate),
	}), nil
}
