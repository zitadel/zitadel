package view

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	applicationTable = "auth.applications"
)

func (v *View) ApplicationByID(projectID, appID string) (*model.ApplicationView, error) {
	return view.ApplicationByID(v.Db, applicationTable, projectID, appID)
}

func (v *View) ApplicationsByProjectID(projectID string) ([]*model.ApplicationView, error) {
	return view.ApplicationsByProjectID(v.Db, applicationTable, projectID)
}

func (v *View) SearchApplications(request *proj_model.ApplicationSearchRequest) ([]*model.ApplicationView, uint64, error) {
	return view.SearchApplications(v.Db, applicationTable, request)
}

func (v *View) PutApplication(app *model.ApplicationView, eventTimestamp time.Time) error {
	err := view.PutApplication(v.Db, applicationTable, app)
	if err != nil {
		return err
	}
	return v.ProcessedApplicationSequence(app.Sequence, eventTimestamp)
}

func (v *View) PutApplications(apps []*model.ApplicationView, sequence uint64, eventTimestamp time.Time) error {
	err := view.PutApplications(v.Db, applicationTable, apps...)
	if err != nil {
		return err
	}
	return v.ProcessedApplicationSequence(sequence, eventTimestamp)
}

func (v *View) DeleteApplication(appID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteApplication(v.Db, applicationTable, appID)
	if err != nil {
		return nil
	}
	return v.ProcessedApplicationSequence(eventSequence, eventTimestamp)
}

func (v *View) DeleteApplicationsByProjectID(projectID string) error {
	return view.DeleteApplicationsByProjectID(v.Db, applicationTable, projectID)
}

func (v *View) GetLatestApplicationSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(applicationTable)
}

func (v *View) ProcessedApplicationSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(applicationTable, eventSequence, eventTimestamp)
}

func (v *View) UpdateApplicationSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(applicationTable)
}

func (v *View) GetLatestApplicationFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(applicationTable, sequence)
}

func (v *View) ProcessedApplicationFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

func (v *View) ApplicationByClientID(_ context.Context, clientID string) (*model.ApplicationView, error) {
	return view.ApplicationByOIDCClientID(v.Db, applicationTable, clientID)
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

func (v *View) AppIDsFromProjectID(ctx context.Context, projectID string) ([]string, error) {
	req := &proj_model.ApplicationSearchRequest{
		Queries: []*proj_model.ApplicationSearchQuery{
			{
				Key:    proj_model.AppSearchKeyProjectID,
				Method: global_model.SearchMethodEquals,
				Value:  projectID,
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
