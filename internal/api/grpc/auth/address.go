package auth

import (
	"context"

	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyAddress(ctx context.Context, _ *auth_pb.GetMyAddressRequest) (*auth_pb.GetMyAddressResponse, error) {
	address, err := s.repo.MyAddress(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyAddressResponse{
		Address: user_grpc.ModelAddressToPb(address),
		Details: object_grpc.ToDetailsPb(
			address.Sequence,
			address.CreationDate,
			address.ChangeDate,
			address.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateMyUserAddress(ctx context.Context, req *auth_pb.UpdateMyAddressRequest) (*auth_pb.UpdateMyAddressResponse, error) {
	address, err := s.command.ChangeHumanAddress(ctx, UpdateAddressToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.UpdateMyAddressResponse{
		Details: object_grpc.ToDetailsPb(
			address.Sequence,
			address.CreationDate,
			address.ChangeDate,
			address.ResourceOwner,
		),
	}, nil
}
