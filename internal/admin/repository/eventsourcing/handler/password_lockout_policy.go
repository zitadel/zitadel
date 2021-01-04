package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	passwordLockoutPolicyTable = "adminapi.password_lockout_policies"
)

type PasswordLockoutPolicy struct {
	handler
	subscription *eventstore.Subscription
}

func newPasswordLockoutPolicy(handler handler) *PasswordLockoutPolicy {
	h := &PasswordLockoutPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (p *PasswordLockoutPolicy) subscribe() {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(p, event)
		}
	}()
}

func (p *PasswordLockoutPolicy) ViewModel() string {
	return passwordLockoutPolicyTable
}

func (p *PasswordLockoutPolicy) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *PasswordLockoutPolicy) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := p.view.GetLatestPasswordLockoutPolicySequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *PasswordLockoutPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestPasswordLockoutPolicySequence("")
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *PasswordLockoutPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processPasswordLockoutPolicy(event)
	}
	return err
}

func (p *PasswordLockoutPolicy) processPasswordLockoutPolicy(event *models.Event) (err error) {
	policy := new(iam_model.PasswordLockoutPolicyView)
	switch event.Type {
	case iam_es_model.PasswordLockoutPolicyAdded, model.PasswordLockoutPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PasswordLockoutPolicyChanged, model.PasswordLockoutPolicyChanged:
		policy, err = p.view.PasswordLockoutPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.PasswordLockoutPolicyRemoved:
		return p.view.DeletePasswordLockoutPolicy(event.AggregateID, event)
	default:
		return p.view.ProcessedPasswordLockoutPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutPasswordLockoutPolicy(policy, event)
}

func (p *PasswordLockoutPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-nD8sie", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordLockout policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestPasswordLockoutPolicyFailedEvent, p.view.ProcessedPasswordLockoutPolicyFailedEvent, p.view.ProcessedPasswordLockoutPolicySequence, p.errorCountUntilSkip)
}

func (p *PasswordLockoutPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdatePasswordLockoutPolicySpoolerRunTimestamp)
}
