package handler

import (
	"context"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/v2/domain"

	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

const (
	loginPolicyTable = "auth.login_policies"
)

type LoginPolicy struct {
	handler
	subscription *eventstore.Subscription
}

func newLoginPolicy(handler handler) *LoginPolicy {
	h := &LoginPolicy{
		handler: handler,
	}

	h.subscribe()

	return h
}

func (p *LoginPolicy) subscribe() {
	p.subscription = p.es.Subscribe(p.AggregateTypes()...)
	go func() {
		for event := range p.subscription.Events {
			query.ReduceEvent(p, event)
		}
	}()
}

func (p *LoginPolicy) ViewModel() string {
	return loginPolicyTable
}

func (_ *LoginPolicy) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{model.OrgAggregate, iam_es_model.IAMAggregate}
}

func (p *LoginPolicy) CurrentSequence(event *models.Event) (uint64, error) {
	sequence, err := p.view.GetLatestLoginPolicySequence(string(event.AggregateType))
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *LoginPolicy) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestLoginPolicySequence("")
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
	case model.OrgAdded:
		policy, err = p.getDefaultLoginPolicy()
		if err != nil {
			return err
		}
		policy.AggregateID = event.AggregateID
		policy.Default = true
	case iam_es_model.LoginPolicyAdded, model.LoginPolicyAdded:
		err = policy.AppendEvent(event)
	case iam_es_model.LoginPolicyChanged,
		iam_es_model.LoginPolicySecondFactorAdded,
		iam_es_model.LoginPolicySecondFactorRemoved,
		iam_es_model.LoginPolicyMultiFactorAdded,
		iam_es_model.LoginPolicyMultiFactorRemoved:
		policies, err := p.view.AllDefaultLoginPolicies()
		if err != nil {
			return err
		}
		for _, policy := range policies {
			err = policy.AppendEvent(event)
			if err != nil {
				return err
			}
		}
		return p.view.PutLoginPolicies(policies, event)
	case model.LoginPolicyChanged,
		model.LoginPolicySecondFactorAdded,
		model.LoginPolicySecondFactorRemoved,
		model.LoginPolicyMultiFactorAdded,
		model.LoginPolicyMultiFactorRemoved:
		policy, err = p.view.LoginPolicyByAggregateID(event.AggregateID)
		if err != nil {
			return err
		}
		err = policy.AppendEvent(event)
	case model.LoginPolicyRemoved:
		return p.view.DeleteLoginPolicy(event.AggregateID, event)
	default:
		return p.view.ProcessedLoginPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutLoginPolicy(policy, event)
}

func (p *LoginPolicy) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-5id9s", "id", event.AggregateID).WithError(err).Warn("something went wrong in login policy handler")
	return spooler.HandleError(event, err, p.view.GetLatestLoginPolicyFailedEvent, p.view.ProcessedLoginPolicyFailedEvent, p.view.ProcessedLoginPolicySequence, p.errorCountUntilSkip)
}

func (p *LoginPolicy) OnSuccess() error {
	return spooler.HandleSuccess(p.view.UpdateLoginPolicySpoolerRunTimestamp)
}

func (p *LoginPolicy) getDefaultLoginPolicy() (*iam_model.LoginPolicyView, error) {
	policy, policyErr := p.view.LoginPolicyByAggregateID(domain.IAMID)
	if policyErr != nil && !caos_errs.IsNotFound(policyErr) {
		return nil, policyErr
	}
	if policy == nil {
		policy = &iam_model.LoginPolicyView{}
	}
	events, err := p.getIAMEvents(policy.Sequence)
	if err != nil {
		return policy, policyErr
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return policy, nil
		}
	}
	return &policyCopy, nil
}

func (p *LoginPolicy) getIAMEvents(sequence uint64) ([]*models.Event, error) {
	query, err := eventsourcing.IAMByIDQuery(domain.IAMID, sequence)
	if err != nil {
		return nil, err
	}

	return p.es.FilterEvents(context.Background(), query)
}
