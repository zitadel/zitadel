package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/repository/view"
	"github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	orgPrgojectMappingTable = "auth.org_project_mapping"
)

func (v *View) OrgProjectMappingByIDs(orgID, projectID string) (*model.OrgProjectMapping, error) {
	return view.OrgProjectMappingByIDs(v.Db, orgPrgojectMappingTable, orgID, projectID)
}

func (v *View) PutOrgProjectMapping(mapping *model.OrgProjectMapping, event *models.Event) error {
	err := view.PutOrgProjectMapping(v.Db, orgPrgojectMappingTable, mapping)
	if err != nil {
		return err
	}
	return v.ProcessedOrgProjectMappingSequence(event)
}

func (v *View) DeleteOrgProjectMapping(orgID, projectID string, event *models.Event) error {
	err := view.DeleteOrgProjectMapping(v.Db, orgPrgojectMappingTable, orgID, projectID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedOrgProjectMappingSequence(event)
}

func (v *View) DeleteOrgProjectMappingsByProjectID(projectID string) error {
	return view.DeleteOrgProjectMappingsByProjectID(v.Db, orgPrgojectMappingTable, projectID)
}

func (v *View) GetLatestOrgProjectMappingSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(orgPrgojectMappingTable)
}

func (v *View) ProcessedOrgProjectMappingSequence(event *models.Event) error {
	return v.saveCurrentSequence(orgPrgojectMappingTable, event)
}

func (v *View) UpdateOrgProjectMappingSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(orgPrgojectMappingTable)
}

func (v *View) GetLatestOrgProjectMappingFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(orgPrgojectMappingTable, sequence)
}

func (v *View) ProcessedOrgProjectMappingFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
