package command

import (
	"context"

	"github.com/caos/zitadel/internal/v2/domain"
)

func (r *CommandSide) GetIAM(ctx context.Context, aggregateID string) (*domain.IAM, error) {
	iamWriteModel := NewIAMWriteModel(aggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	return writeModelToIAM(iamWriteModel), nil
}
