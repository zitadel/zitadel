package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetEmailProvider(ctx context.Context, req *admin_pb.GetEmailProviderRequest) (*admin_pb.GetEmailProviderResponse, error) {
	smtp, err := s.query.SMTPConfigActive(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetEmailProviderResponse{
		Config: emailProviderToProviderPb(smtp),
	}, nil
}

func (s *Server) GetEmailProviderById(ctx context.Context, req *admin_pb.GetEmailProviderByIdRequest) (*admin_pb.GetEmailProviderByIdResponse, error) {
	smtp, err := s.query.SMTPConfigByID(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetEmailProviderByIdResponse{
		Config: emailProviderToProviderPb(smtp),
	}, nil
}

func (s *Server) AddEmailProviderSMTP(ctx context.Context, req *admin_pb.AddEmailProviderSMTPRequest) (*admin_pb.AddEmailProviderSMTPResponse, error) {
	config := addEmailProviderSMTPToConfig(ctx, req)
	if err := s.command.AddSMTPConfig(ctx, config); err != nil {
		return nil, err
	}
	return &admin_pb.AddEmailProviderSMTPResponse{
		Details: object.DomainToChangeDetailsPb(config.Details),
		Id:      config.ID,
	}, nil
}

func (s *Server) UpdateEmailProviderSMTP(ctx context.Context, req *admin_pb.UpdateEmailProviderSMTPRequest) (*admin_pb.UpdateEmailProviderSMTPResponse, error) {
	config := updateEmailProviderSMTPToConfig(ctx, req)
	if err := s.command.ChangeSMTPConfig(ctx, config); err != nil {
		return nil, err
	}
	return &admin_pb.UpdateEmailProviderSMTPResponse{
		Details: object.DomainToChangeDetailsPb(config.Details),
	}, nil
}

func (s *Server) AddEmailProviderHTTP(ctx context.Context, req *admin_pb.AddEmailProviderHTTPRequest) (*admin_pb.AddEmailProviderHTTPResponse, error) {
	config := addEmailProviderHTTPToConfig(ctx, req)
	if err := s.command.AddSMTPConfigHTTP(ctx, config); err != nil {
		return nil, err
	}
	return &admin_pb.AddEmailProviderHTTPResponse{
		Details:    object.DomainToChangeDetailsPb(config.Details),
		Id:         config.ID,
		SigningKey: config.SigningKey,
	}, nil
}

func (s *Server) UpdateEmailProviderHTTP(ctx context.Context, req *admin_pb.UpdateEmailProviderHTTPRequest) (*admin_pb.UpdateEmailProviderHTTPResponse, error) {
	config := updateEmailProviderHTTPToConfig(ctx, req)
	if err := s.command.ChangeSMTPConfigHTTP(ctx, config); err != nil {
		return nil, err
	}
	return &admin_pb.UpdateEmailProviderHTTPResponse{
		Details:    object.DomainToChangeDetailsPb(config.Details),
		SigningKey: config.SigningKey,
	}, nil
}

func (s *Server) RemoveEmailProvider(ctx context.Context, req *admin_pb.RemoveEmailProviderRequest) (*admin_pb.RemoveEmailProviderResponse, error) {
	details, err := s.command.RemoveSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveEmailProviderResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) UpdateEmailProviderSMTPPassword(ctx context.Context, req *admin_pb.UpdateEmailProviderSMTPPasswordRequest) (*admin_pb.UpdateEmailProviderSMTPPasswordResponse, error) {
	details, err := s.command.ChangeSMTPConfigPassword(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, req.Password)
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateEmailProviderSMTPPasswordResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListEmailProviders(ctx context.Context, req *admin_pb.ListEmailProvidersRequest) (*admin_pb.ListEmailProvidersResponse, error) {
	queries, err := listEmailProvidersToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchSMTPConfigs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListEmailProvidersResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.LastRun),
		Result:  emailProvidersToPb(result.Configs),
	}, nil
}

func (s *Server) ActivateEmailProvider(ctx context.Context, req *admin_pb.ActivateEmailProviderRequest) (*admin_pb.ActivateEmailProviderResponse, error) {
	result, err := s.command.ActivateSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ActivateEmailProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) DeactivateEmailProvider(ctx context.Context, req *admin_pb.DeactivateEmailProviderRequest) (*admin_pb.DeactivateEmailProviderResponse, error) {
	result, err := s.command.DeactivateSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.DeactivateEmailProviderResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) TestEmailProviderById(ctx context.Context, req *admin_pb.TestEmailProviderSMTPByIdRequest) (*admin_pb.TestEmailProviderSMTPByIdResponse, error) {
	if err := s.command.TestSMTPConfigById(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, req.ReceiverAddress); err != nil {
		return nil, err
	}
	return &admin_pb.TestEmailProviderSMTPByIdResponse{}, nil
}

func (s *Server) TestEmailProviderSMTP(ctx context.Context, req *admin_pb.TestEmailProviderSMTPRequest) (*admin_pb.TestEmailProviderSMTPResponse, error) {
	if err := s.command.TestSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, req.ReceiverAddress, testEmailProviderSMTPToConfig(req)); err != nil {
		return nil, err
	}
	return &admin_pb.TestEmailProviderSMTPResponse{}, nil
}
