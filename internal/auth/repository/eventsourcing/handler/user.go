package handler

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/zitadel/zitadel/internal/eventstore/v1/sdk"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	org_es_model "github.com/zitadel/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/zitadel/zitadel/internal/org/repository/view"
	query2 "github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	userTable = "auth.users2"
)

type User struct {
	handler
	subscription *v1.Subscription
	queries      *query2.Queries
}

func newUser(
	ctx context.Context,
	handler handler,
	queries *query2.Queries,
) *User {
	h := &User{
		handler: handler,
		queries: queries,
	}

	h.subscribe(ctx)

	return h
}

func (u *User) subscribe(ctx context.Context) {
	u.subscription = u.es.Subscribe(u.AggregateTypes()...)
	go func() {
		for event := range u.subscription.Events {
			query.ReduceEvent(ctx, u, event)
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
	return []es_models.AggregateType{user_repo.AggregateType, org.AggregateType, instance.AggregateType}
}

func (u *User) CurrentSequence(ctx context.Context, instanceID string) (uint64, error) {
	sequence, err := u.view.GetLatestUserSequence(ctx, instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (u *User) EventQuery(ctx context.Context, instanceIDs []string) (*es_models.SearchQuery, error) {
	sequences, err := u.view.GetLatestUserSequences(ctx, instanceIDs)
	if err != nil {
		return nil, err
	}
	return newSearchQuery(sequences, u.AggregateTypes(), instanceIDs), nil
}

func (u *User) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case user_repo.AggregateType:
		return u.ProcessUser(event)
	case org.AggregateType:
		return u.ProcessOrg(event)
	case instance.AggregateType:
		return u.ProcessInstance(event)
	default:
		return nil
	}
}

func (u *User) ProcessUser(event *es_models.Event) (err error) {
	user := new(view_model.UserView)
	switch eventstore.EventType(event.Type) {
	case user_repo.UserV1AddedType,
		user_repo.MachineAddedEventType,
		user_repo.HumanAddedType,
		user_repo.UserV1RegisteredType,
		user_repo.HumanRegisteredType:
		err = user.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillLoginNames(user)
	case user_repo.UserV1ProfileChangedType,
		user_repo.UserV1EmailChangedType,
		user_repo.UserV1EmailVerifiedType,
		user_repo.UserV1PhoneChangedType,
		user_repo.UserV1PhoneVerifiedType,
		user_repo.UserV1PhoneRemovedType,
		user_repo.UserV1AddressChangedType,
		user_repo.UserDeactivatedType,
		user_repo.UserReactivatedType,
		user_repo.UserLockedType,
		user_repo.UserUnlockedType,
		user_repo.UserV1MFAOTPAddedType,
		user_repo.UserV1MFAOTPVerifiedType,
		user_repo.UserV1MFAOTPRemovedType,
		user_repo.UserV1MFAInitSkippedType,
		user_repo.UserV1PasswordChangedType,
		user_repo.HumanProfileChangedType,
		user_repo.HumanEmailChangedType,
		user_repo.HumanEmailVerifiedType,
		user_repo.HumanAvatarAddedType,
		user_repo.HumanAvatarRemovedType,
		user_repo.HumanPhoneChangedType,
		user_repo.HumanPhoneVerifiedType,
		user_repo.HumanPhoneRemovedType,
		user_repo.HumanAddressChangedType,
		user_repo.HumanMFAOTPAddedType,
		user_repo.HumanMFAOTPVerifiedType,
		user_repo.HumanMFAOTPRemovedType,
		user_repo.HumanOTPSMSAddedType,
		user_repo.HumanOTPSMSRemovedType,
		user_repo.HumanOTPEmailAddedType,
		user_repo.HumanOTPEmailRemovedType,
		user_repo.HumanU2FTokenAddedType,
		user_repo.HumanU2FTokenVerifiedType,
		user_repo.HumanU2FTokenRemovedType,
		user_repo.HumanPasswordlessTokenAddedType,
		user_repo.HumanPasswordlessTokenVerifiedType,
		user_repo.HumanPasswordlessTokenRemovedType,
		user_repo.HumanMFAInitSkippedType,
		user_repo.MachineChangedEventType,
		user_repo.HumanPasswordChangedType,
		user_repo.HumanInitialCodeAddedType,
		user_repo.UserV1InitialCodeAddedType,
		user_repo.UserV1InitializedCheckSucceededType,
		user_repo.HumanInitializedCheckSucceededType,
		user_repo.HumanPasswordlessInitCodeAddedType,
		user_repo.HumanPasswordlessInitCodeRequestedType:
		user, err = u.view.UserByID(event.AggregateID, event.InstanceID)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			logging.WithFields(
				"instance", event.InstanceID,
				"userID", event.AggregateID,
				"eventType", event.Type,
			).Info("user not found in view")
			query, err := usr_view.UserByIDQuery(event.AggregateID, event.InstanceID, 0, user.EventTypes())
			if err != nil {
				return err
			}
			events, err := u.es.FilterEvents(context.Background(), query)
			if err != nil {
				return err
			}
			user = &view_model.UserView{}
			for _, e := range events {
				if err = user.AppendEvent(e); err != nil {
					return err
				}
			}
		}
		err = user.AppendEvent(event)
	case user_repo.UserDomainClaimedType,
		user_repo.UserUserNameChangedType:
		user, err = u.view.UserByID(event.AggregateID, event.InstanceID)
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
			logging.WithFields(
				"instance", event.InstanceID,
				"userID", event.AggregateID,
				"eventType", event.Type,
			).Info("user not found in view")
			query, err := usr_view.UserByIDQuery(event.AggregateID, event.InstanceID, 0, user.EventTypes())
			if err != nil {
				return err
			}
			events, err := u.es.FilterEvents(context.Background(), query)
			if err != nil {
				return err
			}
			user = &view_model.UserView{}
			for _, e := range events {
				if err = user.AppendEvent(e); err != nil {
					return err
				}
			}
		}
		err = user.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillLoginNames(user)
	case user_repo.UserRemovedType:
		return u.view.DeleteUser(event.AggregateID, event.InstanceID, event)
	default:
		return u.view.ProcessedUserSequence(event)
	}
	if err != nil {
		return err
	}
	return u.view.PutUser(user, event)
}

