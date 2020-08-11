package management

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

//TODO: should be in user
func (s *Server) GetServiceAccountById(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}

//TODO: should be in user
func (s *Server) IsServiceAccountUnique(ctx context.Context, in *management.ServiceAccountUniqueRequest) (*management.ServiceAccountUniqueResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}

//TODO: should be in user
func (s *Server) SerivceAccountChanges(ctx context.Context, in *management.ServiceAccountChangesRequest) (*management.ServiceAccountChangesResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}

func (s *Server) CreateServiceAccount(ctx context.Context, in *management.CreateServiceAccountRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.serviceAccount.CreateServiceAccount(ctx, createServiceAccountToModel(in))
	if err != nil {
		return nil, err
	}
	return serviceAccountFromModel(serviceAccount), nil
}

func (s *Server) UpdateServiceAccount(ctx context.Context, in *management.UpdateServiceAccountRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.serviceAccount.UpdateServiceAccount(ctx, updateServiceAccountToModel(in))
	if err != nil {
		return nil, err
	}
	return serviceAccountFromModel(serviceAccount), nil
}

func (s *Server) DeactivateServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.serviceAccount.DeactivateServiceAccount(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return serviceAccountFromModel(serviceAccount), nil
}

func (s *Server) ReactivateServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.serviceAccount.ReactivateServiceAccount(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return serviceAccountFromModel(serviceAccount), nil
}

func (s *Server) LockServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.serviceAccount.LockServiceAccount(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return serviceAccountFromModel(serviceAccount), nil
}

func (s *Server) UnlockServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.serviceAccount.UnlockServiceAccount(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return serviceAccountFromModel(serviceAccount), nil
}

func (s *Server) DeleteServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
