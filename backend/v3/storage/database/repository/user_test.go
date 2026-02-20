package repository_test

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/crypto"
)

func Test_user_Get(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()

	instanceID := createInstance(t, tx)

	orgID1 := createOrganization(t, tx, instanceID)
	orgID2 := createOrganization(t, tx, instanceID)

	now := time.Now().Round(time.Millisecond)

	human := &domain.User{
		ID:             gofakeit.UUID(),
		InstanceID:     instanceID,
		OrganizationID: orgID1,
		Username:       "human-user",
		State:          domain.UserStateActive,
		CreatedAt:      now,
		UpdatedAt:      now,
		Human: &domain.HumanUser{
			FirstName:         "John",
			LastName:          "Doe",
			DisplayName:       "JohnD",
			Nickname:          "johnny",
			PreferredLanguage: language.English,
			Gender:            domain.HumanGenderMale,
			AvatarKey:         "https://my.avatar/key",
			Email: domain.HumanEmail{
				Address:           "john@doe.com",
				UnverifiedAddress: "john@doe.com",
				VerifiedAt:        now,
			},
		},
	}

	machine := &domain.User{
		ID:             gofakeit.UUID(),
		InstanceID:     instanceID,
		OrganizationID: orgID2,
		Username:       "machine-user",
		State:          domain.UserStateActive,
		CreatedAt:      now,
		UpdatedAt:      now,
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
		{
			name:      "not found",
			condition: userRepo.PrimaryKeyCondition(instanceID, "not-found"),
			wantErr:   new(database.NoRowFoundError),
		},
		{
			name:      "too many",
			condition: userRepo.InstanceIDCondition(instanceID),
			wantErr:   new(database.MultipleRowsFoundError),
		},
		{
			name:      "get human",
			condition: userRepo.PrimaryKeyCondition(instanceID, human.ID),
			expected:  human,
		},
		{
			name:      "get machine",
			condition: userRepo.PrimaryKeyCondition(instanceID, machine.ID),
			expected:  machine,
		},
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

			require.Equal(t, tt.expected == nil, user == nil)
			if tt.expected == nil {
				return
			}
			assertUser(t, tt.expected, user)
		})
	}
}

