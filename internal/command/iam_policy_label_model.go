package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMLabelPolicyWriteModel struct {
	LabelPolicyWriteModel
}

func NewIAMLabelPolicyWriteModel() *IAMLabelPolicyWriteModel {
	return &IAMLabelPolicyWriteModel{
		LabelPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMLabelPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LabelPolicyAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *iam.LabelPolicyChangedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyChangedEvent)
		case *iam.LabelPolicyActivatedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyActivatedEvent)
		case *iam.LabelPolicyLogoAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoAddedEvent)
		case *iam.LabelPolicyLogoRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoRemovedEvent)
		case *iam.LabelPolicyLogoDarkAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoDarkAddedEvent)
		case *iam.LabelPolicyLogoDarkRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoDarkRemovedEvent)
		case *iam.LabelPolicyIconAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconAddedEvent)
		case *iam.LabelPolicyIconRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconRemovedEvent)
		case *iam.LabelPolicyIconDarkAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconDarkAddedEvent)
		case *iam.LabelPolicyIconDarkRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconDarkRemovedEvent)
		case *iam.LabelPolicyFontAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyFontAddedEvent)
		case *iam.LabelPolicyFontRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyFontRemovedEvent)
		}
	}
}

func (wm *IAMLabelPolicyWriteModel) Reduce() error {
	return wm.LabelPolicyWriteModel.Reduce()
}

func (wm *IAMLabelPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.LabelPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.LabelPolicyAddedEventType,
			iam.LabelPolicyChangedEventType,
			iam.LabelPolicyLogoAddedEventType,
			iam.LabelPolicyLogoRemovedEventType,
			iam.LabelPolicyIconAddedEventType,
			iam.LabelPolicyIconRemovedEventType,
			iam.LabelPolicyLogoDarkAddedEventType,
			iam.LabelPolicyLogoDarkRemovedEventType,
			iam.LabelPolicyIconDarkAddedEventType,
			iam.LabelPolicyIconDarkRemovedEventType,
			iam.LabelPolicyFontAddedEventType,
			iam.LabelPolicyFontRemovedEventType,
		)
}

func (wm *IAMLabelPolicyWriteModel) NewChangedEvent(
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
) (*iam.LabelPolicyChangedEvent, bool) {
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
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewLabelPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
