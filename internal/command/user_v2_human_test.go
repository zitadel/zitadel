package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddUserHuman(t *testing.T) {
	defaultGenerators := &SecretGenerators{
		OTPSMS: &crypto.GeneratorConfig{
			Length:              8,
			Expiry:              time.Hour,
			IncludeLowerLetters: true,
			IncludeUpperLetters: true,
			IncludeDigits:       true,
			IncludeSymbols:      true,
		},
	}
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		idGenerator                 id.Generator
		userPasswordHasher          *crypto.Hasher
		newCode                     encrypedCodeFunc
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		checkPermission             domain.PermissionCheck
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		ctx             context.Context
		orgID           string
		human           *AddHuman
		secretGenerator crypto.Generator
		allowInitMail   bool
		codeAlg         crypto.EncryptionAlgorithm
	}
	type res struct {
		want          *domain.ObjectDetails
		wantID        string
		wantEmailCode string
		err           func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	cryptoAlg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	totpSecret := "TOTPSecret"
	totpSecretEnc, err := crypto.Encrypt([]byte(totpSecret), cryptoAlg)
	require.NoError(t, err)

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "orgid missing, invalid argument error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
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
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMA-095xh8fll1", "Errors.Internal"))
				},
			},
		},
		{
			name: "user invalid, invalid argument error",
			fields: fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
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
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
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
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-7yiox1isql", "Errors.User.AlreadyExisting"))
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
				checkPermission: newMockPermissionCheckAllowed(),
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
				checkPermission: newMockPermissionCheckAllowed(),
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
					return errors.Is(err, zerrors.ThrowInternal(nil, "USER-uQ96e", "Errors.Internal"))
				},
			},
		},
		{
			name: "register human (with initial code), ok",
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
						user.NewHumanRegisteredEvent(context.Background(),
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
							"userAgentID",
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
							"authRequestID",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
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
					Register:          true,
					UserAgentID:       "userAgentID",
					AuthRequestID:     "authRequestID",
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
			name: "add human (with initial code), no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "add human (email not verified, no password), ok (init code)",
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
							"",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
			name: "add human (email not verified, with password), ok (init code)",
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
						newAddHumanEvent("$plain$x$password", false, true, "", language.English),
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
				checkPermission:    newMockPermissionCheckAllowed(),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
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
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (email not verified, no password, no allowInitMail), ok (email verification with passwordInit)",
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
					expectPush(
						newAddHumanEvent("", false, true, "", language.English),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailverify"),
							},
							1*time.Hour,
							"",
							false,
							"",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("emailverify", time.Hour),
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
				allowInitMail:   false,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
						newAddHumanEvent("$plain$x$password", false, true, "", language.English),
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
				checkPermission:    newMockPermissionCheckAllowed(),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
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
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   false,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
						newAddHumanEvent("$plain$x$password", false, true, "", language.English),
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
				checkPermission:    newMockPermissionCheckAllowed(),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
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
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   false,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
			name: "add human email verified and password, ok",
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
						newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
						newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
						newAddHumanEvent("$plain$x$password", true, false, "", language.English),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
				checkPermission:    newMockPermissionCheckAllowed(),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
				checkPermission:    newMockPermissionCheckAllowed(),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
								"",
								"sid",
								"senderNumber",
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
								"",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567", language.English),
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
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				idGenerator:                 id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher:          mockPasswordHasher("x"),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phonecode", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with phone), ok (external)",
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
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
								"",
								"sid",
								"senderNumber",
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
								"verifiyServiceSid",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567", language.English),
						user.NewHumanEmailVerifiedEvent(
							context.Background(),
							&userAgg.Aggregate,
						),
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
							&userAgg.Aggregate,
							nil,
							0,
							"id",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				idGenerator:                 id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher:          mockPasswordHasher("x"),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phonecode", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
						newAddHumanEvent("", false, true, "+41711234567", language.English),
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
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with phone return code), ok",
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
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
								"",
								"sid",
								"senderNumber",
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
								"",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567", language.English),
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
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				idGenerator:                 id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher:          mockPasswordHasher("x"),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneCode", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
						newAddHumanEvent("", false, true, "", language.English),
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
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
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
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "register human with idp, unverified email, allow init mail, ok (verify mail)",
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
							org.NewGoogleIDPAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"idpID",
								"google",
								"clientID",
								nil,
								[]string{"openid"},
								idp.Options{},
							),
						),
					),
					expectPush(
						newRegisterHumanEvent("email@test.ch", "", false, true, "", language.English),
						user.NewHumanEmailCodeAddedEvent(
							context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("mailVerify"),
							},
							time.Hour,
							"authRequestID",
						),
						user.NewUserIDPLinkAddedEvent(
							context.Background(),
							&userAgg.Aggregate,
							"idpID",
							"displayName",
							"externalID",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("mailVerify", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "email@test.ch",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: language.English,
					Register:          true,
					Links: []*AddLink{
						{
							IDPID:         "idpID",
							DisplayName:   "displayName",
							IDPExternalID: "externalID",
						},
					},
					AuthRequestID: "authRequestID",
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "register human with idp, verified email, ok",
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
							org.NewGoogleIDPAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"idpID",
								"google",
								"clientID",
								nil,
								[]string{"openid"},
								idp.Options{},
							),
						),
					),
					expectPush(
						newRegisterHumanEvent("email@test.ch", "", false, true, "", language.English),
						user.NewHumanEmailVerifiedEvent(
							context.Background(),
							&userAgg.Aggregate,
						),
						user.NewUserIDPLinkAddedEvent(
							context.Background(),
							&userAgg.Aggregate,
							"idpID",
							"displayName",
							"externalID",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("mailVerify", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "email@test.ch",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PreferredLanguage: language.English,
					Register:          true,
					Links: []*AddLink{
						{
							IDPID:         "idpID",
							DisplayName:   "displayName",
							IDPExternalID: "externalID",
						},
					},
					AuthRequestID: "authRequestID",
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   false,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "register human with TOTPSecret, ok",
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
						user.NewHumanRegisteredEvent(context.Background(),
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
							"userAgentID",
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
							"authRequestID",
						),
						user.NewHumanOTPAddedEvent(context.Background(),
							&userAgg.Aggregate,
							totpSecretEnc,
						),
						user.NewHumanOTPVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
							"",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
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
					Register:          true,
					UserAgentID:       "userAgentID",
					AuthRequestID:     "authRequestID",
					TOTPSecret:        totpSecret,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
			name: "register human (validate domain), already verified",
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
								&org.NewAggregate("existing").Aggregate,
								"example.com",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username@example.com",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@example.com",
					},
					PreferredLanguage: language.English,
					Register:          true,
					UserAgentID:       "userAgentID",
					AuthRequestID:     "authRequestID",
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-SFd21", "Errors.User.DomainNotAllowedAsUsername"))
				},
			},
		},
		{
			name: "register human (validate domain), ok",
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
								&org.NewAggregate(userAgg.ResourceOwner).Aggregate,
								"example.com",
							),
						),
					),
					expectPush(
						user.NewHumanRegisteredEvent(context.Background(),
							&userAgg.Aggregate,
							"username@example.com",
							"firstname",
							"lastname",
							"",
							"firstname lastname",
							language.English,
							domain.GenderUnspecified,
							"email@example.com",
							false,
							"userAgentID",
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
							"authRequestID",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username@example.com",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@example.com",
					},
					PreferredLanguage: language.English,
					Register:          true,
					UserAgentID:       "userAgentID",
					AuthRequestID:     "authRequestID",
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
			name: "register human (validate domain, domain removed), ok",
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
								&org.NewAggregate("existing").Aggregate,
								"example.com",
							),
						),
						eventFromEventPusher(
							org.NewDomainRemovedEvent(context.Background(),
								&org.NewAggregate("existing").Aggregate,
								"example.com",
								true,
							),
						),
					),
					expectPush(
						user.NewHumanRegisteredEvent(context.Background(),
							&userAgg.Aggregate,
							"username@example.com",
							"firstname",
							"lastname",
							"",
							"firstname lastname",
							language.English,
							domain.GenderUnspecified,
							"email@example.com",
							false,
							"userAgentID",
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
							"authRequestID",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username@example.com",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@example.com",
					},
					PreferredLanguage: language.English,
					Register:          true,
					UserAgentID:       "userAgentID",
					AuthRequestID:     "authRequestID",
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
			name: "register human (validate domain, org removed), ok",
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
								&org.NewAggregate("existing").Aggregate,
								"example.com",
							),
						),
						eventFromEventPusher(
							org.NewOrgRemovedEvent(context.Background(),
								&org.NewAggregate("existing").Aggregate,
								"org",
								[]string{},
								false,
								[]string{},
								[]*domain.UserIDPLink{},
								[]string{},
							),
						),
					),
					expectPush(
						user.NewHumanRegisteredEvent(context.Background(),
							&userAgg.Aggregate,
							"username@example.com",
							"firstname",
							"lastname",
							"",
							"firstname lastname",
							language.English,
							domain.GenderUnspecified,
							"email@example.com",
							false,
							"userAgentID",
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
							"authRequestID",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				idGenerator:     id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:         mockEncryptedCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username@example.com",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@example.com",
					},
					PreferredLanguage: language.English,
					Register:          true,
					UserAgentID:       "userAgentID",
					AuthRequestID:     "authRequestID",
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				userPasswordHasher:          tt.fields.userPasswordHasher,
				idGenerator:                 tt.fields.idGenerator,
				newEncryptedCode:            tt.fields.newCode,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				checkPermission:             tt.fields.checkPermission,
				multifactors: domain.MultifactorConfigs{
					OTP: domain.OTPConfig{
						Issuer:    "zitadel.com",
						CryptoMFA: cryptoAlg,
					},
				},
			}
			err := r.AddUserHuman(tt.args.ctx, tt.args.orgID, tt.args.human, tt.args.allowInitMail, tt.args.codeAlg)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.human.Details)
				assert.Equal(t, tt.res.wantID, tt.args.human.ID)
				assert.Equal(t, tt.res.wantEmailCode, gu.Value(tt.args.human.EmailCode))
			}
		})
	}
}

