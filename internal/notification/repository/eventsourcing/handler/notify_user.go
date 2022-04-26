package handler

import (
	"context"

	"github.com/caos/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	es_sdk "github.com/zitadel/zitadel/internal/eventstore/v1/sdk"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	org_model "github.com/zitadel/zitadel/internal/org/model"
	org_es_model "github.com/zitadel/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/zitadel/zitadel/internal/org/repository/view"
	query2 "github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

const (
	userTable = "notification.notify_users"
)

type NotifyUser struct {
	handler
	subscription *v1.Subscription
	queries      *query2.Queries
}

func newNotifyUser(
	handler handler,
	queries *query2.Queries,
) *NotifyUser {
	h := &NotifyUser{
		handler: handler,
		queries: queries,
	}

	h.subscribe()

	return h
}

func (k *NotifyUser) subscribe() {
	k.subscription = k.es.Subscribe(k.AggregateTypes()...)
	go func() {
		for event := range k.subscription.Events {
			query.ReduceEvent(k, event)
		}
	}()
}

func (p *NotifyUser) ViewModel() string {
	return userTable
}

func (p *NotifyUser) Subscription() *v1.Subscription {
	return p.subscription
}

func (_ *NotifyUser) AggregateTypes() []es_models.AggregateType {
	return []es_models.AggregateType{user.AggregateType, org.AggregateType}
}

func (p *NotifyUser) CurrentSequence(instanceID string) (uint64, error) {
	sequence, err := p.view.GetLatestNotifyUserSequence(instanceID)
	if err != nil {
		return 0, err
	}
	return sequence.CurrentSequence, nil
}

func (p *NotifyUser) EventQuery() (*es_models.SearchQuery, error) {
	sequences, err := p.view.GetLatestNotifyUserSequences()
	if err != nil {
		return nil, err
	}
	query := es_models.NewSearchQuery()
	instances := make([]string, 0)
	for _, sequence := range sequences {
		for _, instance := range instances {
			if sequence.InstanceID == instance {
				break
			}
		}
		instances = append(instances, sequence.InstanceID)
		query.AddQuery().
			AggregateTypeFilter(p.AggregateTypes()...).
			LatestSequenceFilter(sequence.CurrentSequence).
			InstanceIDFilter(sequence.InstanceID)
	}
	return query.AddQuery().
		AggregateTypeFilter(p.AggregateTypes()...).
		LatestSequenceFilter(0).
		ExcludedInstanceIDsFilter(instances...).
		SearchQuery(), nil
}

func (u *NotifyUser) Reduce(event *es_models.Event) (err error) {
	switch event.AggregateType {
	case user.AggregateType:
		return u.ProcessUser(event)
	case org.AggregateType:
		return u.ProcessOrg(event)
	default:
		return nil
	}
}

func (u *NotifyUser) ProcessUser(event *es_models.Event) (err error) {
	notifyUser := new(view_model.NotifyUser)
	switch eventstore.EventType(event.Type) {
	case user.UserV1AddedType,
		user.UserV1RegisteredType,
		user.HumanRegisteredType,
		user.HumanAddedType,
		user.MachineAddedEventType:
		err := notifyUser.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillLoginNames(notifyUser)
	case user.UserV1ProfileChangedType,
		user.UserV1EmailChangedType,
		user.UserV1EmailVerifiedType,
		user.UserV1PhoneChangedType,
		user.UserV1PhoneVerifiedType,
		user.UserV1PhoneRemovedType,
		user.HumanProfileChangedType,
		user.HumanEmailChangedType,
		user.HumanEmailVerifiedType,
		user.HumanPhoneChangedType,
		user.HumanPhoneVerifiedType,
		user.HumanPhoneRemovedType,
		user.MachineChangedEventType:
		notifyUser, err = u.view.NotifyUserByID(event.AggregateID, event.InstanceID)
		if err != nil {
			return err
		}
		err = notifyUser.AppendEvent(event)
	case user.UserDomainClaimedType,
		user.UserUserNameChangedType:
		notifyUser, err = u.view.NotifyUserByID(event.AggregateID, event.InstanceID)
		if err != nil {
			return err
		}
		err = notifyUser.AppendEvent(event)
		if err != nil {
			return err
		}
		err = u.fillLoginNames(notifyUser)
	case user.UserRemovedType:
		return u.view.DeleteNotifyUser(event.AggregateID, event.InstanceID, event)
	default:
		return u.view.ProcessedNotifyUserSequence(event)
	}
	if err != nil {
		return err
	}
	return u.view.PutNotifyUser(notifyUser, event)
}

func (u *NotifyUser) ProcessOrg(event *es_models.Event) (err error) {
	switch eventstore.EventType(event.Type) {
	case org.OrgDomainVerifiedEventType,
		org.OrgDomainRemovedEventType,
		org.DomainPolicyAddedEventType,
		org.DomainPolicyChangedEventType,
		org.DomainPolicyRemovedEventType:
		return u.fillLoginNamesOnOrgUsers(event)
	case org.OrgDomainPrimarySetEventType:
		return u.fillPreferredLoginNamesOnOrgUsers(event)
	default:
		return u.view.ProcessedNotifyUserSequence(event)
	}
}

func (u *NotifyUser) fillLoginNamesOnOrgUsers(event *es_models.Event) error {
	userLoginMustBeDomain, _, domains, err := u.loginNameInformation(context.Background(), event.ResourceOwner)
	if err != nil {
		return err
	}
	users, err := u.view.NotifyUsersByOrgID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.SetLoginNames(userLoginMustBeDomain, domains)
		err := u.view.PutNotifyUser(user, event)
		if err != nil {
			return err
		}
	}
	return u.view.ProcessedNotifyUserSequence(event)
}

