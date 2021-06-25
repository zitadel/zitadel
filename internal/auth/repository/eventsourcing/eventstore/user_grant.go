package eventstore

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	authz_repo "github.com/caos/zitadel/internal/authz/repository/eventsourcing"
	caos_errs "github.com/caos/zitadel/internal/errors"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_view_model "github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type UserGrantRepo struct {
	SearchLimit     uint64
	View            *view.View
	IamID           string
	Auth            authz.Config
	AuthZRepo       *authz_repo.EsRepository
	PrefixAvatarURL string
}

func (repo *UserGrantRepo) SearchMyUserGrants(ctx context.Context, request *grant_model.UserGrantSearchRequest) (*grant_model.UserGrantSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, err := repo.View.GetLatestUserGrantSequence()
	logging.Log("EVENT-Hd7s3").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest user grant sequence")
	request.Queries = append(request.Queries, &grant_model.UserGrantSearchQuery{Key: grant_model.UserGrantSearchKeyUserID, Method: domain.SearchMethodEquals, Value: authz.GetCtxData(ctx).UserID})
	grants, count, err := repo.View.SearchUserGrants(request)
	if err != nil {
		return nil, err
	}
	result := &grant_model.UserGrantSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.UserGrantsToModel(grants, repo.PrefixAvatarURL),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *UserGrantRepo) SearchMyProjectOrgs(ctx context.Context, request *grant_model.UserGrantSearchRequest) (*grant_model.ProjectOrgSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	ctxData := authz.GetCtxData(ctx)
	if ctxData.ProjectID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "APP-7lqva", "Could not get ProjectID")
	}
	err = repo.AuthZRepo.FillIamProjectID(ctx)
	if err != nil {
		return nil, err
	}
	if ctxData.ProjectID == repo.AuthZRepo.UserGrantRepo.IamProjectID {
		isAdmin, err := repo.IsIamAdmin(ctx)
		if err != nil {
			return nil, err
		}
		if isAdmin {
			return repo.SearchAdminOrgs(request)
		}
		return repo.searchZitadelOrgs(ctxData, request)
	}
	request.Queries = append(request.Queries, &grant_model.UserGrantSearchQuery{Key: grant_model.UserGrantSearchKeyProjectID, Method: domain.SearchMethodEquals, Value: ctxData.ProjectID})

	grants, err := repo.SearchMyUserGrants(ctx, request)
	if err != nil {
		return nil, err
	}
	if len(grants.Result) > 0 {
		return grantRespToOrgResp(grants), nil
	}
	return repo.userOrg(ctxData)
}

func membershipsToOrgResp(memberships []*user_view_model.UserMembershipView, count uint64) *grant_model.ProjectOrgSearchResponse {
	orgs := make([]*grant_model.Org, 0, len(memberships))
	for _, m := range memberships {
		if !containsOrg(orgs, m.ResourceOwner) {
			orgs = append(orgs, &grant_model.Org{OrgID: m.ResourceOwner, OrgName: m.ResourceOwnerName})
		}
	}
	return &grant_model.ProjectOrgSearchResponse{
		TotalResult: count,
		Result:      orgs,
	}
}

func (repo *UserGrantRepo) SearchMyUserMemberships(ctx context.Context, request *user_model.UserMembershipSearchRequest) (*user_model.UserMembershipSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, sequenceErr := repo.View.GetLatestUserMembershipSequence()
	logging.Log("EVENT-Dn7sf").OnError(sequenceErr).Warn("could not read latest user sequence")

	memberships, count, err := repo.View.SearchUserMemberships(request)
	if err != nil {
		return nil, err
	}
	result := &user_model.UserMembershipSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      user_view_model.UserMembershipsToModel(memberships),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *UserGrantRepo) SearchMyZitadelPermissions(ctx context.Context) ([]string, error) {
	memberships, err := repo.searchUserMemberships(ctx)
	if err != nil {
		return nil, err
	}
	permissions := &grant_model.Permissions{Permissions: []string{}}
	for _, membership := range memberships {
		for _, role := range membership.Roles {
			permissions = repo.mapRoleToPermission(permissions, membership, role)
		}
	}
	return permissions.Permissions, nil
}

func (repo *UserGrantRepo) searchUserMemberships(ctx context.Context) ([]*user_view_model.UserMembershipView, error) {
	ctxData := authz.GetCtxData(ctx)
	orgMemberships, orgCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: domain.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyResourceOwner,
				Method: domain.SearchMethodEquals,
				Value:  ctxData.OrgID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	iamMemberships, iamCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: domain.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyAggregateID,
				Method: domain.SearchMethodEquals,
				Value:  repo.IamID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if orgCount == 0 && iamCount == 0 {
		return []*user_view_model.UserMembershipView{}, nil
	}
	return append(orgMemberships, iamMemberships...), nil
}

