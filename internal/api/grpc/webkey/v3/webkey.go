package webkey

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/webkey/v3alpha"
)

func (s *Server) GenerateWebKey(ctx context.Context, req *v3alpha.GenerateWebKeyRequest) (_ *v3alpha.GenerateWebKeyResponse, err error) {
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
	details, err := s.command.ActivateWebKey(ctx, req.GetKeyId())
	if err != nil {
		return nil, err
	}

	return &v3alpha.ActivateWebKeyResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) DeleteWebKey(ctx context.Context, req *v3alpha.DeleteWebKeyRequest) (_ *v3alpha.DeleteWebKeyResponse, err error) {
	details, err := s.command.DeleteWebKey(ctx, req.GetKeyId())
	if err != nil {
		return nil, err
	}

	return &v3alpha.DeleteWebKeyResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListWebKeys(ctx context.Context, req *v3alpha.ListWebKeysRequest) (_ *v3alpha.ListWebKeysResponse, err error) {
	list, err := s.query.ListWebKeys(ctx)
	if err != nil {
		return nil, err
	}

	return &v3alpha.ListWebKeysResponse{
		WebKeys: webKeyDetailsListToPb(list),
	}, nil
}
