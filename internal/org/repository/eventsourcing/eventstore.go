package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgEventstore struct {
	eventstore.Eventstore
}

type OrgConfig struct {
	eventstore.Eventstore
}

func StartOrg(conf OrgConfig) (*OrgEventstore, error) {
	return &OrgEventstore{Eventstore: conf.Eventstore}, nil
}

func (es *OrgEventstore) OrgByID(ctx context.Context, org *org_model.Org) (*org_model.Org, error) {
	query, err := OrgByIDQuery(org.ID, org.Sequence)
	if err != nil {
		return nil, err
	}

	esOrg := OrgFromModel(org)
	err = sdk.Filter(ctx, es.FilterEvents, esOrg.AppendEvents, query)
	if err != nil {
		return nil, err
	}

	return OrgToModel(esOrg), nil
}
