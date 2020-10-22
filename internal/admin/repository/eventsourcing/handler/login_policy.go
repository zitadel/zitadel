package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type LoginPolicy struct {
	handler
}

const (
	loginPolicyTable = "adminapi.login_policies"
)

func (m *LoginPolicy) ViewModel() string {
	return loginPolicyTable
}

func (m *LoginPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := m.view.GetLatestLoginPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.IAMAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *LoginPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = m.processLoginPolicy(event)
	}
	return err
}

func (m *LoginPolicy) processLoginPolicy(event *models.Event) (err error) {
	policy := new(iam_model.LoginPolicyView)
	switch event.Type {
	case model.LoginPolicyAdded:
		err = policy.AppendEvent(event)
	case model.LoginPolicyChanged,
		model.LoginPolicySoftwareMFAAdded,
		model.LoginPolicySoftwareMFARemoved,
		model.LoginPolicyHardwareMFAAdded,
		model.LoginPolicyHardwareMFARemoved:
		policy, err = m.view.LoginPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	default:
		return m.view.ProcessedLoginPolicySequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return m.view.PutLoginPolicy(policy, policy.Sequence)
}

func (m *LoginPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Wj8sf", "id", event.AggregateID).WithError(err).Warn("something went wrong in login policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestLoginPolicyFailedEvent, m.view.ProcessedLoginPolicyFailedEvent, m.view.ProcessedLoginPolicySequence, m.errorCountUntilSkip)
}
