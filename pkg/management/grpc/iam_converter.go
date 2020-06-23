package grpc

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func iamFromModel(iam *iam_model.Iam) *Iam {
	return &Iam{
		IamProjectId: iam.IamProjectID,
		GlobalOrgId:  iam.GlobalOrgID,
		SetUpDone:    iam.SetUpDone,
		SetUpStarted: iam.SetUpStarted,
	}
}
