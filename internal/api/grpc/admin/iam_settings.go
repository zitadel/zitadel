package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) ListSecretGenerators(ctx context.Context, req *admin_pb.ListSecretGeneratorsRequest) (*admin_pb.ListSecretGeneratorsResponse, error) {
	queries, err := listSecretGeneratorToModel(req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchSecretGenerators(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListSecretGeneratorsResponse{
		Result:  SecretGeneratorsToPb(result.SecretGenerators),
		Details: object.ToListDetails(result.Count, result.Sequence, result.Timestamp),
	}, nil
}

func (s *Server) GetSecretGenerator(ctx context.Context, req *admin_pb.GetSecretGeneratorRequest) (*admin_pb.GetSecretGeneratorResponse, error) {
	generator, err := s.query.SecretGeneratorByType(ctx, SecretGeneratorTypeToDomain(req.GeneratorType))
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSecretGeneratorResponse{
		SecretGenerator: SecretGeneratorToPb(generator),
	}, nil
}

func (s *Server) UpdateSecretGenerator(ctx context.Context, req *admin_pb.UpdateSecretGeneratorRequest) (*admin_pb.UpdateSecretGeneratorResponse, error) {
	details, err := s.command.ChangeSecretGeneratorConfig(ctx, SecretGeneratorTypeToDomain(req.GeneratorType), UpdateSecretGeneratorToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateSecretGeneratorResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
	}, nil
}

func (s *Server) GetSMTPConfig(ctx context.Context, req *admin_pb.GetSMTPConfigRequest) (*admin_pb.GetSMTPConfigResponse, error) {
	smtp, err := s.query.SMTPConfigByAggregateID(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSMTPConfigResponse{
		SmtpConfig: SMTPConfigToPb(smtp),
	}, nil
}

func (s *Server) AddSMTPConfig(ctx context.Context, req *admin_pb.AddSMTPConfigRequest) (*admin_pb.AddSMTPConfigResponse, error) {
	details, err := s.command.AddSMTPConfig(ctx, AddSMTPToConfig(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddSMTPConfigResponse{
		Details: object.ChangeToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner),
	}, nil
}

func (s *Server) UpdateSMTPConfig(ctx context.Context, req *admin_pb.UpdateSMTPConfigRequest) (*admin_pb.UpdateSMTPConfigResponse, error) {
	details, err := s.command.ChangeSMTPConfig(ctx, UpdateSMTPToConfig(req))
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

func (s *Server) RemoveSMTPConfig(ctx context.Context, _ *admin_pb.RemoveSMTPConfigRequest) (*admin_pb.RemoveSMTPConfigResponse, error) {
	details, err := s.command.RemoveSMTPConfig(ctx)
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
	details, err := s.command.ChangeSMTPConfigPassword(ctx, req.Password)
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

func (s *Server) GetSecurityPolicy(ctx context.Context, req *admin_pb.GetSecurityPolicyRequest) (*admin_pb.GetSecurityPolicyResponse, error) {
	policy, err := s.query.SecurityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSecurityPolicyResponse{
		Policy: SecurityPolicyToPb(policy),
	}, nil
}

func (s *Server) SetSecurityPolicy(ctx context.Context, req *admin_pb.SetSecurityPolicyRequest) (*admin_pb.SetSecurityPolicyResponse, error) {
	details, err := s.command.SetSecurityPolicy(ctx, req.EnableIframeEmbedding, req.AllowedOrigins)
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetSecurityPolicyResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}
