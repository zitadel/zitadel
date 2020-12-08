package login

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
)

const (
	AggregateType = "iam"
)

type LoginPolicyWriteModel struct {
	eventstore.WriteModel
	Policy login.LoginPolicyWriteModel

	iamID string
}

func NewLoginPolicyWriteModel(iamID string) *LoginPolicyWriteModel {
	return &LoginPolicyWriteModel{
		iamID: iamID,
	}
}

func (wm *LoginPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			wm.Policy.AppendEvents(&e.LoginPolicyAddedEvent)
		case *LoginPolicyChangedEvent:
			wm.Policy.AppendEvents(&e.LoginPolicyChangedEvent)
		}
	}
}

func (wm *LoginPolicyWriteModel) Reduce() error {
	if err := wm.Policy.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *LoginPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