func Test_user_List(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()

	instanceID := createInstance(t, tx)

	orgID1 := createOrganization(t, tx, instanceID)
	orgID2 := createOrganization(t, tx, instanceID)

	idpID := createIdentityProvider(t, tx, instanceID, orgID1)

	now := time.Now().Round(time.Millisecond)

	human1 := &domain.User{
		ID:             "h1",
		InstanceID:     instanceID,
		OrganizationID: orgID1,
		Username:       "human-user-1",
		State:          domain.UserStateActive,
		CreatedAt:      now,
		UpdatedAt:      now,
		Human: &domain.HumanUser{
			FirstName:         "John",
			LastName:          "Doe",
			DisplayName:       "JohnD",
			Nickname:          "johnny",
			PreferredLanguage: language.English,
			Gender:            domain.HumanGenderMale,
			AvatarKey:         "https://my.avatar/key",
			Email: domain.HumanEmail{
				Address:           "john@doe.com",
				UnverifiedAddress: "john@doe.com",
				VerifiedAt:        now,
			},
			Phone: &domain.HumanPhone{
				Number:           "+1234567890",
				UnverifiedNumber: "+1234567890",
				VerifiedAt:       now,
			},
			Passkeys: []*domain.Passkey{
				{
					ID:                           "passkey-1",
					CreatedAt:                    now,
					UpdatedAt:                    now,
					VerifiedAt:                   now,
					Challenge:                    []byte("challenge2"),
					KeyID:                        []byte("key-id2"),
					Name:                         "u2f",
					Type:                         domain.PasskeyTypeU2F,
					PublicKey:                    []byte("public key"),
					AuthenticatorAttestationGUID: []byte("aaguid"),
					AttestationType:              "attestation type",
					RelyingPartyID:               "rpid",
				},
			},
			Password: domain.HumanPassword{
				Hash:      "password-hash",
				ChangedAt: now,
			},
			TOTP: &domain.HumanTOTP{
				VerifiedAt: now,
				Secret: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "aes-256",
					KeyID:      "id",
					Crypted:    []byte("secret"),
				},
			},
			IdentityProviderLinks: []*domain.IdentityProviderLink{
				{
					ProviderID:       idpID,
					ProvidedUserID:   "provided-id-h1",
					ProvidedUsername: "provided-username-h1",
					CreatedAt:        now,
					UpdatedAt:        now,
				},
			},
			Verifications: []*domain.Verification{
				{
					ID: "h1-v1",
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeHash,
						Algorithm:  "aes-256",
						KeyID:      "key-id",
						Crypted:    []byte("asdf"),
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
		Metadata: []*domain.Metadata{
			{
				InstanceID: instanceID,
				Key:        "key1",
				Value:      []byte("asdf"),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				InstanceID: instanceID,
				Key:        "key2",
				Value:      []byte("123"),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
	}
	human2 := &domain.User{
		ID:             "h2",
		InstanceID:     instanceID,
		OrganizationID: orgID1,
		Username:       "human-user-2",
		State:          domain.UserStateActive,
		CreatedAt:      now,
		UpdatedAt:      now,
		Human: &domain.HumanUser{
			FirstName:         "Jane",
			LastName:          "Doe",
			DisplayName:       "JaneD",
			Nickname:          "janie",
			PreferredLanguage: language.English,
			Gender:            domain.HumanGenderFemale,
			AvatarKey:         "https://my.avatar/key",
			Email: domain.HumanEmail{
				UnverifiedAddress: "jane@doe.com",
				PendingVerification: &domain.Verification{
					ID:        "h2-email",
					CreatedAt: now,
					UpdatedAt: now,
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeHash,
						Algorithm:  "aes-256",
						KeyID:      "key-id",
						Crypted:    []byte("asdf"),
					},
				},
			},
			Phone: &domain.HumanPhone{
				UnverifiedNumber: "+1234567890",
				PendingVerification: &domain.Verification{
					ID:        "h2-phone",
					CreatedAt: now,
					UpdatedAt: now,
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeHash,
						Algorithm:  "aes-256",
						KeyID:      "key-id",
						Crypted:    []byte("asdf"),
					},
				},
			},
			Passkeys: []*domain.Passkey{
				{
					ID:                           "passkey-2",
					CreatedAt:                    now,
					UpdatedAt:                    now,
					VerifiedAt:                   now,
					Challenge:                    []byte("challenge2"),
					KeyID:                        []byte("key-id3"),
					Name:                         "u2f",
					Type:                         domain.PasskeyTypeU2F,
					PublicKey:                    []byte("public key"),
					AuthenticatorAttestationGUID: []byte("aaguid"),
					AttestationType:              "attestation type",
					RelyingPartyID:               "rpid",
				},
			},
			Password: domain.HumanPassword{
				PendingVerification: &domain.Verification{
					ID:        "h2-pw",
					CreatedAt: now,
					UpdatedAt: now,
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeHash,
						Algorithm:  "aes-256",
						KeyID:      "key-id",
						Crypted:    []byte("asdf"),
					},
				},
			},
			TOTP: &domain.HumanTOTP{
				VerifiedAt: now,
				Secret: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "aes-256",
					KeyID:      "id",
					Crypted:    []byte("secret"),
				},
			},
			IdentityProviderLinks: []*domain.IdentityProviderLink{
				{
					ProviderID:       idpID,
					ProvidedUserID:   "provided-id-h2",
					ProvidedUsername: "provided-username-h2",
					CreatedAt:        now,
					UpdatedAt:        now,
				},
			},
			Verifications: []*domain.Verification{
				{
					ID: "h2-v1",
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeHash,
						Algorithm:  "aes-256",
						KeyID:      "key-id",
						Crypted:    []byte("asdf"),
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
		Metadata: []*domain.Metadata{
			{
				InstanceID: instanceID,
				Key:        "key1",
				Value:      []byte("asdf"),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				InstanceID: instanceID,
				Key:        "key2",
				Value:      []byte("123"),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
	}

	machine1 := &domain.User{
		ID:             "m1",
		InstanceID:     instanceID,
		OrganizationID: orgID2,
		Username:       "machine-user-1",
		State:          domain.UserStateActive,
		CreatedAt:      now,
		UpdatedAt:      now,
		Machine: &domain.MachineUser{
			Name:            "My Machine",
			Description:     "This is my machine user",
			Secret:          "asdf",
			AccessTokenType: domain.AccessTokenTypeJWT,
			PATs: []*domain.PersonalAccessToken{
				{
					ID:        "pat1",
					CreatedAt: now,
					Scopes:    []string{"scope1", "scope2"},
				},
			},
			Keys: []*domain.MachineKey{
				{
					ID:        "mk1",
					PublicKey: []byte("asdf"),
					CreatedAt: now,
					Type:      domain.MachineKeyTypeJSON,
				},
			},
		},
		Metadata: []*domain.Metadata{
			{
				InstanceID: instanceID,
				Key:        "key1",
				Value:      []byte("asdf"),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				InstanceID: instanceID,
				Key:        "key2",
				Value:      []byte("123"),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
	}

	machine2 := &domain.User{
		ID:             "m2",
		InstanceID:     instanceID,
		OrganizationID: orgID2,
		Username:       "machine-user-2",
		State:          domain.UserStateActive,
		CreatedAt:      now,
		UpdatedAt:      now,
		Machine: &domain.MachineUser{
			Name:            "My Machine 2",
			Description:     "This is my machine user 2",
			Secret:          "asdf",
			AccessTokenType: domain.AccessTokenTypeJWT,
			PATs: []*domain.PersonalAccessToken{
				{
					ID:        "pat2",
					CreatedAt: now,
					Scopes:    []string{"scope1", "scope2"},
				},
			},
			Keys: []*domain.MachineKey{
				{
					ID:        "mk2",
					PublicKey: []byte("asdf"),
					CreatedAt: now,
					Type:      domain.MachineKeyTypeJSON,
				},
			},
		},
		Metadata: []*domain.Metadata{
			{
				InstanceID: instanceID,
				Key:        "key1",
				Value:      []byte("asdf"),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				InstanceID: instanceID,
				Key:        "key2",
				Value:      []byte("123"),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
	}

	err := userRepo.Create(t.Context(), tx, human1)
	require.NoError(t, err)
	err = userRepo.Create(t.Context(), tx, human2)
	require.NoError(t, err)
	err = userRepo.Create(t.Context(), tx, machine1)
	require.NoError(t, err)
	err = userRepo.Create(t.Context(), tx, machine2)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		wantErr   error
		expected  []*domain.User
	}{
		{
			name:      "incomplete condition",
			condition: userRepo.IDCondition(human1.ID),
			wantErr:   database.NewMissingConditionError(userRepo.InstanceIDColumn()),
		},
		{
			name:      "get human 1",
			condition: userRepo.PrimaryKeyCondition(instanceID, human1.ID),
			expected:  []*domain.User{human1},
		},
		{
			name: "get humans",
			condition: database.And(
				userRepo.InstanceIDCondition(instanceID),
				userRepo.TypeCondition(domain.UserTypeHuman),
			),
			expected: []*domain.User{human1, human2},
		},
		{
			name:      "get machine 1",
			condition: userRepo.PrimaryKeyCondition(instanceID, machine1.ID),
			expected:  []*domain.User{machine1},
		},
		{
			name: "get machines",
			condition: database.And(
				userRepo.InstanceIDCondition(instanceID),
				userRepo.TypeCondition(domain.UserTypeMachine),
			),
			expected: []*domain.User{machine1, machine2},
		},
		{
			name:      "get all",
			condition: userRepo.InstanceIDCondition(instanceID),
			expected:  []*domain.User{human1, human2, machine1, machine2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotten, err := userRepo.List(t.Context(), tx, database.WithCondition(tt.condition), database.WithOrderByAscending(userRepo.UsernameColumn()))
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr != nil {
				return
			}

			assertUsers(t, tt.expected, gotten)
		})
	}
}

// Test_user_ListConditions verifies the correct behavior of all conditions.
func Test_user_ListConditions(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()
	humanRepo := repository.HumanUserRepository()

	instanceID1 := createInstance(t, tx)
	instance1OrgID1 := createOrganization(t, tx, instanceID1)
	instance1OrgID2 := createOrganization(t, tx, instanceID1)
	idpID := createIdentityProvider(t, tx, instanceID1, instance1OrgID1)

	instanceID2 := createInstance(t, tx)
	instance2OrgID1 := createOrganization(t, tx, instanceID2)

	now := time.Now().Round(time.Millisecond)

	users := []*domain.User{
		{
			ID:             "h1",
			InstanceID:     instanceID1,
			OrganizationID: instance1OrgID1,
			Username:       "human-user-1",
			State:          domain.UserStateActive,
			Human: &domain.HumanUser{
				FirstName:         "aa",
				LastName:          "bb",
				DisplayName:       "aa bb",
				Nickname:          "aabb",
				PreferredLanguage: language.English,
				Gender:            domain.HumanGenderDiverse,
				AvatarKey:         "https://my.avatar/key",
				Email: domain.HumanEmail{
					Address:    "aa.bb@example.com",
					VerifiedAt: now,
				},
			},
		},
		{
			ID:             "h2",
			InstanceID:     instanceID1,
			OrganizationID: instance1OrgID2,
			Username:       "human-user-2",
			State:          domain.UserStateActive,
			Human: &domain.HumanUser{
				FirstName:         "aaa",
				LastName:          "bbb",
				DisplayName:       "aaa bbb",
				Nickname:          "aaabbb",
				PreferredLanguage: language.German,
				Gender:            domain.HumanGenderFemale,
				Email: domain.HumanEmail{
					Address: "aaa@bbb.com",
				},
			},
		},
		{
			ID:             "h3",
			InstanceID:     instanceID2,
			OrganizationID: instance2OrgID1,
			Username:       "human-user-3",
			State:          domain.UserStateActive,
			Human: &domain.HumanUser{
				FirstName:   "aaaa",
				LastName:    "bbbb",
				DisplayName: "aaaa bbbb",
				Email: domain.HumanEmail{
					Address: "aaaa.bbbb@example.com",
				},
			},
		},
		{
			ID:             "m1",
			InstanceID:     instanceID1,
			OrganizationID: instance1OrgID1,
			Username:       "machine-user-1",
			State:          domain.UserStateActive,
			Machine: &domain.MachineUser{
				Name:        "machine-user-1",
				Description: "This is my machine user",
			},
		},
		{
			ID:             "m2",
			InstanceID:     instanceID1,
			OrganizationID: instance1OrgID2,
			Username:       "machine-user-2",
			State:          domain.UserStateActive,
			Machine: &domain.MachineUser{
				Name: "machine-user-2",
			},
		},
		{
			ID:             "m3",
			InstanceID:     instanceID2,
			OrganizationID: instance2OrgID1,
			Username:       "machine-user-3",
			State:          domain.UserStateActive,
			Machine: &domain.MachineUser{
				Name: "machine-user-3",
			},
		},
	}

	for _, user := range users {
		err := userRepo.Create(t.Context(), tx, user)
		require.NoError(t, err)
	}

	type want struct {
		err error
		ids []string
	}
	tests := []struct {
		name  string
		setup func(t *testing.T, tx database.Transaction)
		opts  []database.QueryOption
		want  want
	}{
		// user conditions
		{
			name: "by instance id",
			opts: []database.QueryOption{
				database.WithCondition(userRepo.InstanceIDCondition(instanceID1)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "h2", "m1", "m2"},
			},
		},
		{
			name: "by primary key",
			opts: []database.QueryOption{
				database.WithCondition(userRepo.PrimaryKeyCondition(instanceID1, "h1")),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},
		{
			name: "by organization id",
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					userRepo.OrganizationIDCondition(instance1OrgID1),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "m1"},
			},
		},
		{
			name: "by id",
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					userRepo.IDCondition("h1"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},
		{
			name: "by username",
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					userRepo.UsernameCondition(database.TextOperationStartsWith, "human-user-"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "h2"},
			},
		},
		{
			name: "by state",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					userRepo.SetState(domain.UserStateInactive),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					userRepo.SetState(domain.UserStateInitial),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "m1"),
					userRepo.SetState(domain.UserStateLocked),
				)
				require.NoError(t, err)

			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					userRepo.StateCondition(domain.UserStateActive),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"m2"},
			},
		},
		{
			name: "by type",
			setup: func(t *testing.T, tx database.Transaction) {
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					userRepo.TypeCondition(domain.UserTypeMachine),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"m1", "m2"},
			},
		},
		// {
		// 	name: "by login name",
		// 	setup: func(t *testing.T, tx database.Transaction) {
		// 	},
		// 	opts: []database.QueryOption{
		// 	database.WithCondition(database.And(
		// 		userRepo.InstanceIDCondition(instanceID1),
		// 	)),
		// },
		// 	want: want{
		// 		err: nil,
		// 		ids: []string{},
		// 	},
		// },

		// metadata conditions
		{
			name: "by metadata key",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key1",
							Value:      []byte("value"),
						},
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key2",
							Value:      []byte("value"),
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key2",
							Value:      []byte("value"),
						},
					),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					userRepo.ExistsMetadata(
						userRepo.MetadataConditions().MetadataKeyCondition(database.TextOperationEqual, "key1"),
					),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},
		{
			name: "by metadata value",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key1",
							Value:      []byte("123"),
						},
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key2",
							Value:      []byte("value"),
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key2",
							Value:      []byte("value1"),
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID2, "h3"),
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID2,
							Key:        "key2",
							Value:      []byte("value"),
						},
					),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					userRepo.ExistsMetadata(
						userRepo.MetadataConditions().MetadataValueCondition(database.BytesOperationEqual, []byte("value")),
					),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},
		{
			name: "by metadata key and value",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key1",
							Value:      []byte("123"),
						},
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key2",
							Value:      []byte("value"),
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID1,
							Key:        "key2",
							Value:      []byte("value1"),
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID2, "h3"),
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID2,
							Key:        "key1",
							Value:      []byte("value"),
						},
					),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					userRepo.ExistsMetadata(
						database.And(
							userRepo.MetadataConditions().MetadataKeyCondition(database.TextOperationEqual, "key1"),
							userRepo.MetadataConditions().MetadataValueCondition(database.BytesOperationEqual, []byte("123")),
						),
					),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},

		// human conditions
		{
			name: "by display name",
			setup: func(t *testing.T, tx database.Transaction) {
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.DisplayNameCondition(database.TextOperationEndsWith, " bbb"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h2"},
			},
		},
		{
			name: "by first name",
			setup: func(t *testing.T, tx database.Transaction) {
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.FirstNameCondition(database.TextOperationEqual, "aa"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},
		{
			name: "by last name",
			setup: func(t *testing.T, tx database.Transaction) {
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.LastNameCondition(database.TextOperationNotEqual, "bb"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h2"},
			},
		},
		{
			name: "by nickname",
			setup: func(t *testing.T, tx database.Transaction) {
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.NicknameCondition(database.TextOperationContains, "aabb"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "h2"},
			},
		},
		{
			name: "by preferred language",
			setup: func(t *testing.T, tx database.Transaction) {
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.PreferredLanguageCondition(language.German),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h2"},
			},
		},

		// phone conditions
		{
			name: "by phone",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{}),
				)
				require.NoError(t, err)
				_, err = humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.SetUnverifiedPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeInit{
						Expiry: gu.Ptr(5 * time.Minute),
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.PhoneCondition(database.TextOperationEqual, "+1234567890"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "h2"},
			},
		},
		{
			name: "by unverified phone",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{}),
				)
				require.NoError(t, err)
				_, err = humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.SetUnverifiedPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeInit{
						Expiry: gu.Ptr(5 * time.Minute),
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.UnverifiedPhoneCondition(database.TextOperationEqual, "+1234567890"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h2"},
			},
		},
		{
			name: "by verified phone",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{}),
				)
				require.NoError(t, err)
				_, err = humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.SetUnverifiedPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeInit{
						Expiry: gu.Ptr(5 * time.Minute),
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.VerifiedPhoneCondition(database.TextOperationEqual, "+1234567890"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},

		// human email conditions
		{
			name: "by email",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.SetEmail("test@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeSkipped{}),
				)
				require.NoError(t, err)
				_, err = humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.SetUnverifiedEmail("test@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeInit{
						Expiry: gu.Ptr(5 * time.Minute),
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.EmailCondition(database.TextOperationEqual, "test@email.com"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "h2"},
			},
		},
		{
			name: "by unverified email",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.SetEmail("test@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeSkipped{}),
				)
				require.NoError(t, err)
				_, err = humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.SetUnverifiedEmail("test@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeInit{
						Expiry: gu.Ptr(5 * time.Minute),
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.UnverifiedEmailCondition(database.TextOperationEqual, "test@email.com"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h2"},
			},
		},
		{
			name: "by verified email",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.SetEmail("test@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeSkipped{}),
				)
				require.NoError(t, err)
				_, err = humanRepo.Update(t.Context(), tx, humanRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.SetUnverifiedEmail("test@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeInit{
						Expiry: gu.Ptr(5 * time.Minute),
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.VerifiedEmailCondition(database.TextOperationEqual, "test@email.com"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},

		// passkey conditions
		{
			name: "by passkey same challenge",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-1",
							Challenge:                    []byte("challenge1"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-2",
							Challenge:                    []byte("challenge2"),
							KeyID:                        []byte("key-id2"),
							Name:                         "u2f",
							Type:                         domain.PasskeyTypeU2F,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-3",
							Challenge:                    []byte("challenge1"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.ExistsPasskey(humanRepo.PasskeyConditions().ChallengeCondition([]byte("challenge1"))),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "h2"},
			},
		},
		{
			name: "by passkey challenge",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-1",
							Challenge:                    []byte("challenge1"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-2",
							Challenge:                    []byte("challenge2"),
							KeyID:                        []byte("key-id2"),
							Name:                         "u2f",
							Type:                         domain.PasskeyTypeU2F,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-3",
							Challenge:                    []byte("challenge3"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.ExistsPasskey(humanRepo.PasskeyConditions().ChallengeCondition([]byte("challenge1"))),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},
		{
			name: "by passkey id",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-1",
							Challenge:                    []byte("challenge1"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-2",
							Challenge:                    []byte("challenge2"),
							KeyID:                        []byte("key-id2"),
							Name:                         "u2f",
							Type:                         domain.PasskeyTypeU2F,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-3",
							Challenge:                    []byte("challenge3"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.ExistsPasskey(humanRepo.PasskeyConditions().IDCondition("passkey-1")),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},
		{
			name: "by passkey key id",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-1",
							Challenge:                    []byte("challenge1"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-2",
							Challenge:                    []byte("challenge2"),
							KeyID:                        []byte("key-id2"),
							Name:                         "u2f",
							Type:                         domain.PasskeyTypeU2F,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-3",
							Challenge:                    []byte("challenge3"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.ExistsPasskey(humanRepo.PasskeyConditions().KeyIDCondition("key-id")),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "h2"},
			},
		},
		{
			name: "by passkey type",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-1",
							Challenge:                    []byte("challenge1"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-2",
							Challenge:                    []byte("challenge2"),
							KeyID:                        []byte("key-id2"),
							Name:                         "u2f",
							Type:                         domain.PasskeyTypeU2F,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.AddPasskey(
						&domain.Passkey{
							ID:                           "passkey-3",
							Challenge:                    []byte("challenge3"),
							KeyID:                        []byte("key-id"),
							Name:                         "passwordless",
							Type:                         domain.PasskeyTypePasswordless,
							PublicKey:                    []byte("public key"),
							AuthenticatorAttestationGUID: []byte("aaguid"),
							AttestationType:              "attestation type",
							RelyingPartyID:               "rpid",
						},
					),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.ExistsPasskey(humanRepo.PasskeyConditions().TypeCondition(domain.PasskeyTypeU2F)),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},

		// idp link conditions
		{
			name: "by idp link provider id",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h1",
						ProvidedUsername: "provided-username-h1",
					}),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h1-2",
						ProvidedUsername: "provided-username-h1-2",
					}),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h2-2",
						ProvidedUsername: "provided-username-h2-2",
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.ExistsIdentityProviderLink(humanRepo.IdentityProviderLinkConditions().ProviderIDCondition(idpID)),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1", "h2"},
			},
		},
		{
			name: "by idp link provided user id",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h1",
						ProvidedUsername: "provided-username-h1",
					}),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h1-2",
						ProvidedUsername: "provided-username-h1-2",
					}),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h2-2",
						ProvidedUsername: "provided-username-h2-2",
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.ExistsIdentityProviderLink(humanRepo.IdentityProviderLinkConditions().ProvidedUserIDCondition("provided-id-h1-2")),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},
		{
			name: "by idp link username",
			setup: func(t *testing.T, tx database.Transaction) {
				_, err := userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h1",
						ProvidedUsername: "provided-username-h1",
					}),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h1"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h1-2",
						ProvidedUsername: "provided-username-h1",
					}),
				)
				require.NoError(t, err)
				_, err = userRepo.Update(t.Context(), tx, userRepo.PrimaryKeyCondition(instanceID1, "h2"),
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       idpID,
						ProvidedUserID:   "provided-id-h2-2",
						ProvidedUsername: "provided-username-h2-2",
					}),
				)
				require.NoError(t, err)
			},
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition(instanceID1),
					humanRepo.ExistsIdentityProviderLink(humanRepo.IdentityProviderLinkConditions().ProvidedUsernameCondition(database.TextOperationStartsWith, "provided-username-h1")),
				)),
			},
			want: want{
				err: nil,
				ids: []string{"h1"},
			},
		},

		{
			name: "list none",
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.InstanceIDCondition("non-existing-instance-id"),
				)),
			},
			want: want{
				err: nil,
				ids: []string{},
			},
		},
		{
			name: "missing instance id condition",
			opts: []database.QueryOption{
				database.WithCondition(database.And(
					userRepo.IDCondition("h1"),
				)),
			},
			want: want{
				err: database.NewMissingConditionError(userRepo.InstanceIDColumn()),
				ids: []string{},
			},
		},
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
				tt.setup(t, savepoint)
			}

			gotten, err := userRepo.List(t.Context(), savepoint, append(tt.opts, database.WithOrderByAscending(userRepo.IDColumn()))...)
			require.ErrorIs(t, err, tt.want.err)
			if tt.want.err != nil {
				return
			}

			require.Len(t, gotten, len(tt.want.ids))
			for i, expectedID := range tt.want.ids {
				assert.Equal(t, expectedID, gotten[i].ID)
			}
		})
	}
}

