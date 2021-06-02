package command

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
)

type Step10 struct {
	DefaultMailTemplate domain.MailTemplate
	DefaultMailTexts    []domain.CustomMessageText
}

func (s *Step10) Step() domain.Step {
	return domain.Step10
}

func (s *Step10) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep10(ctx, s)
}

func (c *Commands) SetupStep10(ctx context.Context, step *Step10) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.EventPusher, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		mailTemplateEvent, err := c.addDefaultMailTemplate(ctx, iamAgg, NewIAMMailTemplateWriteModel(), &step.DefaultMailTemplate)
		if err != nil {
			return nil, err
		}
		events := []eventstore.EventPusher{
			mailTemplateEvent,
		}
		//for _, text := range step.DefaultMailTexts {
		//	defaultTextEvent, err := c.addDefaultMailText(ctx, iamAgg, NewIAMMailTextWriteModel(text.Template, text.Language), &text)
		//	if err != nil {
		//		return nil, err
		//	}
		//	events = append(events, defaultTextEvent)
		//}
		logging.Log("SETUP-3N9fs").Info("default mail template/text set up")
		return events, nil
	}
	return c.setup(ctx, step, fn)
}
