package settings

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func (s *Server) SetSecuritySettings(ctx context.Context, req *connect.Request[settings.SetSecuritySettingsRequest]) (*connect.Response[settings.SetSecuritySettingsResponse], error) {
	details, err := s.command.SetSecurityPolicy(ctx, securitySettingsToCommand(req.Msg))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.SetSecuritySettingsResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) SetHostedLoginTranslation(ctx context.Context, req *connect.Request[settings.SetHostedLoginTranslationRequest]) (*connect.Response[settings.SetHostedLoginTranslationResponse], error) {
	res, err := s.command.SetHostedLoginTranslation(ctx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(res), nil
}
