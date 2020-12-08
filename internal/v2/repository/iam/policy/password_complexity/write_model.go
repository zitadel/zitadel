package password_complexity

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_complexity"
)

const (
	AggregateType = "iam"
)

type PasswordComplexityPolicyWriteModel struct {
	eventstore.WriteModel
	Policy password_complexity.PasswordComplexityPolicyWriteModel

	iamID string
}

func NewPasswordComplexityPolicyWriteModel(iamID string) *PasswordComplexityPolicyWriteModel {
	return &PasswordComplexityPolicyWriteModel{
		iamID: iamID,
	}
}

func (wm *PasswordComplexityPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordComplexityPolicyAddedEvent:
			wm.Policy.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *PasswordComplexityPolicyChangedEvent:
			wm.Policy.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		}
	}
}

func (wm *PasswordComplexityPolicyWriteModel) Reduce() error {
	if err := wm.Policy.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *PasswordComplexityPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}
