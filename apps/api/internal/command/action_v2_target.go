package command

import (
	"context"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AddTarget struct {
	models.ObjectRoot

	Name             string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool

	SigningKey string
}

func (a *AddTarget) IsValid() error {
	if a.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-ddqbm9us5p", "Errors.Target.Invalid")
	}
	if a.Timeout == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-39f35d8uri", "Errors.Target.NoTimeout")
	}
	_, err := url.Parse(a.Endpoint)
	if err != nil || a.Endpoint == "" {
		return zerrors.ThrowInvalidArgument(err, "COMMAND-1r2k6qo6wg", "Errors.Target.InvalidURL")
	}

	return nil
}

func (c *Commands) AddTarget(ctx context.Context, add *AddTarget, resourceOwner string) (_ time.Time, err error) {
	if resourceOwner == "" {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, "COMMAND-brml926e2d", "Errors.IDMissing")
	}

	if err := add.IsValid(); err != nil {
		return time.Time{}, err
	}

	if add.AggregateID == "" {
		add.AggregateID, err = c.idGenerator.Next()
		if err != nil {
			return time.Time{}, err
		}
	}
	wm, err := c.getTargetWriteModelByID(ctx, add.AggregateID, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}
	if wm.State.Exists() {
		return time.Time{}, zerrors.ThrowAlreadyExists(nil, "INSTANCE-9axkz0jvzm", "Errors.Target.AlreadyExists")
	}
	code, err := c.newSigningKey(ctx, c.eventstore.Filter, c.targetEncryption) //nolint
	if err != nil {
		return time.Time{}, err
	}
	add.SigningKey = code.PlainCode()
	pushedEvents, err := c.eventstore.Push(ctx, target.NewAddedEvent(
		ctx,
		TargetAggregateFromWriteModel(&wm.WriteModel),
		add.Name,
		add.TargetType,
		add.Endpoint,
		add.Timeout,
		add.InterruptOnError,
		code.Crypted,
	))
	if err != nil {
		return time.Time{}, err
	}
	if err := AppendAndReduce(wm, pushedEvents...); err != nil {
		return time.Time{}, err
	}
	return wm.ChangeDate, nil
}

type ChangeTarget struct {
	models.ObjectRoot

	Name             *string
	TargetType       *domain.TargetType
	Endpoint         *string
	Timeout          *time.Duration
	InterruptOnError *bool

	ExpirationSigningKey bool
	SigningKey           *string
}

func (a *ChangeTarget) IsValid() error {
	if a.AggregateID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-1l6ympeagp", "Errors.IDMissing")
	}
	if a.Name != nil && *a.Name == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-d1wx4lm0zr", "Errors.Target.Invalid")
	}
	if a.Timeout != nil && *a.Timeout == 0 {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-08b39vdi57", "Errors.Target.NoTimeout")
	}
	if a.Endpoint != nil {
		_, err := url.Parse(*a.Endpoint)
		if err != nil || *a.Endpoint == "" {
			return zerrors.ThrowInvalidArgument(err, "COMMAND-jsbaera7b6", "Errors.Target.InvalidURL")
		}
	}
	return nil
}

func (c *Commands) ChangeTarget(ctx context.Context, change *ChangeTarget, resourceOwner string) (time.Time, error) {
	if resourceOwner == "" {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, "COMMAND-zqibgg0wwh", "Errors.IDMissing")
	}
	if err := change.IsValid(); err != nil {
		return time.Time{}, err
	}
	existing, err := c.getTargetWriteModelByID(ctx, change.AggregateID, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}
	if !existing.State.Exists() {
		return time.Time{}, zerrors.ThrowNotFound(nil, "COMMAND-xj14f2cccn", "Errors.Target.NotFound")
	}

	var changedSigningKey *crypto.CryptoValue
	if change.ExpirationSigningKey {
		code, err := c.newSigningKey(ctx, c.eventstore.Filter, c.targetEncryption) //nolint
		if err != nil {
			return time.Time{}, err
		}
		changedSigningKey = code.Crypted
		change.SigningKey = &code.Plain
	}

	changedEvent := existing.NewChangedEvent(
		ctx,
		TargetAggregateFromWriteModel(&existing.WriteModel),
		change.Name,
		change.TargetType,
		change.Endpoint,
		change.Timeout,
		change.InterruptOnError,
		changedSigningKey,
	)
	if changedEvent == nil {
		return existing.WriteModel.ChangeDate, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return time.Time{}, err
	}
	err = AppendAndReduce(existing, pushedEvents...)
	if err != nil {
		return time.Time{}, err
	}
	return existing.WriteModel.ChangeDate, nil
}

func (c *Commands) DeleteTarget(ctx context.Context, id, resourceOwner string) (time.Time, error) {
	if id == "" || resourceOwner == "" {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, "COMMAND-obqos2l3no", "Errors.IDMissing")
	}

	existing, err := c.getTargetWriteModelByID(ctx, id, resourceOwner)
	if err != nil {
		return time.Time{}, err
	}
	if !existing.State.Exists() {
		return existing.WriteModel.ChangeDate, nil
	}

	if err := c.pushAppendAndReduce(ctx,
		existing,
		target.NewRemovedEvent(ctx,
			TargetAggregateFromWriteModel(&existing.WriteModel),
			existing.Name,
		),
	); err != nil {
		return time.Time{}, err
	}
	return existing.WriteModel.ChangeDate, nil
}

func (c *Commands) existsTargetsByIDs(ctx context.Context, ids []string, resourceOwner string) bool {
	wm := NewTargetsExistsWriteModel(ids, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return false
	}
	return wm.AllExists()
}

func (c *Commands) getTargetWriteModelByID(ctx context.Context, id string, resourceOwner string) (*TargetWriteModel, error) {
	wm := NewTargetWriteModel(id, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

func (c *Commands) newSigningKey(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*EncryptedCode, error) {
	return c.newEncryptedCodeWithDefault(ctx, filter, domain.SecretGeneratorTypeSigningKey, alg, c.defaultSecretGenerators.SigningKey)
}
