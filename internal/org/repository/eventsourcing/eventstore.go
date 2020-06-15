package eventsourcing

import (
	"context"
	"encoding/json"
	"log"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/golang/protobuf/ptypes"
)

type OrgEventstore struct {
	eventstore.Eventstore
	idGenerator id.Generator
}

type OrgConfig struct {
	eventstore.Eventstore
}

func StartOrg(conf OrgConfig) *OrgEventstore {
	return &OrgEventstore{
		Eventstore:  conf.Eventstore,
		idGenerator: id.SonyFlakeGenerator,
	}
}

func (es *OrgEventstore) PrepareCreateOrg(ctx context.Context, orgModel *org_model.Org) (*model.Org, []*es_models.Aggregate, error) {
	if orgModel == nil || !orgModel.IsValid() {
		return nil, nil, errors.ThrowInvalidArgument(nil, "EVENT-OeLSk", "org not valid")
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, nil, errors.ThrowInternal(err, "EVENT-OwciI", "id gen failed")
	}
	orgModel.AggregateID = id
	org := model.OrgFromModel(orgModel)

	aggregates, err := orgCreatedAggregates(ctx, es.AggregateCreator(), org)

	return org, aggregates, err
}

func (es *OrgEventstore) CreateOrg(ctx context.Context, orgModel *org_model.Org) (*org_model.Org, error) {
	org, aggregates, err := es.PrepareCreateOrg(ctx, orgModel)
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, org.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	return model.OrgToModel(org), nil
}

func (es *OrgEventstore) OrgByID(ctx context.Context, org *org_model.Org) (*org_model.Org, error) {
	if org == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-gQTYP", "org not set")
	}
	query, err := OrgByIDQuery(org.AggregateID, org.Sequence)
	if err != nil {
		return nil, err
	}

	esOrg := model.OrgFromModel(org)
	err = es_sdk.Filter(ctx, es.FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-kVLb2", "org not found")
	}

	return model.OrgToModel(esOrg), nil
}

func (es *OrgEventstore) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	var found bool
	err = es_sdk.Filter(ctx, es.FilterEvents, isUniqueValidation(&found), OrgNameUniqueQuery(name))
	if (err != nil && !errors.IsNotFound(err)) || found {
		return false, err
	}

	err = es_sdk.Filter(ctx, es.FilterEvents, isUniqueValidation(&found), OrgDomainUniqueQuery(domain))
	if err != nil && !errors.IsNotFound(err) {
		return false, err
	}

	return !found, nil
}

func isUniqueValidation(unique *bool) func(events ...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		if len(events) == 0 {
			return nil
		}
		*unique = *unique || events[0].Type == model.OrgDomainReserved || events[0].Type == model.OrgNameReserved

		return nil
	}
}

func (es *OrgEventstore) DeactivateOrg(ctx context.Context, orgID string) (*org_model.Org, error) {
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-oL9nT", "org not found")
	}
	org := model.OrgFromModel(existingOrg)

	aggregate := orgDeactivateAggregate(es.AggregateCreator(), org)
	err = es_sdk.Push(ctx, es.PushAggregates, org.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}

	return model.OrgToModel(org), nil
}

func (es *OrgEventstore) ReactivateOrg(ctx context.Context, orgID string) (*org_model.Org, error) {
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-oL9nT", "org not set")
	}
	org := model.OrgFromModel(existingOrg)

	aggregate := orgReactivateAggregate(es.AggregateCreator(), org)
	err = es_sdk.Push(ctx, es.PushAggregates, org.AppendEvents, aggregate)
	if err != nil {
		return nil, err

	}
	return model.OrgToModel(org), nil
}

func (es *OrgEventstore) OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64) (*org_model.OrgChanges, error) {
	query := ChangesQuery(id, lastSequence)

	events, err := es.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-328b1", "unable to get current user")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-FpQqK", "no objects found")
	}

	result := make([]*org_model.OrgChange, 0)

	for _, u := range events {
		creationDate, err := ptypes.TimestampProto(u.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &org_model.OrgChange{
			ChangeDate: creationDate,
			EventType:  u.Type.String(),
			Modifier:   u.EditorUser,
			Sequence:   u.Sequence,
		}

		orgDummy := model.Org{}
		if u.Data != nil {
			if err := json.Unmarshal(u.Data, &orgDummy); err != nil {
				log.Println("Error getting data!", err.Error())
			}
		}
		change.Data = orgDummy

		result = append(result, change)
		if lastSequence < u.Sequence {
			lastSequence = u.Sequence
		}
	}

	changes := &org_model.OrgChanges{
		Changes:      result,
		LastSequence: lastSequence,
	}

	return changes, nil
}

func ChangesQuery(orgID string, latestSequence uint64) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate).
		LatestSequenceFilter(latestSequence).
		AggregateIDFilter(orgID)
	return query
}

func (es *OrgEventstore) OrgMemberByIDs(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	if member == nil || member.UserID == "" || member.AggregateID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ld93d", "member not set")
	}

	org, err := es.OrgByID(ctx, &org_model.Org{ObjectRoot: member.ObjectRoot, Members: []*org_model.OrgMember{member}})
	if err != nil {
		return nil, err
	}

	for _, currentMember := range org.Members {
		if currentMember.UserID == member.UserID {
			return currentMember, nil
		}
	}

	return nil, errors.ThrowNotFound(nil, "EVENT-SXji6", "member not found")
}

func (es *OrgEventstore) PrepareAddOrgMember(ctx context.Context, member *org_model.OrgMember, resourceOwner string) (*model.OrgMember, *es_models.Aggregate, error) {
	if member == nil || !member.IsValid() {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}

	repoMember := model.OrgMemberFromModel(member)
	addAggregate, err := orgMemberAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoMember, resourceOwner)

	return repoMember, addAggregate, err
}

func (es *OrgEventstore) AddOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	repoMember, addAggregate, err := es.PrepareAddOrgMember(ctx, member, "")
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoMember.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}

	return model.OrgMemberToModel(repoMember), nil
}

func (es *OrgEventstore) ChangeOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	if member == nil || !member.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}

	existingMember, err := es.OrgMemberByIDs(ctx, member)
	if err != nil {
		return nil, err
	}

	member.ObjectRoot = existingMember.ObjectRoot
	repoMember := model.OrgMemberFromModel(member)
	repoExistingMember := model.OrgMemberFromModel(existingMember)

	orgAggregate := orgMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoExistingMember, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoMember.AppendEvents, orgAggregate)
	if err != nil {
		return nil, err
	}

	return model.OrgMemberToModel(repoMember), nil
}

func (es *OrgEventstore) RemoveOrgMember(ctx context.Context, member *org_model.OrgMember) error {
	if member == nil || member.UserID == "" {
		return errors.ThrowInvalidArgument(nil, "EVENT-d43fs", "UserID is required")
	}

	existingMember, err := es.OrgMemberByIDs(ctx, member)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	member.ObjectRoot = existingMember.ObjectRoot
	repoMember := model.OrgMemberFromModel(member)

	orgAggregate := orgMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoMember)
	return es_sdk.Push(ctx, es.PushAggregates, repoMember.AppendEvents, orgAggregate)
}
