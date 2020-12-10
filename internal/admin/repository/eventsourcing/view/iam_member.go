package view

import (
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	iamMemberTable = "adminapi.iam_members"
)

func (v *View) IAMMemberByIDs(orgID, userID string) (*model.IAMMemberView, error) {
	return view.IAMMemberByIDs(v.Db, iamMemberTable, orgID, userID)
}

func (v *View) SearchIAMMembers(request *iam_model.IAMMemberSearchRequest) ([]*model.IAMMemberView, uint64, error) {
	return view.SearchIAMMembers(v.Db, iamMemberTable, request)
}

func (v *View) IAMMembersByUserID(userID string) ([]*model.IAMMemberView, error) {
	return view.IAMMembersByUserID(v.Db, iamMemberTable, userID)
}

func (v *View) PutIAMMember(org *model.IAMMemberView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutIAMMember(v.Db, iamMemberTable, org)
	if err != nil {
		return err
	}
	return v.ProcessedIAMMemberSequence(sequence, eventTimestamp)
}

func (v *View) PutIAMMembers(members []*model.IAMMemberView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutIAMMembers(v.Db, iamMemberTable, members...)
	if err != nil {
		return err
	}
	return v.ProcessedIAMMemberSequence(sequence, eventTimestamp)
}

func (v *View) DeleteIAMMember(iamID, userID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteIAMMember(v.Db, iamMemberTable, iamID, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIAMMemberSequence(eventSequence, eventTimestamp)
}

func (v *View) DeleteIAMMembersByUserID(userID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteIAMMembersByUserID(v.Db, iamMemberTable, userID)
	if err != nil {
		return err
	}
	return v.ProcessedIAMMemberSequence(eventSequence, eventTimestamp)
}

func (v *View) GetLatestIAMMemberSequence() (*global_view.CurrentSequence, error) {
	return v.latestSequence(iamMemberTable)
}

func (v *View) ProcessedIAMMemberSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(iamMemberTable, eventSequence, eventTimestamp)
}

func (v *View) UpdateIAMMemberSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(iamMemberTable)
}

func (v *View) GetLatestIAMMemberFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(iamMemberTable, sequence)
}

func (v *View) ProcessedIAMMemberFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
