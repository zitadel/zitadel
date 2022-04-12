package command

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func (c *commandNew) AddSecretGeneratorConfig(ctx context.Context, typ domain.SecretGeneratorType, config *crypto.GeneratorConfig) (*domain.ObjectDetails, error) {
	agg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.es.Filter, addSecretGeneratorConfig(agg, typ, config))
	if err != nil {
		return nil, err
	}

	events, err := c.es.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: agg.ResourceOwner,
	}, nil
}

//AddOrg defines the commands to create a new org,
// this includes the verified default domain
func addSecretGeneratorConfig(a *instance.Aggregate, typ domain.SecretGeneratorType, config *crypto.GeneratorConfig) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if !typ.Valid() {
			return nil, errors.ThrowInvalidArgument(nil, "V2-FGqVj", "Errors.InvalidArgument")
		}
		if config.Length < 1 {
			return nil, errors.ThrowInvalidArgument(nil, "V2-jEqCt", "Errors.InvalidArgument")
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
				return nil, errors.ThrowAlreadyExists(nil, "V2-6CqKo", "Errors.SecretGenerator.AlreadyExists")
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
