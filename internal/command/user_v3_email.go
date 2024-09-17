package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ChangeSchemaUserEmail struct {
	Details *domain.ObjectDetails

	ResourceOwner string
	ID            string

	Email      *Email
	ReturnCode string
}

func (s *ChangeSchemaUserEmail) Valid() (err error) {
	if s.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-lvoHfcR8zQ", "Errors.IDMissing")
	}
	if s.Email != nil && s.Email.Address != "" {
		if err := s.Email.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Commands) ChangeSchemaUserEmail(ctx context.Context, user *ChangeSchemaUserEmail, alg crypto.EncryptionAlgorithm) (err error) {
	if err := user.Valid(); err != nil {
		return err
	}

	writeModel, err := c.getSchemaUserEmailWriteModelByID(ctx, user.ResourceOwner, user.ID)
	if err != nil {
		return err
	}
	if !writeModel.Exists() {
		return zerrors.ThrowNotFound(nil, "COMMAND-nJ0TQFuRmP", "Errors.User.NotFound")
	}
	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return err
	}

	events := make([]eventstore.Command, 0)
	events, user.ReturnCode, err = c.updateSchemaUserEmail(ctx, writeModel, events, UserV3AggregateFromWriteModel(&writeModel.WriteModel), user.Email, alg)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		user.Details = writeModelToObjectDetails(&writeModel.WriteModel)
		return nil
	}

	if err := c.pushAppendAndReduce(ctx, writeModel, events...); err != nil {
		return err
	}
	user.Details = writeModelToObjectDetails(&writeModel.WriteModel)
	return nil
}

func (c *Commands) VerifySchemaUserEmail(ctx context.Context, resourceOwner, id, code string, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-0oj2PquNGA", "Errors.IDMissing")
	}
	writeModel, err := c.getSchemaUserEmailWriteModelByID(ctx, resourceOwner, id)
	if err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-bRfVIJnYP6", "Errors.User.NotFound")
	}
	if writeModel.EmailCode == nil {
		return writeModelToObjectDetails(&writeModel.WriteModel), nil
	}
	if err := c.checkPermissionUpdateUser(ctx, writeModel.ResourceOwner, writeModel.AggregateID); err != nil {
		return nil, err
	}
	if err := crypto.VerifyCode(writeModel.EmailCode.CreationDate, writeModel.EmailCode.Expiry, writeModel.EmailCode.Code, code, alg); err != nil {
		return nil, err
	}
	if err := c.pushAppendAndReduce(ctx, writeModel,
		schemauser.NewEmailVerifiedEvent(ctx, UserV3AggregateFromWriteModel(&writeModel.WriteModel)),
	); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) updateSchemaUserEmail(ctx context.Context, existing *UserV3WriteModel, events []eventstore.Command, agg *eventstore.Aggregate, email *Email, alg crypto.EncryptionAlgorithm) (_ []eventstore.Command, plainCode string, err error) {
	if existing.Email == string(email.Address) {
		return events, plainCode, nil
	}

	events = append(events, schemauser.NewEmailUpdatedEvent(ctx,
		agg,
		email.Address,
	))
	if email.Verified {
		events = append(events, schemauser.NewEmailVerifiedEvent(ctx, agg))
	} else {
		cryptoCode, err := c.newEmailCode(ctx, c.eventstore.Filter, alg) //nolint:staticcheck
		if err != nil {
			return nil, "", err
		}
		if email.ReturnCode {
			plainCode = cryptoCode.Plain
		}
		events = append(events, schemauser.NewEmailCodeAddedEvent(ctx, agg,
			cryptoCode.Crypted,
			cryptoCode.Expiry,
			email.URLTemplate,
			email.ReturnCode,
		))
	}
	return events, plainCode, nil
}
