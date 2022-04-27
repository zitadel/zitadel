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
	userMembershipTable = "authz.user_memberships"
)

func (v *View) UserMembershipByIDs(userID, aggregateID, objectID, instanceID string, memberType usr_model.MemberType) (*model.UserMembershipView, error) {
	return view.UserMembershipByIDs(v.Db, userMembershipTable, userID, aggregateID, objectID, instanceID, memberType)
}

func (v *View) UserMembershipsByAggregateID(aggregateID, instanceID string) ([]*model.UserMembershipView, error) {
	return view.UserMembershipsByAggregateID(v.Db, userMembershipTable, aggregateID, instanceID)
}

func (v *View) UserMembershipsByResourceOwner(resourceOwner, instanceID string) ([]*model.UserMembershipView, error) {
	return view.UserMembershipsByResourceOwner(v.Db, userMembershipTable, resourceOwner, instanceID)
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

func (v *View) DeleteUserMembership(userID, aggregateID, objectID, instanceID string, memberType usr_model.MemberType, event *models.Event) error {
	err := view.DeleteUserMembership(v.Db, userMembershipTable, userID, aggregateID, objectID, instanceID, memberType)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) DeleteUserMembershipsByUserID(userID, instanceID string, event *models.Event) error {
	err := view.DeleteUserMembershipsByUserID(v.Db, userMembershipTable, userID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) DeleteUserMembershipsByAggregateID(aggregateID, instanceID string, event *models.Event) error {
	err := view.DeleteUserMembershipsByAggregateID(v.Db, userMembershipTable, aggregateID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) DeleteUserMembershipsByAggregateIDAndObjectID(aggregateID, objectID, instanceID string, event *models.Event) error {
	err := view.DeleteUserMembershipsByAggregateIDAndObjectID(v.Db, userMembershipTable, aggregateID, objectID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserMembershipSequence(event)
}

func (v *View) GetLatestUserMembershipSequence(instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(userMembershipTable, instanceID)
}

func (v *View) GetLatestUserMembershipSequences() ([]*repository.CurrentSequence, error) {
	return v.latestSequences(userMembershipTable)
}

func (v *View) ProcessedUserMembershipSequence(event *models.Event) error {
	return v.saveCurrentSequence(userMembershipTable, event)
}

func (v *View) UpdateUserMembershipSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(userMembershipTable)
}

func (v *View) GetLatestUserMembershipFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userMembershipTable, instanceID, sequence)
}

func (v *View) ProcessedUserMembershipFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
