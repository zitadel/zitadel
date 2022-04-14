package command

import (
	"context"

	"github.com/caos/zitadel/internal/command/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func AddPasswordComplexityPolicy(
	a *instance.Aggregate,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewPasswordComplexityPolicyAddedEvent(ctx, &a.Aggregate,
					minLength,
					hasLowercase,
					hasUppercase,
					hasNumber,
					hasSymbol,
				),
			}, nil
		}, nil
	}
}
