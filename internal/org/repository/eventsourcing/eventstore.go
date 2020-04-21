package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgEventstore struct {
	eventstore.Eventstore
}

type OrgConfig struct {
	eventstore.Eventstore
}

func StartOrg(conf OrgConfig) *OrgEventstore {
	return &OrgEventstore{Eventstore: conf.Eventstore}
}

func (es *OrgEventstore) OrgByID(ctx context.Context, org *org_model.Org) (*org_model.Org, error) {
	query, err := OrgByIDQuery(org.ID, org.Sequence)
	if err != nil {
		return nil, err
	}

	esOrg := OrgFromModel(org)
	err = es_sdk.Filter(ctx, es.FilterEvents, esOrg.AppendEvents, query)
	if err != nil {
		return nil, err
	}

	return OrgToModel(esOrg), nil
}

func (es *OrgEventstore) DeactivateOrg(ctx context.Context, orgModel *org_model.Org) (*org_model.Org, error) {
	org := OrgFromModel(orgModel)
	aggregate := OrgDeactivateAggregate(es.AggregateCreator(), org)
	err := es_sdk.Push(ctx, es.PushAggregates, org.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	return OrgToModel(org), nil
}

func (es *OrgEventstore) ReactivateOrg(ctx context.Context, orgModel *org_model.Org) (*org_model.Org, error) {
	org := OrgFromModel(orgModel)
	aggregate := OrgReactivateAggregate(es.AggregateCreator(), org)
	err := es_sdk.Push(ctx, es.PushAggregates, org.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}
	return OrgToModel(org), nil
}

func (es *OrgEventstore) OrgMemberByIDs(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	if member.UserID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ld93d", "userID missing")
	}

	org, err := es.OrgByID(ctx, org_model.NewOrg(member.ID))
	if err != nil {
		return nil, err
	}

	for _, m := range org.Members {
		if m.UserID == member.UserID {
			return m, nil
		}
	}

	return nil, errors.ThrowInternal(nil, "EVENT-a0Poo", "Could not find member in list")
}

func (es *OrgEventstore) AddOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	if !member.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}

	existing, err := es.OrgByID(ctx, org_model.NewOrg(member.ID))
	if err != nil {
		return nil, err
	}

	if existing.ContainsMember(member.UserID) {
		return nil, errors.ThrowAlreadyExists(nil, "EVENT-idke6", "User is already member of this Org")
	}

	repoOrg := OrgFromModel(existing)
	repoMember := OrgMemberFromModel(member)

	addAggregate := OrgMemberAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	// es.orgCache.cacheOrg(repoOrg)
	for _, m := range repoOrg.Members {
		if m.UserID == member.UserID {
			return OrgMemberToModel(m), nil
		}
	}
	return nil, errors.ThrowInternal(nil, "EVENT-mIyhl", "Could not find member in list")
}

func (es *OrgEventstore) ChangeOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	if !member.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(member.ID))
	if err != nil {
		return nil, err
	}
	if !existing.ContainsMember(member.UserID) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-oe39f", "User is not member of this org")
	}
	repoOrg := OrgFromModel(existing)
	repoMember := OrgMemberFromModel(member)

	orgAggregate := OrgMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
	// es.orgCache.cacheOrg(repoOrg)
	for _, m := range repoOrg.Members {
		if m.UserID == member.UserID {
			return OrgMemberToModel(m), nil
		}
	}
	return nil, errors.ThrowInternal(nil, "EVENT-3udjs", "Could not find member in list")
}

func (es *OrgEventstore) RemoveOrgMember(ctx context.Context, member *org_model.OrgMember) error {
	if member.UserID == "" {
		return errors.ThrowPreconditionFailed(nil, "EVENT-d43fs", "UserID and Roles are required")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(member.ID))
	if err != nil {
		return err
	}
	if !existing.ContainsMember(member.UserID) {
		return errors.ThrowPreconditionFailed(nil, "EVENT-swf34", "User is not member of this org")
	}
	repoOrg := OrgFromModel(existing)
	repoMember := OrgMemberFromModel(member)

	orgAggregate := OrgMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
	// es.orgCache.cacheOrg(repoOrg)
	return err
}
