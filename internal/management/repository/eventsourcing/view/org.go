package view

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	orgTable = "management.orgs"
)

func (v *View) OrgByID(orgID string) (*query.Org, error) {
	return v.query.OrgByID(context.TODO(), orgID)
}

func (v *View) PutOrg(org *model.OrgView, event *models.Event) error {
	err := org_view.PutOrg(v.Db, orgTable, org)
	if err != nil {
		return err
	}
	return v.ProcessedOrgSequence(event)
}

func (v *View) GetLatestOrgFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(orgTable, sequence)
}

func (v *View) ProcessedOrgFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

func (v *View) UpdateOrgSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(orgTable)
}

func (v *View) GetLatestOrgSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(orgTable)
}

func (v *View) ProcessedOrgSequence(event *models.Event) error {
	return v.saveCurrentSequence(orgTable, event)
}
