package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func ServiceAccountCreateAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, account *model.Machine, resourceOwner string) (_ []*es_models.Aggregate, err error) {
	accountAggregate, err := UserAggregate(ctx, aggCreator, &account.ObjectRoot)
	if err != nil {
		return nil, err
	}
	accountAggregate, err = accountAggregate.AppendEvent(model.MachineAdded, account)
	if err != nil {
		return nil, err
	}

	uniqueName, err := reservedUniqueUserNameAggregate(ctx, aggCreator, resourceOwner, account.Name, true)
	if err != nil {
		return nil, err
	}

	return []*es_models.Aggregate{
		accountAggregate,
		uniqueName,
	}, nil
}

func ServiceAccountChangeAggregate(aggCreator *es_models.AggregateCreator, existingAccount *model.Machine, updatedAccount *model.Machine) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if updatedAccount == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "Errors.Internal")
		}
		agg, err := UserAggregate(ctx, aggCreator, &existingAccount.ObjectRoot)
		if err != nil {
			return nil, err
		}
		changes := existingAccount.Changes(updatedAccount)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0spow", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.MachineChanged, changes)
	}
}

func ServiceAccountDeactivateAggregate(aggCreator *es_models.AggregateCreator, aggregate *es_models.ObjectRoot) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, aggregate, model.ServiceAccountDeactivated)
}

func ServiceAccountReactivateAggregate(aggCreator *es_models.AggregateCreator, aggregate *es_models.ObjectRoot) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, aggregate, model.ServiceAccountReactivated)
}

func ServiceAccountLockAggregate(aggCreator *es_models.AggregateCreator, aggregate *es_models.ObjectRoot) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, aggregate, model.ServiceAccountLocked)
}

func ServiceAccountUnlockAggregate(aggCreator *es_models.AggregateCreator, aggregate *es_models.ObjectRoot) func(ctx context.Context) (*es_models.Aggregate, error) {
	return userStateAggregate(aggCreator, aggregate, model.ServiceAccountUnlocked)
}
