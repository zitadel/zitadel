package eventstore

import (
	"context"
	"strings"

	chg_model "github.com/caos/zitadel/internal/changes/model"
	chg_event "github.com/caos/zitadel/internal/changes/repository/eventsourcing"
	"github.com/caos/zitadel/internal/errors"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/view"
)

type OrgRepository struct {
	SearchLimit uint64
	*org_es.OrgEventstore
	View          *mgmt_view.View
	Roles         []string
	ChangesEvents *chg_event.ChangesEventstore
}

func (repo *OrgRepository) OrgByID(ctx context.Context, id string) (*org_model.Org, error) {
	org := org_model.NewOrg(id)
	return repo.OrgEventstore.OrgByID(ctx, org)
}

func (repo *OrgRepository) OrgByDomainGlobal(ctx context.Context, domain string) (*org_model.OrgView, error) {
	org, err := repo.View.OrgByDomain(domain)
	if err != nil {
		return nil, err
	}
	return view.OrgToModel(org), nil
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

func (repo *OrgRepository) OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64) (*chg_model.Changes, error) {
	changes, err := repo.ChangesEvents.Changes(ctx, chg_model.Org, id, 0, 0)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func (repo *OrgRepository) OrgMemberByID(ctx context.Context, orgID, userID string) (member *org_model.OrgMember, err error) {
	member = org_model.NewOrgMember(orgID, userID)
	return repo.OrgEventstore.OrgMemberByIDs(ctx, member)
}

func (repo *OrgRepository) AddOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	return repo.OrgEventstore.AddOrgMember(ctx, member)
}

func (repo *OrgRepository) ChangeOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	return repo.OrgEventstore.ChangeOrgMember(ctx, member)
}

func (repo *OrgRepository) RemoveOrgMember(ctx context.Context, orgID, userID string) error {
	member := org_model.NewOrgMember(orgID, userID)
	return repo.OrgEventstore.RemoveOrgMember(ctx, member)
}

func (repo *OrgRepository) SearchOrgMembers(ctx context.Context, request *org_model.OrgMemberSearchRequest) (*org_model.OrgMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	members, count, err := repo.View.SearchOrgMembers(request)
	if err != nil {
		return nil, err
	}
	return &org_model.OrgMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      view.OrgMembersToModel(members),
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
