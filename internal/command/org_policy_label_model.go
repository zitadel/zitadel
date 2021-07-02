package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type OrgLabelPolicyWriteModel struct {
	LabelPolicyWriteModel
}

func NewOrgLabelPolicyWriteModel(orgID string) *OrgLabelPolicyWriteModel {
	return &OrgLabelPolicyWriteModel{
		LabelPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgLabelPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LabelPolicyAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *org.LabelPolicyChangedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyChangedEvent)
		case *org.LabelPolicyRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyRemovedEvent)
		case *org.LabelPolicyLogoAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoAddedEvent)
		case *org.LabelPolicyLogoRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoRemovedEvent)
		case *org.LabelPolicyLogoDarkAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoDarkAddedEvent)
		case *org.LabelPolicyLogoDarkRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyLogoDarkRemovedEvent)
		case *org.LabelPolicyIconAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconAddedEvent)
		case *org.LabelPolicyIconRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconRemovedEvent)
		case *org.LabelPolicyIconDarkAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconDarkAddedEvent)
		case *org.LabelPolicyIconDarkRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyIconDarkRemovedEvent)
		case *org.LabelPolicyFontAddedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyFontAddedEvent)
		case *org.LabelPolicyFontRemovedEvent:
			wm.LabelPolicyWriteModel.AppendEvents(&e.LabelPolicyFontRemovedEvent)
		}
	}
}

func (wm *OrgLabelPolicyWriteModel) Reduce() error {
	return wm.LabelPolicyWriteModel.Reduce()
}

func (wm *OrgLabelPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.LabelPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.LabelPolicyAddedEventType,
			org.LabelPolicyChangedEventType,
			org.LabelPolicyRemovedEventType,
			org.LabelPolicyLogoAddedEventType,
			org.LabelPolicyLogoRemovedEventType,
			org.LabelPolicyIconAddedEventType,
			org.LabelPolicyIconRemovedEventType,
			org.LabelPolicyLogoDarkAddedEventType,
			org.LabelPolicyLogoDarkRemovedEventType,
			org.LabelPolicyIconDarkAddedEventType,
			org.LabelPolicyIconDarkRemovedEventType,
			org.LabelPolicyFontAddedEventType,
			org.LabelPolicyFontRemovedEventType,
		)
}

func (wm *OrgLabelPolicyWriteModel) NewChangedEvent(
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
) (*org.LabelPolicyChangedEvent, bool) {
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
	changedEvent, err := org.NewLabelPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
