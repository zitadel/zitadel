package view

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	userTable = "auth.users2"
)

func (v *View) UserByID(userID, instanceID string) (*model.UserView, error) {
	return view.UserByID(v.Db, userTable, userID, instanceID)
}

func (v *View) UserByLoginName(loginName, instanceID string) (*model.UserView, error) {
	loginNameQuery, err := query.NewUserLoginNamesSearchQuery(loginName)
	if err != nil {
		return nil, err
	}

	return v.userByID(instanceID, loginNameQuery)
}

func (v *View) UserByLoginNameAndResourceOwner(loginName, resourceOwner, instanceID string) (*model.UserView, error) {
	loginNameQuery, err := query.NewUserLoginNamesSearchQuery(loginName)
	if err != nil {
		return nil, err
	}
	resourceOwnerQuery, err := query.NewUserResourceOwnerSearchQuery(resourceOwner, query.TextEquals)
	if err != nil {
		return nil, err
	}

	return v.userByID(instanceID, loginNameQuery, resourceOwnerQuery)
}

func (v *View) UserByEmail(email, instanceID string) (*model.UserView, error) {
	emailQuery, err := query.NewUserVerifiedEmailSearchQuery(email, query.TextEqualsIgnoreCase)
	if err != nil {
		return nil, err
	}
	return v.userByID(instanceID, emailQuery)
}

func (v *View) UserByEmailAndResourceOwner(email, resourceOwner, instanceID string) (*model.UserView, error) {
	emailQuery, err := query.NewUserVerifiedEmailSearchQuery(email, query.TextEquals)
	if err != nil {
		return nil, err
	}
	resourceOwnerQuery, err := query.NewUserResourceOwnerSearchQuery(resourceOwner, query.TextEquals)
	if err != nil {
		return nil, err
	}

	return v.userByID(instanceID, emailQuery, resourceOwnerQuery)
}

func (v *View) UserByPhone(phone, instanceID string) (*model.UserView, error) {
	phoneQuery, err := query.NewUserVerifiedPhoneSearchQuery(phone, query.TextEquals)
	if err != nil {
		return nil, err
	}
	return v.userByID(instanceID, phoneQuery)
}

func (v *View) UserByPhoneAndResourceOwner(phone, resourceOwner, instanceID string) (*model.UserView, error) {
	phoneQuery, err := query.NewUserVerifiedPhoneSearchQuery(phone, query.TextEquals)
	if err != nil {
		return nil, err
	}
	resourceOwnerQuery, err := query.NewUserResourceOwnerSearchQuery(resourceOwner, query.TextEquals)
	if err != nil {
		return nil, err
	}

	return v.userByID(instanceID, phoneQuery, resourceOwnerQuery)
}

func (v *View) userByID(instanceID string, queries ...query.SearchQuery) (*model.UserView, error) {
	ctx := authz.WithInstanceID(context.Background(), instanceID)

	queriedUser, err := v.query.GetNotifyUser(ctx, true, false, queries...)
	if err != nil {
		return nil, err
	}

	user, err := view.UserByID(v.Db, userTable, queriedUser.ID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	if err != nil {
		user = new(model.UserView)
	}

	query, err := view.UserByIDQuery(queriedUser.ID, instanceID, user.Sequence)
	if err != nil {
		return nil, err
	}
	events, err := v.es.FilterEvents(ctx, query)
	if err != nil && user.Sequence == 0 {
		return nil, err
	} else if err != nil {
		return user, nil
	}

	userCopy := *user

	for _, event := range events {
		if err := user.AppendEvent(event); err != nil {
			return &userCopy, nil
		}
	}

	if user.State == int32(usr_model.UserStateDeleted) {
		return nil, errors.ThrowNotFound(nil, "VIEW-r4y8r", "Errors.User.NotFound")
	}

	return user, nil
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

func (v *View) DeleteInstanceUsers(event *models.Event) error {
	err := view.DeleteInstanceUsers(v.Db, userTable, event.InstanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) UpdateOrgOwnerRemovedUsers(event *models.Event) error {
	err := view.UpdateOrgOwnerRemovedUsers(v.Db, userTable, event.InstanceID, event.AggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) GetLatestUserSequence(instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(userTable, instanceID)
}

func (v *View) GetLatestUserSequences(instanceIDs []string) ([]*repository.CurrentSequence, error) {
	return v.latestSequences(userTable, instanceIDs)
}

func (v *View) ProcessedUserSequence(event *models.Event) error {
	return v.saveCurrentSequence(userTable, event)
}

func (v *View) UpdateUserSpoolerRunTimestamp(instanceIDs []string) error {
	return v.updateSpoolerRunSequence(userTable, instanceIDs)
}

func (v *View) GetLatestUserFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userTable, instanceID, sequence)
}

func (v *View) ProcessedUserFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
