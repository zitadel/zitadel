package repository_test

import (
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
			want:      []*domain.LoginSettings{settings[0], settings[2]},
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
				database.WithOrderByAscending(repo.PrimaryKeyColumns()...),
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