func (u *NotifyUser) fillPreferredLoginNamesOnOrgUsers(event *es_models.Event) error {
	userLoginMustBeDomain, primaryDomain, _, err := u.loginNameInformation(context.Background(), event.ResourceOwner)
	if err != nil {
		return err
	}
	if !userLoginMustBeDomain {
		return nil
	}
	users, err := u.view.NotifyUsersByOrgID(event.AggregateID, event.InstanceID)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.PreferredLoginName = user.GenerateLoginName(primaryDomain, userLoginMustBeDomain)
		err := u.view.PutNotifyUser(user, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *NotifyUser) fillLoginNames(user *view_model.NotifyUser) (err error) {
	userLoginMustBeDomain, primaryDomain, domains, err := u.loginNameInformation(context.Background(), user.ResourceOwner)
	if err != nil {
		return err
	}
	user.SetLoginNames(userLoginMustBeDomain, domains)
	user.PreferredLoginName = user.GenerateLoginName(primaryDomain, userLoginMustBeDomain)
	return nil
}

func (p *NotifyUser) OnError(event *es_models.Event, err error) error {
	logging.LogWithFields("SPOOL-9spwf", "id", event.AggregateID).WithError(err).Warn("something went wrong in notify user handler")
	return spooler.HandleError(event, err, p.view.GetLatestNotifyUserFailedEvent, p.view.ProcessedNotifyUserFailedEvent, p.view.ProcessedNotifyUserSequence, p.errorCountUntilSkip)
}

func (u *NotifyUser) OnSuccess() error {
	return spooler.HandleSuccess(u.view.UpdateNotifyUserSpoolerRunTimestamp)
}

func (u *NotifyUser) getOrgByID(ctx context.Context, orgID string) (*org_model.Org, error) {
	query, err := org_view.OrgByIDQuery(orgID, 0)
	if err != nil {
		return nil, err
	}

	esOrg := &org_es_model.Org{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID: orgID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore().FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !caos_errs.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-kVLb2", "Errors.Org.NotFound")
	}

	return org_es_model.OrgToModel(esOrg), nil
}

func (u *NotifyUser) loginNameInformation(ctx context.Context, orgID string) (userLoginMustBeDomain bool, primaryDomain string, domains []*org_model.OrgDomain, err error) {
	org, err := u.getOrgByID(ctx, orgID)
	if err != nil {
		return false, "", nil, err
	}
	if org.DomainPolicy == nil {
		policy, err := u.queries.DefaultDomainPolicy(authz.WithInstanceID(ctx, org.InstanceID))
		if err != nil {
			return false, "", nil, err
		}
		userLoginMustBeDomain = policy.UserLoginMustBeDomain
	}
	return userLoginMustBeDomain, org.GetPrimaryDomain().Domain, org.Domains, nil
}
