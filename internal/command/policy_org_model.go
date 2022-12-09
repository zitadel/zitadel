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

	PrimaryDomain   string
	VerifiedDomains []string
	Users           []*domainPolicyUsers
}

type domainPolicyUsers struct {
	id       string
	username string
}

func NewDomainPolicyUsernamesWriteModel(orgID string) *DomainPolicyUsernamesWriteModel {
	return &DomainPolicyUsernamesWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: orgID,
		},
		Users: make([]*domainPolicyUsers, 0),
	}
}

func (wm *DomainPolicyUsernamesWriteModel) AppendEvents(events ...eventstore.Event) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *DomainPolicyUsernamesWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.DomainVerifiedEvent:
			wm.VerifiedDomains = append(wm.VerifiedDomains, e.Domain)
		case *org.DomainRemovedEvent:
			wm.removeDomain(e.Domain)
		case *org.DomainPrimarySetEvent:
			wm.PrimaryDomain = e.Domain
		case *user.HumanAddedEvent:
			wm.Users = append(wm.Users, &domainPolicyUsers{id: e.Aggregate().ID, username: e.UserName})
		case *user.HumanRegisteredEvent:
			wm.Users = append(wm.Users, &domainPolicyUsers{id: e.Aggregate().ID, username: e.UserName})
		case *user.MachineAddedEvent:
			wm.Users = append(wm.Users, &domainPolicyUsers{id: e.Aggregate().ID, username: e.UserName})
		case *user.UsernameChangedEvent:
			for _, user := range wm.Users {
				if user.id == e.Aggregate().ID {
					user.username = e.UserName
					break
				}
			}
		case *user.DomainClaimedEvent:
			for _, user := range wm.Users {
				if user.id == e.Aggregate().ID {
					user.username = e.UserName
					break
				}
			}
		case *user.UserRemovedEvent:
			wm.removeUser(e.Aggregate().ID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *DomainPolicyUsernamesWriteModel) removeDomain(domain string) {
	for i, verifiedDomain := range wm.VerifiedDomains {
		if verifiedDomain == domain {
			wm.VerifiedDomains[i] = wm.VerifiedDomains[len(wm.VerifiedDomains)-1]
			wm.VerifiedDomains[len(wm.VerifiedDomains)-1] = ""
			wm.VerifiedDomains = wm.VerifiedDomains[:len(wm.VerifiedDomains)-1]
			return
		}
	}
}

func (wm *DomainPolicyUsernamesWriteModel) removeUser(userID string) {
	for i, user := range wm.Users {
		if user.id == userID {
			wm.Users[i] = wm.Users[len(wm.Users)-1]
			wm.Users[len(wm.Users)-1] = nil
			wm.Users = wm.Users[:len(wm.Users)-1]
			return
		}
	}
}

func (wm *DomainPolicyUsernamesWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType, user.AggregateType).
		EventTypes(
			org.OrgDomainVerifiedEventType,
			org.OrgDomainRemovedEventType,
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
	for _, changeUser := range wm.Users {
		events = append(events, user.NewUsernameChangedEvent(ctx,
			&user.NewAggregate(changeUser.id, wm.ResourceOwner).Aggregate,
			changeUser.username,
			wm.newUsername(changeUser.username, userLoginMustBeDomain),
			userLoginMustBeDomain,
			user.UsernameChangedEventWithPolicyChange()),
		)
	}
	return events
}

func (wm *DomainPolicyUsernamesWriteModel) newUsername(username string, userLoginMustBeDomain bool) string {
	if !userLoginMustBeDomain {
		// if the UserLoginMustBeDomain will be false, then it's currently true
		// which means the usernames must be suffixed to ensure their uniqueness
		// and the preferred login name remains the same
		return username + "@" + wm.PrimaryDomain
	}
	// the UserLoginMustBeDomain is currently false
	// which means the usernames might already be suffixed by a verified domain
	// so let's remove a potential duplicate suffix
	for _, verifiedDomain := range wm.VerifiedDomains {
		if index := strings.LastIndex(username, "@"+verifiedDomain); index > 0 {
			return username[:index]
		}
	}
	return username
}
