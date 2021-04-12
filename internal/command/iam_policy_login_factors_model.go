package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMSecondFactorWriteModel struct {
	SecondFactorWriteModel
}

func NewIAMSecondFactorWriteModel(factorType domain.SecondFactorType) *IAMSecondFactorWriteModel {
	return &IAMSecondFactorWriteModel{
		SecondFactorWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			MFAType: factorType,
		},
	}
}

func (wm *IAMSecondFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LoginPolicySecondFactorAddedEvent:
			if wm.MFAType == e.MFAType {
				wm.WriteModel.AppendEvents(&e.SecondFactorAddedEvent)
			}
		case *iam.LoginPolicySecondFactorRemovedEvent:
			if wm.MFAType == e.MFAType {
				wm.WriteModel.AppendEvents(&e.SecondFactorRemovedEvent)
			}
		}
	}
}

func (wm *IAMSecondFactorWriteModel) Reduce() error {
	return wm.SecondFactorWriteModel.Reduce()
}

func (wm *IAMSecondFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.LoginPolicySecondFactorAddedEventType,
			iam.LoginPolicySecondFactorRemovedEventType)
}

type IAMMultiFactorWriteModel struct {
	MultiFactorWriteModel
}

func NewIAMMultiFactorWriteModel(factorType domain.MultiFactorType) *IAMMultiFactorWriteModel {
	return &IAMMultiFactorWriteModel{
		MultiFactorWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			MFAType: factorType,
		},
	}
}

func (wm *IAMMultiFactorWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LoginPolicyMultiFactorAddedEvent:
			if wm.MFAType == e.MFAType {
				wm.WriteModel.AppendEvents(&e.MultiFactorAddedEvent)
			}
		case *iam.LoginPolicyMultiFactorRemovedEvent:
			if wm.MFAType == e.MFAType {
				wm.WriteModel.AppendEvents(&e.MultiFactorRemovedEvent)
			}
		}
	}
}

func (wm *IAMMultiFactorWriteModel) Reduce() error {
	return wm.MultiFactorWriteModel.Reduce()
}

func (wm *IAMMultiFactorWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.WriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner).
		EventTypes(
			iam.LoginPolicyMultiFactorAddedEventType,
			iam.LoginPolicyMultiFactorRemovedEventType)
}

//
//type IAMAuthFactorsWriteModel struct {
//	AuthFactorsWriteModel
//}
//
//func NewIAMAuthFactorsWriteModel() *IAMAuthFactorsWriteModel {
//	return &IAMAuthFactorsWriteModel{
//		AuthFactorsWriteModel{
//			WriteModel: eventstore.WriteModel{
//				AggregateID:   domain.IAMID,
//				ResourceOwner: domain.IAMID,
//			},
//		},
//	}
//}
//
//func (wm *IAMAuthFactorsWriteModel) AppendEvents(events ...eventstore.EventReader) {
//	for _, event := range events {
//		switch e := event.(type) {
//		case *iam.LoginPolicySecondFactorAddedEvent:
//			wm.AuthFactorsWriteModel.AppendEvents(&e.SecondFactorAddedEvent)
//		case *iam.LoginPolicySecondFactorRemovedEvent:
//			wm.AuthFactorsWriteModel.AppendEvents(&e.SecondFactorRemovedEvent)
//		case *iam.LoginPolicyMultiFactorAddedEvent:
//			wm.AuthFactorsWriteModel.AppendEvents(&e.MultiFactorAddedEvent)
//		case *iam.LoginPolicyMultiFactorRemovedEvent:
//			wm.AuthFactorsWriteModel.AppendEvents(&e.MultiFactorRemovedEvent)
//		}
//	}
//}
//
//func (wm *IAMAuthFactorsWriteModel) Reduce() error {
//	return wm.AuthFactorsWriteModel.Reduce()
//}
//
//func (wm *IAMAuthFactorsWriteModel) Query() *eventstore.SearchQueryBuilder {
//	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
//		AggregateIDs(wm.WriteModel.AggregateID).
//		ResourceOwner(wm.ResourceOwner).
//		EventTypes(
//			iam.LoginPolicySecondFactorAddedEventType,
//			iam.LoginPolicySecondFactorRemovedEventType,
//			iam.LoginPolicyMultiFactorAddedEventType,
//			iam.LoginPolicyMultiFactorRemovedEventType)
//}
