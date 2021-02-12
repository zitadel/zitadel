package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	userTable = "adminapi.users"
)

func (v *View) UserByID(userID string) (*model.UserView, error) {
	return view.UserByID(v.Db, userTable, userID)
}

func (v *View) SearchUsers(request *usr_model.UserSearchRequest) ([]*model.UserView, uint64, error) {
	return view.SearchUsers(v.Db, userTable, request)
}

func (v *View) GetGlobalUserByLoginName(loginName string) (*model.UserView, error) {
	return view.GetGlobalUserByLoginName(v.Db, userTable, loginName)
}

func (v *View) UsersByOrgID(orgID string) ([]*model.UserView, error) {
	return view.UsersByOrgID(v.Db, userTable, orgID)
}

func (v *View) UserIDsByDomain(domain string) ([]string, error) {
	return view.UserIDsByDomain(v.Db, userTable, domain)
}

func (v *View) IsUserUnique(userName, email string) (bool, error) {
	return view.IsUserUnique(v.Db, userTable, userName, email)
}

func (v *View) UserMFAs(userID string) ([]*usr_model.MultiFactor, error) {
	return view.UserMFAs(v.Db, userTable, userID)
}

func (v *View) PutUsers(user []*model.UserView, event *models.Event) error {
	err := view.PutUsers(v.Db, userTable, user...)
	if err != nil {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) PutUser(user *model.UserView, event *models.Event) error {
	err := view.PutUser(v.Db, userTable, user)
	if err != nil {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) DeleteUser(userID string, event *models.Event) error {
	err := view.DeleteUser(v.Db, userTable, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) GetLatestUserSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(userTable)
}

func (v *View) ProcessedUserSequence(event *models.Event) error {
	return v.saveCurrentSequence(userTable, event)
}

func (v *View) UpdateUserSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(userTable)
}

func (v *View) GetLatestUserFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userTable, sequence)
}

func (v *View) ProcessedUserFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
