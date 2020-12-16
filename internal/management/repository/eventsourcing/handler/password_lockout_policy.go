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

type PasswordLockoutPolicy struct {
	handler
}

const (
	passwordLockoutPolicyTable = "management.password_lockout_policies"
)

func (p *PasswordLockoutPolicy) ViewModel() string {
	return passwordLockoutPolicyTable
}

func (_ *PasswordLockoutPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *PasswordLockoutPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestPasswordLockoutPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *PasswordLockoutPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestPasswordLockoutPolicySequence()
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
		return p.view.DeletePasswordLockoutPolicy(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return p.view.ProcessedPasswordLockoutPolicySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutPasswordLockoutPolicy(policy, policy.Sequence, event.CreationDate)
}

func (p *PasswordLockoutPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-Bms8f", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordLockout policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestPasswordLockoutPolicyFailedEvent, p.view.ProcessedPasswordLockoutPolicyFailedEvent, p.view.ProcessedPasswordLockoutPolicySequence, p.errorCountUntilSkip)
}

func (p *PasswordLockoutPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdatePasswordLockoutPolicySpoolerRunTimestamp)
}
