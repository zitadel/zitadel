package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	passwordComplexityPolicyTable = "management.password_complexity_policies"
)

type PasswordComplexityPolicy struct {
	handler
	subscription *eventstore.Subscription
}

func newPasswordComplexityPolicy(handler handler) *PasswordComplexityPolicy {
	h := &PasswordComplexityPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (m *PasswordComplexityPolicy) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (p *PasswordComplexityPolicy) ViewModel() string {
	return passwordComplexityPolicyTable
}

func (_ *PasswordComplexityPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *PasswordComplexityPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestPasswordComplexityPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *PasswordComplexityPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestPasswordComplexityPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *PasswordComplexityPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processPasswordComplexityPolicy(event)
	}
	return err
}

func (p *PasswordComplexityPolicy) processPasswordComplexityPolicy(event *models.Event) (err error) {
	policy := new(iam_model.PasswordComplexityPolicyView)
	switch event.Type {
	case iam_es_model.PasswordComplexityPolicyAdded, model.PasswordComplexityPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PasswordComplexityPolicyChanged, model.PasswordComplexityPolicyChanged:
		policy, err = p.view.PasswordComplexityPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.PasswordComplexityPolicyRemoved:
		return p.view.DeletePasswordComplexityPolicy(event.AggregateID, event)
	default:
		return p.view.ProcessedPasswordComplexityPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutPasswordComplexityPolicy(policy, event)
}

func (p *PasswordComplexityPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Djo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordComplexity policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestPasswordComplexityPolicyFailedEvent, p.view.ProcessedPasswordComplexityPolicyFailedEvent, p.view.ProcessedPasswordComplexityPolicySequence, p.errorCountUntilSkip)
}

func (p *PasswordComplexityPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdatePasswordComplexityPolicySpoolerRunTimestamp)
}
