package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/repository/view/model"
)

type OrgDomain struct {
	handler
}

const (
	orgDomainTable = "management.org_domains"
)

func (d *OrgDomain) ViewModel() string {
	return orgDomainTable
}

func (d *OrgDomain) EventQuery() (*models.SearchQuery, error) {
	sequence, err := d.view.GetLatestOrgDomainSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (d *OrgDomain) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = d.processOrgDomain(event)
	}
	return err
}

func (d *OrgDomain) processOrgDomain(event *models.Event) (err error) {
	domain := new(org_model.OrgDomainView)
	switch event.Type {
	case model.OrgDomainAdded:
		err = domain.AppendEvent(event)
	case model.OrgDomainVerified,
		model.OrgDomainVerificationAdded:
		err = domain.SetData(event)
		if err != nil {
			return err
		}
		domain, err = d.view.OrgDomainByOrgIDAndDomain(event.AggregateID, domain.Domain)
		if err != nil {
			return err
		}
		err = domain.AppendEvent(event)
	case model.OrgDomainPrimarySet:
		err = domain.SetData(event)
		if err != nil {
			return err
		}
		domain, err = d.view.OrgDomainByOrgIDAndDomain(event.AggregateID, domain.Domain)
		if err != nil {
			return err
		}
		existingDomains, err := d.view.OrgDomainsByOrgID(event.AggregateID)
		if err != nil {
			return err
		}
		for _, existingDomain := range existingDomains {
			existingDomain.Primary = false
		}
		err = d.view.PutOrgDomains(existingDomains, 0)
		if err != nil {
			return err
		}
		err = domain.AppendEvent(event)
	case model.OrgDomainRemoved:
		err = domain.SetData(event)
		if err != nil {
			return err
		}
		return d.view.DeleteOrgDomain(event.AggregateID, domain.Domain, event.Sequence)
	default:
		return d.view.ProcessedOrgDomainSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return d.view.PutOrgDomain(domain, domain.Sequence)
}

func (d *OrgDomain) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-us4sj", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgdomain handler")
	return spooler.HandleError(event, err, d.view.GetLatestOrgDomainFailedEvent, d.view.ProcessedOrgDomainFailedEvent, d.view.ProcessedOrgDomainSequence, d.errorCountUntilSkip)
}

func (o *OrgDomain) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateOrgDomainSpoolerRunTimestamp)
}
