package command

import (
	"context"
	"strings"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func (c *Commands) getHuman(ctx context.Context, userID, resourceowner string) (*domain.Human, error) {
	human, err := c.getHumanWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(human.UserState) {
		return nil, errors.ThrowNotFound(nil, "COMMAND-M9dsd", "Errors.User.NotFound")
	}
	return writeModelToHuman(human), nil
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
	Email Email
	// PreferredLanguage is required
	PreferredLanguage language.Tag
	// Gender is required
	Gender domain.Gender
	//Phone represents an international phone number
	Phone Phone
	//Password is optional
	Password string
	//BcryptedPassword is optional
	BcryptedPassword string
	//PasswordChangeRequired is used if the `Password`-field is set
	PasswordChangeRequired bool
	Passwordless           bool
	ExternalIDP            bool
	Register               bool
}

func (c *Commands) AddHumanWithID(ctx context.Context, resourceOwner string, userID string, human *AddHuman) (*domain.HumanDetails, error) {
	existingHuman, err := c.getHumanWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if isUserStateExists(existingHuman.UserState) {
		return nil, errors.ThrowPreconditionFailed(nil, "COMMAND-k2unb", "Errors.User.AlreadyExisting")
	}

	return c.addHumanWithID(ctx, resourceOwner, userID, human)
}

func (c *Commands) addHumanWithID(ctx context.Context, resourceOwner string, userID string, human *AddHuman) (*domain.HumanDetails, error) {
	agg := user.NewAggregate(userID, resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, AddHumanCommand(agg, human, c.userPasswordAlg, c.userEncryption))
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}

	return &domain.HumanDetails{
		ID: userID,
		ObjectDetails: domain.ObjectDetails{
			Sequence:      events[len(events)-1].Sequence(),
			EventDate:     events[len(events)-1].CreationDate(),
			ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
		},
	}, nil
}

func (c *Commands) AddHuman(ctx context.Context, resourceOwner string, human *AddHuman) (*domain.HumanDetails, error) {
	if resourceOwner == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}
	userID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	return c.addHumanWithID(ctx, resourceOwner, userID, human)
}

type humanCreationCommand interface {
	eventstore.Command
	AddPhoneData(phoneNumber domain.PhoneNumber)
	AddPasswordData(secret *crypto.CryptoValue, changeRequired bool)
}

