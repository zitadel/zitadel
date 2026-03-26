package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dbmock"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestGetOTPCryptoGeneratorConfigWithDefault(t *testing.T) {
	tests := []struct {
		name                        string
		instanceID                  string
		secretGeneratorSettingsRepo func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository
		defaultConfig               *crypto.GeneratorConfig
		otpType                     domain.OTPType
		want                        *crypto.GeneratorConfig
		wantErr                     error
	}{
		{
			name:       "no default config",
			instanceID: "instance-1",
			wantErr:    zerrors.ThrowInternal(nil, "DOM-3AcM0U", "missing default config"),
		},
		{
			name:       "failed to get secret generator settings",
			instanceID: "instance-1",
			secretGeneratorSettingsRepo: func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
				repo := domainmock.NewSecretGeneratorSettingsRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.UniqueCondition("instance-1", nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
							),
						),
					).Times(1).
					Return(nil, assert.AnError)
				return repo
			},
			defaultConfig: &crypto.GeneratorConfig{
				Length: 6,
			},
			otpType: domain.OTPTypeEmail,
			wantErr: zerrors.ThrowInternal(assert.AnError, "DOM-x7Yd3E", "failed fetching SecretGeneratorSettings"),
		},
		{
			name:                        "inactive secret generator settings",
			instanceID:                  "instance-1",
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStatePreview, domain.OTPTypeSMS),
			defaultConfig: &crypto.GeneratorConfig{
				Length: 6,
			},
			otpType: domain.OTPTypeSMS,
			want: &crypto.GeneratorConfig{
				Length: 6,
			},
		},
		{
			name:                        "active secret generator settings - type sms",
			instanceID:                  "instance-1",
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeSMS),
			defaultConfig: &crypto.GeneratorConfig{
				Length: 6,
			},
			otpType: domain.OTPTypeSMS,
			want: &crypto.GeneratorConfig{
				Expiry:              30 * time.Minute,
				Length:              6,
				IncludeLowerLetters: true,
				IncludeUpperLetters: false,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			},
		},
		{
			name:                        "active secret generator settings - type email",
			instanceID:                  "instance-1",
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			defaultConfig: &crypto.GeneratorConfig{
				Length: 6,
			},
			otpType: domain.OTPTypeEmail,
			want: &crypto.GeneratorConfig{
				Expiry:              30 * time.Minute,
				Length:              8,
				IncludeLowerLetters: true,
				IncludeUpperLetters: false,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			},
		},
		{
			name:       "active secret generator settings sms type with nil values - return default config",
			instanceID: "instance-1",
			secretGeneratorSettingsRepo: func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
				repo := domainmock.NewSecretGeneratorSettingsRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.UniqueCondition("instance-1", nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
							),
						),
					).Times(1).
					Return(&domain.SecretGeneratorSettings{
						Settings: domain.Settings{
							State: domain.SettingStateActive,
						},
					}, nil)
				return repo
			},
			defaultConfig: &crypto.GeneratorConfig{
				Length: 6,
			},
			otpType: domain.OTPTypeSMS,
			want: &crypto.GeneratorConfig{
				Length: 6,
			},
		},
		{
			name:       "active secret generator settings email type with nil values - return default config",
			instanceID: "instance-1",
			secretGeneratorSettingsRepo: func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
				repo := domainmock.NewSecretGeneratorSettingsRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.UniqueCondition("instance-1", nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
							),
						),
					).Times(1).
					Return(&domain.SecretGeneratorSettings{
						Settings: domain.Settings{
							State: domain.SettingStateActive,
						},
					}, nil)
				return repo
			},
			defaultConfig: &crypto.GeneratorConfig{
				Length: 6,
			},
			otpType: domain.OTPTypeEmail,
			want: &crypto.GeneratorConfig{
				Length: 6,
			},
		},
		{
			name:       "active secret generator settings with nil field values - return default config",
			instanceID: "instance-1",
			secretGeneratorSettingsRepo: func(ctrl *gomock.Controller) domain.SecretGeneratorSettingsRepository {
				repo := domainmock.NewSecretGeneratorSettingsRepo(ctrl)
				repo.EXPECT().
					Get(gomock.Any(), gomock.Any(),
						dbmock.QueryOptions(
							database.WithCondition(
								repo.UniqueCondition("instance-1", nil, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
							),
						),
					).Times(1).
					Return(&domain.SecretGeneratorSettings{
						Settings: domain.Settings{
							State: domain.SettingStateActive,
						},
						SecretGeneratorSettingsAttributes: domain.SecretGeneratorSettingsAttributes{
							OTPEmail: &domain.OTPEmailAttributes{
								SecretGeneratorAttrsWithExpiry: domain.SecretGeneratorAttrsWithExpiry{
									Expiry: gu.Ptr(30 * time.Minute),
									SecretGeneratorAttrs: domain.SecretGeneratorAttrs{
										Length:              nil, // should use value from the default config
										IncludeLowerLetters: nil, // should use value from the default config
										IncludeUpperLetters: gu.Ptr(false),
										IncludeDigits:       gu.Ptr(true),
										IncludeSymbols:      gu.Ptr(false),
									},
								},
							},
						},
					}, nil)
				return repo
			},
			defaultConfig: &crypto.GeneratorConfig{
				Length:              6,
				IncludeUpperLetters: true,
				IncludeLowerLetters: true,
			},
			otpType: domain.OTPTypeEmail,
			want: &crypto.GeneratorConfig{
				Expiry:              30 * time.Minute,
				Length:              6,    // from the default config
				IncludeLowerLetters: true, // from the default config
				IncludeUpperLetters: false,
				IncludeDigits:       true,
				IncludeSymbols:      false,
			},
		},
		{
			name:       "invalid otp type",
			instanceID: "instance-1",
			otpType:    3,
			defaultConfig: &crypto.GeneratorConfig{
				Length:              6,
				IncludeUpperLetters: true,
				IncludeLowerLetters: true,
			},
			secretGeneratorSettingsRepo: secretGeneratorSettingsRepo(domain.SettingStateActive, domain.OTPTypeEmail),
			wantErr:                     zerrors.ThrowInternal(nil, "DOM-3AcM0U", "invalid otp type"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			if tt.secretGeneratorSettingsRepo != nil {
				domain.WithSecretGeneratorSettingsRepo(tt.secretGeneratorSettingsRepo(ctrl))(opts)
			}
			got, err := domain.GetOTPCryptoGeneratorConfigWithDefault(context.Background(), tt.instanceID, opts, tt.defaultConfig, tt.otpType)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
