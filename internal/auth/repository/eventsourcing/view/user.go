package view

import (
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
)

const (
	userTable = "auth.users"
)

func (v *View) UserByID(userID string) (*model.UserView, error) {
	return view.UserByID(v.Db, userTable, userID)
}

func (v *View) UserByUsername(userName string) (*model.UserView, error) {
	return view.UserByUserName(v.Db, userTable, userName)
}

func (v *View) UserByLoginName(loginName string) (*model.UserView, error) {
	return view.UserByLoginName(v.Db, userTable, loginName)
}

func (v *View) UsersByOrgID(orgID string) ([]*model.UserView, error) {
	return view.UsersByOrgID(v.Db, userTable, orgID)
}
func (v *View) SearchUsers(request *usr_model.UserSearchRequest) ([]*model.UserView, int, error) {
	return view.SearchUsers(v.Db, userTable, request)
}

func (v *View) GetGlobalUserByEmail(email string) (*model.UserView, error) {
	return view.GetGlobalUserByEmail(v.Db, userTable, email)
}

func (v *View) IsUserUnique(userName, email string) (bool, error) {
	return view.IsUserUnique(v.Db, userTable, userName, email)
}

func (v *View) UserMfas(userID string) ([]*usr_model.MultiFactor, error) {
	return view.UserMfas(v.Db, userTable, userID)
}

func (v *View) PutUser(user *model.UserView, sequence uint64) error {
	err := view.PutUser(v.Db, userTable, user)
	if err != nil {
		return err
	}
	return v.ProcessedUserSequence(sequence)
}

func (v *View) DeleteUser(userID string, eventSequence uint64) error {
	err := view.DeleteUser(v.Db, userTable, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedUserSequence(eventSequence)
}

func (v *View) GetLatestUserSequence() (uint64, error) {
	return v.latestSequence(userTable)
}

func (v *View) ProcessedUserSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(userTable, eventSequence)
}

func (v *View) GetLatestUserFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(userTable, sequence)
}

func (v *View) ProcessedUserFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