func AddHumanCommand(a *user.Aggregate, human *AddHuman, passwordAlg crypto.HashAlgorithm, codeAlg crypto.EncryptionAlgorithm) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		if err := human.Email.Validate(); err != nil {
			return nil, err
		}
		if human.Username = strings.TrimSpace(human.Username); human.Username == "" {
			return nil, errors.ThrowInvalidArgument(nil, "V2-zzad3", "Errors.Invalid.Argument")
		}

		if human.FirstName = strings.TrimSpace(human.FirstName); human.FirstName == "" {
			return nil, errors.ThrowInvalidArgument(nil, "USER-UCej2", "Errors.User.Profile.FirstNameEmpty")
		}
		if human.LastName = strings.TrimSpace(human.LastName); human.LastName == "" {
			return nil, errors.ThrowInvalidArgument(nil, "USER-4hB7d", "Errors.User.Profile.LastNameEmpty")
		}
		human.ensureDisplayName()

		if human.Phone.Number != "" {
			if human.Phone.Number, err = human.Phone.Number.Normalize(); err != nil {
				return nil, err
			}
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			domainPolicy, err := domainPolicyWriteModel(ctx, filter, a.ResourceOwner)
			if err != nil {
				return nil, err
			}

			if err = userValidateDomain(ctx, a, human.Username, domainPolicy.UserLoginMustBeDomain, filter); err != nil {
				return nil, err
			}

			var createCmd humanCreationCommand
			if human.Register {
				createCmd = user.NewHumanRegisteredEvent(
					ctx,
					&a.Aggregate,
					human.Username,
					human.FirstName,
					human.LastName,
					human.NickName,
					human.DisplayName,
					human.PreferredLanguage,
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
					human.PreferredLanguage,
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

			if human.BcryptedPassword != "" {
				createCmd.AddPasswordData(crypto.FillHash([]byte(human.BcryptedPassword), passwordAlg), human.PasswordChangeRequired)
			}

			cmds := make([]eventstore.Command, 0, 3)
			cmds = append(cmds, createCmd)

			if human.Email.Verified {
				cmds = append(cmds, user.NewHumanEmailVerifiedEvent(ctx, &a.Aggregate))
			}
			//add init code if
			// email not verified or
			// user not registered and password set
			if human.shouldAddInitCode() {
				value, expiry, err := newUserInitCode(ctx, filter, codeAlg)
				if err != nil {
					return nil, err
				}
				cmds = append(cmds, user.NewHumanInitialCodeAddedEvent(ctx, &a.Aggregate, value, expiry))
			} else {
				if !human.Email.Verified {
					value, expiry, err := newEmailCode(ctx, filter, codeAlg)
					if err != nil {
						return nil, err
					}
					cmds = append(cmds, user.NewHumanEmailCodeAddedEvent(ctx, &a.Aggregate, value, expiry))
				}
			}

			if human.Phone.Verified {
				cmds = append(cmds, user.NewHumanPhoneVerifiedEvent(ctx, &a.Aggregate))
			} else if human.Phone.Number != "" {
				value, expiry, err := newPhoneCode(ctx, filter, codeAlg)
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

	index := strings.LastIndex(username, "@")
	if index < 0 {
		return nil
	}

	domainCheck := NewOrgDomainVerifiedWriteModel(username[index+1:])
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

// shouldAddInitCode returns true for all added Humans which:
// - were not added from an external IDP
// - and either:
//   - have no verified email
//     and / or
//   - have no authentication method (password / passwordless)
func (h *AddHuman) shouldAddInitCode() bool {
	return !h.ExternalIDP &&
		!h.Email.Verified ||
		!h.Passwordless &&
			h.Password == ""
}

func (c *Commands) ImportHuman(ctx context.Context, orgID string, human *domain.Human, passwordless bool, links []*domain.UserIDPLink, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessCodeGenerator crypto.Generator) (_ *domain.Human, passwordlessCode *domain.PasswordlessInitCode, err error) {
	if orgID == "" {
		return nil, nil, errors.ThrowInvalidArgument(nil, "COMMAND-5N8fs", "Errors.ResourceOwnerMissing")
	}
	domainPolicy, err := c.getOrgDomainPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, errors.ThrowPreconditionFailed(err, "COMMAND-2N9fs", "Errors.Org.DomainPolicy.NotFound")
	}
	pwPolicy, err := c.getOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, errors.ThrowPreconditionFailed(err, "COMMAND-4N8gs", "Errors.Org.PasswordComplexityPolicy.NotFound")
	}

	if human.AggregateID != "" {
		existing, err := c.getHumanWriteModelByID(ctx, human.AggregateID, human.ResourceOwner)
		if err != nil {
			return nil, nil, err
		}

		if existing.UserState != domain.UserStateUnspecified {
			return nil, nil, errors.ThrowPreconditionFailed(nil, "COMMAND-ziuna", "Errors.User.AlreadyExisting")
		}
	}

	events, addedHuman, addedCode, code, err := c.importHuman(ctx, orgID, human, passwordless, links, domainPolicy, pwPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessCodeGenerator)
	if err != nil {
		return nil, nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, nil, err
	}

	err = AppendAndReduce(addedHuman, pushedEvents...)
	if err != nil {
		return nil, nil, err
	}
	if addedCode != nil {
		err = AppendAndReduce(addedCode, pushedEvents...)
		if err != nil {
			return nil, nil, err
		}
		passwordlessCode = writeModelToPasswordlessInitCode(addedCode, code)
	}

	return writeModelToHuman(addedHuman), passwordlessCode, nil
}

func (c *Commands) RegisterHuman(ctx context.Context, orgID string, human *domain.Human, link *domain.UserIDPLink, orgMemberRoles []string, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator crypto.Generator) (*domain.Human, error) {
	if orgID == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-GEdf2", "Errors.ResourceOwnerMissing")
	}
	domainPolicy, err := c.getOrgDomainPolicy(ctx, orgID)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "COMMAND-33M9f", "Errors.Org.DomainPolicy.NotFound")
	}
	pwPolicy, err := c.getOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "COMMAND-M5Fsd", "Errors.Org.PasswordComplexityPolicy.NotFound")
	}
	loginPolicy, err := c.getOrgLoginPolicy(ctx, orgID)
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "COMMAND-Dfg3g", "Errors.Org.LoginPolicy.NotFound")
	}
	if !loginPolicy.AllowRegister {
		return nil, errors.ThrowPreconditionFailed(err, "COMMAND-SAbr3", "Errors.Org.LoginPolicy.RegistrationNotAllowed")
	}
	userEvents, registeredHuman, err := c.registerHuman(ctx, orgID, human, link, domainPolicy, pwPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator)
	if err != nil {
		return nil, err
	}

	orgMemberWriteModel := NewOrgMemberWriteModel(orgID, registeredHuman.AggregateID)
	orgAgg := OrgAggregateFromWriteModel(&orgMemberWriteModel.WriteModel)
	if len(orgMemberRoles) > 0 {
		orgMember := &domain.Member{
			ObjectRoot: models.ObjectRoot{
				AggregateID: orgID,
			},
			UserID: human.AggregateID,
			Roles:  orgMemberRoles,
		}
		memberEvent, err := c.addOrgMember(ctx, orgAgg, orgMemberWriteModel, orgMember)
		if err != nil {
			return nil, err
		}
		userEvents = append(userEvents, memberEvent)
	}

	pushedEvents, err := c.eventstore.Push(ctx, userEvents...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(registeredHuman, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToHuman(registeredHuman), nil
}

func (c *Commands) importHuman(ctx context.Context, orgID string, human *domain.Human, passwordless bool, links []*domain.UserIDPLink, domainPolicy *domain.DomainPolicy, pwPolicy *domain.PasswordComplexityPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator, passwordlessCodeGenerator crypto.Generator) (events []eventstore.Command, humanWriteModel *HumanWriteModel, passwordlessCodeWriteModel *HumanPasswordlessInitCodeWriteModel, code string, err error) {
	if orgID == "" {
		return nil, nil, nil, "", errors.ThrowInvalidArgument(nil, "COMMAND-00p2b", "Errors.Org.Empty")
	}
	if err := human.Normalize(); err != nil {
		return nil, nil, nil, "", err
	}
	events, humanWriteModel, err = c.createHuman(ctx, orgID, human, links, false, passwordless, domainPolicy, pwPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator)
	if err != nil {
		return nil, nil, nil, "", err
	}
	if passwordless {
		var codeEvent eventstore.Command
		codeEvent, passwordlessCodeWriteModel, code, err = c.humanAddPasswordlessInitCode(ctx, human.AggregateID, orgID, true, passwordlessCodeGenerator)
		if err != nil {
			return nil, nil, nil, "", err
		}
		events = append(events, codeEvent)
	}
	return events, humanWriteModel, passwordlessCodeWriteModel, code, nil
}

func (c *Commands) registerHuman(ctx context.Context, orgID string, human *domain.Human, link *domain.UserIDPLink, domainPolicy *domain.DomainPolicy, pwPolicy *domain.PasswordComplexityPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator crypto.Generator) ([]eventstore.Command, *HumanWriteModel, error) {
	if human == nil {
		return nil, nil, errors.ThrowInvalidArgument(nil, "COMMAND-JKefw", "Errors.User.Invalid")
	}
	if human.Username = strings.TrimSpace(human.Username); human.Username == "" {
		human.Username = string(human.EmailAddress)
	}
	if orgID == "" {
		return nil, nil, errors.ThrowInvalidArgument(nil, "COMMAND-hYsVH", "Errors.Org.Empty")
	}
	if err := human.Normalize(); err != nil {
		return nil, nil, err
	}
	if link == nil && (human.Password == nil || human.Password.SecretString == "") {
		return nil, nil, errors.ThrowInvalidArgument(nil, "COMMAND-X23na", "Errors.User.Password.Empty")
	}
	if human.Password != nil && human.Password.SecretString != "" {
		human.Password.ChangeRequired = false
	}
	var links []*domain.UserIDPLink
	if link != nil {
		links = append(links, link)
	}
	return c.createHuman(ctx, orgID, human, links, true, false, domainPolicy, pwPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator)
}

func (c *Commands) createHuman(ctx context.Context, orgID string, human *domain.Human, links []*domain.UserIDPLink, selfregister, passwordless bool, domainPolicy *domain.DomainPolicy, pwPolicy *domain.PasswordComplexityPolicy, initCodeGenerator, emailCodeGenerator, phoneCodeGenerator crypto.Generator) (events []eventstore.Command, addedHuman *HumanWriteModel, err error) {
	if err := human.CheckDomainPolicy(domainPolicy); err != nil {
		return nil, nil, err
	}
	human.Username = strings.TrimSpace(human.Username)
	human.EmailAddress = human.EmailAddress.Normalize()
	if !domainPolicy.UserLoginMustBeDomain {
		index := strings.LastIndex(human.Username, "@")
		if index > 1 {
			domainCheck := NewOrgDomainVerifiedWriteModel(human.Username[index+1:])
			if err := c.eventstore.FilterToQueryReducer(ctx, domainCheck); err != nil {
				return nil, nil, err
			}
			if domainCheck.Verified && domainCheck.ResourceOwner != orgID {
				return nil, nil, errors.ThrowInvalidArgument(nil, "COMMAND-SFd21", "Errors.User.DomainNotAllowedAsUsername")
			}
		}
	}

	if human.AggregateID == "" {
		userID, err := c.idGenerator.Next()
		if err != nil {
			return nil, nil, err
		}
		human.AggregateID = userID
	}

	human.SetNamesAsDisplayname()
	if human.Password != nil {
		if err := human.HashPasswordIfExisting(pwPolicy, c.userPasswordAlg, human.Password.ChangeRequired); err != nil {
			return nil, nil, err
		}
	}

	addedHuman = NewHumanWriteModel(human.AggregateID, orgID)
	//TODO: adlerhurst maybe we could simplify the code below
	userAgg := UserAggregateFromWriteModel(&addedHuman.WriteModel)

	if selfregister {
		events = append(events, createRegisterHumanEvent(ctx, userAgg, human, domainPolicy.UserLoginMustBeDomain))
	} else {
		events = append(events, createAddHumanEvent(ctx, userAgg, human, domainPolicy.UserLoginMustBeDomain))
	}

	for _, link := range links {
		event, err := c.addUserIDPLink(ctx, userAgg, link)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, event)
	}

	if human.IsInitialState(passwordless, len(links) > 0) {
		initCode, err := domain.NewInitUserCode(initCodeGenerator)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, user.NewHumanInitialCodeAddedEvent(ctx, userAgg, initCode.Code, initCode.Expiry))
	} else {
		if human.Email != nil && human.EmailAddress != "" && human.IsEmailVerified {
			events = append(events, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
		} else {
			emailCode, err := domain.NewEmailCode(emailCodeGenerator)
			if err != nil {
				return nil, nil, err
			}
			events = append(events, user.NewHumanEmailCodeAddedEvent(ctx, userAgg, emailCode.Code, emailCode.Expiry))
		}
	}

	if human.Phone != nil && human.PhoneNumber != "" && !human.IsPhoneVerified {
		phoneCode, err := domain.NewPhoneCode(phoneCodeGenerator)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, user.NewHumanPhoneCodeAddedEvent(ctx, userAgg, phoneCode.Code, phoneCode.Expiry))
	} else if human.Phone != nil && human.PhoneNumber != "" && human.IsPhoneVerified {
		events = append(events, user.NewHumanPhoneVerifiedEvent(ctx, userAgg))
	}

	return events, addedHuman, nil
}

