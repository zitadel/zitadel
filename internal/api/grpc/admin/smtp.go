package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetSMTPConfig(ctx context.Context, req *admin_pb.GetSMTPConfigRequest) (*admin_pb.GetSMTPConfigResponse, error) {
	smtp, err := s.query.SMTPConfigActive(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSMTPConfigResponse{
		SmtpConfig: SMTPConfigToPb(smtp),
	}, nil
}

func (s *Server) GetSMTPConfigById(ctx context.Context, req *admin_pb.GetSMTPConfigByIdRequest) (*admin_pb.GetSMTPConfigByIdResponse, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	resourceOwner := instanceID // Will be replaced when orgs have smtp configs

	smtp, err := s.query.SMTPConfigByID(ctx, instanceID, resourceOwner, req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSMTPConfigByIdResponse{
		SmtpConfig: SMTPConfigToPb(smtp),
	}, nil
}

func (s *Server) AddSMTPConfig(ctx context.Context, req *admin_pb.AddSMTPConfigRequest) (*admin_pb.AddSMTPConfigResponse, error) {
	id, details, err := s.command.AddSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), AddSMTPToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddSMTPConfigResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
		Id: id,
	}, nil
}

func (s *Server) UpdateSMTPConfig(ctx context.Context, req *admin_pb.UpdateSMTPConfigRequest) (*admin_pb.UpdateSMTPConfigResponse, error) {
	details, err := s.command.ChangeSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, UpdateSMTPToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateSMTPConfigResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
	}, nil
}

func (s *Server) RemoveSMTPConfig(ctx context.Context, req *admin_pb.RemoveSMTPConfigRequest) (*admin_pb.RemoveSMTPConfigResponse, error) {
	details, err := s.command.RemoveSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveSMTPConfigResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
	}, nil
}

func (s *Server) UpdateSMTPConfigPassword(ctx context.Context, req *admin_pb.UpdateSMTPConfigPasswordRequest) (*admin_pb.UpdateSMTPConfigPasswordResponse, error) {
	details, err := s.command.ChangeSMTPConfigPassword(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, req.Password)
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateSMTPConfigPasswordResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
	}, nil
}

func (s *Server) ListSMTPConfigs(ctx context.Context, req *admin_pb.ListSMTPConfigsRequest) (*admin_pb.ListSMTPConfigsResponse, error) {
	queries, err := listSMTPConfigsToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchSMTPConfigs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListSMTPConfigsResponse{
		Details: object.ToListDetails(result.Count, result.Sequence, result.LastRun),
		Result:  SMTPConfigsToPb(result.Configs),
	}, nil
}

func (s *Server) ActivateSMTPConfig(ctx context.Context, req *admin_pb.ActivateSMTPConfigRequest) (*admin_pb.ActivateSMTPConfigResponse, error) {
	// Get the ID of current SMTP active provider if any
	currentActiveProviderID := ""
	smtp, err := s.query.SMTPConfigActive(ctx, authz.GetInstance(ctx).InstanceID())
	if err == nil {
		currentActiveProviderID = smtp.ID
	}

	result, err := s.command.ActivateSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, currentActiveProviderID)
	if err != nil {
		return nil, err

	}

	return &admin_pb.ActivateSMTPConfigResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) DeactivateSMTPConfig(ctx context.Context, req *admin_pb.DeactivateSMTPConfigRequest) (*admin_pb.DeactivateSMTPConfigResponse, error) {
	result, err := s.command.DeactivateSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id)
	if err != nil {
		return nil, err

	}
	return &admin_pb.DeactivateSMTPConfigResponse{
		Details: object.DomainToAddDetailsPb(result),
	}, nil
}

func (s *Server) TestSMTPConfigById(ctx context.Context, req *admin_pb.TestSMTPConfigByIdRequest) (*admin_pb.TestSMTPConfigByIdResponse, error) {
	err := s.command.TestSMTPConfigById(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, req.ReceiverAddress)
	if err != nil {
		return nil, err
	}

	return &admin_pb.TestSMTPConfigByIdResponse{}, nil
}

func (s *Server) TestSMTPConfig(ctx context.Context, req *admin_pb.TestSMTPConfigRequest) (*admin_pb.TestSMTPConfigResponse, error) {
	config := smtp.Config{}
	config.Tls = req.Tls
	config.From = req.SenderAddress
	config.FromName = req.SenderName
	config.SMTP.Host = req.Host
	config.SMTP.User = req.User
	config.SMTP.Password = req.Password

	err := s.command.TestSMTPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.Id, req.ReceiverAddress, &config)
	if err != nil {
		return nil, err
	}

	return &admin_pb.TestSMTPConfigResponse{}, nil
}
