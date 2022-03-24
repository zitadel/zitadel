package instance

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

func AddLabelPolicy(
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
