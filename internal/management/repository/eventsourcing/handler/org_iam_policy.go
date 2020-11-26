package handler

import (
	"github.com/caos/logging"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type OrgIAMPolicy struct {
	handler
}

const (
	orgIAMPolicyTable = "management.org_iam_policies"
)

func (m *OrgIAMPolicy) ViewModel() string {
	return orgIAMPolicyTable
}

func (m *OrgIAMPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestOrgIAMPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *OrgIAMPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processOrgIAMPolicy(event)
	}
	return err
}

func (m *OrgIAMPolicy) processOrgIAMPolicy(event *models.Event) (err error) {
	policy := new(iam_model.OrgIAMPolicyView)
	switch event.Type {
	case iam_es_model.OrgIAMPolicyAdded, model.OrgIAMPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.OrgIAMPolicyChanged, model.OrgIAMPolicyChanged:
		policy, err = m.view.OrgIAMPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.OrgIAMPolicyRemoved:
		return m.view.DeleteOrgIAMPolicy(event.AggregateID, event.Sequence)
	default:
		return m.view.ProcessedOrgIAMPolicySequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutOrgIAMPolicy(policy, policy.Sequence)
}

func (m *OrgIAMPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-3Gf9s", "id", event.AggregateID).WithError(err).Warn("something went wrong in orgIAM policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestOrgIAMPolicyFailedEvent, m.view.ProcessedOrgIAMPolicyFailedEvent, m.view.ProcessedOrgIAMPolicySequence, m.errorCountUntilSkip)
}

func (o *OrgIAMPolicy) OnSuccess() error {
	return spooler.HandleSuccess(o.view.UpdateOrgIAMPolicySpoolerRunTimestamp)
}
