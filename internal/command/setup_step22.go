package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Step22 struct {
	DefaultMailTemplate domain.MailTemplate
}

func (s *Step22) Step() domain.Step {
	return domain.Step22
}

func (s *Step22) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep22(ctx, s)
}

func (c *Commands) SetupStep22(ctx context.Context, step *Step22) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		_, mailTemplateEvent, err := c.changeDefaultMailTemplate(ctx, &step.DefaultMailTemplate)
		if err != nil {
			return nil, err
		}
		logging.Log("SETUP-Dvdwq").Info("default mail template/text set up")
		return []eventstore.Command{mailTemplateEvent}, nil
	}
	return c.setup(ctx, step, fn)
}
