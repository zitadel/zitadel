package handler

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type PasswordAgePolicy struct {
	handler
}

const (
	passwordAgePolicyTable = "adminapi.password_age_policies"
)

func (p *PasswordAgePolicy) ViewModel() string {
	return passwordAgePolicyTable
}

func (p *PasswordAgePolicy) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *PasswordAgePolicy) SetSubscription(s eventstore.Subscription) {
}

func (p *PasswordAgePolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestPasswordAgePolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *PasswordAgePolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processPasswordAgePolicy(event)
	}
	return err
}

func (p *PasswordAgePolicy) processPasswordAgePolicy(event *models.Event) (err error) {
	policy := new(iam_model.PasswordAgePolicyView)
	switch event.Type {
	case iam_es_model.PasswordAgePolicyAdded, model.PasswordAgePolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PasswordAgePolicyChanged, model.PasswordAgePolicyChanged:
		policy, err = p.view.PasswordAgePolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.PasswordAgePolicyRemoved:
		return p.view.DeletePasswordAgePolicy(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return p.view.ProcessedPasswordAgePolicySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutPasswordAgePolicy(policy, policy.Sequence, event.CreationDate)
}

func (p *PasswordAgePolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-nD8sie", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordAge policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestPasswordAgePolicyFailedEvent, p.view.ProcessedPasswordAgePolicyFailedEvent, p.view.ProcessedPasswordAgePolicySequence, p.errorCountUntilSkip)
}

func (p *PasswordAgePolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateProcessedPasswordAgePolicySpoolerRunTimestamp)
}
