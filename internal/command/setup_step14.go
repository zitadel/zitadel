package command

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	org_repo "github.com/caos/zitadel/internal/repository/org"
)

type Step14 struct {
	ActivateExistingLabelPolicies bool
}

func (s *Step14) Step() domain.Step {
	return domain.Step14
}

func (s *Step14) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep14(ctx, s)
}

func (c *Commands) SetupStep14(ctx context.Context, step *Step14) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		var events []eventstore.EventPusher
		if step.ActivateExistingLabelPolicies {
			existingPolicies := NewExistingLabelPoliciesReadModel(ctx)
			err := c.eventstore.FilterToQueryReducer(ctx, existingPolicies)
			if err != nil {
				return nil, err
			}
			for _, aggID := range existingPolicies.aggregateIDs {
				if iamAgg.ID == aggID {
					events = append(events, iam_repo.NewLabelPolicyActivatedEvent(ctx, iamAgg))
					continue
				}
				events = append(events, org_repo.NewLabelPolicyActivatedEvent(ctx, &org_repo.NewAggregate(aggID, aggID).Aggregate))
			}
		}
		logging.Log("SETUP-M9fsd").Info("activate login policies")
		return events, nil
	}
	return c.setup(ctx, step, fn)
}
