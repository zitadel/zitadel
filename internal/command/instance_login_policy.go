package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func AddDefaultLoginPolicy(
	a *instance.Aggregate,
	allowUsernamePassword bool,
	allowRegister bool,
	allowExternalIDP bool,
	forceMFA bool,
	hidePasswordReset bool,
	passwordlessType domain.PasswordlessType,
	passwordCheckLifetime time.Duration,
	externalLoginCheckLifetime time.Duration,
	mfaInitSkipLifetime time.Duration,
	secondFactorCheckLifetime time.Duration,
	multiFactorCheckLifetime time.Duration,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewLoginPolicyAddedEvent(ctx, &a.Aggregate,
					allowUsernamePassword,
					allowRegister,
					allowExternalIDP,
					forceMFA,
					hidePasswordReset,
					passwordlessType,
					passwordCheckLifetime,
					externalLoginCheckLifetime,
					mfaInitSkipLifetime,
					secondFactorCheckLifetime,
					multiFactorCheckLifetime,
				),
			}, nil
		}, nil
	}
}

func AddSecondFactorToDefaultLoginPolicy(a *instance.Aggregate, factor domain.SecondFactorType) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewLoginPolicySecondFactorAddedEvent(ctx, &a.Aggregate, factor),
			}, nil
		}, nil
	}
}

func AddMultiFactorToDefaultLoginPolicy(a *instance.Aggregate, factor domain.MultiFactorType) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewLoginPolicyMultiFactorAddedEvent(ctx, &a.Aggregate, factor),
			}, nil
		}, nil
	}
}
