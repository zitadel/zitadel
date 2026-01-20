package repository_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func Test_user_Get(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()
	// humanRepo := repository.HumanUserRepository()
	// machineRepo := repository.MachineUserRepository()

	instanceID := createInstance(t, tx)

	orgID1 := createOrganization(t, tx, instanceID)
	orgID2 := createOrganization(t, tx, instanceID)

	human := &domain.User{
		ID:             gofakeit.UUID(),
		InstanceID:     instanceID,
		OrganizationID: orgID1,
		Username:       "human-user",
		State:          domain.UserStateActive,
		Human: &domain.HumanUser{
			FirstName:         "John",
			LastName:          "Doe",
			DisplayName:       "JohnD",
			Nickname:          "johnny",
			PreferredLanguage: language.English,
			Gender:            domain.HumanGenderMale,
			AvatarKey:         "https://my.avatar/key",
			Email: domain.HumanEmail{
				Address:    "john@doe.com",
				VerifiedAt: time.Now(),
			},
		},
	}

	machine := &domain.User{
		ID:             gofakeit.UUID(),
		InstanceID:     instanceID,
		OrganizationID: orgID2,
		Username:       "machine-user",
		State:          domain.UserStateActive,
		Machine: &domain.MachineUser{
			Name:        "My Machine",
			Description: "This is my machine user",
		},
	}

	err := userRepo.Create(t.Context(), tx, human)
	require.NoError(t, err)
	err = userRepo.Create(t.Context(), tx, machine)
	require.NoError(t, err)

	tests := []struct {
		name      string
		setup     func(t *testing.T, tx database.Transaction) error
		condition database.Condition
		wantErr   error
		expected  *domain.User
	}{
		{
			name:      "incomplete condition",
			condition: userRepo.IDCondition(human.ID),
			wantErr:   database.NewMissingConditionError(userRepo.InstanceIDColumn()),
		},
		// {
		// 	name: "not found",
		// },
		// {
		// 	name: "too many",
		// },
		// {
		// 	name: "get human",
		// },
		// {
		// 	name: "get machine",
		// },
		// {
		// 	name: "wrong type results in not found",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			t.Cleanup(func() {
				err := savepoint.Rollback(context.Background())
				if err != nil {
					t.Log("rollback savepoint failed", err)
				}
			})
			if tt.setup != nil {
				err := tt.setup(t, savepoint)
				require.NoError(t, err)
			}

			user, err := userRepo.Get(t.Context(), savepoint, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.expected, user)
		})
	}
}

func assertUser(t *testing.T, expected, actual *domain.User) {
	t.Helper()
	assert.Equal(t, expected.InstanceID, actual.InstanceID)
	assert.Equal(t, expected.OrganizationID, actual.OrganizationID)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.Equal(t, expected.State, actual.State)
	assert.True(t, expected.CreatedAt.Equal(actual.CreatedAt), "created at")
	assert.True(t, expected.UpdatedAt.Equal(actual.UpdatedAt), "updated at")

	require.Equal(t, expected.Machine == nil, actual.Machine == nil, "machine check")
	require.Equal(t, expected.Human == nil, actual.Human == nil, "human check")

	require.Len(t, actual.Metadata, len(expected.Metadata), "metadata")
	for i, expectedMetadata := range expected.Metadata {
		var actualMetadata *domain.UserMetadata
		actual.Metadata = slices.DeleteFunc(actual.Metadata, func(m *domain.UserMetadata) bool {
			if m.Key == expectedMetadata.Key {
				actualMetadata = m
				return true
			}
			return false
		})
		require.NotNil(t, actualMetadata, "metadata[%d] not found", i)
		assert.Equal(t, expectedMetadata.Key, actualMetadata.Key, "metadata[%d] key", i)
		assert.Equal(t, expectedMetadata.Value, actualMetadata.Value, "metadata[%d] value", i)
		assert.True(t, expectedMetadata.CreatedAt.Equal(actualMetadata.CreatedAt), "metadata[%d] created at", i)
		assert.True(t, expectedMetadata.UpdatedAt.Equal(actualMetadata.UpdatedAt), "metadata[%d] updated at", i)
	}
	assert.Empty(t, actual.Metadata, "unmatched metadata")

	if expected.Machine != nil {
		require.Nil(t, expected.Human)
		assertMachineUser(t, expected.Machine, actual.Machine)
		return
	}

	if expected.Human != nil {
		require.Nil(t, expected.Machine)
		assertHumanUser(t, expected.Human, actual.Human)
		return
	}
	t.Log("either machine or human must be set")
	t.Fail()
}

