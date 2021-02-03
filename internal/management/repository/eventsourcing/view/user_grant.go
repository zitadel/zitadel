package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	userGrantTable = "management.user_grants"
)

func (v *View) UserGrantByID(grantID string) (*model.UserGrantView, error) {
	return view.UserGrantByID(v.Db, userGrantTable, grantID)
}

func (v *View) SearchUserGrants(request *grant_model.UserGrantSearchRequest) ([]*model.UserGrantView, uint64, error) {
	return view.SearchUserGrants(v.Db, userGrantTable, request)
}

func (v *View) UserGrantsByUserID(userID string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByUserID(v.Db, userGrantTable, userID)
}

func (v *View) UserGrantsByProjectID(projectID string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByProjectID(v.Db, userGrantTable, projectID)
}

func (v *View) UserGrantsByProjectAndGrantID(projectID, grantID string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByProjectAndGrantID(v.Db, userGrantTable, projectID, grantID)
}

func (v *View) UserGrantsByOrgID(orgID string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByOrgID(v.Db, userGrantTable, orgID)
}

func (v *View) UserGrantsByProjectIDAndRoleKey(projectID, roleKey string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByProjectIDAndRole(v.Db, userGrantTable, projectID, roleKey)
}

func (v *View) UserGrantsByOrgIDAndProjectID(orgID, projectID string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByOrgIDAndProjectID(v.Db, userGrantTable, orgID, projectID)
}

func (v *View) PutUserGrant(grant *model.UserGrantView, event *models.Event) error {
	err := view.PutUserGrant(v.Db, userGrantTable, grant)
	if err != nil {
		return err
	}
	return v.ProcessedUserGrantSequence(event)
}

func (v *View) PutUserGrants(grants []*model.UserGrantView, event *models.Event) error {
	err := view.PutUserGrants(v.Db, userGrantTable, grants...)
	if err != nil {
		return err
	}
	return v.ProcessedUserGrantSequence(event)
}

func (v *View) DeleteUserGrant(grantID string, event *models.Event) error {
	err := view.DeleteUserGrant(v.Db, userGrantTable, grantID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedUserGrantSequence(event)
}

func (v *View) GetLatestUserGrantSequence(aggregateType string) (*repository.CurrentSequence, error) {
	return v.latestSequence(userGrantTable, aggregateType)
}

func (v *View) ProcessedUserGrantSequence(event *models.Event) error {
	return v.saveCurrentSequence(userGrantTable, event)
}

func (v *View) UpdateUserGrantSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(userGrantTable)
}

func (v *View) GetLatestUserGrantFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userGrantTable, sequence)
}

func (v *View) ProcessedUserGrantFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
