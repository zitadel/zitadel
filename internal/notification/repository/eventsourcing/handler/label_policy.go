package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type LabelPolicy struct {
	handler
}

const (
	labelPolicyTable = "management.label_policies"
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
		AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *LabelPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processLabelPolicy(event)
	}
	return err
}

func (m *LabelPolicy) processLabelPolicy(event *models.Event) (err error) {
	policy := new(iam_model.LabelPolicyView)
	switch event.Type {
	case iam_es_model.LabelPolicyAdded, model.LabelPolicyAdded:
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
	logging.LogWithFields("HANDL-lDvDW", "id", event.AggregateID).WithError(err).Warn("something went wrong in label policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestLabelPolicyFailedEvent, m.view.ProcessedLabelPolicyFailedEvent, m.view.ProcessedLabelPolicySequence, m.errorCountUntilSkip)
}
