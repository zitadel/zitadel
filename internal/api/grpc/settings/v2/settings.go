package settings

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

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

func (s *Server) SetOrganizationSettings(ctx context.Context, req *connect.Request[settings.SetOrganizationSettingsRequest]) (*connect.Response[settings.SetOrganizationSettingsResponse], error) {
	details, err := s.command.SetOrganizationSettings(ctx, organizationSettingsToCommand(req.Msg))
	if err != nil {
		return nil, err
	}
	var setDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		setDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&settings.SetOrganizationSettingsResponse{
		SetDate: setDate,
	}), nil
}

func (s *Server) DeleteOrganizationSettings(ctx context.Context, req *connect.Request[settings.DeleteOrganizationSettingsRequest]) (*connect.Response[settings.DeleteOrganizationSettingsResponse], error) {
	details, err := s.command.DeleteOrganizationSettings(ctx, req.Msg.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		deletionDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&settings.DeleteOrganizationSettingsResponse{
		DeletionDate: deletionDate,
	}), nil
}
