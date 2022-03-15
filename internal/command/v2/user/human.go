package user

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
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
	Password string
	//PasswordChangeRequired is used if the `Password`-field is set
	PasswordChangeRequired bool
}

func AddHumanCommand(a *user.Aggregate, human *AddHuman, passwordAlg crypto.HashAlgorithm) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if !domain.EmailRegex.MatchString(human.Email) {
			return nil, errors.ThrowInvalidArgument(nil, "USER-Ec7dM", "Errors.Invalid.Argument")
		}
		if human.FirstName = strings.TrimSpace(human.FirstName); human.FirstName == "" {
			return nil, errors.ThrowInvalidArgument(nil, "USER-UCej2", "Errors.Invalid.Argument")
		}
		if human.LastName = strings.TrimSpace(human.LastName); human.LastName == "" {
			return nil, errors.ThrowInvalidArgument(nil, "USER-DiAq8", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if exists, err := existsUser(ctx, filter, a.ID, a.ResourceOwner); exists || err != nil {
				return nil, errors.ThrowAlreadyExists(err, "COMMA-CxDKf", "Errors.User.AlreadyExists")
			}

			orgIAMPolicy, err := orgIAMPolicyWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}

			passwordComplexity, err := passwordComplexityPolicyWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}

			if err = passwordComplexity.Validate(human.Password); err != nil {
				return nil, err
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
				orgIAMPolicy.UserLoginMustBeDomain,
			)
			if phone := strings.TrimSpace(human.Phone); phone != "" {
				cmd.AddPhoneData(phone)
			}
			if human.Password != "" {
				secret, err := crypto.Hash([]byte(human.Password), passwordAlg)
				if err != nil {
					return nil, err
				}
				cmd.AddPasswordData(secret, human.PasswordChangeRequired)
			}

			return []eventstore.Command{cmd}, nil
		}, nil
	}
}

func existsUser(ctx context.Context, filter preparation.FilterToQueryReducer, id, resourceOwner string) (bool, error) {
	existing := command.NewHumanWriteModel(id, resourceOwner)
	events, err := filter(ctx, existing.Query())
	if err != nil {
		return false, err
	}
	existing.AppendEvents(events...)
	existing.Reduce()
	return isUserStateExists(existing.UserState), nil
}
