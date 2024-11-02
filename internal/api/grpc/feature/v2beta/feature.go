package feature

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	feature "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
)

func (s *Server) SetSystemFeatures(ctx context.Context, req *feature.SetSystemFeaturesRequest) (_ *feature.SetSystemFeaturesResponse, err error) {
	details, err := s.command.SetSystemFeatures(ctx, systemFeaturesToCommand(req))
	if err != nil {
		return nil, err
	}
	return &feature.SetSystemFeaturesResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ResetSystemFeatures(ctx context.Context, req *feature.ResetSystemFeaturesRequest) (_ *feature.ResetSystemFeaturesResponse, err error) {
	details, err := s.command.ResetSystemFeatures(ctx)
	if err != nil {
		return nil, err
	}
	return &feature.ResetSystemFeaturesResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) GetSystemFeatures(ctx context.Context, req *feature.GetSystemFeaturesRequest) (_ *feature.GetSystemFeaturesResponse, err error) {
	f, err := s.query.GetSystemFeatures(ctx)
	if err != nil {
		return nil, err
	}
	return systemFeaturesToPb(f), nil
}

func (s *Server) SetInstanceFeatures(ctx context.Context, req *feature.SetInstanceFeaturesRequest) (_ *feature.SetInstanceFeaturesResponse, err error) {
	details, err := s.command.SetInstanceFeatures(ctx, instanceFeaturesToCommand(req))
	if err != nil {
		return nil, err
	}
	return &feature.SetInstanceFeaturesResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ResetInstanceFeatures(ctx context.Context, req *feature.ResetInstanceFeaturesRequest) (_ *feature.ResetInstanceFeaturesResponse, err error) {
	details, err := s.command.ResetInstanceFeatures(ctx)
	if err != nil {
		return nil, err
	}
	return &feature.ResetInstanceFeaturesResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) GetInstanceFeatures(ctx context.Context, req *feature.GetInstanceFeaturesRequest) (_ *feature.GetInstanceFeaturesResponse, err error) {
	f, err := s.query.GetInstanceFeatures(ctx, req.GetInheritance())
	if err != nil {
		return nil, err
	}
	return instanceFeaturesToPb(f), nil
}

func (s *Server) SetOrganizationFeatures(ctx context.Context, req *feature.SetOrganizationFeaturesRequest) (_ *feature.SetOrganizationFeaturesResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOrganizationFeatures not implemented")
}
func (s *Server) ResetOrganizationFeatures(ctx context.Context, req *feature.ResetOrganizationFeaturesRequest) (_ *feature.ResetOrganizationFeaturesResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetOrganizationFeatures not implemented")
}
func (s *Server) GetOrganizationFeatures(ctx context.Context, req *feature.GetOrganizationFeaturesRequest) (_ *feature.GetOrganizationFeaturesResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrganizationFeatures not implemented")
}
func (s *Server) SetUserFeatures(ctx context.Context, req *feature.SetUserFeatureRequest) (_ *feature.SetUserFeaturesResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUserFeatures not implemented")
}
func (s *Server) ResetUserFeatures(ctx context.Context, req *feature.ResetUserFeaturesRequest) (_ *feature.ResetUserFeaturesResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetUserFeatures not implemented")
}
func (s *Server) GetUserFeatures(ctx context.Context, req *feature.GetUserFeaturesRequest) (_ *feature.GetUserFeaturesResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFeatures not implemented")
}
