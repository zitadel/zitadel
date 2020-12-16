package handler

import (
	"github.com/caos/logging"
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
	passwordAgePolicyTable = "management.password_age_policies"
)

func (m *PasswordAgePolicy) ViewModel() string {
	return passwordAgePolicyTable
}

func (_ *PasswordAgePolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (o *PasswordAgePolicy) CurrentSequence() (uint64, error) {
	sequence, err := o.view.GetLatestPasswordAgePolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
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

func (m *PasswordAgePolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = m.processPasswordAgePolicy(event)
	}
	return err
}

func (m *PasswordAgePolicy) processPasswordAgePolicy(event *models.Event) (err error) {
	policy := new(iam_model.PasswordAgePolicyView)
	switch event.Type {
	case iam_es_model.PasswordAgePolicyAdded, model.PasswordAgePolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.PasswordAgePolicyChanged, model.PasswordAgePolicyChanged:
		policy, err = m.view.PasswordAgePolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.PasswordAgePolicyRemoved:
		return m.view.DeletePasswordAgePolicy(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return m.view.ProcessedPasswordAgePolicySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return m.view.PutPasswordAgePolicy(policy, policy.Sequence, event.CreationDate)
}

func (m *PasswordAgePolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Bs89f", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordAge policy handler")
	return spooler.HandleError(event, err, m.view.GetLatestPasswordAgePolicyFailedEvent, m.view.ProcessedPasswordAgePolicyFailedEvent, m.view.ProcessedPasswordAgePolicySequence, m.errorCountUntilSkip)
}

func (m *PasswordAgePolicy) OnSuccess() error {
	return spooler.HandleSuccess(m.view.UpdatePasswordAgePolicySpoolerRunTimestamp)
}
