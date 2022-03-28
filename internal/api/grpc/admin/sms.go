package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	settings_pb "github.com/caos/zitadel/pkg/grpc/settings"
)

func (s *Server) ListSMSProviders(ctx context.Context, req *admin_pb.ListSMSProvidersRequest) (*admin_pb.ListSMSProvidersResponse, error) {
	queries, err := listSMSConfigsToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchSMSConfigs(ctx, queries)
	if err != nil {
		return nil, err

	}
	return &admin_pb.ListSMSProvidersResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.Timestamp),
	}, nil
}

func (s *Server) GetSMSProvider(ctx context.Context, req *admin_pb.GetSMSProviderRequest) (*admin_pb.GetSMSProviderResponse, error) {
	result, err := s.query.SMSProviderConfigByID(ctx, req.Id)
	if err != nil {
		return nil, err

	}
	return &admin_pb.GetSMSProviderResponse{
		Config: &settings_pb.SMSProvider{
			Details: object.ToViewDetailsPb(result.Sequence, result.CreationDate, result.ChangeDate, result.ResourceOwner),
			Id:      result.ID,
			State:   smsStateToPb(result.State),
			Config:  SMSConfigToPb(result),
		},
	}, nil
}

func (s *Server) AddSMSProviderTwilio(ctx context.Context, req *admin_pb.AddSMSProviderTwilioRequest) (*admin_pb.AddSMSProviderTwilioResponse, error) {
	id, result, err := s.command.AddSMSConfigTwilio(ctx, authz.GetInstance(ctx).ID, AddSMSConfigTwilioToConfig(req))
	if err != nil {
		return nil, err

	}
	return &admin_pb.AddSMSProviderTwilioResponse{
		Details: object.DomainToAddDetailsPb(result),
		Id:      id,
	}, nil
}

func (s *Server) UpdateSMSProviderTwilio(ctx context.Context, req *admin_pb.UpdateSMSProviderTwilioRequest) (*admin_pb.UpdateSMSProviderTwilioResponse, error) {
	result, err := s.command.ChangeSMSConfigTwilio(ctx, authz.GetInstance(ctx).ID, req.Id, UpdateSMSConfigTwilioToConfig(req))
	if err != nil {
		return nil, err

	}
	return &admin_pb.UpdateSMSProviderTwilioResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) UpdateSMSProviderTwilioToken(ctx context.Context, req *admin_pb.UpdateSMSProviderTwilioTokenRequest) (*admin_pb.UpdateSMSProviderTwilioTokenResponse, error) {
	result, err := s.command.ChangeSMSConfigTwilioToken(ctx, authz.GetInstance(ctx).ID, req.Id, req.Token)
	if err != nil {
		return nil, err

	}
	return &admin_pb.UpdateSMSProviderTwilioTokenResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}
