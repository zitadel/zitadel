package eventstore

import (
	"context"

	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
)

type OrgMemberRepository struct {
	*org_es.OrgEventstore
}

func (repo *OrgMemberRepository) OrgMemberByID(ctx context.Context, orgID, userID string) (member *org_model.OrgMember, err error) {
	member = org_model.NewOrgMember(orgID, userID)
	return repo.OrgEventstore.OrgMemberByIDs(ctx, member)
}

func (repo *OrgMemberRepository) AddOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	return repo.OrgEventstore.AddOrgMember(ctx, member)
}

func (repo *OrgMemberRepository) ChangeOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	return repo.OrgEventstore.ChangeOrgMember(ctx, member)
}

func (repo *OrgMemberRepository) RemoveOrgMember(ctx context.Context, orgID, userID string) error {
	member := org_model.NewOrgMember(orgID, userID)
	return repo.OrgEventstore.RemoveOrgMember(ctx, member)
}
