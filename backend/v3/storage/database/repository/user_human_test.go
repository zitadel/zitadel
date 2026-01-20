package repository_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func Test_humanUser_create(t *testing.T) {
	// tx, rollback := transactionForRollback(t)
	// t.Cleanup(rollback)
	tx := pool

	userRepo := repository.UserRepository()
	userRepo = userRepo.LoadKeys().LoadPATs()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	createdAt := time.Now().Round(time.Second)

	// existingUserID := createMachineUser(t, tx, instanceID, orgID)
	password := []byte("my-password")

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
								Password:         password,
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
							Password: domain.HumanPassword{
								Password:         password,
								IsChangeRequired: true,
							},
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			email := gofakeit.Email()

			return test{
				name: "including optional profile fields",
				args: args{
					user: &domain.User{
						InstanceID:     instanceID,
						OrganizationID: orgID,
						ID:             id,
						Username:       username,
						State:          domain.UserStateActive,
						CreatedAt:      createdAt,
						Human: &domain.HumanUser{
							FirstName:                          firstName,
							LastName:                           lastName,
							Nickname:                           "Nick",
							DisplayName:                        firstName + " " + lastName,
							PreferredLanguage:                  language.Georgian,
							Gender:                             domain.HumanGenderFemale,
							AvatarKey:                          "http://localhost:8080/my/profile/picture",
							MultifactorInitializationSkippedAt: createdAt,
							Password: domain.HumanPassword{
								Password:         password,
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
							FirstName:                          firstName,
							LastName:                           lastName,
							Nickname:                           "Nick",
							DisplayName:                        firstName + " " + lastName,
							PreferredLanguage:                  language.Georgian,
							Gender:                             domain.HumanGenderFemale,
							AvatarKey:                          "http://localhost:8080/my/profile/picture",
							MultifactorInitializationSkippedAt: createdAt,
							Email: domain.HumanEmail{
								Address: email,
							},
							Password: domain.HumanPassword{
								Password:         password,
								IsChangeRequired: true,
							},
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			email := gofakeit.Email()

			return test{
				name: "with unverified email without expiry",
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
								Password:         password,
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								Address: email,
								Unverified: &domain.Verification{
									Code: []byte("verification-code"),
								},
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
								Unverified: &domain.Verification{
									Value: &email,
									Code:  []byte("verification-code"),
								},
							},
							Password: domain.HumanPassword{
								Password:         password,
								IsChangeRequired: true,
							},
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			email := gofakeit.Email()
			phone := gofakeit.Phone()

			return test{
				name: "with unverified phone without expiry",
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
								Password:         password,
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								Address: email,
							},
							Phone: &domain.HumanPhone{
								Number: phone,
								Unverified: &domain.Verification{
									Code: []byte("verification-code"),
								},
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
							Password: domain.HumanPassword{
								Password:         password,
								IsChangeRequired: true,
							},
							Phone: &domain.HumanPhone{
								Unverified: &domain.Verification{
									Value: &phone,
									Code:  []byte("verification-code"),
								},
							},
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			email := gofakeit.Email()
			passkeys := []*domain.Passkey{
				{
					ID:                           gofakeit.UUID(),
					KeyID:                        []byte(gofakeit.UUID()),
					Type:                         domain.PasskeyTypePasswordless,
					Name:                         "passwordless",
					PublicKey:                    []byte("public key"),
					AttestationType:              "don't know",
					AuthenticatorAttestationGUID: []byte("aaguid"),
					CreatedAt:                    createdAt,
					UpdatedAt:                    createdAt,
					Challenge:                    []byte("challenge"),
					VerifiedAt:                   createdAt,
					RelyingPartyID:               "rpid",
				},
				{
					ID:                           gofakeit.UUID(),
					KeyID:                        []byte(gofakeit.UUID()),
					Type:                         domain.PasskeyTypeU2F,
					Name:                         "u2f",
					PublicKey:                    []byte("public key"),
					AttestationType:              "don't know",
					AuthenticatorAttestationGUID: []byte("aaguid"),
					CreatedAt:                    createdAt,
					UpdatedAt:                    createdAt,
					Challenge:                    []byte("challenge"),
					VerifiedAt:                   createdAt,
					RelyingPartyID:               "rpid",
				},
			}

			return test{
				name: "with verified passkeys",
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
								Password:         password,
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								Address: email,
							},
							Passkeys: passkeys,
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
							Password: domain.HumanPassword{
								Password:         password,
								IsChangeRequired: true,
							},
							Passkeys: passkeys,
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			email := gofakeit.Email()
			verifications := []*domain.Verification{
				{
					ID:   gofakeit.UUID(),
					Code: []byte("code"),
				},
				{
					ID:   gofakeit.UUID(),
					Code: []byte("code"),
				},
			}

			return test{
				name: "with verifications",
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
								Password:         password,
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								Address: email,
							},
							Verifications: verifications,
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
							Password: domain.HumanPassword{
								Password:         password,
								IsChangeRequired: true,
							},
							Verifications: verifications,
						},
					},
				},
			}
		}(),
		func() test {
			username := gofakeit.Username()
			id := gofakeit.UUID()
			firstName := gofakeit.FirstName()
			lastName := gofakeit.LastName()
			email := gofakeit.Email()

			return test{
				name: "with totp",
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
								Password:         password,
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								Address: email,
							},
							TOTP: domain.HumanTOTP{
								Unverified: &domain.Verification{
									ID:   gofakeit.UUID(),
									Code: []byte("code"),
								},
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
							Password: domain.HumanPassword{
								Password:         password,
								IsChangeRequired: true,
							},
							TOTP: domain.HumanTOTP{
								Unverified: &domain.Verification{
									ID:   gofakeit.UUID(),
									Code: []byte("code"),
								},
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
