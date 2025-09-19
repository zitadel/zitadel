package group

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/query"
	group "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

// ListGroups returns a list of groups that match the search criteria
func (s *Server) ListGroups(ctx context.Context, req *connect.Request[group.ListGroupsRequest]) (*connect.Response[group.ListGroupsResponse], error) {
	resp, err := s.query.SearchGroups(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group.ListGroupsResponse{
		Groups: groupsToPb(resp.Groups),
	}), nil
}

func groupsToPb(groups []*query.Group) []*group.Group {
	pbGroups := make([]*group.Group, len(groups))
	for i, g := range groups {
		pbGroups[i] = groupToPb(g)
	}
	return pbGroups
}

func groupToPb(g *query.Group) *group.Group {
	return &group.Group{
		Id:           g.ID,
		Name:         g.Name,
		CreationDate: timestamppb.New(g.CreationDate),
		ChangeDate:   timestamppb.New(g.ChangeDate),
	}
}
