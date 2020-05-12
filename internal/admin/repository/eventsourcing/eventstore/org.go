package eventstore

import (
	"context"

	admin_model "github.com/caos/zitadel/internal/admin/model"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type OrgRepo struct {
	Eventstore     eventstore.Eventstore
	OrgEventstore  *org_es.OrgEventstore
	UserEventstore *usr_es.UserEventstore
}

func (repo *OrgRepo) SetUpOrg(ctx context.Context, setUp *admin_model.SetupOrg) (*admin_model.SetupOrg, error) {
	org, aggregates, err := repo.OrgEventstore.PrepareCreateOrg(ctx, setUp.Org)
	if err != nil {
		return nil, err
	}

	user, userAggregate, err := repo.UserEventstore.PrepareCreateUser(ctx, setUp.User, org.AggregateID)
	if err != nil {
		return nil, err
	}

	aggregates = append(aggregates, userAggregate)
	setupModel := &Setup{Org: org, User: user}

	member := org_model.NewOrgMemberWithRoles(org.AggregateID, user.AggregateID, "ORG_ADMIN") //TODO: role as const
	_, memberAggregate, err := repo.OrgEventstore.PrepareAddOrgMember(ctx, member)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, memberAggregate)

	err = sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, setupModel.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	return SetupToModel(setupModel), nil
}

func (repo *OrgRepo) OrgByID(ctx context.Context, id string) (*org_model.Org, error) {
	return repo.OrgEventstore.OrgByID(ctx, org_model.NewOrg(id))
}

func (repo *OrgRepo) SearchOrgs(ctx context.Context) ([]*org_model.Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "EVENT-hFIHK", "search not implemented")
}

func (repo *OrgRepo) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	var found bool
	err = sdk.Filter(ctx, repo.Eventstore.FilterEvents, isUniqueValidation(&found), org_es.OrgNameUniqueQuery(name))
	if (err != nil && !errors.IsNotFound(err)) || found {
		return false, err
	}

	err = sdk.Filter(ctx, repo.Eventstore.FilterEvents, isUniqueValidation(&found), org_es.OrgDomainUniqueQuery(domain))
	if err != nil && !errors.IsNotFound(err) {
		return false, err
	}

	return !found, nil
}