func Test_user_Update(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()

	humanRepo := repository.HumanUserRepository()
	machineRepo := repository.MachineUserRepository()

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
				Address: "john@doe.com",
			},
		},
	}
	humanCondition := userRepo.PrimaryKeyCondition(instanceID, human.ID)

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
	machineCondition := userRepo.PrimaryKeyCondition(instanceID, machine.ID)

	err := userRepo.Create(t.Context(), tx, human)
	require.NoError(t, err)
	err = userRepo.Create(t.Context(), tx, machine)
	require.NoError(t, err)

	now := time.Now().Round(time.Millisecond)

	type args struct {
		condition database.Condition
		changes   []database.Change
	}
	type want struct {
		err  error
		user *domain.User
	}
	type test struct {
		name  string
		setup func(t *testing.T, tx database.QueryExecutor) error
		args  args
		want  want
	}
	tests := []test{
		{
			name: "set username",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					userRepo.SetUsername("new-human-username"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Username = "new-human-username"
					return u
				}(),
			},
		},
		{
			name: "set state",
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					userRepo.SetState(domain.UserStateInactive),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(machineCondition))
					require.NoError(t, err)
					u.State = domain.UserStateInactive
					return u
				}(),
			},
		},
		{
			name: "add metadata",
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					userRepo.SetMetadata(
						&domain.Metadata{
							InstanceID: instanceID,
							Key:        "key1",
							Value:      []byte("value"),
						},
						&domain.Metadata{
							InstanceID: instanceID,
							Key:        "key2",
							Value:      []byte("42"),
						},
					),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(machineCondition))
					require.NoError(t, err)
					u.Metadata = append(u.Metadata, &domain.Metadata{
						InstanceID: instanceID,
						Key:        "key1",
						Value:      []byte("value"),
					}, &domain.Metadata{
						InstanceID: instanceID,
						Key:        "key2",
						Value:      []byte("42"),
					})
					return u
				}(),
			},
		},
		{
			name: "remove metadata",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, machineCondition,
					userRepo.SetMetadata(&domain.Metadata{
						InstanceID: instanceID,
						Key:        "key1",
						Value:      []byte("value"),
					}),
				)
				return err
			},
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					userRepo.RemoveMetadata(userRepo.MetadataConditions().MetadataKeyCondition(database.TextOperationEqual, "key1")),
				},
			},
			want: want{
				user: machine,
			},
		},
		{
			name: "set human first name",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetFirstName("Joanne"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.FirstName = "Joanne"
					return u
				}(),
			},
		},
		{
			name: "set human last name",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetLastName("Known"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.LastName = "Known"
					return u
				}(),
			},
		},
		{
			name: "set human nickname",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetNickname("Ghost"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Nickname = "Ghost"
					return u
				}(),
			},
		},
		{
			name: "set human display name",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetDisplayName("Johnny the Ghost"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.DisplayName = "Johnny the Ghost"
					return u
				}(),
			},
		},
		{
			name: "set human preferred language",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetPreferredLanguage(language.Afrikaans),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.PreferredLanguage = language.Afrikaans
					return u
				}(),
			},
		},
		{
			name: "set human gender",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetGender(domain.HumanGenderMale),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Gender = domain.HumanGenderMale
					return u
				}(),
			},
		},
		{
			name: "set human avatar key",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetAvatarKey(gu.Ptr("https://new.avatar/key")),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.AvatarKey = "https://new.avatar/key"
					return u
				}(),
			},
		},
		{
			name: "set human skip mfa initialization at",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SkipMultifactorInitializationAt(now),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.MultifactorInitializationSkippedAt = now
					return u
				}(),
			},
		},
		{
			name: "set human verification init",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetVerification(&domain.VerificationTypeInit{
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
						ID: gu.Ptr("verification-id"),
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Verifications = append(u.Human.Verifications, &domain.Verification{
						ID:        "verification-id",
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
					})
					return u
				}(),
			},
		},
		{
			name: "set human verification update",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := humanRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetVerification(&domain.VerificationTypeInit{
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
						ID: gu.Ptr("verification-id"),
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetVerification(&domain.VerificationTypeUpdate{
						ID:     gu.Ptr("verification-id"),
						Expiry: gu.Ptr(10 * time.Minute),
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Verifications = append(u.Human.Verifications, &domain.Verification{
						ID:        "verification-id",
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
						ExpiresAt: gu.Ptr(now.Add(10 * time.Minute)),
					})
					return u
				}(),
			},
		},
		{
			name: "set human verification verified",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := humanRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetVerification(&domain.VerificationTypeInit{
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
						ID: gu.Ptr("verification-id"),
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetVerification(&domain.VerificationTypeSucceeded{
						ID:         gu.Ptr("verification-id"),
						VerifiedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Verifications = nil
					return u
				}(),
			},
		},
		{
			name: "set human verification failed",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := humanRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetVerification(&domain.VerificationTypeInit{
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
						ID: gu.Ptr("verification-id"),
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetVerification(&domain.VerificationTypeFailed{
						ID:       gu.Ptr("verification-id"),
						FailedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Verifications = append(u.Human.Verifications, &domain.Verification{
						ID:        "verification-id",
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
						UpdatedAt:      now,
						FailedAttempts: 1,
					})
					return u
				}(),
			},
		},
		{
			name: "set human password verification init",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetResetPasswordVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("pw-id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedPassword"),
						},
						Expiry: gu.Ptr(24 * time.Hour),
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Password.PendingVerification = &domain.Verification{
						ID:        "pw-id",
						CreatedAt: now,
						UpdatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedPassword"),
						},
						ExpiresAt: gu.Ptr(now.Add(24 * time.Hour)),
					}
					return u
				}(),
			},
		},
		{
			name: "set human password verification update",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetResetPasswordVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("pw-id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedPassword"),
						},
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetResetPasswordVerification(&domain.VerificationTypeUpdate{
						Expiry:    gu.Ptr(48 * time.Hour),
						UpdatedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Password.PendingVerification = &domain.Verification{
						ID:        "pw-id",
						CreatedAt: now,
						UpdatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedPassword"),
						},
						ExpiresAt: gu.Ptr(now.Add(48 * time.Hour)),
					}
					return u
				}(),
			},
		},
		{
			name: "set human password verified",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetResetPasswordVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("pw-id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedPassword"),
						},
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetResetPasswordVerification(&domain.VerificationTypeSucceeded{
						VerifiedAt: now,
					}),
					humanRepo.SetPassword("hashedPassword"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Password = domain.HumanPassword{
						Hash:             "hashedPassword",
						IsChangeRequired: false,
						ChangedAt:        now,
					}
					return u
				}(),
			},
		},
		{
			name: "set human password failed",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetResetPasswordVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("pw-id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedPassword"),
						},
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetResetPasswordVerification(&domain.VerificationTypeFailed{
						FailedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Password.PendingVerification = &domain.Verification{
						ID:        "pw-id",
						CreatedAt: now,
						UpdatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedPassword"),
						},
						FailedAttempts: 1,
					}
					return u
				}(),
			},
		},
		{
			name: "set human password change required",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetPasswordChangeRequired(true),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Password.IsChangeRequired = true
					return u
				}(),
			},
		},
		{
			name: "set human last successful password check",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetLastSuccessfulPasswordCheck(now),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Password.LastSuccessfullyCheckedAt = &now
					return u
				}(),
			},
		},
		{
			name: "set human increment password failed attempts",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.IncrementPasswordFailedAttempts(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Password.FailedAttempts++
					return u
				}(),
			},
		},
		{
			name: "set human reset password failed attempts",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.IncrementPasswordFailedAttempts(),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.ResetPasswordFailedAttempts(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Password.FailedAttempts = 0
					return u
				}(),
			},
		},
		{
			name: "set human set email verification init",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetUnverifiedEmail("new@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("verification-id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Email.UnverifiedAddress = "new@email.com"
					u.Human.Email.PendingVerification = &domain.Verification{
						ID:        "verification-id",
						CreatedAt: now,
						UpdatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human set email verified",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				rowCount, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetUnverifiedEmail("new@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("pw-id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
					}),
				)
				assert.EqualValues(t, 1, rowCount)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetEmailVerification(&domain.VerificationTypeSucceeded{
						VerifiedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Email.PendingVerification = nil
					u.Human.Email.Address = "new@email.com"
					u.Human.Email.UnverifiedAddress = "new@email.com"
					u.Human.Email.VerifiedAt = now
					return u
				}(),
			},
		},
		{
			name: "set human set email verification update",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				rowCount, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetEmailVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
					}),
					humanRepo.SetUnverifiedEmail("new@email.com"),
				)
				assert.EqualValues(t, 1, rowCount)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetEmailVerification(&domain.VerificationTypeUpdate{
						Expiry:    gu.Ptr(10 * time.Minute),
						UpdatedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Email.UnverifiedAddress = "new@email.com"
					u.Human.Email.PendingVerification = &domain.Verification{
						ID:        "id",
						CreatedAt: now,
						UpdatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
						ExpiresAt: gu.Ptr(now.Add(10 * time.Minute)),
					}
					return u
				}(),
			},
		},
		{
			name: "set human set email verification skipped",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetEmail("new@email.com"),
					humanRepo.SetEmailVerification(&domain.VerificationTypeSkipped{
						SkippedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Email.PendingVerification = nil
					u.Human.Email.Address = "new@email.com"
					u.Human.Email.VerifiedAt = now
					return u
				}(),
			},
		},
		{
			name: "set human set email verification failed",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				rowCount, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetEmailVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
					}),
					humanRepo.SetUnverifiedEmail("new@email.com"),
				)
				assert.EqualValues(t, 1, rowCount)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetEmailVerification(&domain.VerificationTypeFailed{
						FailedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Email.UnverifiedAddress = "new@email.com"
					u.Human.Email.PendingVerification = &domain.Verification{
						ID:        "id",
						CreatedAt: now,
						UpdatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
						FailedAttempts: 1,
					}
					return u
				}(),
			},
		},
		{
			name: "set human enable email otp at",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.EnableEmailOTPAt(now),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					assert.False(t, u.Human.Email.OTP.EnabledAt.Equal(now))
					u.Human.Email.OTP.EnabledAt = now
					return u
				}(),
			},
		},
		{
			name: "set human disable email otp",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.DisableEmailOTP(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Email.OTP.EnabledAt = time.Time{}
					return u
				}(),
			},
		},
		{
			name: "set human last successful email otp check",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetLastSuccessfulEmailOTPCheck(now),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					if u.Human.Email.OTP.LastSuccessfullyCheckedAt != nil {
						assert.False(t, u.Human.Email.OTP.LastSuccessfullyCheckedAt.Equal(now))
					}
					u.Human.Email.OTP.LastSuccessfullyCheckedAt = &now
					return u
				}(),
			},
		},
		{
			name: "set human increment email otp failed attempts",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.IncrementEmailOTPFailedAttempts(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Email.OTP.FailedAttempts = 1
					return u
				}(),
			},
		},
		{
			name: "set human reset email otp failed attempts",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				rowCount, err := userRepo.Update(t.Context(), tx, humanCondition, humanRepo.IncrementEmailOTPFailedAttempts())
				assert.EqualValues(t, 1, rowCount)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.ResetEmailOTPFailedAttempts(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Email.OTP.FailedAttempts = 0
					return u
				}(),
			},
		},
		{
			name: "set human set phone verification init",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetPhoneVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
					}),
					humanRepo.SetUnverifiedPhone("+0000000000"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						UnverifiedNumber: "+0000000000",
						PendingVerification: &domain.Verification{
							ID:        "id",
							CreatedAt: now,
							UpdatedAt: now,
							Code: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "sha-1",
								KeyID:      "keyID",
								Crypted:    []byte("cryptedCode"),
							},
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human set phone verified",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetPhoneVerification(&domain.VerificationTypeInit{}),
					humanRepo.SetUnverifiedPhone("+1234567890"),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSucceeded{
						VerifiedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						Number:           "+1234567890",
						UnverifiedNumber: "+1234567890",
						VerifiedAt:       now,
					}
					return u
				}(),
			},
		},
		{
			name: "set human set phone verification update",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				rowCount, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetPhoneVerification(&domain.VerificationTypeInit{
						ID:        gu.Ptr("id"),
						CreatedAt: now,
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "sha-1",
							KeyID:      "keyID",
							Crypted:    []byte("cryptedCode"),
						},
					}),
					humanRepo.SetUnverifiedPhone("+1234567890"),
				)
				assert.EqualValues(t, 1, rowCount)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetPhoneVerification(&domain.VerificationTypeUpdate{
						Expiry:    gu.Ptr(10 * time.Minute),
						UpdatedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						UnverifiedNumber: "+1234567890",
						PendingVerification: &domain.Verification{
							ID:        "id",
							CreatedAt: now,
							UpdatedAt: now,
							Code: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "sha-1",
								KeyID:      "keyID",
								Crypted:    []byte("cryptedCode"),
							},
							ExpiresAt: gu.Ptr(now.Add(10 * time.Minute)),
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human set phone verification skipped",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{
						SkippedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						Number:     "+1234567890",
						VerifiedAt: now,
					}
					return u
				}(),
			},
		},
		{
			name: "set human enable sms otp at",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{
						SkippedAt: now,
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.EnableSMSOTPAt(now),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						Number:     "+1234567890",
						VerifiedAt: now,
						OTP: domain.OTP{
							EnabledAt: now,
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human disable sms otp",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{
						SkippedAt: now,
					}),
					humanRepo.EnableSMSOTPAt(now),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.DisableSMSOTP(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						Number:     "+1234567890",
						VerifiedAt: now,
					}
					return u
				}(),
			},
		},
		{
			name: "set human last successful sms otp check",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{
						SkippedAt: now,
					}),
					humanRepo.EnableSMSOTPAt(now),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetLastSuccessfulSMSOTPCheck(now),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						Number:     "+1234567890",
						VerifiedAt: now,
						OTP: domain.OTP{
							EnabledAt:                 now,
							LastSuccessfullyCheckedAt: &now,
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human increment sms otp failed attempts",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{
						SkippedAt: now,
					}),
					humanRepo.EnableSMSOTPAt(now),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.IncrementSMSOTPFailedAttempts(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						Number:     "+1234567890",
						VerifiedAt: now,
						OTP: domain.OTP{
							EnabledAt:      now,
							FailedAttempts: 1,
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human reset sms otp failed attempts",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{
						SkippedAt: now,
					}),
					humanRepo.EnableSMSOTPAt(now),
				)
				if err != nil {
					return err
				}
				_, err = userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.IncrementSMSOTPFailedAttempts(),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.ResetSMSOTPFailedAttempts(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = &domain.HumanPhone{
						Number:     "+1234567890",
						VerifiedAt: now,
						OTP: domain.OTP{
							EnabledAt:      now,
							FailedAttempts: 0,
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human remove phone",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetPhone("+1234567890"),
					humanRepo.SetPhoneVerification(&domain.VerificationTypeSkipped{
						SkippedAt: now,
					}),
					humanRepo.EnableSMSOTPAt(now),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.RemovePhone(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Phone = nil
					return u
				}(),
			},
		},
		{
			name: "set human totp secret",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetTOTPSecret(&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "aes-256",
						KeyID:      "id",
						Crypted:    []byte("new-secret"),
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.TOTP = &domain.HumanTOTP{
						Secret: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "aes-256",
							KeyID:      "id",
							Crypted:    []byte("new-secret"),
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human totp verified at",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := humanRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetTOTPSecret(&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "aes-256",
						KeyID:      "id",
						Crypted:    []byte("secret"),
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetTOTPVerifiedAt(now),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.TOTP = &domain.HumanTOTP{
						Secret: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "aes-256",
							KeyID:      "id",
							Crypted:    []byte("secret"),
						},
						VerifiedAt: now,
					}
					return u
				}(),
			},
		},
		{
			name: "set human remove totp",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := humanRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetTOTPSecret(&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "aes-256",
						KeyID:      "id",
						Crypted:    []byte("secret"),
					}),
					humanRepo.SetTOTPVerifiedAt(now),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.RemoveTOTP(),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.TOTP = nil
					return u
				}(),
			},
		},
		{
			name: "set human last successful totp check",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := humanRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetTOTPSecret(&crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "aes-256",
						KeyID:      "id",
						Crypted:    []byte("secret"),
					}),
					humanRepo.SetTOTPVerifiedAt(now),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetLastSuccessfulTOTPCheck(now),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.TOTP = &domain.HumanTOTP{
						Secret: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "aes-256",
							KeyID:      "id",
							Crypted:    []byte("secret"),
						},
						VerifiedAt:                now,
						LastSuccessfullyCheckedAt: &now,
					}
					return u
				}(),
			},
		},
		{
			name: "set human add identity provider link",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				return repository.IDProviderRepository().Create(t.Context(), tx, &domain.IdentityProvider{
					InstanceID: instanceID,
					OrgID:      &orgID1,
					ID:         "provider-id",
					Name:       "idp",
					State:      domain.IDPStateActive,
					Payload:    json.RawMessage("{}"),
				})
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       "provider-id",
						ProvidedUserID:   "provided-user-id",
						ProvidedUsername: "provided-username",
						CreatedAt:        now,
						UpdatedAt:        now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.IdentityProviderLinks = []*domain.IdentityProviderLink{
						{
							ProviderID:       "provider-id",
							ProvidedUserID:   "provided-user-id",
							ProvidedUsername: "provided-username",
							CreatedAt:        now,
							UpdatedAt:        now,
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human update identity provider link",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				err := repository.IDProviderRepository().Create(t.Context(), tx, &domain.IdentityProvider{
					InstanceID: instanceID,
					OrgID:      &orgID1,
					ID:         "provider-id",
					Name:       "idp",
					State:      domain.IDPStateActive,
					Payload:    json.RawMessage("{}"),
				})
				if err != nil {
					return err
				}
				_, err = userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       "provider-id",
						ProvidedUserID:   "provided-user-id",
						ProvidedUsername: "provided-username",
						CreatedAt:        now,
						UpdatedAt:        now,
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.UpdateIdentityProviderLink(
						database.And(
							humanRepo.IdentityProviderLinkConditions().ProviderIDCondition("provider-id"),
						),
						humanRepo.SetIdentityProviderLinkUsername("updated-username"),
						humanRepo.SetIdentityProviderLinkProvidedID("new-id"),
					),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.IdentityProviderLinks = []*domain.IdentityProviderLink{
						{
							ProviderID:       "provider-id",
							ProvidedUserID:   "new-id",
							ProvidedUsername: "updated-username",
							CreatedAt:        now,
							UpdatedAt:        now,
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human remove identity provider link",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				err := repository.IDProviderRepository().Create(t.Context(), tx, &domain.IdentityProvider{
					InstanceID: instanceID,
					OrgID:      &orgID1,
					ID:         "provider-id",
					Name:       "idp",
					State:      domain.IDPStateActive,
					Payload:    json.RawMessage("{}"),
				})
				if err != nil {
					return err
				}
				_, err = userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.AddIdentityProviderLink(&domain.IdentityProviderLink{
						ProviderID:       "provider-id",
						ProvidedUserID:   "provided-user-id",
						ProvidedUsername: "provided-username",
						CreatedAt:        now,
						UpdatedAt:        now,
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.RemoveIdentityProviderLink("provider-id", "provided-user-id"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.IdentityProviderLinks = []*domain.IdentityProviderLink{}
					return u
				}(),
			},
		},
		{
			name: "set human add passkey",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.AddPasskey(&domain.Passkey{
						ID:                           "passkey-id",
						PublicKey:                    []byte("public-key"),
						Type:                         domain.PasskeyTypePasswordless,
						CreatedAt:                    now,
						UpdatedAt:                    now,
						VerifiedAt:                   now,
						KeyID:                        []byte("keyID"),
						Name:                         "name",
						AttestationType:              "attestation-type",
						AuthenticatorAttestationGUID: []byte("aa guid"),
						Challenge:                    []byte("challenge"),
						RelyingPartyID:               "rp id",
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Passkeys = append(u.Human.Passkeys, &domain.Passkey{
						ID:                           "passkey-id",
						PublicKey:                    []byte("public-key"),
						Type:                         domain.PasskeyTypePasswordless,
						CreatedAt:                    now,
						UpdatedAt:                    now,
						VerifiedAt:                   now,
						KeyID:                        []byte("keyID"),
						Name:                         "name",
						AttestationType:              "attestation-type",
						AuthenticatorAttestationGUID: []byte("aa guid"),
						Challenge:                    []byte("challenge"),
						RelyingPartyID:               "rp id",
					})
					return u
				}(),
			},
		},
		{
			name: "set human update passkey",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.AddPasskey(&domain.Passkey{
						ID:                           "passkey-id",
						PublicKey:                    []byte("public-key"),
						Type:                         domain.PasskeyTypePasswordless,
						CreatedAt:                    now,
						UpdatedAt:                    now,
						VerifiedAt:                   now,
						KeyID:                        []byte("keyID"),
						Name:                         "name",
						AttestationType:              "attestation-type",
						AuthenticatorAttestationGUID: []byte("aa guid"),
						Challenge:                    []byte("challenge"),
						RelyingPartyID:               "rp id",
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.UpdatePasskey(
						database.And(
							humanRepo.PasskeyConditions().IDCondition("passkey-id"),
						),
						humanRepo.SetPasskeyAttestationType("new-attestation-type"),
						humanRepo.SetPasskeyAuthenticatorAttestationGUID([]byte("new-aaguid")),
						humanRepo.SetPasskeyKeyID([]byte("new-key-id")),
						humanRepo.SetPasskeyName("updated-name"),
						humanRepo.SetPasskeyPublicKey([]byte("new-public-key")),
						humanRepo.SetPasskeySignCount(123),
						humanRepo.SetPasskeyUpdatedAt(now.Add(500*time.Millisecond)),
					),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Passkeys = []*domain.Passkey{
						{
							ID:                           "passkey-id",
							PublicKey:                    []byte("new-public-key"),
							Type:                         domain.PasskeyTypePasswordless,
							CreatedAt:                    now,
							UpdatedAt:                    now.Add(500 * time.Millisecond),
							VerifiedAt:                   now,
							KeyID:                        []byte("new-key-id"),
							Name:                         "updated-name",
							AttestationType:              "new-attestation-type",
							AuthenticatorAttestationGUID: []byte("new-aaguid"),
							Challenge:                    []byte("challenge"),
							RelyingPartyID:               "rp id",
							SignCount:                    123,
						},
					}
					return u
				}(),
			},
		},
		{
			name: "set human remove passkey",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.AddPasskey(&domain.Passkey{
						ID:                           "passkey-id",
						PublicKey:                    []byte("public-key"),
						Type:                         domain.PasskeyTypePasswordless,
						CreatedAt:                    now,
						UpdatedAt:                    now,
						VerifiedAt:                   now,
						KeyID:                        []byte("keyID"),
						Name:                         "name",
						AttestationType:              "attestation-type",
						AuthenticatorAttestationGUID: []byte("aa guid"),
						Challenge:                    []byte("challenge"),
						RelyingPartyID:               "rp id",
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.RemovePasskey(humanRepo.PasskeyConditions().IDCondition("passkey-id")),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Passkeys = nil
					return u
				}(),
			},
		},
		{
			name: "set machine name",
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					machineRepo.SetName("Updated Machine"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(machineCondition))
					require.NoError(t, err)
					u.Machine.Name = "Updated Machine"
					return u
				}(),
			},
		},
		{
			name: "set machine description",
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					machineRepo.SetDescription("Updated description"),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(machineCondition))
					require.NoError(t, err)
					u.Machine.Description = "Updated description"
					return u
				}(),
			},
		},
		{
			name: "set machine set secret",
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					machineRepo.SetSecret(gu.Ptr("new-secret")),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(machineCondition))
					require.NoError(t, err)
					u.Machine.Secret = "new-secret"
					return u
				}(),
			},
		},
		{
			name: "set machine access token type",
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					machineRepo.SetAccessTokenType(domain.AccessTokenTypeJWT),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(machineCondition))
					require.NoError(t, err)
					u.Machine.AccessTokenType = domain.AccessTokenTypeJWT
					return u
				}(),
			},
		},
		func() test {
			key := &domain.MachineKey{
				ID:        gofakeit.UUID(),
				PublicKey: []byte("public-key"),
				Type:      domain.MachineKeyTypeJSON,
				ExpiresAt: now.Add(24 * time.Hour),
				CreatedAt: now,
			}
			return test{
				name: "set machine add key",
				args: args{
					condition: machineCondition,
					changes: []database.Change{
						machineRepo.AddKey(key),
					},
				},
				want: want{
					user: func() *domain.User {
						u, err := userRepo.Get(t.Context(), tx, database.WithCondition(machineCondition))
						require.NoError(t, err)
						u.Machine.Keys = append(u.Machine.Keys, key)
						return u
					}(),
				},
			}
		}(),
		{
			name: "set machine remove key",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				rowCount, err := userRepo.Update(t.Context(), tx, machineCondition,
					machineRepo.AddKey(&domain.MachineKey{
						ID:        "key-to-remove",
						PublicKey: []byte("public-key"),
						Type:      domain.MachineKeyTypeJSON,
						ExpiresAt: now.Add(24 * time.Hour),
						CreatedAt: now,
					}),
				)
				assert.EqualValues(t, 1, rowCount)
				return err
			},
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					machineRepo.RemoveKey("key-to-remove"),
				},
			},
			want: want{
				user: machine,
			},
		},
		func() test {
			pat := &domain.PersonalAccessToken{
				ID:        gofakeit.UUID(),
				Scopes:    []string{"openid", "profile"},
				ExpiresAt: now.Add(24 * time.Hour),
				CreatedAt: now,
			}
			return test{
				name: "set machine add personal access token",
				args: args{
					condition: machineCondition,
					changes: []database.Change{
						machineRepo.AddPersonalAccessToken(pat),
					},
				},
				want: want{
					user: func() *domain.User {
						u, err := userRepo.Get(t.Context(), tx, database.WithCondition(machineCondition))
						require.NoError(t, err)
						u.Machine.PATs = append(u.Machine.PATs, pat)
						return u
					}(),
				},
			}
		}(),
		{
			name: "set machine remove personal access token",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				rowCount, err := machineRepo.Update(t.Context(), tx, machineCondition,
					machineRepo.AddPersonalAccessToken(&domain.PersonalAccessToken{
						ID:        "pat-to-remove",
						Scopes:    []string{"openid", "profile"},
						ExpiresAt: now.Add(24 * time.Hour),
						CreatedAt: now,
					}),
				)
				assert.EqualValues(t, 1, rowCount)
				return err
			},
			args: args{
				condition: machineCondition,
				changes: []database.Change{
					machineRepo.RemovePersonalAccessToken("pat-to-remove"),
				},
			},
			want: want{
				user: machine,
			},
		},
		{
			name: "init invite verification",
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetInviteVerification(&domain.VerificationTypeInit{
						CreatedAt: now,
						Expiry:    gu.Ptr(10 * time.Minute),
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "aes-256",
							KeyID:      "key-id",
							Crypted:    []byte("invite code"),
						},
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Invite = &domain.HumanInvite{
						PendingVerification: &domain.Verification{
							CreatedAt: now,
							UpdatedAt: now,
							Code: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "aes-256",
								KeyID:      "key-id",
								Crypted:    []byte("invite code"),
							},
							ExpiresAt: gu.Ptr(now.Add(10 * time.Minute)),
						},
					}
					return u
				}(),
			},
		},
		{
			name: "invite verification succeeded",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetInviteVerification(&domain.VerificationTypeInit{
						CreatedAt: now,
						Expiry:    gu.Ptr(10 * time.Minute),
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "aes-256",
							KeyID:      "key-id",
							Crypted:    []byte("invite code"),
						},
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetInviteVerification(&domain.VerificationTypeSucceeded{
						VerifiedAt: now,
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Invite = &domain.HumanInvite{
						PendingVerification: nil,
						AcceptedAt:          now,
					}
					return u
				}(),
			},
		},
		{
			name: "invite verification updated",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetInviteVerification(&domain.VerificationTypeInit{
						CreatedAt: now,
						Expiry:    gu.Ptr(10 * time.Minute),
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "aes-256",
							KeyID:      "key-id",
							Crypted:    []byte("invite code"),
						},
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetInviteVerification(&domain.VerificationTypeUpdate{
						Expiry: gu.Ptr(1 * time.Hour),
					}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Invite = &domain.HumanInvite{
						PendingVerification: &domain.Verification{
							CreatedAt: now,
							UpdatedAt: now,
							Code: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "aes-256",
								KeyID:      "key-id",
								Crypted:    []byte("invite code"),
							},
							ExpiresAt: gu.Ptr(now.Add(1 * time.Hour)),
						},
					}
					return u
				}(),
			},
		},
		{
			name: "invite verification failed",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := userRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetInviteVerification(&domain.VerificationTypeInit{
						CreatedAt: now,
						Expiry:    gu.Ptr(10 * time.Minute),
						Code: &crypto.CryptoValue{
							CryptoType: crypto.TypeHash,
							Algorithm:  "aes-256",
							KeyID:      "key-id",
							Crypted:    []byte("invite code"),
						},
					}),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.SetInviteVerification(&domain.VerificationTypeFailed{}),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.Invite = &domain.HumanInvite{
						FailedAttempts: 1,
						PendingVerification: &domain.Verification{
							CreatedAt: now,
							UpdatedAt: now,
							Code: &crypto.CryptoValue{
								CryptoType: crypto.TypeHash,
								Algorithm:  "aes-256",
								KeyID:      "key-id",
								Crypted:    []byte("invite code"),
							},
							ExpiresAt: gu.Ptr(now.Add(1 * time.Hour)),
						},
					}
					return u
				}(),
			},
		},
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

			_, err = userRepo.Update(t.Context(), savepoint, tt.args.condition, tt.args.changes...)
			require.ErrorIs(t, err, tt.want.err)

			user, err := userRepo.Get(t.Context(), savepoint, database.WithCondition(tt.args.condition))
			require.NoError(t, err)
			assertUser(t, tt.want.user, user)
		})
	}
}

