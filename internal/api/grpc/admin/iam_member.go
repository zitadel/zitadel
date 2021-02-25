package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) AddIAMMember(ctx context.Context, req *admin_pb.AddIAMMemberRequest) (*admin_pb.AddIAMMemberResponse, error) {
	member, err := s.command.AddIAMMember(ctx, AddIAMMemberToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddIAMMemberResponse{
		Details: object.ToDetailsPb(
			member.Sequence,
			member.CreationDate,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateIAMMember(ctx context.Context, req *admin_pb.UpdateIAMMemberRequest) (*admin_pb.UpdateIAMMemberResponse, error) {
	member, err := s.command.ChangeIAMMember(ctx, UpdateIAMMemberToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateIAMMemberResponse{
		Details: object.ToDetailsPb(
			member.Sequence,
			member.CreationDate,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveIAMMember(ctx context.Context, req *admin_pb.RemoveIAMMemberRequest) (*admin_pb.RemoveIAMMemberResponse, error) {
	err := s.command.RemoveIAMMember(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveIAMMemberResponse{
		//TODO: return value
		// 	Details: object.ToDetailsPb(
		// 		member.Sequence,
		// 		member.CreationDate,
		// 		member.ChangeDate,
		// 		member.ResourceOwner,
		// 	),
	}, nil
}
