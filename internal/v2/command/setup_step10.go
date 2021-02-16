package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step10 struct {
	DefaultMailTemplate domain.MailTemplate
	DefaultMailTexts    []domain.MailText
}

func (s *Step10) Step() domain.Step {
	return domain.Step10
}

func (s *Step10) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep10(ctx, s)
}

func (r *CommandSide) SetupStep10(ctx context.Context, step *Step10) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		iamAgg := IAMAggregateFromWriteModel(&iam.WriteModel)
		err := r.addDefaultMailTemplate(ctx, iamAgg, NewIAMMailTemplateWriteModel(), &step.DefaultMailTemplate)
		if err != nil {
			return nil, err
		}
		for _, text := range step.DefaultMailTexts {
			r.addDefaultMailText(ctx, iamAgg, NewIAMMailTextWriteModel(text.MailTextType, text.Language), &text)
			if err != nil {
				return nil, err
			}
		}
		logging.Log("SETUP-3N9fs").Info("default mail template/text set up")
		return iamAgg, nil
	}
	return r.setup(ctx, step, fn)
}
