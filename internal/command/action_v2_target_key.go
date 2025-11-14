package command

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ErrTargetIDMissing = "ErrTargetIDMissing"
	ErrPublicKeyFormat = "ErrPublicKeyFormat"
)

type TargetPublicKey struct {
	TargetID  string
	PublicKey []byte

	// KeyID is generated and returned after adding the public key
	KeyID string
}

func (c *Commands) AddTargetPublicKey(ctx context.Context, key *TargetPublicKey, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if key.TargetID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing")
	}
	_, err = checkPublicKey(key.PublicKey)
	if err != nil {
		return nil, err
	}
	id, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	key.KeyID = id
	writeModel, err := c.getTargetKeyWriteModelByID(ctx, key.TargetID, key.KeyID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !writeModel.TargetExists {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-nd3fsd", "Errors.Target.NotFound")
	}
	err = c.pushAppendAndReduce(ctx, writeModel, target.NewKeyAddedEvent(
		ctx,
		TargetAggregateFromWriteModel(&writeModel.WriteModel),
		key.KeyID,
		key.PublicKey,
	))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) RemoveTargetPublicKey(ctx context.Context, targetID, keyID, resourceOwner string) (_ *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if targetID == "" || keyID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing")
	}
	writeModel, err := c.getTargetKeyWriteModelByID(ctx, targetID, keyID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !writeModel.TargetExists || !writeModel.KeyExists {
		return writeModelToObjectDetails(&writeModel.WriteModel), nil
	}
	err = c.pushAppendAndReduce(ctx, writeModel, target.NewKeyRemovedEvent(
		ctx,
		TargetAggregateFromWriteModel(&writeModel.WriteModel),
		keyID,
	))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func checkPublicKey(key []byte) (crypto.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, zerrors.ThrowInvalidArgument(nil, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey")
	}
	switch publicKey := pub.(type) {
	case *rsa.PublicKey,
		*ecdsa.PublicKey,
		ed25519.PublicKey:
		return publicKey, nil
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey")
	}
}

func (c *Commands) getTargetKeyWriteModelByID(ctx context.Context, targetID, keyID string, resourceOwner string) (*TargetKeyWriteModel, error) {
	wm := NewTargetKeyWriteModel(targetID, keyID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}
