package app

import (
	"context"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func (s *Server) CreateApplication(ctx context.Context, req *app.CreateApplicationRequest) (*app.CreateApplicationResponse, error) {
	switch t := req.GetCreationRequestType().(type) {
	case *app.CreateApplicationRequest_ApiRequest:
		// should also pass the Id in case we don't want to regen it
		apiApp, err := s.command.AddAPIApplication(ctx, convert.CreateAPIApplicationRequestToDomain(req.GetName(), req.GetProjectId(), t.ApiRequest), "")
		if err != nil {
			return nil, err
		}

		return &app.CreateApplicationResponse{
			AppId:        apiApp.AppID,
			CreationDate: timestamppb.New(apiApp.ChangeDate),
			CreationResponseType: &app.CreateApplicationResponse_ApiResponse{
				ApiResponse: &app.CreateAPIApplicationResponse{
					ClientId:     apiApp.ClientID,
					ClientSecret: apiApp.ClientSecretString,
				},
			},
		}, nil

	case *app.CreateApplicationRequest_OidcRequest:
		oidcAppRequest, err := convert.CreateOIDCAppRequestToDomain(req.GetName(), req.GetProjectId(), req.GetOidcRequest())
		if err != nil {
			return nil, err
		}

		oidcApp, err := s.command.AddOIDCApplication(ctx, oidcAppRequest, "")
		if err != nil {
			return nil, err
		}

		return &app.CreateApplicationResponse{
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
		}, nil

	case *app.CreateApplicationRequest_SamlRequest:
		samlAppRequest, err := convert.CreateSAMLAppRequestToDomain(req.GetName(), req.GetProjectId(), req.GetSamlRequest())
		if err != nil {
			return nil, err
		}

		samlApp, err := s.command.AddSAMLApplication(ctx, samlAppRequest, "")
		if err != nil {
			return nil, err
		}

		return &app.CreateApplicationResponse{
			AppId:        samlApp.AppID,
			CreationDate: timestamppb.New(samlApp.ChangeDate),
			CreationResponseType: &app.CreateApplicationResponse_SamlResponse{
				SamlResponse: &app.CreateSAMLApplicationResponse{},
			},
		}, nil
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "APP-0iiN46", "unknown app type")
	}
}

func (s *Server) UpdateApplication(ctx context.Context, req *app.UpdateApplicationRequest) (*app.UpdateApplicationResponse, error) {
	var changedTime time.Time

	if name := strings.TrimSpace(req.GetName()); name != "" {
		updatedDetails, err := s.command.PatchApplication(
			ctx,
			req.GetProjectId(),
			&domain.ChangeApp{
				AppID:   req.GetId(),
				AppName: name,
			},
			"",
		)
		if err != nil {
			return nil, err
		}

		changedTime = updatedDetails.EventDate
	}

	switch t := req.GetUpdateRequestType().(type) {
	case *app.UpdateApplicationRequest_ApiConfigurationRequest:
		updatedAPIApp, err := s.command.PatchAPIApplication(ctx, convert.PatchAPIApplicationConfigurationRequestToDomain(req.GetId(), req.GetProjectId(), t.ApiConfigurationRequest), "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedAPIApp.ChangeDate

	case *app.UpdateApplicationRequest_OidcConfigurationRequest:
		oidcApp, err := convert.PatchOIDCAppConfigRequestToDomain(req.GetId(), req.GetProjectId(), t.OidcConfigurationRequest)
		if err != nil {
			return nil, err
		}

		updatedOIDCApp, err := s.command.PatchOIDCApplication(ctx, oidcApp, "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedOIDCApp.ChangeDate

	case *app.UpdateApplicationRequest_SamlConfigurationRequest:
		samlApp, err := convert.PatchSAMLAppConfigRequestToDomain(req.GetId(), req.GetProjectId(), t.SamlConfigurationRequest)
		if err != nil {
			return nil, err
		}

		updatedSAMLApp, err := s.command.PatchSAMLApplication(ctx, samlApp, "")
		if err != nil {
			return nil, err
		}

		changedTime = updatedSAMLApp.ChangeDate
	}

	return &app.UpdateApplicationResponse{
		ChangeDate: timestamppb.New(changedTime),
	}, nil
}

func (s *Server) DeleteApplication(ctx context.Context, req *app.DeleteApplicationRequest) (*app.DeleteApplicationResponse, error) {
	details, err := s.command.RemoveApplication(ctx, req.GetProjectId(), req.GetId(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}

	return &app.DeleteApplicationResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) DeactivateApplication(ctx context.Context, req *app.DeactivateApplicationRequest) (*app.DeactivateApplicationResponse, error) {
	details, err := s.command.DeactivateApplication(ctx, req.GetProjectId(), req.GetId(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}

	return &app.DeactivateApplicationResponse{
		DeactivationDate: timestamppb.New(details.EventDate),
	}, nil

}

func (s *Server) ReactivateApplication(ctx context.Context, req *app.ReactivateApplicationRequest) (*app.ReactivateApplicationResponse, error) {
	details, err := s.command.ReactivateApplication(ctx, req.GetProjectId(), req.GetId(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}

	return &app.ReactivateApplicationResponse{
		ReactivationDate: timestamppb.New(details.EventDate),
	}, nil

}

func (s *Server) RegenerateClientSecret(ctx context.Context, req *app.RegenerateClientSecretRequest) (*app.RegenerateClientSecretResponse, error) {
	var secret string
	var changeDate time.Time

	switch req.GetAppType().(type) {
	case *app.RegenerateClientSecretRequest_IsApi:
		config, err := s.command.ChangeAPIApplicationSecret(ctx, req.GetProjectId(), req.GetApplicationId(), authz.GetCtxData(ctx).OrgID)
		if err != nil {
			return nil, err
		}
		secret = config.ClientSecretString
		changeDate = config.ChangeDate

	case *app.RegenerateClientSecretRequest_IsOidc:
		config, err := s.command.ChangeOIDCApplicationSecret(ctx, req.GetProjectId(), req.GetApplicationId(), authz.GetCtxData(ctx).OrgID)
		if err != nil {
			return nil, err
		}

		secret = config.ClientSecretString
		changeDate = config.ChangeDate

	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "APP-aLWIzw", "unknown app type")
	}

	return &app.RegenerateClientSecretResponse{
		ClientSecret: secret,
		CreationDate: timestamppb.New(changeDate),
	}, nil
}
