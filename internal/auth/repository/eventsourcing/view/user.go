package view

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	userTable = "auth.users2"
)

func (v *View) UserByID(userID, instanceID string) (*model.UserView, error) {
	return view.UserByID(v.Db, userTable, userID, instanceID)
}

func (v *View) UserByLoginName(ctx context.Context, loginName, instanceID string) (*model.UserView, error) {
	queriedUser, err := v.query.GetNotifyUserByLoginName(ctx, true, loginName)
	if err != nil {
		return nil, err
	}

	//nolint: contextcheck // no lint was added because refactor would change too much code
	return view.UserByID(v.Db, userTable, queriedUser.ID, instanceID)
}

func (v *View) UserByLoginNameAndResourceOwner(ctx context.Context, loginName, resourceOwner, instanceID string) (*model.UserView, error) {
	queriedUser, err := v.query.GetNotifyUserByLoginName(ctx, true, loginName)
	if err != nil {
		return nil, err
	}

	//nolint: contextcheck // no lint was added because refactor would change too much code
	user, err := view.UserByID(v.Db, userTable, queriedUser.ID, instanceID)
	if err != nil {
		return nil, err
	}
	if user.ResourceOwner != resourceOwner {
		return nil, zerrors.ThrowNotFound(nil, "VIEW-qScmi", "Errors.User.NotFound")
	}

	return user, nil
}

func (v *View) UserByEmail(ctx context.Context, email, instanceID string) (*model.UserView, error) {
	emailQuery, err := query.NewUserVerifiedEmailSearchQuery(email, query.TextEqualsIgnoreCase)
	if err != nil {
		return nil, err
	}
	return v.userByID(ctx, instanceID, emailQuery)
}

func (v *View) UserByEmailAndResourceOwner(ctx context.Context, email, resourceOwner, instanceID string) (*model.UserView, error) {
	emailQuery, err := query.NewUserVerifiedEmailSearchQuery(email, query.TextEquals)
	if err != nil {
		return nil, err
	}
	resourceOwnerQuery, err := query.NewUserResourceOwnerSearchQuery(resourceOwner, query.TextEquals)
	if err != nil {
		return nil, err
	}

	return v.userByID(ctx, instanceID, emailQuery, resourceOwnerQuery)
}

func (v *View) UserByPhone(ctx context.Context, phone, instanceID string) (*model.UserView, error) {
	phoneQuery, err := query.NewUserVerifiedPhoneSearchQuery(phone, query.TextEquals)
	if err != nil {
		return nil, err
	}
	return v.userByID(ctx, instanceID, phoneQuery)
}

func (v *View) UserByPhoneAndResourceOwner(ctx context.Context, phone, resourceOwner, instanceID string) (*model.UserView, error) {
	phoneQuery, err := query.NewUserVerifiedPhoneSearchQuery(phone, query.TextEquals)
	if err != nil {
		return nil, err
	}
	resourceOwnerQuery, err := query.NewUserResourceOwnerSearchQuery(resourceOwner, query.TextEquals)
	if err != nil {
		return nil, err
	}

	return v.userByID(ctx, instanceID, phoneQuery, resourceOwnerQuery)
}

func (v *View) userByID(ctx context.Context, instanceID string, queries ...query.SearchQuery) (*model.UserView, error) {
	queriedUser, err := v.query.GetNotifyUser(ctx, true, queries...)
	if err != nil {
		return nil, err
	}

	// always load the latest sequence first, so in case the user was not found by id,
	// the sequence will be equal or lower than the actual projection and no events are lost
	sequence, err := v.GetLatestUserSequence(ctx, instanceID)
	logging.WithFields("instanceID", instanceID).
		OnError(err).
		Errorf("could not get current sequence for userByID")

	user, err := view.UserByID(v.Db, userTable, queriedUser.ID, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, err
	}

	if err != nil {
		user = new(model.UserView)
		if sequence != nil {
			user.ChangeDate = sequence.EventCreatedAt
		}
	}

	query, err := view.UserByIDQuery(queriedUser.ID, instanceID, user.ChangeDate, user.EventTypes())
	if err != nil {
		return nil, err
	}
	events, err := v.es.Filter(ctx, query)
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
		return nil, zerrors.ThrowNotFound(nil, "VIEW-r4y8r", "Errors.User.NotFound")
	}

	return user, nil
}

func (v *View) UsersByOrgID(orgID, instanceID string) ([]*model.UserView, error) {
	return view.UsersByOrgID(v.Db, userTable, orgID, instanceID)
}

func (v *View) PutUser(user *model.UserView, event eventstore.Event) error {
	return view.PutUser(v.Db, userTable, user)
}

func (v *View) PutUsers(users []*model.UserView, event eventstore.Event) error {
	return view.PutUsers(v.Db, userTable, users...)
}

func (v *View) DeleteUser(userID, instanceID string, event eventstore.Event) error {
	err := view.DeleteUser(v.Db, userTable, userID, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) DeleteInstanceUsers(event eventstore.Event) error {
	err := view.DeleteInstanceUsers(v.Db, userTable, event.Aggregate().InstanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) UpdateOrgOwnerRemovedUsers(event eventstore.Event) error {
	err := view.UpdateOrgOwnerRemovedUsers(v.Db, userTable, event.Aggregate().InstanceID, event.Aggregate().ID)
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (v *View) GetLatestUserSequence(ctx context.Context, instanceID string) (_ *query.CurrentState, err error) {
	q := &query.CurrentStateSearchQueries{
		Queries: make([]query.SearchQuery, 2),
	}
	q.Queries[0], err = query.NewCurrentStatesInstanceIDSearchQuery(instanceID)
	if err != nil {
		return nil, err
	}
	q.Queries[1], err = query.NewCurrentStatesProjectionSearchQuery(userTable)
	if err != nil {
		return nil, err
	}
	states, err := v.query.SearchCurrentStates(ctx, q)
	if err != nil || states.SearchResponse.Count == 0 {
		return nil, err
	}
	return states.CurrentStates[0], nil
}
