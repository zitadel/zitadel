package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
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

func (r *CommandSide) setGlobalOrg(ctx context.Context, iamAgg *eventstore.Aggregate, iamWriteModel *IAMWriteModel, orgID string) (eventstore.EventPusher, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	if iamWriteModel.GlobalOrgID != "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-HGG24", "Errors.IAM.GlobalOrgAlreadySet")
	}
	return iam.NewGlobalOrgSetEventEvent(ctx, iamAgg, orgID), nil
}

func (r *CommandSide) setIAMProject(ctx context.Context, iamAgg *eventstore.Aggregate, iamWriteModel *IAMWriteModel, projectID string) (eventstore.EventPusher, error) {
	err := r.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	if iamWriteModel.ProjectID != "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-EGbw2", "Errors.IAM.IAMProjectAlreadySet")
	}
	return iam.NewIAMProjectSetEvent(ctx, iamAgg, projectID), nil
}
