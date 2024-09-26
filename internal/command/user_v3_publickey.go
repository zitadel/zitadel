package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddPublicKey struct {
	ResourceOwner string
	UserID        string

	ExpirationDate time.Time
	PublicKey      []byte
	PrivateKey     []byte
}

func (wm *AddPublicKey) GetExpirationDate() time.Time {
	return wm.ExpirationDate
}

func (wm *AddPublicKey) SetExpirationDate(date time.Time) {
	wm.ExpirationDate = date
}

func (wm *AddPublicKey) SetPublicKey(data []byte) {
	wm.PublicKey = data
}

func (wm *AddPublicKey) SetPrivateKey(data []byte) {
	wm.PrivateKey = data
}

func (c *Commands) AddPublicKey(ctx context.Context, pk *AddPublicKey) (*domain.ObjectDetails, error) {
	if pk.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-14sGR7lTaj", "Errors.IDMissing")
	}
	schemauser, err := existingSchemaUser(ctx, c, pk.ResourceOwner, pk.UserID)
	if err != nil {
		return nil, err
	}

	_, err = existingSchema(ctx, c, "", schemauser.SchemaID)
	if err != nil {
		return nil, err
	}
	// TODO check for possible authenticators

	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	writeModel, err := c.getSchemaPublicKeyWM(ctx, schemauser.ResourceOwner, schemauser.AggregateID, id)
	if err != nil {
		return nil, err
	}

	if len(pk.PublicKey) == 0 {
		if err := domain.SetNewAuthNKeyPair(pk, c.machineKeySize); err != nil {
			return nil, err
		}
	}

	events, err := writeModel.NewCreate(ctx, pk.ExpirationDate, pk.PublicKey)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) DeletePublicKey(ctx context.Context, resourceOwner, userID, id string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-hzqeAXW1qP", "Errors.IDMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-BNNYJz6Yxt", "Errors.IDMissing")
	}

	writeModel, err := c.getSchemaPublicKeyWM(ctx, resourceOwner, userID, id)
	if err != nil {
		return nil, err
	}

	events, err := writeModel.NewDelete(ctx)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) getSchemaPublicKeyWM(ctx context.Context, resourceOwner, userID, id string) (*PublicKeyV3WriteModel, error) {
	writeModel := NewPublicKeyV3WriteModel(resourceOwner, userID, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}
