package handler

import (
	"context"
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	loginPolicyTable = "adminapi.login_policies"
)

type LoginPolicy struct {
	handler
	subscription *v1.Subscription
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

func (p *LoginPolicy) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{iam_es_model.IAMAggregate, model.OrgAggregate}
}

func (p *LoginPolicy) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := p.view.GetLatestLoginPolicySequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (p *LoginPolicy) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestLoginPolicySequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *LoginPolicy) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case model.OrgAggregate, iam_es_model.IAMAggregate:
		err = p.processLoginPolicy(event)
	}
	return err
}

func (p *LoginPolicy) processLoginPolicy(event *es_models.Event) (err error) {
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
		policy, err = p.getDefaultLoginPolicy()
		if err != nil {
			return err
		}
		policy.AggregateID = event.AggregateID
		policy.Default = true
	default:
		return p.view.ProcessedLoginPolicySequence(event)
	}
	if err != nil {
		return err
	}
	return p.view.PutLoginPolicy(policy, event)
}

func (p *LoginPolicy) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-Wj8sf", "id", event.AggregateID).WithError(err).Warn("something went wrong in login policy handler")
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

func (p *LoginPolicy) getIAMEvents(sequence uint64) ([]*es_models.Event, error) {
	query, err := eventsourcing.IAMByIDQuery(domain.IAMID, sequence)
	if err != nil {
		return nil, err
	}

	return p.es.FilterEvents(context.Background(), query)
}
