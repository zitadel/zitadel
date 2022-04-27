package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func AddDefaultLabelPolicy(
	a *instance.Aggregate,
	primaryColor,
	backgroundColor,
	warnColor,
	fontColor,
	primaryColorDark,
	backgroundColorDark,
	warnColorDark,
	fontColorDark string,
	hideLoginNameSuffix,
	errorMsgPopup,
	disableWatermark bool,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewLabelPolicyAddedEvent(ctx, &a.Aggregate,
					primaryColor,
					backgroundColor,
					warnColor,
					fontColor,
					primaryColorDark,
					backgroundColorDark,
					warnColorDark,
					fontColorDark,
					hideLoginNameSuffix,
					errorMsgPopup,
					disableWatermark,
				),
			}, nil
		}, nil
	}
}

func ActivateDefaultLabelPolicy(
	a *instance.Aggregate,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			//TODO: check if already exists
			return []eventstore.Command{
				instance.NewLabelPolicyActivatedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}
