package view

import (
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	orgMemberTable = "management.org_members"
)

func (v *View) OrgMemberByIDs(orgID, userID string) (*model.OrgMemberView, error) {
	return view.OrgMemberByIDs(v.Db, orgMemberTable, orgID, userID)
}

func (v *View) SearchOrgMembers(request *org_model.OrgMemberSearchRequest) ([]*model.OrgMemberView, uint64, error) {
	return view.SearchOrgMembers(v.Db, orgMemberTable, request)
}

func (v *View) OrgMembersByUserID(userID string) ([]*model.OrgMemberView, error) {
	return view.OrgMembersByUserID(v.Db, orgMemberTable, userID)
}

func (v *View) PutOrgMember(member *model.OrgMemberView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutOrgMember(v.Db, orgMemberTable, member)
	if err != nil {
		return err
	}
	return v.ProcessedOrgMemberSequence(sequence, eventTimestamp)
}

func (v *View) PutOrgMembers(members []*model.OrgMemberView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutOrgMembers(v.Db, orgMemberTable, members...)
	if err != nil {
		return err
	}
	return v.ProcessedOrgMemberSequence(sequence, eventTimestamp)
}

func (v *View) DeleteOrgMember(orgID, userID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteOrgMember(v.Db, orgMemberTable, orgID, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedOrgMemberSequence(eventSequence, eventTimestamp)
}

func (v *View) DeleteOrgMembersByUserID(userID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteOrgMembersByUserID(v.Db, orgMemberTable, userID)
	if err != nil {
		return nil
	}
	return v.ProcessedOrgMemberSequence(eventSequence, eventTimestamp)
}

func (v *View) GetLatestOrgMemberSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(orgMemberTable)
}

func (v *View) ProcessedOrgMemberSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(orgMemberTable, eventSequence, eventTimestamp)
}

func (v *View) UpdateOrgMemberSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(orgMemberTable)
}

func (v *View) GetLatestOrgMemberFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(orgMemberTable, sequence)
}

func (v *View) ProcessedOrgMemberFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
