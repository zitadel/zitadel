package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
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

func (r *CommandSide) setGlobalOrg(ctx context.Context, iamAgg *iam.Aggregate, iamWriteModel *IAMWriteModel, orgID string) error {
	err := r.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return err
	}
	if iamWriteModel.GlobalOrgID != "" {
		return caos_errs.ThrowPreconditionFailed(nil, "IAM-HGG24", "Errors.IAM.GlobalOrgAlreadySet")
	}
	iamAgg.PushEvents(iam.NewGlobalOrgSetEventEvent(ctx, orgID))
	return nil
}

func (r *CommandSide) setIAMProject(ctx context.Context, iamAgg *iam.Aggregate, iamWriteModel *IAMWriteModel, projectID string) error {
	err := r.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return err
	}
	if iamWriteModel.ProjectID != "" {
		return caos_errs.ThrowPreconditionFailed(nil, "IAM-EGbw2", "Errors.IAM.IAMProjectAlreadySet")
	}
	iamAgg.PushEvents(iam.NewIAMProjectSetEvent(ctx, projectID))
	return nil
}
