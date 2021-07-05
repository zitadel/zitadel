package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/v1"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	privacyPolicyTable = "auth.privacy_policies"
)

type PrivacyPolicy struct {
	handler
	subscription *v1.Subscription
}

func newPrivacyPolicy(handler handler) *PrivacyPolicy {
	h := &PrivacyPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (p *PrivacyPolicy) subscribe() {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(p, event)
		}
	}()
}

func (p *PrivacyPolicy) ViewModel() string {
	return privacyPolicyTable
}

func (p *PrivacyPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *PrivacyPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestPrivacyPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *PrivacyPolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestPrivacyPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *PrivacyPolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processPrivacyPolicy(event)
	}
	return err
}

func (p *PrivacyPolicy) processPrivacyPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.PrivacyPolicyView)
	switch event.Type {
	case iam_es_model.PrivacyPolicyAdded, model.PrivacyPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PrivacyPolicyChanged, model.PrivacyPolicyChanged:
		policy, err = p.view.PrivacyPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.PrivacyPolicyRemoved:
		return p.view.DeletePrivacyPolicy(event.AggregateID, event)
	default:
		return p.view.ProcessedPrivacyPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutPrivacyPolicy(policy, event)
}

func (p *PrivacyPolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-4N8sw", "id", event.AggregateID).WithError(err).Warn("something went wrong in privacy policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestPrivacyPolicyFailedEvent, p.view.ProcessedPrivacyPolicyFailedEvent, p.view.ProcessedPrivacyPolicySequence, p.errorCountUntilSkip)
}

func (p *PrivacyPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdatePrivacyPolicySpoolerRunTimestamp)
}
