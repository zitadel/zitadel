package repository_test

import (
	"context"
	"encoding/base64"
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

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	createdAt := time.Now().Round(time.Second)

	// machine := &domain.User{
	// 	ID:             gofakeit.UUID(),
	// 	InstanceID:     instanceID,
	// 	OrganizationID: orgID,
	// 	Username:       "machine-user",
	// 	State:          domain.UserStateActive,
	// 	Machine: &domain.MachineUser{
	// 		Name:        "My Machine",
	// 		Description: "This is my machine user",
	// 	},
	// }

	// err := userRepo.Create(t.Context(), tx, human)
	// require.NoError(t, err)
	// err = userRepo.Create(t.Context(), tx, machine)
	// require.NoError(t, err)

	type args struct {
		user *domain.User
	}
	type want struct {
		err  error
		user *domain.User
	}
	type test struct {
		name  string
		setup func(t *testing.T, tx database.Transaction) error
		args  args
		want  want
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
			secret := []byte(base64.StdEncoding.EncodeToString([]byte("secret")))
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
		// {
		// 	name: "with keys",
		// },
		// {
		// 	name: "with personal access tokens",
		// },
		// {
		// 	name: "all fields set",
		// },
		// {
		// 	name: "already exists",
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

			err = userRepo.Create(t.Context(), savepoint, tt.args.user)
			require.ErrorIs(t, err, tt.want.err)
			got, err := userRepo.Get(t.Context(), savepoint, database.WithCondition(
				userRepo.PrimaryKeyCondition(tt.want.user.InstanceID, tt.want.user.ID),
			))
			if tt.want.user == nil {
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.user, got)
		})
	}
}
