package password_age

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_age"
)

const (
	AggregateType = "iam"
)

type PasswordAgePolicyWriteModel struct {
	eventstore.WriteModel
	Policy password_age.PasswordAgePolicyWriteModel

	iamID string
}

func NewPasswordAgePolicyWriteModel(iamID string) *PasswordAgePolicyWriteModel {
	return &PasswordAgePolicyWriteModel{
		iamID: iamID,
	}
}

func (wm *PasswordAgePolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordAgePolicyAddedEvent:
			wm.Policy.AppendEvents(&e.PasswordAgePolicyAddedEvent)
		case *PasswordAgePolicyChangedEvent:
			wm.Policy.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		}
	}
}

func (wm *PasswordAgePolicyWriteModel) Reduce() error {
	if err := wm.Policy.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *PasswordAgePolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
