package eventsourcing

import (
	"context"

	admin_model "github.com/caos/zitadel/internal/admin/model"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
)

type OrgRepo struct {
	*org_es.OrgEventstore
}

func (s *OrgRepo) GetOrgByID(ctx context.Context, orgID string) (_ *org_model.Org, err error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mvn3R", "Not implemented")
}

// func (s *OrgRepo) SearchOrgs(ctx context.Context, request *OrgSearchRequest) (_ *OrgSearchResponse, err error) {
// 	return nil, errors.ThrowUnimplemented(nil, "GRPC-Po9Hd", "Not implemented")
// }

// func (s *OrgRepo) IsOrgUnique(ctx context.Context, request *UniqueOrgRequest) (org *UniqueOrgResponse, err error) {
// 	return nil, errors.ThrowUnimplemented(nil, "GRPC-0p6Fw", "Not implemented")
// }

func (repo *OrgRepo) SetUpOrg(ctx context.Context, setUp *admin_model.SetupOrg) (*admin_model.SetupOrg, error) {
	eventstoreOrg := eventsourcing.OrgFromModel(setUp.Org)
	aggregates, err := eventsourcing.OrgCreateAggregates(ctx, repo.AggregateCreator(), eventstoreOrg)
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
