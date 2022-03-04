package iam

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func AddLoginPolicy(
	a *iam.Aggregate,
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
			return []eventstore.Command{
				iam.NewLoginPolicyAddedEvent(ctx, &a.Aggregate,
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
