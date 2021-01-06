package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) AddHuman(ctx context.Context, orgID, username string, human *domain.Human) (*domain.Human, error) {
	if !human.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-4M90d", "Errors.User.Invalid")
	}
	userID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	human.AggregateID = userID
	orgIAMPolicy, err := r.GetOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return nil, err
	}
	pwPolicy, err := r.GetOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, err
	}

	addedHuman := NewHumanWriteModel(human.AggregateID)
	//TODO: Check Unique Username
	human.CheckOrgIAMPolicy(username, orgIAMPolicy)
	human.SetNamesAsDisplayname()
	human.HashPasswordIfExisting(pwPolicy, r.userPasswordAlg, true)

	userAgg := UserAggregateFromWriteModel(&addedHuman.WriteModel)
	addEvent := user.NewHumanAddedEvent(
		ctx,
		username,
		human.FirstName,
		human.LastName,
		human.NickName,
		human.DisplayName,
		human.PreferredLanguage,
		human.Gender,
		human.EmailAddress,
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
	userAgg.PushEvents(addEvent)

	if human.IsInitialState() {
		initCode, err := domain.NewInitUserCode(r.initializeUserCode)
		if err != nil {
			return nil, err
		}
		user.NewHumanInitialCodeAddedEvent(ctx, initCode.Code, initCode.Expiry)
	}
	if human.Email != nil && human.EmailAddress != "" && human.IsEmailVerified {
		userAgg.PushEvents(user.NewHumanEmailVerifiedEvent(ctx))
	}
	if human.Phone != nil && human.PhoneNumber != "" && !human.IsPhoneVerified {
		phoneCode, err := domain.NewPhoneCode(r.phoneVerificationCode)
		if err != nil {
			return nil, err
		}
		user.NewHumanPhoneCodeAddedEvent(ctx, phoneCode.Code, phoneCode.Expiry)
	} else if human.Phone != nil && human.PhoneNumber != "" && human.IsPhoneVerified {
		userAgg.PushEvents(user.NewHumanPhoneVerifiedEvent(ctx))
	}

	err = r.eventstore.PushAggregate(ctx, addedHuman, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToHuman(addedHuman), nil
}

func (r *CommandSide) RegisterHuman(ctx context.Context, orgID, username string, human *domain.Human, externalIDP *domain.ExternalIDP) (*domain.Human, error) {
	if !human.IsValid() || externalIDP == nil && (human.Password == nil || human.SecretString == "") {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-9dk45", "Errors.User.Invalid")
	}
	userID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	human.AggregateID = userID
	orgIAMPolicy, err := r.GetOrgIAMPolicy(ctx, orgID)
	if err != nil {
		return nil, err
	}
	pwPolicy, err := r.GetOrgPasswordComplexityPolicy(ctx, orgID)
	if err != nil {
		return nil, err
	}

	addedHuman := NewHumanWriteModel(human.AggregateID)
	//TODO: Check Unique Username or unique external idp
	human.CheckOrgIAMPolicy(username, orgIAMPolicy)
	human.SetNamesAsDisplayname()
	human.HashPasswordIfExisting(pwPolicy, r.userPasswordAlg, true)

	userAgg := UserAggregateFromWriteModel(&addedHuman.WriteModel)
	addEvent := user.NewHumanRegisteredEvent(
		ctx,
		username,
		human.FirstName,
		human.LastName,
		human.NickName,
		human.DisplayName,
		human.PreferredLanguage,
		human.Gender,
		human.EmailAddress,
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
	userAgg.PushEvents(addEvent)
	//TODO: Add External IDP Event
	if human.IsInitialState() {
		initCode, err := domain.NewInitUserCode(r.initializeUserCode)
		if err != nil {
			return nil, err
		}
		user.NewHumanInitialCodeAddedEvent(ctx, initCode.Code, initCode.Expiry)
	}

	if human.Email != nil && human.EmailAddress != "" && human.IsEmailVerified {
		userAgg.PushEvents(user.NewHumanEmailVerifiedEvent(ctx))
	}
	if human.Phone != nil && human.PhoneNumber != "" && !human.IsPhoneVerified {
		phoneCode, err := domain.NewPhoneCode(r.phoneVerificationCode)
		if err != nil {
			return nil, err
		}
		user.NewHumanPhoneCodeAddedEvent(ctx, phoneCode.Code, phoneCode.Expiry)
	}

	err = r.eventstore.PushAggregate(ctx, addedHuman, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToHuman(addedHuman), nil
}
