package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
)

func readModelToIAM(readModel *ReadModel) *model.IAM {
	return &model.IAM{
		ObjectRoot:   readModelToObjectRoot(readModel.ReadModel),
		GlobalOrgID:  readModel.GlobalOrgID,
		IAMProjectID: readModel.ProjectID,
		SetUpDone:    readModel.SetUpDone,
		SetUpStarted: readModel.SetUpStarted,
	}
}

func readModelToObjectRoot(readModel eventstore.ReadModel) models.ObjectRoot {
	return models.ObjectRoot{
		AggregateID:   readModel.AggregateID,
		ChangeDate:    readModel.ChangeDate,
		CreationDate:  readModel.CreationDate,
		ResourceOwner: readModel.ResourceOwner,
		Sequence:      readModel.ProcessedSequence,
	}
}
