package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddUserHuman(ctx context.Context, resourceOwner string, human *AddHuman, allowInitMail bool) (err error) {
	if resourceOwner == "" {
		return errors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal")
	}

	if err := human.Validate(c.userPasswordHasher); err != nil {
		return err
	}

	if human.ID == "" {
		human.ID, err = c.idGenerator.Next()
		if err != nil {
			return err
		}
	}

	existingHuman, err := c.userHumanWriteModel(ctx, human.ID, resourceOwner)
	if err != nil {
		return err
	}
	if isUserStateExists(existingHuman.UserState) {
		return errors.ThrowPreconditionFailed(nil, "COMMAND-k2unb", "Errors.User.AlreadyExisting")
	}

	domainPolicy, err := c.domainPolicyWriteModel(ctx, resourceOwner)
	if err != nil {
		return err
	}

	if err = c.userValidateDomain(ctx, resourceOwner, human.Username, domainPolicy.UserLoginMustBeDomain); err != nil {
		return err
	}

	var createCmd humanCreationCommand
	if human.Register {
		createCmd = user.NewHumanRegisteredEvent(
			ctx,
			UserAggregateFromWriteModel(&existingHuman.WriteModel),
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
			UserAggregateFromWriteModel(&existingHuman.WriteModel),
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

	if err := c.addHumanCommandPassword(ctx, createCmd, human, c.userPasswordHasher); err != nil {
		return err
	}

	cmds := make([]eventstore.Command, 0, 3)
	cmds = append(cmds, createCmd)
	filter := c.eventstore.Filter

	cmds, err = c.addHumanCommandEmail(ctx, filter, cmds, existingHuman.Aggregate(), human, c.userEncryption, allowInitMail)
	if err != nil {
		return err
	}

	cmds, err = c.addHumanCommandPhone(ctx, filter, cmds, existingHuman.Aggregate(), human, c.userEncryption)
	if err != nil {
		return err
	}

	for _, metadataEntry := range human.Metadata {
		cmds = append(cmds, user.NewMetadataSetEvent(
			ctx,
			UserAggregateFromWriteModel(&existingHuman.WriteModel),
			metadataEntry.Key,
			metadataEntry.Value,
		))
	}
	for _, link := range human.Links {
		cmd, err := addLink(ctx, filter, existingHuman.Aggregate(), link)
		if err != nil {
			return err
		}
		cmds = append(cmds, cmd)
	}

	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return err
	}
	human.Details = &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().ResourceOwner,
	}
	return nil
}

func (c *Commands) userHumanWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *UserHumanWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewUserHumanWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) orgDomainVerifiedWriteModel(ctx context.Context, domain string) (writeModel *OrgDomainVerifiedWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewOrgDomainVerifiedWriteModel(domain)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}
