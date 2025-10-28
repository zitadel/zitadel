package instance

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func (s *Server) AddCustomDomain(ctx context.Context, req *connect.Request[instance.AddCustomDomainRequest]) (*connect.Response[instance.AddCustomDomainResponse], error) {
	// Adding a custom domain is currently only allowed with system permissions,
	// so we directly check for them in the auth interceptor and do not check here again.
	details, err := s.command.AddInstanceDomain(ctx, req.Msg.GetCustomDomain())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&instance.AddCustomDomainResponse{
		CreationDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) RemoveCustomDomain(ctx context.Context, req *connect.Request[instance.RemoveCustomDomainRequest]) (*connect.Response[instance.RemoveCustomDomainResponse], error) {
	// Removing a custom domain is currently only allowed with system permissions,
	// so we directly check for them in the auth interceptor and do not check here again.
	details, err := s.command.RemoveInstanceDomain(ctx, req.Msg.GetCustomDomain())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&instance.RemoveCustomDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) AddTrustedDomain(ctx context.Context, req *connect.Request[instance.AddTrustedDomainRequest]) (*connect.Response[instance.AddTrustedDomainResponse], error) {
	if err := s.checkPermission(ctx, domain.PermissionSystemInstanceWrite, domain.PermissionInstanceWrite); err != nil {
		return nil, err
	}
	details, err := s.command.AddTrustedDomain(ctx, req.Msg.GetTrustedDomain())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&instance.AddTrustedDomainResponse{
		CreationDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) RemoveTrustedDomain(ctx context.Context, req *connect.Request[instance.RemoveTrustedDomainRequest]) (*connect.Response[instance.RemoveTrustedDomainResponse], error) {
	if err := s.checkPermission(ctx, domain.PermissionSystemInstanceWrite, domain.PermissionInstanceWrite); err != nil {
		return nil, err
	}
	details, err := s.command.RemoveTrustedDomain(ctx, req.Msg.GetTrustedDomain())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.RemoveTrustedDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}
