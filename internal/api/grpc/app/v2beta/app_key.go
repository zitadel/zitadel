package app

import (
	"context"

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
