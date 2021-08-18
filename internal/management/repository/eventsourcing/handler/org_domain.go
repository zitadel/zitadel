package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_model "github.com/caos/zitadel/internal/org/repository/view/model"
)

const (
	orgDomainTable = "management.org_domains"
)

type OrgDomain struct {
	handler
	subscription *v1.Subscription
}

func newOrgDomain(handler handler) *OrgDomain {
	h := &OrgDomain{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *OrgDomain) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (d *OrgDomain) ViewModel() string {
	return orgDomainTable
}

func (d *OrgDomain) Subscription() *v1.Subscription {
	return d.subscription
}

func (_ *OrgDomain) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate}
}

func (p *OrgDomain) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestOrgDomainSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (d *OrgDomain) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := d.view.GetLatestOrgDomainSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(d.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (d *OrgDomain) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate:
		err = d.processOrgDomain(event)
	}
	return err
}

func (d *OrgDomain) processOrgDomain(event *es_models.Event) (err error) {
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
		err = d.view.PutOrgDomains(existingDomains, event)
		if err != nil {
			return err
		}
		err = domain.AppendEvent(event)
	case model.OrgDomainRemoved:
		err = domain.SetData(event)
		if err != nil {
			return err
		}
		return d.view.DeleteOrgDomain(event.AggregateID, domain.Domain, event)
	default:
		return d.view.ProcessedOrgDomainSequence(event)
	}
	if err != nil {
		return err
	}
	return d.view.PutOrgDomain(domain, event)
}

func (d *OrgDomain) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-us4sj", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgdomain handler")
	return spooler.HandleError(event, err, d.view.GetLatestOrgDomainFailedEvent, d.view.ProcessedOrgDomainFailedEvent, d.view.ProcessedOrgDomainSequence, d.errorCountUntilSkip)
}

func (o *OrgDomain) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateOrgDomainSpoolerRunTimestamp)
}
