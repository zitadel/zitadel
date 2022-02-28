package command

import (
	"context"
	"errors"
	"strings"

	"github.com/caos/zitadel/internal/domain"
	errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
	"golang.org/x/text/language"
)

type OrgSetup struct {
	Name   string
	Domain string
	Human  AddHuman
}

var (
	ErrOrgNameEmpty   = errors.New("Errors.Invalid.Argument")
	ErrOrgDomainEmpty = errors.New("Errors.Invalid.Argument")
)

type Command struct {
	es *eventstore.Eventstore
}

func (command *Command) SetUpOrg(ctx context.Context, o *OrgSetup) (*domain.ObjectDetails, error) {
	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	orgAgg := org.NewAggregate(orgID, orgID)
	userAgg := user.NewAggregate(userID, orgID)

	cmds, err := NewCommander(&orgAgg.Aggregate, addOrgCommand(o.Name)).
		Next(addOrgDomainCommand(o.Domain)).
		Next(verifyOrgDomainCommand(o.Domain)).
		Next(setOrgDomainPrimaryCommand(o.Domain)).
		// TODO: default domain
		Next(WithAggregate(&userAgg.Aggregate), addHumanCommand(command.es, &o.Human)).
		Next(WithAggregate(&orgAgg.Aggregate), addOrgMemberCommand(userID, orgID, domain.RoleOrgOwner)).
		Commands(ctx)
	if err != nil {
		return nil, err
	}

	events, err := command.es.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: orgID,
	}, nil
}

func addOrgCommand(name string) commanderOption {
	// TODO: should we always remove the spaces?
	return func(c *commander) {
		if name = strings.TrimSpace(name); name == "" {
			c.err = ErrOrgNameEmpty
			return
		}
		c.command = func(ctx context.Context, a *eventstore.Aggregate) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewOrgAddedEvent(ctx, a, name)}, nil
		}
	}
}

func addOrgDomainCommand(domain string) commanderOption {
	// TODO: should we always remove the spaces?
	return func(c *commander) {
		if domain = strings.TrimSpace(domain); domain == "" {
			c.err = ErrOrgDomainEmpty
			return
		}
		c.command = func(ctx context.Context, a *eventstore.Aggregate) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainAddedEvent(ctx, a, domain)}, nil
		}
	}
}

func verifyOrgDomainCommand(domain string) commanderOption {
	//TODO: check exists, but when: if domain will be added in this transaction?
	return func(c *commander) {
		if domain = strings.TrimSpace(domain); domain == "" {
			c.err = ErrOrgDomainEmpty
			return
		}
		c.command = func(ctx context.Context, a *eventstore.Aggregate) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainVerifiedEvent(ctx, a, domain)}, nil
		}
	}
}

func setOrgDomainPrimaryCommand(domain string) commanderOption {
	//TODO: check exists, but when if domain will be added in this transaction?
	return func(c *commander) {
		if domain = strings.TrimSpace(domain); domain == "" {
			c.err = ErrOrgDomainEmpty
			return
		}
		c.command = func(ctx context.Context, a *eventstore.Aggregate) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainPrimarySetEvent(ctx, a, domain)}, nil
		}
	}
}

type AddHuman struct {
	// Username is required
	Username string
	// FirstName is required
	FirstName string
	// LastName is required
	LastName string
	// NickName is required
	NickName string
	// DisplayName is required
	DisplayName string
	// Email is required
	Email string
	// PreferredLang is required
	PreferredLang language.Tag
	// Gender is required
	Gender domain.Gender
	//TODO: can it also be verified?
	Phone string
	//Password is optional
	//TODO: should we use the domain object?
	Password *domain.Password
}

func addHumanCommand(es *eventstore.Eventstore, human *AddHuman) commanderOption {
	return func(c *commander) {
		c.command = func(ctx context.Context, a *eventstore.Aggregate) ([]eventstore.Command, error) {
			existing := NewHumanWriteModel(a.ID, a.ResourceOwner)
			err := es.FilterToQueryReducer(ctx, existing)
			if err != nil {
				return nil, err
			}
			if isUserStateExists(existing.UserState) {
				return nil, errs.ThrowAlreadyExists(nil, "COMMA-CxDKf", "Errors.User.AlreadyExists")
			}

			cmd := user.NewHumanAddedEvent(
				ctx,
				a,
				human.Username,
				human.FirstName,
				human.LastName,
				human.NickName,
				human.DisplayName,
				human.PreferredLang,
				human.Gender,
				human.Email,
				true, //TODO: depends on policy
			)
			if phone := strings.TrimSpace(human.Phone); phone != "" {
				cmd.AddPhoneData(phone)
			}
			if human.Password != nil {
				cmd.AddPasswordData(human.Password.SecretCrypto, false) //TOOD: when is it false when true?
			}

			return []eventstore.Command{cmd}, nil
		}
	}
}

func addOrgMemberCommand(userID, orgID string, roles ...string) commanderOption {
	return func(c *commander) {}
}
