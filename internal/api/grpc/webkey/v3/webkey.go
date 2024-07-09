package webkey

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/webkey/v3alpha"
)

func (s *Server) GenerateWebKey(ctx context.Context, req *v3alpha.GenerateWebKeyRequest) (_ *v3alpha.GenerateWebKeyResponse, err error) {
	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	webKey, err := s.command.GenerateWebKey(ctx, generateWebKeyRequestToConfig(req))
	if err != nil {
		return nil, err
	}

	return &v3alpha.GenerateWebKeyResponse{
		KeyId:   webKey.KeyID,
		Details: object.DomainToDetailsPb(webKey.ObjectDetails),
	}, nil
}

func (s *Server) ActivateWebKey(ctx context.Context, req *v3alpha.ActivateWebKeyRequest) (_ *v3alpha.ActivateWebKeyResponse, err error) {
	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.ActivateWebKey(ctx, req.GetKeyId())
	if err != nil {
		return nil, err
	}

	return &v3alpha.ActivateWebKeyResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeleteWebKey(ctx context.Context, req *v3alpha.DeleteWebKeyRequest) (_ *v3alpha.DeleteWebKeyResponse, err error) {
	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.DeleteWebKey(ctx, req.GetKeyId())
	if err != nil {
		return nil, err
	}

	return &v3alpha.DeleteWebKeyResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListWebKeys(ctx context.Context, req *v3alpha.ListWebKeysRequest) (_ *v3alpha.ListWebKeysResponse, err error) {
	if err = checkWebKeyFeature(ctx); err != nil {
		return nil, err
	}
	list, err := s.query.ListWebKeys(ctx)
	if err != nil {
		return nil, err
	}

	return &v3alpha.ListWebKeysResponse{
		WebKeys: webKeyDetailsListToPb(list),
	}, nil
}

func checkWebKeyFeature(ctx context.Context) error {
	if !authz.GetFeatures(ctx).TokenExchange {
		return zerrors.ThrowPreconditionFailed(nil, "WEBKEY-Ohx6E", "Errors.WebKey.FeatureDisabled")
	}
	return nil
}
