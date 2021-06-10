package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	userMembershipTable = "management.user_memberships"
)

func (v *View) UserMembershipByIDs(userID, aggregateID, objectID string, memberType usr_model.MemberType) (*model.UserMembershipView, error) {
	return view.UserMembershipByIDs(v.Db, userMembershipTable, userID, aggregateID, objectID, memberType)
}

func (v *View) UserMembershipsByAggregateID(aggregateID string) ([]*model.UserMembershipView, error) {
	return view.UserMembershipsByAggregateID(v.Db, userMembershipTable, aggregateID)
}

func (v *View) UserMembershipsByUserID(userID string) ([]*model.UserMembershipView, error) {
	return view.UserMembershipsByUserID(v.Db, userMembershipTable, userID)
}

func (v *View) SearchUserMemberships(request *usr_model.UserMembershipSearchRequest) ([]*model.UserMembershipView, uint64, error) {
	return view.SearchUserMemberships(v.Db, userMembershipTable, request)
}

func (v *View) PutUserMembership(membership *model.UserMembershipView, event *models.Event) error {
	err := view.PutUserMembership(v.Db, userMembershipTable, membership)
	if err != nil {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) BulkPutUserMemberships(memberships []*model.UserMembershipView, event *models.Event) error {
	err := view.PutUserMemberships(v.Db, userMembershipTable, memberships...)
	if err != nil {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) DeleteUserMembership(userID, aggregateID, objectID string, memberType usr_model.MemberType, event *models.Event) error {
	err := view.DeleteUserMembership(v.Db, userMembershipTable, userID, aggregateID, objectID, memberType)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) DeleteUserMembershipsByUserID(userID string, event *models.Event) error {
	err := view.DeleteUserMembershipsByUserID(v.Db, userMembershipTable, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) DeleteUserMembershipsByAggregateID(aggregateID string, event *models.Event) error {
	err := view.DeleteUserMembershipsByAggregateID(v.Db, userMembershipTable, aggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) DeleteUserMembershipsByAggregateIDAndObjectID(aggregateID, objectID string, event *models.Event) error {
	err := view.DeleteUserMembershipsByAggregateIDAndObjectID(v.Db, userMembershipTable, aggregateID, objectID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) GetLatestUserMembershipSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(userMembershipTable)
}

func (v *View) ProcessedUserMembershipSequence(event *models.Event) error {
	return v.saveCurrentSequence(userMembershipTable, event)
}

func (v *View) UpdateUserMembershipSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(userMembershipTable)
}

func (v *View) GetLatestUserMembershipFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userMembershipTable, sequence)
}

func (v *View) ProcessedUserMembershipFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
