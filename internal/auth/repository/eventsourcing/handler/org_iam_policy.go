package handler

import (
	"github.com/caos/logging"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type OrgIAMPolicy struct {
	handler
}

const (
	orgIAMPolicyTable = "auth.org_iam_policies"
)

func (p *OrgIAMPolicy) ViewModel() string {
	return orgIAMPolicyTable
}

func (_ *OrgIAMPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{org_es_model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *OrgIAMPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestOrgIAMPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *OrgIAMPolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestOrgIAMPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *OrgIAMPolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case org_es_model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processOrgIAMPolicy(event)
	}
	return err
}

func (p *OrgIAMPolicy) processOrgIAMPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.OrgIAMPolicyView)
	switch event.Type {
	case iam_es_model.OrgIAMPolicyAdded, org_es_model.OrgIAMPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.OrgIAMPolicyChanged, org_es_model.OrgIAMPolicyChanged:
		policy, err = p.view.OrgIAMPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case org_es_model.OrgIAMPolicyRemoved:
		return p.view.DeleteOrgIAMPolicy(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return p.view.ProcessedOrgIAMPolicySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutOrgIAMPolicy(policy, policy.Sequence, event.CreationDate)
}

func (p *OrgIAMPolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-3Gj8s", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgIAM policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestOrgIAMPolicyFailedEvent, p.view.ProcessedOrgIAMPolicyFailedEvent, p.view.ProcessedOrgIAMPolicySequence, p.errorCountUntilSkip)
}

func (p *OrgIAMPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateOrgIAMPolicySpoolerRunTimestamp)
}
