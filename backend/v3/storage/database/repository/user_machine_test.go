package repository_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func Test_machineUser_create(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()
	userRepo = userRepo.LoadKeys().LoadPATs()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	createdAt := time.Now().Round(time.Second)

	existingUserID := createMachineUser(t, tx, instanceID, orgID)

	type args struct {
		user *domain.User
	}
	type want struct {
		err  error
		user *domain.User
	}
	type test struct {
		name string
		args args
		want want
	}
	tests := []test{
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			return test{
				name: "minimal representation",
				args: args{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:        "machine",
							Description: "description",
						},
					},
				},
				want: want{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						UpdatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:        "machine",
							Description: "description",
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			secret := []byte("secret")
			return test{
				name: "with secret",
				args: args{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:        "machine",
							Description: "description",
							Secret:      secret,
						},
					},
				},
				want: want{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						UpdatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:        "machine",
							Description: "description",
							Secret:      secret,
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			return test{
				name: "with access token type",
				args: args{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:            "machine",
							Description:     "description",
							AccessTokenType: domain.PersonalAccessTokenTypeBearer,
						},
					},
				},
				want: want{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						UpdatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:            "machine",
							Description:     "description",
							AccessTokenType: domain.PersonalAccessTokenTypeBearer,
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			keys := []*domain.MachineKey{
				{
					ID:        gofakeit.UUID(),
					PublicKey: []byte("public key"),
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
					Type:      domain.MachineKeyTypeJSON,
				},
				{
					ID:        gofakeit.UUID(),
					PublicKey: []byte("public key"),
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
					Type:      domain.MachineKeyTypeJSON,
				},
			}
			return test{
				name: "with keys",
				args: args{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:        "machine",
							Description: "description",
							Keys:        keys,
						},
					},
				},
				want: want{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						UpdatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:        "machine",
							Description: "description",
							Keys:        keys,
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			pats := []*domain.PersonalAccessToken{
				{
					ID:        gofakeit.UUID(),
					PublicKey: []byte("public key"),
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
					Type:      domain.PersonalAccessTokenTypeBearer,
				},
				{
					ID:        gofakeit.UUID(),
					PublicKey: []byte("public key"),
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
					Type:      domain.PersonalAccessTokenTypeJWT,
					Scopes:    []string{"some", "scopes"},
				},
			}
			return test{
				name: "with keys",
				args: args{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:        "machine",
							Description: "description",
							PATs:        pats,
						},
					},
				},
				want: want{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						UpdatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:        "machine",
							Description: "description",
							PATs:        pats,
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			pats := []*domain.PersonalAccessToken{
				{
					ID:        gofakeit.UUID(),
					PublicKey: []byte("public key"),
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
					Type:      domain.PersonalAccessTokenTypeBearer,
				},
			}
			keys := []*domain.MachineKey{
				{
					ID:        gofakeit.UUID(),
					PublicKey: []byte("public key"),
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
					Type:      domain.MachineKeyTypeJSON,
				},
			}
			return test{
				name: "all fields set",
				args: args{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:            "machine",
							Description:     "description",
							PATs:            pats,
							Keys:            keys,
							Secret:          []byte("secret"),
							AccessTokenType: domain.PersonalAccessTokenTypeJWT,
						},
					},
				},
				want: want{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						UpdatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name:            "machine",
							Description:     "description",
							PATs:            pats,
							Keys:            keys,
							Secret:          []byte("secret"),
							AccessTokenType: domain.PersonalAccessTokenTypeJWT,
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			return test{
				name: "already exists",
				args: args{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             existingUserID,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						Machine: &domain.MachineUser{
							Name: "machine",
						},
					},
				},
				want: want{
					err: new(database.UniqueError),
				},
			}
		}(),
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

			err = userRepo.Create(t.Context(), savepoint, tt.args.user)
			require.ErrorIs(t, err, tt.want.err)
			if tt.want.err != nil {
				return
			}
			got, err := userRepo.Get(t.Context(), savepoint, database.WithCondition(
				userRepo.PrimaryKeyCondition(tt.want.user.InstanceID, tt.want.user.ID),
			))
			require.NoError(t, err)
			assertUser(t, tt.want.user, got)
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

	if expected.Machine != nil {
		assert.Equal(t, expected.Machine.Name, actual.Machine.Name, "machine name")
		assert.Equal(t, expected.Machine.Description, actual.Machine.Description, "machine description")
		assert.Equal(t, expected.Machine.AccessTokenType, actual.Machine.AccessTokenType, "machine access token type")
		assert.Equal(t, expected.Machine.Secret, actual.Machine.Secret, "machine secret")
		require.Len(t, actual.Machine.PATs, len(expected.Machine.PATs), "machine pats")
		for i, expectedPAT := range expected.Machine.PATs {
			var actualPAT *domain.PersonalAccessToken
			actual.Machine.PATs = slices.DeleteFunc(actual.Machine.PATs, func(p *domain.PersonalAccessToken) bool {
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
		assert.Empty(t, actual.Machine.PATs, "unmatched machine pats")
		require.Len(t, actual.Machine.Keys, len(expected.Machine.Keys), "machine keys")
		for i, expectedKey := range expected.Machine.Keys {
			var actualKey *domain.MachineKey
			actual.Machine.Keys = slices.DeleteFunc(actual.Machine.Keys, func(k *domain.MachineKey) bool {
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
		assert.Empty(t, actual.Machine.Keys, "unmatched machine keys")
	}
}
