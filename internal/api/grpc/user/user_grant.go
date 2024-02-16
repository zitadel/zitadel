package user

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
)

func UserGrantsToPb(assetPrefix string, grants []*query.UserGrant) []*user_pb.UserGrant {
	u := make([]*user_pb.UserGrant, len(grants))
	for i, grant := range grants {
		u[i] = UserGrantToPb(assetPrefix, grant)
	}
	return u
}

func UserGrantToPb(assetPrefix string, grant *query.UserGrant) *user_pb.UserGrant {
	return &user_pb.UserGrant{
		Id:                 grant.ID,
		UserId:             grant.UserID,
		State:              user_pb.UserGrantState_USER_GRANT_STATE_ACTIVE,
		RoleKeys:           grant.Roles,
		ProjectId:          grant.ProjectID,
		OrgId:              grant.ResourceOwner,
		ProjectGrantId:     grant.GrantID,
		UserName:           grant.Username,
		FirstName:          grant.FirstName,
		LastName:           grant.LastName,
		Email:              grant.Email,
		DisplayName:        grant.DisplayName,
		OrgDomain:          grant.OrgPrimaryDomain,
		OrgName:            grant.OrgName,
		ProjectName:        grant.ProjectName,
		AvatarUrl:          domain.AvatarURL(assetPrefix, grant.UserResourceOwner, grant.AvatarURL),
		PreferredLoginName: grant.PreferredLoginName,
		UserType:           TypeToPb(grant.UserType),
		GrantedOrgId:       grant.GrantedOrgID,
		GrantedOrgName:     grant.GrantedOrgName,
		GrantedOrgDomain:   grant.GrantedOrgDomain,
		Details: object.ToViewDetailsPb(
			grant.Sequence,
			grant.CreationDate,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}
}

func UserGrantQueriesToQuery(ctx context.Context, queries []*user_pb.UserGrantQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = UserGrantQueryToQuery(ctx, query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func UserGrantQueryToQuery(ctx context.Context, query *user_pb.UserGrantQuery) (query.SearchQuery, error) {
	switch q := query.Query.(type) {
	case *user_pb.UserGrantQuery_DisplayNameQuery:
		return UserGrantDisplayNameQueryToModel(q.DisplayNameQuery)
	case *user_pb.UserGrantQuery_EmailQuery:
		return UserGrantEmailQueryToModel(q.EmailQuery)
	case *user_pb.UserGrantQuery_FirstNameQuery:
		return UserGrantFirstNameQueryToModel(q.FirstNameQuery)
	case *user_pb.UserGrantQuery_LastNameQuery:
		return UserGrantLastNameQueryToModel(q.LastNameQuery)
	case *user_pb.UserGrantQuery_OrgDomainQuery:
		return UserGrantOrgDomainQueryToModel(q.OrgDomainQuery)
	case *user_pb.UserGrantQuery_OrgNameQuery:
		return UserGrantOrgNameQueryToModel(q.OrgNameQuery)
	case *user_pb.UserGrantQuery_ProjectGrantIdQuery:
		return UserGrantProjectGrantIDQueryToModel(q.ProjectGrantIdQuery)
	case *user_pb.UserGrantQuery_ProjectIdQuery:
		return UserGrantProjectIDQueryToModel(q.ProjectIdQuery)
	case *user_pb.UserGrantQuery_ProjectNameQuery:
		return UserGrantProjectNameQueryToModel(q.ProjectNameQuery)
	case *user_pb.UserGrantQuery_RoleKeyQuery:
		return UserGrantRoleKeyQueryToModel(q.RoleKeyQuery)
	case *user_pb.UserGrantQuery_UserIdQuery:
		return UserGrantUserIDQueryToModel(q.UserIdQuery)
	case *user_pb.UserGrantQuery_UserNameQuery:
		return UserGrantUserNameQueryToModel(q.UserNameQuery)
	case *user_pb.UserGrantQuery_WithGrantedQuery:
		return UserGrantWithGrantedQueryToModel(ctx, q.WithGrantedQuery)
	case *user_pb.UserGrantQuery_UserTypeQuery:
		return UserGrantUserTypeQueryToModel(q.UserTypeQuery)
	default:
		return nil, errors.New("invalid query")
	}
}

func UserGrantDisplayNameQueryToModel(q *user_pb.UserGrantDisplayNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantDisplayNameQuery(q.DisplayName, object.TextMethodToQuery(q.Method))
}

func UserGrantEmailQueryToModel(q *user_pb.UserGrantEmailQuery) (query.SearchQuery, error) {
	return query.NewUserGrantEmailQuery(q.Email, object.TextMethodToQuery(q.Method))
}

func UserGrantFirstNameQueryToModel(q *user_pb.UserGrantFirstNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantFirstNameQuery(q.FirstName, object.TextMethodToQuery(q.Method))
}

func UserGrantLastNameQueryToModel(q *user_pb.UserGrantLastNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantLastNameQuery(q.LastName, object.TextMethodToQuery(q.Method))
}

func UserGrantOrgDomainQueryToModel(q *user_pb.UserGrantOrgDomainQuery) (query.SearchQuery, error) {
	return query.NewUserGrantDomainQuery(q.OrgDomain, object.TextMethodToQuery(q.Method))
}

func UserGrantOrgNameQueryToModel(q *user_pb.UserGrantOrgNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantOrgNameQuery(q.OrgName, object.TextMethodToQuery(q.Method))
}

func UserGrantProjectIDQueryToModel(q *user_pb.UserGrantProjectIDQuery) (query.SearchQuery, error) {
	return query.NewUserGrantProjectIDSearchQuery(q.ProjectId)
}

func UserGrantProjectGrantIDQueryToModel(q *user_pb.UserGrantProjectGrantIDQuery) (query.SearchQuery, error) {
	return query.NewUserGrantGrantIDSearchQuery(q.ProjectGrantId)
}

func UserGrantProjectNameQueryToModel(q *user_pb.UserGrantProjectNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantProjectNameQuery(q.ProjectName, object.TextMethodToQuery(q.Method))
}

func UserGrantRoleKeyQueryToModel(q *user_pb.UserGrantRoleKeyQuery) (query.SearchQuery, error) {
	return query.NewUserGrantRoleQuery(q.RoleKey)
}

func UserGrantUserIDQueryToModel(q *user_pb.UserGrantUserIDQuery) (query.SearchQuery, error) {
	return query.NewUserGrantUserIDSearchQuery(q.UserId)
}

func UserGrantUserNameQueryToModel(q *user_pb.UserGrantUserNameQuery) (query.SearchQuery, error) {
	return query.NewUserGrantUsernameQuery(q.UserName, object.TextMethodToQuery(q.Method))
}

func UserGrantWithGrantedQueryToModel(ctx context.Context, q *user_pb.UserGrantWithGrantedQuery) (query.SearchQuery, error) {
	return query.NewUserGrantWithGrantedQuery(authz.GetCtxData(ctx).OrgID)
}

func UserGrantUserTypeQueryToModel(q *user_pb.UserGrantUserTypeQuery) (query.SearchQuery, error) {
	return query.NewUserGrantUserTypeQuery(grantTypeToDomain(q.Type))
}

func grantTypeToDomain(typ user_pb.Type) domain.UserType {
	switch typ {
	case user_pb.Type_TYPE_HUMAN:
		return domain.UserTypeHuman
	case user_pb.Type_TYPE_MACHINE:
		return domain.UserTypeMachine
	default:
		return domain.UserTypeUnspecified
	}
}
