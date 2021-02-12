package view

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
)

const (
	orgDomainTable = "management.org_domains"
)

func (v *View) OrgDomainByOrgIDAndDomain(orgID, domain string) (*model.OrgDomainView, error) {
	return view.OrgDomainByOrgIDAndDomain(v.Db, orgDomainTable, orgID, domain)
}

func (v *View) OrgDomainsByOrgID(domain string) ([]*model.OrgDomainView, error) {
	return view.OrgDomainsByOrgID(v.Db, orgDomainTable, domain)
}

func (v *View) VerifiedOrgDomain(domain string) (*model.OrgDomainView, error) {
	return view.VerifiedOrgDomain(v.Db, orgDomainTable, domain)
}

func (v *View) SearchOrgDomains(request *org_model.OrgDomainSearchRequest) ([]*model.OrgDomainView, uint64, error) {
	return view.SearchOrgDomains(v.Db, orgDomainTable, request)
}

func (v *View) PutOrgDomain(org *model.OrgDomainView, event *models.Event) error {
	err := view.PutOrgDomain(v.Db, orgDomainTable, org)
	if err != nil {
		return err
	}
	return v.ProcessedOrgDomainSequence(event)
}

func (v *View) PutOrgDomains(domains []*model.OrgDomainView, event *models.Event) error {
	err := view.PutOrgDomains(v.Db, orgDomainTable, domains...)
	if err != nil {
		return err
	}
	return v.ProcessedUserSequence(event)
}

func (v *View) DeleteOrgDomain(orgID, domain string, event *models.Event) error {
	err := view.DeleteOrgDomain(v.Db, orgDomainTable, orgID, domain)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return v.ProcessedOrgDomainSequence(event)
}

func (v *View) GetLatestOrgDomainSequence() (*repository.CurrentSequence, error) {
	return v.latestSequence(orgDomainTable)
}

func (v *View) ProcessedOrgDomainSequence(event *models.Event) error {
	return v.saveCurrentSequence(orgDomainTable, event)
}

func (v *View) UpdateOrgDomainSpoolerRunTimestamp() error {
	return v.updateSpoolerRunSequence(orgDomainTable)
}

func (v *View) GetLatestOrgDomainFailedEvent(sequence uint64) (*repository.FailedEvent, error) {
	return v.latestFailedEvent(orgDomainTable, sequence)
}

func (v *View) ProcessedOrgDomainFailedEvent(failedEvent *repository.FailedEvent) error {
	return v.saveFailedEvent(failedEvent)
}
