package management

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
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

func iamSetupStepFromModel(step iam_model.Step) management.IamSetupStep {
	switch step {
	case iam_model.Step1:
		return management.IamSetupStep_iam_setup_step_1
	case iam_model.Step2:
		return management.IamSetupStep_iam_setup_step_2
	default:
		return management.IamSetupStep_iam_setup_step_UNDEFINED
	}
}
