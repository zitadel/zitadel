package command

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) SetInstanceCustomText(ctx context.Context, customText *domain.CustomText) (*domain.CustomText, error) {
	setText := NewInstanceCustomTextWriteModel(ctx, customText.Key, customText.Language)
	instanceAgg := InstanceAggregateFromWriteModel(&setText.CustomTextWriteModel.WriteModel)
	event, err := c.setDefaultCustomText(ctx, instanceAgg, setText, customText)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(setText, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToCustomText(&setText.CustomTextWriteModel), nil
}

func (c *Commands) setDefaultCustomText(ctx context.Context, instanceAgg *eventstore.Aggregate, addedPolicy *InstanceCustomTextWriteModel, text *domain.CustomText) (eventstore.Command, error) {
	if !text.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-3MN0s", "Errors.CustomText.Invalid")
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	return instance.NewCustomTextSetEvent(
		ctx,
		instanceAgg,
		text.Template,
		text.Key,
		text.Text,
		text.Language), nil
}

func (c *Commands) defaultCustomTextWriteModelByID(ctx context.Context, key string, language language.Tag) (policy *InstanceCustomTextWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceCustomTextWriteModel(ctx, key, language)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
