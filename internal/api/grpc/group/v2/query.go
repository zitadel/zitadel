package group

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

// ListGroups returns a list of groups that match the search criteria
func (s *Server) ListGroups(ctx context.Context, req *connect.Request[group_v2.ListGroupsRequest]) (*connect.Response[group_v2.ListGroupsResponse], error) {
	queries, err := s.listGroupsRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchGroups(ctx, queries)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&group_v2.ListGroupsResponse{
		Groups:     groupsToPb(resp.Groups),
		Pagination: filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}), nil
}

func (s *Server) listGroupsRequestToModel(req *group_v2.ListGroupsRequest) (*query.GroupSearchQuery, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}
	queries, err := groupSearchFiltersToQuery(req.GetFilters())
	if err != nil {
		return nil, err
	}
	return &query.GroupSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: groupFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func groupSearchFiltersToQuery(filters []*group_v2.GroupsSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(filters))
	for i, f := range filters {
		q[i], err = groupFilterToQuery(f)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func groupFilterToQuery(f *group_v2.GroupsSearchFilter) (query.SearchQuery, error) {
	switch q := f.Filter.(type) {
	case *group_v2.GroupsSearchFilter_GroupIds:
		return query.NewGroupIDsSearchQuery(q.GroupIds.GetIds())
	case *group_v2.GroupsSearchFilter_NameQuery:
		return query.NewGroupNameSearchQuery(q.NameQuery.GetName(), filter.TextMethodPbToQuery(q.NameQuery.GetMethod()))
	case *group_v2.GroupsSearchFilter_OrganizationId:
		return query.NewGroupOrganizationIdSearchQuery(q.OrganizationId.GetId())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-g3f4g", "List.Query.Invalid")
	}
}

func groupFieldNameToSortingColumn(field *group_v2.FieldName) query.Column {
	if field == nil {
		return query.GroupColumnCreationDate
	}
	switch *field {
	case group_v2.FieldName_FIELD_NAME_CREATION_DATE, group_v2.FieldName_FIELD_NAME_UNSPECIFIED:
		return query.GroupColumnCreationDate
	case group_v2.FieldName_FIELD_NAME_ID:
		return query.GroupColumnID
	case group_v2.FieldName_FIELD_NAME_NAME:
		return query.GroupColumnName
	case group_v2.FieldName_FIELD_NAME_CHANGE_DATE:
		return query.GroupColumnChangeDate
	default:
		return query.GroupColumnCreationDate
	}
}

func groupsToPb(groups []*query.Group) []*group_v2.Group {
	pbGroups := make([]*group_v2.Group, len(groups))
	for i, g := range groups {
		pbGroups[i] = groupToPb(g)
	}
	return pbGroups
}

func groupToPb(g *query.Group) *group_v2.Group {
	return &group_v2.Group{
		Id:           g.ID,
		Name:         g.Name,
		CreationDate: timestamppb.New(g.CreationDate),
		ChangeDate:   timestamppb.New(g.ChangeDate),
	}
}
