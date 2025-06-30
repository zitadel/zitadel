package webkey

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	webkey "github.com/zitadel/zitadel/pkg/grpc/webkey/v2beta"
)

func (s *Server) CreateWebKey(ctx context.Context, req *webkey.CreateWebKeyRequest) (_ *webkey.CreateWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	webKey, err := s.command.CreateWebKey(ctx, createWebKeyRequestToConfig(req))
	if err != nil {
		return nil, err
	}

	return &webkey.CreateWebKeyResponse{
		Id:           webKey.KeyID,
		CreationDate: timestamppb.New(webKey.ObjectDetails.EventDate),
	}, nil
}

func (s *Server) ActivateWebKey(ctx context.Context, req *webkey.ActivateWebKeyRequest) (_ *webkey.ActivateWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	details, err := s.command.ActivateWebKey(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &webkey.ActivateWebKeyResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) DeleteWebKey(ctx context.Context, req *webkey.DeleteWebKeyRequest) (_ *webkey.DeleteWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	deletedAt, err := s.command.DeleteWebKey(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	var deletionDate *timestamppb.Timestamp
	if !deletedAt.IsZero() {
		deletionDate = timestamppb.New(deletedAt)
	}
	return &webkey.DeleteWebKeyResponse{
		DeletionDate: deletionDate,
	}, nil
}

func (s *Server) ListWebKeys(ctx context.Context, _ *webkey.ListWebKeysRequest) (_ *webkey.ListWebKeysResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	list, err := s.query.ListWebKeys(ctx)
	if err != nil {
		return nil, err
	}

	return &webkey.ListWebKeysResponse{
		WebKeys: webKeyDetailsListToPb(list),
	}, nil
}
