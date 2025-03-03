package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	change_grpc "github.com/zitadel/zitadel/internal/api/grpc/change"
	group_grpc "github.com/zitadel/zitadel/internal/api/grpc/group"
	groupmember_grpc "github.com/zitadel/zitadel/internal/api/grpc/groupmember"
	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetGroupByID(ctx context.Context, req *mgmt_pb.GetGroupByIDRequest) (*mgmt_pb.GetGroupByIDResponse, error) {
	group, err := s.query.GroupByID(ctx, true, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetGroupByIDResponse{
		Group: group_grpc.GroupViewToPb(group),
	}, nil
}

func (s *Server) ListGroups(ctx context.Context, req *mgmt_pb.ListGroupsRequest) (*mgmt_pb.ListGroupsResponse, error) {
	queries, err := listGroupRequestToModel(req)
	if err != nil {
		return nil, err
	}
	err = queries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	err = queries.AppendPermissionQueries(authz.GetRequestPermissionsFromCtx(ctx))
	if err != nil {
		return nil, err
	}
	groups, err := s.query.SearchGroups(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListGroupsResponse{
		Result:  group_grpc.GroupViewsToPb(groups.Groups),
		Details: object_grpc.ToListDetails(groups.Count, groups.Sequence, groups.LastRun),
	}, nil
}

func (s *Server) ListGroupGrantChanges(ctx context.Context, req *mgmt_pb.ListGroupGrantChangesRequest) (*mgmt_pb.ListGroupGrantChangesResponse, error) {
	var (
		limit    uint64
		sequence uint64
		asc      bool
	)
	if req.Query != nil {
		limit = uint64(req.Query.Limit)
		sequence = req.Query.Sequence
		asc = req.Query.Asc
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		Limit(limit).
		OrderDesc().
		ResourceOwner(authz.GetCtxData(ctx).OrgID).
		AwaitOpenTransactions().
		SequenceGreater(sequence).
		AddQuery().
		AggregateTypes(groupgrant.AggregateType).
		AggregateIDs(req.GrantId).
		Builder()
	if asc {
		query.OrderAsc()
	}

	changes, err := s.query.SearchEvents(ctx, query)
	if err != nil {
		return nil, err
	}

	return &mgmt_pb.ListGroupGrantChangesResponse{
		Result: change_grpc.EventsToChangesPb(changes, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) ListGroupChanges(ctx context.Context, req *mgmt_pb.ListGroupChangesRequest) (*mgmt_pb.ListGroupChangesResponse, error) {
	var (
		limit    uint64
		sequence uint64
		asc      bool
	)
	if req.Query != nil {
		limit = uint64(req.Query.Limit)
		sequence = req.Query.Sequence
		asc = req.Query.Asc
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		Limit(limit).
		AwaitOpenTransactions().
		OrderDesc().
		ResourceOwner(authz.GetCtxData(ctx).OrgID).
		SequenceGreater(sequence).
		AddQuery().
		AggregateTypes(group.AggregateType).
		AggregateIDs(req.GroupId).
		Builder()
	if asc {
		query.OrderAsc()
	}

	changes, err := s.query.SearchEvents(ctx, query)
	if err != nil {
		return nil, err
	}

	return &mgmt_pb.ListGroupChangesResponse{
		Result: change_grpc.EventsToChangesPb(changes, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) AddGroup(ctx context.Context, req *mgmt_pb.AddGroupRequest) (*mgmt_pb.AddGroupResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	group, err := s.command.AddGroup(ctx, GroupCreateToDomain(req), ctxData.OrgID, ctxData.UserID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGroupResponse{
		Id:      group.AggregateID,
		Details: object_grpc.AddToDetailsPb(group.Sequence, group.ChangeDate, group.ResourceOwner),
	}, nil
}

func (s *Server) UpdateGroup(ctx context.Context, req *mgmt_pb.UpdateGroupRequest) (*mgmt_pb.UpdateGroupResponse, error) {
	group, err := s.command.ChangeGroup(ctx, GroupUpdateToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGroupResponse{
		Details: object_grpc.ChangeToDetailsPb(
			group.Sequence,
			group.ChangeDate,
			group.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateGroup(ctx context.Context, req *mgmt_pb.DeactivateGroupRequest) (*mgmt_pb.DeactivateGroupResponse, error) {
	details, err := s.command.DeactivateGroup(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateGroupResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateGroup(ctx context.Context, req *mgmt_pb.ReactivateGroupRequest) (*mgmt_pb.ReactivateGroupResponse, error) {
	details, err := s.command.ReactivateGroup(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateGroupResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) removeGroupDependencies(ctx context.Context, groupID string) ([]*command.CascadingMembership, []string, error) {
	groupGrantGroupQuery, err := query.NewGroupGrantGroupIDSearchQuery(groupID)
	if err != nil {
		return nil, nil, err
	}
	grants, err := s.query.GroupGrants(ctx, &query.GroupGrantsQueries{
		Queries: []query.SearchQuery{groupGrantGroupQuery},
	}, true)
	if err != nil {
		return nil, nil, err
	}
	membershipsGroupQuery, err := query.NewGroupMembershipGroupIDQuery(groupID)
	if err != nil {
		return nil, nil, err
	}
	memberships, err := s.query.Memberships(ctx, &query.MembershipSearchQuery{
		Queries: []query.SearchQuery{membershipsGroupQuery},
	}, false)
	if err != nil {
		return nil, nil, err
	}
	return cascadingMemberships(memberships.Memberships), groupGrantsToIDs(grants.GroupGrants), nil
}

func (s *Server) RemoveGroup(ctx context.Context, req *mgmt_pb.RemoveGroupRequest) (*mgmt_pb.RemoveGroupResponse, error) {
	projectQuery, err := query.NewGroupGrantProjectIDSearchQuery(req.Id)
	if err != nil {
		return nil, err
	}
	grants, err := s.query.GroupGrants(ctx, &query.GroupGrantsQueries{
		Queries: []query.SearchQuery{projectQuery},
	}, true)
	if err != nil {
		return nil, err
	}
	details, err := s.command.RemoveGroup(ctx, req.Id, authz.GetCtxData(ctx).OrgID, groupGrantsToIDs(grants.GroupGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveGroupResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListGroupMembers(ctx context.Context, req *mgmt_pb.ListGroupMembersRequest) (*mgmt_pb.ListGroupMembersResponse, error) {
	queries, err := ListGroupMembersRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	members, err := s.query.GroupMembers(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListGroupMembersResponse{
		Result:  groupmember_grpc.MembersToPb(s.assetAPIPrefix(ctx), members.GroupMembers),
		Details: object_grpc.ToListDetails(members.Count, members.Sequence, members.LastRun),
	}, nil
}

func (s *Server) AddGroupMember(ctx context.Context, req *mgmt_pb.AddGroupMemberRequest) (*mgmt_pb.AddGroupMemberResponse, error) {
	member, err := s.command.AddGroupMember(ctx, AddGroupMemberRequestToDomain(ctx, req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGroupMemberResponse{
		Details: object_grpc.AddToDetailsPb(member.Sequence, member.ChangeDate, member.ResourceOwner),
	}, nil
}

func (s *Server) UpdateGroupMember(ctx context.Context, req *mgmt_pb.UpdateGroupMemberRequest) (*mgmt_pb.UpdateGroupMemberResponse, error) {
	member, err := s.command.ChangeGroupMember(ctx, UpdateGroupMemberRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGroupMemberResponse{
		Details: object_grpc.ChangeToDetailsPb(
			member.Sequence,
			member.ChangeDate,
			member.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveGroupMember(ctx context.Context, req *mgmt_pb.RemoveGroupMemberRequest) (*mgmt_pb.RemoveGroupMemberResponse, error) {
	details, err := s.command.RemoveGroupMember(ctx, req.GroupId, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveGroupMemberResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func groupGrantsToIDs(groupGrants []*query.GroupGrant) []string {
	converted := make([]string, len(groupGrants))
	for i, grant := range groupGrants {
		converted[i] = grant.ID
	}
	return converted
}
