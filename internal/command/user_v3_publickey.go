package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddPublicKey struct {
	ResourceOwner string
	UserID        string

	PublicKey *PublicKey
}

func (wm *AddPublicKey) GetPublicKey() []byte {
	if wm.PublicKey == nil {
		return nil
	}
	return wm.PublicKey.PublicKey
}

func (wm *AddPublicKey) GetPrivateKey() []byte {
	if wm.PublicKey == nil {
		return nil
	}
	return wm.PublicKey.PrivateKey
}

func (wm *AddPublicKey) GetExpirationDate() time.Time {
	if wm.PublicKey == nil {
		return time.Time{}
	}
	return wm.PublicKey.GetExpirationDate()
}

func (wm *AddPublicKey) SetExpirationDate(date time.Time) {
	if wm.PublicKey == nil {
		wm.PublicKey = &PublicKey{}
	}
	wm.PublicKey.SetExpirationDate(date)
}

type PublicKey struct {
	ExpirationDate time.Time
	PublicKey      []byte
	PrivateKey     []byte
}

func (wm *PublicKey) GetExpirationDate() time.Time {
	return wm.ExpirationDate
}

func (wm *PublicKey) SetExpirationDate(date time.Time) {
	wm.ExpirationDate = date
}

func (wm *PublicKey) SetPublicKey(data []byte) {
	wm.PublicKey = data
}

func (wm *PublicKey) SetPrivateKey(data []byte) {
	wm.PrivateKey = data
}

func (c *Commands) AddPublicKey(ctx context.Context, add *AddPublicKey) (*domain.ObjectDetails, error) {
	if add.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-14sGR7lTaj", "Errors.IDMissing")
	}
	if publicKey := add.GetPublicKey(); publicKey != nil {
		if _, err := crypto.BytesToPublicKey(publicKey); err != nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-WdWlhUSVqK", "Errors.User.Machine.Key.Invalid")
		}
	}
	if err := domain.EnsureValidExpirationDate(add.PublicKey); err != nil {
		return nil, err
	}

	schemauser, err := existingSchemaUser(ctx, c, add.ResourceOwner, add.UserID)
	if err != nil {
		return nil, err
	}
	add.ResourceOwner = schemauser.ResourceOwner

	_, err = existingSchema(ctx, c, "", schemauser.SchemaID)
	if err != nil {
		return nil, err
	}
	// TODO check for possible authenticators

	writeModel, events, err := c.addPublicKey(ctx, add.ResourceOwner, add.UserID, add.PublicKey)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) addPublicKey(ctx context.Context, resourceOwner, userID string, add *PublicKey) (*PublicKeyV3WriteModel, []eventstore.Command, error) {
	if add == nil {
		return nil, nil, nil
	}
	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	writeModel, err := c.getSchemaPublicKeyWM(ctx, resourceOwner, userID, id)
	if err != nil {
		return nil, nil, err
	}
	if len(add.PublicKey) == 0 {
		if err := domain.SetNewAuthNKeyPair(add, c.machineKeySize); err != nil {
			return nil, nil, err
		}
	}
	events, err := writeModel.NewCreate(ctx, add.ExpirationDate, add.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return writeModel, events, nil
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
