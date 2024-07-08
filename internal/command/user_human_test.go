package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddHuman(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		idGenerator        id_generator.Generator
		userPasswordHasher *crypto.Hasher
		codeAlg            crypto.EncryptionAlgorithm
		newCode            encrypedCodeFunc
	}
	type args struct {
		ctx             context.Context
		orgID           string
		human           *AddHuman
		secretGenerator crypto.Generator
		allowInitMail   bool
	}
	type res struct {
		want          *domain.ObjectDetails
		wantID        string
		wantEmailCode string
		err           func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "orgid missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal"))
				},
			},
		},
		{
			name: "user invalid, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"))
				},
			},
		},
		{
			name: "with id, already exists, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", AllowedLanguage),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					ID:        "user1",
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: AllowedLanguage,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-k2unb", "Errors.User.AlreadyExisting"))
				},
			},
		},
		{
			name: "domain policy not found, precondition error",

			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: AllowedLanguage,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal"))
				},
			},
		},
		{
			name: "password policy not found, precondition error",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Password:  "pass",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PreferredLanguage: AllowedLanguage,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInternal(nil, "USER-uQ96e", "Errors.Internal"))
				},
			},
		},
		{
			name: "add human with undefined preferred language, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewHumanAddedEvent(context.Background(),
							&userAgg.Aggregate,
							"username",
							"firstname",
							"lastname",
							"",
							"firstname lastname",
							language.Und,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
						user.NewHumanInitialCodeAddedEvent(context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							time.Hour*1,
							"",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				codeAlg:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:     mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human with unsupported preferred language, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewHumanAddedEvent(context.Background(),
							&userAgg.Aggregate,
							"username",
							"firstname",
							"lastname",
							"",
							"firstname lastname",
							UnsupportedLanguage,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
						user.NewHumanInitialCodeAddedEvent(context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							time.Hour*1,
							"",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				codeAlg:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:     mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: UnsupportedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with initial code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewHumanAddedEvent(context.Background(),
							&userAgg.Aggregate,
							"username",
							"firstname",
							"lastname",
							"",
							"firstname lastname",
							AllowedLanguage,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
						user.NewHumanInitialCodeAddedEvent(context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							time.Hour*1,
							"",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				codeAlg:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:     mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with password and initial code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "", AllowedLanguage),
						user.NewHumanInitialCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							1*time.Hour,
							"",
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:            mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with password and email code custom template), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "", AllowedLanguage),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailCode"),
							},
							1*time.Hour,
							"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}",
							false,
							"",
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:            mockEncryptedCode("emailCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:     "email@test.ch",
						URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}",
					},
					PreferredLanguage: AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   false,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with password and return email code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "", AllowedLanguage),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailCode"),
							},
							1*time.Hour,
							"",
							true,
							"",
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:            mockEncryptedCode("emailCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:    "email@test.ch",
						ReturnCode: true,
					},
					PreferredLanguage: AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   false,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID:        "user1",
				wantEmailCode: "emailCode",
			},
		},
		{
			name: "add human email verified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", true, true, "", AllowedLanguage),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PasswordChangeRequired: true,
					PreferredLanguage:      AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human email verified, trim spaces, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", true, true, "", AllowedLanguage),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  " username ",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PasswordChangeRequired: true,
					PreferredLanguage:      AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human, email verified, userLoginMustBeDomain false, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								false,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", true, false, "", AllowedLanguage),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PasswordChangeRequired: true,
					PreferredLanguage:      AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human claimed domain, userLoginMustBeDomain false, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								false,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								"test.ch",
							),
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username@test.ch",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PasswordChangeRequired: true,
					PreferredLanguage:      AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-SFd21", "Errors.User.DomainNotAllowedAsUsername"))
				},
			},
		},
		{
			name: "add human domain, userLoginMustBeDomain false, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								false,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"test.ch",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						func() eventstore.Command {
							event := user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username@test.ch",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								AllowedLanguage,
								domain.GenderUnspecified,
								"email@test.ch",
								false,
							)
							event.AddPasswordData("$plain$x$password", true)
							return event
						}(),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username@test.ch",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PasswordChangeRequired: true,
					PreferredLanguage:      AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},

			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with phone), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567", AllowedLanguage),
						user.NewHumanEmailVerifiedEvent(
							context.Background(),
							&userAgg.Aggregate,
						),
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phonecode"),
							},
							time.Hour*1,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:            mockEncryptedCode("phonecode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Password:  "password",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					Phone: Phone{
						Number: "+41711234567",
					},
					PreferredLanguage: AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with verified phone), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						newAddHumanEvent("", false, true, "+41711234567", AllowedLanguage),
						user.NewHumanInitialCodeAddedEvent(
							context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							1*time.Hour,
							"",
						),
						user.NewHumanPhoneVerifiedEvent(
							context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				codeAlg:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:     mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					Phone: Phone{
						Number:   "+41711234567",
						Verified: true,
					},
					PreferredLanguage: AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		}, {
			name: "add human (with return code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567", AllowedLanguage),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate),
						user.NewHumanPhoneCodeAddedEventV2(
							context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							1*time.Hour,
							true,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:            mockEncryptedCode("phoneCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					Phone: Phone{
						Number:     "+41711234567",
						ReturnCode: true,
					},
					PreferredLanguage: AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human with metadata, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						newAddHumanEvent("", false, true, "", AllowedLanguage),
						user.NewHumanInitialCodeAddedEvent(
							context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							1*time.Hour,
							"",
						),
						user.NewMetadataSetEvent(
							context.Background(),
							&userAgg.Aggregate,
							"testKey",
							[]byte("testValue"),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				codeAlg:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:     mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					Metadata: []*AddMetadataEntry{
						{
							Key:   "testKey",
							Value: []byte("testValue"),
						},
					},
					PreferredLanguage: AllowedLanguage,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				userPasswordHasher: tt.fields.userPasswordHasher,
				userEncryption:     tt.fields.codeAlg,
				newEncryptedCode:   tt.fields.newCode,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			err := r.AddHuman(tt.args.ctx, tt.args.orgID, tt.args.human, tt.args.allowInitMail)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, tt.args.human.Details)
				assert.Equal(t, tt.res.wantID, tt.args.human.ID)
				assert.Equal(t, tt.res.wantEmailCode, gu.Value(tt.args.human.EmailCode))
			}
		})
	}
}

