package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/features/repository/view/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	org_repo "github.com/caos/zitadel/internal/repository/org"
)

const (
	featuresTable = "auth.features"
)

type Features struct {
	handler
	subscription *v1.Subscription
}

func newFeatures(handler handler) *Features {
	h := &Features{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (p *Features) subscribe() {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(p, event)
		}
	}()
}

func (p *Features) ViewModel() string {
	return featuresTable
}

func (p *Features) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{iam_es_model.IAMAggregate, org_es_model.OrgAggregate}
}

func (p *Features) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestFeaturesSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *Features) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestFeaturesSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *Features) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case org_es_model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processFeatures(event)
	}
	return err
}

func (p *Features) processFeatures(event *es_models.Event) (err error) {
	features := new(model.FeaturesView)
	switch string(event.Type) {
	case string(org_es_model.OrgAdded):
		features, err = p.getDefaultFeatures()
		if err != nil {
			return err
		}
		features.AggregateID = event.AggregateID
		features.Default = true
	case string(iam_repo.FeaturesSetEventType):
		defaultFeatures, err := p.view.AllDefaultFeatures()
		if err != nil {
			return err
		}
		for _, features := range defaultFeatures {
			err = features.AppendEvent(event)
			if err != nil {
				return err
			}
		}
		return p.view.PutFeaturesList(defaultFeatures, event)
	case string(org_repo.FeaturesSetEventType):
		features, err = p.view.FeaturesByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = features.AppendEvent(event)
	case string(org_repo.FeaturesRemovedEventType):
		features, err = p.getDefaultFeatures()
		if err != nil {
			return err
		}
		features.AggregateID = event.AggregateID
		features.Default = true
	default:
		return p.view.ProcessedFeaturesSequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutFeatures(features, event)
}

func (p *Features) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Wj8sf", "id", event.AggregateID).WithError(err).Warn("something went wrong in login features handler")
	return spooler.HandleError(event, err, p.view.GetLatestFeaturesFailedEvent, p.view.ProcessedFeaturesFailedEvent, p.view.ProcessedFeaturesSequence, p.errorCountUntilSkip)
}

func (p *Features) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateFeaturesSpoolerRunTimestamp)
}

func (p *Features) getDefaultFeatures() (*model.FeaturesView, error) {
	features, featuresErr := p.view.FeaturesByAggregateID(domain.IAMID)
	if featuresErr != nil && !caos_errs.IsNotFound(featuresErr) {
		return nil, featuresErr
	}
	if features == nil {
		features = &model.FeaturesView{}
	}
	events, err := p.getIAMEvents(features.Sequence)
	if err != nil {
		return features, featuresErr
	}
	featuresCopy := *features
	for _, event := range events {
		if err := featuresCopy.AppendEvent(event); err != nil {
			return features, nil
		}
	}
	return &featuresCopy, nil
}

func (p *Features) getIAMEvents(sequence uint64) ([]*es_models.Event, error) {
	query, err := eventsourcing.IAMByIDQuery(domain.IAMID, sequence)
	if err != nil {
		return nil, err
	}

	return p.es.FilterEvents(context.Background(), query)
}
