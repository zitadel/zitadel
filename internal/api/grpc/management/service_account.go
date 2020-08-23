package management

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) CreateServiceAccount(ctx context.Context, in *management.CreateServiceAccountRequest) (*management.ServiceAccountResponse, error) {
	user, err := s.user.CreateUser(ctx, createServiceAccountToUserModel(in))
	if err != nil {
		return nil, err
	}
	return serviceAccountFromUserModel(user), nil
}

func (s *Server) UpdateServiceAccount(ctx context.Context, in *management.UpdateServiceAccountRequest) (*management.ServiceAccountResponse, error) {
	// serviceAccount, err := s.user.User UpdateUser(ctx, updateServiceAccountToUserModel(in))
	// if err != nil {
	// 	return nil, err
	// }
	// return serviceAccountFromUserModel(serviceAccount), nil
	return nil, errors.ThrowUnimplemented(nil, "MANAG-7HlGy", "unimplemented")
}

func (s *Server) DeactivateServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.user.DeactivateUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return serviceAccountFromUserModel(serviceAccount), nil
}

func (s *Server) ReactivateServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.user.ReactivateUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return serviceAccountFromUserModel(serviceAccount), nil
}

func (s *Server) LockServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.user.LockUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return serviceAccountFromUserModel(serviceAccount), nil
}

func (s *Server) UnlockServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	serviceAccount, err := s.user.UnlockUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return serviceAccountFromUserModel(serviceAccount), nil
}

func (s *Server) DeleteServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