func assertMachineUser(t *testing.T, expected, actual *domain.MachineUser) {
	t.Helper()

	assert.Equal(t, expected.Name, actual.Name, "machine name")
	assert.Equal(t, expected.Description, actual.Description, "machine description")
	assert.Equal(t, expected.AccessTokenType, actual.AccessTokenType, "machine access token type")
	assert.Equal(t, expected.Secret, actual.Secret, "machine secret")
	require.Len(t, actual.PATs, len(expected.PATs), "machine pats")
	for i, expectedPAT := range expected.PATs {
		var actualPAT *domain.PersonalAccessToken
		actual.PATs = slices.DeleteFunc(actual.PATs, func(p *domain.PersonalAccessToken) bool {
			if p.ID == expectedPAT.ID {
				actualPAT = p
				return true
			}
			return false
		})
		require.NotNil(t, actualPAT, "machine pat[%d] not found", i)
		assert.Equal(t, expectedPAT.ID, actualPAT.ID, "machine pat[%d] id", i)
		assert.True(t, expectedPAT.CreatedAt.Equal(actualPAT.CreatedAt), "machine pat[%d] created at", i)
		assert.True(t, expectedPAT.ExpiresAt.Equal(actualPAT.ExpiresAt), "machine pat[%d] expires at", i)
		assert.Equal(t, expectedPAT.Scopes, actualPAT.Scopes, "machine pat[%d] scopes", i)
	}
	assert.Empty(t, actual.PATs, "unmatched machine pats")
	require.Len(t, actual.Keys, len(expected.Keys), "machine keys")
	for i, expectedKey := range expected.Keys {
		var actualKey *domain.MachineKey
		actual.Keys = slices.DeleteFunc(actual.Keys, func(k *domain.MachineKey) bool {
			if k.ID == expectedKey.ID {
				actualKey = k
				return true
			}
			return false
		})
		require.NotNil(t, actualKey, "machine key[%d] not found", i)
		assert.Equal(t, expectedKey.ID, actualKey.ID, "machine key[%d] id", i)
		assert.Equal(t, expectedKey.PublicKey, actualKey.PublicKey, "machine key[%d] public key", i)
		assert.True(t, expectedKey.CreatedAt.Equal(actualKey.CreatedAt), "machine key[%d] created at", i)
		assert.True(t, expectedKey.ExpiresAt.Equal(actualKey.ExpiresAt), "machine key[%d] expires at", i)
		assert.Equal(t, expectedKey.Type, actualKey.Type, "machine key[%d] type", i)
	}
	assert.Empty(t, actual.Keys, "unmatched machine keys")
}

