package grpc

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/management/grpc"
)

func iamFromModel(iam *iam_model.Iam) *grpc.Iam {
	return &grpc.Iam{
		IamProjectId: iam.IamProjectID,
		GlobalOrgId:  iam.GlobalOrgID,
		SetUpDone:    iam.SetUpDone,
		SetUpStarted: iam.SetUpStarted,
	}
}
