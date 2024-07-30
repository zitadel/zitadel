package webkey

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
	webkey "github.com/zitadel/zitadel/pkg/grpc/webkey/v3alpha"
)

func (s *Server) GenerateWebKey(ctx context.Context, req *webkey.GenerateWebKeyRequest) (_ *webkey.GenerateWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	webKey, err := s.command.GenerateWebKey(ctx, generateWebKeyRequestToConfig(req))
	if err != nil {
		return nil, err
	}

	return &webkey.GenerateWebKeyResponse{
		KeyId:   webKey.KeyID,
		Details: object.DomainToDetailsPb(webKey.ObjectDetails),
	}, nil
}

func (s *Server) ActivateWebKey(ctx context.Context, req *webkey.ActivateWebKeyRequest) (_ *webkey.ActivateWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.ActivateWebKey(ctx, req.GetKeyId())
	if err != nil {
		return nil, err
	}

	return &webkey.ActivateWebKeyResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeleteWebKey(ctx context.Context, req *webkey.DeleteWebKeyRequest) (_ *webkey.DeleteWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteWebKey(ctx, req.GetKeyId())
	if err != nil {
		return nil, err
	}

	return &webkey.DeleteWebKeyResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListWebKeys(ctx context.Context, req *webkey.ListWebKeysRequest) (_ *webkey.ListWebKeysResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	list, err := s.query.ListWebKeys(ctx)
	if err != nil {
		return nil, err
	}

	return &webkey.ListWebKeysResponse{
		WebKeys: webKeyDetailsListToPb(list),
	}, nil
}

func checkWebKeyFeature(ctx context.Context) error {
	if !authz.GetFeatures(ctx).WebKey {
		return zerrors.ThrowPreconditionFailed(nil, "WEBKEY-Ohx6E", "Errors.WebKey.FeatureDisabled")
	}
	return nil
}
