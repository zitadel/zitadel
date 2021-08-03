package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/user"
)

func (c *Commands) getHuman(ctx context.Context, userID, resourceowner string) (*domain.Human, error) {
	human, err := c.getHumanWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(human.UserState) {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-M9dsd", "Errors.User.NotFound")
	}
	return writeModelToHuman(human), nil
}

func (c *Commands) AddHuman(ctx context.Context, orgID string, human *domain.Human) (*domain.Human, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-XYFk9", "Errors.ResourceOwnerMissing")
	}
	orgIAMPolicy, err := c.getOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-33M9f", "Errors.Org.OrgIAMPolicy.NotFound")
	}
	pwPolicy, err := c.getOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-M5Fsd", "Errors.Org.PasswordComplexity.NotFound")
	}
	events, addedHuman, err := c.addHuman(ctx, orgID, human, orgIAMPolicy, pwPolicy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(addedHuman, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToHuman(addedHuman), nil
}

func (c *Commands) ImportHuman(ctx context.Context, orgID string, human *domain.Human, passwordless bool) (_ *domain.Human, passwordlessCode *domain.PasswordlessInitCode, err error) {
	if orgID == "" {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-5N8fs", "Errors.ResourceOwnerMissing")
	}
	orgIAMPolicy, err := c.getOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-2N9fs", "Errors.Org.OrgIAMPolicy.NotFound")
	}
	pwPolicy, err := c.getOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-4N8gs", "Errors.Org.PasswordComplexity.NotFound")
	}
	events, addedHuman, addedCode, code, err := c.importHuman(ctx, orgID, human, passwordless, orgIAMPolicy, pwPolicy)
	if err != nil {
		return nil, nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
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

func (c *Commands) addHuman(ctx context.Context, orgID string, human *domain.Human, orgIAMPolicy *domain.OrgIAMPolicy, pwPolicy *domain.PasswordComplexityPolicy) ([]eventstore.EventPusher, *HumanWriteModel, error) {
	if orgID == "" || !human.IsValid() {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-67Ms8", "Errors.User.Invalid")
	}
	if human.Password != nil && human.SecretString != "" {
		human.ChangeRequired = true
	}
	return c.createHuman(ctx, orgID, human, nil, false, false, orgIAMPolicy, pwPolicy)
}

func (c *Commands) importHuman(ctx context.Context, orgID string, human *domain.Human, passwordless bool, orgIAMPolicy *domain.OrgIAMPolicy, pwPolicy *domain.PasswordComplexityPolicy) (events []eventstore.EventPusher, humanWriteModel *HumanWriteModel, passwordlessCodeWriteModel *HumanPasswordlessInitCodeWriteModel, code string, err error) {
	if orgID == "" || !human.IsValid() {
		return nil, nil, nil, "", caos_errs.ThrowInvalidArgument(nil, "COMMAND-00p2b", "Errors.User.Invalid")
	}
	events, humanWriteModel, err = c.createHuman(ctx, orgID, human, nil, false, passwordless, orgIAMPolicy, pwPolicy)
	if err != nil {
		return nil, nil, nil, "", err
	}
	if passwordless {
		var codeEvent eventstore.EventPusher
		codeEvent, passwordlessCodeWriteModel, code, err = c.humanAddPasswordlessInitCode(ctx, human.AggregateID, orgID, true)
		if err != nil {
			return nil, nil, nil, "", err
		}
		events = append(events, codeEvent)
	}
	return events, humanWriteModel, passwordlessCodeWriteModel, code, nil
}

func (c *Commands) RegisterHuman(ctx context.Context, orgID string, human *domain.Human, externalIDP *domain.ExternalIDP, orgMemberRoles []string) (*domain.Human, error) {
	userEvents, registeredHuman, err := c.registerHuman(ctx, orgID, human, externalIDP)
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

	pushedEvents, err := c.eventstore.PushEvents(ctx, userEvents...)
	if err != nil {
		return nil, err
	}

	err = AppendAndReduce(registeredHuman, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToHuman(registeredHuman), nil
}

func (c *Commands) registerHuman(ctx context.Context, orgID string, human *domain.Human, externalIDP *domain.ExternalIDP) ([]eventstore.EventPusher, *HumanWriteModel, error) {
	if human != nil && human.Username == "" {
		human.Username = human.EmailAddress
	}
	if orgID == "" || !human.IsValid() || externalIDP == nil && (human.Password == nil || human.SecretString == "") {
		return nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-9dk45", "Errors.User.Invalid")
	}
	orgIAMPolicy, err := c.getOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-33M9f", "Errors.Org.OrgIAMPolicy.NotFound")
	}
	pwPolicy, err := c.getOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, caos_errs.ThrowPreconditionFailed(err, "COMMAND-M5Fsd", "Errors.Org.PasswordComplexity.NotFound")
	}
	if human.Password != nil && human.SecretString != "" {
		human.ChangeRequired = false
	}
	return c.createHuman(ctx, orgID, human, externalIDP, true, false, orgIAMPolicy, pwPolicy)
}

