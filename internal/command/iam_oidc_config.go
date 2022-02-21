package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) AddOIDCConfig(ctx context.Context, config *domain.OIDCConfig) (*domain.ObjectDetails, error) {
	secretConfigWriteModel, err := c.getOIDCConfig(ctx)
	if err != nil {
		return nil, err
	}
	if secretConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-d9nlw", "Errors.OIDCConfig.AlreadyExists")
	}
	iamAgg := IAMAggregateFromWriteModel(&secretConfigWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewOIDCConfigAddedEvent(
		ctx,
		iamAgg,
		config.AccessTokenLifetime,
		config.IdTokenLifetime,
		config.RefreshTokenIdleExpiration,
		config.RefreshTokenExpiration))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(secretConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&secretConfigWriteModel.WriteModel), nil
}

func (c *Commands) ChangeOIDCConfig(ctx context.Context, config *domain.OIDCConfig) (*domain.ObjectDetails, error) {
	secretConfigWriteModel, err := c.getOIDCConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !secretConfigWriteModel.State.Exists() {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-8snEr", "Errors.OIDCConfig.NotExisting")
	}
	iamAgg := IAMAggregateFromWriteModel(&secretConfigWriteModel.WriteModel)

	changedEvent, hasChanged, err := secretConfigWriteModel.NewChangedEvent(
		ctx,
		iamAgg,
		config.AccessTokenLifetime,
		config.IdTokenLifetime,
		config.RefreshTokenIdleExpiration,
		config.RefreshTokenExpiration)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-398uF", "Errors.NoChangesFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(secretConfigWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&secretConfigWriteModel.WriteModel), nil
}

func (c *Commands) getOIDCConfig(ctx context.Context) (_ *IAMOIDCConfigWriteModel, err error) {
	writeModel := NewIAMOIDCConfigWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
