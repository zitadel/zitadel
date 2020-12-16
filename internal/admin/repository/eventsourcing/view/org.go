package view

import (
	org_model "github.com/caos/zitadel/internal/org/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
	"time"
)

const (
	orgTable = "adminapi.orgs"
)

func (v *View) OrgByID(orgID string) (*model.OrgView, error) {
	return org_view.OrgByID(v.Db, orgTable, orgID)
}

func (v *View) SearchOrgs(query *org_model.OrgSearchRequest) ([]*model.OrgView, uint64, error) {
	return org_view.SearchOrgs(v.Db, orgTable, query)
}

func (v *View) PutOrg(org *model.OrgView, eventTimestamp time.Time) error {
	err := org_view.PutOrg(v.Db, orgTable, org)
	if err != nil {
		return err
	}
	return v.ProcessedOrgSequence(org.Sequence, eventTimestamp)
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

func (v *View) ProcessedOrgSequence(eventSequence uint64, eventTimestamp time.Time) error {
	return v.saveCurrentSequence(orgTable, eventSequence, eventTimestamp)
}
