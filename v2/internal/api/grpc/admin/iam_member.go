package admin

import (
	"context"
	"time"

	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/caos/zitadel/v2/internal/api/grpc/member"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func (s *Server) ListIAMMemberRoles(ctx context.Context, req *admin_pb.ListIAMMemberRolesRequest) (*admin_pb.ListIAMMemberRolesResponse, error) {
	roles := s.query.GetIAMMemberRoles()
	return &admin_pb.ListIAMMemberRolesResponse{
		Roles:   roles,
		Details: object.ToListDetails(uint64(len(roles)), 0, time.Now()),
	}, nil
}

func (s *Server) ListIAMMembers(ctx context.Context, req *admin_pb.ListIAMMembersRequest) (*admin_pb.ListIAMMembersResponse, error) {
	queries, err := ListIAMMembersRequestToQuery(req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.IAMMembers(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListIAMMembersResponse{
		Details: object.ToListDetails(res.Count, res.Sequence, res.Timestamp),
		//TODO: resource owner of user of the member instead of the membership resource owner
		Result: member.MembersToPb("", res.Members),
	}, nil
}

func (s *Server) AddIAMMember(ctx context.Context, req *admin_pb.AddIAMMemberRequest) (*admin_pb.AddIAMMemberResponse, error) {
	member, err := s.command.AddIAMMember(ctx, AddIAMMemberToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddIAMMemberResponse{
		Details: object.AddToDetailsPb(
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
		Details: object.ChangeToDetailsPb(
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
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
