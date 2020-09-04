package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

// ToDo Michi
type LabelPolicy struct {
	handler
}

const (
	labelPolicyTable = "adminapi.label_policies"
)

func (m *LabelPolicy) ViewModel() string {
	return labelPolicyTable
}

func (m *LabelPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestLabelPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *LabelPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = m.processLabelPolicy(event)
	}
	return err
}

func (m *LabelPolicy) processLabelPolicy(event *models.Event) (err error) {
	policy := new(iam_model.LabelPolicyView)
	switch event.Type {
	case model.LabelPolicyAdded:
		err = policy.AppendEvent(event)
	case model.LabelPolicyChanged:
		policy, err = m.view.LabelPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	default:
		return m.view.ProcessedLabelPolicySequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutLabelPolicy(policy, policy.Sequence)
}

func (m *LabelPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Wj8sf", "id", event.AggregateID).WithError(err).Warn("something went wrong in label policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestLabelPolicyFailedEvent, m.view.ProcessedLabelPolicyFailedEvent, m.view.ProcessedLabelPolicySequence, m.errorCountUntilSkip)
}
