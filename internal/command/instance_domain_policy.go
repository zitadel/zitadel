package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func AddDefaultDomainPolicy(
	a *instance.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomains bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewDomainPolicyAddedEvent(ctx, &a.Aggregate,
					userLoginMustBeDomain,
					validateOrgDomains,
				),
			}, nil
		}, nil
	}
}
