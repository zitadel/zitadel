package handler

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	org_es_model "github.com/zitadel/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/zitadel/zitadel/internal/org/repository/view"
	query2 "github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	userTable = "auth.users2"
)

type User struct {
	view    *auth_view.View
	queries *query2.Queries
	es      handler.EventStore
}

var _ handler.Projection = (*User)(nil)

func newUser(
	ctx context.Context,
	config handler.Config,
	view *auth_view.View,
	queries *query2.Queries,
) *handler.Handler {
	return handler.NewHandler(
		ctx,
		&config,
		&User{
			view:    view,
			queries: queries,
			es:      config.Eventstore,
		},
	)
}

func (*User) Name() string {
	return userTable
}
func (u *User) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user_repo.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user_repo.HumanOTPSMSAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanOTPSMSRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanOTPEmailAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanOTPEmailRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.MachineAddedEventType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1AddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1RegisteredType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanRegisteredType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1ProfileChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1EmailChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1EmailVerifiedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1PhoneChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1PhoneVerifiedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1PhoneRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1AddressChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserDeactivatedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserReactivatedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserLockedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserUnlockedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1MFAOTPAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1MFAOTPVerifiedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1MFAOTPRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1MFAInitSkippedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1PasswordChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanProfileChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanEmailChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanEmailVerifiedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanAvatarAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanAvatarRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPhoneChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPhoneVerifiedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPhoneRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanAddressChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanMFAOTPAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanMFAOTPVerifiedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanMFAOTPRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanU2FTokenAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanU2FTokenVerifiedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanU2FTokenRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPasswordlessTokenAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPasswordlessTokenVerifiedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPasswordlessTokenRemovedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanMFAInitSkippedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.MachineChangedEventType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPasswordChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanInitialCodeAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1InitialCodeAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserV1InitializedCheckSucceededType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanInitializedCheckSucceededType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPasswordlessInitCodeAddedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.HumanPasswordlessInitCodeRequestedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserDomainClaimedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserUserNameChangedType,
					Reduce: u.ProcessUser,
				},
				{
					Event:  user_repo.UserRemovedType,
					Reduce: u.ProcessUser,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgDomainVerifiedEventType,
					Reduce: u.ProcessOrg,
				},
				{
					Event:  org.OrgDomainRemovedEventType,
					Reduce: u.ProcessOrg,
				},
				{
					Event:  org.DomainPolicyAddedEventType,
					Reduce: u.ProcessOrg,
				},
				{
					Event:  org.DomainPolicyChangedEventType,
					Reduce: u.ProcessOrg,
				},
				{
					Event:  org.DomainPolicyRemovedEventType,
					Reduce: u.ProcessOrg,
				},
				{
					Event:  org.OrgDomainPrimarySetEventType,
					Reduce: u.ProcessOrg,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: u.ProcessOrg,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: u.ProcessInstance,
				},
			},
		},
	}
}

//nolint:gocognit
func (u *User) ProcessUser(event eventstore.Event) (_ *handler.Statement, err error) {
	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		user := new(view_model.UserView)
		switch event.Type() {
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
			user, err = u.view.UserByID(event.Aggregate().ID, event.Aggregate().InstanceID)
			if err != nil {
				if !zerrors.IsNotFound(err) {
					return err
				}
				user, err = u.userFromEventstore(event.Aggregate(), user.EventTypes())
				if err != nil {
					return err
				}
			}
			err = user.AppendEvent(event)
		case user_repo.UserDomainClaimedType,
			user_repo.UserUserNameChangedType:
			user, err = u.view.UserByID(event.Aggregate().ID, event.Aggregate().InstanceID)
			if err != nil {
				if !zerrors.IsNotFound(err) {
					return err
				}
				user, err = u.userFromEventstore(event.Aggregate(), user.EventTypes())
				if err != nil {
					return err
				}
			}
			err = user.AppendEvent(event)
			if err != nil {
				return err
			}
			err = u.fillLoginNames(user)
		case user_repo.UserRemovedType:
			return u.view.DeleteUser(event.Aggregate().ID, event.Aggregate().InstanceID, event)
		default:
			return nil
		}
		if err != nil {
			return err
		}
		return u.view.PutUser(user, event)
	}), nil
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

func (u *User) ProcessOrg(event eventstore.Event) (_ *handler.Statement, err error) {
	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		switch event.Type() {
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
			return nil
		}
	}), nil
}

func (u *User) ProcessInstance(event eventstore.Event) (_ *handler.Statement, err error) {
	switch event.Type() {
	case instance.InstanceRemovedEventType:
		return handler.NewStatement(event,
			func(ex handler.Executer, projectionName string) error {
				return u.view.DeleteInstanceUsers(event)
			},
		), nil
	default:
		return handler.NewNoOpStatement(event), nil
	}
}

func (u *User) fillLoginNamesOnOrgUsers(event eventstore.Event) error {
	userLoginMustBeDomain, _, domains, err := u.loginNameInformation(context.Background(), event.Aggregate().ResourceOwner, event.Aggregate().InstanceID)
	if err != nil {
		return err
	}
	users, err := u.view.UsersByOrgID(event.Aggregate().ID, event.Aggregate().InstanceID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.SetLoginNames(userLoginMustBeDomain, domains)
	}
	return u.view.PutUsers(users, event)
}

func (u *User) fillPreferredLoginNamesOnOrgUsers(event eventstore.Event) error {
	userLoginMustBeDomain, primaryDomain, _, err := u.loginNameInformation(context.Background(), event.Aggregate().ResourceOwner, event.Aggregate().InstanceID)
	if err != nil {
		return err
	}
	if !userLoginMustBeDomain {
		return nil
	}
	users, err := u.view.UsersByOrgID(event.Aggregate().ID, event.Aggregate().InstanceID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.PreferredLoginName = user.GenerateLoginName(primaryDomain, userLoginMustBeDomain)
	}
	return u.view.PutUsers(users, event)
}

func (u *User) getOrgByID(ctx context.Context, orgID, instanceID string) (*org_model.Org, error) {
	query, err := org_view.OrgByIDQuery(orgID, instanceID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	events, err := u.es.Filter(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = esOrg.AppendEvents(events...); err != nil {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-3m9vs", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *User) loginNameInformation(ctx context.Context, orgID string, instanceID string) (userLoginMustBeDomain bool, primaryDomain string, domains []*org_model.OrgDomain, err error) {
	org, err := u.getOrgByID(ctx, orgID, instanceID)
	if err != nil {
		return false, "", nil, err
	}
	primaryDomain, err = org.GetPrimaryDomain()
	if err != nil {
		return false, "", nil, err
	}
	if org.DomainPolicy != nil {
		return org.DomainPolicy.UserLoginMustBeDomain, primaryDomain, org.Domains, nil
	}
	policy, err := u.queries.DefaultDomainPolicy(authz.WithInstanceID(ctx, org.InstanceID))
	if err != nil {
		return false, "", nil, err
	}
	return policy.UserLoginMustBeDomain, primaryDomain, org.Domains, nil
}

func (u *User) userFromEventstore(agg *eventstore.Aggregate, eventTypes []eventstore.EventType) (*view_model.UserView, error) {
	query, err := usr_view.UserByIDQuery(agg.ID, agg.InstanceID, time.Time{}, eventTypes)
	if err != nil {
		return nil, err
	}
	events, err := u.es.Filter(context.Background(), query)
	if err != nil {
		return nil, err
	}
	user := &view_model.UserView{}
	for _, e := range events {
		if err = user.AppendEvent(e); err != nil {
			return nil, err
		}
	}
	return user, nil
}
