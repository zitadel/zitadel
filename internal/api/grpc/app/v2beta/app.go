package app

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	"github.com/zitadel/zitadel/internal/zerrors"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func (s *Server) CreateApplication(ctx context.Context, req *app.CreateApplicationRequest) (*app.CreateApplicationResponse, error) {
	switch t := req.GetCreationRequestType().(type) {
	case *app.CreateApplicationRequest_ApiRequest:
		apiApp, err := s.command.AddAPIApplication(ctx, convert.CreateAPIApplicationRequestToDomain(req.GetName(), req.GetProjectId(), t.ApiRequest), authz.GetCtxData(ctx).OrgID)
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

		oidcApp, err := s.command.AddOIDCApplication(ctx, oidcAppRequest, authz.GetCtxData(ctx).OrgID)
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

		samlApp, err := s.command.AddSAMLApplication(ctx, samlAppRequest, authz.GetCtxData(ctx).OrgID)
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
