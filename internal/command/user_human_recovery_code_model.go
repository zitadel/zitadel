package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type HumanRecoveryCodeWriteModel struct {
	eventstore.WriteModel

	State          domain.MFAState
	codes          []domain.RecoveryCode
	userLocked     bool
	FailedAttempts uint64
}

func (wm *HumanRecoveryCodeWriteModel) Codes() []domain.RecoveryCode {
	if wm.codes == nil {
		return []domain.RecoveryCode{}
	}
	return wm.codes
}

func (wm *HumanRecoveryCodeWriteModel) UserLocked() bool {
	return wm.userLocked
}

func NewHumanRecoveryCodeWriteModel(userID, resourceOwner string) *HumanRecoveryCodeWriteModel {
	return &HumanRecoveryCodeWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanRecoveryCodeWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanRecoveryCodesAddedEvent:
			recoveryCodes := make([]domain.RecoveryCode, len(e.Codes))
			for i, code := range e.Codes {
				recoveryCodes[i] = domain.RecoveryCode{
					HashedCode: code,
					CheckAt:    time.Time{},
				}
			}
			wm.codes = recoveryCodes
			wm.State = domain.MFAStateReady
		case *user.HumanRecoveryCodeCheckSucceededEvent:
			wm.FailedAttempts = 0
			wm.codes[e.CodeIndex].CheckAt = e.CreatedAt()
		case *user.HumanRecoveryCodeCheckFailedEvent:
			wm.FailedAttempts += 1
		case *user.HumanRecoveryCodeRemovedEvent:
			wm.State = domain.MFAStateRemoved
			wm.codes = nil
		case *user.UserLockedEvent:
			wm.userLocked = true
		case *user.UserUnlockedEvent:
			wm.userLocked = false
			wm.FailedAttempts = 0
		case *user.UserRemovedEvent:
			wm.State = domain.MFAStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanRecoveryCodeWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanRecoveryCodesAddedType,
			user.HumanRecoveryCodesRemovedType,
			user.HumanRecoveryCodeCheckSucceededType,
			user.HumanRecoveryCodeCheckFailedType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserRemovedType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}