func Test_user_Delete(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()

	instanceID := createInstance(t, tx)
	instance2ID := createInstance(t, tx)

	orgID1 := createOrganization(t, tx, instanceID)
	orgID2 := createOrganization(t, tx, instanceID)
	orgID3 := createOrganization(t, tx, instance2ID)

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
				Address: "john@doe.com",
			},
		},
	}
	humanCondition := userRepo.PrimaryKeyCondition(instanceID, human.ID)

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
	machineCondition := userRepo.PrimaryKeyCondition(instanceID, machine.ID)

	instance2Human := &domain.User{
		ID:             human.ID,
		InstanceID:     instance2ID,
		OrganizationID: orgID3,
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
				Address: "john@doe.com",
			},
		},
	}

	err := userRepo.Create(t.Context(), tx, human)
	require.NoError(t, err)
	err = userRepo.Create(t.Context(), tx, instance2Human)
	require.NoError(t, err)
	err = userRepo.Create(t.Context(), tx, machine)
	require.NoError(t, err)

	type args struct {
		condition database.Condition
	}
	type want struct {
		err      error
		rowCount int
	}
	type test struct {
		name  string
		setup func(t *testing.T, tx database.QueryExecutor) error
		args  args
		want  want
	}
	tests := []test{
		{
			name: "delete human user",
			args: args{
				condition: humanCondition,
			},
			want: want{
				err:      nil,
				rowCount: 1,
			},
		},
		{
			name: "delete machine user",
			args: args{
				condition: machineCondition,
			},
			want: want{
				err:      nil,
				rowCount: 1,
			},
		},
		{
			name: "delete non-existing user",
			args: args{
				condition: userRepo.PrimaryKeyCondition(instanceID, "non-existing-id"),
			},
			want: want{
				err:      nil,
				rowCount: 0,
			},
		},
		{
			name: "wrong condition",
			args: args{
				condition: userRepo.TypeCondition(domain.UserTypeHuman),
			},
			want: want{
				err:      new(database.MissingConditionError),
				rowCount: 0,
			},
		},
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

			rowCount, err := userRepo.Delete(t.Context(), savepoint, tt.args.condition)
			require.ErrorIs(t, err, tt.want.err)
			assert.EqualValues(t, tt.want.rowCount, rowCount)
			if tt.want.err != nil {
				return
			}

			user, err := userRepo.Get(t.Context(), savepoint, database.WithCondition(tt.args.condition))
			assert.ErrorIs(t, err, new(database.NoRowFoundError))
			assert.Nil(t, user)
		})
	}
}

