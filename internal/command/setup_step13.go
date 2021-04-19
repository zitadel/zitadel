package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step13 struct {
	DefaultMailTemplate domain.MailTemplate
}

func (s *Step13) Step() domain.Step {
	return domain.Step13
}

func (s *Step13) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep13(ctx, s)
}

func (c *Commands) SetupStep13(ctx context.Context, step *Step13) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		_, mailTemplateEvent, err := c.changeDefaultMailTemplate(ctx, &step.DefaultMailTemplate)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-4insR").Info("default mail template/text set up")
		return []eventstore.EventPusher{mailTemplateEvent}, nil
	}
	return c.setup(ctx, step, fn)
}
