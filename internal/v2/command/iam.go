package command

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func (r *CommandSide) GetIAM(ctx context.Context, aggregateID string) (*iam_model.IAM, error) {
	iamWriteModel := NewIAMWriteModel(aggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	return writeModelToIAM(iamWriteModel), nil
}
