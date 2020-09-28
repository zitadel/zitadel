package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type PasswordComplexityPolicy struct {
	handler
}

const (
	passwordComplexityPolicyTable = "adminapi.password_complexity_policies"
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
		AggregateTypeFilter(model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *PasswordComplexityPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = m.processPasswordComplexityPolicy(event)
	}
	return err
}

func (m *PasswordComplexityPolicy) processPasswordComplexityPolicy(event *models.Event) (err error) {
	policy := new(iam_model.PasswordComplexityView)
	switch event.Type {
	case model.PasswordComplexityPolicyAdded:
		err = policy.AppendEvent(event)
	case model.PasswordComplexityPolicyChanged:
		policy, err = m.view.PasswordComplexityPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	default:
		return m.view.ProcessedPasswordComplexityPolicySequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutPasswordComplexityPolicy(policy, policy.Sequence)
}

func (m *PasswordComplexityPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Sj9d", "id", event.AggregateID).WithError(err).Warn("something went wrong in login policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestPasswordComplexityPolicyFailedEvent, m.view.ProcessedPasswordComplexityPolicyFailedEvent, m.view.ProcessedPasswordComplexityPolicySequence, m.errorCountUntilSkip)
}
