package eventsourcing

import (
	"context"
	"strconv"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/sony/sonyflake"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

func OrgByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dke74", "id should be filled")
	}
	return OrgQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func OrgQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(org_model.OrgAggregate).
		LatestSequenceFilter(latestSequence)
}

func OrgAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, id string, sequence uint64) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, id, org_model.OrgAggregate, orgVersion, sequence)
}

func OrgCreateAggregates(ctx context.Context, aggCreator *es_models.AggregateCreator, org *Org) (_ []*es_models.Aggregate, err error) {
	if org == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "project should not be nil")
	}

	domainAgrregate, err := uniqueDomainAggregate(ctx, aggCreator, org.Domain)
	if err != nil {
		return nil, err
	}

	nameAggregate, err := uniqueNameAggregate(ctx, aggCreator, org.Name)
	if err != nil {
		return nil, err
	}

	id, err := idGenerator.NextID()
	if err != nil {
		return nil, err
	}
	org.ID = strconv.FormatUint(id, 10)

	agg, err := OrgAggregate(ctx, aggCreator, org.ID, org.Sequence)
	if err != nil {
		return nil, err
	}
	agg.AppendEvent(org_model.OrgAdded, org)

	return []*es_models.Aggregate{
		domainAgrregate,
		nameAggregate,
		agg,
	}, nil

}

func OrgUpdateAggregates(ctx context.Context, aggCreator *es_models.AggregateCreator, existing *Org, updated *Org) ([]*es_models.Aggregate, error) {
	if existing == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
	}
	if updated == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
	}
	changes := existing.Changes(updated)

	aggregates := make([]*es_models.Aggregate, 0, 3)

	if name, ok := changes["name"]; ok {
		nameAggregate, err := uniqueNameAggregate(ctx, aggCreator, name.(string))
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, nameAggregate)
	}

	if name, ok := changes["domain"]; ok {
		domainAggregate, err := uniqueDomainAggregate(ctx, aggCreator, name.(string))
		if err != nil {
			return nil, err
		}
		aggregates = append(aggregates, domainAggregate)
	}

	orgAggregate, err := OrgAggregate(ctx, aggCreator, existing.ID, existing.Sequence)
	if err != nil {
		return nil, err
	}
	orgAggregate.AppendEvent(org_model.OrgChanged, changes)
	aggregates = append(aggregates, orgAggregate)

	return aggregates, nil
}

func OrgDeactivateAggregate(aggCreator *es_models.AggregateCreator, project *Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return projectStateAggregate(aggCreator, project, org_model.OrgDeactivated)
}

func OrgReactivateAggregate(aggCreator *es_models.AggregateCreator, project *Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return projectStateAggregate(aggCreator, project, org_model.OrgReactivated)
}

func projectStateAggregate(aggCreator *es_models.AggregateCreator, org *Org, state es_models.EventType) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if org == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-37dur", "existing project should not be nil")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.ID, org.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(state, nil)
	}
}

func uniqueDomainAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, domain string) (*es_models.Aggregate, error) {
	aggregate, err := aggCreator.NewAggregate(ctx, domain, org_model.OrgDomainAggregate, orgVersion, 0)
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(org_model.OrgDomainReserved, nil)
	if err != nil {
		return nil, err
	}

	query := es_models.NewSearchQuery().
		AggregateIDFilter(domain).
		AggregateTypeFilter(org_model.OrgDomainAggregate).
		OrderDesc().
		SetLimit(1)

	validation := func(events ...*es_models.Event) error {
		if len(events) == 0 {
			aggregate.PreviousSequence = 0
			return nil
		}
		if events[0].Type == org_model.OrgDomainReserved {
			return errors.ThrowPreconditionFailed(nil, "EVENT-WMKO4", "domain already reseved")
		}
		aggregate.PreviousSequence = events[0].Sequence
		return nil
	}

	aggregate.SetPrecondition(query, validation)

	return aggregate, nil
}

func uniqueNameAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, name string) (*es_models.Aggregate, error) {
	aggregate, err := aggCreator.NewAggregate(ctx, name, org_model.OrgNameAggregate, orgVersion, 0)
	if err != nil {
		return nil, err
	}
	aggregate, err = aggregate.AppendEvent(org_model.OrgNameReserved, nil)
	if err != nil {
		return nil, err
	}

	query := es_models.NewSearchQuery().
		AggregateIDFilter(name).
		AggregateTypeFilter(org_model.OrgNameAggregate).
		OrderDesc().
		SetLimit(1)

	validation := func(events ...*es_models.Event) error {
		if len(events) == 0 {
			aggregate.PreviousSequence = 0
			return nil
		}
		if events[0].Type == org_model.OrgNameReserved {
			return errors.ThrowPreconditionFailed(nil, "EVENT-WMKO4", "domain already reseved")
		}
		aggregate.PreviousSequence = events[0].Sequence
		return nil
	}

	aggregate.SetPrecondition(query, validation)

	return aggregate, nil
}