func (repo *UserGrantRepo) SearchMyProjectPermissions(ctx context.Context) ([]string, error) {
	ctxData := authz.GetCtxData(ctx)
	usergrant, err := repo.View.UserGrantByIDs(ctxData.OrgID, ctxData.ProjectID, ctxData.UserID)
	if err != nil {
		return nil, err
	}
	permissions := make([]string, len(usergrant.RoleKeys))
	for i, role := range usergrant.RoleKeys {
		permissions[i] = role
	}
	return permissions, nil
}

func (repo *UserGrantRepo) SearchAdminOrgs(request *grant_model.UserGrantSearchRequest) (*grant_model.ProjectOrgSearchResponse, error) {
	searchRequest := &org_model.OrgSearchRequest{
		SortingColumn: org_model.OrgSearchKeyOrgNameIgnoreCase,
		Asc:           true,
	}
	if len(request.Queries) > 0 {
		for _, q := range request.Queries {
			if q.Key == grant_model.UserGrantSearchKeyOrgName {
				searchRequest.Queries = append(searchRequest.Queries, &org_model.OrgSearchQuery{Key: org_model.OrgSearchKeyOrgName, Method: q.Method, Value: q.Value})
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
			{Key: grant_model.UserGrantSearchKeyResourceOwner, Method: domain.SearchMethodEquals, Value: repo.IamID},
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

func (repo *UserGrantRepo) UserGrantsByProjectAndUserID(projectID, userID string) ([]*grant_model.UserGrantView, error) {
	grants, err := repo.View.UserGrantsByProjectAndUserID(projectID, userID)
	if err != nil {
		return nil, err
	}
	return model.UserGrantsToModel(grants, repo.PrefixAvatarURL), nil
}

func (repo *UserGrantRepo) userOrg(ctxData authz.CtxData) (*grant_model.ProjectOrgSearchResponse, error) {
	user, err := repo.View.UserByID(ctxData.UserID)
	if err != nil {
		return nil, err
	}
	org, err := repo.View.OrgByID(user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	return &grant_model.ProjectOrgSearchResponse{Result: []*grant_model.Org{&grant_model.Org{
		OrgID:   org.ID,
		OrgName: org.Name,
	}}}, nil
}

func (repo *UserGrantRepo) searchZitadelOrgs(ctxData authz.CtxData, request *grant_model.UserGrantSearchRequest) (*grant_model.ProjectOrgSearchResponse, error) {
	memberships, count, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Offset: request.Offset,
		Limit:  request.Limit,
		Asc:    request.Asc,
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: domain.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(memberships) > 0 {
		return membershipsToOrgResp(memberships, count), nil
	}
	return repo.userOrg(ctxData)
}

func (repo *UserGrantRepo) mapRoleToPermission(permissions *grant_model.Permissions, membership *user_view_model.UserMembershipView, role string) *grant_model.Permissions {
	for _, mapping := range repo.Auth.RolePermissionMappings {
		if mapping.Role == role {
			ctxID := ""
			if membership.MemberType == int32(user_model.MemberTypeProject) || membership.MemberType == int32(user_model.MemberTypeProjectGrant) {
				ctxID = membership.ObjectID
			}
			permissions.AppendPermissions(ctxID, mapping.Permissions...)
		}
	}
	return permissions
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

func orgRespToOrgResp(orgs []*org_view_model.OrgView, count uint64) *grant_model.ProjectOrgSearchResponse {
	resp := &grant_model.ProjectOrgSearchResponse{
		TotalResult: count,
	}
	resp.Result = make([]*grant_model.Org, len(orgs))
	for i, o := range orgs {
		resp.Result[i] = &grant_model.Org{OrgID: o.ID, OrgName: o.Name}
	}
	return resp
}

func containsOrg(orgs []*grant_model.Org, resourceOwner string) bool {
	for _, org := range orgs {
		if org.OrgID == resourceOwner {
			return true
		}
	}
	return false
}

func userMembershipToMembership(membership *user_view_model.UserMembershipView) *authz.Membership {
	return &authz.Membership{
		MemberType:  authz.MemberType(membership.MemberType),
		AggregateID: membership.AggregateID,
		ObjectID:    membership.ObjectID,
		Roles:       membership.Roles,
	}
}

func userMembershipsToMemberships(memberships []*user_view_model.UserMembershipView) []*authz.Membership {
	result := make([]*authz.Membership, len(memberships))
	for i, m := range memberships {
		result[i] = userMembershipToMembership(m)
	}
	return result
}
