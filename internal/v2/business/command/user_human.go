package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) AddHuman(ctx context.Context, orgID, username string, human *usr_model.Human) (*usr_model.Human, error) {
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
	//pwPolicy, err := r.GetOrgPasswordComplexityPolicy(ctx, orgID)
	//if err != nil {
	//	return nil, err
	//}

	addedHuman := NewHumanWriteModel(human.AggregateID)
	//TODO: Check Unique Username
	human.CheckOrgIAMPolicy(username, orgIAMPolicy)
	human.SetNamesAsDisplayname()
	//human.HashPasswordIfExisting(pwPolicy, r.userPasswordAlg, true)

	userAgg := UserAggregateFromWriteModel(&addedHuman.WriteModel)
	userAgg.PushEvents(
		user.NewHumanAddedEvent(
			ctx,
			username,
			human.FirstName,
			human.LastName,
			human.NickName,
			human.DisplayName,
			human.PreferredLanguage,
			domain.Gender(human.Gender),
			human.EmailAddress,
			human.PhoneNumber,
			human.Country,
			human.Locality,
			human.PostalCode,
			human.Region,
			human.StreetAddress,
		),
	)
	//TODO: HashPassword If existing
	//TODO: Generate Init Code if needed
	//TODO: Generate Phone Code if needed
	if human.Email != nil && human.EmailAddress != "" && human.IsEmailVerified {
		userAgg.PushEvents(user.NewHumanEmailVerifiedEvent(ctx))
	}
	if human.Phone != nil && human.PhoneNumber != "" && human.IsPhoneVerified {
		userAgg.PushEvents(user.NewHumanPhoneVerifiedEvent(ctx))
	}

	err = r.eventstore.PushAggregate(ctx, addedHuman, userAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToHuman(addedHuman), nil
}
