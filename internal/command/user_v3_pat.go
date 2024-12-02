package command

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	PATPrefix = "pat_"
)

type AddPAT struct {
	ResourceOwner string
	UserID        string

	PAT *PAT
}

type PAT struct {
	ExpirationDate time.Time
	Scopes         []string
	Token          string
}

func (wm *PAT) GetExpirationDate() time.Time {
	return wm.ExpirationDate
}

func (wm *PAT) SetExpirationDate(date time.Time) {
	wm.ExpirationDate = date
}

func (c *Commands) AddPAT(ctx context.Context, add *AddPAT) (*domain.ObjectDetails, error) {
	if add.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-14sGR7lTaj", "Errors.IDMissing")
	}
	if err := domain.EnsureValidExpirationDate(add.PAT); err != nil {
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

	writeModel, events, err := c.addPAT(ctx, add.ResourceOwner, add.UserID, add.PAT)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) addPAT(ctx context.Context, resourceOwner, userID string, add *PAT) (*PATV3WriteModel, []eventstore.Command, error) {
	if add == nil {
		return nil, nil, nil
	}
	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, err
	}
	writeModel, err := c.getSchemaPATWM(ctx, resourceOwner, userID, id)
	if err != nil {
		return nil, nil, err
	}
	events, err := writeModel.NewCreate(ctx, add.ExpirationDate, add.Scopes)
	if err != nil {
		return nil, nil, err
	}
	add.Token, err = createSchemaUserPAT(c.keyAlgorithm, writeModel.AggregateID, writeModel.UserID)
	if err != nil {
		return nil, nil, err
	}
	return writeModel, events, nil
}

func (c *Commands) DeletePAT(ctx context.Context, resourceOwner, userID, id string) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-hzqeAXW1qP", "Errors.IDMissing")
	}
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-BNNYJz6Yxt", "Errors.IDMissing")
	}

	writeModel, err := c.getSchemaPATWM(ctx, resourceOwner, userID, id)
	if err != nil {
		return nil, err
	}

	events, err := writeModel.NewDelete(ctx)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
}

func (c *Commands) getSchemaPATWM(ctx context.Context, resourceOwner, userID, id string) (*PATV3WriteModel, error) {
	writeModel := NewPATV3WriteModel(resourceOwner, userID, id, c.checkPermission)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	return writeModel, nil
}

func createSchemaUserPAT(algorithm crypto.EncryptionAlgorithm, tokenID, userID string) (string, error) {
	encrypted, err := algorithm.Encrypt([]byte(tokenID + ":" + userID))
	if err != nil {
		return "", err
	}
	return PATPrefix + base64.RawURLEncoding.EncodeToString(encrypted), nil
}
