package management

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func iamFromModel(iam *iam_model.Iam) *management.Iam {
	return &management.Iam{
		IamProjectId: iam.IamProjectID,
		GlobalOrgId:  iam.GlobalOrgID,
		SetUpDone:    iam.SetUpDone,
		SetUpStarted: iam.SetUpStarted,
	}
}
