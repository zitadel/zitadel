package command

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"

	"github.com/caos/zitadel/internal/domain"
	errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
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

	cmd, err := NewCommander(&orgAgg.Aggregate, command.es.Filter,
		addOrgCommand(orgAgg, o.Name),
		addOrgDomainCommand(orgAgg, o.Domain),
		verifyOrgDomainCommand(orgAgg, o.Domain),
		setOrgDomainPrimaryCommand(orgAgg, o.Domain),
		// TODO: default domain
		addHumanCommand(userAgg, &o.Human),
		addOrgMemberCommand(orgAgg, userID, domain.RoleOrgOwner),
	)
	if err != nil {
		return nil, err
	}
	cmds, err := cmd.Commands(ctx)
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

func addOrgCommand(a *org.Aggregate, name string) commanderOption {
	return func(filter FilterToQueryReducer) (createCommands, error) {
		if name = strings.TrimSpace(name); name == "" {
			return nil, ErrOrgNameEmpty
		}
		return func(ctx context.Context) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewOrgAddedEvent(ctx, &a.Aggregate, name)}, nil
		}, nil
	}
}

func addOrgDomainCommand(a *org.Aggregate, domain string) commanderOption {
	return func(filter FilterToQueryReducer) (createCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			//c.err = ErrOrgDomainEmpty
			return nil, ErrOrgDomainEmpty
		}
		return func(ctx context.Context) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainAddedEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func verifyOrgDomainCommand(a *org.Aggregate, domain string) commanderOption {
	return func(filter FilterToQueryReducer) (createCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, ErrOrgDomainEmpty
		}
		return func(ctx context.Context) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainVerifiedEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func setOrgDomainPrimaryCommand(a *org.Aggregate, domain string) commanderOption {
	return func(filter FilterToQueryReducer) (createCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, ErrOrgDomainEmpty
		}
		return func(ctx context.Context) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainPrimarySetEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
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

type FilterToQueryReducer func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error)

type Generator interface {
	Gen() (*crypto.CryptoValue, error)
}

func addHumanCommand(a *user.Aggregate, human *AddHuman) commanderOption {
	return func(filter FilterToQueryReducer) (createCommands, error) {
		return func(ctx context.Context) ([]eventstore.Command, error) {
			existing := NewHumanWriteModel(a.ID, a.ResourceOwner)
			events, err := filter(ctx, existing.Query())
			if err != nil {
				return nil, err
			}
			existing.AppendEvents(events...)
			existing.Reduce()
			if isUserStateExists(existing.UserState) {
				return nil, errs.ThrowAlreadyExists(nil, "COMMA-CxDKf", "Errors.User.AlreadyExists")
			}

			cmd := user.NewHumanAddedEvent(
				ctx,
				&a.Aggregate,
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
		}, nil
	}
}

func changeHumanProfileCommand(a *org.Aggregate, human *AddHuman) commanderOption {
	return func(filter FilterToQueryReducer) (createCommands, error) {
		return func(ctx context.Context) ([]eventstore.Command, error) {
			existing := NewHumanProfileWriteModel(a.ID, a.ResourceOwner)
			events, err := filter(ctx, existing.Query())
			if err != nil {
				return nil, err
			}
			existing.AppendEvents(events...)
			existing.Reduce()
			if !isUserStateExists(existing.UserState) {
				return nil, errs.ThrowAlreadyExists(nil, "COMMA-CxDKf", "Errors.User.AlreadyExists")
			}

			changedEvent, hasChanged, err := existing.NewChangedEvent(ctx, &a.Aggregate, human.FirstName, human.LastName, human.NickName, human.DisplayName, human.PreferredLang, human.Gender)
			if err != nil {
				return nil, err
			}
			if !hasChanged {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.User.Profile.NotChanged")
			}

			return []eventstore.Command{changedEvent}, nil
		}, nil
	}
}

func addOrgMemberCommand(a *org.Aggregate, userID string, roles ...string) commanderOption {
	return func(filter FilterToQueryReducer) (createCommands, error) { return nil, nil }
}
