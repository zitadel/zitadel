package instance

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func (s *Server) AddCustomDomain(ctx context.Context, req *instance.AddCustomDomainRequest) (*instance.AddCustomDomainResponse, error) {
	details, err := s.command.AddInstanceDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &instance.AddCustomDomainResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) RemoveCustomDomain(ctx context.Context, req *instance.RemoveCustomDomainRequest) (*instance.RemoveCustomDomainResponse, error) {
	details, err := s.command.RemoveInstanceDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &instance.RemoveCustomDomainResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) AddTrustedDomain(ctx context.Context, req *instance.AddTrustedDomainRequest) (*instance.AddTrustedDomainResponse, error) {
	details, err := s.command.AddTrustedDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	return &instance.AddTrustedDomainResponse{
		Details: object.AddToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}

func (s *Server) RemoveTrustedDomain(ctx context.Context, req *instance.RemoveTrustedDomainRequest) (*instance.RemoveTrustedDomainResponse, error) {
	domain := strings.TrimSpace(req.GetDomain())
	err := validateParam(domain, "domain")
	if err != nil {
		return nil, err
	}

	details, err := s.command.RemoveTrustedDomain(ctx, domain)
	if err != nil {
		return nil, err
	}

	return &instance.RemoveTrustedDomainResponse{
		Details: object.ChangeToDetailsPb(details.Sequence, details.EventDate, details.ResourceOwner),
	}, nil
}
