package command

import (
	"context"

	"github.com/caos/zitadel/internal/v2/domain"
)

//TODO: private
func (r *CommandSide) GetIAM(ctx context.Context) (*domain.IAM, error) {
	iamWriteModel := NewIAMWriteModel()
	err := r.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	return writeModelToIAM(iamWriteModel), nil
}
