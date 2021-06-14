package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step17 struct {
	DefaultLoginTexts []domain.CustomLoginText
}

func (s *Step17) Step() domain.Step {
	return domain.Step17
}

func (s *Step17) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep17(ctx, s)
}

func (c *Commands) SetupStep17(ctx context.Context, step *Step17) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		events := make([]eventstore.EventPusher, 0)

		for _, text := range step.DefaultLoginTexts {
			mailEvents, _, err := c.setDefaultLoginText(ctx, iamAgg, &text)
			if err != nil {
				return nil, err
			}
			events = append(events, mailEvents...)
		}

		logging.Log("SETUP-m9Wrf").Info("default login text set up")
		return events, nil
	}
	return c.setup(ctx, step, fn)
}