func TestCommandSide_ImportHuman(t *testing.T) {
	type fields struct {
		eventstore         *eventstore.Eventstore
		idGenerator        id_generator.Generator
		userPasswordHasher *crypto.Hasher
	}
	type args struct {
		ctx                  context.Context
		orgID                string
		human                *domain.Human
		passwordless         bool
		links                []*domain.UserIDPLink
		secretGenerator      crypto.Generator
		passwordlessInitCode crypto.Generator
	}
	type res struct {
		wantHuman *domain.Human
		wantCode  *domain.PasswordlessInitCode
		err       func(error) bool
	}
	tests := []struct {
		name  string
		given func(t *testing.T) (fields, args)
		res   res
	}{
		{
			name: "orgid missing, invalid argument error",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
						),
					},
					args{
						ctx:   context.Background(),
						orgID: "",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Email: &domain.Email{
								EmailAddress: "email@test.ch",
							},
						},
					}
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "org policy not found, precondition error",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(),
							expectFilter(),
						),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Email: &domain.Email{
								EmailAddress: "email@test.ch",
							},
						},
					}
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "password policy not found, precondition error",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(),
							expectFilter(),
						),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Email: &domain.Email{
								EmailAddress: "email@test.ch",
							},
						},
					}
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "user invalid, invalid argument error",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
						),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								PreferredLanguage: AllowedLanguage,
							},
						},
					}
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		}, {
			name: "add human (with password and initial code), ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectPush(
								newAddHumanEvent("$plain$x$password", true, true, "", AllowedLanguage),
								user.NewHumanInitialCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour*1,
									"",
								),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Password: &domain.Password{
								SecretString:   "password",
								ChangeRequired: true,
							},
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Email: &domain.Email{
								EmailAddress: "email@test.ch",
							},
						},
						secretGenerator: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: AllowedLanguage,
					},
					Email: &domain.Email{
						EmailAddress: "email@test.ch",
					},
					State: domain.UserStateInitial,
				},
			},
		},
		{
			name: "add human email verified password change not required, ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectPush(
								newAddHumanEvent("$plain$x$password", false, true, "", AllowedLanguage),
								user.NewHumanEmailVerifiedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
								),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Password: &domain.Password{
								SecretString:   "password",
								ChangeRequired: false,
							},
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Email: &domain.Email{
								EmailAddress:    "email@test.ch",
								IsEmailVerified: true,
							},
						},
						secretGenerator: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: AllowedLanguage,
					},
					Email: &domain.Email{
						EmailAddress:    "email@test.ch",
						IsEmailVerified: true,
					},
					State: domain.UserStateActive,
				},
			},
		},
		{
			name: "add human email verified passwordless only, ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectFilter(),
							expectPush(
								newAddHumanEvent("", false, true, "", AllowedLanguage),
								user.NewHumanEmailVerifiedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate),
								user.NewHumanPasswordlessInitCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"code1",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour,
								),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1", "code1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Email: &domain.Email{
								EmailAddress:    "email@test.ch",
								IsEmailVerified: true,
							},
						},
						passwordless:         true,
						secretGenerator:      GetMockSecretGenerator(t),
						passwordlessInitCode: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: AllowedLanguage,
					},
					Email: &domain.Email{
						EmailAddress:    "email@test.ch",
						IsEmailVerified: true,
					},
					State: domain.UserStateActive,
				},
				wantCode: &domain.PasswordlessInitCode{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Expiration: time.Hour,
					CodeID:     "code1",
					Code:       "a",
					State:      domain.PasswordlessInitCodeStateActive,
				},
			},
		},
		{
			name: "add human email verified passwordless and password change not required, ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectFilter(),
							expectPush(
								newAddHumanEvent("$plain$x$password", false, true, "", AllowedLanguage),
								user.NewHumanEmailVerifiedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate),
								user.NewHumanPasswordlessInitCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"code1",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour,
								),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1", "code1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Password: &domain.Password{
								SecretString:   "password",
								ChangeRequired: false,
							},
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Email: &domain.Email{
								EmailAddress:    "email@test.ch",
								IsEmailVerified: true,
							},
						},
						passwordless:         true,
						secretGenerator:      GetMockSecretGenerator(t),
						passwordlessInitCode: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: AllowedLanguage,
					},
					Email: &domain.Email{
						EmailAddress:    "email@test.ch",
						IsEmailVerified: true,
					},
					State: domain.UserStateActive,
				},
				wantCode: &domain.PasswordlessInitCode{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Expiration: time.Hour,
					CodeID:     "code1",
					Code:       "a",
					State:      domain.PasswordlessInitCodeStateActive,
				},
			},
		},
		{
			name: "add human (with phone), ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectPush(
								newAddHumanEvent("$plain$x$password", false, true, "+41711234567", AllowedLanguage),
								user.NewHumanInitialCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour*1,
									"",
								),
								user.NewHumanPhoneCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour*1),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Password: &domain.Password{
								SecretString:   "password",
								ChangeRequired: false,
							},
							Email: &domain.Email{
								EmailAddress: "email@test.ch",
							},
							Phone: &domain.Phone{
								PhoneNumber: "+41711234567",
							},
						},
						secretGenerator: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: AllowedLanguage,
					},
					Email: &domain.Email{
						EmailAddress: "email@test.ch",
					},
					Phone: &domain.Phone{
						PhoneNumber: "+41711234567",
					},
					State: domain.UserStateInitial,
				},
			},
		},
		{
			name: "add human (with verified phone), ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectPush(
								newAddHumanEvent("$plain$x$password", false, true, "+41711234567", AllowedLanguage),
								user.NewHumanInitialCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour*1,
									"",
								),
								user.NewHumanPhoneVerifiedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Password: &domain.Password{
								SecretString:   "password",
								ChangeRequired: false,
							},
							Email: &domain.Email{
								EmailAddress: "email@test.ch",
							},
							Phone: &domain.Phone{
								PhoneNumber:     "+41711234567",
								IsPhoneVerified: true,
							},
						},
						secretGenerator: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: AllowedLanguage,
					},
					Email: &domain.Email{
						EmailAddress: "email@test.ch",
					},
					Phone: &domain.Phone{
						PhoneNumber: "+41711234567",
					},
					State: domain.UserStateInitial,
				},
			},
		},
		{
			name: "add human (with undefined preferred language), ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectPush(
								newAddHumanEvent("$plain$x$password", false, true, "", language.Und),
								user.NewHumanInitialCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour*1,
									"",
								),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName: "firstname",
								LastName:  "lastname",
							},
							Password: &domain.Password{
								SecretString:   "password",
								ChangeRequired: false,
							},
							Email: &domain.Email{
								EmailAddress: "email@test.ch",
							},
						},
						secretGenerator: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: language.Und,
					},
					Email: &domain.Email{
						EmailAddress: "email@test.ch",
					},
					State: domain.UserStateInitial,
				},
			},
		},
		{
			name: "add human (with unsupported preferred language), ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectPush(
								newAddHumanEvent("$plain$x$password", false, true, "", UnsupportedLanguage),
								user.NewHumanInitialCodeAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("a"),
									},
									time.Hour*1,
									"",
								),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: UnsupportedLanguage,
							},
							Password: &domain.Password{
								SecretString:   "password",
								ChangeRequired: false,
							},
							Email: &domain.Email{
								EmailAddress: "email@test.ch",
							},
						},
						secretGenerator: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: UnsupportedLanguage,
					},
					Email: &domain.Email{
						EmailAddress: "email@test.ch",
					},
					State: domain.UserStateInitial,
				},
			},
		},
		{
			name: "add human (with idp), ok",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewIDPConfigAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"idpID",
										"name",
										domain.IDPConfigTypeOIDC,
										domain.IDPConfigStylingTypeUnspecified,
										false,
									),
								),
								eventFromEventPusher(
									org.NewIDPOIDCConfigAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"clientID",
										"idpID",
										"issuer",
										"authEndpoint",
										"tokenEndpoint",
										nil,
										domain.OIDCMappingFieldUnspecified,
										domain.OIDCMappingFieldUnspecified,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewIDPConfigAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"idpID",
										"name",
										domain.IDPConfigTypeOIDC,
										domain.IDPConfigStylingTypeUnspecified,
										false,
									),
								),
								eventFromEventPusher(
									org.NewIDPOIDCConfigAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"clientID",
										"idpID",
										"issuer",
										"authEndpoint",
										"tokenEndpoint",
										nil,
										domain.OIDCMappingFieldUnspecified,
										domain.OIDCMappingFieldUnspecified,
									),
								),
								eventFromEventPusher(
									org.NewIdentityProviderAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"idpID",
										domain.IdentityProviderTypeOrg,
									),
								),
							),
							expectPush(
								newAddHumanEvent("", false, true, "", AllowedLanguage),
								user.NewUserIDPLinkAddedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"idpID",
									"name",
									"externalID",
								),
								user.NewHumanEmailVerifiedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName:         "firstname",
								LastName:          "lastname",
								PreferredLanguage: AllowedLanguage,
							},
							Email: &domain.Email{
								EmailAddress:    "email@test.ch",
								IsEmailVerified: true,
							},
						},
						links: []*domain.UserIDPLink{
							{
								IDPConfigID:    "idpID",
								ExternalUserID: "externalID",
								DisplayName:    "name",
							},
						},
						secretGenerator: GetMockSecretGenerator(t),
					}
			},
			res: res{
				wantHuman: &domain.Human{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Username: "username",
					Profile: &domain.Profile{
						FirstName:         "firstname",
						LastName:          "lastname",
						DisplayName:       "firstname lastname",
						PreferredLanguage: AllowedLanguage,
					},
					Email: &domain.Email{
						EmailAddress:    "email@test.ch",
						IsEmailVerified: true,
					},
					State: domain.UserStateActive,
				},
			},
		},
		{
			name: "add human (with idp, creation not allowed), precondition error",
			given: func(t *testing.T) (fields, args) {
				return fields{
						eventstore: eventstoreExpect(
							t,
							expectFilter(
								eventFromEventPusher(
									org.NewDomainPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										true,
										true,
										true,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
										&user.NewAggregate("user1", "org1").Aggregate,
										1,
										false,
										false,
										false,
										false,
									),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewIDPConfigAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"idpID",
										"name",
										domain.IDPConfigTypeOIDC,
										domain.IDPConfigStylingTypeUnspecified,
										false,
									),
								),
								eventFromEventPusher(
									org.NewIDPOIDCConfigAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"clientID",
										"idpID",
										"issuer",
										"authEndpoint",
										"tokenEndpoint",
										nil,
										domain.OIDCMappingFieldUnspecified,
										domain.OIDCMappingFieldUnspecified,
									),
								),
								eventFromEventPusher(
									func() eventstore.Command {
										e, _ := org.NewOIDCIDPChangedEvent(context.Background(),
											&org.NewAggregate("org1").Aggregate,
											"config1",
											[]idp.OIDCIDPChanges{
												idp.ChangeOIDCOptions(idp.OptionChanges{IsCreationAllowed: gu.Ptr(false)}),
											},
										)
										return e
									}(),
								),
							),
							expectFilter(
								eventFromEventPusher(
									org.NewIDPConfigAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"idpID",
										"name",
										domain.IDPConfigTypeOIDC,
										domain.IDPConfigStylingTypeUnspecified,
										false,
									),
								),
								eventFromEventPusher(
									org.NewIDPOIDCConfigAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"clientID",
										"idpID",
										"issuer",
										"authEndpoint",
										"tokenEndpoint",
										nil,
										domain.OIDCMappingFieldUnspecified,
										domain.OIDCMappingFieldUnspecified,
									),
								),
								eventFromEventPusher(
									func() eventstore.Command {
										e, _ := org.NewOIDCIDPChangedEvent(context.Background(),
											&org.NewAggregate("org1").Aggregate,
											"config1",
											[]idp.OIDCIDPChanges{
												idp.ChangeOIDCOptions(idp.OptionChanges{IsCreationAllowed: gu.Ptr(false)}),
											},
										)
										return e
									}(),
								),
								eventFromEventPusher(
									org.NewIdentityProviderAddedEvent(context.Background(),
										&org.NewAggregate("org1").Aggregate,
										"idpID",
										domain.IdentityProviderTypeOrg,
									),
								),
							),
						),
						idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
						userPasswordHasher: mockPasswordHasher("x"),
					},
					args{
						ctx:   context.Background(),
						orgID: "org1",
						human: &domain.Human{
							Username: "username",
							Profile: &domain.Profile{
								FirstName: "firstname",
								LastName:  "lastname",
							},
							Email: &domain.Email{
								EmailAddress:    "email@test.ch",
								IsEmailVerified: true,
							},
						},
						links: []*domain.UserIDPLink{
							{
								IDPConfigID:    "idpID",
								ExternalUserID: "externalID",
								DisplayName:    "name",
							},
						},
						secretGenerator: GetMockSecretGenerator(t),
					}
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, a := tt.given(t)
			r := &Commands{
				eventstore:         f.eventstore,
				userPasswordHasher: f.userPasswordHasher,
			}
			id_generator.SetGenerator(f.idGenerator)
			gotHuman, gotCode, err := r.ImportHuman(a.ctx, a.orgID, a.human, a.passwordless, a.links, a.secretGenerator, a.secretGenerator, a.secretGenerator, a.secretGenerator)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.wantHuman.PreferredLanguage, gotHuman.PreferredLanguage)
				assert.Equal(t, tt.res.wantCode, gotCode)
			}
		})
	}
}

