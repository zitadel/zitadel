package view

import (
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/usergrant/repository/view"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	userGrantTable = "auth.user_grants"
)

func (v *View) UserGrantByID(grantID string) (*model.UserGrantView, error) {
	return view.UserGrantByID(v.Db, userGrantTable, grantID)
}

func (v *View) UserGrantByIDs(resourceOwnerID, projectID, userID string) (*model.UserGrantView, error) {
	return view.UserGrantByIDs(v.Db, userGrantTable, resourceOwnerID, projectID, userID)
}

func (v *View) UserGrantsByUserID(userID string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByUserID(v.Db, userGrantTable, userID)
}

func (v *View) UserGrantsByProjectID(projectID string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByProjectID(v.Db, userGrantTable, projectID)
}

func (v *View) UserGrantsByProjectAndUserID(projectID, userID string) ([]*model.UserGrantView, error) {
	return view.UserGrantsByProjectAndUserID(v.Db, userGrantTable, projectID, userID)
}

func (v *View) SearchUserGrants(request *grant_model.UserGrantSearchRequest) ([]*model.UserGrantView, int, error) {
	return view.SearchUserGrants(v.Db, userGrantTable, request)
}

func (v *View) PutUserGrant(grant *model.UserGrantView, sequence uint64) error {
	err := view.PutUserGrant(v.Db, userGrantTable, grant)
	if err != nil {
		return err
	}
	return v.ProcessedUserGrantSequence(sequence)
}

func (v *View) DeleteUserGrant(grantID string, eventSequence uint64) error {
	err := view.DeleteUserGrant(v.Db, userGrantTable, grantID)
	if err != nil {
		return nil
	}
	return v.ProcessedUserGrantSequence(eventSequence)
}

func (v *View) GetLatestUserGrantSequence() (uint64, error) {
	return v.latestSequence(userGrantTable)
}

func (v *View) ProcessedUserGrantSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(userGrantTable, eventSequence)
}

func (v *View) GetLatestUserGrantFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(userGrantTable, sequence)
}

func (v *View) ProcessedUserGrantFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
