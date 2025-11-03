package instance

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func (s *Server) AddCustomDomain(ctx context.Context, req *connect.Request[instance.AddCustomDomainRequest]) (*connect.Response[instance.AddCustomDomainResponse], error) {
	details, err := s.command.AddInstanceDomain(ctx, req.Msg.GetDomain())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&instance.AddCustomDomainResponse{
		CreationDate: timestamppb.New(details.CreationDate),
	}), nil
}

func (s *Server) RemoveCustomDomain(ctx context.Context, req *connect.Request[instance.RemoveCustomDomainRequest]) (*connect.Response[instance.RemoveCustomDomainResponse], error) {
	details, err := s.command.RemoveInstanceDomain(ctx, req.Msg.GetDomain())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&instance.RemoveCustomDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) AddTrustedDomain(ctx context.Context, req *connect.Request[instance.AddTrustedDomainRequest]) (*connect.Response[instance.AddTrustedDomainResponse], error) {
	details, err := s.command.AddTrustedDomain(ctx, req.Msg.GetDomain())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&instance.AddTrustedDomainResponse{
		CreationDate: timestamppb.New(details.CreationDate),
	}), nil
}

func (s *Server) RemoveTrustedDomain(ctx context.Context, req *connect.Request[instance.RemoveTrustedDomainRequest]) (*connect.Response[instance.RemoveTrustedDomainResponse], error) {
	details, err := s.command.RemoveTrustedDomain(ctx, req.Msg.GetDomain())
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&instance.RemoveTrustedDomainResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}
