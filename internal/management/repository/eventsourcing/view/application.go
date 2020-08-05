package view

import (
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	applicationTable = "management.applications"
)

func (v *View) ApplicationByID(appID string) (*model.ApplicationView, error) {
	return view.ApplicationByID(v.Db, applicationTable, appID)
}

func (v *View) ApplicationsByProjectID(ProjectID string) ([]*model.ApplicationView, error) {
	return view.ApplicationsByProjectID(v.Db, applicationTable, ProjectID)
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

func (v *View) DeleteApplicationsByProjectID(ProjectID string) error {
	return view.DeleteApplicationsByProjectID(v.Db, applicationTable, ProjectID)
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
