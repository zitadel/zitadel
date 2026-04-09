package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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
							Secret:      "secret",
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
							Secret:      "secret",
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
							AccessTokenType: domain.AccessTokenTypeBearer,
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
							AccessTokenType: domain.AccessTokenTypeBearer,
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
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
				},
				{
					ID:        gofakeit.UUID(),
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
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
					CreatedAt: createdAt,
					ExpiresAt: createdAt.Add(24 * time.Hour),
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
						Metadata: []*domain.Metadata{
							{
								Key:       "meta1",
								Value:     []byte("value1"),
								CreatedAt: createdAt,
								UpdatedAt: createdAt,
							},
							{
								Key:       "meta2",
								Value:     []byte("123"),
								CreatedAt: createdAt,
								UpdatedAt: createdAt,
							},
						},
						Machine: &domain.MachineUser{
							Name:            "machine",
							Description:     "description",
							PATs:            pats,
							Keys:            keys,
							Secret:          "secret",
							AccessTokenType: domain.AccessTokenTypeJWT,
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
						Metadata: []*domain.Metadata{
							{
								InstanceID: instanceID,
								Key:        "meta1",
								Value:      []byte("value1"),
								CreatedAt:  createdAt,
								UpdatedAt:  createdAt,
							},
							{
								InstanceID: instanceID,
								Key:        "meta2",
								Value:      []byte("123"),
								CreatedAt:  createdAt,
								UpdatedAt:  createdAt,
							},
						},
						Machine: &domain.MachineUser{
							Name:            "machine",
							Description:     "description",
							PATs:            pats,
							Keys:            keys,
							Secret:          "secret",
							AccessTokenType: domain.AccessTokenTypeJWT,
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
