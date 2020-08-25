package handler

import (
	"context"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_events "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type User struct {
	handler
	eventstore eventstore.Eventstore
	orgEvents  *org_events.OrgEventstore
}

const (
	userTable = "auth.users"
)

func (u *User) ViewModel() string {
	return userTable
}

func (u *User) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(es_model.UserAggregate, org_es_model.OrgAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *User) Reduce(event *models.Event) (err error) {
	switch event.AggregateType {
	case es_model.UserAggregate:
		return u.ProcessUser(event)
	case org_es_model.OrgAggregate:
		return u.ProcessOrg(event)
	default:
		return nil
	}
}

func (u *User) ProcessUser(event *models.Event) (err error) {
	user := new(view_model.UserView)
	switch event.Type {
	case es_model.UserAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered,
		es_model.MachineAdded,
		es_model.HumanAdded:
		err = user.AppendEvent(event)
		if err != nil {
			return err
		}
		u.fillLoginNames(user)
	case es_model.UserProfileChanged,
		es_model.UserEmailChanged,
		es_model.UserEmailVerified,
		es_model.UserPhoneChanged,
		es_model.UserPhoneVerified,
		es_model.UserPhoneRemoved,
		es_model.UserAddressChanged,
		es_model.UserDeactivated,
		es_model.UserReactivated,
		es_model.UserLocked,
		es_model.UserUnlocked,
		es_model.MfaOtpAdded,
		es_model.MfaOtpVerified,
		es_model.MfaOtpRemoved,
		es_model.MfaInitSkipped,
		es_model.UserPasswordChanged,
		es_model.HumanProfileChanged,
		es_model.HumanEmailChanged,
		es_model.HumanEmailVerified,
		es_model.HumanPhoneChanged,
		es_model.HumanPhoneVerified,
		es_model.HumanPhoneRemoved,
		es_model.HumanAddressChanged,
		es_model.HumanMfaOtpAdded,
		es_model.HumanMfaOtpVerified,
		es_model.HumanMfaOtpRemoved,
		es_model.MachineChanged,
		es_model.HumanPasswordChanged:
		user, err = u.view.UserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
	case es_model.DomainClaimed:
		user, err = u.view.UserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillLoginNames(user)
	case es_model.UserRemoved:
		err = u.view.DeleteUser(event.AggregateID, event.Sequence)
	default:
		return u.view.ProcessedUserSequence(event.Sequence)
	}
	if err != nil {
		return err
	}
	return u.view.PutUser(user, user.Sequence)
}

func (u *User) fillLoginNames(user *view_model.UserView) (err error) {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(user.ResourceOwner))
	if err != nil {
		return err
	}
	policy, err := u.orgEvents.GetOrgIamPolicy(context.Background(), user.ResourceOwner)
	if err != nil {
		return err
	}
	user.SetLoginNames(policy, org.Domains)
	user.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
	return nil
}

func (u *User) ProcessOrg(event *models.Event) (err error) {
	switch event.Type {
	case org_es_model.OrgDomainVerified,
		org_es_model.OrgDomainRemoved,
		org_es_model.OrgIamPolicyAdded,
		org_es_model.OrgIamPolicyChanged,
		org_es_model.OrgIamPolicyRemoved:
		return u.fillLoginNamesOnOrgUsers(event)
	case org_es_model.OrgDomainPrimarySet:
		return u.fillPreferredLoginNamesOnOrgUsers(event)
	default:
		return u.view.ProcessedUserSequence(event.Sequence)
	}
}

func (u *User) fillLoginNamesOnOrgUsers(event *models.Event) error {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(event.ResourceOwner))
	if err != nil {
		return err
	}
	policy, err := u.orgEvents.GetOrgIamPolicy(context.Background(), event.ResourceOwner)
	if err != nil {
		return err
	}
	users, err := u.view.UsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.SetLoginNames(policy, org.Domains)
		err := u.view.PutUser(user, event.Sequence)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) fillPreferredLoginNamesOnOrgUsers(event *models.Event) error {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(event.ResourceOwner))
	if err != nil {
		return err
	}
	policy, err := u.orgEvents.GetOrgIamPolicy(context.Background(), event.ResourceOwner)
	if err != nil {
		return err
	}
	if !policy.UserLoginMustBeDomain {
		return nil
	}
	users, err := u.view.UsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
		err := u.view.PutUser(user, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-is8wa", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserFailedEvent, u.view.ProcessedUserFailedEvent, u.view.ProcessedUserSequence, u.errorCountUntilSkip)
}
