package eventstore

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/errors"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/view/model"
)

type OrgRepository struct {
	SearchLimit uint64
	*org_es.OrgEventstore
	View  *mgmt_view.View
	Roles []string
}

func (repo *OrgRepository) OrgByID(ctx context.Context, id string) (*org_model.OrgView, error) {
	org, err := repo.View.OrgByID(id)
	if err != nil {
		return nil, err
	}
	return model.OrgToModel(org), nil
}

func (repo *OrgRepository) OrgByDomainGlobal(ctx context.Context, domain string) (*org_model.OrgView, error) {
	verifiedDomain, err := repo.View.VerifiedOrgDomain(domain)
	if err != nil {
		return nil, err
	}
	return repo.OrgByID(ctx, verifiedDomain.OrgID)
}

func (repo *OrgRepository) UpdateOrg(ctx context.Context, org *org_model.Org) (*org_model.Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "EVENT-RkurR", "not implemented")
}

func (repo *OrgRepository) DeactivateOrg(ctx context.Context, id string) (*org_model.Org, error) {
	return repo.OrgEventstore.DeactivateOrg(ctx, id)
}

func (repo *OrgRepository) ReactivateOrg(ctx context.Context, id string) (*org_model.Org, error) {
	return repo.OrgEventstore.ReactivateOrg(ctx, id)
}

func (repo *OrgRepository) GetMyOrgIamPolicy(ctx context.Context) (*org_model.OrgIamPolicy, error) {
	return repo.OrgEventstore.GetOrgIamPolicy(ctx, authz.GetCtxData(ctx).OrgID)
}

func (repo *OrgRepository) SearchMyOrgDomains(ctx context.Context, request *org_model.OrgDomainSearchRequest) (*org_model.OrgDomainSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.Queries = append(request.Queries, &org_model.OrgDomainSearchQuery{Key: org_model.OrgDomainSearchKeyOrgID, Method: global_model.SearchMethodEquals, Value: authz.GetCtxData(ctx).OrgID})
	domains, count, err := repo.View.SearchOrgDomains(request)
	if err != nil {
		return nil, err
	}
	return &org_model.OrgDomainSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.OrgDomainsToModel(domains),
	}, nil
}

func (repo *OrgRepository) AddMyOrgDomain(ctx context.Context, domain *org_model.OrgDomain) (*org_model.OrgDomain, error) {
	domain.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddOrgDomain(ctx, domain)
}

func (repo *OrgRepository) RemoveMyOrgDomain(ctx context.Context, domain string) error {
	d := org_model.NewOrgDomain(authz.GetCtxData(ctx).OrgID, domain)
	return repo.OrgEventstore.RemoveOrgDomain(ctx, d)
}

func (repo *OrgRepository) OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64) (*org_model.OrgChanges, error) {
	changes, err := repo.OrgEventstore.OrgChanges(ctx, id, lastSequence, limit)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func (repo *OrgRepository) OrgMemberByID(ctx context.Context, orgID, userID string) (*org_model.OrgMemberView, error) {
	member, err := repo.View.OrgMemberByIDs(orgID, userID)
	if err != nil {
		return nil, err
	}
	return model.OrgMemberToModel(member), nil
}

func (repo *OrgRepository) AddMyOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	member.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddOrgMember(ctx, member)
}

func (repo *OrgRepository) ChangeMyOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	member.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangeOrgMember(ctx, member)
}

func (repo *OrgRepository) RemoveMyOrgMember(ctx context.Context, userID string) error {
	member := org_model.NewOrgMember(authz.GetCtxData(ctx).OrgID, userID)
	return repo.OrgEventstore.RemoveOrgMember(ctx, member)
}

func (repo *OrgRepository) SearchMyOrgMembers(ctx context.Context, request *org_model.OrgMemberSearchRequest) (*org_model.OrgMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.Queries[len(request.Queries)-1] = &org_model.OrgMemberSearchQuery{Key: org_model.OrgMemberSearchKeyOrgID, Method: global_model.SearchMethodEquals, Value: authz.GetCtxData(ctx).OrgID}
	members, count, err := repo.View.SearchOrgMembers(request)
	if err != nil {
		return nil, err
	}
	return &org_model.OrgMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.OrgMembersToModel(members),
	}, nil
}

func (repo *OrgRepository) GetOrgMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "ORG") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}
