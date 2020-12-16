package view

import (
	"context"
	"time"

	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	applicationTable = "authz.applications"
)

func (v *View) ApplicationByID(projectID, appID string) (*model.ApplicationView, error) {
	return view.ApplicationByID(v.Db, applicationTable, projectID, appID)
}

func (v *View) ApplicationByOIDCClientID(clientID string) (*model.ApplicationView, error) {
	return view.ApplicationByOIDCClientID(v.Db, applicationTable, clientID)
}

func (v *View) ApplicationByProjecIDAndAppName(ctx context.Context, projectID, appName string) (_ *model.ApplicationView, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	return view.ApplicationByProjectIDAndAppName(v.Db, applicationTable, projectID, appName)
}

func (v *View) SearchApplications(request *proj_model.ApplicationSearchRequest) ([]*model.ApplicationView, uint64, error) {
	return view.SearchApplications(v.Db, applicationTable, request)
}

func (v *View) PutApplication(project *model.ApplicationView, eventTimestamp time.Time) error {
	err := view.PutApplication(v.Db, applicationTable, project)
	if err != nil {
		return err
	}
	return v.ProcessedApplicationSequence(project.Sequence, eventTimestamp)
}

func (v *View) DeleteApplication(appID string, eventSequence uint64, eventTimestamp time.Time) error {
	err := view.DeleteApplication(v.Db, applicationTable, appID)
	if err != nil {
		return nil
	}
	return v.ProcessedApplicationSequence(eventSequence, eventTimestamp)
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
