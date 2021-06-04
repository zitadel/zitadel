package handler

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/static"
)

const (
	labelPolicyTable = "management.label_policies"
)

type LabelPolicy struct {
	handler
	subscription *v1.Subscription
	static       static.Storage
}

func newLabelPolicy(handler handler, static static.Storage) *LabelPolicy {
	h := &LabelPolicy{
		handler: handler,
		static:  static,
	}

	h.subscribe()

	return h
}

func (m *LabelPolicy) subscribe() {
	m.subscription = m.es.Subscribe(m.AggregateTypes()...)
	go func() {
		for event := range m.subscription.Events {
			query.ReduceEvent(m, event)
		}
	}()
}

func (m *LabelPolicy) ViewModel() string {
	return labelPolicyTable
}

func (_ *LabelPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (m *LabelPolicy) CurrentSequence() (uint64, error) {
	sequence, err := m.view.GetLatestLabelPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (m *LabelPolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := m.view.GetLatestLabelPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(m.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (m *LabelPolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processLabelPolicy(event)
	}
	return err
}

func (m *LabelPolicy) processLabelPolicy(event *es_models.Event) (err error) {
	policy := new(iam_model.LabelPolicyView)
	switch event.Type {
	case iam_es_model.LabelPolicyAdded, model.LabelPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.LabelPolicyChanged, model.LabelPolicyChanged,
		iam_es_model.LabelPolicyLogoAdded, model.LabelPolicyLogoAdded,
		iam_es_model.LabelPolicyLogoRemoved, model.LabelPolicyLogoRemoved,
		iam_es_model.LabelPolicyIconAdded, model.LabelPolicyIconAdded,
		iam_es_model.LabelPolicyIconRemoved, model.LabelPolicyIconRemoved,
		iam_es_model.LabelPolicyLogoDarkAdded, model.LabelPolicyLogoDarkAdded,
		iam_es_model.LabelPolicyLogoDarkRemoved, model.LabelPolicyLogoDarkRemoved,
		iam_es_model.LabelPolicyIconDarkAdded, model.LabelPolicyIconDarkAdded,
		iam_es_model.LabelPolicyIconDarkRemoved, model.LabelPolicyIconDarkRemoved,
		iam_es_model.LabelPolicyFontAdded, model.LabelPolicyFontAdded,
		iam_es_model.LabelPolicyFontRemoved, model.LabelPolicyFontRemoved:
		policy, err = m.view.LabelPolicyByAggregateIDAndState(event.AggregateID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case iam_es_model.LabelPolicyActivated, model.LabelPolicyActivated:
		policy, err = m.view.LabelPolicyByAggregateIDAndState(event.AggregateID, int32(domain.LabelPolicyStatePreview))
		if err != nil {
			return err
		}
		go m.CleanUpBucket(policy)
		err = policy.AppendEvent(event)
	default:
		return m.view.ProcessedLabelPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return m.view.PutLabelPolicy(policy, event)
}

func (m *LabelPolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Djo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in label policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestLabelPolicyFailedEvent, m.view.ProcessedLabelPolicyFailedEvent, m.view.ProcessedLabelPolicySequence, m.errorCountUntilSkip)
}

func (m *LabelPolicy) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdateLabelPolicySpoolerRunTimestamp)
}

func (p *LabelPolicy) CleanUpBucket(policy *iam_model.LabelPolicyView) {
	if p.static == nil {
		return
	}
	ctx := context.Background()
	objects, err := p.static.ListObjectInfos(ctx, policy.AggregateID, domain.LabelPolicyPrefix+"/", false)
	if err != nil {
		return
	}
	for _, object := range objects {
		if !deleteableObject(object, policy) {
			continue
		}
		p.static.RemoveObject(ctx, policy.AggregateID, object.Key)
	}
}

func deleteableObject(object *domain.AssetInfo, policy *iam_model.LabelPolicyView) bool {
	if object.Key == policy.LogoURL || object.Key == policy.LogoDarkURL || object.Key == policy.IconURL || object.Key == policy.IconDarkURL || object.Key == policy.FontURL {
		return false
	}
	return true
}
