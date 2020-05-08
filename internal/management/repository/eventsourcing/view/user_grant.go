package view

import (
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
)

const (
	userGrantTable = "management.user_grants"
)

func (v *View) UserGrantByID(grantID string) (*model.UserGrantView, error) {
	return view.UserGrantByID(v.Db, userGrantTable, grantID)
}

func (v *View) SearchUserGrants(request *grant_model.UserGrantSearchRequest) ([]*model.UserGrantView, int, error) {
	return view.SearchUserGrants(v.Db, userGrantTable, request)
}

func (v *View) PutUserGrant(user *model.UserGrantView) error {
	err := view.PutUserGrant(v.Db, userGrantTable, user)
	if err != nil {
		return err
	}
	return v.ProcessedUserGrantSequence(user.Sequence)
}

func (v *View) DeleteUserGrant(grantID string, eventSequence uint64) error {
	err := view.DeleteUserGrant(v.Db, userGrantTable, grantID)
	if err != nil {
		return nil
	}
	return v.ProcessedUserGrantSequence(eventSequence)
}

func (v *View) GetLatestUserGrantSequence() (uint64, error) {
	return v.latestSequence(userGrantTable)
}

func (v *View) ProcessedUserGrantSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(userGrantTable, eventSequence)
}

func (v *View) GetLatestUserGrantFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(userGrantTable, sequence)
}

func (v *View) ProcessedUserGrantFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
