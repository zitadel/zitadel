package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type Step19 struct{}

func (s *Step19) Step() domain.Step {
	return domain.Step19
}

func (s *Step19) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep19(ctx, s)
}

func (c *Commands) SetupStep19(ctx context.Context, step *Step19) error {
	fn := func(iam *IAMWriteModel) ([]eventstore.Command, error) {
		events := make([]eventstore.Command, 0)
		orgs := newOrgsWithUsernameNotDomain()
		if err := c.eventstore.FilterToQueryReducer(ctx, orgs); err != nil {
			return nil, err
		}
		for orgID, usernameCheck := range orgs.orgs {
			if !usernameCheck {
				continue
			}
			users := newDomainClaimedUsernames(orgID)
			if err := c.eventstore.FilterToQueryReducer(ctx, users); err != nil {
				return nil, err
			}
			for userID, username := range users.users {
				split := strings.Split(username, "@")
				if len(split) != 2 {
					continue
				}
				domainVerified := NewOrgDomainVerifiedWriteModel(split[1])
				if err := c.eventstore.FilterToQueryReducer(ctx, domainVerified); err != nil {
					return nil, err
				}
				if domainVerified.Verified && domainVerified.ResourceOwner != orgID {
					id, err := c.idGenerator.Next()
					if err != nil {
						return nil, err
					}
					events = append(events, user.NewDomainClaimedEvent(
						ctx,
						&user.NewAggregate(userID, orgID).Aggregate,
						fmt.Sprintf("%s@temporary.%s", id, c.iamDomain),
						username,
						false))
				}
			}
		}

		if length := len(events); length > 0 {
			logging.Log("SETUP-dFG2t").WithField("count", length).Info("domain claimed events created")
		}
		return events, nil
	}
	return c.setup(ctx, step, fn)
}

func newOrgsWithUsernameNotDomain() *orgsWithUsernameNotDomain {
	return &orgsWithUsernameNotDomain{
		orgEvents: make(map[string][]eventstore.Event),
		orgs:      make(map[string]bool),
	}
}

type orgsWithUsernameNotDomain struct {
	eventstore.WriteModel

	orgEvents map[string][]eventstore.Event
	orgs      map[string]bool
}

func (wm *orgsWithUsernameNotDomain) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.OrgAddedEvent:
			wm.orgEvents[e.Aggregate().ID] = append(wm.orgEvents[e.Aggregate().ID], e)
		case *org.OrgRemovedEvent:
			delete(wm.orgEvents, e.Aggregate().ID)
		case *org.OrgIAMPolicyAddedEvent:
			wm.orgEvents[e.Aggregate().ID] = append(wm.orgEvents[e.Aggregate().ID], e)
		case *org.OrgIAMPolicyChangedEvent:
			if e.UserLoginMustBeDomain == nil {
				continue
			}
			wm.orgEvents[e.Aggregate().ID] = append(wm.orgEvents[e.Aggregate().ID], e)
		case *org.OrgIAMPolicyRemovedEvent:
			delete(wm.orgEvents, e.Aggregate().ID)
		}
	}
}

func (wm *orgsWithUsernameNotDomain) Reduce() error {
	for _, events := range wm.orgEvents {
		for _, event := range events {
			switch e := event.(type) {
			case *org.OrgIAMPolicyAddedEvent:
				if !e.UserLoginMustBeDomain {
					wm.orgs[e.Aggregate().ID] = true
				}
			case *org.OrgIAMPolicyChangedEvent:
				if !*e.UserLoginMustBeDomain {
					wm.orgs[e.Aggregate().ID] = true
				}
				delete(wm.orgs, e.Aggregate().ID)
			}
		}
	}
	return nil
}

func (wm *orgsWithUsernameNotDomain) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgRemovedEventType,
			org.OrgIAMPolicyAddedEventType,
			org.OrgIAMPolicyChangedEventType,
			org.OrgIAMPolicyRemovedEventType).
		Builder()
}

func newDomainClaimedUsernames(orgID string) *domainClaimedUsernames {
	return &domainClaimedUsernames{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: orgID,
		},
		userEvents: make(map[string][]eventstore.Event),
		users:      make(map[string]string),
	}
}

type domainClaimedUsernames struct {
	eventstore.WriteModel

	userEvents map[string][]eventstore.Event
	users      map[string]string
}

func (wm *domainClaimedUsernames) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			if !strings.Contains(e.UserName, "@") {
				continue
			}
			wm.userEvents[e.Aggregate().ID] = append(wm.userEvents[e.Aggregate().ID], e)
		case *user.HumanRegisteredEvent:
			if !strings.Contains(e.UserName, "@") {
				continue
			}
			wm.userEvents[e.Aggregate().ID] = append(wm.userEvents[e.Aggregate().ID], e)
		case *user.UsernameChangedEvent:
			if !strings.Contains(e.UserName, "@") {
				delete(wm.userEvents, e.Aggregate().ID)
				continue
			}
			wm.userEvents[e.Aggregate().ID] = append(wm.userEvents[e.Aggregate().ID], e)
		case *user.DomainClaimedEvent:
			delete(wm.userEvents, e.Aggregate().ID)
		case *user.UserRemovedEvent:
			delete(wm.userEvents, e.Aggregate().ID)
		}
	}
}

func (wm *domainClaimedUsernames) Reduce() error {
	for _, events := range wm.userEvents {
		for _, event := range events {
			switch e := event.(type) {
			case *user.HumanAddedEvent:
				wm.users[e.Aggregate().ID] = e.UserName
			case *user.HumanRegisteredEvent:
				wm.users[e.Aggregate().ID] = e.UserName
			case *user.UsernameChangedEvent:
				wm.users[e.Aggregate().ID] = e.UserName
			}
		}
	}
	return nil
}

func (wm *domainClaimedUsernames) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.UserUserNameChangedType,
			user.UserDomainClaimedType,
			user.UserRemovedType).
		Builder()
}
