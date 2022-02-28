package command

import (
	"context"
	"errors"
	"strings"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgSetup struct {
	Name   string
	Domain string
}

var (
	ErrOrgNameEmpty   = errors.New("Errors.Invalid.Argument")
	ErrOrgDomainEmpty = errors.New("Errors.Invalid.Argument")
)

func SetUpOrg(ctx context.Context, o *OrgSetup) error {
	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return err
	}

	//TODO: add default domain

	// cmd :=

	cmd := addOrgDomainCommand(
		addOrgCommand(
			NewCommander(&org.NewAggregate(orgID, orgID).Aggregate),
			o.Name,
		),
		o.Domain,
	)

	_ = cmd

	return nil
}

func addOrgCommand(c *commander, name string) *commander {
	//TODO: should we always remove the spaces?
	if strings.TrimSpace(name) == "" {
		return c.Next(nil, WithErr(ErrOrgNameEmpty))
	}
	return c.Next(
		func(ctx context.Context, a *eventstore.Aggregate) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewOrgAddedEvent(ctx, a, name)}, nil
		},
	)
}

func addOrgDomainCommand(c *commander, domain string) *commander {
	//TODO: should we always remove the spaces?
	if strings.TrimSpace(domain) == "" {
		return c.Next(nil, WithErr(ErrOrgDomainEmpty))
	}

	return c.Next(
		func(ctx context.Context, a *eventstore.Aggregate) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainAddedEvent(ctx, a, domain)}, nil
		},
	)
}

/*
AddOrg(name string)
AddDomain(domain string)
SetDomainPrimary(domain string)
VerifyDomain(domain string)
AddMember(userID, roles ...string)

AddUser(firstName, lastName, username string) => validate password, validate username on orgiam policy
*/

// import (
// 	"context"
// 	"errors"
// 	"fmt"

// 	"github.com/caos/zitadel/internal/api/authz"
// 	"github.com/caos/zitadel/internal/domain"
// 	caos_errs "github.com/caos/zitadel/internal/errors"
// 	errs "github.com/caos/zitadel/internal/errors"
// 	"github.com/caos/zitadel/internal/eventstore"
// 	"github.com/caos/zitadel/internal/id"
// 	"github.com/caos/zitadel/internal/repository/org"
// 	"github.com/caos/zitadel/internal/repository/user"
// )

// type SetUpOrg struct {
// 	Name string
// }

// var (
// 	ErrInvalidArg      = errors.New("Errors.Invalid.Argument")
// 	ErrInvalidUsername = fmt.Errorf("%w.User.Username", ErrInvalidArg)
// )

// func SetupOrg(ctx context.Context, o *SetUpOrg) error {
// 	orgID, err := id.SonyFlakeGenerator.Next()
// 	if err != nil {
// 		return errs.ThrowInternal(err, "COMMA-41QOL", "Errors.Internal")
// 	}
// 	userID, err := id.SonyFlakeGenerator.Next()
// 	if err != nil {
// 		return errs.ThrowInternal(err, "COMMA-mSzIb", "Errors.Internal")
// 	}

// 	orgAgg := eventstore.NewAggregate(ctx, orgID, org.AggregateType, org.AggregateVersion)
// 	userAgg := eventstore.NewAggregate(ctx, userID, user.AggregateType, user.AggregateVersion, eventstore.WithResourceOwner(orgID))

// 	defaultDomain := domain.NewIAMDomainName(o.Name, "TODO:")

// 	cmds := []commands{
// 		addOrg(o.Name),
// 		addDomain(&AddDomain{Domain: defaultDomain, Verified: true}),
// 		verifyDomain(defaultDomain),
// 	}
// 	addMember, err := AddMember(userID, domain.RoleOrgOwner)
// 	if err != nil {
// 		return err
// 	}
// 	cmds = append(cmds, addMember)

// 	// cmds = append(cmds, addOrgCommands(ctx, o.Name)

// 	return nil
// }

// type cmds struct {
// 	err  error
// 	cmds []commands
// }

// type commands func(ctx context.Context, agg *eventstore.Aggregate) (cmds []eventstore.Command)

// func (c *commands) chain(...commands) {
// 	if c.err != nil {
// 		return
// 	}

// 	for _, cmmand := range commands {

// 	}
// }

// func (c *commands) Err() error {
// 	return c.err
// }

// func addOrg(name string, domains ...*AddDomain) (commands, error) {
// 	return func(ctx context.Context, agg *eventstore.Aggregate) []eventstore.Command {
// 		cmds := make([]eventstore.Command, 0, len(domains)+3)
// 		cmds = append(cmds, org.NewOrgAddedEvent(ctx, agg, name))
// 		for _, domain := range domains {
// 			cmds = append(cmds, addDomain(domain)(ctx, agg)...)
// 		}
// 		cmds = append(cmds)
// 		return cmds
// 	}, nil
// }

// type AddDomain struct {
// 	Domain   string
// 	Verified bool
// }

// func addDomain(domain *AddDomain) commands {
// 	//TODO: should the event has a unique constraint check for verified domains?
// 	return func(ctx context.Context, agg *eventstore.Aggregate) []eventstore.Command {
// 		cmds := make([]eventstore.Command, 0, 3)

// 		cmds = append(cmds, org.NewDomainAddedEvent(ctx, agg, domain.Domain))
// 		if domain.Verified {
// 			cmds = append(cmds, org.NewDomainVerifiedEvent(ctx, agg, domain.Domain))
// 		}
// 		//TODO: claimed users

// 		return cmds
// 	}
// }

// func verifyDomain(domain string) commands {
// 	return func(ctx context.Context, agg *eventstore.Aggregate) []eventstore.Command {
// 		return []eventstore.Command{org.NewDomainVerifiedEvent(ctx, agg, domain)}
// 	}
// }

// func AddMember(userID string, roles ...string) (commands, error) {
// 	//TODO: rolemapping
// 	if len(domain.CheckForInvalidRoles(roles, domain.OrgRolePrefix, []authz.RoleMapping{})) > 0 &&
// 		len(domain.CheckForInvalidRoles(roles, domain.RoleSelfManagementGlobal, []authz.RoleMapping{})) > 0 {
// 		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMA-4N8es", "Errors.Org.MemberInvalid")
// 	}
// 	//TODO: unique constraint for member and roles?
// 	return func(ctx context.Context, agg *eventstore.Aggregate) []eventstore.Command {
// 		return []eventstore.Command{org.NewMemberAddedEvent(ctx, agg, userID, roles...)}
// 	}, nil
// }
