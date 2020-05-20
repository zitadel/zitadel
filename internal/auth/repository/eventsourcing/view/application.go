package view

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	global_view "github.com/caos/zitadel/internal/view"
)

const (
	applicationTable = "auth.applications"
)

func (v *View) ApplicationByID(appID string) (*model.ApplicationView, error) {
	return view.ApplicationByID(v.Db, applicationTable, appID)
}

func (v *View) SearchApplications(request *proj_model.ApplicationSearchRequest) ([]*model.ApplicationView, int, error) {
	return view.SearchApplications(v.Db, applicationTable, request)
}

func (v *View) PutApplication(project *model.ApplicationView) error {
	err := view.PutApplication(v.Db, applicationTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedApplicationSequence(project.Sequence)
}

func (v *View) DeleteApplication(appID string, eventSequence uint64) error {
	err := view.DeleteApplication(v.Db, applicationTable, appID)
	if err != nil {
		return nil
	}
	return v.ProcessedApplicationSequence(eventSequence)
}

func (v *View) GetLatestApplicationSequence() (uint64, error) {
	return v.latestSequence(applicationTable)
}

func (v *View) ProcessedApplicationSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(applicationTable, eventSequence)
}

func (v *View) GetLatestApplicationFailedEvent(sequence uint64) (*global_view.FailedEvent, error) {
	return v.latestFailedEvent(applicationTable, sequence)
}

func (v *View) ProcessedApplicationFailedEvent(failedEvent *global_view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

func (v *View) ApplicationByClientID(_ context.Context, clientID string) (*model.ApplicationView, error) {
	req := &proj_model.ApplicationSearchRequest{
		Limit: 1,
		Queries: []*proj_model.ApplicationSearchQuery{
			{
				Key:    proj_model.APPLICATIONSEARCHKEY_OIDC_CLIENT_ID,
				Method: global_model.SEARCHMETHOD_EQUALS,
				Value:  clientID,
			},
		},
	}
	apps, count, err := view.SearchApplications(v.Db, applicationTable, req)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "VIEW-sd6JQ", "cannot find client")
	}
	if count != 1 {
		return nil, errors.ThrowPreconditionFailed(nil, "VIEW-dfw3as", "cannot find client")
	}
	return apps[0], nil
}
