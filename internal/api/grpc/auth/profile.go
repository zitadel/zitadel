package auth

import (
	"context"

	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyProfile(ctx context.Context, req *auth_pb.GetMyProfileRequest) (*auth_pb.GetMyProfileResponse, error) {
	profile, err := s.repo.MyProfile(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyProfileResponse{
		Profile: user_grpc.ProfileToPb(profile),
		Details: object_grpc.ToDetailsPb(
			profile.Sequence,
			profile.ChangeDate,
			profile.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateMyProfile(ctx context.Context, req *auth_pb.UpdateMyProfileRequest) (*auth_pb.UpdateMyProfileResponse, error) {
	profile, err := s.command.ChangeHumanProfile(ctx, UpdateProfileToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &auth_pb.UpdateMyProfileResponse{
		Details: object_grpc.ToDetailsPb(
			profile.Sequence,
			profile.ChangeDate,
			profile.ResourceOwner,
		),
	}, nil
}
