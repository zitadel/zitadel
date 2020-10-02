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

type PasswordLockoutPolicy struct {
	handler
}

const (
	passwordLockoutPolicyTable = "management.password_lockout_policies"
)

func (m *PasswordLockoutPolicy) ViewModel() string {
	return passwordLockoutPolicyTable
}

func (m *PasswordLockoutPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestPasswordLockoutPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *PasswordLockoutPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processPasswordLockoutPolicy(event)
	}
	return err
}

func (m *PasswordLockoutPolicy) processPasswordLockoutPolicy(event *models.Event) (err error) {
	policy := new(iam_model.PasswordLockoutPolicyView)
	switch event.Type {
	case iam_es_model.PasswordLockoutPolicyAdded, model.PasswordLockoutPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PasswordLockoutPolicyChanged, model.PasswordLockoutPolicyChanged:
		policy, err = m.view.PasswordLockoutPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	default:
		return m.view.ProcessedPasswordLockoutPolicySequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutPasswordLockoutPolicy(policy, policy.Sequence)
}

func (m *PasswordLockoutPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Bms8f", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordLockout policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestPasswordLockoutPolicyFailedEvent, m.view.ProcessedPasswordLockoutPolicyFailedEvent, m.view.ProcessedPasswordLockoutPolicySequence, m.errorCountUntilSkip)
}
