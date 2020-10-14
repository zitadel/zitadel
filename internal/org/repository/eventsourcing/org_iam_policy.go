package eventsourcing

import (
	"context"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func OrgIAMPolicyAddedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, policy *iam_es_model.OrgIAMPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-i9sJS", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgIAMPolicyAdded, policy)
	}
}

func OrgIAMPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org, policy *iam_es_model.OrgIAMPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9Ksie", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		changes := org.OrgIAMPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Js6Vs", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.OrgIAMPolicyChanged, changes)
	}
}

func OrgIamPolicyRemovedAggregate(aggCreator *es_models.AggregateCreator, org *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := OrgAggregate(ctx, aggCreator, org.AggregateID, org.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.OrgIAMPolicyRemoved, nil)
	}
}
