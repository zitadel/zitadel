package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
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
