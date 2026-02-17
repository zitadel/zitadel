package repository_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

var now = time.Now().Round(time.Millisecond)

func TestGetLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LoginSettings()

	settings := []*domain.LoginSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			AllowUsernamePassword:       true,
			AllowRegister:               true,
			AllowExternalIDP:            true,
			ForceMultiFactor:            true,
			ForceMultiFactorLocalOnly:   true,
			HidePasswordReset:           true,
			IgnoreUnknownUsernames:      true,
			AllowDomainDiscovery:        true,
			DisableLoginWithEmail:       true,
			DisableLoginWithPhone:       true,
			PasswordlessType:            domain.PasswordlessTypeAllowed,
			DefaultRedirectURI:          "uri",
			PasswordCheckLifetime:       time.Hour,
			ExternalLoginCheckLifetime:  time.Hour,
			MultiFactorInitSkipLifetime: time.Hour,
			SecondFactorCheckLifetime:   time.Hour,
			MultiFactorCheckLifetime:    time.Hour,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			AllowUsernamePassword:       false,
			AllowRegister:               false,
			AllowExternalIDP:            false,
			ForceMultiFactor:            false,
			ForceMultiFactorLocalOnly:   false,
			HidePasswordReset:           false,
			IgnoreUnknownUsernames:      false,
			AllowDomainDiscovery:        false,
			DisableLoginWithEmail:       false,
			DisableLoginWithPhone:       false,
			PasswordlessType:            domain.PasswordlessTypeNotAllowed,
			DefaultRedirectURI:          "uri2",
			PasswordCheckLifetime:       time.Minute,
			ExternalLoginCheckLifetime:  time.Minute,
			MultiFactorInitSkipLifetime: time.Minute,
			SecondFactorCheckLifetime:   time.Minute,
			MultiFactorCheckLifetime:    time.Minute,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: &firstOrgID,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			AllowUsernamePassword:       true,
			AllowRegister:               true,
			AllowExternalIDP:            true,
			ForceMultiFactor:            true,
			ForceMultiFactorLocalOnly:   true,
			HidePasswordReset:           true,
			IgnoreUnknownUsernames:      true,
			AllowDomainDiscovery:        true,
			DisableLoginWithEmail:       true,
			DisableLoginWithPhone:       true,
			PasswordlessType:            domain.PasswordlessTypeAllowed,
			DefaultRedirectURI:          "uri",
			PasswordCheckLifetime:       time.Hour,
			ExternalLoginCheckLifetime:  time.Hour,
			MultiFactorInitSkipLifetime: time.Hour,
			SecondFactorCheckLifetime:   time.Hour,
			MultiFactorCheckLifetime:    time.Hour,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: &secondOrgID,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			AllowUsernamePassword:       false,
			AllowRegister:               false,
			AllowExternalIDP:            false,
			ForceMultiFactor:            false,
			ForceMultiFactorLocalOnly:   false,
			HidePasswordReset:           false,
			IgnoreUnknownUsernames:      false,
			AllowDomainDiscovery:        false,
			DisableLoginWithEmail:       false,
			DisableLoginWithPhone:       false,
			PasswordlessType:            domain.PasswordlessTypeNotAllowed,
			DefaultRedirectURI:          "uri2",
			PasswordCheckLifetime:       time.Minute,
			ExternalLoginCheckLifetime:  time.Minute,
			MultiFactorInitSkipLifetime: time.Minute,
			SecondFactorCheckLifetime:   time.Minute,
			MultiFactorCheckLifetime:    time.Minute,
		},
	}

	for i, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err, "failed to create %d. setting", i)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.LoginSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(nil)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(secondInstanceID), repo.OrganizationIDCondition(&secondOrgID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LoginSettings()

	settings := []*domain.LoginSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			AllowUsernamePassword:       true,
			AllowRegister:               true,
			AllowExternalIDP:            true,
			ForceMultiFactor:            true,
			ForceMultiFactorLocalOnly:   true,
			HidePasswordReset:           true,
			IgnoreUnknownUsernames:      true,
			AllowDomainDiscovery:        true,
			DisableLoginWithEmail:       true,
			DisableLoginWithPhone:       true,
			PasswordlessType:            domain.PasswordlessTypeAllowed,
			DefaultRedirectURI:          "uri",
			PasswordCheckLifetime:       time.Hour,
			ExternalLoginCheckLifetime:  time.Hour,
			MultiFactorInitSkipLifetime: time.Hour,
			SecondFactorCheckLifetime:   time.Hour,
			MultiFactorCheckLifetime:    time.Hour,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			AllowUsernamePassword:       false,
			AllowRegister:               false,
			AllowExternalIDP:            false,
			ForceMultiFactor:            false,
			ForceMultiFactorLocalOnly:   false,
			HidePasswordReset:           false,
			IgnoreUnknownUsernames:      false,
			AllowDomainDiscovery:        false,
			DisableLoginWithEmail:       false,
			DisableLoginWithPhone:       false,
			PasswordlessType:            domain.PasswordlessTypeNotAllowed,
			DefaultRedirectURI:          "uri2",
			PasswordCheckLifetime:       time.Minute,
			ExternalLoginCheckLifetime:  time.Minute,
			MultiFactorInitSkipLifetime: time.Minute,
			SecondFactorCheckLifetime:   time.Minute,
			MultiFactorCheckLifetime:    time.Minute,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			AllowUsernamePassword:       true,
			AllowRegister:               true,
			AllowExternalIDP:            true,
			ForceMultiFactor:            true,
			ForceMultiFactorLocalOnly:   true,
			HidePasswordReset:           true,
			IgnoreUnknownUsernames:      true,
			AllowDomainDiscovery:        true,
			DisableLoginWithEmail:       true,
			DisableLoginWithPhone:       true,
			PasswordlessType:            domain.PasswordlessTypeAllowed,
			DefaultRedirectURI:          "uri",
			PasswordCheckLifetime:       time.Hour,
			ExternalLoginCheckLifetime:  time.Hour,
			MultiFactorInitSkipLifetime: time.Hour,
			SecondFactorCheckLifetime:   time.Hour,
			MultiFactorCheckLifetime:    time.Hour,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			AllowUsernamePassword:       false,
			AllowRegister:               false,
			AllowExternalIDP:            false,
			ForceMultiFactor:            false,
			ForceMultiFactorLocalOnly:   false,
			HidePasswordReset:           false,
			IgnoreUnknownUsernames:      false,
			AllowDomainDiscovery:        false,
			DisableLoginWithEmail:       false,
			DisableLoginWithPhone:       false,
			PasswordlessType:            domain.PasswordlessTypeNotAllowed,
			DefaultRedirectURI:          "uri2",
			PasswordCheckLifetime:       time.Minute,
			ExternalLoginCheckLifetime:  time.Minute,
			MultiFactorInitSkipLifetime: time.Minute,
			SecondFactorCheckLifetime:   time.Minute,
			MultiFactorCheckLifetime:    time.Minute,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.LoginSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.LoginSetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.LoginSetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.LoginSetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.LoginSetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LoginSettings()

	existingSettings := &domain.LoginSetting{
		Setting: domain.Setting{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
			ID:             gofakeit.UUID(),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		AllowUsernamePassword:       false,
		AllowRegister:               false,
		AllowExternalIDP:            false,
		ForceMultiFactor:            false,
		ForceMultiFactorLocalOnly:   false,
		HidePasswordReset:           false,
		IgnoreUnknownUsernames:      false,
		AllowDomainDiscovery:        false,
		DisableLoginWithEmail:       false,
		DisableLoginWithPhone:       false,
		PasswordlessType:            domain.PasswordlessTypeNotAllowed,
		DefaultRedirectURI:          "uri2",
		PasswordCheckLifetime:       time.Minute,
		ExternalLoginCheckLifetime:  time.Minute,
		MultiFactorInitSkipLifetime: time.Minute,
		SecondFactorCheckLifetime:   time.Minute,
		MultiFactorCheckLifetime:    time.Minute,
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.LoginSetting
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.LoginSetting{
				Setting: domain.Setting{
					InstanceID:     instanceID,
					OrganizationID: nil,
					ID:             gofakeit.UUID(),
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				AllowUsernamePassword:       true,
				AllowRegister:               true,
				AllowExternalIDP:            true,
				ForceMultiFactor:            true,
				ForceMultiFactorLocalOnly:   true,
				HidePasswordReset:           true,
				IgnoreUnknownUsernames:      true,
				AllowDomainDiscovery:        true,
				DisableLoginWithEmail:       true,
				DisableLoginWithPhone:       true,
				PasswordlessType:            domain.PasswordlessTypeAllowed,
				DefaultRedirectURI:          "uri",
				PasswordCheckLifetime:       time.Hour,
				ExternalLoginCheckLifetime:  time.Hour,
				MultiFactorInitSkipLifetime: time.Hour,
				SecondFactorCheckLifetime:   time.Hour,
				MultiFactorCheckLifetime:    time.Hour,
			},
		},
		{
			name: "update organization",
			settings: &domain.LoginSetting{
				Setting: domain.Setting{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
					ID:             gofakeit.UUID(),
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				AllowUsernamePassword:       true,
				AllowRegister:               true,
				AllowExternalIDP:            true,
				ForceMultiFactor:            true,
				ForceMultiFactorLocalOnly:   true,
				HidePasswordReset:           true,
				IgnoreUnknownUsernames:      true,
				AllowDomainDiscovery:        true,
				DisableLoginWithEmail:       true,
				DisableLoginWithPhone:       true,
				PasswordlessType:            domain.PasswordlessTypeAllowed,
				DefaultRedirectURI:          "uri",
				PasswordCheckLifetime:       time.Hour,
				ExternalLoginCheckLifetime:  time.Hour,
				MultiFactorInitSkipLifetime: time.Hour,
				SecondFactorCheckLifetime:   time.Hour,
				MultiFactorCheckLifetime:    time.Hour,
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.LoginSetting{
				Setting: domain.Setting{
					InstanceID:     "foo",
					OrganizationID: nil,
					ID:             gofakeit.UUID(),
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				AllowUsernamePassword:       true,
				AllowRegister:               true,
				AllowExternalIDP:            true,
				ForceMultiFactor:            true,
				ForceMultiFactorLocalOnly:   true,
				HidePasswordReset:           true,
				IgnoreUnknownUsernames:      true,
				AllowDomainDiscovery:        true,
				DisableLoginWithEmail:       true,
				DisableLoginWithPhone:       true,
				PasswordlessType:            domain.PasswordlessTypeAllowed,
				DefaultRedirectURI:          "uri",
				PasswordCheckLifetime:       time.Hour,
				ExternalLoginCheckLifetime:  time.Hour,
				MultiFactorInitSkipLifetime: time.Hour,
				SecondFactorCheckLifetime:   time.Hour,
				MultiFactorCheckLifetime:    time.Hour,
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.LoginSetting{
				Setting: domain.Setting{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
					ID:             gofakeit.UUID(),
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				AllowUsernamePassword:       true,
				AllowRegister:               true,
				AllowExternalIDP:            true,
				ForceMultiFactor:            true,
				ForceMultiFactorLocalOnly:   true,
				HidePasswordReset:           true,
				IgnoreUnknownUsernames:      true,
				AllowDomainDiscovery:        true,
				DisableLoginWithEmail:       true,
				DisableLoginWithPhone:       true,
				PasswordlessType:            domain.PasswordlessTypeAllowed,
				DefaultRedirectURI:          "uri",
				PasswordCheckLifetime:       time.Hour,
				ExternalLoginCheckLifetime:  time.Hour,
				MultiFactorInitSkipLifetime: time.Hour,
				SecondFactorCheckLifetime:   time.Hour,
				MultiFactorCheckLifetime:    time.Hour,
			},
			wantErr: new(database.ForeignKeyError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			t.Cleanup(rollback)
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LoginSettings()

	existingInstanceSettings := &domain.LoginSetting{
		Setting: domain.Setting{
			InstanceID:     instanceID,
			OrganizationID: nil,
			ID:             gofakeit.UUID(),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		AllowUsernamePassword:       false,
		AllowRegister:               false,
		AllowExternalIDP:            false,
		ForceMultiFactor:            false,
		ForceMultiFactorLocalOnly:   false,
		HidePasswordReset:           false,
		IgnoreUnknownUsernames:      false,
		AllowDomainDiscovery:        false,
		DisableLoginWithEmail:       false,
		DisableLoginWithPhone:       false,
		PasswordlessType:            domain.PasswordlessTypeNotAllowed,
		DefaultRedirectURI:          "uri",
		PasswordCheckLifetime:       time.Minute,
		ExternalLoginCheckLifetime:  time.Minute,
		MultiFactorInitSkipLifetime: time.Minute,
		SecondFactorCheckLifetime:   time.Minute,
		MultiFactorCheckLifetime:    time.Minute,
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.LoginSetting{
		Setting: domain.Setting{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
			ID:             gofakeit.UUID(),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		AllowUsernamePassword:       false,
		AllowRegister:               false,
		AllowExternalIDP:            false,
		ForceMultiFactor:            false,
		ForceMultiFactorLocalOnly:   false,
		HidePasswordReset:           false,
		IgnoreUnknownUsernames:      false,
		AllowDomainDiscovery:        false,
		DisableLoginWithEmail:       false,
		DisableLoginWithPhone:       false,
		PasswordlessType:            domain.PasswordlessTypeNotAllowed,
		DefaultRedirectURI:          "uri",
		PasswordCheckLifetime:       time.Minute,
		ExternalLoginCheckLifetime:  time.Minute,
		MultiFactorInitSkipLifetime: time.Minute,
		SecondFactorCheckLifetime:   time.Minute,
		MultiFactorCheckLifetime:    time.Minute,
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.OrganizationIDCondition(gu.Ptr(orgID)),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(instanceID, gu.Ptr("foo"), domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeLogin, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.BrandingSettings()

	settings := []*domain.BrandingSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PrimaryColorLight:    "color",
			BackgroundColorLight: "color",
			WarnColorLight:       "color",
			FontColorLight:       "color",
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     "color",
			BackgroundColorDark:  "color",
			WarnColorDark:        "color",
			FontColorDark:        "color",
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  true,
			ErrorMessagePopup:    true,
			DisableWatermark:     true,
			ThemeMode:            domain.BrandingPolicyThemeAuto,
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PrimaryColorLight:    "color",
			BackgroundColorLight: "color",
			WarnColorLight:       "color",
			FontColorLight:       "color",
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     "color",
			BackgroundColorDark:  "color",
			WarnColorDark:        "color",
			FontColorDark:        "color",
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  true,
			ErrorMessagePopup:    true,
			DisableWatermark:     true,
			ThemeMode:            domain.BrandingPolicyThemeAuto,
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PrimaryColorLight:    "color",
			BackgroundColorLight: "color",
			WarnColorLight:       "color",
			FontColorLight:       "color",
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     "color",
			BackgroundColorDark:  "color",
			WarnColorDark:        "color",
			FontColorDark:        "color",
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  true,
			ErrorMessagePopup:    true,
			DisableWatermark:     true,
			ThemeMode:            domain.BrandingPolicyThemeAuto,
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PrimaryColorLight:    "color",
			BackgroundColorLight: "color",
			WarnColorLight:       "color",
			FontColorLight:       "color",
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     "color",
			BackgroundColorDark:  "color",
			WarnColorDark:        "color",
			FontColorDark:        "color",
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  true,
			ErrorMessagePopup:    true,
			DisableWatermark:     true,
			ThemeMode:            domain.BrandingPolicyThemeAuto,
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}
	_, err := repo.ActivateAt(t.Context(), tx, repo.UniqueCondition(settings[0].InstanceID, settings[0].OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview), now)
	require.NoError(t, err)
	activeInstanceSetting := *settings[0]
	activeInstanceSetting.State = domain.SettingStateActive

	_, err = repo.ActivateAt(t.Context(), tx, repo.UniqueCondition(settings[3].InstanceID, settings[3].OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview), now)
	require.NoError(t, err)
	activeOrganizationSetting := *settings[3]
	activeOrganizationSetting.State = domain.SettingStateActive

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.BrandingSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(&firstOrgID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name: "ok, instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: settings[0],
		},
		{
			name: "ok, instance, active",
			condition: database.And(
				repo.InstanceIDCondition(activeInstanceSetting.InstanceID),
				repo.OrganizationIDCondition(nil),
				repo.StateCondition(domain.SettingStateActive),
			),
			want: &activeInstanceSetting,
		},
		{
			name: "ok, organization",
			condition: database.And(
				repo.InstanceIDCondition(settings[3].InstanceID),
				repo.OrganizationIDCondition(settings[3].OrganizationID),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: settings[3],
		},
		{
			name: "ok, organization, active",
			condition: database.And(
				repo.InstanceIDCondition(activeOrganizationSetting.InstanceID),
				repo.OrganizationIDCondition(activeOrganizationSetting.OrganizationID),
				repo.StateCondition(domain.SettingStateActive),
			),
			want: &activeOrganizationSetting,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.BrandingSettings()

	settings := []*domain.BrandingSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PrimaryColorLight:    "color",
			BackgroundColorLight: "color",
			WarnColorLight:       "color",
			FontColorLight:       "color",
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     "color",
			BackgroundColorDark:  "color",
			WarnColorDark:        "color",
			FontColorDark:        "color",
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  true,
			ErrorMessagePopup:    true,
			DisableWatermark:     true,
			ThemeMode:            domain.BrandingPolicyThemeAuto,
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PrimaryColorLight:    "color",
			BackgroundColorLight: "color",
			WarnColorLight:       "color",
			FontColorLight:       "color",
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     "color",
			BackgroundColorDark:  "color",
			WarnColorDark:        "color",
			FontColorDark:        "color",
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  true,
			ErrorMessagePopup:    true,
			DisableWatermark:     true,
			ThemeMode:            domain.BrandingPolicyThemeAuto,
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PrimaryColorLight:    "color",
			BackgroundColorLight: "color",
			WarnColorLight:       "color",
			FontColorLight:       "color",
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     "color",
			BackgroundColorDark:  "color",
			WarnColorDark:        "color",
			FontColorDark:        "color",
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  true,
			ErrorMessagePopup:    true,
			DisableWatermark:     true,
			ThemeMode:            domain.BrandingPolicyThemeAuto,
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PrimaryColorLight:    "color",
			BackgroundColorLight: "color",
			WarnColorLight:       "color",
			FontColorLight:       "color",
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     "color",
			BackgroundColorDark:  "color",
			WarnColorDark:        "color",
			FontColorDark:        "color",
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  true,
			ErrorMessagePopup:    true,
			DisableWatermark:     true,
			ThemeMode:            domain.BrandingPolicyThemeAuto,
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}
	_, err := repo.ActivateAt(t.Context(), tx, repo.UniqueCondition(settings[0].InstanceID, settings[0].OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview), now)
	require.NoError(t, err)
	activeInstanceSetting := *settings[0]
	activeInstanceSetting.State = domain.SettingStateActive

	_, err = repo.ActivateAt(t.Context(), tx, repo.UniqueCondition(settings[2].InstanceID, settings[2].OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview), now)
	require.NoError(t, err)
	activeOrganizationSetting := *settings[2]
	activeOrganizationSetting.State = domain.SettingStateActive

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.BrandingSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.BrandingSetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.BrandingSetting{&activeOrganizationSetting, settings[2], &activeInstanceSetting, settings[0]},
		},
		{
			name: "all from instance, preview",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: []*domain.BrandingSetting{settings[2], settings[0]},
		},
		{
			name: "all from instance, active",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.StateCondition(domain.SettingStateActive),
			),
			want: []*domain.BrandingSetting{&activeOrganizationSetting, &activeInstanceSetting},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.BrandingSetting{&activeInstanceSetting, settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.BrandingSetting{&activeOrganizationSetting, settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn(), repo.StateColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)

			if !assert.Len(t, got, len(tt.want)) {
				return
			}
			for i := range got {
				assert.EqualExportedValues(t, tt.want[i], got[i])
			}
		})
	}
}

func TestSetBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettings()

	existingSettings := &domain.BrandingSetting{
		Setting: domain.Setting{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
			ID:             gofakeit.UUID(),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		PrimaryColorLight:    "color",
		BackgroundColorLight: "color",
		WarnColorLight:       "color",
		FontColorLight:       "color",
		LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
		IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
		PrimaryColorDark:     "color",
		BackgroundColorDark:  "color",
		WarnColorDark:        "color",
		FontColorDark:        "color",
		LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
		IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
		HideLoginNameSuffix:  true,
		ErrorMessagePopup:    true,
		DisableWatermark:     true,
		ThemeMode:            domain.BrandingPolicyThemeAuto,
		FontURL:              &url.URL{Scheme: "https", Host: "host"},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.BrandingSetting
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.BrandingSetting{
				Setting: domain.Setting{
					InstanceID:     instanceID,
					OrganizationID: nil,
					ID:             gofakeit.UUID(),
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				PrimaryColorLight:    "color",
				BackgroundColorLight: "color",
				WarnColorLight:       "color",
				FontColorLight:       "color",
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     "color",
				BackgroundColorDark:  "color",
				WarnColorDark:        "color",
				FontColorDark:        "color",
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  true,
				ErrorMessagePopup:    true,
				DisableWatermark:     true,
				ThemeMode:            domain.BrandingPolicyThemeAuto,
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			name: "update organization",
			settings: &domain.BrandingSetting{
				Setting: domain.Setting{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
					ID:             gofakeit.UUID(),
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				PrimaryColorLight:    "color",
				BackgroundColorLight: "color",
				WarnColorLight:       "color",
				FontColorLight:       "color",
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     "color",
				BackgroundColorDark:  "color",
				WarnColorDark:        "color",
				FontColorDark:        "color",
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  true,
				ErrorMessagePopup:    true,
				DisableWatermark:     true,
				ThemeMode:            domain.BrandingPolicyThemeAuto,
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.BrandingSetting{
				Setting: domain.Setting{
					InstanceID:     "foo",
					OrganizationID: nil,
					ID:             gofakeit.UUID(),
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				PrimaryColorLight:    "color",
				BackgroundColorLight: "color",
				WarnColorLight:       "color",
				FontColorLight:       "color",
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     "color",
				BackgroundColorDark:  "color",
				WarnColorDark:        "color",
				FontColorDark:        "color",
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  true,
				ErrorMessagePopup:    true,
				DisableWatermark:     true,
				ThemeMode:            domain.BrandingPolicyThemeAuto,
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.BrandingSetting{
				Setting: domain.Setting{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
					ID:             gofakeit.UUID(),
					CreatedAt:      now,
					UpdatedAt:      now,
				},
				PrimaryColorLight:    "color",
				BackgroundColorLight: "color",
				WarnColorLight:       "color",
				FontColorLight:       "color",
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     "color",
				BackgroundColorDark:  "color",
				WarnColorDark:        "color",
				FontColorDark:        "color",
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  true,
				ErrorMessagePopup:    true,
				DisableWatermark:     true,
				ThemeMode:            domain.BrandingPolicyThemeAuto,
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
			wantErr: new(database.ForeignKeyError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestActivateBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettings()

	existingInstanceSettings := &domain.BrandingSetting{
		Setting: domain.Setting{
			InstanceID:     instanceID,
			OrganizationID: nil,
			ID:             gofakeit.UUID(),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		PrimaryColorLight:    "color",
		BackgroundColorLight: "color",
		WarnColorLight:       "color",
		FontColorLight:       "color",
		LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
		IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
		PrimaryColorDark:     "color",
		BackgroundColorDark:  "color",
		WarnColorDark:        "color",
		FontColorDark:        "color",
		LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
		IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
		HideLoginNameSuffix:  true,
		ErrorMessagePopup:    true,
		DisableWatermark:     true,
		ThemeMode:            domain.BrandingPolicyThemeAuto,
		FontURL:              &url.URL{Scheme: "https", Host: "host"},
	}

	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.BrandingSetting{
		Setting: domain.Setting{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
			ID:             gofakeit.UUID(),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		PrimaryColorLight:    "color",
		BackgroundColorLight: "color",
		WarnColorLight:       "color",
		FontColorLight:       "color",
		LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
		IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
		PrimaryColorDark:     "color",
		BackgroundColorDark:  "color",
		WarnColorDark:        "color",
		FontColorDark:        "color",
		LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
		IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
		HideLoginNameSuffix:  true,
		ErrorMessagePopup:    true,
		DisableWatermark:     true,
		ThemeMode:            domain.BrandingPolicyThemeAuto,
		FontURL:              &url.URL{Scheme: "https", Host: "host"},
	}

	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		wantRowsAffected int64
		condition        database.Condition
		wantErr          error
	}{
		{
			name: "activate instance",
			condition: repo.UniqueCondition(
				existingInstanceSettings.InstanceID,
				nil,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			wantRowsAffected: 1,
		},
		{
			name: "activate organization",
			condition: repo.UniqueCondition(
				existingInstanceSettings.InstanceID,
				existingOrganizationSettings.OrganizationID,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			wantRowsAffected: 1,
		},
		{
			name: "non-existing instance",
			condition: repo.UniqueCondition(
				"foo",
				nil,
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			wantRowsAffected: 0,
		},

		{
			name: "non-existing organization",
			condition: repo.UniqueCondition(
				existingInstanceSettings.InstanceID,
				gu.Ptr("foo"),
				domain.SettingTypeBranding,
				domain.SettingStatePreview,
			),
			wantRowsAffected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			rowsAffected, err := repo.Activate(t.Context(), savepoint, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestDeleteBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettings()

	existingInstanceSettings := &domain.BrandingSetting{
		Setting: domain.Setting{
			InstanceID:     instanceID,
			OrganizationID: nil,
			ID:             gofakeit.UUID(),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		PrimaryColorLight:    "color",
		BackgroundColorLight: "color",
		WarnColorLight:       "color",
		FontColorLight:       "color",
		LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
		IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
		PrimaryColorDark:     "color",
		BackgroundColorDark:  "color",
		WarnColorDark:        "color",
		FontColorDark:        "color",
		LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
		IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
		HideLoginNameSuffix:  true,
		ErrorMessagePopup:    true,
		DisableWatermark:     true,
		ThemeMode:            domain.BrandingPolicyThemeAuto,
		FontURL:              &url.URL{Scheme: "https", Host: "host"},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	_, err = repo.Activate(t.Context(), tx,
		repo.UniqueCondition(
			existingInstanceSettings.InstanceID,
			existingInstanceSettings.OrganizationID,
			domain.SettingTypeBranding,
			domain.SettingStatePreview,
		),
	)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.BrandingSetting{
		Setting: domain.Setting{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
			ID:             gofakeit.UUID(),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		PrimaryColorLight:    "color",
		BackgroundColorLight: "color",
		WarnColorLight:       "color",
		FontColorLight:       "color",
		LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
		IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
		PrimaryColorDark:     "color",
		BackgroundColorDark:  "color",
		WarnColorDark:        "color",
		FontColorDark:        "color",
		LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
		IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
		HideLoginNameSuffix:  true,
		ErrorMessagePopup:    true,
		DisableWatermark:     true,
		ThemeMode:            domain.BrandingPolicyThemeAuto,
		FontURL:              &url.URL{Scheme: "https", Host: "host"},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	_, err = repo.Activate(t.Context(), tx,
		repo.UniqueCondition(
			existingOrganizationSettings.InstanceID,
			existingOrganizationSettings.OrganizationID,
			domain.SettingTypeBranding,
			domain.SettingStatePreview,
		),
	)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.OrganizationIDCondition(&orgID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeBranding, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance, active",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization, active",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice, active",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeBranding, domain.SettingStatePreview),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func TestGetPasswordComplexitySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordComplexitySettings()

	settings := []*domain.PasswordComplexitySetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MinLength:    uint64(8),
			HasLowercase: true,
			HasUppercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MinLength:    uint64(10),
			HasLowercase: true,
			HasUppercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MinLength:    uint64(12),
			HasLowercase: true,
			HasUppercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MinLength:    uint64(14),
			HasLowercase: true,
			HasUppercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.PasswordComplexitySetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListPasswordComplexitySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordComplexitySettings()

	settings := []*domain.PasswordComplexitySetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MinLength:    uint64(8),
			HasLowercase: true,
			HasUppercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MinLength:    uint64(10),
			HasLowercase: true,
			HasUppercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MinLength:    uint64(10),
			HasLowercase: true,
			HasUppercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MinLength:    uint64(12),
			HasLowercase: true,
			HasUppercase: true,
			HasNumber:    true,
			HasSymbol:    true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.PasswordComplexitySetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.PasswordComplexitySetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.PasswordComplexitySetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.PasswordComplexitySetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.PasswordComplexitySetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGetPasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordExpirySettings()

	settings := []*domain.PasswordExpirySetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			ExpireWarnDays: uint64(10),
			MaxAgeDays:     uint64(30),
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			ExpireWarnDays: uint64(11),
			MaxAgeDays:     uint64(31),
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			ExpireWarnDays: uint64(12),
			MaxAgeDays:     uint64(32),
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			ExpireWarnDays: uint64(13),
			MaxAgeDays:     uint64(33),
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.PasswordExpirySetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListPasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordExpirySettings()

	settings := []*domain.PasswordExpirySetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			ExpireWarnDays: uint64(10),
			MaxAgeDays:     uint64(30),
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			ExpireWarnDays: uint64(11),
			MaxAgeDays:     uint64(31),
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			ExpireWarnDays: uint64(12),
			MaxAgeDays:     uint64(32),
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			ExpireWarnDays: uint64(13),
			MaxAgeDays:     uint64(33),
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.PasswordExpirySetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.PasswordExpirySetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.PasswordExpirySetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.PasswordExpirySetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.PasswordExpirySetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGetLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LockoutSettings()

	settings := []*domain.LockoutSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MaxPasswordAttempts: uint64(3),
			MaxOTPAttempts:      uint64(3),
			ShowLockOutFailures: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MaxPasswordAttempts: uint64(4),
			MaxOTPAttempts:      uint64(4),
			ShowLockOutFailures: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MaxPasswordAttempts: uint64(5),
			MaxOTPAttempts:      uint64(5),
			ShowLockOutFailures: false,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MaxPasswordAttempts: uint64(6),
			MaxOTPAttempts:      uint64(6),
			ShowLockOutFailures: true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.LockoutSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LockoutSettings()

	settings := []*domain.LockoutSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MaxPasswordAttempts: uint64(3),
			MaxOTPAttempts:      uint64(3),
			ShowLockOutFailures: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MaxPasswordAttempts: uint64(4),
			MaxOTPAttempts:      uint64(4),
			ShowLockOutFailures: false,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MaxPasswordAttempts: uint64(5),
			MaxOTPAttempts:      uint64(5),
			ShowLockOutFailures: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			MaxPasswordAttempts: uint64(6),
			MaxOTPAttempts:      uint64(6),
			ShowLockOutFailures: true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.LockoutSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.LockoutSetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.LockoutSetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.LockoutSetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.LockoutSetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGetSecuritySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.SecuritySettings()

	settings := []*domain.SecuritySetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"origin1", "origin2"},
			EnableImpersonation:   true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"origin3", "origin4"},
			EnableImpersonation:   false,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"origin5", "origin6"},
			EnableImpersonation:   true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			EnableIframeEmbedding: false,
			AllowedOrigins:        []string{"origin7", "origin8"},
			EnableImpersonation:   true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.SecuritySetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListSecuritySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.SecuritySettings()

	settings := []*domain.SecuritySetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"origin1", "origin2"},
			EnableImpersonation:   true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"origin3", "origin4"},
			EnableImpersonation:   false,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"origin5", "origin6"},
			EnableImpersonation:   true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			EnableIframeEmbedding: false,
			AllowedOrigins:        []string{"origin7", "origin8"},
			EnableImpersonation:   true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.SecuritySetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.SecuritySetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.SecuritySetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.SecuritySetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.SecuritySetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGetDomainSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.DomainSettings()

	settings := []*domain.DomainSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			LoginNameIncludesDomain:                true,
			RequireOrgDomainVerification:           true,
			SMTPSenderAddressMatchesInstanceDomain: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			LoginNameIncludesDomain:                false,
			RequireOrgDomainVerification:           true,
			SMTPSenderAddressMatchesInstanceDomain: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			LoginNameIncludesDomain:                true,
			RequireOrgDomainVerification:           false,
			SMTPSenderAddressMatchesInstanceDomain: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			LoginNameIncludesDomain:                true,
			RequireOrgDomainVerification:           true,
			SMTPSenderAddressMatchesInstanceDomain: false,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.DomainSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListDomainSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.DomainSettings()

	settings := []*domain.DomainSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			LoginNameIncludesDomain:                true,
			RequireOrgDomainVerification:           true,
			SMTPSenderAddressMatchesInstanceDomain: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			LoginNameIncludesDomain:                false,
			RequireOrgDomainVerification:           true,
			SMTPSenderAddressMatchesInstanceDomain: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			LoginNameIncludesDomain:                true,
			RequireOrgDomainVerification:           false,
			SMTPSenderAddressMatchesInstanceDomain: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			LoginNameIncludesDomain:                true,
			RequireOrgDomainVerification:           true,
			SMTPSenderAddressMatchesInstanceDomain: false,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.DomainSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.DomainSetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.DomainSetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.DomainSetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.DomainSetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGetOrganizationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.OrganizationSettings()

	settings := []*domain.OrganizationSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			OrganizationScopedUsernames: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			OrganizationScopedUsernames: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			OrganizationScopedUsernames: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			OrganizationScopedUsernames: true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.OrganizationSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListOrganizationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.OrganizationSettings()

	settings := []*domain.OrganizationSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			OrganizationScopedUsernames: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			OrganizationScopedUsernames: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			OrganizationScopedUsernames: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			OrganizationScopedUsernames: true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.OrganizationSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.OrganizationSetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.OrganizationSetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.OrganizationSetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.OrganizationSetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGetNotificationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.NotificationSettings()

	settings := []*domain.NotificationSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PasswordChange: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PasswordChange: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PasswordChange: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PasswordChange: true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.NotificationSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListNotificationSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.NotificationSettings()

	settings := []*domain.NotificationSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PasswordChange: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PasswordChange: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PasswordChange: true,
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			PasswordChange: true,
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.NotificationSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.NotificationSetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.NotificationSetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.NotificationSetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.NotificationSetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGetLegalAndSupportSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LegalAndSupportSettings()

	settings := []*domain.LegalAndSupportSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			TOSLink:           "https://host",
			PrivacyPolicyLink: "https://host",
			HelpLink:          "https://host",
			SupportEmail:      "email",
			DocsLink:          "https://host",
			CustomLink:        "https://host",
			CustomLinkText:    "linktext",
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			TOSLink:           "https://host",
			PrivacyPolicyLink: "https://host",
			HelpLink:          "https://host",
			SupportEmail:      "email",
			DocsLink:          "https://host",
			CustomLink:        "https://host",
			CustomLinkText:    "linktext",
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			TOSLink:           "https://host",
			PrivacyPolicyLink: "https://host",
			HelpLink:          "https://host",
			SupportEmail:      "email",
			DocsLink:          "https://host",
			CustomLink:        "https://host",
			CustomLinkText:    "linktext",
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			TOSLink:           "https://host",
			PrivacyPolicyLink: "https://host",
			HelpLink:          "https://host",
			SupportEmail:      "email",
			DocsLink:          "https://host",
			CustomLink:        "https://host",
			CustomLinkText:    "linktext",
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.LegalAndSupportSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListLegalAndSupportSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LegalAndSupportSettings()

	settings := []*domain.LegalAndSupportSetting{
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			TOSLink:           "https://host",
			PrivacyPolicyLink: "https://host",
			HelpLink:          "https://host",
			SupportEmail:      "email",
			DocsLink:          "https://host",
			CustomLink:        "https://host",
			CustomLinkText:    "linktext",
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			TOSLink:           "https://host",
			PrivacyPolicyLink: "https://host",
			HelpLink:          "https://host",
			SupportEmail:      "email",
			DocsLink:          "https://host",
			CustomLink:        "https://host",
			CustomLinkText:    "linktext",
		},
		{
			Setting: domain.Setting{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			TOSLink:           "https://host",
			PrivacyPolicyLink: "https://host",
			HelpLink:          "https://host",
			SupportEmail:      "email",
			DocsLink:          "https://host",
			CustomLink:        "https://host",
			CustomLinkText:    "linktext",
		},
		{
			Setting: domain.Setting{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				ID:             gofakeit.UUID(),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			TOSLink:           "https://host",
			PrivacyPolicyLink: "https://host",
			HelpLink:          "https://host",
			SupportEmail:      "email",
			DocsLink:          "https://host",
			CustomLink:        "https://host",
			CustomLinkText:    "linktext",
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.LegalAndSupportSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.LegalAndSupportSetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.LegalAndSupportSetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.LegalAndSupportSetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.LegalAndSupportSetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestGetSecretGeneratorSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.SecretGeneratorSettings()

	settings := []*domain.SecretGeneratorSetting{
		createSecretGeneratorSettings(firstInstanceID, nil),
		createSecretGeneratorSettings(secondInstanceID, nil),
		createSecretGeneratorSettings(firstInstanceID, gu.Ptr(firstOrgID)),
		createSecretGeneratorSettings(secondInstanceID, gu.Ptr(secondOrgID)),
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.SecretGeneratorSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(settings[0].OrganizationID),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: database.And(repo.InstanceIDCondition(settings[0].InstanceID), repo.OrganizationIDCondition(settings[0].OrganizationID)),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: database.And(repo.InstanceIDCondition(settings[3].InstanceID), repo.OrganizationIDCondition(settings[3].OrganizationID)),
			want:      settings[3],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestListSecretGeneratorSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.SecretGeneratorSettings()

	settings := []*domain.SecretGeneratorSetting{
		createSecretGeneratorSettings(firstInstanceID, nil),
		createSecretGeneratorSettings(secondInstanceID, nil),
		createSecretGeneratorSettings(firstInstanceID, gu.Ptr(firstOrgID)),
		createSecretGeneratorSettings(secondInstanceID, gu.Ptr(secondOrgID)),
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.SecretGeneratorSetting
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: database.And(repo.InstanceIDCondition(firstInstanceID), repo.OrganizationIDCondition(gu.Ptr("nix"))),
			want:      []*domain.SecretGeneratorSetting{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.SecretGeneratorSetting{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.SecretGeneratorSetting{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.SecretGeneratorSetting{settings[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(repo.InstanceIDColumn(), repo.OrganizationIDColumn()),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.EqualExportedValues(t, tt.want, got)
		})
	}
}

func TestSetSecretGeneratorSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.SecretGeneratorSettings()

	existingSettings := createSecretGeneratorSettings(instanceID, gu.Ptr(orgID))

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.SecretGeneratorSetting
		wantErr  error
	}{
		{
			name:     "create instance",
			settings: createSecretGeneratorSettings(instanceID, nil),
		},
		{
			name:     "update organization",
			settings: createSecretGeneratorSettings(instanceID, gu.Ptr(orgID)),
		},
		{
			name:     "non-existing instance",
			settings: createSecretGeneratorSettings("foo", nil),
			wantErr:  new(database.ForeignKeyError),
		},
		{
			name:     "non-existing org",
			settings: createSecretGeneratorSettings(instanceID, gu.Ptr("foo")),
			wantErr:  new(database.ForeignKeyError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := repo.Set(t.Context(), savepoint, tt.settings)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestDeleteSecretGeneratorSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.SecretGeneratorSettings()

	existingInstanceSettings := createSecretGeneratorSettings(instanceID, nil)
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := createSecretGeneratorSettings(instanceID, gu.Ptr(orgID))
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.OrganizationIDCondition(gu.Ptr(orgID)),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.InstanceIDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, gu.Ptr("foo"), domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.UniqueCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.OrganizationID, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.UniqueCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.OrganizationID, domain.SettingTypeSecretGenerator, domain.SettingStateActive),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := repo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}

func createDefaultSecretGeneratorAttrs() domain.SecretGeneratorAttrs {
	return domain.SecretGeneratorAttrs{
		Length:              uint(32),
		IncludeUpperLetters: true,
		IncludeLowerLetters: true,
		IncludeDigits:       true,
		IncludeSymbols:      false,
	}
}

func createDefaultSecretGeneratorAttrsWithExpiry(expiry time.Duration) domain.SecretGeneratorAttrsWithExpiry {
	return domain.SecretGeneratorAttrsWithExpiry{
		Expiry: gu.Ptr(expiry),
		SecretGeneratorAttrs: domain.SecretGeneratorAttrs{
			Length:              uint(32),
			IncludeUpperLetters: true,
			IncludeLowerLetters: true,
			IncludeDigits:       true,
			IncludeSymbols:      false,
		},
	}
}

func createDefaultSecretGeneratorSettingsAttributes() domain.SecretGeneratorSetting {
	attrs := createDefaultSecretGeneratorAttrs()
	oneDay := 24 * time.Hour
	thirtyMinutes := 30 * time.Minute
	tenMinutes := 10 * time.Minute

	return domain.SecretGeneratorSetting{
		ClientSecret: &domain.ClientSecretAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(tenMinutes),
		},
		InitializeUserCode: &domain.InitializeUserCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(oneDay),
		},
		EmailVerificationCode: &domain.EmailVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
		PhoneVerificationCode: &domain.PhoneVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
		PasswordVerificationCode: &domain.PasswordVerificationCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
		PasswordlessInitCode: &domain.PasswordlessInitCodeAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(tenMinutes),
		},
		DomainVerification: &domain.DomainVerificationAttributes{
			SecretGeneratorAttrs: attrs,
		},
		OTPSMS: &domain.OTPSMSAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
		OTPEmail: &domain.OTPEmailAttributes{
			SecretGeneratorAttrsWithExpiry: createDefaultSecretGeneratorAttrsWithExpiry(thirtyMinutes),
		},
	}
}

func createSecretGeneratorSettings(instanceID string, organizationID *string) *domain.SecretGeneratorSetting {
	setting := createDefaultSecretGeneratorSettingsAttributes()
	setting.Setting = domain.Setting{
		InstanceID:     instanceID,
		OrganizationID: organizationID,
		ID:             gofakeit.UUID(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return &setting
}
