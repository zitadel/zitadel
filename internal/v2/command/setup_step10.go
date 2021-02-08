package command

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"

	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step10 struct {
	DefaultMailTemplate iam_model.MailTemplate
	DefaultMailTexts    []iam_model.MailText
}

func (s *Step10) Step() domain.Step {
	return domain.Step10
}

func (s *Step10) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep10(ctx, s)
}

func (r *CommandSide) SetupStep10(ctx context.Context, step *Step10) error {
	fn := func(iam *IAMWriteModel) (*iam_repo.Aggregate, error) {
		//multiFactorModel := NewIAMMultiFactorWriteModel()
		//iamAgg := IAMAggregateFromWriteModel(&multiFactorModel.MultiFactoryWriteModel.WriteModel)
		//if !step.Passwordless {
		//	return iamAgg, nil
		//}
		//err := setPasswordlessAllowedInPolicy(ctx, r, iamAgg)
		//if err != nil {
		//	return nil, err
		//}
		//logging.Log("SETUP-AEG2t").Info("allowed passwordless in login policy")
		//err = r.addMultiFactorToDefaultLoginPolicy(ctx, iamAgg, multiFactorModel, domain.MultiFactorTypeU2FWithPIN)
		//if err != nil {
		//	return nil, err
		//}
		//logging.Log("SETUP-ADfng").Info("added passwordless to MFA login policy")
		//return iamAgg, err
		return nil, nil
	}
	return r.setup(ctx, step, fn)
}
