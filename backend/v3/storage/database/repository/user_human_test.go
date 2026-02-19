package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/crypto"
)

func Test_humanUser_create(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	t.Cleanup(rollback)
	// tx := pool

	userRepo := repository.UserRepository()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	createdAt := time.Now().Round(time.Millisecond).UTC()

	existingUserID := createMachineUser(t, tx, instanceID, orgID)
	password := gofakeit.Password(true, true, true, true, false, 16)

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
							Email: domain.HumanEmail{
								UnverifiedAddress: email,
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
								UnverifiedAddress: email,
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
								Hash:             password,
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
								Address:           email,
								UnverifiedAddress: email,
								VerifiedAt:        createdAt,
							},
							Password: domain.HumanPassword{
								Hash:             password,
								IsChangeRequired: true,
								ChangedAt:        createdAt,
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
			code := &crypto.CryptoValue{
				CryptoType: crypto.TypeEncryption,
				Algorithm:  "aes256",
				KeyID:      "key-id",
				Crypted:    []byte("crypted"),
			}

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
								Hash:             password,
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								UnverifiedAddress: email,
								PendingVerification: &domain.Verification{
									ID:   "test",
									Code: code,
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
								UnverifiedAddress: email,
								PendingVerification: &domain.Verification{
									ID:        "test",
									Code:      code,
									CreatedAt: createdAt,
									UpdatedAt: createdAt,
								},
							},
							Password: domain.HumanPassword{
								Hash:             password,
								IsChangeRequired: true,
								ChangedAt:        createdAt,
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
			code := &crypto.CryptoValue{
				CryptoType: crypto.TypeEncryption,
				Algorithm:  "aes256",
				KeyID:      "key-id",
				Crypted:    []byte("crypted"),
			}

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
								Hash:             password,
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								Address:           email,
								UnverifiedAddress: email,
								VerifiedAt:        createdAt,
							},
							Phone: &domain.HumanPhone{
								UnverifiedNumber: phone,
								PendingVerification: &domain.Verification{
									ID:   "id",
									Code: code,
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
								Address:           email,
								UnverifiedAddress: email,
								VerifiedAt:        createdAt,
							},
							Password: domain.HumanPassword{
								Hash:             password,
								IsChangeRequired: true,
								ChangedAt:        createdAt,
							},
							Phone: &domain.HumanPhone{
								UnverifiedNumber: phone,
								PendingVerification: &domain.Verification{
									CreatedAt: createdAt,
									UpdatedAt: createdAt,
									ID:        "id",
									Code:      code,
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
								Hash:             password,
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
								Address:           email,
								UnverifiedAddress: email,
								VerifiedAt:        createdAt,
							},
							Password: domain.HumanPassword{
								Hash:             password,
								IsChangeRequired: true,
								ChangedAt:        createdAt,
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
			code := &crypto.CryptoValue{
				CryptoType: crypto.TypeEncryption,
				Algorithm:  "aes256",
				KeyID:      "key-id",
				Crypted:    []byte("crypted"),
			}
			verifications := []*domain.Verification{
				{
					ID:   gofakeit.UUID(),
					Code: code,
				},
				{
					ID:   gofakeit.UUID(),
					Code: code,
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
								Hash:             password,
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
								Address:           email,
								UnverifiedAddress: email,
								VerifiedAt:        createdAt,
							},
							Password: domain.HumanPassword{
								Hash:             password,
								IsChangeRequired: true,
								ChangedAt:        createdAt,
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
								Hash:             password,
								IsChangeRequired: true,
							},
							Email: domain.HumanEmail{
								Address: email,
							},
							TOTP: &domain.HumanTOTP{
								VerifiedAt: createdAt,
								Secret: &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "aes-256",
									KeyID:      "id",
									Crypted:    []byte("secret"),
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
								Address:           email,
								UnverifiedAddress: email,
								VerifiedAt:        createdAt,
							},
							Password: domain.HumanPassword{
								Hash:             password,
								IsChangeRequired: true,
								ChangedAt:        createdAt,
							},
							TOTP: &domain.HumanTOTP{
								VerifiedAt: createdAt,
								Secret: &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "aes-256",
									KeyID:      "id",
									Crypted:    []byte("secret"),
								},
							},
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
