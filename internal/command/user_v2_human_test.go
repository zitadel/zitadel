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

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommandSide_AddUserHuman(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		idGenerator        id.Generator
		userPasswordHasher *crypto.PasswordHasher
		codeAlg            crypto.EncryptionAlgorithm
		newCode            cryptoCodeFunc
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
					return errors.Is(err, caos_errs.ThrowInvalidArgument(nil, "COMMA-5Ky74", "Errors.Internal"))
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
					return errors.Is(err, caos_errs.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"))
				},
			},
		},
		{
			name: "with id, already exists, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, ""),
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
					PreferredLanguage: language.English,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-k2unb", "Errors.User.AlreadyExisting"))
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
					PreferredLanguage: language.English,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal"))
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
					PreferredLanguage: language.English,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInternal(nil, "USER-uQ96e", "Errors.Internal"))
				},
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
							language.English,
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
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				codeAlg:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:     mockCode("userinit", time.Hour),
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
					PreferredLanguage: language.English,
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
						newAddHumanEvent("$plain$x$password", false, true, ""),
						user.NewHumanInitialCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							1*time.Hour,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:            mockCode("userinit", time.Hour),
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
					PreferredLanguage: language.English,
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
						newAddHumanEvent("$plain$x$password", false, true, ""),
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
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:            mockCode("emailCode", time.Hour),
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
					PreferredLanguage: language.English,
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
						newAddHumanEvent("$plain$x$password", false, true, ""),
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
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:            mockCode("emailCode", time.Hour),
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
					PreferredLanguage: language.English,
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
						newAddHumanEvent("$plain$x$password", true, true, ""),
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
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
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
						newAddHumanEvent("$plain$x$password", true, true, ""),
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
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
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
						newAddHumanEvent("$plain$x$password", true, false, ""),
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
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
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
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInvalidArgument(nil, "COMMAND-SFd21", "Errors.User.DomainNotAllowedAsUsername"))
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
								language.English,
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
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
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
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567"),
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
				newCode:            mockCode("phonecode", time.Hour),
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
					PreferredLanguage: language.English,
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
						newAddHumanEvent("", false, true, "+41711234567"),
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
						),
						user.NewHumanPhoneVerifiedEvent(
							context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				codeAlg:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				newCode:     mockCode("userinit", time.Hour),
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
					PreferredLanguage: language.English,
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
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567"),
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
				newCode:            mockCode("phoneCode", time.Hour),
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
					PreferredLanguage: language.English,
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
						newAddHumanEvent("", false, true, ""),
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
				newCode:     mockCode("userinit", time.Hour),
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
					PreferredLanguage: language.English,
					Metadata: []*AddMetadataEntry{
						{
							Key:   "testKey",
							Value: []byte("testValue"),
						},
					},
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
				idGenerator:        tt.fields.idGenerator,
				newCode:            tt.fields.newCode,
			}
			err := r.AddUserHuman(tt.args.ctx, tt.args.orgID, tt.args.human, tt.args.allowInitMail)
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
