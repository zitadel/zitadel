package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	query2 "github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/repository/org"
	user_repo "github.com/caos/zitadel/internal/repository/user"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	userTable = "auth.users"
)

type User struct {
	handler
	subscription *v1.Subscription
	queries      *query2.Queries
}

func newUser(
	handler handler,
	queries *query2.Queries,
) *User {
	h := &User{
		handler: handler,
		queries: queries,
	}

	h.subscribe()

	return h
}

func (k *User) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			query.ReduceEvent(k, event)
		}
	}()
}

func (u *User) ViewModel() string {
	return userTable
}

func (u *User) Subscription() *v1.Subscription {
	return u.subscription
}
func (_ *User) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{es_model.UserAggregate, org_es_model.OrgAggregate}
}

func (u *User) CurrentSequence() (uint64, error) {
	sequence, err := u.view.GetLatestUserSequence()
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *User) EventQuery() (*es_models.SearchQuery, error) {
	sequence, err := u.view.GetLatestUserSequence()
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(u.AggregateTypes()...).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *User) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case es_model.UserAggregate:
		return u.ProcessUser(event)
	case org_es_model.OrgAggregate:
		return u.ProcessOrg(event)
	default:
		return nil
	}
}

func (u *User) ProcessUser(event *es_models.Event) (err error) {
	user := new(view_model.UserView)
	switch event.Type {
	case es_model.UserAdded,
		es_model.MachineAdded,
		es_model.HumanAdded,
		es_model.UserRegistered,
		es_model.HumanRegistered:
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
		es_model.MFAInitSkipped,
		es_model.UserPasswordChanged,
		es_model.HumanProfileChanged,
		es_model.HumanEmailChanged,
		es_model.HumanEmailVerified,
		es_model.HumanAvatarAdded,
		es_model.HumanAvatarRemoved,
		es_model.HumanPhoneChanged,
		es_model.HumanPhoneVerified,
		es_model.HumanPhoneRemoved,
		es_model.HumanAddressChanged,
		es_model.HumanMFAOTPAdded,
		es_model.HumanMFAOTPVerified,
		es_model.HumanMFAOTPRemoved,
		es_model.HumanMFAU2FTokenAdded,
		es_model.HumanMFAU2FTokenVerified,
		es_model.HumanMFAU2FTokenRemoved,
		es_model.HumanPasswordlessTokenAdded,
		es_model.HumanPasswordlessTokenVerified,
		es_model.HumanPasswordlessTokenRemoved,
		es_model.HumanMFAInitSkipped,
		es_model.MachineChanged,
		es_model.HumanPasswordChanged,
		es_models.EventType(user_repo.HumanPasswordlessInitCodeAddedType),
		es_models.EventType(user_repo.HumanPasswordlessInitCodeRequestedType):
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
		return u.view.DeleteUser(event.AggregateID, event)
	default:
		return u.view.ProcessedUserSequence(event)
	}
	if err != nil {
		return err
	}
	return u.view.PutUser(user, event)
}

func (u *User) fillLoginNames(user *view_model.UserView) (err error) {
	userLoginMustBeDomain, primaryDomain, domains, err := u.loginNameInformation(context.Background(), user.ResourceOwner)
	if err != nil {
		return err
	}
	user.SetLoginNames(userLoginMustBeDomain, domains)
	user.PreferredLoginName = user.GenerateLoginName(primaryDomain, userLoginMustBeDomain)
	return nil
}

func (u *User) ProcessOrg(event *es_models.Event) (err error) {
	switch event.Type {
	case org_es_model.OrgDomainVerified,
		org_es_model.OrgDomainRemoved,
		es_models.EventType(org.OrgDomainPolicyAddedEventType),
		es_models.EventType(org.OrgDomainPolicyChangedEventType),
		es_models.EventType(org.OrgDomainPolicyRemovedEventType):
		return u.fillLoginNamesOnOrgUsers(event)
	case org_es_model.OrgDomainPrimarySet:
		return u.fillPreferredLoginNamesOnOrgUsers(event)
	default:
		return u.view.ProcessedUserSequence(event)
	}
}

func (u *User) fillLoginNamesOnOrgUsers(event *es_models.Event) error {
	userLoginMustBeDomain, _, domains, err := u.loginNameInformation(context.Background(), event.ResourceOwner)
	if err != nil {
		return err
	}
	users, err := u.view.UsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.SetLoginNames(userLoginMustBeDomain, domains)
	}
	return u.view.PutUsers(users, event)
}

func (u *User) fillPreferredLoginNamesOnOrgUsers(event *es_models.Event) error {
	userLoginMustBeDomain, primaryDomain, _, err := u.loginNameInformation(context.Background(), event.ResourceOwner)
	if err != nil {
		return err
	}
	if !userLoginMustBeDomain {
		return nil
	}
	users, err := u.view.UsersByOrgID(event.AggregateID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.PreferredLoginName = user.GenerateLoginName(primaryDomain, userLoginMustBeDomain)
	}
	return u.view.PutUsers(users, event)
}

func (u *User) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-is8aAWima", "id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserFailedEvent, u.view.ProcessedUserFailedEvent, u.view.ProcessedUserSequence, u.errorCountUntilSkip)
}

func (u *User) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateUserSpoolerRunTimestamp)
}

func (u *User) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-3m9vs", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *User) loginNameInformation(ctx context.Context, orgID string) (userLoginMustBeDomain bool, primaryDomain string, domains []*org_model.OrgDomain, err error) {
	org, err := u.getOrgByID(ctx, orgID)
	if err != nil {
		return false, "", nil, err
	}
	if org.DomainPolicy == nil {
		policy, err := u.queries.DefaultDomainPolicy(ctx)
		if err != nil {
			return false, "", nil, err
		}
		userLoginMustBeDomain = policy.UserLoginMustBeDomain
	}
	return userLoginMustBeDomain, org.GetPrimaryDomain().Domain, org.Domains, nil
}
