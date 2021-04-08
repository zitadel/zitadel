package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type OrgLoginPolicyWriteModel struct {
	LoginPolicyWriteModel
}

func NewOrgLoginPolicyWriteModel(orgID string) *OrgLoginPolicyWriteModel {
	return &OrgLoginPolicyWriteModel{
		LoginPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgLoginPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LoginPolicyAddedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *org.LoginPolicyChangedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyChangedEvent)
		case *org.LoginPolicyRemovedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyRemovedEvent)
		}
	}
}

func (wm *OrgLoginPolicyWriteModel) IsValid() bool {
	return wm.AggregateID != ""
}

func (wm *OrgLoginPolicyWriteModel) Reduce() error {
	return wm.LoginPolicyWriteModel.Reduce()
}

func (wm *OrgLoginPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.LoginPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			org.LoginPolicyAddedEventType,
			org.LoginPolicyChangedEventType,
			org.LoginPolicyRemovedEventType)
}

func (wm *OrgLoginPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType domain.PasswordlessType,
) (*org.LoginPolicyChangedEvent, bool) {

	changes := make([]policy.LoginPolicyChanges, 0)
	if wm.AllowUserNamePassword != allowUsernamePassword {
		changes = append(changes, policy.ChangeAllowUserNamePassword(allowUsernamePassword))
	}
	if wm.AllowRegister != allowRegister {
		changes = append(changes, policy.ChangeAllowRegister(allowRegister))
	}
	if wm.AllowExternalIDP != allowExternalIDP {
		changes = append(changes, policy.ChangeAllowExternalIDP(allowExternalIDP))
	}
	if wm.ForceMFA != forceMFA {
		changes = append(changes, policy.ChangeForceMFA(forceMFA))
	}
	if passwordlessType.Valid() && wm.PasswordlessType != passwordlessType {
		changes = append(changes, policy.ChangePasswordlessType(passwordlessType))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewLoginPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}

//type AllowedOrgLoginPolicyWriteModel struct {
//	LoginPolicyWriteModel
//	SecondFactorWriteModel
//	MultiFactorWriteModel
//}
//
//func NewAllowedOrgLoginPolicyWriteModel(orgID string) *AllowedOrgLoginPolicyWriteModel {
//	return &AllowedOrgLoginPolicyWriteModel{
//		LoginPolicyWriteModel{
//			WriteModel: eventstore.WriteModel{
//				AggregateID:   orgID,
//				ResourceOwner: orgID,
//			},
//		},
//		SecondFactorWriteModel{
//			WriteModel: eventstore.WriteModel{
//				AggregateID:   orgID,
//				ResourceOwner: orgID,
//			},
//		},
//		MultiFactorWriteModel{
//			WriteModel: eventstore.WriteModel{
//				AggregateID:   orgID,
//				ResourceOwner: orgID,
//			},
//		},
//	}
//}
//
//func (wm *AllowedOrgLoginPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
//	for _, event := range events {
//		switch e := event.(type) {
//		case *org.LoginPolicyAddedEvent:
//			wm.AppendEvents(&e.LoginPolicyAddedEvent)
//		case *org.LoginPolicyChangedEvent:
//			wm.AppendEvents(&e.LoginPolicyChangedEvent)
//		case *org.LoginPolicyRemovedEvent:
//			wm.AppendEvents(&e.LoginPolicyRemovedEvent)
//		case *org.LoginPolicySecondFactorAddedEvent:
//			wm.AppendEvents(&e.SecondFactorAddedEvent)
//		case *org.LoginPolicySecondFactorRemovedEvent:
//			wm.AppendEvents(&e.SecondFactorRemovedEvent)
//		case *org.LoginPolicyMultiFactorAddedEvent:
//			wm.AppendEvents(&e.MultiFactorAddedEvent)
//		case *org.LoginPolicyMultiFactorRemovedEvent:
//			wm.AppendEvents(&e.MultiFactorRemovedEvent)
//		}
//	}
//}
//
//func (wm *AllowedOrgLoginPolicyWriteModel) Reduce() error {
//	err := wm.LoginPolicyWriteModel.Reduce()
//	if err != nil {
//		return err
//	}
//	err = wm.SecondFactorWriteModel.Reduce()
//	if err != nil {
//		return err
//	}
//	return wm.MultiFactorWriteModel.Reduce()
//}
//
//func (wm *AllowedOrgLoginPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
//	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
//		AggregateIDs(wm.LoginPolicyWriteModel.AggregateID).
//		ResourceOwner(wm.ResourceOwner).
//		EventTypes(
//			org.LoginPolicyAddedEventType,
//			org.LoginPolicyChangedEventType,
//			org.LoginPolicyRemovedEventType)
//}
