package view

import (
	org_model "github.com/caos/zitadel/internal/org/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view"
)

const (
	orgTable = "admin_api.orgs"
)

func (v *View) OrgByID(orgID string) (*model.OrgView, error) {
	return org_view.OrgByID(v.Db, orgTable, orgID)
}

func (v *View) SearchOrgs(query *org_model.OrgSearchRequest) ([]*model.OrgView, int, error) {
	return org_view.SearchOrgs(v.Db, orgTable, query)
}

func (v *View) PutOrg(org *model.OrgView) error {
	err := org_view.PutOrg(v.Db, orgTable, org)
	if err != nil {
		return err
	}
	return v.ProcessedOrgSequence(org.Sequence)
}

func (v *View) GetLatestOrgFailedEvent(sequence uint64) (*view.FailedEvent, error) {
	return v.latestFailedEvent(orgTable, sequence)
}

func (v *View) ProcessedOrgFailedEvent(failedEvent *view.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}

func (v *View) GetLatestOrgSequence() (uint64, error) {
	return v.latestSequence(orgTable)
}

func (v *View) ProcessedOrgSequence(eventSequence uint64) error {
	return v.saveCurrentSequence(orgTable, eventSequence)
}
