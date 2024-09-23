package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SetSchemaUserPassword struct {
	ResourceOwner string
	UserID        string

	Password            string
	EncodedPasswordHash string
	ChangeRequired      bool
}

func (p *SetSchemaUserPassword) Validate(hasher *crypto.Hasher) (err error) {
	if p.UserID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing")
	}

	if p.EncodedPasswordHash != "" {
		if !hasher.EncodingSupported(p.EncodedPasswordHash) {
			return zerrors.ThrowInvalidArgument(nil, "COMMAND-oz74onzvqr", "Errors.User.Password.NotSupported")
		}
	}
	if p.Password == "" && p.EncodedPasswordHash == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty")
	}

	return nil
}

func (c *Commands) SetSchemaUserPassword(ctx context.Context, username *SetSchemaUserPassword) (*domain.ObjectDetails, error) {
	if err := username.Validate(c.userPasswordHasher); err != nil {
		return nil, err
	}

	existing, err := c.getPasswordExistsWithVerification(ctx, username.ResourceOwner, username.UserID)
	if err != nil {
		return nil, err
	}
	resourceOwner := existing.ResourceOwner
	if existing.EncodedHash == "" {
		existingUser, err := c.getSchemaUserExists(ctx, username.ResourceOwner, username.UserID)
		if err != nil {
			return nil, err
		}
		if !existingUser.Exists() {
			return nil, zerrors.ThrowNotFound(nil, "TODO", "TODO")
		}
		resourceOwner = existingUser.ResourceOwner
	}

	// If password is provided, let's check if is compliant with the policy.
	// If only a encodedPassword is passed, we can skip this.
	if username.Password != "" {
		if err = c.checkPasswordComplexity(ctx, username.Password, resourceOwner); err != nil {
			return nil, err
		}
	}

	encodedPassword := username.EncodedPasswordHash
	if username.Password != "" {
		encodedPassword, err = c.userPasswordHasher.Hash(username.Password)
		if err = convertPasswapErr(err); err != nil {
			return nil, err
		}
	}

	events, err := c.eventstore.Push(ctx,
		authenticator.NewPasswordCreatedEvent(ctx,
			&authenticator.NewAggregate(username.UserID, resourceOwner).Aggregate,
			existing.UserID,
			encodedPassword,
			username.ChangeRequired,
		),
	)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) DeleteSchemaUserPassword(ctx context.Context, resourceOwner, id string) (_ *domain.ObjectDetails, err error) {
	existing, err := c.getPasswordExistsWithVerification(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if existing.EncodedHash == "" {
		return nil, zerrors.ThrowNotFound(nil, "TODO", "TODO")
	}

	events, err := c.eventstore.Push(ctx,
		authenticator.NewPasswordDeletedEvent(ctx,
			&authenticator.NewAggregate(id, existing.ResourceOwner).Aggregate,
		),
	)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func (c *Commands) getPasswordExistsWithVerification(ctx context.Context, resourceOwner, id string) (*PasswordV3WriteModel, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing")
	}
	writeModel := NewPasswordV3WriteModel(resourceOwner, id)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}

	// TODO permission through old password and password code
	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.UserID); err != nil {
		return nil, err
	}
	return writeModel, nil
}
