package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

const (
	labelPolicyTable = "adminapi.label_policies"
)

type LabelPolicy struct {
	handler
	subscription *v1.Subscription
}

func newLabelPolicy(handler handler) *LabelPolicy {
	h := &LabelPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (p *LabelPolicy) subscribe() {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(p, event)
		}
	}()
}

func (p *LabelPolicy) ViewModel() string {
	return labelPolicyTable
}

func (p *LabelPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.IAMAggregate}
}

func (p *LabelPolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestLabelPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *LabelPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestLabelPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *LabelPolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.IAMAggregate:
		err = p.processLabelPolicy(event)
	}
	return err
}

func (p *LabelPolicy) processLabelPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.LabelPolicyView)
	switch event.Type {
	case model.LabelPolicyAdded:
		err = policy.AppendEvent(event)
	case model.LabelPolicyChanged,
		model.LabelPolicyLogoAdded,
		model.LabelPolicyLogoRemoved,
		model.LabelPolicyIconAdded,
		model.LabelPolicyIconRemoved,
		model.LabelPolicyLogoDarkAdded,
		model.LabelPolicyLogoDarkRemoved,
		model.LabelPolicyIconDarkAdded,
		model.LabelPolicyIconDarkRemoved,
		model.LabelPolicyActivated:
		policy, err = p.view.LabelPolicyByAggregateIDAndState(event.AggregateID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	default:
		return p.view.ProcessedLabelPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutLabelPolicy(policy, event)
}

func (p *LabelPolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Wj8sf", "id", event.AggregateID).WithError(err).Warn("something went wrong in label policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestLabelPolicyFailedEvent, p.view.ProcessedLabelPolicyFailedEvent, p.view.ProcessedLabelPolicySequence, p.errorCountUntilSkip)
}

func (p *LabelPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateLabelPolicySpoolerRunTimestamp)
}
