package view

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/project/repository/view"
	"github.com/zitadel/zitadel/internal/project/repository/view/model"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	orgProjectMappingTable = "auth.org_project_mapping2"
)

func (v *View) OrgProjectMappingByIDs(orgID, projectID, instanceID string) (*model.OrgProjectMapping, error) {
	return view.OrgProjectMappingByIDs(v.Db, orgProjectMappingTable, orgID, projectID, instanceID)
}

func (v *View) PutOrgProjectMapping(mapping *model.OrgProjectMapping, event *models.Event) error {
	err := view.PutOrgProjectMapping(v.Db, orgProjectMappingTable, mapping)
	if err != nil {
		return err
	}
	return v.ProcessedOrgProjectMappingSequence(event)
}

func (v *View) DeleteOrgProjectMapping(orgID, projectID, instanceID string, event *models.Event) error {
	err := view.DeleteOrgProjectMapping(v.Db, orgProjectMappingTable, orgID, projectID, instanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedOrgProjectMappingSequence(event)
}

func (v *View) DeleteInstanceOrgProjectMappings(event *models.Event) error {
	err := view.DeleteInstanceOrgProjectMappings(v.Db, orgProjectMappingTable, event.InstanceID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedOrgProjectMappingSequence(event)
}

func (v *View) UpdateOwnerRemovedOrgProjectMappings(event *models.Event) error {
	err := view.UpdateOwnerRemovedOrgProjectMappings(v.Db, orgProjectMappingTable, event.InstanceID, event.AggregateID)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedOrgProjectMappingSequence(event)
}

func (v *View) DeleteOrgProjectMappingsByProjectID(projectID, instanceID string) error {
	return view.DeleteOrgProjectMappingsByProjectID(v.Db, orgProjectMappingTable, projectID, instanceID)
}

func (v *View) DeleteOrgProjectMappingsByProjectGrantID(projectGrantID, instanceID string) error {
	return view.DeleteOrgProjectMappingsByProjectGrantID(v.Db, orgProjectMappingTable, projectGrantID, instanceID)
}

func (v *View) GetLatestOrgProjectMappingSequence(instanceID string) (*repository.CurrentSequence, error) {
	return v.latestSequence(orgProjectMappingTable, instanceID)
}

func (v *View) GetLatestOrgProjectMappingSequences(instanceIDs []string) ([]*repository.CurrentSequence, error) {
	return v.latestSequences(orgProjectMappingTable, instanceIDs)
}

func (v *View) ProcessedOrgProjectMappingSequence(event *models.Event) error {
	return v.saveCurrentSequence(orgProjectMappingTable, event)
}

func (v *View) UpdateOrgProjectMappingSpoolerRunTimestamp(instanceIDs []string) error {
	return v.updateSpoolerRunSequence(orgProjectMappingTable, instanceIDs)
}

func (v *View) GetLatestOrgProjectMappingFailedEvent(sequence uint64, instanceID string) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(orgProjectMappingTable, instanceID, sequence)
}

func (v *View) ProcessedOrgProjectMappingFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
