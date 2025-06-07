package group

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	group_pb "github.com/zitadel/zitadel/pkg/grpc/group"
)

func GroupGrantsToPb(assetPrefix string, grants []*query.GroupGrant) []*group_pb.GroupGrant {
	u := make([]*group_pb.GroupGrant, len(grants))
	for i, grant := range grants {
		u[i] = GroupGrantToPb(assetPrefix, grant)
	}
	return u
}

func GroupGrantToPb(assetPrefix string, grant *query.GroupGrant) *group_pb.GroupGrant {
	return &group_pb.GroupGrant{
		Id:               grant.ID,
		GroupId:          grant.GroupID,
		State:            GroupGrantStateToPb(grant.State),
		RoleKeys:         grant.Roles,
		ProjectId:        grant.ProjectID,
		OrgId:            grant.ResourceOwner,
		ProjectGrantId:   grant.GrantID,
		OrgDomain:        grant.OrgPrimaryDomain,
		OrgName:          grant.OrgName,
		ProjectName:      grant.ProjectName,
		GrantedOrgId:     grant.GrantedOrgID,
		GrantedOrgName:   grant.GrantedOrgName,
		GrantedOrgDomain: grant.GrantedOrgDomain,
		GroupName:        grant.GroupName,
		Details: object.ToViewDetailsPb(
			grant.Sequence,
			grant.CreationDate,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}
}

func GroupGrantStateToPb(state domain.GroupGrantState) group_pb.GroupGrantState {
	switch state {
	case domain.GroupGrantStateActive:
		return group_pb.GroupGrantState_GROUP_GRANT_STATE_ACTIVE
	case domain.GroupGrantStateInactive:
		return group_pb.GroupGrantState_GROUP_GRANT_STATE_INACTIVE
	case domain.GroupGrantStateRemoved,
		domain.GroupGrantStateUnspecified:
		// these states should never occur here and are mainly listed for linting purposes
		fallthrough
	default:
		return group_pb.GroupGrantState_GROUP_GRANT_STATE_UNSPECIFIED
	}
}

func GroupGrantQueriesToQuery(ctx context.Context, queries []*group_pb.GroupGrantQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = GroupGrantQueryToQuery(ctx, query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func GroupGrantQueryToQuery(ctx context.Context, query *group_pb.GroupGrantQuery) (query.SearchQuery, error) {
	switch q := query.Query.(type) {
	case *group_pb.GroupGrantQuery_OrgDomainQuery:
		return GroupGrantOrgDomainQueryToModel(q.OrgDomainQuery)
	case *group_pb.GroupGrantQuery_OrgNameQuery:
		return GroupGrantOrgNameQueryToModel(q.OrgNameQuery)
	case *group_pb.GroupGrantQuery_ProjectGrantIdQuery:
		return GroupGrantProjectGrantIDQueryToModel(q.ProjectGrantIdQuery)
	case *group_pb.GroupGrantQuery_ProjectIdQuery:
		return GroupGrantProjectIDQueryToModel(q.ProjectIdQuery)
	case *group_pb.GroupGrantQuery_ProjectNameQuery:
		return GroupGrantProjectNameQueryToModel(q.ProjectNameQuery)
	case *group_pb.GroupGrantQuery_GroupNameQuery:
		return GroupGrantGroupNameQueryToModel(q.GroupNameQuery)
	case *group_pb.GroupGrantQuery_GroupIdQuery:
		return GroupGrantGroupIDQueryToModel(q.GroupIdQuery)
	case *group_pb.GroupGrantQuery_WithGrantedQuery:
		return GroupGrantWithGrantedQueryToModel(ctx, q.WithGrantedQuery)
	default:
		return nil, errors.New("invalid query")
	}
}

func GroupGrantGroupNameQueryToModel(q *group_pb.GroupGrantGroupNameQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantGroupNameQuery(q.GroupName, object.TextMethodToQuery(q.Method))
}

func GroupGrantOrgDomainQueryToModel(q *group_pb.GroupGrantOrgDomainQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantDomainQuery(q.OrgDomain, object.TextMethodToQuery(q.Method))
}

func GroupGrantOrgNameQueryToModel(q *group_pb.GroupGrantOrgNameQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantOrgNameQuery(q.OrgName, object.TextMethodToQuery(q.Method))
}

func GroupGrantProjectIDQueryToModel(q *group_pb.GroupGrantProjectIDQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantProjectIDSearchQuery(q.ProjectId)
}

func GroupGrantProjectGrantIDQueryToModel(q *group_pb.GroupGrantProjectGrantIDQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantGrantIDSearchQuery(q.ProjectGrantId)
}

func GroupGrantProjectNameQueryToModel(q *group_pb.GroupGrantProjectNameQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantProjectNameQuery(q.ProjectName, object.TextMethodToQuery(q.Method))
}

func GroupGrantRoleKeyQueryToModel(q *group_pb.GroupGrantRoleKeyQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantRoleQuery(q.RoleKey)
}

func GroupGrantGroupIDQueryToModel(q *group_pb.GroupGrantGroupIDQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantGroupIDSearchQuery(q.GroupId)
}

func GroupGrantWithGrantedQueryToModel(ctx context.Context, q *group_pb.GroupGrantWithGrantedQuery) (query.SearchQuery, error) {
	return query.NewGroupGrantWithGrantedQuery(authz.GetCtxData(ctx).OrgID)
}
