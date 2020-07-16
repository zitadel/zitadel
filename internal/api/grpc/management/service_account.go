package management

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetServiceAccountById(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) IsServiceAccountUnique(ctx context.Context, in *management.ServiceAccountUniqueRequest) (*management.ServiceAccountUniqueResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) CreateServiceAccount(ctx context.Context, in *management.CreateServiceAccountRequest) (*management.ServiceAccountResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) UpdateServiceAccount(ctx context.Context, in *management.UpdateServiceAccountRequest) (*management.ServiceAccountResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) DeactivateServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) ReactivateServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) LockServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) UnlockServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) DeleteServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
func (s *Server) SerivceAccountChanges(ctx context.Context, in *management.ServiceAccountChangesRequest) (*management.ServiceAccountChangesResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
