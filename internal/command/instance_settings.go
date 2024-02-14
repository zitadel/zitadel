package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddSecretGeneratorConfig(ctx context.Context, typ domain.SecretGeneratorType, config *crypto.GeneratorConfig) (*domain.ObjectDetails, error) {
	agg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddSecretGeneratorConfig(agg, typ, config))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: agg.ResourceOwner,
	}, nil
}

func prepareAddSecretGeneratorConfig(a *instance.Aggregate, typ domain.SecretGeneratorType, config *crypto.GeneratorConfig) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if !typ.Valid() {
			return nil, zerrors.ThrowInvalidArgument(nil, "V2-FGqVj", "Errors.InvalidArgument")
		}
		if config.Length < 1 {
			return nil, zerrors.ThrowInvalidArgument(nil, "V2-jEqCt", "Errors.InvalidArgument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceSecretGeneratorConfigWriteModel(ctx, typ)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}

			if writeModel.State == domain.SecretGeneratorStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "V2-6CqKo", "Errors.SecretGenerator.AlreadyExists")
			}

			return []eventstore.Command{
				instance.NewSecretGeneratorAddedEvent(
					ctx,
					&a.Aggregate,
					typ,
					config.Length,
					config.Expiry,
					config.IncludeLowerLetters,
					config.IncludeUpperLetters,
					config.IncludeDigits,
					config.IncludeSymbols,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) ChangeSecretGeneratorConfig(ctx context.Context, generatorType domain.SecretGeneratorType, config *crypto.GeneratorConfig) (*domain.ObjectDetails, error) {
	if generatorType == domain.SecretGeneratorTypeUnspecified {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-33k9f", "Errors.SecretGenerator.TypeMissing")
	}

	generatorWriteModel, err := c.getSecretConfig(ctx, generatorType)
	if err != nil {
		return nil, err
	}
	instanceAgg := InstanceAggregateFromWriteModel(&generatorWriteModel.WriteModel)
	if generatorWriteModel.State == domain.SecretGeneratorStateUnspecified || generatorWriteModel.State == domain.SecretGeneratorStateRemoved {
		err = c.pushAppendAndReduce(ctx, generatorWriteModel,
			instance.NewSecretGeneratorAddedEvent(
				ctx,
				instanceAgg,
				generatorType,
				config.Length,
				config.Expiry,
				config.IncludeLowerLetters,
				config.IncludeUpperLetters,
				config.IncludeDigits,
				config.IncludeSymbols,
			),
		)
		if err != nil {
			return nil, err
		}
		return writeModelToObjectDetails(&generatorWriteModel.WriteModel), nil
	}

	changedEvent, hasChanged, err := generatorWriteModel.NewChangedEvent(
		ctx,
		instanceAgg,
		generatorType,
		config.Length,
		config.Expiry,
		config.IncludeLowerLetters,
		config.IncludeUpperLetters,
		config.IncludeDigits,
		config.IncludeSymbols)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-m0o3f", "Errors.NoChangesFound")
	}
	if err = c.pushAppendAndReduce(ctx, generatorWriteModel, changedEvent); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&generatorWriteModel.WriteModel), nil
}

func (c *Commands) RemoveSecretGeneratorConfig(ctx context.Context, generatorType domain.SecretGeneratorType) (*domain.ObjectDetails, error) {
	if generatorType == domain.SecretGeneratorTypeUnspecified {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-2j9lw", "Errors.SecretGenerator.TypeMissing")
	}

	generatorWriteModel, err := c.getSecretConfig(ctx, generatorType)
	if err != nil {
		return nil, err
	}
	if generatorWriteModel.State == domain.SecretGeneratorStateUnspecified || generatorWriteModel.State == domain.SecretGeneratorStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-b8les", "Errors.SecretGenerator.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&generatorWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewSecretGeneratorRemovedEvent(ctx, instanceAgg, generatorType))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(generatorWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&generatorWriteModel.WriteModel), nil
}

func (c *Commands) getSecretConfig(ctx context.Context, generatorType domain.SecretGeneratorType) (_ *InstanceSecretGeneratorConfigWriteModel, err error) {
	writeModel := NewInstanceSecretGeneratorConfigWriteModel(ctx, generatorType)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
