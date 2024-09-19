package webkey

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	resource_object "github.com/zitadel/zitadel/v2/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/v2/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
	object "github.com/zitadel/zitadel/v2/pkg/grpc/object/v3alpha"
	webkey "github.com/zitadel/zitadel/v2/pkg/grpc/resources/webkey/v3alpha"
)

func (s *Server) CreateWebKey(ctx context.Context, req *webkey.CreateWebKeyRequest) (_ *webkey.CreateWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	webKey, err := s.command.CreateWebKey(ctx, createWebKeyRequestToConfig(req))
	if err != nil {
		return nil, err
	}

	return &webkey.CreateWebKeyResponse{
		Details: resource_object.DomainToDetailsPb(webKey.ObjectDetails, object.OwnerType_OWNER_TYPE_INSTANCE, authz.GetInstance(ctx).InstanceID()),
	}, nil
}

func (s *Server) ActivateWebKey(ctx context.Context, req *webkey.ActivateWebKeyRequest) (_ *webkey.ActivateWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.ActivateWebKey(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &webkey.ActivateWebKeyResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, authz.GetInstance(ctx).InstanceID()),
	}, nil
}

func (s *Server) DeleteWebKey(ctx context.Context, req *webkey.DeleteWebKeyRequest) (_ *webkey.DeleteWebKeyResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteWebKey(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &webkey.DeleteWebKeyResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_INSTANCE, authz.GetInstance(ctx).InstanceID()),
	}, nil
}

func (s *Server) ListWebKeys(ctx context.Context, _ *webkey.ListWebKeysRequest) (_ *webkey.ListWebKeysResponse, err error) {
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
		WebKeys: webKeyDetailsListToPb(list, authz.GetInstance(ctx).InstanceID()),
	}, nil
}

func checkWebKeyFeature(ctx context.Context) error {
	if !authz.GetFeatures(ctx).WebKey {
		return zerrors.ThrowPreconditionFailed(nil, "WEBKEY-Ohx6E", "Errors.WebKey.FeatureDisabled")
	}
	return nil
}
