package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
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
		Details: object.ToListDetails(result.Count, result.Sequence, result.LastRun),
		Result:  SMSConfigsToPb(result.Configs),
	}, nil
}

func (s *Server) GetSMSProvider(ctx context.Context, req *admin_pb.GetSMSProviderRequest) (*admin_pb.GetSMSProviderResponse, error) {
	result, err := s.query.SMSProviderConfigByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSMSProviderResponse{
		Config: SMSConfigToProviderPb(result),
	}, nil
}

func (s *Server) AddSMSProviderTwilio(ctx context.Context, req *admin_pb.AddSMSProviderTwilioRequest) (*admin_pb.AddSMSProviderTwilioResponse, error) {
	smsConfig := addSMSConfigTwilioToConfig(ctx, req)
	if err := s.command.AddSMSConfigTwilio(ctx, smsConfig); err != nil {
		return nil, err
	}
	return &admin_pb.AddSMSProviderTwilioResponse{
		Details: object.DomainToAddDetailsPb(smsConfig.Details),
		Id:      smsConfig.ID,
	}, nil
}

func (s *Server) UpdateSMSProviderTwilio(ctx context.Context, req *admin_pb.UpdateSMSProviderTwilioRequest) (*admin_pb.UpdateSMSProviderTwilioResponse, error) {
	smsConfig := updateSMSConfigTwilioToConfig(ctx, req)
	if err := s.command.ChangeSMSConfigTwilio(ctx, smsConfig); err != nil {
		return nil, err
	}
	return &admin_pb.UpdateSMSProviderTwilioResponse{
		Details: object.DomainToChangeDetailsPb(smsConfig.Details),
	}, nil
}

func (s *Server) UpdateSMSProviderTwilioToken(ctx context.Context, req *admin_pb.UpdateSMSProviderTwilioTokenRequest) (*admin_pb.UpdateSMSProviderTwilioTokenResponse, error) {
	result, err := s.command.ChangeSMSConfigTwilioToken(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, req.Token)
	if err != nil {
		return nil, err

	}
	return &admin_pb.UpdateSMSProviderTwilioTokenResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) AddSMSProviderHTTP(ctx context.Context, req *admin_pb.AddSMSProviderHTTPRequest) (*admin_pb.AddSMSProviderHTTPResponse, error) {
	smsConfig := addSMSConfigHTTPToConfig(ctx, req)
	if err := s.command.AddSMSConfigHTTP(ctx, smsConfig); err != nil {
		return nil, err
	}
	return &admin_pb.AddSMSProviderHTTPResponse{
		Details:    object.DomainToAddDetailsPb(smsConfig.Details),
		Id:         smsConfig.ID,
		SigningKey: smsConfig.SigningKey,
	}, nil
}

func (s *Server) UpdateSMSProviderHTTP(ctx context.Context, req *admin_pb.UpdateSMSProviderHTTPRequest) (*admin_pb.UpdateSMSProviderHTTPResponse, error) {
	smsConfig := updateSMSConfigHTTPToConfig(ctx, req)
	if err := s.command.ChangeSMSConfigHTTP(ctx, smsConfig); err != nil {
		return nil, err
	}
	return &admin_pb.UpdateSMSProviderHTTPResponse{
		Details:    object.DomainToChangeDetailsPb(smsConfig.Details),
		SigningKey: smsConfig.SigningKey,
	}, nil
}

func (s *Server) ActivateSMSProvider(ctx context.Context, req *admin_pb.ActivateSMSProviderRequest) (*admin_pb.ActivateSMSProviderResponse, error) {
	result, err := s.command.ActivateSMSConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ActivateSMSProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) DeactivateSMSProvider(ctx context.Context, req *admin_pb.DeactivateSMSProviderRequest) (*admin_pb.DeactivateSMSProviderResponse, error) {
	result, err := s.command.DeactivateSMSConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err

	}
	return &admin_pb.DeactivateSMSProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) RemoveSMSProvider(ctx context.Context, req *admin_pb.RemoveSMSProviderRequest) (*admin_pb.RemoveSMSProviderResponse, error) {
	result, err := s.command.RemoveSMSConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveSMSProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}
