package command

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) AddSecretGeneratorConfig(ctx context.Context, generatorType string, config *crypto.GeneratorConfig) (*domain.ObjectDetails, error) {
	if generatorType == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-0pkwf", "Errors.SecretGenerator.TypeMissing")
	}

	generatorWriteModel, err := c.getSecretConfig(ctx, generatorType)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&generatorWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewSecretGeneratorAddedEvent(
		ctx,
		iamAgg,
		generatorType,
		config.Length,
		config.Expiry.Duration,
		config.IncludeLowerLetters,
		config.IncludeUpperLetters,
		config.IncludeDigits,
		config.IncludeSymbols))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(generatorWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&generatorWriteModel.WriteModel), nil
}

func (c *Commands) ChangeSecretGeneratorConfig(ctx context.Context, generatorType string, config *crypto.GeneratorConfig) (*domain.ObjectDetails, error) {
	if generatorType == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-33k9f", "Errors.SecretGenerator.TypeMissing")
	}

	generatorWriteModel, err := c.getSecretConfig(ctx, generatorType)
	if err != nil {
		return nil, err
	}
	if generatorWriteModel.State == domain.SecretGeneratorStateUnspecified || generatorWriteModel.State == domain.SecretGeneratorStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-3n9ls", "Errors.SecretGenerator.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&generatorWriteModel.WriteModel)

	changedEvent, hasChanged, err := generatorWriteModel.NewChangedEvent(
		ctx,
		iamAgg,
		generatorType,
		config.Length,
		config.Expiry.Duration,
		config.IncludeLowerLetters,
		config.IncludeUpperLetters,
		config.IncludeDigits,
		config.IncludeSymbols)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-m0o3f", "Errors.NoChangesFound")
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(generatorWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&generatorWriteModel.WriteModel), nil
}

func (c *Commands) RemoveSecretGeneratorConfig(ctx context.Context, generatorType string) (*domain.ObjectDetails, error) {
	if generatorType == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-2j9lw", "Errors.SecretGenerator.TypeMissing")
	}

	generatorWriteModel, err := c.getSecretConfig(ctx, generatorType)
	if err != nil {
		return nil, err
	}
	if generatorWriteModel.State == domain.SecretGeneratorStateUnspecified || generatorWriteModel.State == domain.SecretGeneratorStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-b8les", "Errors.SecretGenerator.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&generatorWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewSecretGeneratorRemovedEvent(ctx, iamAgg, generatorType))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(generatorWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&generatorWriteModel.WriteModel), nil
}

func (c *Commands) getSecretConfig(ctx context.Context, generatorType string) (_ *IAMSecretGeneratorConfigWriteModel, err error) {
	writeModel := NewIAMSecretGeneratorConfigWriteModel(generatorType)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
