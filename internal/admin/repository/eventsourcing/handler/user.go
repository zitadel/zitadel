package handler

import (
	"context"
	"github.com/caos/zitadel/internal/config/systemdefaults"
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

type User struct {
	handler
	eventstore     eventstore.Eventstore
	orgEvents      *org_events.OrgEventstore
	iamEvents      *iam_es.IAMEventstore
	systemDefaults systemdefaults.SystemDefaults
}

const (
	userTable = "adminapi.users"
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
		err = u.fillLoginNames(user)
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
		es_model.MFAOTPAdded,
		es_model.MFAOTPVerified,
		es_model.MFAOTPRemoved,
		es_model.HumanProfileChanged,
		es_model.HumanEmailChanged,
		es_model.HumanEmailVerified,
		es_model.HumanPhoneChanged,
		es_model.HumanPhoneVerified,
		es_model.HumanPhoneRemoved,
		es_model.HumanAddressChanged,
		es_model.HumanMFAOTPAdded,
		es_model.HumanMFAOTPVerified,
		es_model.HumanMFAOTPRemoved,
		es_model.MachineChanged:
		user, err = u.view.UserByID(event.AggregateID)
		if err != nil {
			return err
		}
		err = user.AppendEvent(event)
	case es_model.DomainClaimed,
		es_model.UserUserNameChanged:
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
		return u.view.DeleteUser(event.AggregateID, event.Sequence, event.CreationDate)
	default:
		return u.view.ProcessedUserSequence(event.Sequence, event.CreationDate)
	}
	if err != nil {
		return err
	}
	return u.view.PutUser(user, user.Sequence, event.CreationDate)
}

func (u *User) ProcessOrg(event *models.Event) (err error) {
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
		return u.view.ProcessedUserSequence(event.Sequence, event.CreationDate)
	}
}

func (u *User) fillLoginNamesOnOrgUsers(event *models.Event) error {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(event.ResourceOwner))
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.iamEvents.GetOrgIAMPolicy(context.Background(), u.systemDefaults.IamID)
		if err != nil {
			return err
		}
	}
	users, err := u.view.UsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.SetLoginNames(policy, org.Domains)
	}
	return u.view.PutUsers(users, event.Sequence, event.CreationDate)
}

func (u *User) fillPreferredLoginNamesOnOrgUsers(event *models.Event) error {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(event.ResourceOwner))
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.iamEvents.GetOrgIAMPolicy(context.Background(), u.systemDefaults.IamID)
		if err != nil {
			return err
		}
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
	}
	return u.view.PutUsers(users, event.Sequence, event.CreationDate)
}

func (u *User) fillLoginNames(user *view_model.UserView) (err error) {
	org, err := u.orgEvents.OrgByID(context.Background(), org_model.NewOrg(user.ResourceOwner))
	if err != nil {
		return err
	}
	policy := org.OrgIamPolicy
	if policy == nil {
		policy, err = u.iamEvents.GetOrgIAMPolicy(context.Background(), u.systemDefaults.IamID)
		if err != nil {
			return err
		}
	}
	user.SetLoginNames(policy, org.Domains)
	user.PreferredLoginName = user.GenerateLoginName(org.GetPrimaryDomain().Domain, policy.UserLoginMustBeDomain)
	return nil
}

func (u *User) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-vLmwQ", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserFailedEvent, u.view.ProcessedUserFailedEvent, u.view.ProcessedUserSequence, u.errorCountUntilSkip)
}

func (u *User) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserSpoolerRunTimestamp)
}
