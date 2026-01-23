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
	"github.com/zitadel/zitadel/internal/crypto"
)

func Test_userHuman_CheckEmailOTP(t *testing.T) {
	// tx, rollback := transactionForRollback(t)
	// t.Cleanup(rollback)
	tx := pool

	userRepo := repository.UserRepository()
	userRepo = userRepo.LoadKeys().LoadPATs()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	createdAt := time.Now().Round(time.Second)

	code := &crypto.CryptoValue{
		Crypted:    []byte("crypted-code"),
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "aes256",
		KeyID:      "key-id",
	}

	type args struct {
		check domain.CheckType
	}
	type want struct {
		err   error
		email domain.HumanEmail
	}
	type test struct {
		name  string
		setup func(t *testing.T, tx database.QueryExecutor) (userID string)
		args  args
		want  want
	}
	tests := []test{
		{
			name: "init without expiry",
			setup: func(t *testing.T, tx database.QueryExecutor) (userID string) {
				return createHumanUser(t, tx, instanceID, orgID).ID
			},
			args: args{
				check: &domain.CheckTypeInit{
					CreatedAt: createdAt,
					Code:      code,
				},
			},
			want: want{
				email: domain.HumanEmail{
					OTP: domain.OTP{
						Check: &domain.Check{
							Code: code,
						},
					},
				},
			},
		},
		{
			name: "init with expiry",
			setup: func(t *testing.T, tx database.QueryExecutor) (userID string) {
				return createHumanUser(t, tx, instanceID, orgID).ID
			},
			args: args{
				check: &domain.CheckTypeInit{
					CreatedAt: createdAt,
					Code:      code,
					Expiry:    gu.Ptr(24 * time.Hour),
				},
			},
			want: want{
				email: domain.HumanEmail{
					OTP: domain.OTP{
						Check: &domain.Check{
							Code:      code,
							ExpiresAt: gu.Ptr(createdAt.Add(24 * time.Hour)),
						},
					},
				},
			},
		},
		{
			name: "overwrite init",
			setup: func(t *testing.T, tx database.QueryExecutor) (userID string) {
				userID = createHumanUser(t, tx, instanceID, orgID).ID

				_, err := userRepo.Update(t.Context(), tx,
					userRepo.PrimaryKeyCondition(instanceID, userID),
					userRepo.Human().CheckEmailOTP(&domain.CheckTypeInit{
						CreatedAt: createdAt,
						Code: &crypto.CryptoValue{
							Crypted:    []byte("previous-crypted-code"),
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "aes256",
							KeyID:      "previous-key-id",
						},
					}),
				)
				require.NoError(t, err)
				return userID
			},
			args: args{
				check: &domain.CheckTypeInit{
					CreatedAt: createdAt,
					Code:      code,
				},
			},
			want: want{
				email: domain.HumanEmail{
					OTP: domain.OTP{
						Check: &domain.Check{
							Code: code,
						},
					},
				},
			},
		},
		{
			name: "check succeeded",
			setup: func(t *testing.T, tx database.QueryExecutor) (userID string) {
				userID = createHumanUser(t, tx, instanceID, orgID).ID

				_, err := userRepo.Update(t.Context(), tx,
					userRepo.PrimaryKeyCondition(instanceID, userID),
					userRepo.Human().CheckEmailOTP(&domain.CheckTypeInit{
						CreatedAt: createdAt,
						Code: &crypto.CryptoValue{
							Crypted:    []byte("crypted-code"),
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "aes256",
							KeyID:      "key-id",
						},
					}),
				)
				require.NoError(t, err)
				return userID
			},
			args: args{
				check: &domain.CheckTypeSucceeded{},
			},
			want: want{
				email: domain.HumanEmail{
					OTP: domain.OTP{
						LastSuccessfullyCheckedAt: createdAt,
					},
				},
			},
		},
		// {
		// 	name: "check succeeded at time",
		// 	args: args{},
		// 	want: want{},
		// },
		// {
		// 	name: "check failed",
		// 	args: args{},
		// 	want: want{},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// savepoint, err := tx.Begin(t.Context())
			// require.NoError(t, err)
			// t.Cleanup(func() {
			// 	err := savepoint.Rollback(context.Background())
			// 	if err != nil {
			// 		t.Log("rollback savepoint failed", err)
			// 	}
			// })
			savepoint := tx
			var err error

			userID := tt.setup(t, savepoint)
			_, err = userRepo.Update(t.Context(), savepoint,
				userRepo.PrimaryKeyCondition(instanceID, userID),
				userRepo.Human().CheckEmailOTP(tt.args.check),
			)
			require.ErrorIs(t, err, tt.want.err)
			if tt.want.err != nil {
				return
			}

			got, err := userRepo.Get(t.Context(), savepoint, database.WithCondition(
				userRepo.PrimaryKeyCondition(instanceID, userID),
			))
			require.NoError(t, err)
			require.NotNil(t, got.Human)

			assert.True(t, tt.want.email.OTP.LastSuccessfullyCheckedAt.Equal(got.Human.Email.OTP.LastSuccessfullyCheckedAt))
			assertCheck(t, tt.want.email.OTP.Check, got.Human.Email.OTP.Check)
		})
	}
}