func TestCommandSide_HumanMFASkip(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
			userID string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "skip mfa init, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								AllowedLanguage,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						user.NewHumanMFAInitSkippedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := r.HumanSkipMFAInit(tt.args.ctx, tt.args.userID, tt.args.orgID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_HumanSignOut(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx     context.Context
			agentID string
			userIDs []string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "agentid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:     context.Background(),
				agentID: "",
				userIDs: []string{"user1"},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "userids missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:     context.Background(),
				agentID: "agent1",
				userIDs: []string{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:     context.Background(),
				agentID: "agent1",
				userIDs: []string{"user1"},
			},
			res: res{},
		},
		{
			name: "human sign out, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								AllowedLanguage,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						user.NewHumanSignedOutEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"agent1",
						),
					),
				),
			},
			args: args{
				ctx:     context.Background(),
				agentID: "agent1",
				userIDs: []string{"user1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "human sign out multiple users, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								AllowedLanguage,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user2", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								AllowedLanguage,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						user.NewHumanSignedOutEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"agent1",
						),
						user.NewHumanSignedOutEvent(context.Background(),
							&user.NewAggregate("user2", "org1").Aggregate,
							"agent1",
						),
					),
				),
			},
			args: args{
				ctx:     context.Background(),
				agentID: "agent1",
				userIDs: []string{"user1", "user2"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := r.HumansSignOut(tt.args.ctx, tt.args.agentID, tt.args.userIDs)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newAddMachineEvent(userLoginMustBeDomain bool, accessTokenType domain.OIDCTokenType) *user.MachineAddedEvent {
	return user.NewMachineAddedEvent(context.Background(),
		&user.NewAggregate("user1", "org1").Aggregate,
		"username",
		"name",
		"description",
		userLoginMustBeDomain,
		accessTokenType,
	)
}

func newAddHumanEvent(password string, changeRequired, userLoginMustBeDomain bool, phone string, preferredLanguage language.Tag) *user.HumanAddedEvent {
	event := user.NewHumanAddedEvent(context.Background(),
		&user.NewAggregate("user1", "org1").Aggregate,
		"username",
		"firstname",
		"lastname",
		"",
		"firstname lastname",
		preferredLanguage,
		domain.GenderUnspecified,
		"email@test.ch",
		userLoginMustBeDomain,
	)
	if password != "" {
		event.AddPasswordData(password, changeRequired)
	}
	if phone != "" {
		event.AddPhoneData(domain.PhoneNumber(phone))
	}
	return event
}

func newRegisterHumanEvent(username, password string, changeRequired, userLoginMustBeUnique bool, phone string, preferredLanguage language.Tag) *user.HumanRegisteredEvent {
	event := user.NewHumanRegisteredEvent(context.Background(),
		&user.NewAggregate("user1", "org1").Aggregate,
		username,
		"firstname",
		"lastname",
		"",
		"firstname lastname",
		preferredLanguage,
		domain.GenderUnspecified,
		"email@test.ch",
		userLoginMustBeUnique,
		"",
	)
	if password != "" {
		event.AddPasswordData(password, changeRequired)
	}
	if phone != "" {
		event.AddPhoneData(domain.PhoneNumber(phone))
	}
	return event
}

func TestAddHumanCommand(t *testing.T) {
	type fields struct {
		idGenerator id_generator.Generator
	}
	type args struct {
		human         *AddHuman
		orgID         string
		hasher        *crypto.Hasher
		filter        preparation.FilterToQueryReducer
		codeAlg       crypto.EncryptionAlgorithm
		allowInitMail bool
	}
	agg := user.NewAggregate("id", "ro")
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Want
	}{
		{
			name: "invalid email",
			args: args{
				human: &AddHuman{
					Email: Email{
						Address: "invalid",
					},
					PreferredLanguage: AllowedLanguage,
				},
				orgID: "ro",
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "EMAIL-599BI", "Errors.User.Email.Invalid"),
			},
		},
		{
			name: "invalid first name",
			args: args{
				human: &AddHuman{
					Username: "username",
					Email: Email{
						Address: "support@zitadel.com",
					},
					PreferredLanguage: AllowedLanguage,
				},
				orgID: "ro",
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "USER-UCej2", "Errors.User.Profile.FirstNameEmpty"),
			},
		},
		{
			name: "invalid last name",
			args: args{
				human: &AddHuman{
					Username:          "username",
					FirstName:         "hurst",
					Email:             Email{Address: "support@zitadel.com"},
					PreferredLanguage: AllowedLanguage,
				},
				orgID: "ro",
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "USER-4hB7d", "Errors.User.Profile.LastNameEmpty"),
			},
		},
		{
			name: "unsupported password hash encoding",
			args: args{
				human: &AddHuman{
					Email:               Email{Address: "support@zitadel.com", Verified: true},
					FirstName:           "gigi",
					LastName:            "giraffe",
					EncodedPasswordHash: "$foo$x$password",
					Username:            "username",
					PreferredLanguage:   AllowedLanguage,
				},
				orgID:  "ro",
				hasher: mockPasswordHasher("x"),
			},
			want: Want{
				ValidationErr: zerrors.ThrowInvalidArgument(nil, "USER-JDk4t", "Errors.User.Password.NotSupported"),
			},
		},
		{
			name: "invalid password",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args: args{
				human: &AddHuman{
					Email:             Email{Address: "support@zitadel.com"},
					FirstName:         "gigi",
					LastName:          "giraffe",
					Password:          "short",
					Username:          "username",
					PreferredLanguage: AllowedLanguage,
				},
				orgID: "ro",
				filter: NewMultiFilter().Append(
					func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{}, nil
					}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewDomainPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									true,
									true,
									true,
								),
							}, nil
						}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewPasswordComplexityPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									8,
									true,
									true,
									true,
									true,
								),
							}, nil
						}).
					Filter(),
			},
			want: Want{
				CreateErr: zerrors.ThrowInvalidArgument(nil, "COMMA-HuJf6", "Errors.User.PasswordComplexityPolicy.MinLength"),
			},
		},
		{
			name: "correct",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args: args{
				human: &AddHuman{
					Email:             Email{Address: "support@zitadel.com", Verified: true},
					FirstName:         "gigi",
					LastName:          "giraffe",
					Password:          "password",
					Username:          "username",
					PreferredLanguage: AllowedLanguage,
				},
				orgID:   "ro",
				hasher:  mockPasswordHasher("x"),
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				filter: NewMultiFilter().Append(
					func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{}, nil
					}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewDomainPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									true,
									true,
									true,
								),
							}, nil
						}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewPasswordComplexityPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									2,
									false,
									false,
									false,
									false,
								),
							}, nil
						}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					func() *user.HumanAddedEvent {
						event := user.NewHumanAddedEvent(
							context.Background(),
							&agg.Aggregate,
							"username",
							"gigi",
							"giraffe",
							"",
							"gigi giraffe",
							AllowedLanguage,
							0,
							"support@zitadel.com",
							true,
						)
						event.AddPasswordData("$plain$x$password", false)
						return event
					}(),
					user.NewHumanEmailVerifiedEvent(context.Background(), &agg.Aggregate),
				},
			},
		},
		{
			name: "undefined preferred language, ok",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args: args{
				human: &AddHuman{
					Email:     Email{Address: "support@zitadel.com", Verified: true},
					FirstName: "gigi",
					LastName:  "giraffe",
					Password:  "password",
					Username:  "username",
				},
				orgID:   "ro",
				hasher:  mockPasswordHasher("x"),
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				filter: NewMultiFilter().Append(
					func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{}, nil
					}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewDomainPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									true,
									true,
									true,
								),
							}, nil
						}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewPasswordComplexityPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									2,
									false,
									false,
									false,
									false,
								),
							}, nil
						}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					func() *user.HumanAddedEvent {
						event := user.NewHumanAddedEvent(
							context.Background(),
							&agg.Aggregate,
							"username",
							"gigi",
							"giraffe",
							"",
							"gigi giraffe",
							language.Und,
							0,
							"support@zitadel.com",
							true,
						)
						event.AddPasswordData("$plain$x$password", false)
						return event
					}(),
					user.NewHumanEmailVerifiedEvent(context.Background(), &agg.Aggregate),
				},
			},
		},
		{
			name: "unsupported preferred language, ok",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args: args{
				human: &AddHuman{
					Email:             Email{Address: "support@zitadel.com", Verified: true},
					FirstName:         "gigi",
					LastName:          "giraffe",
					Password:          "password",
					Username:          "username",
					PreferredLanguage: UnsupportedLanguage,
				},
				orgID:   "ro",
				hasher:  mockPasswordHasher("x"),
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				filter: NewMultiFilter().Append(
					func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{}, nil
					}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewDomainPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									true,
									true,
									true,
								),
							}, nil
						}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewPasswordComplexityPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									2,
									false,
									false,
									false,
									false,
								),
							}, nil
						}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					func() *user.HumanAddedEvent {
						event := user.NewHumanAddedEvent(
							context.Background(),
							&agg.Aggregate,
							"username",
							"gigi",
							"giraffe",
							"",
							"gigi giraffe",
							UnsupportedLanguage,
							0,
							"support@zitadel.com",
							true,
						)
						event.AddPasswordData("$plain$x$password", false)
						return event
					}(),
					user.NewHumanEmailVerifiedEvent(context.Background(), &agg.Aggregate),
				},
			},
		},
		{
			name: "hashed password",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args: args{
				human: &AddHuman{
					Email:               Email{Address: "support@zitadel.com", Verified: true},
					FirstName:           "gigi",
					LastName:            "giraffe",
					EncodedPasswordHash: "$plain$x$password",
					Username:            "username",
					PreferredLanguage:   AllowedLanguage,
				},
				orgID:   "ro",
				hasher:  mockPasswordHasher("x"),
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				filter: NewMultiFilter().Append(
					func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{}, nil
					}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewDomainPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									true,
									true,
									true,
								),
							}, nil
						}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewPasswordComplexityPolicyAddedEvent(
									ctx,
									&org.NewAggregate("id").Aggregate,
									2,
									false,
									false,
									false,
									false,
								),
							}, nil
						}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					func() *user.HumanAddedEvent {
						event := user.NewHumanAddedEvent(
							context.Background(),
							&agg.Aggregate,
							"username",
							"gigi",
							"giraffe",
							"",
							"gigi giraffe",
							AllowedLanguage,
							0,
							"support@zitadel.com",
							true,
						)
						event.AddPasswordData("$plain$x$password", false)
						return event
					}(),
					user.NewHumanEmailVerifiedEvent(context.Background(), &agg.Aggregate),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{}
			id_generator.SetGenerator(tt.fields.idGenerator)
			AssertValidation(t, context.Background(), c.AddHumanCommand(tt.args.human, tt.args.orgID, tt.args.hasher, tt.args.codeAlg, tt.args.allowInitMail), tt.args.filter, tt.want)
		})
	}
}
