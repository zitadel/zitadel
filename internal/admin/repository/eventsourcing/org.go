package eventsourcing

import (
	"context"
	"strings"

	admin_model "github.com/caos/zitadel/internal/admin/model"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
)

type OrgRepo struct {
	*org_es.OrgEventstore
}

func (repo *OrgRepo) SetUpOrg(ctx context.Context, setUp *admin_model.SetupOrg) (*admin_model.SetupOrg, error) {
	eventstoreOrg := eventsourcing.OrgFromModel(setUp.Org)
	aggregates, err := eventsourcing.OrgCreatedAggregates(ctx, repo.AggregateCreator(), eventstoreOrg)
	if err != nil {
		return nil, err
	}
	//TODO: create user with org as resource owner
	//TODO: add member(user) to org

	err = sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, eventstoreOrg.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	setUp.Org = eventsourcing.OrgToModel(eventstoreOrg)
	return setUp, nil
}

func (repo *OrgRepo) OrgByID(ctx context.Context, id string) (*org_model.Org, error) {
	return repo.OrgEventstore.OrgByID(ctx, org_model.NewOrg(id))
}

func (repo *OrgRepo) SearchOrgs(ctx context.Context) ([]*org_model.Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "EVENT-hFIHK", "search not implemented")
}

func (repo *OrgRepo) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	var found bool
	err = sdk.Filter(ctx, repo.FilterEvents, isUniqueValidation(&found), eventsourcing.OrgNameUniqueQuery(name))
	if err != nil && !errors.IsNotFound(err) {
		return false, err
	}

	err = sdk.Filter(ctx, repo.FilterEvents, isUniqueValidation(&found), eventsourcing.OrgDomainUniqueQuery(domain))
	if err != nil && !errors.IsNotFound(err) {
		return false, err
	}

	return !found, nil
}

func isUniqueValidation(unique *bool) func(events ...*models.Event) error {
	return func(events ...*models.Event) error {
		if len(events) == 0 {
			return nil
		}
		*unique = *unique || strings.HasSuffix(string(events[0].Type), "reserved")

		return nil
	}
}
