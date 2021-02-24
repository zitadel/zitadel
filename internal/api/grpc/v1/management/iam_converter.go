package management

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func iamFromModel(iam *iam_model.IAM) *management.Iam {
	return &management.Iam{
		IamProjectId: iam.IAMProjectID,
		GlobalOrgId:  iam.GlobalOrgID,
		SetUpDone:    iamSetupStepFromModel(iam.SetUpDone),
		SetUpStarted: iamSetupStepFromModel(iam.SetUpStarted),
	}
}

func iamSetupStepFromModel(step domain.Step) management.IamSetupStep {
	switch step {
	case domain.Step1:
		return management.IamSetupStep_iam_setup_step_1
	case domain.Step2:
		return management.IamSetupStep_iam_setup_step_2
	// case iam_model.Step3:
	// 	return management.IamSetupStep_iam_setup_step_3
	// case iam_model.Step4:
	// 	return management.IamSetupStep_iam_setup_step_4
	// case iam_model.Step5:
	// 	return management.IamSetupStep_iam_setup_step_5
	// case iam_model.Step6:
	// 	return management.IamSetupStep_iam_setup_step_6

	default:
		return management.IamSetupStep_iam_setup_step_UNDEFINED
	}
}
