package command

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ErrTargetIDMissing          = "ErrTargetIDMissing"
	ErrPublicKeyFormat          = "ErrPublicKeyFormat"
	ErrExpirationDateBeforeNow  = "ErrExpirationDateBeforeNow"
	ErrPublicKeyDeleteActiveKey = "ErrPublicKeyDeleteActiveKey"
)

type TargetPublicKey struct {
	TargetID   string
	PublicKey  []byte
	Expiration time.Time

	// KeyID is generated and returned after adding the public key
	KeyID string
}

func (c *Commands) AddTargetPublicKey(ctx context.Context, key *TargetPublicKey, resourceOwner string) (_ time.Time, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if key.TargetID == "" {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing")
	}
	if !key.Expiration.IsZero() && key.Expiration.Before(time.Now()) {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, ErrExpirationDateBeforeNow, "Errors.Target.InvalidExpirationDate")
	}
	fingerprint, err := checkPublicKeyAndComputeFingerprint(key.PublicKey)
	if err != nil {
		return time.Time{}, err
	}
	writeModel, err := c.getTargetKeyWriteModelByID(ctx, key.TargetID, key.KeyID, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}
	if !writeModel.TargetExists {
		return time.Time{}, zerrors.ThrowPreconditionFailed(nil, "COMMAND-nd3fsd", "Errors.Target.NotFound")
	}
	key.KeyID, err = c.idGenerator.Next()
	if err != nil {
		return time.Time{}, err
	}
	err = c.pushAppendAndReduce(ctx, writeModel, target.NewKeyAddedEvent(
		ctx,
		TargetAggregateFromWriteModel(&writeModel.WriteModel),
		key.KeyID,
		key.PublicKey,
		fingerprint,
		key.Expiration,
	))
	if err != nil {
		return time.Time{}, err
	}
	return writeModel.WriteModel.ChangeDate, nil
}

func (c *Commands) ActivateTargetPublicKey(ctx context.Context, targetID, keyID, resourceOwner string) (_ time.Time, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if targetID == "" || keyID == "" {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing")
	}
	writeModel, err := c.getTargetKeyWriteModelByID(ctx, targetID, keyID, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}
	if !writeModel.TargetExists || !writeModel.KeyExists {
		return time.Time{}, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SAF4g", "Errors.Target.NotFound")
	}
	if !writeModel.ExpirationDate.IsZero() && writeModel.ExpirationDate.Before(time.Now()) {
		return time.Time{}, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SAF4g", "Errors.Target.PublicKeyExpired")
	}
	if writeModel.Active {
		return writeModel.WriteModel.ChangeDate, nil
	}
	err = c.pushAppendAndReduce(ctx, writeModel, target.NewKeyActivatedEvent(
		ctx,
		TargetAggregateFromWriteModel(&writeModel.WriteModel),
		keyID,
	))
	if err != nil {
		return time.Time{}, err
	}
	return writeModel.WriteModel.ChangeDate, nil
}

func (c *Commands) DeactivateTargetPublicKey(ctx context.Context, targetID, keyID, resourceOwner string) (_ time.Time, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if targetID == "" || keyID == "" {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing")
	}
	writeModel, err := c.getTargetKeyWriteModelByID(ctx, targetID, keyID, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}
	if !writeModel.TargetExists || !writeModel.KeyExists {
		return time.Time{}, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SAF4g", "Errors.Target.NotFound")
	}
	if !writeModel.Active {
		return writeModel.WriteModel.ChangeDate, nil
	}
	err = c.pushAppendAndReduce(ctx, writeModel, target.NewKeyDeactivatedEvent(
		ctx,
		TargetAggregateFromWriteModel(&writeModel.WriteModel),
		keyID,
	))
	if err != nil {
		return time.Time{}, err
	}
	return writeModel.WriteModel.ChangeDate, nil
}

func (c *Commands) RemoveTargetPublicKey(ctx context.Context, targetID, keyID, resourceOwner string) (_ time.Time, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if targetID == "" || keyID == "" {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing")
	}
	writeModel, err := c.getTargetKeyWriteModelByID(ctx, targetID, keyID, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}
	if !writeModel.TargetExists || !writeModel.KeyExists {
		return writeModel.WriteModel.ChangeDate, nil
	}
	if writeModel.Active {
		return time.Time{}, zerrors.ThrowPreconditionFailed(nil, ErrPublicKeyDeleteActiveKey, "Errors.Target.PublicKeyActive")
	}
	err = c.pushAppendAndReduce(ctx, writeModel, target.NewKeyRemovedEvent(
		ctx,
		TargetAggregateFromWriteModel(&writeModel.WriteModel),
		keyID,
	))
	if err != nil {
		return time.Time{}, err
	}
	return writeModel.WriteModel.ChangeDate, nil
}

func checkPublicKeyAndComputeFingerprint(key []byte) (string, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return "", zerrors.ThrowInvalidArgument(nil, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", zerrors.ThrowInvalidArgument(err, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey")
	}
	switch pub.(type) {
	case *rsa.PublicKey,
		*ecdsa.PublicKey:
		fingerprint := sha256.Sum256(block.Bytes)
		return fmt.Sprintf("SHA256:%s", base64.RawStdEncoding.EncodeToString(fingerprint[:])), nil
	default:
		return "", zerrors.ThrowInvalidArgument(nil, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey")
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
