package app

import (
	"context"
	"strings"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/application/v2/convert"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func (s *Server) CreateApplication(ctx context.Context, req *connect.Request[application.CreateApplicationRequest]) (*connect.Response[application.CreateApplicationResponse], error) {
	switch t := req.Msg.GetApplicationType().(type) {
	case *application.CreateApplicationRequest_ApiConfiguration:
		apiApp, err := s.command.AddAPIApplication(ctx, convert.CreateAPIApplicationRequestToDomain(req.Msg.GetName(), req.Msg.GetProjectId(), req.Msg.GetApplicationId(), t.ApiConfiguration), "")
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&application.CreateApplicationResponse{
			ApplicationId: apiApp.AppID,
			CreationDate:  timestamppb.New(apiApp.ChangeDate),
			ApplicationType: &application.CreateApplicationResponse_ApiConfiguration{
				ApiConfiguration: &application.CreateAPIApplicationResponse{
					ClientId:     apiApp.ClientID,
					ClientSecret: apiApp.ClientSecretString,
				},
			},
		}), nil

	case *application.CreateApplicationRequest_OidcConfiguration:
		oidcAppRequest, err := convert.CreateOIDCAppRequestToDomain(req.Msg.GetName(), req.Msg.GetApplicationId(), req.Msg.GetProjectId(), req.Msg.GetOidcConfiguration())
		if err != nil {
			return nil, err
		}

		oidcApp, err := s.command.AddOIDCApplication(ctx, oidcAppRequest, "")
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&application.CreateApplicationResponse{
			ApplicationId: oidcApp.AppID,
			CreationDate:  timestamppb.New(oidcApp.ChangeDate),
			ApplicationType: &application.CreateApplicationResponse_OidcConfiguration{
				OidcConfiguration: &application.CreateOIDCApplicationResponse{
					ClientId:           oidcApp.ClientID,
					ClientSecret:       oidcApp.ClientSecretString,
					NonCompliant:       oidcApp.Compliance.NoneCompliant,
					ComplianceProblems: convert.ComplianceProblemsToLocalizedMessages(oidcApp.Compliance.Problems),
				},
			},
		}), nil

	case *application.CreateApplicationRequest_SamlConfiguration:
		samlAppRequest, err := convert.CreateSAMLAppRequestToDomain(req.Msg.GetName(), req.Msg.GetProjectId(), req.Msg.GetSamlConfiguration())
		if err != nil {
			return nil, err
		}

		samlApp, err := s.command.AddSAMLApplication(ctx, samlAppRequest, "")
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&application.CreateApplicationResponse{
			ApplicationId: samlApp.AppID,
			CreationDate:  timestamppb.New(samlApp.ChangeDate),
			ApplicationType: &application.CreateApplicationResponse_SamlConfiguration{
				SamlConfiguration: &application.CreateSAMLApplicationResponse{},
			},
		}), nil
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "APP-0iiN46", "unknown app type")
	}
}

func (s *Server) UpdateApplication(ctx context.Context, req *connect.Request[application.UpdateApplicationRequest]) (*connect.Response[application.UpdateApplicationResponse], error) {
	var changedTime time.Time

	if name := strings.TrimSpace(req.Msg.GetName()); name != "" {
		updatedDetails, err := s.command.UpdateApplicationName(
			ctx,
			req.Msg.GetProjectId(),
			&domain.ChangeApp{
				AppID:   req.Msg.GetApplicationId(),
				AppName: name,
			},
			"",
		)
		if err != nil {
			return nil, err
		}

		changedTime = updatedDetails.EventDate
	}

	switch t := req.Msg.GetApplicationType().(type) {
	case *application.UpdateApplicationRequest_ApiConfiguration:
		updatedAPIApp, err := s.command.UpdateAPIApplication(ctx, convert.UpdateAPIApplicationConfigurationRequestToDomain(req.Msg.GetApplicationId(), req.Msg.GetProjectId(), t.ApiConfiguration), "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedAPIApp.ChangeDate

	case *application.UpdateApplicationRequest_OidcConfiguration:
		oidcApp, err := convert.UpdateOIDCAppConfigRequestToDomain(req.Msg.GetApplicationId(), req.Msg.GetProjectId(), t.OidcConfiguration)
		if err != nil {
			return nil, err
		}

		updatedOIDCApp, err := s.command.UpdateOIDCApplication(ctx, oidcApp, "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedOIDCApp.ChangeDate

	case *application.UpdateApplicationRequest_SamlConfiguration:
		samlApp, err := convert.UpdateSAMLAppConfigRequestToDomain(req.Msg.GetApplicationId(), req.Msg.GetProjectId(), t.SamlConfiguration)
		if err != nil {
			return nil, err
		}

		updatedSAMLApp, err := s.command.UpdateSAMLApplication(ctx, samlApp, "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedSAMLApp.ChangeDate
	}

	return connect.NewResponse(&application.UpdateApplicationResponse{
		ChangeDate: timestamppb.New(changedTime),
	}), nil
}

func (s *Server) DeleteApplication(ctx context.Context, req *connect.Request[application.DeleteApplicationRequest]) (*connect.Response[application.DeleteApplicationResponse], error) {
	details, err := s.command.RemoveApplication(ctx, req.Msg.GetProjectId(), req.Msg.GetApplicationId(), "")
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.DeleteApplicationResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) DeactivateApplication(ctx context.Context, req *connect.Request[application.DeactivateApplicationRequest]) (*connect.Response[application.DeactivateApplicationResponse], error) {
	details, err := s.command.DeactivateApplication(ctx, req.Msg.GetProjectId(), req.Msg.GetApplicationId(), "")
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.DeactivateApplicationResponse{
		DeactivationDate: timestamppb.New(details.EventDate),
	}), nil

}

func (s *Server) ReactivateApplication(ctx context.Context, req *connect.Request[application.ReactivateApplicationRequest]) (*connect.Response[application.ReactivateApplicationResponse], error) {
	details, err := s.command.ReactivateApplication(ctx, req.Msg.GetProjectId(), req.Msg.GetApplicationId(), "")
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&application.ReactivateApplicationResponse{
		ReactivationDate: timestamppb.New(details.EventDate),
	}), nil

}

func (s *Server) GenerateClientSecret(ctx context.Context, req *connect.Request[application.GenerateClientSecretRequest]) (*connect.Response[application.GenerateClientSecretResponse], error) {
	var secret string
	var changeDate time.Time

	secret, changeDate, err := s.command.ChangeApplicationSecret(ctx, req.Msg.GetProjectId(), req.Msg.GetApplicationId(), "")
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&application.GenerateClientSecretResponse{
		ClientSecret: secret,
		CreationDate: timestamppb.New(changeDate),
	}), nil
}
