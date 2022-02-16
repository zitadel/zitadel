package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"golang.org/x/text/language"
)

//TODO: private as soon as setup uses query
func (c *Commands) GetIAM(ctx context.Context) (*domain.IAM, error) {
	iamWriteModel := NewIAMWriteModel()
	err := c.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	return writeModelToIAM(iamWriteModel), nil
}

func (c *Commands) setGlobalOrg(ctx context.Context, iamAgg *eventstore.Aggregate, iamWriteModel *IAMWriteModel, orgID string) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	if iamWriteModel.GlobalOrgID != "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-HGG24", "Errors.IAM.GlobalOrgAlreadySet")
	}
	return iam.NewGlobalOrgSetEventEvent(ctx, iamAgg, orgID), nil
}

func (c *Commands) setIAMProject(ctx context.Context, iamAgg *eventstore.Aggregate, iamWriteModel *IAMWriteModel, projectID string) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	if iamWriteModel.ProjectID != "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-EGbw2", "Errors.IAM.IAMProjectAlreadySet")
	}
	return iam.NewIAMProjectSetEvent(ctx, iamAgg, projectID), nil
}

func (c *Commands) SetDefaultLanguage(ctx context.Context, language language.Tag) (*domain.ObjectDetails, error) {
	iamWriteModel, err := c.getIAMWriteModel(ctx)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&iamWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewDefaultLanguageSetEvent(ctx, iamAgg, language))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(iamWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&iamWriteModel.WriteModel), nil
}

func (c *Commands) getIAMWriteModel(ctx context.Context) (_ *IAMWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
