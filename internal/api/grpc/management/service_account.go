package management

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) CreateServiceAccount(ctx context.Context, in *management.CreateServiceAccountRequest) (*management.ServiceAccountResponse, error) {
	// serviceAccount, err := s.serviceAccount.CreateServiceAccount(ctx, createServiceAccountToModel(in))
	// if err != nil {
	// 	return nil, err
	// }
	// return serviceAccountFromModel(serviceAccount), nil
	return nil, errors.ThrowUnimplemented(nil, "MANAG-kzeI6", "unimplemented")
}

func (s *Server) UpdateServiceAccount(ctx context.Context, in *management.UpdateServiceAccountRequest) (*management.ServiceAccountResponse, error) {
	// serviceAccount, err := s.serviceAccount.UpdateServiceAccount(ctx, updateServiceAccountToModel(in))
	// if err != nil {
	// 	return nil, err
	// }
	// return serviceAccountFromModel(serviceAccount), nil
	return nil, errors.ThrowUnimplemented(nil, "MANAG-7HlGy", "unimplemented")
}

func (s *Server) DeactivateServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	// serviceAccount, err := s.serviceAccount.DeactivateServiceAccount(ctx, in.Id)
	// if err != nil {
	// 	return nil, err
	// }
	// return serviceAccountFromModel(serviceAccount), nil
	return nil, errors.ThrowUnimplemented(nil, "MANAG-ByopO", "unimplemented")
}

func (s *Server) ReactivateServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	// serviceAccount, err := s.serviceAccount.ReactivateServiceAccount(ctx, in.Id)
	// if err != nil {
	// 	return nil, err
	// }
	// return serviceAccountFromModel(serviceAccount), nil
	return nil, errors.ThrowUnimplemented(nil, "MANAG-L58LU", "unimplemented")
}

func (s *Server) LockServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	// serviceAccount, err := s.serviceAccount.LockServiceAccount(ctx, in.Id)
	// if err != nil {
	// 	return nil, err
	// }
	// return serviceAccountFromModel(serviceAccount), nil
	return nil, errors.ThrowUnimplemented(nil, "MANAG-lZ21v", "unimplemented")
}

func (s *Server) UnlockServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*management.ServiceAccountResponse, error) {
	// serviceAccount, err := s.serviceAccount.UnlockServiceAccount(ctx, in.Id)
	// if err != nil {
	// 	return nil, err
	// }
	// return serviceAccountFromModel(serviceAccount), nil
	return nil, errors.ThrowUnimplemented(nil, "MANAG-4col1", "unimplemented")
}

func (s *Server) DeleteServiceAccount(ctx context.Context, in *management.ServiceAccountIDRequest) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "ID", "Errors.*")
}
