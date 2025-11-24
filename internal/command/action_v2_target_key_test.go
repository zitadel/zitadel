package command

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	rsaPublicKey = func() *rsa.PublicKey {
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		return &privateKey.PublicKey
	}()
	rsaPublicKeyBytes = func() []byte {
		data, _ := crypto.PublicKeyToBytes(rsaPublicKey)
		return data
	}()
	rsaFingerprint = func() string {
		block, _ := x509.MarshalPKIXPublicKey(rsaPublicKey)
		hash := sha256.Sum256(block)
		return fmt.Sprintf("SHA256:%s", base64.RawStdEncoding.EncodeToString(hash[:]))
	}()
)

func TestCommands_AddTargetPublicKey(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		key           *TargetPublicKey
		resourceOwner string
	}
	type res struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing target id",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				key: &TargetPublicKey{},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing"),
			},
		},
		{
			name: "expiration date before now",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				key: &TargetPublicKey{
					TargetID:   "target-1",
					Expiration: testNow.Add(-time.Hour),
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrExpirationDateBeforeNow, "Errors.Target.InvalidExpirationDate"),
			},
		},
		{
			name: "invalid public key",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				key: &TargetPublicKey{
					TargetID:  "target-1",
					PublicKey: []byte("invalid-key"),
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey"),
			},
		},
		{
			name: "invalid public key format (PKCS1)",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				key: &TargetPublicKey{
					TargetID:  "target-1",
					PublicKey: pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(rsaPublicKey)}),
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey"),
			},
		},
		{
			name: "invalid public key type (ed25519)",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
				key: &TargetPublicKey{
					TargetID: "target-1",
					PublicKey: func() []byte {
						publicKey, _, _ := ed25519.GenerateKey(rand.Reader)
						ed25519PublicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
						return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: ed25519PublicKeyBytes})
					}(),
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrPublicKeyFormat, "Errors.Target.InvalidPublicKey"),
			},
		},
		{
			name: "target not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				idGenerator: mock.ExpectID(t, "key-1"),
			},
			args: args{
				ctx: context.Background(),
				key: &TargetPublicKey{
					TargetID:  "target-1",
					PublicKey: rsaPublicKeyBytes,
				},
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-nd3fsd", "Errors.Target.NotFound"),
			},
		},
		{
			name: "success",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target-1", "instance1"),
						),
					),
					expectPush(
						eventFromEventPusher(
							target.NewKeyAddedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
								rsaPublicKeyBytes,
								rsaFingerprint,
								time.Time{},
							),
						),
					),
				),
				idGenerator: mock.ExpectID(t, "key-1"),
			},
			args: args{
				ctx: context.Background(),
				key: &TargetPublicKey{
					TargetID:  "target-1",
					PublicKey: rsaPublicKeyBytes,
				},
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore(t),
				idGenerator: tt.fields.idGenerator,
			}
			_, err := c.AddTargetPublicKey(tt.args.ctx, tt.args.key, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommands_ActivateTargetPublicKey(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		targetID      string
		keyID         string
		resourceOwner string
	}
	type res struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing target id",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing"),
			},
		},
		{
			name: "missing key id",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing"),
			},
		},
		{
			name: "target not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SAF4g", "Errors.Target.NotFound"),
			},
		},
		{
			name: "already expired",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target-1", "instance1"),
						),
						eventFromEventPusher(
							target.NewKeyAddedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
								rsaPublicKeyBytes,
								rsaFingerprint,
								time.Now().Add(-time.Hour),
							),
						),
					),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SAF4g", "Errors.Target.PublicKeyExpired"),
			},
		},
		{
			name: "already active",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target-1", "instance1"),
						),
						eventFromEventPusher(
							target.NewKeyAddedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
								rsaPublicKeyBytes,
								rsaFingerprint,
								time.Time{},
							),
						),
						eventFromEventPusher(
							target.NewKeyActivatedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
							),
						),
					),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{},
		},
		{
			name: "success",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target-1", "instance1"),
						),
						eventFromEventPusher(
							target.NewKeyAddedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
								rsaPublicKeyBytes,
								rsaFingerprint,
								time.Time{},
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							target.NewKeyActivatedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
							),
						),
					),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			_, err := c.ActivateTargetPublicKey(tt.args.ctx, tt.args.targetID, tt.args.keyID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommands_DeactivateTargetPublicKey(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		targetID      string
		keyID         string
		resourceOwner string
	}
	type res struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing target id",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing"),
			},
		},
		{
			name: "missing key id",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing"),
			},
		},
		{
			name: "target not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SAF4g", "Errors.Target.NotFound"),
			},
		},
		{
			name: "already inactive",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target-1", "instance1"),
						),
						eventFromEventPusher(
							target.NewKeyAddedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
								rsaPublicKeyBytes,
								rsaFingerprint,
								time.Time{},
							),
						),
					),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{},
		},
		{
			name: "success",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target-1", "instance1"),
						),
						eventFromEventPusher(
							target.NewKeyAddedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
								rsaPublicKeyBytes,
								rsaFingerprint,
								time.Time{},
							),
						),
						eventFromEventPusher(
							target.NewKeyActivatedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							target.NewKeyDeactivatedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
							),
						),
					),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			_, err := c.DeactivateTargetPublicKey(tt.args.ctx, tt.args.targetID, tt.args.keyID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommands_RemoveTargetPublicKey(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		targetID      string
		keyID         string
		resourceOwner string
	}
	type res struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing target id",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing"),
			},
		},
		{
			name: "missing key id",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, ErrTargetIDMissing, "Errors.IDMissing"),
			},
		},
		{
			name: "target not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{
				err: nil,
			},
		},
		{
			name: "key active",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target-1", "instance1"),
						),
						eventFromEventPusher(
							target.NewKeyAddedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
								rsaPublicKeyBytes,
								rsaFingerprint,
								time.Time{},
							),
						),
						eventFromEventPusher(
							target.NewKeyActivatedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
							),
						),
					),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, ErrPublicKeyDeleteActiveKey, "Errors.Target.PublicKeyActive"),
			},
		},
		{
			name: "success",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target-1", "instance1"),
						),
						eventFromEventPusher(
							target.NewKeyAddedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
								rsaPublicKeyBytes,
								rsaFingerprint,
								time.Time{},
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							target.NewKeyRemovedEvent(context.Background(),
								target.NewAggregate("target-1", "instance1"),
								"key-1",
							),
						),
					),
				),
			},
			args: args{
				ctx:      context.Background(),
				targetID: "target-1",
				keyID:    "key-1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			_, err := c.RemoveTargetPublicKey(tt.args.ctx, tt.args.targetID, tt.args.keyID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}
