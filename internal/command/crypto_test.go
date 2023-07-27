package command

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func mockCode(code string, exp time.Duration) cryptoCodeFunc {
	return func(ctx context.Context, filter preparation.FilterToQueryReducer, _ domain.SecretGeneratorType, alg crypto.Crypto) (*CryptoCode, error) {
		return &CryptoCode{
			Crypted: &crypto.CryptoValue{
				CryptoType: crypto.TypeEncryption,
				Algorithm:  "enc",
				KeyID:      "id",
				Crypted:    []byte(code),
			},
			Plain:  code,
			Expiry: exp,
		}, nil
	}
}

var (
	testGeneratorConfig = crypto.GeneratorConfig{
		Length:              12,
		Expiry:              60000000000,
		IncludeLowerLetters: true,
		IncludeUpperLetters: true,
		IncludeDigits:       true,
		IncludeSymbols:      true,
	}
)

func testSecretGeneratorAddedEvent(typ domain.SecretGeneratorType) *instance.SecretGeneratorAddedEvent {
	return instance.NewSecretGeneratorAddedEvent(context.Background(),
		&instance.NewAggregate("inst1").Aggregate, typ,
		testGeneratorConfig.Length,
		testGeneratorConfig.Expiry,
		testGeneratorConfig.IncludeLowerLetters,
		testGeneratorConfig.IncludeUpperLetters,
		testGeneratorConfig.IncludeDigits,
		testGeneratorConfig.IncludeSymbols,
	)
}

func Test_newCryptoCode(t *testing.T) {
	type args struct {
		typ domain.SecretGeneratorType
		alg crypto.Crypto
	}
	tests := []struct {
		name       string
		eventstore *eventstore.Eventstore
		args       args
		wantErr    error
	}{
		{
			name:       "filter config error",
			eventstore: eventstoreExpect(t, expectFilterError(io.ErrClosedPipe)),
			args: args{
				typ: domain.SecretGeneratorTypeVerifyEmailCode,
				alg: crypto.CreateMockHashAlg(gomock.NewController(t)),
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			eventstore: eventstoreExpect(t, expectFilter(
				eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypeVerifyEmailCode)),
			)),
			args: args{
				typ: domain.SecretGeneratorTypeVerifyEmailCode,
				alg: crypto.CreateMockHashAlg(gomock.NewController(t)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newCryptoCode(context.Background(), tt.eventstore.Filter, tt.args.typ, tt.args.alg)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				require.NotNil(t, got)
				assert.NotNil(t, got.Crypted)
				assert.NotEmpty(t, got)
				assert.Equal(t, testGeneratorConfig.Expiry, got.Expiry)
			}
		})
	}
}

func Test_verifyCryptoCode(t *testing.T) {
	es := eventstoreExpect(t, expectFilter(
		eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypeVerifyEmailCode)),
	))
	code, err := newCryptoCode(context.Background(), es.Filter, domain.SecretGeneratorTypeVerifyEmailCode, crypto.CreateMockHashAlg(gomock.NewController(t)))
	require.NoError(t, err)

	type args struct {
		typ     domain.SecretGeneratorType
		alg     crypto.Crypto
		expiry  time.Duration
		crypted *crypto.CryptoValue
		plain   string
	}
	tests := []struct {
		name      string
		eventsore *eventstore.Eventstore
		args      args
		wantErr   bool
	}{
		{
			name:      "filter config error",
			eventsore: eventstoreExpect(t, expectFilterError(io.ErrClosedPipe)),
			args: args{
				typ:     domain.SecretGeneratorTypeVerifyEmailCode,
				alg:     crypto.CreateMockHashAlg(gomock.NewController(t)),
				expiry:  code.Expiry,
				crypted: code.Crypted,
				plain:   code.Plain,
			},
			wantErr: true,
		},
		{
			name: "success",
			eventsore: eventstoreExpect(t, expectFilter(
				eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypeVerifyEmailCode)),
			)),
			args: args{
				typ:     domain.SecretGeneratorTypeVerifyEmailCode,
				alg:     crypto.CreateMockHashAlg(gomock.NewController(t)),
				expiry:  code.Expiry,
				crypted: code.Crypted,
				plain:   code.Plain,
			},
		},
		{
			name: "wrong plain",
			eventsore: eventstoreExpect(t, expectFilter(
				eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypeVerifyEmailCode)),
			)),
			args: args{
				typ:     domain.SecretGeneratorTypeVerifyEmailCode,
				alg:     crypto.CreateMockHashAlg(gomock.NewController(t)),
				expiry:  code.Expiry,
				crypted: code.Crypted,
				plain:   "wrong",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyCryptoCode(context.Background(), tt.eventsore.Filter, tt.args.typ, tt.args.alg, time.Now(), tt.args.expiry, tt.args.crypted, tt.args.plain)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func Test_secretGenerator(t *testing.T) {
	type args struct {
		typ domain.SecretGeneratorType
		alg crypto.Crypto
	}
	tests := []struct {
		name      string
		eventsore *eventstore.Eventstore
		args      args
		want      crypto.Generator
		wantConf  *crypto.GeneratorConfig
		wantErr   error
	}{
		{
			name:      "filter config error",
			eventsore: eventstoreExpect(t, expectFilterError(io.ErrClosedPipe)),
			args: args{
				typ: domain.SecretGeneratorTypeVerifyEmailCode,
				alg: crypto.CreateMockHashAlg(gomock.NewController(t)),
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "hash generator",
			eventsore: eventstoreExpect(t, expectFilter(
				eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypeVerifyEmailCode)),
			)),
			args: args{
				typ: domain.SecretGeneratorTypeVerifyEmailCode,
				alg: crypto.CreateMockHashAlg(gomock.NewController(t)),
			},
			want:     crypto.NewHashGenerator(testGeneratorConfig, crypto.CreateMockHashAlg(gomock.NewController(t))),
			wantConf: &testGeneratorConfig,
		},
		{
			name: "encryption generator",
			eventsore: eventstoreExpect(t, expectFilter(
				eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypeVerifyEmailCode)),
			)),
			args: args{
				typ: domain.SecretGeneratorTypeVerifyEmailCode,
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			want:     crypto.NewEncryptionGenerator(testGeneratorConfig, crypto.CreateMockEncryptionAlg(gomock.NewController(t))),
			wantConf: &testGeneratorConfig,
		},
		{
			name: "unsupported type",
			eventsore: eventstoreExpect(t, expectFilter(
				eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypeVerifyEmailCode)),
			)),
			args: args{
				typ: domain.SecretGeneratorTypeVerifyEmailCode,
				alg: nil,
			},
			wantErr: errors.ThrowInternalf(nil, "COMMA-RreV6", "Errors.Internal unsupported crypto algorithm type %T", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotConf, err := secretGenerator(context.Background(), tt.eventsore.Filter, tt.args.typ, tt.args.alg)
			require.ErrorIs(t, err, tt.wantErr)
			assert.IsType(t, tt.want, got)
			assert.Equal(t, tt.wantConf, gotConf)
		})
	}
}
