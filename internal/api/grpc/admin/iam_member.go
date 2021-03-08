package admin

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/grpc/member"
	"github.com/caos/zitadel/internal/api/grpc/object"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) ListIAMMemberRoles(ctx context.Context, req *admin_pb.ListIAMMemberRolesRequest) (*admin_pb.ListIAMMemberRolesResponse, error) {
	roles := s.iam.GetIAMMemberRoles()
	return &admin_pb.ListIAMMemberRolesResponse{
		Details: object.ToListDetails(uint64(len(roles)), 0, time.Now()),
	}, nil
}

func (s *Server) ListIAMMembers(ctx context.Context, req *admin_pb.ListIAMMembersRequest) (*admin_pb.ListIAMMembersResponse, error) {
	res, err := s.iam.SearchIAMMembers(ctx, ListIAMMemberRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListIAMMembersResponse{
		Details: object.ToListDetails(res.TotalResult, res.Sequence, res.Timestamp),
		Result:  member.IAMMembersToPb(res.Result),
	}, nil
}

func (s *Server) AddIAMMember(ctx context.Context, req *admin_pb.AddIAMMemberRequest) (*admin_pb.AddIAMMemberResponse, error) {
	member, err := s.command.AddIAMMember(ctx, AddIAMMemberToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddIAMMemberResponse{
		Details: object.ToDetailsPb(
			member.Sequence,
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
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveIAMMember(ctx context.Context, req *admin_pb.RemoveIAMMemberRequest) (*admin_pb.RemoveIAMMemberResponse, error) {
	objectDetails, err := s.command.RemoveIAMMember(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveIAMMemberResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}
