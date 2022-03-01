package user

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
	"golang.org/x/text/language"
)

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

func AddHumanCommand(a *user.Aggregate, human *AddHuman) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existing := command.NewHumanWriteModel(a.ID, a.ResourceOwner)
			events, err := filter(ctx, existing.Query())
			if err != nil {
				return nil, err
			}
			existing.AppendEvents(events...)
			existing.Reduce()
			if isUserStateExists(existing.UserState) {
				return nil, errors.ThrowAlreadyExists(nil, "COMMA-CxDKf", "Errors.User.AlreadyExists")
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

func ChangeHumanProfileCommand(a *user.Aggregate, human *AddHuman) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existing := command.NewHumanProfileWriteModel(a.ID, a.ResourceOwner)
			events, err := filter(ctx, existing.Query())
			if err != nil {
				return nil, err
			}
			existing.AppendEvents(events...)
			existing.Reduce()
			if !isUserStateExists(existing.UserState) {
				return nil, errors.ThrowAlreadyExists(nil, "COMMA-CxDKf", "Errors.User.AlreadyExists")
			}

			changedEvent, hasChanged, err := existing.NewChangedEvent(ctx, &a.Aggregate, human.FirstName, human.LastName, human.NickName, human.DisplayName, human.PreferredLang, human.Gender)
			if err != nil {
				return nil, err
			}
			if !hasChanged {
				return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-2M0fs", "Errors.User.Profile.NotChanged")
			}

			return []eventstore.Command{changedEvent}, nil
		}, nil
	}
}
