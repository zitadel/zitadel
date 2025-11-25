package repository_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestGetLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LoginSettingsRepository()

	settings := []*domain.LoginSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(true),
				AllowRegister:              gu.Ptr(true),
				AllowExternalIDP:           gu.Ptr(true),
				ForceMFA:                   gu.Ptr(true),
				ForceMFALocalOnly:          gu.Ptr(true),
				HidePasswordReset:          gu.Ptr(true),
				IgnoreUnknownUsernames:     gu.Ptr(true),
				AllowDomainDiscovery:       gu.Ptr(true),
				DisableLoginWithEmail:      gu.Ptr(true),
				DisableLoginWithPhone:      gu.Ptr(true),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
				DefaultRedirectURI:         gu.Ptr("uri"),
				PasswordCheckLifetime:      gu.Ptr(time.Hour),
				ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
				MFAInitSkipLifetime:        gu.Ptr(time.Hour),
				SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
				MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(false),
				AllowRegister:              gu.Ptr(false),
				AllowExternalIDP:           gu.Ptr(false),
				ForceMFA:                   gu.Ptr(false),
				ForceMFALocalOnly:          gu.Ptr(false),
				HidePasswordReset:          gu.Ptr(false),
				IgnoreUnknownUsernames:     gu.Ptr(false),
				AllowDomainDiscovery:       gu.Ptr(false),
				DisableLoginWithEmail:      gu.Ptr(false),
				DisableLoginWithPhone:      gu.Ptr(false),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
				DefaultRedirectURI:         gu.Ptr("uri2"),
				PasswordCheckLifetime:      gu.Ptr(time.Minute),
				ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
				MFAInitSkipLifetime:        gu.Ptr(time.Minute),
				SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
				MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(true),
				AllowRegister:              gu.Ptr(true),
				AllowExternalIDP:           gu.Ptr(true),
				ForceMFA:                   gu.Ptr(true),
				ForceMFALocalOnly:          gu.Ptr(true),
				HidePasswordReset:          gu.Ptr(true),
				IgnoreUnknownUsernames:     gu.Ptr(true),
				AllowDomainDiscovery:       gu.Ptr(true),
				DisableLoginWithEmail:      gu.Ptr(true),
				DisableLoginWithPhone:      gu.Ptr(true),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
				DefaultRedirectURI:         gu.Ptr("uri"),
				PasswordCheckLifetime:      gu.Ptr(time.Hour),
				ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
				MFAInitSkipLifetime:        gu.Ptr(time.Hour),
				SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
				MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(false),
				AllowRegister:              gu.Ptr(false),
				AllowExternalIDP:           gu.Ptr(false),
				ForceMFA:                   gu.Ptr(false),
				ForceMFALocalOnly:          gu.Ptr(false),
				HidePasswordReset:          gu.Ptr(false),
				IgnoreUnknownUsernames:     gu.Ptr(false),
				AllowDomainDiscovery:       gu.Ptr(false),
				DisableLoginWithEmail:      gu.Ptr(false),
				DisableLoginWithPhone:      gu.Ptr(false),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
				DefaultRedirectURI:         gu.Ptr("uri2"),
				PasswordCheckLifetime:      gu.Ptr(time.Minute),
				ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
				MFAInitSkipLifetime:        gu.Ptr(time.Minute),
				SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
				MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.LoginSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
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
	repo := repository.LoginSettingsRepository()

	settings := []*domain.LoginSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(true),
				AllowRegister:              gu.Ptr(true),
				AllowExternalIDP:           gu.Ptr(true),
				ForceMFA:                   gu.Ptr(true),
				ForceMFALocalOnly:          gu.Ptr(true),
				HidePasswordReset:          gu.Ptr(true),
				IgnoreUnknownUsernames:     gu.Ptr(true),
				AllowDomainDiscovery:       gu.Ptr(true),
				DisableLoginWithEmail:      gu.Ptr(true),
				DisableLoginWithPhone:      gu.Ptr(true),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
				DefaultRedirectURI:         gu.Ptr("uri"),
				PasswordCheckLifetime:      gu.Ptr(time.Hour),
				ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
				MFAInitSkipLifetime:        gu.Ptr(time.Hour),
				SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
				MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(false),
				AllowRegister:              gu.Ptr(false),
				AllowExternalIDP:           gu.Ptr(false),
				ForceMFA:                   gu.Ptr(false),
				ForceMFALocalOnly:          gu.Ptr(false),
				HidePasswordReset:          gu.Ptr(false),
				IgnoreUnknownUsernames:     gu.Ptr(false),
				AllowDomainDiscovery:       gu.Ptr(false),
				DisableLoginWithEmail:      gu.Ptr(false),
				DisableLoginWithPhone:      gu.Ptr(false),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
				DefaultRedirectURI:         gu.Ptr("uri2"),
				PasswordCheckLifetime:      gu.Ptr(time.Minute),
				ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
				MFAInitSkipLifetime:        gu.Ptr(time.Minute),
				SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
				MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(true),
				AllowRegister:              gu.Ptr(true),
				AllowExternalIDP:           gu.Ptr(true),
				ForceMFA:                   gu.Ptr(true),
				ForceMFALocalOnly:          gu.Ptr(true),
				HidePasswordReset:          gu.Ptr(true),
				IgnoreUnknownUsernames:     gu.Ptr(true),
				AllowDomainDiscovery:       gu.Ptr(true),
				DisableLoginWithEmail:      gu.Ptr(true),
				DisableLoginWithPhone:      gu.Ptr(true),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
				DefaultRedirectURI:         gu.Ptr("uri"),
				PasswordCheckLifetime:      gu.Ptr(time.Hour),
				ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
				MFAInitSkipLifetime:        gu.Ptr(time.Hour),
				SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
				MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LoginSettingsAttributes: domain.LoginSettingsAttributes{
				AllowUserNamePassword:      gu.Ptr(false),
				AllowRegister:              gu.Ptr(false),
				AllowExternalIDP:           gu.Ptr(false),
				ForceMFA:                   gu.Ptr(false),
				ForceMFALocalOnly:          gu.Ptr(false),
				HidePasswordReset:          gu.Ptr(false),
				IgnoreUnknownUsernames:     gu.Ptr(false),
				AllowDomainDiscovery:       gu.Ptr(false),
				DisableLoginWithEmail:      gu.Ptr(false),
				DisableLoginWithPhone:      gu.Ptr(false),
				PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
				DefaultRedirectURI:         gu.Ptr("uri2"),
				PasswordCheckLifetime:      gu.Ptr(time.Minute),
				ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
				MFAInitSkipLifetime:        gu.Ptr(time.Minute),
				SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
				MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.LoginSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.LoginSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.LoginSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.LoginSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.LoginSettings{settings[2]},
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
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LoginSettingsRepository()

	existingSettings := &domain.LoginSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LoginSettingsAttributes: domain.LoginSettingsAttributes{
			AllowUserNamePassword:      gu.Ptr(false),
			AllowRegister:              gu.Ptr(false),
			AllowExternalIDP:           gu.Ptr(false),
			ForceMFA:                   gu.Ptr(false),
			ForceMFALocalOnly:          gu.Ptr(false),
			HidePasswordReset:          gu.Ptr(false),
			IgnoreUnknownUsernames:     gu.Ptr(false),
			AllowDomainDiscovery:       gu.Ptr(false),
			DisableLoginWithEmail:      gu.Ptr(false),
			DisableLoginWithPhone:      gu.Ptr(false),
			PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
			DefaultRedirectURI:         gu.Ptr("uri2"),
			PasswordCheckLifetime:      gu.Ptr(time.Minute),
			ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
			MFAInitSkipLifetime:        gu.Ptr(time.Minute),
			SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
			MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.LoginSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.LoginSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				LoginSettingsAttributes: domain.LoginSettingsAttributes{
					AllowUserNamePassword:      gu.Ptr(true),
					AllowRegister:              gu.Ptr(true),
					AllowExternalIDP:           gu.Ptr(true),
					ForceMFA:                   gu.Ptr(true),
					ForceMFALocalOnly:          gu.Ptr(true),
					HidePasswordReset:          gu.Ptr(true),
					IgnoreUnknownUsernames:     gu.Ptr(true),
					AllowDomainDiscovery:       gu.Ptr(true),
					DisableLoginWithEmail:      gu.Ptr(true),
					DisableLoginWithPhone:      gu.Ptr(true),
					PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
					DefaultRedirectURI:         gu.Ptr("uri"),
					PasswordCheckLifetime:      gu.Ptr(time.Hour),
					ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
					MFAInitSkipLifetime:        gu.Ptr(time.Hour),
					SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
					MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.LoginSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				LoginSettingsAttributes: domain.LoginSettingsAttributes{
					AllowUserNamePassword:      gu.Ptr(true),
					AllowRegister:              gu.Ptr(true),
					AllowExternalIDP:           gu.Ptr(true),
					ForceMFA:                   gu.Ptr(true),
					ForceMFALocalOnly:          gu.Ptr(true),
					HidePasswordReset:          gu.Ptr(true),
					IgnoreUnknownUsernames:     gu.Ptr(true),
					AllowDomainDiscovery:       gu.Ptr(true),
					DisableLoginWithEmail:      gu.Ptr(true),
					DisableLoginWithPhone:      gu.Ptr(true),
					PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
					DefaultRedirectURI:         gu.Ptr("uri"),
					PasswordCheckLifetime:      gu.Ptr(time.Hour),
					ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
					MFAInitSkipLifetime:        gu.Ptr(time.Hour),
					SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
					MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.LoginSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				LoginSettingsAttributes: domain.LoginSettingsAttributes{
					AllowUserNamePassword:      gu.Ptr(true),
					AllowRegister:              gu.Ptr(true),
					AllowExternalIDP:           gu.Ptr(true),
					ForceMFA:                   gu.Ptr(true),
					ForceMFALocalOnly:          gu.Ptr(true),
					HidePasswordReset:          gu.Ptr(true),
					IgnoreUnknownUsernames:     gu.Ptr(true),
					AllowDomainDiscovery:       gu.Ptr(true),
					DisableLoginWithEmail:      gu.Ptr(true),
					DisableLoginWithPhone:      gu.Ptr(true),
					PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
					DefaultRedirectURI:         gu.Ptr("uri"),
					PasswordCheckLifetime:      gu.Ptr(time.Hour),
					ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
					MFAInitSkipLifetime:        gu.Ptr(time.Hour),
					SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
					MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.LoginSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				LoginSettingsAttributes: domain.LoginSettingsAttributes{
					AllowUserNamePassword:      gu.Ptr(true),
					AllowRegister:              gu.Ptr(true),
					AllowExternalIDP:           gu.Ptr(true),
					ForceMFA:                   gu.Ptr(true),
					ForceMFALocalOnly:          gu.Ptr(true),
					HidePasswordReset:          gu.Ptr(true),
					IgnoreUnknownUsernames:     gu.Ptr(true),
					AllowDomainDiscovery:       gu.Ptr(true),
					DisableLoginWithEmail:      gu.Ptr(true),
					DisableLoginWithPhone:      gu.Ptr(true),
					PasswordlessType:           gu.Ptr(domain.PasswordlessTypeAllowed),
					DefaultRedirectURI:         gu.Ptr("uri"),
					PasswordCheckLifetime:      gu.Ptr(time.Hour),
					ExternalLoginCheckLifetime: gu.Ptr(time.Hour),
					MFAInitSkipLifetime:        gu.Ptr(time.Hour),
					SecondFactorCheckLifetime:  gu.Ptr(time.Hour),
					MultiFactorCheckLifetime:   gu.Ptr(time.Hour),
				},
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

func TestDeleteLoginSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LoginSettingsRepository()

	existingInstanceSettings := &domain.LoginSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		LoginSettingsAttributes: domain.LoginSettingsAttributes{
			AllowUserNamePassword:      gu.Ptr(false),
			AllowRegister:              gu.Ptr(false),
			AllowExternalIDP:           gu.Ptr(false),
			ForceMFA:                   gu.Ptr(false),
			ForceMFALocalOnly:          gu.Ptr(false),
			HidePasswordReset:          gu.Ptr(false),
			IgnoreUnknownUsernames:     gu.Ptr(false),
			AllowDomainDiscovery:       gu.Ptr(false),
			DisableLoginWithEmail:      gu.Ptr(false),
			DisableLoginWithPhone:      gu.Ptr(false),
			PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
			DefaultRedirectURI:         gu.Ptr("uri"),
			PasswordCheckLifetime:      gu.Ptr(time.Minute),
			ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
			MFAInitSkipLifetime:        gu.Ptr(time.Minute),
			SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
			MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.LoginSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LoginSettingsAttributes: domain.LoginSettingsAttributes{
			AllowUserNamePassword:      gu.Ptr(false),
			AllowRegister:              gu.Ptr(false),
			AllowExternalIDP:           gu.Ptr(false),
			ForceMFA:                   gu.Ptr(false),
			ForceMFALocalOnly:          gu.Ptr(false),
			HidePasswordReset:          gu.Ptr(false),
			IgnoreUnknownUsernames:     gu.Ptr(false),
			AllowDomainDiscovery:       gu.Ptr(false),
			DisableLoginWithEmail:      gu.Ptr(false),
			DisableLoginWithPhone:      gu.Ptr(false),
			PasswordlessType:           gu.Ptr(domain.PasswordlessTypeNotAllowed),
			DefaultRedirectURI:         gu.Ptr("uri"),
			PasswordCheckLifetime:      gu.Ptr(time.Minute),
			ExternalLoginCheckLifetime: gu.Ptr(time.Minute),
			MFAInitSkipLifetime:        gu.Ptr(time.Minute),
			SecondFactorCheckLifetime:  gu.Ptr(time.Minute),
			MultiFactorCheckLifetime:   gu.Ptr(time.Minute),
		},
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
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.PrimaryKeyCondition(instanceID, "foo"),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
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
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.BrandingSettingsRepository()

	settings := []*domain.BrandingSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
				State:          domain.SettingStatePreview,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
				State:          domain.SettingStatePreview,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
				State:          domain.SettingStatePreview,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
				State:          domain.SettingStatePreview,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}
	_, err := repo.Activate(t.Context(), tx, repo.PrimaryKeyCondition(settings[0].InstanceID, settings[0].ID))
	require.NoError(t, err)
	activeInstanceSetting := *settings[0]
	activeInstanceSetting.State = domain.SettingStateActive

	_, err = repo.Activate(t.Context(), tx, repo.PrimaryKeyCondition(settings[3].InstanceID, settings[3].ID))
	require.NoError(t, err)
	activeOrganizationSetting := *settings[3]
	activeOrganizationSetting.State = domain.SettingStateActive

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.BrandingSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
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
				repo.PrimaryKeyCondition(settings[0].InstanceID, settings[0].ID),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: settings[0],
		},
		{
			name: "ok, instance, active",
			condition: database.And(
				repo.PrimaryKeyCondition(activeInstanceSetting.InstanceID, activeInstanceSetting.ID),
				repo.StateCondition(domain.SettingStateActive),
			),
			want: &activeInstanceSetting,
		},
		{
			name: "ok, organization",
			condition: database.And(
				repo.PrimaryKeyCondition(settings[3].InstanceID, settings[3].ID),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: settings[3],
		},
		{
			name: "ok, organization",
			condition: database.And(
				repo.PrimaryKeyCondition(activeOrganizationSetting.InstanceID, activeOrganizationSetting.ID),
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
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.BrandingSettingsRepository()

	settings := []*domain.BrandingSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
				PrimaryColorLight:    gu.Ptr("color"),
				BackgroundColorLight: gu.Ptr("color"),
				WarnColorLight:       gu.Ptr("color"),
				FontColorLight:       gu.Ptr("color"),
				LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
				IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
				PrimaryColorDark:     gu.Ptr("color"),
				BackgroundColorDark:  gu.Ptr("color"),
				WarnColorDark:        gu.Ptr("color"),
				FontColorDark:        gu.Ptr("color"),
				LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
				IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
				HideLoginNameSuffix:  gu.Ptr(true),
				ErrorMsgPopup:        gu.Ptr(true),
				DisableWatermark:     gu.Ptr(true),
				ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
				FontURL:              &url.URL{Scheme: "https", Host: "host"},
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}
	_, err := repo.Activate(t.Context(), tx, repo.PrimaryKeyCondition(settings[0].InstanceID, settings[0].ID))
	require.NoError(t, err)
	activeInstanceSetting := *settings[0]
	activeInstanceSetting.State = domain.SettingStateActive

	_, err = repo.Activate(t.Context(), tx, repo.PrimaryKeyCondition(settings[2].InstanceID, settings[2].ID))
	require.NoError(t, err)
	activeOrganizationSetting := *settings[2]
	activeOrganizationSetting.State = domain.SettingStateActive

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.BrandingSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.BrandingSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.BrandingSettings{&activeOrganizationSetting, settings[2], &activeInstanceSetting, settings[0]},
		},
		{
			name: "all from instance, preview",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.StateCondition(domain.SettingStatePreview),
			),
			want: []*domain.BrandingSettings{settings[2], settings[0]},
		},
		{
			name: "all from instance, active",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.StateCondition(domain.SettingStateActive),
			),
			want: []*domain.BrandingSettings{&activeOrganizationSetting, &activeInstanceSetting},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.BrandingSettings{&activeInstanceSetting, settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.BrandingSettings{&activeOrganizationSetting, settings[2]},
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
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettingsRepository()

	existingSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.BrandingSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.BrandingSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
					PrimaryColorLight:    gu.Ptr("color"),
					BackgroundColorLight: gu.Ptr("color"),
					WarnColorLight:       gu.Ptr("color"),
					FontColorLight:       gu.Ptr("color"),
					LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
					IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
					PrimaryColorDark:     gu.Ptr("color"),
					BackgroundColorDark:  gu.Ptr("color"),
					WarnColorDark:        gu.Ptr("color"),
					FontColorDark:        gu.Ptr("color"),
					LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
					IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
					HideLoginNameSuffix:  gu.Ptr(true),
					ErrorMsgPopup:        gu.Ptr(true),
					DisableWatermark:     gu.Ptr(true),
					ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
					FontURL:              &url.URL{Scheme: "https", Host: "host"},
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.BrandingSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
					PrimaryColorLight:    gu.Ptr("color"),
					BackgroundColorLight: gu.Ptr("color"),
					WarnColorLight:       gu.Ptr("color"),
					FontColorLight:       gu.Ptr("color"),
					LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
					IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
					PrimaryColorDark:     gu.Ptr("color"),
					BackgroundColorDark:  gu.Ptr("color"),
					WarnColorDark:        gu.Ptr("color"),
					FontColorDark:        gu.Ptr("color"),
					LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
					IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
					HideLoginNameSuffix:  gu.Ptr(true),
					ErrorMsgPopup:        gu.Ptr(true),
					DisableWatermark:     gu.Ptr(true),
					ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
					FontURL:              &url.URL{Scheme: "https", Host: "host"},
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.BrandingSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
					PrimaryColorLight:    gu.Ptr("color"),
					BackgroundColorLight: gu.Ptr("color"),
					WarnColorLight:       gu.Ptr("color"),
					FontColorLight:       gu.Ptr("color"),
					LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
					IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
					PrimaryColorDark:     gu.Ptr("color"),
					BackgroundColorDark:  gu.Ptr("color"),
					WarnColorDark:        gu.Ptr("color"),
					FontColorDark:        gu.Ptr("color"),
					LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
					IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
					HideLoginNameSuffix:  gu.Ptr(true),
					ErrorMsgPopup:        gu.Ptr(true),
					DisableWatermark:     gu.Ptr(true),
					ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
					FontURL:              &url.URL{Scheme: "https", Host: "host"},
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.BrandingSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
					PrimaryColorLight:    gu.Ptr("color"),
					BackgroundColorLight: gu.Ptr("color"),
					WarnColorLight:       gu.Ptr("color"),
					FontColorLight:       gu.Ptr("color"),
					LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
					IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
					PrimaryColorDark:     gu.Ptr("color"),
					BackgroundColorDark:  gu.Ptr("color"),
					WarnColorDark:        gu.Ptr("color"),
					FontColorDark:        gu.Ptr("color"),
					LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
					IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
					HideLoginNameSuffix:  gu.Ptr(true),
					ErrorMsgPopup:        gu.Ptr(true),
					DisableWatermark:     gu.Ptr(true),
					ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
					FontURL:              &url.URL{Scheme: "https", Host: "host"},
				},
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

func TestDeleteBrandingSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.BrandingSettingsRepository()

	existingInstanceSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	_, err = repo.Activate(t.Context(), tx, repo.PrimaryKeyCondition(existingInstanceSettings.InstanceID, existingInstanceSettings.ID))
	require.NoError(t, err)
	activeInstanceSetting := *existingInstanceSettings
	activeInstanceSetting.State = domain.SettingStateActive

	existingOrganizationSettings := &domain.BrandingSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		BrandingSettingsAttributes: domain.BrandingSettingsAttributes{
			PrimaryColorLight:    gu.Ptr("color"),
			BackgroundColorLight: gu.Ptr("color"),
			WarnColorLight:       gu.Ptr("color"),
			FontColorLight:       gu.Ptr("color"),
			LogoURLLight:         &url.URL{Scheme: "https", Host: "host"},
			IconURLLight:         &url.URL{Scheme: "https", Host: "host"},
			PrimaryColorDark:     gu.Ptr("color"),
			BackgroundColorDark:  gu.Ptr("color"),
			WarnColorDark:        gu.Ptr("color"),
			FontColorDark:        gu.Ptr("color"),
			LogoURLDark:          &url.URL{Scheme: "https", Host: "host"},
			IconURLDark:          &url.URL{Scheme: "https", Host: "host"},
			HideLoginNameSuffix:  gu.Ptr(true),
			ErrorMsgPopup:        gu.Ptr(true),
			DisableWatermark:     gu.Ptr(true),
			ThemeMode:            gu.Ptr(domain.BrandingPolicyThemeAuto),
			FontURL:              &url.URL{Scheme: "https", Host: "host"},
		},
	}
	err = repo.Set(t.Context(), tx, existingOrganizationSettings)
	require.NoError(t, err)

	_, err = repo.Activate(t.Context(), tx, repo.PrimaryKeyCondition(existingOrganizationSettings.InstanceID, existingOrganizationSettings.ID))
	require.NoError(t, err)
	activeOrganizationSetting := *existingOrganizationSettings
	activeOrganizationSetting.State = domain.SettingStateActive

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.PrimaryKeyCondition(instanceID, "foo"),
			wantRowsAffected: 0,
		},
		{
			name: "delete instance, active",
			condition: database.And(
				repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
				repo.StateCondition(domain.SettingStateActive),
			),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
			wantRowsAffected: 2,
		},
		{
			name:             "delete organization twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
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
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordComplexitySettingsRepository()

	settings := []*domain.PasswordComplexitySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(8)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(10)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(12)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(14)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.PasswordComplexitySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
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
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordComplexitySettingsRepository()

	settings := []*domain.PasswordComplexitySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(8)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(10)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(10)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
				MinLength:    gu.Ptr(uint64(12)),
				HasLowercase: gu.Ptr(true),
				HasUppercase: gu.Ptr(true),
				HasNumber:    gu.Ptr(true),
				HasSymbol:    gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.PasswordComplexitySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.PasswordComplexitySettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.PasswordComplexitySettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.PasswordComplexitySettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.PasswordComplexitySettings{settings[2]},
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

func TestSetPasswordComplexitySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.PasswordComplexitySettingsRepository()

	existingSettings := &domain.PasswordComplexitySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
			MinLength:    gu.Ptr(uint64(8)),
			HasLowercase: gu.Ptr(true),
			HasUppercase: gu.Ptr(true),
			HasNumber:    gu.Ptr(true),
			HasSymbol:    gu.Ptr(true),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.PasswordComplexitySettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.PasswordComplexitySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
					MinLength:    gu.Ptr(uint64(8)),
					HasLowercase: gu.Ptr(true),
					HasUppercase: gu.Ptr(true),
					HasNumber:    gu.Ptr(true),
					HasSymbol:    gu.Ptr(true),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.PasswordComplexitySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
					MinLength:    gu.Ptr(uint64(10)),
					HasLowercase: gu.Ptr(true),
					HasUppercase: gu.Ptr(true),
					HasNumber:    gu.Ptr(true),
					HasSymbol:    gu.Ptr(true),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.PasswordComplexitySettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
					MinLength:    gu.Ptr(uint64(8)),
					HasLowercase: gu.Ptr(true),
					HasUppercase: gu.Ptr(true),
					HasNumber:    gu.Ptr(true),
					HasSymbol:    gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.PasswordComplexitySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
					MinLength:    gu.Ptr(uint64(8)),
					HasLowercase: gu.Ptr(true),
					HasUppercase: gu.Ptr(true),
					HasNumber:    gu.Ptr(true),
					HasSymbol:    gu.Ptr(true),
				},
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

func TestDeletePasswordComplexitySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.PasswordComplexitySettingsRepository()

	existingInstanceSettings := &domain.PasswordComplexitySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
			MinLength:    gu.Ptr(uint64(8)),
			HasLowercase: gu.Ptr(true),
			HasUppercase: gu.Ptr(true),
			HasNumber:    gu.Ptr(true),
			HasSymbol:    gu.Ptr(true),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.PasswordComplexitySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		PasswordComplexitySettingsAttributes: domain.PasswordComplexitySettingsAttributes{
			MinLength:    gu.Ptr(uint64(8)),
			HasLowercase: gu.Ptr(true),
			HasUppercase: gu.Ptr(true),
			HasNumber:    gu.Ptr(true),
			HasSymbol:    gu.Ptr(true),
		},
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
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.PrimaryKeyCondition(instanceID, "foo"),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
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

func TestGetPasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordExpirySettingsRepository()

	settings := []*domain.PasswordExpirySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(10)),
				MaxAgeDays:     gu.Ptr(uint64(30)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(11)),
				MaxAgeDays:     gu.Ptr(uint64(31)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(12)),
				MaxAgeDays:     gu.Ptr(uint64(32)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(13)),
				MaxAgeDays:     gu.Ptr(uint64(33)),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.PasswordExpirySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
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
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.PasswordExpirySettingsRepository()

	settings := []*domain.PasswordExpirySettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(10)),
				MaxAgeDays:     gu.Ptr(uint64(30)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(11)),
				MaxAgeDays:     gu.Ptr(uint64(31)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(12)),
				MaxAgeDays:     gu.Ptr(uint64(32)),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
				ExpireWarnDays: gu.Ptr(uint64(13)),
				MaxAgeDays:     gu.Ptr(uint64(33)),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.PasswordExpirySettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.PasswordExpirySettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.PasswordExpirySettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.PasswordExpirySettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.PasswordExpirySettings{settings[2]},
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

func TestSetPasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.PasswordExpirySettingsRepository()

	existingSettings := &domain.PasswordExpirySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
			ExpireWarnDays: gu.Ptr(uint64(10)),
			MaxAgeDays:     gu.Ptr(uint64(30)),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.PasswordExpirySettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.PasswordExpirySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
					ExpireWarnDays: gu.Ptr(uint64(10)),
					MaxAgeDays:     gu.Ptr(uint64(30)),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.PasswordExpirySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
					ExpireWarnDays: gu.Ptr(uint64(10)),
					MaxAgeDays:     gu.Ptr(uint64(30)),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.PasswordExpirySettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
					ExpireWarnDays: gu.Ptr(uint64(10)),
					MaxAgeDays:     gu.Ptr(uint64(30)),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.PasswordExpirySettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
					ExpireWarnDays: gu.Ptr(uint64(10)),
					MaxAgeDays:     gu.Ptr(uint64(30)),
				},
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

func TestDeletePasswordExpirySettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.PasswordExpirySettingsRepository()

	existingInstanceSettings := &domain.PasswordExpirySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
			ExpireWarnDays: gu.Ptr(uint64(10)),
			MaxAgeDays:     gu.Ptr(uint64(30)),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.PasswordExpirySettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		PasswordExpirySettingsAttributes: domain.PasswordExpirySettingsAttributes{
			ExpireWarnDays: gu.Ptr(uint64(10)),
			MaxAgeDays:     gu.Ptr(uint64(30)),
		},
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
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.PrimaryKeyCondition(instanceID, "foo"),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
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

func TestGetLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LockoutSettingsRepository()

	settings := []*domain.LockoutSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(3)),
				MaxOTPAttempts:      gu.Ptr(uint64(3)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(4)),
				MaxOTPAttempts:      gu.Ptr(uint64(4)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(5)),
				MaxOTPAttempts:      gu.Ptr(uint64(5)),
				ShowLockOutFailures: gu.Ptr(false),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(6)),
				MaxOTPAttempts:      gu.Ptr(uint64(6)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.LockoutSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.IDCondition(settings[0].ID),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "not found",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: repo.InstanceIDCondition(firstInstanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok, instance",
			condition: repo.PrimaryKeyCondition(firstInstanceID, settings[0].ID),
			want:      settings[0],
		},
		{
			name:      "ok, organization",
			condition: repo.PrimaryKeyCondition(secondInstanceID, settings[3].ID),
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
	defer rollback()

	firstInstanceID := createInstance(t, tx)
	secondInstanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, firstInstanceID)
	secondOrgID := createOrganization(t, tx, secondInstanceID)
	repo := repository.LockoutSettingsRepository()

	settings := []*domain.LockoutSettings{
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: nil,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(3)),
				MaxOTPAttempts:      gu.Ptr(uint64(3)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: nil,
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(4)),
				MaxOTPAttempts:      gu.Ptr(uint64(4)),
				ShowLockOutFailures: gu.Ptr(false),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     firstInstanceID,
				OrganizationID: gu.Ptr(firstOrgID),
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(5)),
				MaxOTPAttempts:      gu.Ptr(uint64(5)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
		{
			Settings: domain.Settings{
				InstanceID:     secondInstanceID,
				OrganizationID: gu.Ptr(secondOrgID),
			},
			LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
				MaxPasswordAttempts: gu.Ptr(uint64(6)),
				MaxOTPAttempts:      gu.Ptr(uint64(6)),
				ShowLockOutFailures: gu.Ptr(true),
			},
		},
	}

	for _, setting := range settings {
		err := repo.Set(t.Context(), tx, setting)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.LockoutSettings
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			wantErr:   database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: repo.PrimaryKeyCondition(firstInstanceID, "nix"),
			want:      []*domain.LockoutSettings{},
		},
		{
			name:      "all from instance",
			condition: repo.InstanceIDCondition(firstInstanceID),
			want:      []*domain.LockoutSettings{settings[2], settings[0]},
		},
		{
			name: "only from instance",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(nil),
			),
			want: []*domain.LockoutSettings{settings[0]},
		},
		{
			name: "all from first org",
			condition: database.And(
				repo.InstanceIDCondition(firstInstanceID),
				repo.OrganizationIDCondition(gu.Ptr(firstOrgID)),
			),
			want: []*domain.LockoutSettings{settings[2]},
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

func TestSetLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LockoutSettingsRepository()

	existingSettings := &domain.LockoutSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
			MaxPasswordAttempts: gu.Ptr(uint64(3)),
			MaxOTPAttempts:      gu.Ptr(uint64(3)),
			ShowLockOutFailures: gu.Ptr(true),
		},
	}

	err := repo.Set(t.Context(), tx, existingSettings)
	require.NoError(t, err)

	tests := []struct {
		name     string
		settings *domain.LockoutSettings
		wantErr  error
	}{
		{
			name: "create instance",
			settings: &domain.LockoutSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: nil,
				},
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxPasswordAttempts: gu.Ptr(uint64(4)),
					MaxOTPAttempts:      gu.Ptr(uint64(4)),
					ShowLockOutFailures: gu.Ptr(true),
				},
			},
		},
		{
			name: "update organization",
			settings: &domain.LockoutSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr(orgID),
				},
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxPasswordAttempts: gu.Ptr(uint64(5)),
					MaxOTPAttempts:      gu.Ptr(uint64(5)),
					ShowLockOutFailures: gu.Ptr(true),
				},
			},
		},
		{
			name: "non-existing instance",
			settings: &domain.LockoutSettings{
				Settings: domain.Settings{
					InstanceID:     "foo",
					OrganizationID: nil,
				},
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxPasswordAttempts: gu.Ptr(uint64(3)),
					MaxOTPAttempts:      gu.Ptr(uint64(3)),
					ShowLockOutFailures: gu.Ptr(true),
				},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			settings: &domain.LockoutSettings{
				Settings: domain.Settings{
					InstanceID:     instanceID,
					OrganizationID: gu.Ptr("foo"),
				},
				LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
					MaxPasswordAttempts: gu.Ptr(uint64(3)),
					MaxOTPAttempts:      gu.Ptr(uint64(3)),
					ShowLockOutFailures: gu.Ptr(true),
				},
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

func TestDeleteLockoutSettings(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	repo := repository.LockoutSettingsRepository()

	existingInstanceSettings := &domain.LockoutSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: nil,
		},
		LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
			MaxPasswordAttempts: gu.Ptr(uint64(3)),
			MaxOTPAttempts:      gu.Ptr(uint64(3)),
			ShowLockOutFailures: gu.Ptr(true),
		},
	}
	err := repo.Set(t.Context(), tx, existingInstanceSettings)
	require.NoError(t, err)

	existingOrganizationSettings := &domain.LockoutSettings{
		Settings: domain.Settings{
			InstanceID:     instanceID,
			OrganizationID: gu.Ptr(orgID),
		},
		LockoutSettingsAttributes: domain.LockoutSettingsAttributes{
			MaxPasswordAttempts: gu.Ptr(uint64(4)),
			MaxOTPAttempts:      gu.Ptr(uint64(4)),
			ShowLockOutFailures: gu.Ptr(true),
		},
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
			condition:        repo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(repo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        repo.PrimaryKeyCondition(instanceID, "foo"),
			wantRowsAffected: 0,
		},
		{
			name:             "delete instance",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete instance twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingInstanceSettings.ID),
			wantRowsAffected: 0,
		},
		{
			name:             "delete organization",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete organization twice",
			condition:        repo.PrimaryKeyCondition(instanceID, existingOrganizationSettings.ID),
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
