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
					userRepo.AddMetadata(
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
					userRepo.AddMetadata(&domain.Metadata{
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
					humanRepo.IncrementPhoneOTPFailedAttempts(),
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
					humanRepo.IncrementPhoneOTPFailedAttempts(),
				)
				return err
			},
			args: args{
				condition: humanCondition,
				changes: []database.Change{
					humanRepo.ResetPhoneOTPFailedAttempts(),
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
					humanRepo.SetTOTPSecret([]byte("new-secret")),
				},
			},
			want: want{
				user: func() *domain.User {
					u, err := userRepo.Get(t.Context(), tx, database.WithCondition(humanCondition))
					require.NoError(t, err)
					u.Human.TOTP = &domain.HumanTOTP{
						Secret: []byte("new-secret"),
					}
					return u
				}(),
			},
		},
		{
			name: "set human totp verified at",
			setup: func(t *testing.T, tx database.QueryExecutor) error {
				_, err := humanRepo.Update(t.Context(), tx, humanCondition,
					humanRepo.SetTOTPSecret([]byte("secret")),
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
						Secret:     []byte("secret"),
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
					humanRepo.SetTOTPSecret([]byte("secret")),
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
					humanRepo.SetTOTPSecret([]byte("secret")),
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
						Secret:                    []byte("secret"),
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
							ProvidedUserID:   "provided-user-id",
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
						humanRepo.SetPasskeyName("updated-name"),
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
							PublicKey:                    []byte("public-key"),
							Type:                         domain.PasskeyTypePasswordless,
							CreatedAt:                    now,
							UpdatedAt:                    now,
							VerifiedAt:                   now,
							KeyID:                        []byte("keyID"),
							Name:                         "updated-name",
							AttestationType:              "attestation-type",
							AuthenticatorAttestationGUID: []byte("aa guid"),
							Challenge:                    []byte("challenge"),
							RelyingPartyID:               "rp id",
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
					u.Machine.Secret = []byte("new-secret")
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
				ExpiresAt: time.Now().Add(24 * time.Hour),
				CreatedAt: time.Now(),
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
						ExpiresAt: time.Now().Add(24 * time.Hour),
						CreatedAt: time.Now(),
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
				ExpiresAt: time.Now().Add(24 * time.Hour),
				CreatedAt: time.Now(),
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
				rowCount, err := userRepo.Update(t.Context(), tx, machineCondition,
					machineRepo.AddPersonalAccessToken(&domain.PersonalAccessToken{
						ID:        "pat-to-remove",
						Scopes:    []string{"openid", "profile"},
						ExpiresAt: time.Now().Add(24 * time.Hour),
						CreatedAt: time.Now(),
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
