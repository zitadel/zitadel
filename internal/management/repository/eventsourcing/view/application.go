package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	applicationTable = "management.applications"
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

func (v *View) PutApplication(app *model.ApplicationView, event *models.Event) error {
	err := view.PutApplication(v.Db, applicationTable, app)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedApplicationSequence(event)
}

func (v *View) PutApplications(apps []*model.ApplicationView, event *models.Event) error {
	err := view.PutApplications(v.Db, applicationTable, apps...)
	if err != nil {
		return err
	}
	return v.ProcessedApplicationSequence(event)
}

func (v *View) DeleteApplication(appID string, event *models.Event) error {
	err := view.DeleteApplication(v.Db, applicationTable, appID)
	if err != nil {
		return err
	}
	return v.ProcessedApplicationSequence(event)
}

func (v *View) DeleteApplicationsByProjectID(projectID string) error {
	return view.DeleteApplicationsByProjectID(v.Db, applicationTable, projectID)
}

func (v *View) GetLatestApplicationSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(applicationTable)
}

func (v *View) ProcessedApplicationSequence(event *models.Event) error {
	return v.saveCurrentSequence(applicationTable, event)
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
