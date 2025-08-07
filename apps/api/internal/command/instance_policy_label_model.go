package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type InstanceLabelPolicyWriteModel struct {
	LabelPolicyWriteModel
}

func NewInstanceLabelPolicyWriteModel(ctx context.Context) *InstanceLabelPolicyWriteModel {
	return &InstanceLabelPolicyWriteModel{
		LabelPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstanceLabelPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.LabelPolicyAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *instance.LabelPolicyChangedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyChangedEvent)
		case *instance.LabelPolicyActivatedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyActivatedEvent)
		case *instance.LabelPolicyLogoAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoAddedEvent)
		case *instance.LabelPolicyLogoRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoRemovedEvent)
		case *instance.LabelPolicyLogoDarkAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoDarkAddedEvent)
		case *instance.LabelPolicyLogoDarkRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoDarkRemovedEvent)
		case *instance.LabelPolicyIconAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconAddedEvent)
		case *instance.LabelPolicyIconRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconRemovedEvent)
		case *instance.LabelPolicyIconDarkAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconDarkAddedEvent)
		case *instance.LabelPolicyIconDarkRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconDarkRemovedEvent)
		case *instance.LabelPolicyFontAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyFontAddedEvent)
		case *instance.LabelPolicyFontRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyFontRemovedEvent)
		}
	}
}

func (wm *InstanceLabelPolicyWriteModel) Reduce() error {
	return wm.LabelPolicyWriteModel.Reduce()
}

func (wm *InstanceLabelPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.LabelPolicyWriteModel.AggregateID).
		EventTypes(
			instance.LabelPolicyAddedEventType,
			instance.LabelPolicyChangedEventType,
			instance.LabelPolicyLogoAddedEventType,
			instance.LabelPolicyLogoRemovedEventType,
			instance.LabelPolicyIconAddedEventType,
			instance.LabelPolicyIconRemovedEventType,
			instance.LabelPolicyLogoDarkAddedEventType,
			instance.LabelPolicyLogoDarkRemovedEventType,
			instance.LabelPolicyIconDarkAddedEventType,
			instance.LabelPolicyIconDarkRemovedEventType,
			instance.LabelPolicyFontAddedEventType,
			instance.LabelPolicyFontRemovedEventType).
		Builder()
}

func (wm *InstanceLabelPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
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
	themeMode domain.LabelPolicyThemeMode,
) (*instance.LabelPolicyChangedEvent, bool) {
	changes := make([]policy.LabelPolicyChanges, 0)
	if wm.PrimaryColor != primaryColor {
		changes = append(changes, policy.ChangePrimaryColor(primaryColor))
	}
	if wm.BackgroundColor != backgroundColor {
		changes = append(changes, policy.ChangeBackgroundColor(backgroundColor))
	}
	if wm.WarnColor != warnColor {
		changes = append(changes, policy.ChangeWarnColor(warnColor))
	}
	if wm.FontColor != fontColor {
		changes = append(changes, policy.ChangeFontColor(fontColor))
	}
	if wm.PrimaryColorDark != primaryColorDark {
		changes = append(changes, policy.ChangePrimaryColorDark(primaryColorDark))
	}
	if wm.BackgroundColorDark != backgroundColorDark {
		changes = append(changes, policy.ChangeBackgroundColorDark(backgroundColorDark))
	}
	if wm.WarnColorDark != warnColorDark {
		changes = append(changes, policy.ChangeWarnColorDark(warnColorDark))
	}
	if wm.FontColorDark != fontColorDark {
		changes = append(changes, policy.ChangeFontColorDark(fontColorDark))
	}
	if wm.HideLoginNameSuffix != hideLoginNameSuffix {
		changes = append(changes, policy.ChangeHideLoginNameSuffix(hideLoginNameSuffix))
	}
	if wm.ErrorMsgPopup != errorMsgPopup {
		changes = append(changes, policy.ChangeErrorMsgPopup(errorMsgPopup))
	}
	if wm.DisableWatermark != disableWatermark {
		changes = append(changes, policy.ChangeDisableWatermark(disableWatermark))
	}
	if wm.ThemeMode != themeMode {
		changes = append(changes, policy.ChangeThemeMode(themeMode))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewLabelPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
