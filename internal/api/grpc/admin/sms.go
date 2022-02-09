package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) ListSMSProviderConfigs(ctx context.Context, req *admin_pb.ListSMSProviderConfigsRequest) (*admin_pb.ListSMSProviderConfigsResponse, error) {
	queries, err := listSMSConfigsToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchSMSConfigs(ctx, queries)
	if err != nil {
		return nil, err

	}
	return &admin_pb.ListSMSProviderConfigsResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.Timestamp),
	}, nil
}

func (s *Server) GetSMSProviderConfig(ctx context.Context, req *admin_pb.GetSMSProviderConfigRequest) (*admin_pb.GetSMSProviderConfigResponse, error) {
	result, err := s.query.SMSProviderConfigByID(ctx, req.Id)
	if err != nil {
		return nil, err

	}
	return &admin_pb.GetSMSProviderConfigResponse{
		Config: SMTPConfigToPb(result),
	}, nil
}

func (s *Server) AddSMSProviderConfigTwilio(ctx context.Context, req *admin_pb.AddSMSProviderConfigTwilioRequest) (*admin_pb.AddSMSProviderConfigTwilioResponse, error) {
	id, result, err := s.command.AddSMSConfigTwilio(ctx, AddSMSConfigTwilioToConfig(req))
	if err != nil {
		return nil, err

	}
	return &admin_pb.AddSMSProviderConfigTwilioResponse{
		Details: object.DomainToAddDetailsPb(result),
		Id:      id,
	}, nil
}

func (s *Server) UpdateSMSProviderConfigTwilio(ctx context.Context, req *admin_pb.UpdateSMSProviderConfigTwilioRequest) (*admin_pb.UpdateSMSProviderConfigTwilioResponse, error) {
	result, err := s.command.ChangeSMSConfigTwilio(ctx, req.Id, UpdateSMSConfigTwilioToConfig(req))
	if err != nil {
		return nil, err

	}
	return &admin_pb.UpdateSMSProviderConfigTwilioResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) UpdateSMSProviderConfigTwilioToken(ctx context.Context, req *admin_pb.UpdateSMSProviderConfigTwilioTokenRequest) (*admin_pb.UpdateSMSProviderConfigTwilioTokenResponse, error) {
	result, err := s.command.ChangeSMSConfigTwilioToken(ctx, req.Id, req.Token)
	if err != nil {
		return nil, err

	}
	return &admin_pb.UpdateSMSProviderConfigTwilioTokenResponse{
		Details: object.DomainToChangeDetailsPb(result),
	}, nil
}
