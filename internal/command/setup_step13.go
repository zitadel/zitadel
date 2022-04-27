package command

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
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
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		_, mailTemplateEvent, err := c.changeDefaultMailTemplate(ctx, &step.DefaultMailTemplate)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-4insR").Info("default mail template/text set up")
		return []eventstore.Command{mailTemplateEvent}, nil
	}
	return c.setup(ctx, step, fn)
}
