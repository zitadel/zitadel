package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	userTable = "auth.users"
)

func (v *View) UserByID(userID, instanceID string) (*model.UserView, error) {
	return view.UserByID(v.Db, userTable, userID, instanceID)
}

func (v *View) UserByUsername(userName, instanceID string) (*model.UserView, error) {
	return view.UserByUserName(v.Db, userTable, userName, instanceID)
}

func (v *View) UserByLoginName(loginName, instanceID string) (*model.UserView, error) {
	return view.UserByLoginName(v.Db, userTable, loginName, instanceID)
}

func (v *View) UserByLoginNameAndResourceOwner(loginName, resourceOwner, instanceID string) (*model.UserView, error) {
	return view.UserByLoginNameAndResourceOwner(v.Db, userTable, loginName, resourceOwner, instanceID)
}

func (v *View) UsersByOrgID(orgID, instanceID string) ([]*model.UserView, error) {
	return view.UsersByOrgID(v.Db, userTable, orgID, instanceID)
}

func (v *View) UserIDsByDomain(domain, instanceID string) ([]string, error) {
	return view.UserIDsByDomain(v.Db, userTable, domain, instanceID)
}

func (v *View) SearchUsers(request *usr_model.UserSearchRequest) ([]*model.UserView, uint64, error) {
	return view.SearchUsers(v.Db, userTable, request)
}

func (v *View) GetGlobalUserByLoginName(email, instanceID string) (*model.UserView, error) {
	return view.GetGlobalUserByLoginName(v.Db, userTable, email, instanceID)
}

func (v *View) UserMFAs(userID, instanceID string) ([]*usr_model.MultiFactor, error) {
	return view.UserMFAs(v.Db, userTable, userID, instanceID)
}

func (v *View) PutUser(user *model.UserView, event *models.Event) error {
	err := view.PutUser(v.Db, userTable, user)
	if err != nil {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) PutUsers(users []*model.UserView, event *models.Event) error {
	err := view.PutUsers(v.Db, userTable, users...)
	if err != nil {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) DeleteUser(userID, instanceID string, event *models.Event) error {
	err := view.DeleteUser(v.Db, userTable, userID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) GetLatestUserSequence(instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(userTable, instanceID)
}

func (v *View) GetLatestUserSequences() ([]*repository.CurrentSequence, error) {
	return v.latestSequences(userTable)
}

func (v *View) ProcessedUserSequence(event *models.Event) error {
	return v.saveCurrentSequence(userTable, event)
}

func (v *View) UpdateUserSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(userTable)
}

func (v *View) GetLatestUserFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userTable, instanceID, sequence)
}

func (v *View) ProcessedUserFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
