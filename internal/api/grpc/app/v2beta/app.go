package app

import (
	"context"
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func (s *Server) CreateApplication(ctx context.Context, req *connect.Request[app.CreateApplicationRequest]) (*connect.Response[app.CreateApplicationResponse], error) {
	switch t := req.Msg.GetCreationRequestType().(type) {
	case *app.CreateApplicationRequest_ApiRequest:
		apiApp, err := s.command.AddAPIApplication(ctx, convert.CreateAPIApplicationRequestToDomain(req.Msg.GetName(), req.Msg.GetProjectId(), req.Msg.GetId(), t.ApiRequest), "")
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&app.CreateApplicationResponse{
			AppId:        apiApp.AppID,
			CreationDate: timestamppb.New(apiApp.ChangeDate),
			CreationResponseType: &app.CreateApplicationResponse_ApiResponse{
				ApiResponse: &app.CreateAPIApplicationResponse{
					ClientId:     apiApp.ClientID,
					ClientSecret: apiApp.ClientSecretString,
				},
			},
		}), nil

	case *app.CreateApplicationRequest_OidcRequest:
		oidcAppRequest, err := convert.CreateOIDCAppRequestToDomain(req.Msg.GetName(), req.Msg.GetProjectId(), req.Msg.GetOidcRequest())
		if err != nil {
			return nil, err
		}

		oidcApp, err := s.command.AddOIDCApplication(ctx, oidcAppRequest, "")
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&app.CreateApplicationResponse{
			AppId:        oidcApp.AppID,
			CreationDate: timestamppb.New(oidcApp.ChangeDate),
			CreationResponseType: &app.CreateApplicationResponse_OidcResponse{
				OidcResponse: &app.CreateOIDCApplicationResponse{
					ClientId:           oidcApp.ClientID,
					ClientSecret:       oidcApp.ClientSecretString,
					NoneCompliant:      oidcApp.Compliance.NoneCompliant,
					ComplianceProblems: convert.ComplianceProblemsToLocalizedMessages(oidcApp.Compliance.Problems),
				},
			},
		}), nil

	case *app.CreateApplicationRequest_SamlRequest:
		samlAppRequest, err := convert.CreateSAMLAppRequestToDomain(req.Msg.GetName(), req.Msg.GetProjectId(), req.Msg.GetSamlRequest())
		if err != nil {
			return nil, err
		}

		samlApp, err := s.command.AddSAMLApplication(ctx, samlAppRequest, "")
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&app.CreateApplicationResponse{
			AppId:        samlApp.AppID,
			CreationDate: timestamppb.New(samlApp.ChangeDate),
			CreationResponseType: &app.CreateApplicationResponse_SamlResponse{
				SamlResponse: &app.CreateSAMLApplicationResponse{},
			},
		}), nil
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "APP-0iiN46", "unknown app type")
	}
}

func (s *Server) UpdateApplication(ctx context.Context, req *connect.Request[app.UpdateApplicationRequest]) (*connect.Response[app.UpdateApplicationResponse], error) {
	var changedTime time.Time

	if name := strings.TrimSpace(req.Msg.GetName()); name != "" {
		updatedDetails, err := s.command.UpdateApplicationName(
			ctx,
			req.Msg.GetProjectId(),
			&domain.ChangeApp{
				AppID:   req.Msg.GetId(),
				AppName: name,
			},
			"",
		)
		if err != nil {
			return nil, err
		}

		changedTime = updatedDetails.EventDate
	}

	switch t := req.Msg.GetUpdateRequestType().(type) {
	case *app.UpdateApplicationRequest_ApiConfigurationRequest:
		updatedAPIApp, err := s.command.UpdateAPIApplication(ctx, convert.UpdateAPIApplicationConfigurationRequestToDomain(req.Msg.GetId(), req.Msg.GetProjectId(), t.ApiConfigurationRequest), "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedAPIApp.ChangeDate

	case *app.UpdateApplicationRequest_OidcConfigurationRequest:
		oidcApp, err := convert.UpdateOIDCAppConfigRequestToDomain(req.Msg.GetId(), req.Msg.GetProjectId(), t.OidcConfigurationRequest)
		if err != nil {
			return nil, err
		}

		updatedOIDCApp, err := s.command.UpdateOIDCApplication(ctx, oidcApp, "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedOIDCApp.ChangeDate

	case *app.UpdateApplicationRequest_SamlConfigurationRequest:
		samlApp, err := convert.UpdateSAMLAppConfigRequestToDomain(req.Msg.GetId(), req.Msg.GetProjectId(), t.SamlConfigurationRequest)
		if err != nil {
			return nil, err
		}

		updatedSAMLApp, err := s.command.UpdateSAMLApplication(ctx, samlApp, "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedSAMLApp.ChangeDate
	}

	return connect.NewResponse(&app.UpdateApplicationResponse{
		ChangeDate: timestamppb.New(changedTime),
	}), nil
}

func (s *Server) DeleteApplication(ctx context.Context, req *connect.Request[app.DeleteApplicationRequest]) (*connect.Response[app.DeleteApplicationResponse], error) {
	details, err := s.command.RemoveApplication(ctx, req.Msg.GetProjectId(), req.Msg.GetId(), "")
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.DeleteApplicationResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) DeactivateApplication(ctx context.Context, req *connect.Request[app.DeactivateApplicationRequest]) (*connect.Response[app.DeactivateApplicationResponse], error) {
	details, err := s.command.DeactivateApplication(ctx, req.Msg.GetProjectId(), req.Msg.GetId(), "")
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.DeactivateApplicationResponse{
		DeactivationDate: timestamppb.New(details.EventDate),
	}), nil

}

func (s *Server) ReactivateApplication(ctx context.Context, req *connect.Request[app.ReactivateApplicationRequest]) (*connect.Response[app.ReactivateApplicationResponse], error) {
	details, err := s.command.ReactivateApplication(ctx, req.Msg.GetProjectId(), req.Msg.GetId(), "")
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&app.ReactivateApplicationResponse{
		ReactivationDate: timestamppb.New(details.EventDate),
	}), nil

}

func (s *Server) RegenerateClientSecret(ctx context.Context, req *connect.Request[app.RegenerateClientSecretRequest]) (*connect.Response[app.RegenerateClientSecretResponse], error) {
	var secret string
	var changeDate time.Time

	switch req.Msg.GetAppType().(type) {
	case *app.RegenerateClientSecretRequest_IsApi:
		config, err := s.command.ChangeAPIApplicationSecret(ctx, req.Msg.GetProjectId(), req.Msg.GetApplicationId(), "")
		if err != nil {
			return nil, err
		}
		secret = config.ClientSecretString
		changeDate = config.ChangeDate

	case *app.RegenerateClientSecretRequest_IsOidc:
		config, err := s.command.ChangeOIDCApplicationSecret(ctx, req.Msg.GetProjectId(), req.Msg.GetApplicationId(), "")
		if err != nil {
			return nil, err
		}

		secret = config.ClientSecretString
		changeDate = config.ChangeDate

	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "APP-aLWIzw", "unknown app type")
	}

	return connect.NewResponse(&app.RegenerateClientSecretResponse{
		ClientSecret: secret,
		CreationDate: timestamppb.New(changeDate),
	}), nil
}
