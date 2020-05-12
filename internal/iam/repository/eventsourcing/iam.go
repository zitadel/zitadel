package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func IamByIDQuery(id string, latestSequence uint64) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0soe4", "id should be filled")
	}
	return IamQuery(latestSequence).
		AggregateIDFilter(id), nil
}

func IamQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IamAggregate).
		LatestSequenceFilter(latestSequence)
}

func IamAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.Iam) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lo04e", "existing iam should not be nil")
	}
	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IamAggregate, model.IamVersion, iam.Sequence)
}

func IamAggregateOverwriteContext(ctx context.Context, aggCreator *es_models.AggregateCreator, iam *model.Iam, resourceOwnerID string, userID string) (*es_models.Aggregate, error) {
	if iam == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dis83", "existing iam should not be nil")
	}

	return aggCreator.NewAggregate(ctx, iam.AggregateID, model.IamAggregate, model.IamVersion, iam.Sequence, es_models.OverwriteResourceOwner(resourceOwnerID), es_models.OverwriteEditorUser(userID))
}

func IamSetupStartedAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamSetupStarted, nil)
	}
}

func IamSetupDoneAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.IamSetupDone, nil)
	}
}

func IamSetGlobalOrgAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam, globalOrg string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if globalOrg == "" {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8siwa", "globalOrg must be set")
		}
		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.GlobalOrgSet, &model.Iam{GlobalOrgID: globalOrg})
	}
}

func IamSetIamProjectAggregate(aggCreator *es_models.AggregateCreator, iam *model.Iam, projectID string) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if projectID == "" {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-sjuw3", "projectID must be set")
		}
		agg, err := IamAggregate(ctx, aggCreator, iam)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamProjectSet, &model.Iam{IamProjectID: projectID})
	}
}

func IamMemberAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, member *model.IamMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9sope", "member should not be nil")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamMemberAdded, member)
	}
}

func IamMemberChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, member *model.IamMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-38skf", "member should not be nil")
		}

		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamMemberChanged, member)
	}
}

func IamMemberRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Iam, member *model.IamMember) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if member == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-90lsw", "member should not be nil")
		}
		agg, err := IamAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.IamMemberRemoved, member)
	}
}
