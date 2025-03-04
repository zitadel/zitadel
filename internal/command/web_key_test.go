package command

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/webkey"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateWebKey(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)

	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		idGenerator     id.Generator
		webKeyGenerator func(keyID string, alg crypto.EncryptionAlgorithm, genConfig crypto.WebKeyConfig) (encryptedPrivate *crypto.CryptoValue, public *jose.JSONWebKey, err error)
	}
	type args struct {
		conf crypto.WebKeyConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *WebKeyDetails
		wantErr error
	}{
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				&crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "generate error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "key1"),
				webKeyGenerator: func(string, crypto.EncryptionAlgorithm, crypto.WebKeyConfig) (*crypto.CryptoValue, *jose.JSONWebKey, error) {
					return nil, nil, io.ErrClosedPipe
				},
			},
			args: args{
				&crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "generate key, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
					),
					expectPush(
						mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key2", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key2",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "key2"),
				webKeyGenerator: func(keyID string, _ crypto.EncryptionAlgorithm, _ crypto.WebKeyConfig) (*crypto.CryptoValue, *jose.JSONWebKey, error) {
					return &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						}, &jose.JSONWebKey{
							Key:       &key.PublicKey,
							KeyID:     keyID,
							Algorithm: string(jose.ES384),
							Use:       crypto.KeyUsageSigning.String(),
						}, nil
				},
			},
			args: args{
				&crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			},
			want: &WebKeyDetails{
				KeyID: "key2",
				ObjectDetails: &domain.ObjectDetails{
					ResourceOwner: "instance1",
					ID:            "key2",
				},
			},
		},
		{
			name: "generate and activate key, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						),
						webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "key1"),
				webKeyGenerator: func(keyID string, _ crypto.EncryptionAlgorithm, _ crypto.WebKeyConfig) (*crypto.CryptoValue, *jose.JSONWebKey, error) {
					return &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						}, &jose.JSONWebKey{
							Key:       &key.PublicKey,
							KeyID:     keyID,
							Algorithm: string(jose.ES384),
							Use:       crypto.KeyUsageSigning.String(),
						}, nil
				},
			},
			args: args{
				&crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			},
			want: &WebKeyDetails{
				KeyID: "key1",
				ObjectDetails: &domain.ObjectDetails{
					ResourceOwner: "instance1",
					ID:            "key1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				idGenerator:     tt.fields.idGenerator,
				webKeyGenerator: tt.fields.webKeyGenerator,
			}
			got, err := c.CreateWebKey(ctx, tt.args.conf)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_GenerateInitialWebKeys(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)

	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		idGenerator     id.Generator
		webKeyGenerator func(keyID string, alg crypto.EncryptionAlgorithm, genConfig crypto.WebKeyConfig) (encryptedPrivate *crypto.CryptoValue, public *jose.JSONWebKey, err error)
	}
	type args struct {
		conf crypto.WebKeyConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				&crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "key found, noop",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
					),
				),
			},
			args: args{
				&crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			},
			wantErr: nil,
		},
		{
			name: "id generator error",
			fields: fields{
				eventstore:  expectEventstore(expectFilter()),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrUnexpectedEOF),
			},
			args: args{
				&crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			},
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name: "keys generated and activated",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						),
						webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						),
						mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key2", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key2",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "key1", "key2"),
				webKeyGenerator: func(keyID string, _ crypto.EncryptionAlgorithm, _ crypto.WebKeyConfig) (*crypto.CryptoValue, *jose.JSONWebKey, error) {
					return &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "encKey",
							Crypted:    []byte("crypted"),
						}, &jose.JSONWebKey{
							Key:       &key.PublicKey,
							KeyID:     keyID,
							Algorithm: string(jose.ES384),
							Use:       crypto.KeyUsageSigning.String(),
						}, nil
				},
			},
			args: args{
				&crypto.WebKeyECDSAConfig{
					Curve: crypto.EllipticCurveP384,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				idGenerator:     tt.fields.idGenerator,
				webKeyGenerator: tt.fields.webKeyGenerator,
			}
			err := c.GenerateInitialWebKeys(ctx, tt.args.conf)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_ActivateWebKey(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)

	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		webKeyGenerator func(keyID string, alg crypto.EncryptionAlgorithm, genConfig crypto.WebKeyConfig) (encryptedPrivate *crypto.CryptoValue, public *jose.JSONWebKey, err error)
	}
	type args struct {
		keyID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr error
	}{
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args:    args{"key2"},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
					),
				),
			},
			args: args{"key1"},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
				ID:            "key1",
			},
		},
		{
			name: "not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
					),
				),
			},
			args:    args{"key2"},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-teiG3", "Errors.WebKey.NotFound"),
		},
		{
			name: "activate next, de-activate old, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key2", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key2",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
					),
					expectPush(
						webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key2", "instance1"),
						),
						webkey.NewDeactivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						),
					),
				),
			},
			args: args{"key2"},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
				ID:            "key2",
			},
		},
		{
			name: "activate next, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
					),
					expectPush(
						webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						),
					),
				),
			},
			args: args{"key1"},
			want: &domain.ObjectDetails{
				ResourceOwner: "instance1",
				ID:            "key1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				webKeyGenerator: tt.fields.webKeyGenerator,
			}
			got, err := c.ActivateWebKey(ctx, tt.args.keyID)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_DeleteWebKey(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		keyID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    time.Time
		wantErr error
	}{
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args:    args{"key1"},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{"key1"},
			want: time.Time{},
		},
		{
			name: "previously deleted",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key2", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key2",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key2", "instance1"),
						)),
						eventFromEventPusher(webkey.NewDeactivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
						eventFromEventPusher(webkey.NewRemovedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
					),
				),
			},
			args: args{"key1"},
			want: time.Time{},
		},
		{
			name: "key active error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
					),
				),
			},
			args:    args{"key1"},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Chai1", "Errors.WebKey.ActiveDelete"),
		},
		{
			name: "delete deactivated key",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key2", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key2",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewActivatedEvent(ctx,
							webkey.NewAggregate("key2", "instance1"),
						)),
						eventFromEventPusher(webkey.NewDeactivatedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
					),
					expectPush(
						webkey.NewRemovedEvent(ctx, webkey.NewAggregate("key1", "instance1")),
					),
				),
			},
			args: args{"key1"},
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.DeleteWebKey(ctx, tt.args.keyID)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func mustNewWebkeyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	privateKey *crypto.CryptoValue,
	publicKey *jose.JSONWebKey,
	config crypto.WebKeyConfig) *webkey.AddedEvent {
	event, err := webkey.NewAddedEvent(ctx, aggregate, privateKey, publicKey, config)
	if err != nil {
		panic(err)
	}
	return event
}