func (c *Commands) createHuman(ctx context.Context, orgID string, human *domain.Human, externalIDP *domain.ExternalIDP, selfregister, passwordless bool, orgIAMPolicy *domain.OrgIAMPolicy, pwPolicy *domain.PasswordComplexityPolicy) ([]eventstore.EventPusher, *HumanWriteModel, error) {
	if err := human.CheckOrgIAMPolicy(orgIAMPolicy); err != nil {
		return nil, nil, err
	}
	userID, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	human.AggregateID = userID
	human.SetNamesAsDisplayname()
	if human.Password != nil {
		if err := human.HashPasswordIfExisting(pwPolicy, c.userPasswordAlg, human.ChangeRequired); err != nil {
			return nil, nil, err
		}
	}

	addedHuman := NewHumanWriteModel(human.AggregateID, orgID)
	//TODO: adlerhurst maybe we could simplify the code below
	userAgg := UserAggregateFromWriteModel(&addedHuman.WriteModel)
	var events []eventstore.EventPusher

	if selfregister {
		events = append(events, createRegisterHumanEvent(ctx, userAgg, human, orgIAMPolicy.UserLoginMustBeDomain))
	} else {
		events = append(events, createAddHumanEvent(ctx, userAgg, human, orgIAMPolicy.UserLoginMustBeDomain))
	}

	if externalIDP != nil {
		event, err := c.addHumanExternalIDP(ctx, userAgg, externalIDP)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, event)
	}

	if human.IsInitialState(passwordless) {
		initCode, err := domain.NewInitUserCode(c.initializeUserCode)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, user.NewHumanInitialCodeAddedEvent(ctx, userAgg, initCode.Code, initCode.Expiry))
	}

	if human.Email != nil && human.EmailAddress != "" && human.IsEmailVerified {
		events = append(events, user.NewHumanEmailVerifiedEvent(ctx, userAgg))
	}

	if human.Phone != nil && human.PhoneNumber != "" && !human.IsPhoneVerified {
		phoneCode, err := domain.NewPhoneCode(c.phoneVerificationCode)
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
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-2xpX9", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.getHumanWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return caos_errs.ThrowNotFound(nil, "COMMAND-m9cV8", "Errors.User.NotFound")
	}

	_, err = c.eventstore.PushEvents(ctx,
		user.NewHumanMFAInitSkippedEvent(ctx, UserAggregateFromWriteModel(&existingHuman.WriteModel)))
	return err
}

///TODO: adlerhurst maybe we can simplify createAddHumanEvent and createRegisterHumanEvent
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
		addEvent.AddPasswordData(human.SecretCrypto, human.ChangeRequired)
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
		addEvent.AddPasswordData(human.SecretCrypto, human.ChangeRequired)
	}
	return addEvent
}

func (c *Commands) HumansSignOut(ctx context.Context, agentID string, userIDs []string) error {
	if agentID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-2M0ds", "Errors.User.UserIDMissing")
	}
	if len(userIDs) == 0 {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-M0od3", "Errors.User.UserIDMissing")
	}
	events := make([]eventstore.EventPusher, 0)
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
	_, err := c.eventstore.PushEvents(ctx, events...)
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
