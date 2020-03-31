package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) SearchApplications(ctx context.Context, request *ApplicationSearchRequest) (*ApplicationSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-yW23f", "Not implemented")
}

func (s *Server) ApplicationByID(ctx context.Context, request *ApplicationID) (*Application, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-bmr6X", "Not implemented")
}

func (s *Server) CreateOIDCApplication(ctx context.Context, in *OIDCApplicationCreate) (*Application, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-poe4d", "Not implemented")
}
func (s *Server) UpdateApplication(ctx context.Context, in *ApplicationUpdate) (*Application, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-bmt6J", "Not implemented")
}
func (s *Server) DeactivateApplication(ctx context.Context, in *ApplicationID) (*Application, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-cD34f", "Not implemented")
}
func (s *Server) ReactivateApplication(ctx context.Context, in *ApplicationID) (*Application, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mo98S", "Not implemented")
}
func (s *Server) UpdateApplicationOIDCConfig(ctx context.Context, in *OIDCConfigUpdate) (*OIDCConfig, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-xm56g", "Not implemented")
}
func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, in *ApplicationID) (*ClientSecret, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dlwp3", "Not implemented")
}

func (s *Server) ApplicationChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-due45", "Not implemented")
}
