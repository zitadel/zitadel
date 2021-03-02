package user

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	usr_grant_model "github.com/caos/zitadel/internal/usergrant/model"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
)

func UserGrantsToPb(grants []*usr_grant_model.UserGrantView) []*user_pb.UserGrant {
	u := make([]*user_pb.UserGrant, len(grants))
	for i, grant := range grants {
		u[i] = UserGrantToPb(grant)
	}
	return u
}

func UserGrantToPb(grant *usr_grant_model.UserGrantView) *user_pb.UserGrant {
	return &user_pb.UserGrant{
		GrantId:     grant.ID,
		UserId:      grant.UserID,
		State:       ModelUserGrantStateToPb(grant.State),
		RoleKeys:    grant.RoleKeys,
		UserName:    grant.UserName,
		FirstName:   grant.FirstName,
		LastName:    grant.LastName,
		Email:       grant.Email,
		DisplayName: grant.DisplayName,
		OrgId:       grant.ResourceOwner,
		OrgDomain:   grant.OrgPrimaryDomain,
		OrgName:     grant.OrgName,
		ProjectId:   grant.ProjectID,
		ProjectName: grant.ProjectName,
		Details: object.ToDetailsPb(
			grant.Sequence,
			grant.ChangeDate,
			grant.ResourceOwner,
		),
	}
}

func UserGrantQueriesToModel(queries []*user_pb.UserGrantQuery) []*usr_grant_model.UserGrantSearchQuery {
	q := make([]*usr_grant_model.UserGrantSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = UserGrantQueryToModel(query)
	}
	return q
}

func UserGrantQueryToModel(query *user_pb.UserGrantQuery) *usr_grant_model.UserGrantSearchQuery {
	switch q := query.Query.(type) {
	case *user_pb.UserGrantQuery_DisplayName:
		return UserGrantDisplayNameQueryToModel(q.DisplayName)
	case *user_pb.UserGrantQuery_Email:
		return UserGrantEmailQueryToModel(q.Email)
	case *user_pb.UserGrantQuery_FirstName:
		return UserGrantFirstNameQueryToModel(q.FirstName)
	case *user_pb.UserGrantQuery_LastName:
		return UserGrantLastNameQueryToModel(q.LastName)
	case *user_pb.UserGrantQuery_OrgDomain:
		return UserGrantOrgDomainQueryToModel(q.OrgDomain)
	case *user_pb.UserGrantQuery_OrgName:
		return UserGrantOrgNameQueryToModel(q.OrgName)
	case *user_pb.UserGrantQuery_ProjectGrantId:
		return UserGrantProjectGrantIDQueryToModel(q.ProjectGrantId)
	case *user_pb.UserGrantQuery_ProjectId:
		return UserGrantProjectIDQueryToModel(q.ProjectId)
	case *user_pb.UserGrantQuery_ProjectName:
		return UserGrantProjectNameQueryToModel(q.ProjectName)
	case *user_pb.UserGrantQuery_RoleKey:
		return UserGrantRoleKeyQueryToModel(q.RoleKey)
	case *user_pb.UserGrantQuery_UserId:
		return UserGrantUserIDQueryToModel(q.UserId)
	case *user_pb.UserGrantQuery_UserName:
		return UserGrantUserNameQueryToModel(q.UserName)
	case *user_pb.UserGrantQuery_WithGranted:
		return UserGrantWithGrantedQueryToModel(q.WithGranted)
	default:
		return nil
	}
}

func UserGrantDisplayNameQueryToModel(q *user_pb.UserGrantDisplayNameQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyDisplayName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.DisplayName,
	}
}

func UserGrantEmailQueryToModel(q *user_pb.UserGrantEmailQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyEmail,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.Email,
	}
}

func UserGrantFirstNameQueryToModel(q *user_pb.UserGrantFirstNameQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyFirstName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.FirstName,
	}
}

func UserGrantLastNameQueryToModel(q *user_pb.UserGrantLastNameQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyLastName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.LastName,
	}
}

func UserGrantOrgDomainQueryToModel(q *user_pb.UserGrantOrgDomainQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyOrgDomain,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.OrgDomain,
	}
}

func UserGrantOrgNameQueryToModel(q *user_pb.UserGrantOrgNameQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyOrgName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.OrgName,
	}
}

func UserGrantProjectIDQueryToModel(q *user_pb.UserGrantProjectIDQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyProjectID,
		Method: domain.SearchMethodEquals,
		Value:  q.ProjectId,
	}
}

func UserGrantProjectGrantIDQueryToModel(q *user_pb.UserGrantProjectGrantIDQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyGrantID,
		Method: domain.SearchMethodEquals,
		Value:  q.ProjectGrantId,
	}
}

func UserGrantProjectNameQueryToModel(q *user_pb.UserGrantProjectNameQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyProjectName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.ProjectName,
	}
}

func UserGrantRoleKeyQueryToModel(q *user_pb.UserGrantRoleKeyQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyRoleKey,
		Method: domain.SearchMethodListContains,
		Value:  q.RoleKey,
	}
}

func UserGrantUserIDQueryToModel(q *user_pb.UserGrantUserIDQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  q.UserId,
	}
}

func UserGrantUserNameQueryToModel(q *user_pb.UserGrantUserNameQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyUserName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.UserName,
	}
}

func UserGrantWithGrantedQueryToModel(q *user_pb.UserGrantWithGrantedQuery) *usr_grant_model.UserGrantSearchQuery {
	return &usr_grant_model.UserGrantSearchQuery{
		Key:    usr_grant_model.UserGrantSearchKeyWithGranted,
		Method: domain.SearchMethodEquals,
		Value:  q.WithGranted,
	}
}
