package command

import (
	"context"
	"strings"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
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
	Email Email
	// PreferredLang is required
	PreferredLang language.Tag
	// Gender is required
	Gender domain.Gender
	//TODO: can it also be verified?
	Phone Phone
	//Password is optional
	Password string
	//PasswordChangeRequired is used if the `Password`-field is set
	PasswordChangeRequired bool
	Passwordless           bool
	ExternalIDP            bool
	Register               bool
}

type addCommand interface {
	eventstore.Command
	AddPhoneData(phoneNumber string)
	AddPasswordData(secret *crypto.CryptoValue, changeRequired bool)
}

func AddHumanCommand(a *user.Aggregate, human *AddHuman, passwordAlg crypto.HashAlgorithm, phoneAlg, initCodeAlg crypto.EncryptionAlgorithm) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if !human.Email.Valid() {
			return nil, errors.ThrowInvalidArgument(nil, "USER-Ec7dM", "Errors.Invalid.Argument")
		}
		if human.Username = strings.TrimSpace(human.Username); human.Username == "" {
			return nil, errors.ThrowInvalidArgument(nil, "V2-zzad3", "Errors.Invalid.Argument")
		}

		if human.FirstName = strings.TrimSpace(human.FirstName); human.FirstName == "" {
			return nil, errors.ThrowInvalidArgument(nil, "USER-UCej2", "Errors.Invalid.Argument")
		}
		if human.LastName = strings.TrimSpace(human.LastName); human.LastName == "" {
			return nil, errors.ThrowInvalidArgument(nil, "USER-DiAq8", "Errors.Invalid.Argument")
		}
		human.ensureDisplayName()

		//TODO: verify check
		if human.Phone.Number, err = FormatPhoneNumber(human.Phone.Number); err != nil {
			return nil, errors.ThrowInvalidArgument(nil, "USER-tD6ax", "Errors.Invalid.Argument")
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			domainPolicy, err := domainPolicyWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}

			if err = userValidateDomain(ctx, a, human.Username, domainPolicy.UserLoginMustBeDomain, filter); err != nil {
				return nil, err
			}

			var createCmd addCommand
			if human.Register {
				createCmd = user.NewHumanRegisteredEvent(
					ctx,
					&a.Aggregate,
					human.Username,
					human.FirstName,
					human.LastName,
					human.NickName,
					human.DisplayName,
					human.PreferredLang,
					human.Gender,
					human.Email.Address,
					domainPolicy.UserLoginMustBeDomain,
				)
			} else {
				createCmd = user.NewHumanAddedEvent(
					ctx,
					&a.Aggregate,
					human.Username,
					human.FirstName,
					human.LastName,
					human.NickName,
					human.DisplayName,
					human.PreferredLang,
					human.Gender,
					human.Email.Address,
					domainPolicy.UserLoginMustBeDomain,
				)
			}

			if human.Phone.Number != "" {
				createCmd.AddPhoneData(human.Phone.Number)
			}

			if human.Password != "" {
				if err = humanValidatePassword(ctx, filter, human.Password); err != nil {
					return nil, err
				}

				secret, err := crypto.Hash([]byte(human.Password), passwordAlg)
				if err != nil {
					return nil, err
				}
				createCmd.AddPasswordData(secret, human.PasswordChangeRequired)
			}

			if human.shouldAddInitCode() {
				value, expiry, err := newUserInitCode(ctx, filter, initCodeAlg)
				if err != nil {
					return nil, err
				}
				user.NewHumanInitialCodeAddedEvent(ctx, &a.Aggregate, value, expiry)
			}

			cmds := make([]eventstore.Command, 0, 3)
			cmds = append(cmds, createCmd)

			if human.Email.Verified {
				cmds = append(cmds, user.NewHumanEmailVerifiedEvent(ctx, &a.Aggregate))
			} else {
				value, expiry, err := newEmailCode(ctx, filter, phoneAlg)
				if err != nil {
					return nil, err
				}
				cmds = append(cmds, user.NewHumanEmailCodeAddedEvent(ctx, &a.Aggregate, value, expiry))
			}

			if human.Phone.Verified {
				cmds = append(cmds, user.NewHumanPhoneVerifiedEvent(ctx, &a.Aggregate))
			} else if human.Phone.Number != "" {
				value, expiry, err := newPhoneCode(ctx, filter, phoneAlg)
				if err != nil {
					return nil, err
				}
				cmds = append(cmds, user.NewHumanPhoneCodeAddedEvent(ctx, &a.Aggregate, value, expiry))
			}

			return cmds, nil
		}, nil
	}
}

func userValidateDomain(ctx context.Context, a *user.Aggregate, username string, mustBeDomain bool, filter preparation.FilterToQueryReducer) error {
	if mustBeDomain {
		return nil
	}

	usernameSplit := strings.Split(username, "@")
	if len(usernameSplit) != 2 {
		return errors.ThrowInvalidArgument(nil, "COMMAND-Dfd21", "Errors.User.Invalid")
	}

	domainCheck := command.NewOrgDomainVerifiedWriteModel(usernameSplit[1])
	events, err := filter(ctx, domainCheck.Query())
	if err != nil {
		return err
	}
	domainCheck.AppendEvents(events...)
	if err = domainCheck.Reduce(); err != nil {
		return err
	}

	if domainCheck.Verified && domainCheck.ResourceOwner != a.ResourceOwner {
		return errors.ThrowInvalidArgument(nil, "COMMAND-SFd21", "Errors.User.DomainNotAllowedAsUsername")
	}

	return nil
}

func humanValidatePassword(ctx context.Context, filter preparation.FilterToQueryReducer, password string) error {
	passwordComplexity, err := passwordComplexityPolicyWriteModel(ctx, filter)
	if err != nil {
		return err
	}

	return passwordComplexity.Validate(password)
}

func (h *AddHuman) ensureDisplayName() {
	if strings.TrimSpace(h.DisplayName) != "" {
		return
	}
	h.DisplayName = h.FirstName + " " + h.LastName
}

func (h *AddHuman) shouldAddInitCode() bool {
	//user without idp
	return !h.Email.Verified ||
		//user with idp
		!h.ExternalIDP &&
			!h.Passwordless &&
			h.Password != ""
}
