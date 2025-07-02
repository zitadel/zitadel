package settings

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func (s *Server) SetSecuritySettings(ctx context.Context, req *settings.SetSecuritySettingsRequest) (*settings.SetSecuritySettingsResponse, error) {
	details, err := s.command.SetSecurityPolicy(ctx, securitySettingsToCommand(req))
	if err != nil {
		return nil, err
	}
	return &settings.SetSecuritySettingsResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) SetHostedLoginTranslation(ctx context.Context, req *settings.SetHostedLoginTranslationRequest) (*settings.SetHostedLoginTranslationResponse, error) {
	res, err := s.command.SetHostedLoginTranslation(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
