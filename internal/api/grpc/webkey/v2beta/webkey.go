package webkey

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	webkey "github.com/zitadel/zitadel/pkg/grpc/webkey/v2beta"
)

func (s *Server) CreateWebKey(ctx context.Context, req *connect.Request[webkey.CreateWebKeyRequest]) (_ *connect.Response[webkey.CreateWebKeyResponse], err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	webKey, err := s.command.CreateWebKey(ctx, createWebKeyRequestToConfig(req.Msg))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&webkey.CreateWebKeyResponse{
		Id:           webKey.KeyID,
		CreationDate: timestamppb.New(webKey.ObjectDetails.EventDate),
	}), nil
}

func (s *Server) ActivateWebKey(ctx context.Context, req *connect.Request[webkey.ActivateWebKeyRequest]) (_ *connect.Response[webkey.ActivateWebKeyResponse], err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	details, err := s.command.ActivateWebKey(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&webkey.ActivateWebKeyResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) DeleteWebKey(ctx context.Context, req *connect.Request[webkey.DeleteWebKeyRequest]) (_ *connect.Response[webkey.DeleteWebKeyResponse], err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	deletedAt, err := s.command.DeleteWebKey(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	var deletionDate *timestamppb.Timestamp
	if !deletedAt.IsZero() {
		deletionDate = timestamppb.New(deletedAt)
	}
	return connect.NewResponse(&webkey.DeleteWebKeyResponse{
		DeletionDate: deletionDate,
	}), nil
}

func (s *Server) ListWebKeys(ctx context.Context, _ *connect.Request[webkey.ListWebKeysRequest]) (_ *connect.Response[webkey.ListWebKeysResponse], err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	list, err := s.query.ListWebKeys(ctx)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&webkey.ListWebKeysResponse{
		WebKeys: webKeyDetailsListToPb(list),
	}), nil
}
