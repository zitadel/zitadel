package command

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/passwap"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/senders/mock"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_SetOneTimePassword(t *testing.T) {
	type fields struct {
		eventstore         func(*testing.T) *eventstore.Eventstore
		userPasswordHasher *crypto.Hasher
		checkPermission    domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		password      string
		oneTime       bool
	}
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "missing permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				checkPermission:    newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				oneTime:       true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change password onetime, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
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
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							true,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				checkPermission:    newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				oneTime:       true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change password no one time, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
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
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							false,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				checkPermission:    newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				oneTime:       false,
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
				eventstore:         tt.fields.eventstore(t),
				userPasswordHasher: tt.fields.userPasswordHasher,
				checkPermission:    tt.fields.checkPermission,
			}
			got, err := r.SetPassword(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.password, tt.args.oneTime)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_SetPasswordWithVerifyCode(t *testing.T) {
	type fields struct {
		eventstore         func(*testing.T) *eventstore.Eventstore
		userEncryption     crypto.EncryptionAlgorithm
		userPasswordHasher *crypto.Hasher
		phoneCodeVerifier  func(ctx context.Context, id string) (senders.CodeGenerator, error)
	}
	type args struct {
		ctx            context.Context
		userID         string
		code           string
		resourceOwner  string
		password       string
		userAgentID    string
		changeRequired bool
	}
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "password missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "code not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				code:          "aa",
				resourceOwner: "org1",
				password:      "string",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid code, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								"",
							),
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				code:          "test",
				resourceOwner: "org1",
				password:      "password",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "set password, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
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
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							false,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				userEncryption:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				code:          "a",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set password with userAgentID, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
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
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							false,
							"userAgent1",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				userEncryption:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				code:          "a",
				userAgentID:   "userAgent1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set password with changeRequired, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
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
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							true,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				userEncryption:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:            context.Background(),
				userID:         "user1",
				resourceOwner:  "org1",
				password:       "password",
				code:           "a",
				userAgentID:    "",
				changeRequired: true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set password (external code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil,
								0,
								domain.NotificationTypeSms,
								"",
								"id",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeSentEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&senders.CodeGeneratorInfo{
									ID:             "id",
									VerificationID: "verificationID",
								},
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
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							false,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				userEncryption:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				phoneCodeVerifier: func(ctx context.Context, id string) (senders.CodeGenerator, error) {
					sender := mock.NewMockCodeGenerator(gomock.NewController(t))
					sender.EXPECT().VerifyCode("verificationID", "a").Return(nil)
					return sender, nil
				},
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				code:          "a",
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
				eventstore:         tt.fields.eventstore(t),
				userPasswordHasher: tt.fields.userPasswordHasher,
				userEncryption:     tt.fields.userEncryption,
				phoneCodeVerifier:  tt.fields.phoneCodeVerifier,
			}
			got, err := r.SetPasswordWithVerifyCode(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.code, tt.args.password, tt.args.userAgentID, tt.args.changeRequired)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangePassword(t *testing.T) {
	type fields struct {
		userPasswordHasher *crypto.Hasher
	}
	type args struct {
		ctx            context.Context
		userID         string
		resourceOwner  string
		oldPassword    string
		newPassword    string
		userAgentID    string
		changeRequired bool
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect []expect
		res    res
	}{
		{
			name:   "userid missing, invalid argument error",
			fields: fields{},
			args: args{
				ctx:           context.Background(),
				oldPassword:   "password",
				newPassword:   "password1",
				resourceOwner: "org1",
			},
			expect: []expect{},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "old password missing, invalid argument error",
			fields: fields{},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				newPassword:   "password1",
				resourceOwner: "org1",
			},
			expect: []expect{},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "new password missing, invalid argument error",
			fields: fields{},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				oldPassword:   "password",
				resourceOwner: "org1",
			},
			expect: []expect{},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "user not existing, precondition error",
			fields: fields{},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				oldPassword:   "password",
				newPassword:   "password1",
			},
			expect: []expect{
				expectFilter(),
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "existing password empty, precondition error",
			fields: fields{
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				oldPassword:   "password",
				newPassword:   "password1",
				resourceOwner: "org1",
			},
			expect: []expect{
				expectFilter(
					eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					),
				),
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "password not matching complexity policy, invalid argument error",
			fields: fields{
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				oldPassword:   "password-old",
				newPassword:   "password1",
				resourceOwner: "org1",
			},
			expect: []expect{
				expectFilter(
					eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					),
					eventFromEventPusher(
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
					eventFromEventPusher(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password-old",
							false,
							"")),
				),
				expectFilter(
					eventFromEventPusher(
						org.NewPasswordComplexityPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("org1").Aggregate,
							1,
							true,
							true,
							true,
							true,
						),
					),
				),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "password not matching, invalid argument error",
			fields: fields{
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				oldPassword:   "password-old",
				newPassword:   "password1",
				resourceOwner: "org1",
			},
			expect: []expect{
				expectFilter(
					eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					),
					eventFromEventPusher(
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
					eventFromEventPusher(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							false,
							"")),
				),
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "change password, ok",
			fields: fields{
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				oldPassword:   "password",
				newPassword:   "password1",
			},
			expect: []expect{
				expectFilter(
					eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					),
					eventFromEventPusher(
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
					eventFromEventPusher(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							false,
							"")),
				),
				expectFilter(
					eventFromEventPusher(
						org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
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
						&user.NewAggregate("user1", "org1").Aggregate,
						"$plain$x$password1",
						false,
						"",
					),
				),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change password with userAgentID, ok",
			fields: fields{
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				oldPassword:   "password",
				newPassword:   "password1",
				userAgentID:   "userAgent1",
			},
			expect: []expect{
				expectFilter(
					eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					),
					eventFromEventPusher(
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
					eventFromEventPusher(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							false,
							"")),
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
						&user.NewAggregate("user1", "org1").Aggregate,
						"$plain$x$password1",
						false,
						"userAgent1",
					),
				),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change password with changeRequired, ok",
			fields: fields{
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:            context.Background(),
				userID:         "user1",
				resourceOwner:  "org1",
				oldPassword:    "password",
				newPassword:    "password1",
				userAgentID:    "",
				changeRequired: true,
			},
			expect: []expect{
				expectFilter(
					eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					),
					eventFromEventPusher(
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
					eventFromEventPusher(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
							false,
							"")),
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
						&user.NewAggregate("user1", "org1").Aggregate,
						"$plain$x$password1",
						true,
						"",
					),
				),
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
				eventstore:         eventstoreExpect(t, tt.expect...),
				userPasswordHasher: tt.fields.userPasswordHasher,
			}
			got, err := r.ChangePassword(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.oldPassword, tt.args.newPassword, tt.args.userAgentID, tt.args.changeRequired)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RequestSetPassword(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
		newCode    encrypedCodeFunc
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		notifyType    domain.NotificationType
		authRequestID string
	}
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "user initial, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil, time.Hour*1,
								"",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "new code, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate)),
					),
					expectPush(
						user.NewHumanPasswordCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							domain.NotificationTypeEmail,
							"",
							"",
						),
					),
				),
				newCode: mockEncryptedCode("a", 1*time.Hour),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "new code with authRequestID, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate)),
					),
					expectPush(
						user.NewHumanPasswordCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
							time.Hour*1,
							domain.NotificationTypeEmail,
							"authRequestID",
							"",
						),
					),
				),
				newCode: mockEncryptedCode("a", 1*time.Hour),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				authRequestID: "authRequestID",
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
				eventstore:       tt.fields.eventstore(t),
				newEncryptedCode: tt.fields.newCode,
			}
			got, err := r.RequestSetPassword(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.notifyType, tt.args.authRequestID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_PasswordCodeSent(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		generatorInfo *senders.CodeGeneratorInfo
	}
	type res struct {
		err func(error) bool
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "code sent, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectPush(
						user.NewHumanPasswordCodeSentEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&senders.CodeGeneratorInfo{},
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				generatorInfo: &senders.CodeGeneratorInfo{},
			},
			res: res{},
		},
		{
			name: "code sent (external code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"+411234567",
							),
						),
					),
					expectPush(
						user.NewHumanPasswordCodeSentEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&senders.CodeGeneratorInfo{
								ID:             "generatorID",
								VerificationID: "verificationID",
							},
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				generatorInfo: &senders.CodeGeneratorInfo{
					ID:             "generatorID",
					VerificationID: "verificationID",
				},
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := r.PasswordCodeSent(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.generatorInfo)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_CheckPassword(t *testing.T) {
	type fields struct {
		eventstore         func(*testing.T) *eventstore.Eventstore
		userPasswordHasher *crypto.Hasher
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		password      string
		authReq       *domain.AuthRequest
	}
	type res struct {
		err func(error) bool
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				password:      "password",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "password missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "login policy not found, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "login policy login password not allowed, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "user locked, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "existing password empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				password:      "password",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "password not matching lockout policy not relevant, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"$plain$x$password",
								false,
								"")),
					),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								0, 0, false,
							)),
					),
					expectPush(
						user.NewHumanPasswordCheckFailedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "request1",
								UserAgentID: "agent1",
							},
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				password:      "password1",
				resourceOwner: "org1",
				authReq: &domain.AuthRequest{
					ID:      "request1",
					AgentID: "agent1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "password not matching, max password attempts reached - user locked, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"$plain$x$password",
								false,
								""),
						),
					),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								1, 1, false,
							)),
					),
					expectPush(
						user.NewHumanPasswordCheckFailedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "request1",
								UserAgentID: "agent1",
							},
						),
						user.NewUserLockedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				password:      "password1",
				resourceOwner: "org1",
				authReq: &domain.AuthRequest{
					ID:      "request1",
					AgentID: "agent1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "check password, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"$plain$x$password",
								false,
								"")),
					),
					expectFilter(),
					expectPush(
						user.NewHumanPasswordCheckSucceededEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "request1",
								UserAgentID: "agent1",
							},
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				authReq: &domain.AuthRequest{
					ID:      "request1",
					AgentID: "agent1",
				},
			},
			res: res{},
		},
		{
			name: "check password, ok, updated hash",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"$plain$v$password",
								false,
								"")),
					),
					expectFilter(),
					expectPush(
						user.NewHumanPasswordCheckSucceededEvent(
							context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "request1",
								UserAgentID: "agent1",
							},
						),
						user.NewHumanPasswordHashUpdatedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				authReq: &domain.AuthRequest{
					ID:      "request1",
					AgentID: "agent1",
				},
			},
			res: res{},
		},
		{
			name: "check password ok, locked in the mean time",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"$plain$x$password",
								false,
								"")),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewUserLockedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				authReq: &domain.AuthRequest{
					ID:      "request1",
					AgentID: "agent1",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "regression test old version event",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
								"",
								time.Hour*1,
								time.Hour*2,
								time.Hour*3,
								time.Hour*4,
								time.Hour*5,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							&user.HumanPasswordChangedEvent{
								BaseEvent: *eventstore.NewBaseEventForPush(
									context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									user.HumanPasswordChangedType,
								),
								Secret: &crypto.CryptoValue{
									CryptoType: crypto.TypeHash,
									Algorithm:  "plain",
									KeyID:      "",
									Crypted:    []byte("$plain$v$password"),
								},
								ChangeRequired: false,
							},
						),
					),
					expectFilter(),
					expectPush(
						user.NewHumanPasswordCheckSucceededEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "request1",
								UserAgentID: "agent1",
							},
						),
						user.NewHumanPasswordHashUpdatedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$x$password",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				authReq: &domain.AuthRequest{
					ID:      "request1",
					AgentID: "agent1",
				},
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				userPasswordHasher: tt.fields.userPasswordHasher,
			}
			err := r.HumanCheckPassword(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.password, tt.args.authReq)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func Test_convertPasswapErr(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "nil",
			args:    args{nil},
			wantErr: nil,
		},
		{
			name:    "mismatch",
			args:    args{passwap.ErrPasswordMismatch},
			wantErr: zerrors.ThrowInvalidArgument(passwap.ErrPasswordMismatch, "COMMAND-3M0fs", "Errors.User.Password.Invalid"),
		},
		{
			name:    "no change",
			args:    args{passwap.ErrPasswordNoChange},
			wantErr: zerrors.ThrowPreconditionFailed(passwap.ErrPasswordNoChange, "COMMAND-Aesh5", "Errors.User.Password.NotChanged"),
		},
		{
			name:    "other",
			args:    args{io.ErrClosedPipe},
			wantErr: zerrors.ThrowInternal(io.ErrClosedPipe, "COMMAND-CahN2", "Errors.Internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertPasswapErr(tt.args.err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
