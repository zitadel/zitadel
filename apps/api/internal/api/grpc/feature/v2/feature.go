package feature

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
)

func (s *Server) SetSystemFeatures(ctx context.Context, req *connect.Request[feature.SetSystemFeaturesRequest]) (_ *connect.Response[feature.SetSystemFeaturesResponse], err error) {
	features, err := systemFeaturesToCommand(req.Msg)
	if err != nil {
		return nil, err
	}
	details, err := s.command.SetSystemFeatures(ctx, features)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&feature.SetSystemFeaturesResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) ResetSystemFeatures(ctx context.Context, req *connect.Request[feature.ResetSystemFeaturesRequest]) (_ *connect.Response[feature.ResetSystemFeaturesResponse], err error) {
	details, err := s.command.ResetSystemFeatures(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&feature.ResetSystemFeaturesResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) GetSystemFeatures(ctx context.Context, req *connect.Request[feature.GetSystemFeaturesRequest]) (_ *connect.Response[feature.GetSystemFeaturesResponse], err error) {
	f, err := s.query.GetSystemFeatures(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(systemFeaturesToPb(f)), nil
}

func (s *Server) SetInstanceFeatures(ctx context.Context, req *connect.Request[feature.SetInstanceFeaturesRequest]) (_ *connect.Response[feature.SetInstanceFeaturesResponse], err error) {
	features, err := instanceFeaturesToCommand(req.Msg)
	if err != nil {
		return nil, err
	}
	details, err := s.command.SetInstanceFeatures(ctx, features)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&feature.SetInstanceFeaturesResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) ResetInstanceFeatures(ctx context.Context, req *connect.Request[feature.ResetInstanceFeaturesRequest]) (_ *connect.Response[feature.ResetInstanceFeaturesResponse], err error) {
	details, err := s.command.ResetInstanceFeatures(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&feature.ResetInstanceFeaturesResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) GetInstanceFeatures(ctx context.Context, req *connect.Request[feature.GetInstanceFeaturesRequest]) (_ *connect.Response[feature.GetInstanceFeaturesResponse], err error) {
	f, err := s.query.GetInstanceFeatures(ctx, req.Msg.GetInheritance())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(instanceFeaturesToPb(f)), nil
}

func (s *Server) SetOrganizationFeatures(ctx context.Context, req *connect.Request[feature.SetOrganizationFeaturesRequest]) (_ *connect.Response[feature.SetOrganizationFeaturesResponse], err error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOrganizationFeatures not implemented")
}
func (s *Server) ResetOrganizationFeatures(ctx context.Context, req *connect.Request[feature.ResetOrganizationFeaturesRequest]) (_ *connect.Response[feature.ResetOrganizationFeaturesResponse], err error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetOrganizationFeatures not implemented")
}
func (s *Server) GetOrganizationFeatures(ctx context.Context, req *connect.Request[feature.GetOrganizationFeaturesRequest]) (_ *connect.Response[feature.GetOrganizationFeaturesResponse], err error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrganizationFeatures not implemented")
}
func (s *Server) SetUserFeatures(ctx context.Context, req *connect.Request[feature.SetUserFeatureRequest]) (_ *connect.Response[feature.SetUserFeaturesResponse], err error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUserFeatures not implemented")
}
func (s *Server) ResetUserFeatures(ctx context.Context, req *connect.Request[feature.ResetUserFeaturesRequest]) (_ *connect.Response[feature.ResetUserFeaturesResponse], err error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetUserFeatures not implemented")
}
func (s *Server) GetUserFeatures(ctx context.Context, req *connect.Request[feature.GetUserFeaturesRequest]) (_ *connect.Response[feature.GetUserFeaturesResponse], err error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFeatures not implemented")
}
