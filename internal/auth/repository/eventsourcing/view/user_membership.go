package view

import (
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	userMembershipTable = "auth.user_memberships"
)

func (v *View) UserMembershipByIDs(userID, aggregateID, objectID string, memberType usr_model.MemberType) (*model.UserMembershipView, error) {
	return view.UserMembershipByIDs(v.Db, userMembershipTable, userID, aggregateID, objectID, memberType)
}

func (v *View) UserMembershipsByAggregateID(aggregateID string) ([]*model.UserMembershipView, error) {
	return view.UserMembershipsByAggregateID(v.Db, userMembershipTable, aggregateID)
}

func (v *View) SearchUserMemberships(request *usr_model.UserMembershipSearchRequest) ([]*model.UserMembershipView, uint64, error) {
	return view.SearchUserMemberships(v.Db, userMembershipTable, request)
}

func (v *View) PutUserMembership(membership *model.UserMembershipView, sequence uint64) error {
	err := view.PutUserMembership(v.Db, userMembershipTable, membership)
	if err != nil {
		return err
	}
	return v.ProcessedUserMembershipSequence(sequence)
}

func (v *View) BulkPutUserMemberships(memberships []*model.UserMembershipView, sequence uint64) error {
	err := view.PutUserMemberships(v.Db, userTable, memberships...)
	if err != nil {
		return err
	}
	return v.ProcessedUserMembershipSequence(sequence)
}

func (v *View) DeleteUserMembership(userID, aggregateID, objectID string, memberType usr_model.MemberType, eventSequence uint64) error {
	err := view.DeleteUserMembership(v.Db, userMembershipTable, userID, aggregateID, objectID, memberType)
	if err != nil {
		return nil
	}
	return v.ProcessedUserMembershipSequence(eventSequence)
}

func (v *View) GetLatestUserMembershipSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(userMembershipTable)
}

func (v *View) ProcessedUserMembershipSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(userMembershipTable, eventSequence)
}

func (v *View) GetLatestUserMembershipFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userMembershipTable, sequence)
}

func (v *View) ProcessedUserMembershipFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
