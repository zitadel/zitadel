package command

import (
	"slices"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type HumanRecoveryCodeWriteModel struct {
	eventstore.WriteModel

	State          domain.MFAState
	FailedAttempts uint64
	codes          []string
	userLocked     bool
}

func (wm *HumanRecoveryCodeWriteModel) Codes() []string {
	if wm.codes == nil {
		return make([]string, 0)
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
			wm.codes = append(wm.Codes(), e.Codes...)
			wm.State = domain.MFAStateReady
		case *user.HumanRecoveryCodeCheckSucceededEvent:
			wm.FailedAttempts = 0
			codeCheckedIndex := slices.Index(wm.codes, e.CodeChecked)
			if codeCheckedIndex == -1 {
				// NB: this should typically never happen, but a race-condition could
				// possibly lead to inconsistent state and duplicate code checks
				continue
			}
			wm.codes = slices.Delete(wm.codes, codeCheckedIndex, codeCheckedIndex+1)
		case *user.HumanRecoveryCodeCheckFailedEvent:
			wm.FailedAttempts += 1
		case *user.HumanRecoveryCodesRemovedEvent:
			wm.codes = make([]string, 0)
			wm.State = domain.MFAStateRemoved
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