func assertUsers(t *testing.T, expected, actual []*domain.User) {
	t.Helper()
	require.Len(t, actual, len(expected))
	for i := range expected {
		assertUser(t, expected[i], actual[i])
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

	assertMetadata(t, expected.Metadata, actual.Metadata)

	if expected.Machine != nil {
		assertMachineUser(t, expected.Machine, actual.Machine)
		return
	}

	if expected.Human != nil {
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

	assertPersonalAccessTokens(t, expected.PATs, actual.PATs)
	assertMachineKeys(t, expected.Keys, actual.Keys)
}

func assertHumanUser(t *testing.T, expected, actual *domain.HumanUser) {
	t.Helper()

	assert.Equal(t, expected.FirstName, actual.FirstName, "human first name")
	assert.Equal(t, expected.LastName, actual.LastName, "human last name")
	assert.Equal(t, expected.Nickname, actual.Nickname, "human nickname")
	assert.Equal(t, expected.DisplayName, actual.DisplayName, "human display name")
	assert.Equal(t, expected.PreferredLanguage, actual.PreferredLanguage, "human preferred language")
	assert.Equal(t, expected.Gender, actual.Gender, "human gender")
	assert.Equal(t, expected.AvatarKey, actual.AvatarKey, "human avatar key")
	assert.True(t, expected.MultifactorInitializationSkippedAt.Equal(actual.MultifactorInitializationSkippedAt), "human multifactor initialization skipped at")

	assertHumanEmail(t, expected.Email, actual.Email)
	assertHumanPhone(t, expected.Phone, actual.Phone)
	assertPasskeys(t, expected.Passkeys, actual.Passkeys)
	assertHumanPassword(t, expected.Password, actual.Password)
	assertHumanTOTP(t, expected.TOTP, actual.TOTP)

	assertIdentityProviderLinks(t, expected.IdentityProviderLinks, actual.IdentityProviderLinks)
	assertVerifications(t, expected.Verifications, actual.Verifications)
}

func assertHumanTOTP(t *testing.T, expected, actual *domain.HumanTOTP) {
	t.Helper()

	require.Equal(t, expected == nil, actual == nil, "human totp nil check")
	if expected == nil {
		return
	}

	assert.True(t, expected.VerifiedAt.Equal(actual.VerifiedAt), "human totp verified at")
	assert.Equal(t, expected.Secret, actual.Secret)
	if expected.LastSuccessfullyCheckedAt != nil {
		require.NotNil(t, actual.LastSuccessfullyCheckedAt, "human totp last successfully checked at not nil")
		assert.True(t, expected.LastSuccessfullyCheckedAt.Equal(*actual.LastSuccessfullyCheckedAt), "human totp last successfully checked at")
	} else {
		assert.Nil(t, actual.LastSuccessfullyCheckedAt, "human totp last successfully checked at nil")
	}
	assert.Equal(t, expected.FailedAttempts, actual.FailedAttempts, "human totp failed attempts")
}

func assertHumanPassword(t *testing.T, expected, actual domain.HumanPassword) {
	t.Helper()

	assert.Equal(t, expected.Hash, actual.Hash, "human password hash")
	assert.Equal(t, expected.IsChangeRequired, actual.IsChangeRequired, "human password is change required")
	assert.True(t, expected.ChangedAt.Equal(actual.ChangedAt), "human password changed at")
	assertVerification(t, expected.PendingVerification, actual.PendingVerification, "human password pending verification")
	if expected.LastSuccessfullyCheckedAt != nil {
		require.NotNil(t, actual.LastSuccessfullyCheckedAt, "human password last successfully checked at not nil")
		assert.True(t, expected.LastSuccessfullyCheckedAt.Equal(*actual.LastSuccessfullyCheckedAt), "human password last successfully checked at")
	} else {
		assert.Nil(t, actual.LastSuccessfullyCheckedAt, "human password last successfully checked at nil")
	}
	assert.Equal(t, expected.FailedAttempts, actual.FailedAttempts, "human password failed attempts")
}

func assertHumanPhone(t *testing.T, expected, actual *domain.HumanPhone) {
	t.Helper()

	require.Equal(t, expected == nil, actual == nil, "human phone number nil check")
	if expected == nil {
		return
	}

	assert.Equal(t, expected.Number, actual.Number, "human phone number")
	assert.Equal(t, expected.UnverifiedNumber, actual.UnverifiedNumber, "human unverified phone number")
	assert.True(t, expected.VerifiedAt.Equal(actual.VerifiedAt), "human phone verified at")
	assertOTP(t, expected.OTP, actual.OTP, "human phone otp")
	assertVerification(t, expected.PendingVerification, actual.PendingVerification, "human phone pending verification")
}

func assertHumanEmail(t *testing.T, expected, actual domain.HumanEmail) {
	t.Helper()

	assert.Equal(t, expected.Address, actual.Address, "human email address")
	assert.Equal(t, expected.UnverifiedAddress, actual.UnverifiedAddress, "human unverified email address")
	assert.True(t, expected.VerifiedAt.Equal(actual.VerifiedAt), "human email verified at")
	assertOTP(t, expected.OTP, actual.OTP, "human email otp")
	assertVerification(t, expected.PendingVerification, actual.PendingVerification, "human email pending verification")
}

func assertOTP(t *testing.T, expected, actual domain.OTP, field string) {
	t.Helper()

	assert.True(t, expected.EnabledAt.Equal(actual.EnabledAt), field+" enabled at")
	if expected.LastSuccessfullyCheckedAt != nil {
		require.NotNil(t, actual.LastSuccessfullyCheckedAt, field+" last successfully checked at not nil")
		assert.True(t, expected.LastSuccessfullyCheckedAt.Equal(*actual.LastSuccessfullyCheckedAt), field+" last successfully checked at")
	} else {
		assert.Nil(t, actual.LastSuccessfullyCheckedAt, field+" last successfully checked at nil")
	}
	assert.Equal(t, expected.FailedAttempts, actual.FailedAttempts, field+" failed attempts")
}

func assertCryptoValue(t *testing.T, expected, actual *crypto.CryptoValue, field string) {
	t.Helper()

	if expected == nil {
		assert.Nil(t, actual, field+" nil")
		return
	}

	require.NotNil(t, actual, field+" not nil")
	assert.Equal(t, expected.Crypted, actual.Crypted, field+" crypted")
	assert.Equal(t, expected.CryptoType, actual.CryptoType, field+" crypto type")
	assert.Equal(t, expected.Algorithm, actual.Algorithm, field+" algorithm")
	assert.Equal(t, expected.KeyID, actual.KeyID, field+" key id")
}

func assertPersonalAccessTokens(t *testing.T, expected, actual []*domain.PersonalAccessToken) {
	assert.Len(t, actual, len(expected), "machine pats")
	for _, expectedPAT := range expected {
		var actualPAT *domain.PersonalAccessToken
		actual = slices.DeleteFunc(actual, func(p *domain.PersonalAccessToken) bool {
			if p.ID == expectedPAT.ID {
				actualPAT = p
				return true
			}
			return false
		})
		require.NotNil(t, actualPAT, "pat[%s] not found", expectedPAT.ID)
		assert.Equal(t, expectedPAT.ID, actualPAT.ID, "pat[%s] id", expectedPAT.ID)
		assert.True(t, expectedPAT.CreatedAt.Equal(actualPAT.CreatedAt), "pat[%s] created at", expectedPAT.ID)
		assert.True(t, expectedPAT.ExpiresAt.Equal(actualPAT.ExpiresAt), "pat[%s] expires at", expectedPAT.ID)
		assert.ElementsMatch(t, expectedPAT.Scopes, actualPAT.Scopes, "pat[%s] scopes", expectedPAT.ID)
	}
	assert.Empty(t, actual, "unmatched machine pats")
}

func assertMachineKeys(t *testing.T, expected, actual []*domain.MachineKey) {
	t.Helper()

	require.Len(t, actual, len(expected), "machine keys")
	for _, expectedKey := range expected {
		var actualKey *domain.MachineKey
		actual = slices.DeleteFunc(actual, func(k *domain.MachineKey) bool {
			if k.ID == expectedKey.ID {
				actualKey = k
				return true
			}
			return false
		})
		require.NotNil(t, actualKey, "machine key[%s] not found", expectedKey.ID)
		assert.Equal(t, expectedKey.PublicKey, actualKey.PublicKey, "machine key[%s] public key", expectedKey.ID)
		assert.True(t, expectedKey.CreatedAt.Equal(actualKey.CreatedAt), "machine key[%s] created at", expectedKey.ID)
		assert.True(t, expectedKey.ExpiresAt.Equal(actualKey.ExpiresAt), "machine key[%s] expires at", expectedKey.ID)
		assert.Equal(t, expectedKey.Type, actualKey.Type, "machine key[%s] type", expectedKey.ID)
	}
	assert.Empty(t, actual, "unmatched machine keys")
}

func assertPasskeys(t *testing.T, expected, actual []*domain.Passkey) {
	t.Helper()

	for i, expectedPasskey := range expected {
		var actualPasskey *domain.Passkey
		actual = slices.DeleteFunc(actual, func(p *domain.Passkey) bool {
			if p.ID == expectedPasskey.ID {
				actualPasskey = p
				return true
			}
			return false
		})
		require.NotNil(t, actualPasskey, "human passkey[%d] not found", i)
		assert.Equal(t, expectedPasskey.KeyID, actualPasskey.KeyID, "human passkey[%d] key id", i)
		assert.Equal(t, expectedPasskey.Name, actualPasskey.Name, "human passkey[%d] name", i)
		assert.Equal(t, expectedPasskey.SignCount, actualPasskey.SignCount, "human passkey[%d] sign count", i)
		assert.Equal(t, expectedPasskey.PublicKey, actualPasskey.PublicKey, "human passkey[%d] public key", i)
		assert.Equal(t, expectedPasskey.AttestationType, actualPasskey.AttestationType, "human passkey[%d] attestation type", i)
		assert.Equal(t, expectedPasskey.AuthenticatorAttestationGUID, actualPasskey.AuthenticatorAttestationGUID, "human passkey[%d] authenticator attestation guid", i)
		assert.Equal(t, expectedPasskey.Type, actualPasskey.Type, "human passkey[%d] type", i)
		assert.True(t, expectedPasskey.CreatedAt.Equal(actualPasskey.CreatedAt), "human passkey[%d] created at", i)
		assert.True(t, expectedPasskey.UpdatedAt.Equal(actualPasskey.UpdatedAt), "human passkey[%d] updated at", i)
		assert.True(t, expectedPasskey.VerifiedAt.Equal(actualPasskey.VerifiedAt), "human passkey[%d] verified at", i)
		assert.Equal(t, expectedPasskey.Challenge, actualPasskey.Challenge, "human passkey[%d] challenge", i)
		assert.Equal(t, expectedPasskey.RelyingPartyID, actualPasskey.RelyingPartyID, "human passkey[%d] relying party id", i)
		assertVerification(t, expectedPasskey.Initialization, actualPasskey.Initialization, fmt.Sprintf("human passkey[%d] initialization", i))
	}
	assert.Empty(t, actual, "unmatched human passkeys")
}

func assertVerifications(t *testing.T, expected, actual []*domain.Verification) {
	t.Helper()

	for i, expectedVerification := range expected {
		var actualVerification *domain.Verification
		actual = slices.DeleteFunc(actual, func(v *domain.Verification) bool {
			if v.ID == expectedVerification.ID {
				actualVerification = v
				return true
			}
			return false
		})
		require.NotNil(t, actualVerification, "human verification[%d] not found", i)
		assert.Equal(t, expectedVerification.ID, actualVerification.ID, "human verification[%d] id", i)
		assert.Equal(t, expectedVerification.Code, actualVerification.Code, "human verification[%d] code", i)
		if expectedVerification.ExpiresAt != nil {
			require.NotNil(t, actualVerification.ExpiresAt, "human verification[%d] expires at not nil", i)
			assert.True(t, expectedVerification.ExpiresAt.Equal(*actualVerification.ExpiresAt), "human verification[%d] expires at", i)
		} else {
			assert.Nil(t, actualVerification.ExpiresAt, "human verification[%d] expires at nil", i)
		}
		assert.Equal(t, expectedVerification.FailedAttempts, actualVerification.FailedAttempts, "human verification[%d] failed attempts", i)
	}
	assert.Empty(t, actual, "unmatched human verifications")
}

func assertVerification(t *testing.T, expected, actual *domain.Verification, field string) {
	t.Helper()

	require.Equal(t, expected == nil, actual == nil, "%s nil check", field)
	if expected == nil {
		return
	}

	assert.Equal(t, expected.ID, actual.ID, "%s verification id", field+" id")
	assertCryptoValue(t, expected.Code, actual.Code, field+" code")
	assert.True(t, expected.CreatedAt.Equal(actual.CreatedAt), field+" created at")
	assert.True(t, expected.UpdatedAt.Equal(actual.UpdatedAt), field+" updated at")

	if expected.ExpiresAt != nil {
		require.NotNil(t, actual.ExpiresAt, "%s expires at not nil", field)
		assert.True(t, expected.ExpiresAt.Equal(*actual.ExpiresAt), "%s expires at", field)
	} else {
		assert.Nil(t, actual.ExpiresAt, "%s expires at nil", field)
	}
	assert.Equal(t, expected.FailedAttempts, actual.FailedAttempts, "%s failed attempts", field)
}

func assertIdentityProviderLinks(t *testing.T, expected, actual []*domain.IdentityProviderLink) {
	t.Helper()

	require.Len(t, actual, len(expected), "identity provider links")
	for i, expectedLink := range expected {
		assert.Equal(t, expectedLink.ProviderID, actual[i].ProviderID)
		assert.Equal(t, expectedLink.ProvidedUserID, actual[i].ProvidedUserID)
		assert.Equal(t, expectedLink.ProvidedUsername, actual[i].ProvidedUsername)
		assert.True(t, expectedLink.CreatedAt.Equal(actual[i].CreatedAt), "identity provider link[%d] created at", i)
		assert.True(t, expectedLink.UpdatedAt.Equal(actual[i].UpdatedAt), "identity provider link[%d] updated at", i)
	}
}
