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

type PasswordComplexityPolicy struct {
	handler
}

const (
	passwordComplexityPolicyTable = "management.password_complexity_policies"
)

func (p *PasswordComplexityPolicy) ViewModel() string {
	return passwordComplexityPolicyTable
}

func (p *PasswordComplexityPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestPasswordComplexityPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate, iam_es_model.IAMAggregate).
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
		return p.view.DeletePasswordComplexityPolicy(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return p.view.ProcessedPasswordComplexityPolicySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutPasswordComplexityPolicy(policy, policy.Sequence, event.CreationDate)
}

func (p *PasswordComplexityPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-4Djo9", "id", event.AggregateID).WithError(err).Warn("something went wrong in passwordComplexity policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestPasswordComplexityPolicyFailedEvent, p.view.ProcessedPasswordComplexityPolicyFailedEvent, p.view.ProcessedPasswordComplexityPolicySequence, p.errorCountUntilSkip)
}

func (p *PasswordComplexityPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdatePasswordComplexityPolicySpoolerRunTimestamp)
}
