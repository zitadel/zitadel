package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	change_grpc "github.com/zitadel/zitadel/internal/api/grpc/change"
	group_grpc "github.com/zitadel/zitadel/internal/api/grpc/group"
	groupuser_grpc "github.com/zitadel/zitadel/internal/api/grpc/groupuser"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata"
	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/zerrors"
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

func (s *Server) GetGroupByUserID(ctx context.Context, req *mgmt_pb.GetGroupByUserIDRequest) (*mgmt_pb.ListUserGroupsResponse, error) {
	groups, err := s.query.GroupByUserID(ctx, false, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListUserGroupsResponse{
		Result:  group_grpc.GroupViewsToPb(groups.Groups),
		Details: object_grpc.ToListDetails(groups.Count, groups.Sequence, groups.LastRun),
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
		// AllowTimeTravel().
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
		// AllowTimeTravel().
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

func (s *Server) RemoveGroupDependencies(ctx context.Context, groupID string) ([]*command.CascadingMembership, []string, error) {
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
	groupUser, err := s.query.GroupUsers(ctx, &query.GroupUsersQuery{
		GroupID: req.Id,
	})
	if err != nil {
		return nil, zerrors.ThrowInvalidArgumentf(err, "GROUP-IDasq", "GroupMember %v", groupUser)
	}
	details, err := s.command.RemoveGroup(ctx, req.Id, authz.GetCtxData(ctx).OrgID, memberToUserIDs(groupUser.GroupUsers), groupGrantsToIDs(grants.GroupGrants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveGroupResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ListGroupUsers(ctx context.Context, req *mgmt_pb.ListGroupUsersRequest) (*mgmt_pb.ListGroupUsersResponse, error) {
	queries, err := ListGroupUsersRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	members, err := s.query.GroupUsers(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListGroupUsersResponse{
		Result:  groupuser_grpc.MembersToPb(s.assetAPIPrefix(ctx), members.GroupUsers),
		Details: object_grpc.ToListDetails(members.Count, members.Sequence, members.LastRun),
	}, nil
}

func (s *Server) AddGroupUser(ctx context.Context, req *mgmt_pb.AddGroupUserRequest) (*mgmt_pb.AddGroupUserResponse, error) {
	groupuser, err := s.command.AddGroupUser(ctx, AddGroupUserRequestToDomain(ctx, req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGroupUserResponse{
		Details: object_grpc.AddToDetailsPb(groupuser.Sequence, groupuser.ChangeDate, groupuser.ResourceOwner),
	}, nil
}

func (s *Server) UpdateGroupUser(ctx context.Context, req *mgmt_pb.UpdateGroupUserRequest) (*mgmt_pb.UpdateGroupUserResponse, error) {
	groupuser, err := s.command.ChangeGroupUser(ctx, UpdateGroupUserRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGroupUserResponse{
		Details: object_grpc.ChangeToDetailsPb(
			groupuser.Sequence,
			groupuser.ChangeDate,
			groupuser.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveGroupUser(ctx context.Context, req *mgmt_pb.RemoveGroupUserRequest) (*mgmt_pb.RemoveGroupUserResponse, error) {
	details, err := s.command.RemoveGroupUser(ctx, req.GroupId, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveGroupUserResponse{
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

func memberToUserIDs(groupUser []*query.GroupUser) []string {
	converted := make([]string, len(groupUser))
	for i, group := range groupUser {
		converted[i] = group.UserID
	}
	return converted
}

// Group Metadata

func (s *Server) ListGroupMetadata(ctx context.Context, req *mgmt_pb.ListGroupMetadataRequest) (*mgmt_pb.ListGroupMetadataResponse, error) {
	metadataQueries, err := ListGroupMetadataToDomain(req)
	if err != nil {
		return nil, err
	}
	err = metadataQueries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchGroupMetadata(ctx, true, req.Id, metadataQueries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListGroupMetadataResponse{
		Result:  metadata.GroupMetadataListToPb(res.Metadata),
		Details: obj_grpc.ToListDetails(res.Count, res.Sequence, res.LastRun),
	}, nil
}

func (s *Server) GetGroupMetadata(ctx context.Context, req *mgmt_pb.GetGroupMetadataRequest) (*mgmt_pb.GetGroupMetadataResponse, error) {
	owner, err := query.NewGroupMetadataResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	data, err := s.query.GetGroupMetadataByKey(ctx, true, req.Id, req.Key, false, owner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetGroupMetadataResponse{
		Metadata: metadata.GroupMetadataToPb(data),
	}, nil
}

func (s *Server) SetGroupMetadata(ctx context.Context, req *mgmt_pb.SetGroupMetadataRequest) (*mgmt_pb.SetGroupMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.SetGroupMetadata(ctx, &domain.Metadata{Key: req.Key, Value: req.Value}, req.Id, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetGroupMetadataResponse{
		Details: obj_grpc.AddToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) BulkSetGroupMetadata(ctx context.Context, req *mgmt_pb.BulkSetGroupMetadataRequest) (*mgmt_pb.BulkSetGroupMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.BulkSetGroupMetadata(ctx, req.Id, ctxData.OrgID, BulkSetGroupMetadataToDomain(req)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.BulkSetGroupMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) RemoveGroupMetadata(ctx context.Context, req *mgmt_pb.RemoveGroupMetadataRequest) (*mgmt_pb.RemoveGroupMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.RemoveGroupMetadata(ctx, req.Key, req.Id, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveGroupMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) BulkRemoveGroupMetadata(ctx context.Context, req *mgmt_pb.BulkRemoveGroupMetadataRequest) (*mgmt_pb.BulkRemoveGroupMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.BulkRemoveGroupMetadata(ctx, req.Id, ctxData.OrgID, req.Keys...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.BulkRemoveGroupMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}
