package instance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func (s *Server) AddCustomDomain(ctx context.Context, req *instance.AddCustomDomainRequest) (*instance.AddCustomDomainResponse, error) {
	details, err := s.command.AddInstanceDomain(ctx, req.GetDomain())
	if err != nil {
		return nil, err
	}
	return &instance.AddCustomDomainResponse{
		CreationDate: timestamppb.New(details.CreationDate),
	}, nil
}

func (s *Server) RemoveCustomDomain(ctx context.Context, req *instance.RemoveCustomDomainRequest) (*instance.RemoveCustomDomainResponse, error) {
	details, err := s.command.RemoveInstanceDomain(ctx, req.GetDomain())
	if err != nil {
		return nil, err
	}
	return &instance.RemoveCustomDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}, nil
}

func (s *Server) AddTrustedDomain(ctx context.Context, req *instance.AddTrustedDomainRequest) (*instance.AddTrustedDomainResponse, error) {
	details, err := s.command.AddTrustedDomain(ctx, req.GetDomain())
	if err != nil {
		return nil, err
	}
	return &instance.AddTrustedDomainResponse{
		CreationDate: timestamppb.New(details.CreationDate),
	}, nil
}

func (s *Server) RemoveTrustedDomain(ctx context.Context, req *instance.RemoveTrustedDomainRequest) (*instance.RemoveTrustedDomainResponse, error) {
	details, err := s.command.RemoveTrustedDomain(ctx, req.GetDomain())
	if err != nil {
		return nil, err
	}

	return &instance.RemoveTrustedDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}, nil
}
