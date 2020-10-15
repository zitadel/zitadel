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

type PasswordComplexityPolicy struct {
	handler
}

const (
	passwordComplexityPolicyTable = "auth.password_complexity_policies"
)

func (m *PasswordComplexityPolicy) ViewModel() string {
	return passwordComplexityPolicyTable
}

func (m *PasswordComplexityPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestPasswordComplexityPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *PasswordComplexityPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processPasswordComplexityPolicy(event)
	}
	return err
}

func (m *PasswordComplexityPolicy) processPasswordComplexityPolicy(event *models.Event) (err error) {
	policy := new(iam_model.PasswordComplexityPolicyView)
	switch event.Type {
	case iam_es_model.PasswordComplexityPolicyAdded, model.PasswordComplexityPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PasswordComplexityPolicyChanged, model.PasswordComplexityPolicyChanged:
		policy, err = m.view.PasswordComplexityPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.PasswordComplexityPolicyRemoved:
		return m.view.DeletePasswordComplexityPolicy(event.AggregateID, event.Sequence)
	default:
		return m.view.ProcessedPasswordComplexityPolicySequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutPasswordComplexityPolicy(policy, policy.Sequence)
}

func (m *PasswordComplexityPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Djo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordComplexity policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestPasswordComplexityPolicyFailedEvent, m.view.ProcessedPasswordComplexityPolicyFailedEvent, m.view.ProcessedPasswordComplexityPolicySequence, m.errorCountUntilSkip)
}
