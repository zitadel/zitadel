package admin

import (
	"context"

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
		Details: object.ToListDetails(result.Count, result.Sequence, result.LastRun),
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
	details, err := s.command.SetSecurityPolicy(ctx, securityPolicyToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.SetSecurityPolicyResponse{
		Details: object.DomainToChangeDetailsPb(details),
	}, nil
}
