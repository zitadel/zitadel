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
	events := make([]eventstore.EventPusher, 0)
	events = append(events, c.getSelectLoginTextEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getLoginTextEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getPasswordTextEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getPasswordResetTextEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getInitUserEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getInitDoneEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getInitMFAPromptEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getInitMFAOTPEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getInitMFAU2FEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getInitMFADoneEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getVerifyMFAOTPEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getVerifyMFAU2FEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getRegistrationOptionEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getRegistrationUserEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getRegistrationOrgEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getPasswordlessEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
	events = append(events, c.getSuccessLoginEvents(ctx, iamAgg, &existingLoginText.CustomLoginTextReadModel, text, false)...)
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
