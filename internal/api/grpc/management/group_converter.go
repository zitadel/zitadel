package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	grp_grpc "github.com/zitadel/zitadel/internal/api/grpc/group"
	group_member_grpc "github.com/zitadel/zitadel/internal/api/grpc/groupmember"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/group"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func GroupCreateToDomain(req *mgmt_pb.AddGroupRequest) *domain.Group {
	return &domain.Group{
		Name:        req.Name,
		Description: req.Description,
	}
}

func GroupUpdateToDomain(req *mgmt_pb.UpdateGroupRequest) *domain.Group {
	return &domain.Group{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.Id,
		},
		Name:        req.Name,
		Description: req.Description,
	}
}

func GroupGrantsToIDs(groupGrants *query.GroupGrants) []string {
	converted := make([]string, len(groupGrants.GroupGrants))
	for i, grant := range groupGrants.GroupGrants {
		converted[i] = grant.GrantID
	}
	return converted
}

func AddGroupMemberRequestToDomain(ctx context.Context, req *mgmt_pb.AddGroupMemberRequest) *domain.Member {
	return domain.NewMember(req.GroupId, req.UserId)
}

func UpdateGroupMemberRequestToDomain(req *mgmt_pb.UpdateGroupMemberRequest) *domain.Member {
	return domain.NewMember(req.GroupId, req.UserId)
}

func listGroupRequestToModel(req *mgmt_pb.ListGroupsRequest) (*query.GroupSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := grp_grpc.GroupQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.GroupSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: GroupFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

// func listGrantedGroupsRequestToModel(req *mgmt_pb.ListGrantedGroupsRequest) (*query.GroupGrantSearchQueries, error) {
// 	offset, limit, asc := object.ListQueryToModel(req.Query)
// 	queries, err := grp_grpc.GroupQueriesToModel(req.Queries)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &query.ProjectGrantSearchQueries{
// 		SearchRequest: query.SearchRequest{
// 			Offset: offset,
// 			Limit:  limit,
// 			Asc:    asc,
// 		},
// 		Queries: queries,
// 	}, nil
// }

// func listProjectRolesRequestToModel(req *mgmt_pb.ListProjectRolesRequest) (*query.ProjectRoleSearchQueries, error) {
// 	offset, limit, asc := object.ListQueryToModel(req.Query)
// 	queries, err := proj_grpc.RoleQueriesToModel(req.Queries)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &query.ProjectRoleSearchQueries{
// 		SearchRequest: query.SearchRequest{
// 			Offset: offset,
// 			Limit:  limit,
// 			Asc:    asc,
// 		},
// 		Queries: queries,
// 	}, nil
// }

// func listGrantedProjectRolesRequestToModel(req *mgmt_pb.ListGrantedProjectRolesRequest) (*query.ProjectRoleSearchQueries, error) {
// 	offset, limit, asc := object.ListQueryToModel(req.Query)
// 	queries, err := proj_grpc.RoleQueriesToModel(req.Queries)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &query.ProjectRoleSearchQueries{
// 		SearchRequest: query.SearchRequest{
// 			Offset: offset,
// 			Limit:  limit,
// 			Asc:    asc,
// 		},
// 		Queries: queries,
// 	}, nil
// }

func ListGroupMembersRequestToModel(ctx context.Context, req *mgmt_pb.ListGroupMembersRequest) (*query.GroupMembersQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := group_member_grpc.MemberQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	ownerQuery, err := query.NewMemberResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, ownerQuery)
	return &query.GroupMembersQuery{
		MembersQuery: query.MembersQuery{
			SearchRequest: query.SearchRequest{
				Offset: offset,
				Limit:  limit,
				Asc:    asc,
				//SortingColumn: //TODO: sorting
			},
			Queries: queries,
		},
		GroupID: req.GroupId,
	}, nil
}

/*
func ListGroupMembershipsRequestToModel(ctx context.Context, req *mgmt_pb.ListGroupMembershipsRequest) (*query.MembershipSearchQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := user_grpc.MembershipQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	userQuery, err := query.NewMembershipUserIDQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	ownerQuery, err := query.NewMembershipResourceOwnersSearchQuery(authz.GetInstance(ctx).InstanceID(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, userQuery, ownerQuery)
	return &query.MembershipSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}
*/

func GroupFieldNameToSortingColumn(field group.GroupFieldName) query.Column {
	switch field {
	case group.GroupFieldName_GROUP_FIELD_NAME_NAME:
		return query.GroupColumnName
	case group.GroupFieldName_GROUP_FIELD_NAME_DESCRIPTION:
		return query.GroupColumnDescription
	case group.GroupFieldName_GROUP_FIELD_NAME_CREATION_DATE:
		return query.GroupColumnCreationDate
	default:
		return query.GroupColumnID
	}
}
