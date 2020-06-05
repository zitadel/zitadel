package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_view_model "github.com/caos/zitadel/internal/org/repository/view/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type UserGrantRepo struct {
	SearchLimit uint64
	View        *view.View
	IamID       string
	Auth        auth.Config
	AuthZRepo   *authz_repo.EsRepository
}

func (repo *UserGrantRepo) SearchMyUserGrants(ctx context.Context, request *grant_model.UserGrantSearchRequest) (*grant_model.UserGrantSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.Queries = append(request.Queries, &grant_model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_USER_ID, Method: global_model.SEARCHMETHOD_EQUALS, Value: auth.GetCtxData(ctx).UserID})
	grants, count, err := repo.View.SearchUserGrants(request)
	if err != nil {
		return nil, err
	}
	return &grant_model.UserGrantSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.UserGrantsToModel(grants),
	}, nil
}

func (repo *UserGrantRepo) SearchMyProjectOrgs(ctx context.Context, request *grant_model.UserGrantSearchRequest) (*grant_model.ProjectOrgSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	ctxData := auth.GetCtxData(ctx)
	if ctxData.ProjectID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "APP-7lqva", "Could not get ProjectID")
	}
	if ctxData.ProjectID == repo.AuthZRepo.IamProjectID {
		isAdmin, err := repo.IsIamAdmin(ctx)
		if err != nil {
			return nil, err
		}
		if isAdmin {
			return repo.SearchAdminOrgs(request)
		}
	}
	request.Queries = append(request.Queries, &grant_model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_PROJECT_ID, Method: global_model.SEARCHMETHOD_EQUALS, Value: ctxData.ProjectID})

	grants, err := repo.SearchMyUserGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	return grantRespToOrgResp(grants), nil
}

func (repo *UserGrantRepo) SearchMyZitadelPermissions(ctx context.Context) ([]string, error) {
	grant, err := repo.AuthZRepo.ResolveGrants(ctx)
	if err != nil {
		return nil, err
	}
	permissions := &grant_model.Permissions{Permissions: []string{}}
	for _, role := range grant.Roles {
		roleName, ctxID := auth.SplitPermission(role)
		for _, mapping := range repo.Auth.RolePermissionMappings {
			if mapping.Role == roleName {
				permissions.AppendPermissions(ctxID, mapping.Permissions...)
			}
		}
	}
	return permissions.Permissions, nil
}

func (repo *UserGrantRepo) SearchAdminOrgs(request *grant_model.UserGrantSearchRequest) (*grant_model.ProjectOrgSearchResponse, error) {
	searchRequest := &org_model.OrgSearchRequest{}
	if len(request.Queries) > 0 {
		for _, q := range request.Queries {
			if q.Key == grant_model.USERGRANTSEARCHKEY_ORG_NAME {
				searchRequest.Queries = append(searchRequest.Queries, &org_model.OrgSearchQuery{Key: org_model.ORGSEARCHKEY_ORG_NAME, Method: q.Method, Value: q.Value})
			}
		}
	}
	orgs, count, err := repo.View.SearchOrgs(searchRequest)
	if err != nil {
		return nil, err
	}
	return orgRespToOrgResp(orgs, count), nil
}

func (repo *UserGrantRepo) IsIamAdmin(ctx context.Context) (bool, error) {
	grantSearch := &grant_model.UserGrantSearchRequest{
		Queries: []*grant_model.UserGrantSearchQuery{
			&grant_model.UserGrantSearchQuery{Key: grant_model.USERGRANTSEARCHKEY_RESOURCEOWNER, Method: global_model.SEARCHMETHOD_EQUALS, Value: repo.IamID},
		}}
	result, err := repo.SearchMyUserGrants(ctx, grantSearch)
	if err != nil {
		return false, err
	}
	if result.TotalResult == 0 {
		return false, nil
	}
	return true, nil
}

func grantRespToOrgResp(grants *grant_model.UserGrantSearchResponse) *grant_model.ProjectOrgSearchResponse {
	resp := &grant_model.ProjectOrgSearchResponse{
		TotalResult: grants.TotalResult,
	}
	resp.Result = make([]*grant_model.Org, len(grants.Result))
	for i, g := range grants.Result {
		resp.Result[i] = &grant_model.Org{OrgID: g.ResourceOwner, OrgName: g.OrgName}
	}
	return resp
}

func orgRespToOrgResp(orgs []*org_view_model.OrgView, count int) *grant_model.ProjectOrgSearchResponse {
	resp := &grant_model.ProjectOrgSearchResponse{
		TotalResult: uint64(count),
	}
	resp.Result = make([]*grant_model.Org, len(orgs))
	for i, o := range orgs {
		resp.Result[i] = &grant_model.Org{OrgID: o.ID, OrgName: o.Name}
	}
	return resp
}

func mergeOrgAndAdminGrant(ctxData auth.CtxData, orgGrant, iamAdminGrant *model.UserGrantView) (grant *auth.Grant) {
	if orgGrant != nil {
		roles := orgGrant.RoleKeys
		if iamAdminGrant != nil {
			roles = addIamAdminRoles(roles, iamAdminGrant.RoleKeys)
		}
		grant = &auth.Grant{OrgID: orgGrant.ResourceOwner, Roles: roles}
	} else if iamAdminGrant != nil {
		grant = &auth.Grant{
			OrgID: ctxData.OrgID,
			Roles: iamAdminGrant.RoleKeys,
		}
	}
	return grant
}

func addIamAdminRoles(orgRoles, iamAdminRoles []string) []string {
	result := make([]string, 0)
	result = append(result, iamAdminRoles...)
	for _, role := range orgRoles {
		if !auth.ExistsPerm(result, role) {
			result = append(result, role)
		}
	}
	return result
}
