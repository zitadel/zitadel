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

func Test_humanUser_create(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)

	userRepo := repository.UserRepository()
	userRepo = userRepo.LoadKeys().LoadPATs()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	createdAt := time.Now().Round(time.Second)

	// existingUserID := createMachineUser(t, tx, instanceID, orgID)

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
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			email := gofakeit.Email()

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
						Human: &domain.HumanUser{
							FirstName: firstName,
							LastName:  lastName,
							Password: domain.HumanPassword{
								Password:         "my-password",
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								Address: email,
							},
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
						Human: &domain.HumanUser{
							FirstName: firstName,
							LastName:  lastName,
							Email: domain.HumanEmail{
								Address: email,
							},
						},
					},
				},
			}
		}(),
		// func() test {
		// 	username := gofakeit.Username()
		// 	return test{
		// 		name: "already exists",
		// 		args: args{
		// 			user: &domain.User{
		// 				InstanceID:     instanceID,
		// 				OrganizationID: orgID,
		// 				ID:             existingUserID,
		// 				Username:       username,
		// 				State:          domain.UserStateActive,
		// 				CreatedAt:      createdAt,
		// 				Machine: &domain.MachineUser{
		// 					Name: "machine",
		// 				},
		// 			},
		// 		},
		// 		want: want{
		// 			err: new(database.UniqueError),
		// 		},
		// 	}
		// }(),
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
