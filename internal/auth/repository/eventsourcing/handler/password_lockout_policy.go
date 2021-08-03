package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	passwordLockoutPolicyTable = "auth.password_complexity_policies"
)

type PasswordLockoutPolicy struct {
	handler
	subscription *v1.Subscription
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

func (p *PasswordLockoutPolicy) Subscription() *v1.Subscription {
	return p.subscription
}

func (_ *PasswordLockoutPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{org_es_model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *PasswordLockoutPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestLockoutPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *PasswordLockoutPolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestLockoutPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *PasswordLockoutPolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case org_es_model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processPasswordLockoutPolicy(event)
	}
	return err
}

func (p *PasswordLockoutPolicy) processPasswordLockoutPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.LockoutPolicyView)
	switch event.Type {
	case iam_es_model.PasswordLockoutPolicyAdded, org_es_model.PasswordLockoutPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PasswordLockoutPolicyChanged, org_es_model.PasswordLockoutPolicyChanged:
		policy, err = p.view.LockoutPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case org_es_model.PasswordLockoutPolicyRemoved:
		return p.view.DeleteLockoutPolicy(event.AggregateID, event)
	default:
		return p.view.ProcessedLockoutPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutLockoutPolicy(policy, event)
}

func (p *PasswordLockoutPolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-0pos2", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordLockout policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestLockoutPolicyFailedEvent, p.view.ProcessedLockoutPolicyFailedEvent, p.view.ProcessedLockoutPolicySequence, p.errorCountUntilSkip)
}

func (p *PasswordLockoutPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateLockoutPolicySpoolerRunTimestamp)
}
