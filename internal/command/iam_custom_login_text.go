package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

func (c *Commands) SetDefaultLoginText(ctx context.Context, loginText *domain.CustomLoginText) (*domain.ObjectDetails, error) {
	iamAgg := iam.NewAggregate()
	events, existingMailText, err := c.setDefaultLoginText(ctx, &iamAgg.Aggregate, loginText)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingMailText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingMailText.WriteModel), nil
}

func (c *Commands) setDefaultLoginText(ctx context.Context, iamAgg *eventstore.Aggregate, text *domain.CustomLoginText) ([]eventstore.EventPusher, *IAMCustomLoginTextReadModel, error) {
	if !text.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "IAM-kd9fs", "Errors.CustomText.Invalid")
	}

	existingLoginText, err := c.defaultLoginTextWriteModelByID(ctx, text.Language)
	if err != nil {
		return nil, nil, err
	}
	events := c.getAllLoginTextEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, true)
	return events, existingLoginText, nil
}

func (c *Commands) defaultLoginTextWriteModelByID(ctx context.Context, lang language.Tag) (*IAMCustomLoginTextReadModel, error) {
	writeModel := NewIAMCustomLoginTextReadModel(lang)
	err := c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