func assertHumanUser(t *testing.T, expected, actual *domain.HumanUser) {
	t.Helper()

	assert.Equal(t, expected.FirstName, actual.FirstName, "human first name")
	assert.Equal(t, expected.LastName, actual.LastName, "human last name")
	assert.Equal(t, expected.Nickname, actual.Nickname, "human nickname")
	assert.Equal(t, expected.DisplayName, actual.DisplayName, "human display name")
	assert.Equal(t, expected.AvatarKey, actual.AvatarKey, "human avatar key")
	assert.Equal(t, expected.PreferredLanguage, actual.PreferredLanguage, "human preferred language")
	assert.Equal(t, expected.Gender, actual.Gender, "human gender")
	assert.True(t, expected.MultifactorInitializationSkippedAt.Equal(actual.MultifactorInitializationSkippedAt), "human multifactor initialization skipped at")

	assert.Equal(t, expected.Email.Address, actual.Email.Address, "human email address")
	assert.True(t, expected.Email.VerifiedAt.Equal(actual.Email.VerifiedAt), "human email verified at")
	assert.True(t, expected.Email.OTP.EnabledAt.Equal(actual.Email.OTP.EnabledAt), "human email otp enabled at")
	assert.True(t, expected.Email.OTP.LastSuccessfullyCheckedAt.Equal(actual.Email.OTP.LastSuccessfullyCheckedAt), "human email otp last successfully checked at")
	if expected.Email.OTP.Check != nil {
		t.Error("human.email.otp.check not asserted")
	}
	assertVerification(t, expected.Email.Unverified, actual.Email.Unverified, "human email unverified")

	if expected.Phone == nil {
		assert.Nil(t, actual.Phone, "human phone nil")
	}
	if expected.Phone != nil {
		require.NotNil(t, actual.Phone, "human phone not nil")
		assert.Equal(t, expected.Phone.Number, actual.Phone.Number, "human phone number")
		assert.True(t, expected.Phone.VerifiedAt.Equal(actual.Phone.VerifiedAt), "human phone verified at")
		assert.True(t, expected.Phone.OTP.EnabledAt.Equal(actual.Phone.OTP.EnabledAt), "human phone otp enabled at")
		assert.True(t, expected.Phone.OTP.LastSuccessfullyCheckedAt.Equal(actual.Phone.OTP.LastSuccessfullyCheckedAt), "human phone otp last successfully checked at")
		if expected.Phone.OTP.Check != nil {
			t.Error("human.phone.otp.check not asserted")
		}
		assertVerification(t, expected.Phone.Unverified, actual.Phone.Unverified, "human phone unverified")
	}

	for i, expectedPasskey := range expected.Passkeys {
		var actualPasskey *domain.Passkey
		actual.Passkeys = slices.DeleteFunc(actual.Passkeys, func(p *domain.Passkey) bool {
			if p.ID == expectedPasskey.ID {
				actualPasskey = p
				return true
			}
			return false
		})
		require.NotNil(t, actualPasskey, "human passkey[%d] not found", i)
		assert.Equal(t, expectedPasskey.ID, actualPasskey.ID, "human passkey[%d] id", i)
		assert.Equal(t, expectedPasskey.KeyID, actualPasskey.KeyID, "human passkey[%d] key id", i)
		assert.Equal(t, expectedPasskey.Type, actualPasskey.Type, "human passkey[%d] type", i)
		assert.Equal(t, expectedPasskey.Name, actualPasskey.Name, "human passkey[%d] name", i)
		assert.Equal(t, expectedPasskey.PublicKey, actualPasskey.PublicKey, "human passkey[%d] public key", i)
		assert.Equal(t, expectedPasskey.AttestationType, actualPasskey.AttestationType, "human passkey[%d] attestation type", i)
		assert.Equal(t, expectedPasskey.AuthenticatorAttestationGUID, actualPasskey.AuthenticatorAttestationGUID, "human passkey[%d] authenticator attestation guid", i)
		assert.Equal(t, expectedPasskey.RelyingPartyID, actualPasskey.RelyingPartyID, "human passkey[%d] relying party id", i)
		assert.True(t, expectedPasskey.CreatedAt.Equal(actualPasskey.CreatedAt), "human passkey[%d] created at", i)
		assert.True(t, expectedPasskey.UpdatedAt.Equal(actualPasskey.UpdatedAt), "human passkey[%d] updated at", i)
		assert.True(t, expectedPasskey.VerifiedAt.Equal(actualPasskey.VerifiedAt), "human passkey[%d] verified at", i)
		assert.Equal(t, expectedPasskey.SignCount, actualPasskey.SignCount, "human passkey[%d] sign count", i)
		assert.Equal(t, expectedPasskey.Challenge, actualPasskey.Challenge, "human passkey[%d] challenge", i)
	}
	assert.Empty(t, actual.Passkeys, "unmatched human passkeys")

	assert.Equal(t, expected.Password.Password, actual.Password.Password, "human password password")
	assert.Equal(t, expected.Password.IsChangeRequired, actual.Password.IsChangeRequired, "human password is change required")
	assert.True(t, expected.Password.VerifiedAt.Equal(actual.Password.VerifiedAt), "human password verified at")
	assertVerification(t, expected.Password.Unverified, actual.Password.Unverified, "human password unverified")
	assert.Equal(t, expected.Password.FailedAttempts, actual.Password.FailedAttempts, "human password failed attempts")

	assert.True(t, expected.TOTP.VerifiedAt.Equal(actual.TOTP.VerifiedAt), "human totp verified at")
	assert.True(t, expected.TOTP.LastSuccessfullyCheckedAt.Equal(actual.TOTP.LastSuccessfullyCheckedAt), "human totp last successfully checked at")
	if expected.TOTP.Check != nil {
		require.NotNil(t, actual.TOTP.Check, "human totp check not found")
		assert.Equal(t, expected.TOTP.Check.Code, actual.TOTP.Check.Code, "human totp check code")
		if expected.TOTP.Check.ExpiresAt != nil {
			require.NotNil(t, actual.TOTP.Check.ExpiresAt, "human totp check expires at not nil")
			assert.True(t, expected.TOTP.Check.ExpiresAt.Equal(*actual.TOTP.Check.ExpiresAt), "human totp check expires at")
		} else {
			assert.Nil(t, actual.TOTP.Check.ExpiresAt, "human totp check expires at nil")
		}
		assert.Equal(t, expected.TOTP.Check.FailedAttempts, actual.TOTP.Check.FailedAttempts, "human totp check failed attempts")
	} else {
		assert.Nil(t, actual.TOTP.Check, "human totp check nil")
	}
	assertVerification(t, expected.TOTP.Unverified, actual.TOTP.Unverified, "human totp unverified")

	if len(expected.IdentityProviderLinks) > 0 {
		t.Error("human.identityProviders not asserted")
	}

	for i, expectedVerification := range expected.Verifications {
		var actualVerification *domain.Verification
		actual.Verifications = slices.DeleteFunc(actual.Verifications, func(v *domain.Verification) bool {
			if v.ID == expectedVerification.ID {
				actualVerification = v
				return true
			}
			return false
		})
		require.NotNil(t, actualVerification, "human verification[%d] not found", i)
		assert.Equal(t, expectedVerification.ID, actualVerification.ID, "human verification[%d] id", i)
		if expectedVerification.Value != nil {
			require.NotNil(t, actualVerification.Value, "human verification[%d] value not nil", i)
			assert.Equal(t, *expectedVerification.Value, *actualVerification.Value, "human verification[%d] value", i)
		} else {
			assert.Nil(t, actualVerification.Value, "human verification[%d] value nil", i)
		}
		assert.Equal(t, expectedVerification.Code, actualVerification.Code, "human verification[%d] code", i)
		if expectedVerification.ExpiresAt != nil {
			require.NotNil(t, actualVerification.ExpiresAt, "human verification[%d] expires at not nil", i)
			assert.True(t, expectedVerification.ExpiresAt.Equal(*actualVerification.ExpiresAt), "human verification[%d] expires at", i)
		} else {
			assert.Nil(t, actualVerification.ExpiresAt, "human verification[%d] expires at nil", i)
		}
		assert.Equal(t, expectedVerification.FailedAttempts, actualVerification.FailedAttempts, "human verification[%d] failed attempts", i)
		assert.True(t, expectedVerification.VerifiedAt.Equal(actualVerification.VerifiedAt), "human verification[%d] verified at", i)
	}
	assert.Empty(t, actual.Verifications, "unmatched human verifications")
}

func assertVerification(t *testing.T, expected, actual *domain.Verification, field string) {
	t.Helper()

	if expected == nil {
		assert.Nil(t, actual, "verification nil")
		return
	}

	require.NotNil(t, actual, "human phone unverified not nil")
	assert.Equal(t, expected.Value, actual.Value, "human phone unverified value")
	assert.Equal(t, expected.Code, actual.Code, "human phone unverified code")
	if expected.ExpiresAt != nil {
		require.NotNil(t, actual.ExpiresAt, "human phone unverified expires at not nil")
		assert.True(t, expected.ExpiresAt.Equal(*actual.ExpiresAt), "human phone unverified expires at")
	} else {
		assert.Nil(t, actual.ExpiresAt, "human phone unverified expires at nil")
	}
	assert.Equal(t, expected.FailedAttempts, actual.FailedAttempts, "human phone unverified failed attempts")
	assert.True(t, expected.VerifiedAt.Equal(actual.VerifiedAt), "human phone unverified verified at")
}