func (u *User) fillLoginNames(user *view_model.UserView) (err error) {
	userLoginMustBeDomain, primaryDomain, domains, err := u.loginNameInformation(context.Background(), user.ResourceOwner, user.InstanceID)
	if err != nil {
		return err
	}
	user.SetLoginNames(userLoginMustBeDomain, domains)
	user.PreferredLoginName = user.GenerateLoginName(primaryDomain, userLoginMustBeDomain)
	return nil
}

func (u *User) ProcessOrg(event *es_models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case org.OrgDomainVerifiedEventType,
		org.OrgDomainRemovedEventType,
		org.DomainPolicyAddedEventType,
		org.DomainPolicyChangedEventType,
		org.DomainPolicyRemovedEventType:
		return u.fillLoginNamesOnOrgUsers(event)
	case org.OrgDomainPrimarySetEventType:
		return u.fillPreferredLoginNamesOnOrgUsers(event)
	case org.OrgRemovedEventType:
		return u.view.UpdateOrgOwnerRemovedUsers(event)
	default:
		return u.view.ProcessedUserSequence(event)
	}
}

func (u *User) ProcessInstance(event *es_models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case instance.InstanceRemovedEventType:
		return u.view.DeleteInstanceUsers(event)
	default:
		return u.view.ProcessedUserSequence(event)
	}
}

func (u *User) fillLoginNamesOnOrgUsers(event *es_models.Event) error {
	userLoginMustBeDomain, _, domains, err := u.loginNameInformation(context.Background(), event.ResourceOwner, event.InstanceID)
	if err != nil {
		return err
	}
	users, err := u.view.UsersByOrgID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.SetLoginNames(userLoginMustBeDomain, domains)
	}
	return u.view.PutUsers(users, event)
}

func (u *User) fillPreferredLoginNamesOnOrgUsers(event *es_models.Event) error {
	userLoginMustBeDomain, primaryDomain, _, err := u.loginNameInformation(context.Background(), event.ResourceOwner, event.InstanceID)
	if err != nil {
		return err
	}
	if !userLoginMustBeDomain {
		return nil
	}
	users, err := u.view.UsersByOrgID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.PreferredLoginName = user.GenerateLoginName(primaryDomain, userLoginMustBeDomain)
	}
	return u.view.PutUsers(users, event)
}

func (u *User) OnError(event *es_models.Event, err error) error {
	logging.WithFields("id", event.AggregateID).WithError(err).Warn("something went wrong in user handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserFailedEvent, u.view.ProcessedUserFailedEvent, u.view.ProcessedUserSequence, u.errorCountUntilSkip)
}

func (u *User) OnSuccess(instanceIDs []string) error {
	return spooler.HandleSuccess(u.view.UpdateUserSpoolerRunTimestamp, instanceIDs)
}

func (u *User) getOrgByID(ctx context.Context, orgID, instanceID string) (*org_model.Org, error) {
	orgQuery, err := view.OrgByIDQuery(orgID, instanceID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, orgQuery)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-3m9vs", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *User) loginNameInformation(ctx context.Context, orgID string, instanceID string) (userLoginMustBeDomain bool, primaryDomain string, domains []*org_model.OrgDomain, err error) {
	org, err := u.getOrgByID(ctx, orgID, instanceID)
	if err != nil {
		return false, "", nil, err
	}
	if org.DomainPolicy != nil {
		return org.DomainPolicy.UserLoginMustBeDomain, org.GetPrimaryDomain().Domain, org.Domains, nil
	}
	policy, err := u.queries.DefaultDomainPolicy(withInstanceID(ctx, org.InstanceID))
	if err != nil {
		return false, "", nil, err
	}
	return policy.UserLoginMustBeDomain, org.GetPrimaryDomain().Domain, org.Domains, nil
}
