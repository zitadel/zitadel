package group

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

// CreateGroupGrant authorizes all members of a group for a project with the given roles
func (s *Server) CreateGroupGrant(ctx context.Context, req *connect.Request[group_v2.CreateGroupGrantRequest]) (*connect.Response[group_v2.CreateGroupGrantResponse], error) {
	details, err := s.command.AddGroupGrant(ctx, &command.AddGroupGrant{
		GroupID:        req.Msg.GetGroupId(),
		ProjectID:      req.Msg.GetProjectId(),
		ProjectGrantID: req.Msg.GetProjectGrantId(),
		RoleKeys:       req.Msg.GetRoleKeys(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.CreateGroupGrantResponse{
		Id:           details.ID,
		CreationDate: timestamppb.New(details.EventDate),
	}), nil
}

// UpdateGroupGrant updates the roles of a group grant
func (s *Server) UpdateGroupGrant(ctx context.Context, req *connect.Request[group_v2.UpdateGroupGrantRequest]) (*connect.Response[group_v2.UpdateGroupGrantResponse], error) {
	details, err := s.command.ChangeGroupGrant(ctx, req.Msg.GetId(), req.Msg.GetRoleKeys())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.UpdateGroupGrantResponse{
		ChangeDate: timestamppb.New(details.EventDate),
	}), nil
}

// DeleteGroupGrant deletes a group grant
func (s *Server) DeleteGroupGrant(ctx context.Context, req *connect.Request[group_v2.DeleteGroupGrantRequest]) (*connect.Response[group_v2.DeleteGroupGrantResponse], error) {
	details, err := s.command.RemoveGroupGrant(ctx, req.Msg.GetId())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.DeleteGroupGrantResponse{
		DeletionDate: timestamppb.New(details.EventDate),
	}), nil
}

// ListGroupGrants returns the group grants matching the search criteria
func (s *Server) ListGroupGrants(ctx context.Context, req *connect.Request[group_v2.ListGroupGrantsRequest]) (*connect.Response[group_v2.ListGroupGrantsResponse], error) {
	queries, err := listGroupGrantsRequestToModel(req.Msg, s.systemDefaults)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchGroupGrants(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.ListGroupGrantsResponse{
		GroupGrants: groupGrantsToPb(resp.GroupGrants),
		Pagination:  filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}), nil
}

func listGroupGrantsRequestToModel(req *group_v2.ListGroupGrantsRequest, systemDefaults systemdefaults.SystemDefaults) (*query.GroupGrantsSearchQuery, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(systemDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}
	queries, err := groupGrantsSearchFiltersToQuery(req.GetFilters())
	if err != nil {
		return nil, err
	}
	return &query.GroupGrantsSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: query.GroupGrantColumnCreationDate,
		},
		Queries: queries,
	}, nil
}

func groupGrantsSearchFiltersToQuery(filters []*group_v2.GroupGrantsSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(filters))
	for i, f := range filters {
		q[i], err = groupGrantFilterToQuery(f)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func groupGrantFilterToQuery(f *group_v2.GroupGrantsSearchFilter) (query.SearchQuery, error) {
	switch q := f.Filter.(type) {
	case *group_v2.GroupGrantsSearchFilter_GroupIds:
		return query.NewGroupGrantGroupIDsSearchQuery(q.GroupIds.GetIds())
	case *group_v2.GroupGrantsSearchFilter_ProjectId:
		return query.NewGroupGrantProjectIDSearchQuery(q.ProjectId.GetId())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-x4Fmq2", "List.Query.Invalid")
	}
}

func groupGrantsToPb(grants []*query.GroupGrant) []*group_v2.GroupGrant {
	pbGrants := make([]*group_v2.GroupGrant, len(grants))
	for i, grant := range grants {
		pbGrants[i] = &group_v2.GroupGrant{
			Id:             grant.ID,
			GroupId:        grant.GroupID,
			GroupName:      grant.GroupName,
			OrganizationId: grant.ResourceOwner,
			ProjectId:      grant.ProjectID,
			ProjectGrantId: grant.ProjectGrantID,
			RoleKeys:       grant.RoleKeys,
			CreationDate:   timestamppb.New(grant.CreationDate),
			ChangeDate:     timestamppb.New(grant.ChangeDate),
		}
	}
	return pbGrants
}
