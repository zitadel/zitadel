package command

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	PATPrefix = "pat_"
)

type AddPAT struct {
	ResourceOwner string
	UserID        string

	ExpirationDate time.Time
	Scope          []string
	Token          string
}

func (wm *AddPAT) GetExpirationDate() time.Time {
	return wm.ExpirationDate
}

func (wm *AddPAT) SetExpirationDate(date time.Time) {
	wm.ExpirationDate = date
}

func (c *Commands) AddPAT(ctx context.Context, pat *AddPAT) (*domain.ObjectDetails, error) {
	if pat.UserID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-14sGR7lTaj", "Errors.IDMissing")
	}
	schemauser, err := existingSchemaUser(ctx, c, pat.ResourceOwner, pat.UserID)
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
	writeModel, err := c.getSchemaPATWM(ctx, schemauser.ResourceOwner, schemauser.AggregateID, id)
	if err != nil {
		return nil, err
	}

	events, err := writeModel.NewCreate(ctx, pat.ExpirationDate, pat.Scope)
	if err != nil {
		return nil, err
	}
	pat.Token, err = createSchemaUserPAT(c.keyAlgorithm, writeModel.AggregateID, pat.UserID)
	if err != nil {
		return nil, err
	}
	return c.pushAppendAndReduceDetails(ctx, writeModel, events...)
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