func (c *Commands) HumanSkipMFAInit(ctx context.Context, userID, resourceowner string) (err error) {
	if userID == "" {
		return errors.ThrowInvalidArgument(nil, "COMMAND-2xpX9", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.getHumanWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return errors.ThrowNotFound(nil, "COMMAND-m9cV8", "Errors.User.NotFound")
	}

	_, err = c.eventstore.Push(ctx,
		user.NewHumanMFAInitSkippedEvent(ctx, UserAggregateFromWriteModel(&existingHuman.WriteModel)))
	return err
}

// TODO: adlerhurst maybe we can simplify createAddHumanEvent and createRegisterHumanEvent
func createAddHumanEvent(ctx context.Context, aggregate *eventstore.Aggregate, human *domain.Human, userLoginMustBeDomain bool) *user.HumanAddedEvent {
	addEvent := user.NewHumanAddedEvent(
		ctx,
		aggregate,
		human.Username,
		human.FirstName,
		human.LastName,
		human.NickName,
		human.DisplayName,
		human.PreferredLanguage,
		human.Gender,
		human.EmailAddress,
		userLoginMustBeDomain,
	)
	if human.Phone != nil {
		addEvent.AddPhoneData(human.PhoneNumber)
	}
	if human.Address != nil {
		addEvent.AddAddressData(
			human.Country,
			human.Locality,
			human.PostalCode,
			human.Region,
			human.StreetAddress)
	}
	if human.Password != nil {
		addEvent.AddPasswordData(human.Password.SecretCrypto, human.Password.ChangeRequired)
	}
	if human.HashedPassword != nil {
		addEvent.AddPasswordData(human.HashedPassword.SecretCrypto, false)
	}
	return addEvent
}

func createRegisterHumanEvent(ctx context.Context, aggregate *eventstore.Aggregate, human *domain.Human, userLoginMustBeDomain bool) *user.HumanRegisteredEvent {
	addEvent := user.NewHumanRegisteredEvent(
		ctx,
		aggregate,
		human.Username,
		human.FirstName,
		human.LastName,
		human.NickName,
		human.DisplayName,
		human.PreferredLanguage,
		human.Gender,
		human.EmailAddress,
		userLoginMustBeDomain,
	)
	if human.Phone != nil {
		addEvent.AddPhoneData(human.PhoneNumber)
	}
	if human.Address != nil {
		addEvent.AddAddressData(
			human.Country,
			human.Locality,
			human.PostalCode,
			human.Region,
			human.StreetAddress)
	}
	if human.Password != nil {
		addEvent.AddPasswordData(human.Password.SecretCrypto, human.Password.ChangeRequired)
	}
	if human.HashedPassword != nil {
		addEvent.AddPasswordData(human.HashedPassword.SecretCrypto, false)
	}
	return addEvent
}

func (c *Commands) HumansSignOut(ctx context.Context, agentID string, userIDs []string) error {
	if agentID == "" {
		return errors.ThrowInvalidArgument(nil, "COMMAND-2M0ds", "Errors.User.UserIDMissing")
	}
	if len(userIDs) == 0 {
		return errors.ThrowInvalidArgument(nil, "COMMAND-M0od3", "Errors.User.UserIDMissing")
	}
	events := make([]eventstore.Command, 0)
	for _, userID := range userIDs {
		existingUser, err := c.getHumanWriteModelByID(ctx, userID, "")
		if err != nil {
			return err
		}
		if !isUserStateExists(existingUser.UserState) {
			continue
		}
		events = append(events, user.NewHumanSignedOutEvent(
			ctx,
			UserAggregateFromWriteModel(&existingUser.WriteModel),
			agentID))
	}
	if len(events) == 0 {
		return nil
	}
	_, err := c.eventstore.Push(ctx, events...)
	return err
}

func (c *Commands) getHumanWriteModelByID(ctx context.Context, userID, resourceowner string) (*HumanWriteModel, error) {
	humanWriteModel := NewHumanWriteModel(userID, resourceowner)
	err := c.eventstore.FilterToQueryReducer(ctx, humanWriteModel)
	if err != nil {
		return nil, err
	}
	return humanWriteModel, nil
}
