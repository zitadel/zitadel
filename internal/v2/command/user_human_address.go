package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
)

func (r *CommandSide) ChangeHumanAddress(ctx context.Context, address *domain.Address) (*domain.Address, error) {
	existingAddress, err := r.addressWriteModel(ctx, address.AggregateID, address.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if existingAddress.State == domain.AddressStateUnspecified || existingAddress.State == domain.AddressStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-0pLdo", "Errors.User.Address.NotFound")
	}
	changedEvent, hasChanged := existingAddress.NewChangedEvent(ctx, address.Country, address.Locality, address.PostalCode, address.Region, address.StreetAddress)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0cs", "Errors.User.Address.NotChanged")
	}
	userAgg := UserAggregateFromWriteModel(&existingAddress.WriteModel)
	userAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingAddress, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToAddress(existingAddress), nil
}

func (r *CommandSide) addressWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *HumanAddressWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanAddressWriteModel(userID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
