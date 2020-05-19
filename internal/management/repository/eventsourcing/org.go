package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
)

type OrgRepository struct {
	*org_es.OrgEventstore
}

func (repo *OrgRepository) OrgByID(ctx context.Context, id string) (*org_model.Org, error) {
	org := org_model.NewOrg(id)
	return repo.OrgEventstore.OrgByID(ctx, org)
}

func (repo *OrgRepository) OrgByDomainGlobal(ctx context.Context, domain string) (*org_model.Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "EVENT-GQoS8", "not implemented")
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
