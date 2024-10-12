package command

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/bcrypt"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func mockEncryptedCode(code string, exp time.Duration) encrypedCodeFunc {
	return func(ctx context.Context, filter preparation.FilterToQueryReducer, _ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm) (*EncryptedCode, error) {
		return &EncryptedCode{
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

func mockEncryptedCodeWithDefault(code string, exp time.Duration) encryptedCodeWithDefaultFunc {
	return func(ctx context.Context, filter preparation.FilterToQueryReducer, _ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, _ *crypto.GeneratorConfig) (*EncryptedCode, error) {
		return &EncryptedCode{
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

func mockEncryptedCodeGeneratorWithDefault(code string, exp time.Duration) encryptedCodeGeneratorWithDefaultFunc {
	return func(ctx context.Context, filter preparation.FilterToQueryReducer, _ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, _ *crypto.GeneratorConfig) (*EncryptedCode, string, error) {
		return &EncryptedCode{
			Crypted: &crypto.CryptoValue{
				CryptoType: crypto.TypeEncryption,
				Algorithm:  "enc",
				KeyID:      "id",
				Crypted:    []byte(code),
			},
			Plain:  code,
			Expiry: exp,
		}, "", nil
	}
}

func mockEncryptedCodeGeneratorWithDefaultExternal(id string) encryptedCodeGeneratorWithDefaultFunc {
	return func(ctx context.Context, filter preparation.FilterToQueryReducer, _ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, _ *crypto.GeneratorConfig) (*EncryptedCode, string, error) {
		return nil, id, nil
	}
}

func mockHashedSecret(secret string) hashedSecretFunc {
	return func(_ context.Context, _ preparation.FilterToQueryReducer) (encodedHash string, plain string, err error) {
		return secret, secret, nil
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
		alg crypto.EncryptionAlgorithm
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
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
				alg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newEncryptedCode(context.Background(), tt.eventstore.Filter, tt.args.typ, tt.args.alg) //nolint:staticcheck
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

func Test_cryptoCodeGenerator(t *testing.T) {
	type args struct {
		typ           domain.SecretGeneratorType
		alg           crypto.EncryptionAlgorithm
		defaultConfig *crypto.GeneratorConfig
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
				typ:           domain.SecretGeneratorTypeVerifyEmailCode,
				alg:           crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				defaultConfig: emptyConfig,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "encryption generator",
			eventsore: eventstoreExpect(t, expectFilter(
				eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypeVerifyEmailCode)),
			)),
			args: args{
				typ:           domain.SecretGeneratorTypeVerifyEmailCode,
				alg:           crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				defaultConfig: emptyConfig,
			},
			want:     crypto.NewEncryptionGenerator(testGeneratorConfig, crypto.CreateMockEncryptionAlg(gomock.NewController(t))),
			wantConf: &testGeneratorConfig,
		},
		{
			name:      "encryption generator with default config",
			eventsore: eventstoreExpect(t, expectFilter()),
			args: args{
				typ:           domain.SecretGeneratorTypeVerifyEmailCode,
				alg:           crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				defaultConfig: &testGeneratorConfig,
			},
			want:     crypto.NewEncryptionGenerator(testGeneratorConfig, crypto.CreateMockEncryptionAlg(gomock.NewController(t))),
			wantConf: &testGeneratorConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotConf, err := encryptedCodeGenerator(context.Background(), tt.eventsore.Filter, tt.args.typ, tt.args.alg, tt.args.defaultConfig) //nolint:staticcheck
			require.ErrorIs(t, err, tt.wantErr)
			assert.IsType(t, tt.want, got)
			assert.Equal(t, tt.wantConf, gotConf)
		})
	}
}

func Test_newHashedSecretWithDefault(t *testing.T) {
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		wantLen    int
		wantErr    bool
	}{
		{
			name: "filter error",
			eventstore: expectEventstore(
				expectFilterError(io.ErrClosedPipe),
			),
			wantErr: true,
		},
		{
			name: "default config",
			eventstore: expectEventstore(
				expectFilter(),
			),
			wantLen: 32,
		},
		{
			name: "instance config",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(
						instance.NewSecretGeneratorAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							domain.SecretGeneratorTypeAppSecret,
							24,
							time.Hour*1,
							true,
							true,
							true,
							true,
						),
					),
				),
			),
			wantLen: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := &crypto.Hasher{
				Swapper: passwap.NewSwapper(bcrypt.New(bcrypt.MinCost)),
			}
			defaultConfig := &crypto.GeneratorConfig{
				Length:              32,
				Expiry:              time.Minute,
				IncludeLowerLetters: true,
			}
			generate := newHashedSecretWithDefault(hasher, defaultConfig)
			encodedHash, plain, err := generate(context.Background(), tt.eventstore(t).Filter) //nolint:staticcheck
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, plain, tt.wantLen)
			_, err = hasher.Verify(encodedHash, plain)
			require.NoError(t, err)
		})
	}
}