func TestCommandSide_ChangeUserHuman(t *testing.T) {
	defaultGenerators := &SecretGenerators{
		OTPSMS: &crypto.GeneratorConfig{
			Length:              8,
			Expiry:              time.Hour,
			IncludeLowerLetters: true,
			IncludeUpperLetters: true,
			IncludeDigits:       true,
			IncludeSymbols:      true,
		},
	}
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		userPasswordHasher          *crypto.Hasher
		newCode                     encrypedCodeFunc
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		checkPermission             domain.PermissionCheck
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		ctx     context.Context
		orgID   string
		human   *ChangeHuman
		codeAlg crypto.EncryptionAlgorithm
	}
	type res struct {
		want          *domain.ObjectDetails
		wantEmailCode *string
		wantPhoneCode *string
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
			name: "domain policy not found, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectFilter(),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-79pv6e1q62", "Errors.Org.DomainPolicy.NotExisting"))
				},
			},
		},
		{
			name: "change human username, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human username, not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-ugjs0upun6", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "change human username, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
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
						user.NewUsernameChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"username",
							"changed",
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human username, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Username: gu.Ptr("username"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human profile, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Profile: &Profile{
						FirstName:         gu.Ptr("changedfn"),
						LastName:          gu.Ptr("changedln"),
						NickName:          gu.Ptr("changednn"),
						DisplayName:       gu.Ptr("changeddn"),
						PreferredLanguage: gu.Ptr(language.Afrikaans),
						Gender:            gu.Ptr(domain.GenderDiverse),
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human profile, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						func() eventstore.Command {
							cmd, _ := user.NewHumanProfileChangedEvent(context.Background(),
								&userAgg.Aggregate,
								[]user.ProfileChanges{
									user.ChangeFirstName("changedfn"),
									user.ChangeLastName("changedln"),
									user.ChangeNickName("changednn"),
									user.ChangeDisplayName("changeddn"),
									user.ChangePreferredLanguage(language.Afrikaans),
									user.ChangeGender(domain.GenderDiverse),
								},
							)
							return cmd
						}(),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Profile: &Profile{
						FirstName:         gu.Ptr("changedfn"),
						LastName:          gu.Ptr("changedln"),
						NickName:          gu.Ptr("changednn"),
						DisplayName:       gu.Ptr("changeddn"),
						PreferredLanguage: gu.Ptr(language.Afrikaans),
						Gender:            gu.Ptr(domain.GenderDiverse),
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human profile, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Profile: &Profile{
						FirstName:         gu.Ptr("firstname"),
						LastName:          gu.Ptr("lastname"),
						NickName:          gu.Ptr(""),
						DisplayName:       gu.Ptr("firstname lastname"),
						PreferredLanguage: gu.Ptr(language.English),
						Gender:            gu.Ptr(domain.GenderUnspecified),
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"changed@example.com",
						),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailCode"),
							},
							time.Hour,
							"",
							false,
							"",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("emailCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address: "changed@example.com",
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address: "email@test.ch",
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email verified, not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address:  "changed@example.com",
						Verified: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human email verified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"changed@example.com",
						),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address:  "changed@example.com",
						Verified: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email isVerified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Verified: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email returnCode, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"changed@test.com",
						),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailCode"),
							},
							time.Hour,
							"",
							true,
							"",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("emailCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address:    "changed@test.com",
						ReturnCode: true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
				wantEmailCode: gu.Ptr("emailCode"),
			},
		},
		{
			name: "change human phone, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
								"",
								"sid",
								"senderNumber",
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
								"",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"+41791234567",
						),
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							time.Hour,
							false,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneCode", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number: "+41791234567",
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human phone, ok (external)",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
								"",
								"sid",
								"senderNumber",
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
								"verifyServiceSid",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"+41791234567",
						),
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							nil,
							0,
							false,
							"id",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneCode", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number: "+41791234567",
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human phone verified, not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number:   "+41791234567",
						Verified: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human phone verified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"+41791234567",
						),
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number:   "+41791234567",
						Verified: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human phone isVerified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Verified: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human phone returnCode, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							instance.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
								"",
								"sid",
								"senderNumber",
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "id", Crypted: []byte("crypted")},
								"",
							),
						),
						eventFromEventPusher(
							instance.NewSMSConfigActivatedEvent(
								context.Background(),
								&instance.NewAggregate("instanceID").Aggregate,
								"id",
							),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"+41791234567",
						),
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							time.Hour,
							true,
							"",
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("phoneCode", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number:     "+41791234567",
						ReturnCode: true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
				wantPhoneCode: gu.Ptr("phoneCode"),
			},
		},
		{
			name: "password change, no password, invalid argument error",
			fields: fields{
				eventstore:         expectEventstore(),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty"))
				},
			},
		},
		{
			name: "change human password, not initialized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&userAgg.Aggregate,
								nil, time.Hour*1,
								"",
							),
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						OldPassword:    "password",
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-M9dse", "Errors.User.NotInitialised"))
				},
			},
		},
		{
			name: "change human password, not in complexity",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								true,
								false,
								false,
							),
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						OldPassword:    "password",
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "DOMAIN-VoaRj", "Errors.User.PasswordComplexityPolicy.HasUpper"))
				},
			},
		},
		{
			name: "change human password, empty",
			fields: fields{
				eventstore:         expectEventstore(),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						OldPassword:    "password",
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty"))
				},
			},
		},
		{
			name: "change human password, not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				checkPermission:    newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human password, permission, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
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
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				checkPermission:    newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						ChangeRequired: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human password, old password, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
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
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						OldPassword:    "password",
						ChangeRequired: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human password, old password, failed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						OldPassword:    "wrong",
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Invalid"))
				},
			},
		},
		{
			name: "change human password, password code, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
								"",
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
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						PasswordCode:   "code",
						ChangeRequired: true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human password, password code, wrong code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
								"",
							),
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						PasswordCode:   "wrong",
						ChangeRequired: true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "CODE-woT0xc", "Errors.User.Code.Invalid"))
				},
			},
		},
		{
			name: "change human password, password code, not matching policy",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
								"",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								true,
								true,
								true,
								true,
							),
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       "password2",
						PasswordCode:   "code",
						ChangeRequired: true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "change human password encoded, password code, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
								"",
							),
						),
					),
					expectPush(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						EncodedPasswordHash: "$plain$x$password2",
						PasswordCode:        "code",
						ChangeRequired:      true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human password and password encoded, password code, encoded used",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
								"",
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
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:            "passwordnotused",
						EncodedPasswordHash: "$plain$x$password2",
						PasswordCode:        "code",
						ChangeRequired:      true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				userPasswordHasher:          tt.fields.userPasswordHasher,
				newEncryptedCode:            tt.fields.newCode,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				checkPermission:             tt.fields.checkPermission,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				userEncryption:              tt.args.codeAlg,
			}
			err := r.ChangeUserHuman(tt.args.ctx, tt.args.human, tt.args.codeAlg)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.human.Details)
				assert.Equal(t, tt.res.wantEmailCode, tt.args.human.EmailCode)
				assert.Equal(t, tt.res.wantPhoneCode, tt.args.human.PhoneCode)
			}
		})
	}
}
