package handler

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
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

func (p *LoginPolicy) ViewModel() string {
	return loginPolicyTable
}

func (p *LoginPolicy) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.IAMAggregate}
}

func (p *LoginPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestLoginPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *LoginPolicy) SetSubscription(s eventstore.Subscription) {
}

func (p *LoginPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = p.processLoginPolicy(event)
	}
	return err
}

func (p *LoginPolicy) processLoginPolicy(event *models.Event) (err error) {
	policy := new(iam_model.LoginPolicyView)
	switch event.Type {
	case model.LoginPolicyAdded:
		err = policy.AppendEvent(event)
	case model.LoginPolicyChanged,
		model.LoginPolicySecondFactorAdded,
		model.LoginPolicySecondFactorRemoved,
		model.LoginPolicyMultiFactorAdded,
		model.LoginPolicyMultiFactorRemoved:
		policy, err = p.view.LoginPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	default:
		return p.view.ProcessedLoginPolicySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutLoginPolicy(policy, policy.Sequence, event.CreationDate)
}

func (p *LoginPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Wj8sf", "id", event.AggregateID).WithError(err).Warn("something went wrong in login policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestLoginPolicyFailedEvent, p.view.ProcessedLoginPolicyFailedEvent, p.view.ProcessedLoginPolicySequence, p.errorCountUntilSkip)
}

func (p *LoginPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateLoginPolicySpoolerRunTimestamp)
}
