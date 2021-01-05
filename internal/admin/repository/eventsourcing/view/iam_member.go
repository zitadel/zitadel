package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view/repository"
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

func (v *View) PutIAMMember(org *model.IAMMemberView, event *models.Event) error {
	err := view.PutIAMMember(v.Db, iamMemberTable, org)
	if err != nil {
		return err
	}
	return v.ProcessedIAMMemberSequence(event)
}

func (v *View) PutIAMMembers(members []*model.IAMMemberView, event *models.Event) error {
	err := view.PutIAMMembers(v.Db, iamMemberTable, members...)
	if err != nil {
		return err
	}
	return v.ProcessedIAMMemberSequence(event)
}

func (v *View) DeleteIAMMember(iamID, userID string, event *models.Event) error {
	err := view.DeleteIAMMember(v.Db, iamMemberTable, iamID, userID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedIAMMemberSequence(event)
}

func (v *View) DeleteIAMMembersByUserID(userID string, event *models.Event) error {
	err := view.DeleteIAMMembersByUserID(v.Db, iamMemberTable, userID)
	if err != nil {
		return err
	}
	return v.ProcessedIAMMemberSequence(event)
}

func (v *View) GetLatestIAMMemberSequence(aggregateType string) (*global_view.CurrentSequence, error) {
	return v.latestSequence(iamMemberTable, aggregateType)
}

func (v *View) ProcessedIAMMemberSequence(event *models.Event) error {
	return v.saveCurrentSequence(iamMemberTable, event)
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
