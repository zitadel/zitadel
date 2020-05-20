package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/view"
)

const (
	orgTable = "management.orgs"
)

func (v *View) OrgByID(orgID string) (*org_view.OrgView, error) {
	return org_view.OrgByID(v.Db, orgTable, orgID)
}

func (v *View) OrgByDomain(domain string) (*org_view.OrgView, error) {
	orgs, _, err := org_view.SearchOrgs(v.Db, orgTable, &org_model.OrgSearchRequest{
		Limit: 1,
		Queries: []*org_model.OrgSearchQuery{
			{
				Key:    org_model.ORGSEARCHKEY_ORG_DOMAIN,
				Method: model.SEARCHMETHOD_EQUALS,
				Value:  domain,
			},
		}})
	if err != nil {
		return nil, err
	}
	if len(orgs) == 0 {
		return nil, errors.ThrowNotFound(nil, "VIEW-ByecF", "no org found")
	}
	return orgs[0], nil
}

func (v *View) PutOrg(org *org_view.OrgView) error {
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
