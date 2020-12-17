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

type LoginPolicy struct {
	handler
}

const (
	loginPolicyTable = "auth.login_policies"
)

func (p *LoginPolicy) ViewModel() string {
	return loginPolicyTable
}

func (_ *LoginPolicy) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *LoginPolicy) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := p.view.GetLatestLoginPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *LoginPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestLoginPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *LoginPolicy) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processLoginPolicy(event)
	}
	return err
}

func (p *LoginPolicy) processLoginPolicy(event *models.Event) (err error) {
	policy := new(iam_model.LoginPolicyView)
	switch event.Type {
	case iam_es_model.LoginPolicyAdded, model.LoginPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.LoginPolicyChanged, model.LoginPolicyChanged,
		iam_es_model.LoginPolicySecondFactorAdded, model.LoginPolicySecondFactorAdded,
		iam_es_model.LoginPolicySecondFactorRemoved, model.LoginPolicySecondFactorRemoved,
		iam_es_model.LoginPolicyMultiFactorAdded, model.LoginPolicyMultiFactorAdded,
		iam_es_model.LoginPolicyMultiFactorRemoved, model.LoginPolicyMultiFactorRemoved:
		policy, err = p.view.LoginPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.LoginPolicyRemoved:
		return p.view.DeleteLoginPolicy(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return p.view.ProcessedLoginPolicySequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return p.view.PutLoginPolicy(policy, policy.Sequence, event.CreationDate)
}

func (p *LoginPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-5id9s", "id", event.AggregateID).WithError(err).Warn("something went wrong in login policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestLoginPolicyFailedEvent, p.view.ProcessedLoginPolicyFailedEvent, p.view.ProcessedLoginPolicySequence, p.errorCountUntilSkip)
}

func (p *LoginPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateLoginPolicySpoolerRunTimestamp)
}
