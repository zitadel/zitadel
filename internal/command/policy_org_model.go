package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type PolicyDomainWriteModel struct {
	eventstore.WriteModel

	UserLoginMustBeDomain                  bool
	ValidateOrgDomains                     bool
	SMTPSenderAddressMatchesInstanceDomain bool
	State                                  domain.PolicyState
}

func (wm *PolicyDomainWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.DomainPolicyAddedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
			wm.ValidateOrgDomains = e.ValidateOrgDomains
			wm.SMTPSenderAddressMatchesInstanceDomain = e.SMTPSenderAddressMatchesInstanceDomain
			wm.State = domain.PolicyStateActive
		case *policy.DomainPolicyChangedEvent:
			if e.UserLoginMustBeDomain != nil {
				wm.UserLoginMustBeDomain = *e.UserLoginMustBeDomain
			}
			if e.ValidateOrgDomains != nil {
				wm.ValidateOrgDomains = *e.ValidateOrgDomains
			}
			if e.SMTPSenderAddressMatchesInstanceDomain != nil {
				wm.SMTPSenderAddressMatchesInstanceDomain = *e.SMTPSenderAddressMatchesInstanceDomain
			}
		case *policy.DomainPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

type DomainPolicyUsernamesWriteModel struct {
	eventstore.WriteModel

	Domain string
	Users  map[string]string
}

func NewDomainPolicyUsernamesWriteModel(orgID string) *DomainPolicyUsernamesWriteModel {
	return &DomainPolicyUsernamesWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: orgID,
		},
		Users: make(map[string]string),
	}
}

func (wm *DomainPolicyUsernamesWriteModel) AppendEvents(events ...eventstore.Event) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *DomainPolicyUsernamesWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.DomainPrimarySetEvent:
			wm.Domain = e.Domain
		case *user.HumanAddedEvent:
			wm.Users[e.Aggregate().ID] = e.UserName
		case *user.HumanRegisteredEvent:
			wm.Users[e.Aggregate().ID] = e.UserName
		case *user.MachineAddedEvent:
			wm.Users[e.Aggregate().ID] = e.UserName
		case *user.UsernameChangedEvent:
			wm.Users[e.Aggregate().ID] = e.UserName
		case *user.DomainClaimedEvent:
			wm.Users[e.Aggregate().ID] = e.UserName
		case *user.UserRemovedEvent:
			delete(wm.Users, e.Aggregate().ID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *DomainPolicyUsernamesWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType, user.AggregateType).
		EventTypes(
			org.OrgDomainPrimarySetEventType,
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.MachineAddedEventType,
			user.UserUserNameChangedType,
			user.UserDomainClaimedType,
			user.UserRemovedType,
		).
		Builder()
}

func (wm *DomainPolicyUsernamesWriteModel) NewUsernameChangedEvents(ctx context.Context, userLoginMustBeDomain bool) []eventstore.Command {
	events := make([]eventstore.Command, 0, len(wm.Users))
	for id, name := range wm.Users {
		var newName string
		if userLoginMustBeDomain {
			// if the UserLoginMustBeDomain will be true, then it's currently false
			// which means the usernames might already be suffixed by the domain
			// so let's remove a potential duplicate suffix
			newName = strings.TrimSuffix(name, "@"+wm.Domain)
		} else {
			// the UserLoginMustBeDomain is currently true
			// which means the usernames must be suffixed to ensure their uniqueness
			// and the preferred login name remains the same
			newName = name + "@" + wm.Domain
		}
		events = append(events, user.NewUsernameChangedEvent(ctx,
			&user.NewAggregate(id, wm.ResourceOwner).Aggregate,
			name,
			newName,
			userLoginMustBeDomain,
			user.UsernameChangedEventWithPolicyChange()),
		)
	}
	return events
}
