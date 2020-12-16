package handler

import (
	"context"

	iam_es "github.com/caos/zitadel/internal/iam/repository/eventsourcing"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_events "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type NotifyUser struct {
	handler
	eventstore eventstore.Eventstore
	orgEvents  *org_events.OrgEventstore
	iamEvents  *iam_es.IAMEventstore
	iamID      string
}

const (
	userTable = "notification.notify_users"
)

func (p *NotifyUser) ViewModel() string {
	return userTable
}

func (_ *NotifyUser) AggregateTypes() []models.AggregateType {
	return []models.AggregateType{es_model.UserAggregate, org_es_model.OrgAggregate}
}

func (p *NotifyUser) CurrentSequence() (uint64, error) {
	sequence, err := p.view.GetLatestNotifyUserSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *NotifyUser) EventQuery() (*models.SearchQuery, error) {
	sequence, err := p.view.GetLatestNotifyUserSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *NotifyUser) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case es_model.UserAggregate:
		return u.ProcessUser(event)
	case org_es_model.OrgAggregate:
		return u.ProcessOrg(event)
	default:
		return nil
	}
}

func (u *NotifyUser) ProcessUser(event *models.Event) (err error) {
	user := new(view_model.NotifyUser)
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered,
		es_model.HumanAdded:
		user.AppendEvent(event)
		u.fillLoginNames(user)
	case es_model.UserProfileChanged,
		es_model.UserEmailChanged,
		es_model.UserEmailVerified,
		es_model.UserPhoneChanged,
		es_model.UserPhoneVerified,
		es_model.UserPhoneRemoved,
		es_model.HumanProfileChanged,
		es_model.HumanEmailChanged,
		es_model.HumanEmailVerified,
		es_model.HumanPhoneChanged,
		es_model.HumanPhoneVerified,
		es_model.HumanPhoneRemoved,
		es_model.MachineChanged:
		user, err = u.view.NotifyUserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
	case es_model.DomainClaimed,
		es_model.UserUserNameChanged:
		user, err = u.view.NotifyUserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
		if err != nil {
			return err
		}
		u.fillLoginNames(user)
	case es_model.UserRemoved:
		return u.view.DeleteNotifyUser(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return u.view.ProcessedNotifyUserSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return u.view.PutNotifyUser(user, user.Sequence, event.CreationDate)
}

func (u *NotifyUser) ProcessOrg(event *models.Event) (err error) {
	switch event.Type {
	case org_es_model.OrgDomainVerified,
		org_es_model.OrgDomainRemoved,
		org_es_model.OrgIAMPolicyAdded,
		org_es_model.OrgIAMPolicyChanged,
		org_es_model.OrgIAMPolicyRemoved:
		return u.fillLoginNamesOnOrgUsers(event)
	case org_es_model.OrgDomainPrimarySet:
		return u.fillPreferredLoginNamesOnOrgUsers(event)
	default:
		return u.view.ProcessedNotifyUserSequence(event.Sequence, event.CreationDate)
	}
}

func (u *NotifyUser) fillLoginNamesOnOrgUsers(event *models.Event) error {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(event.ResourceOwner))
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.iamEvents.GetOrgIAMPolicy(context.Background(), u.iamID)
		if err != nil {
			return err
		}
	}
	users, err := u.view.NotifyUsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.SetLoginNames(policy, org.Domains)
		err := u.view.PutNotifyUser(user, 0, event.CreationDate)
		if err != nil {
			return err
		}
	}
	return u.view.ProcessedNotifyUserSequence(event.Sequence, event.CreationDate)
}

func (u *NotifyUser) fillPreferredLoginNamesOnOrgUsers(event *models.Event) error {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(event.ResourceOwner))
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.iamEvents.GetOrgIAMPolicy(context.Background(), u.iamID)
		if err != nil {
			return err
		}
	}
	if !policy.UserLoginMustBeDomain {
		return nil
	}
	users, err := u.view.NotifyUsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
		err := u.view.PutNotifyUser(user, 0, event.CreationDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *NotifyUser) fillLoginNames(user *view_model.NotifyUser) (err error) {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(user.ResourceOwner))
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.iamEvents.GetOrgIAMPolicy(context.Background(), u.iamID)
		if err != nil {
			return err
		}
	}
	user.SetLoginNames(policy, org.Domains)
	user.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
	return nil
}

func (p *NotifyUser) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-9spwf", "id", event.AggregateID).WithError(err).Warn("something went wrong in notify user handler")
	return spooler.HandleError(event, err, p.view.GetLatestNotifyUserFailedEvent, p.view.ProcessedNotifyUserFailedEvent, p.view.ProcessedNotifyUserSequence, p.errorCountUntilSkip)
}

func (u *NotifyUser) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateNotifyUserSpoolerRunTimestamp)
}
