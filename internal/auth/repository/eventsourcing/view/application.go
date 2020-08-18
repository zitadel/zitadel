package view

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	applicationTable = "auth.applications"
)

func (v *View) ApplicationByID(projectID, appID string) (*model.ApplicationView, error) {
	return view.ApplicationByID(v.Db, applicationTable, projectID, appID)
}

func (v *View) SearchApplications(request *proj_model.ApplicationSearchRequest) ([]*model.ApplicationView, uint64, error) {
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

func (v *View) GetLatestApplicationSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(applicationTable)
}

func (v *View) ProcessedApplicationSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(applicationTable, eventSequence)
}

func (v *View) GetLatestApplicationFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(applicationTable, sequence)
}

func (v *View) ProcessedApplicationFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

func (v *View) ApplicationByClientID(_ context.Context, clientID string) (*model.ApplicationView, error) {
	req := &proj_model.ApplicationSearchRequest{
		Limit: 1,
		Queries: []*proj_model.ApplicationSearchQuery{
			{
				Key:    proj_model.AppSearchKeyOIDCClientID,
				Method: global_model.SearchMethodEquals,
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

func (v *View) AppIDsFromProjectByClientID(ctx context.Context, clientID string) ([]string, error) {
	app, err := v.ApplicationByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	req := &proj_model.ApplicationSearchRequest{
		Queries: []*proj_model.ApplicationSearchQuery{
			{
				Key:    proj_model.AppSearchKeyProjectID,
				Method: global_model.SearchMethodEquals,
				Value:  app.ProjectID,
			},
		},
	}
	apps, _, err := view.SearchApplications(v.Db, applicationTable, req)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "VIEW-Gd24q", "cannot find applications")
	}
	ids := make([]string, 0, len(apps))
	for _, app := range apps {
		if !app.IsOIDC {
			continue
		}
		ids = append(ids, app.OIDCClientID)
	}
	return ids, nil
}
