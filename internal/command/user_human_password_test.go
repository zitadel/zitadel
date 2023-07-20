package command

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/passwap"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommandSide_SetOneTimePassword(t *testing.T) {
	type fields struct {
		eventstore         *eventstore.Eventstore
		userPasswordHasher *crypto.PasswordHasher
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "missing permission, error",
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
					return errors.Is(err, caos_errs.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change password onetime, ok",
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
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordChangedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"$plain$x$password",
									true,
									"",
								),
							),
						},
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
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordChangedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"$plain$x$password",
									false,
									"",
								),
							),
						},
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
				eventstore:         tt.fields.eventstore,
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
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_SetPasswordWithVerifyCode(t *testing.T) {
	type fields struct {
		eventstore         *eventstore.Eventstore
		userEncryption     crypto.EncryptionAlgorithm
		userPasswordHasher *crypto.PasswordHasher
	}
	type args struct {
		ctx           context.Context
		userID        string
		code          string
		resourceOwner string
		password      string
		agentID       string
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "password missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "code not existing, precondition error",
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid code, invalid argument error",
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "set password, ok",
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
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordChangedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"$plain$x$password",
									false,
									"",
								),
							),
						},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore,
				userPasswordHasher: tt.fields.userPasswordHasher,
				userEncryption:     tt.fields.userEncryption,
			}
			got, err := r.SetPasswordWithVerifyCode(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.code, tt.args.password, tt.args.agentID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangePassword(t *testing.T) {
	type fields struct {
		userPasswordHasher *crypto.PasswordHasher
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		oldPassword   string
		newPassword   string
		agentID       string
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsPreconditionFailed,
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
				err: caos_errs.IsPreconditionFailed,
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
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
					[]*repository.Event{
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"$plain$x$password1",
								false,
								"",
							),
						),
					},
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
			got, err := r.ChangePassword(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.oldPassword, tt.args.newPassword, tt.args.agentID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RequestSetPassword(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx             context.Context
		userID          string
		resourceOwner   string
		notifyType      domain.NotificationType
		secretGenerator crypto.Generator
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "user initial, precondition error",
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "new code, ok",
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
						[]*repository.Event{
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
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:             context.Background(),
				userID:          "user1",
				resourceOwner:   "org1",
				secretGenerator: GetMockSecretGenerator(t),
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
			got, err := r.RequestSetPassword(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.notifyType, tt.args.secretGenerator)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_PasswordCodeSent(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "code sent, ok",
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
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordCodeSentEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := r.PasswordCodeSent(tt.args.ctx, tt.args.resourceOwner, tt.args.userID)
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
		eventstore         *eventstore.Eventstore
		userPasswordHasher *crypto.PasswordHasher
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		password      string
		authReq       *domain.AuthRequest
		lockoutPolicy *domain.LockoutPolicy
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				password:      "password",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "password missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "login policy not found, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "login policy login password not allowed, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing password empty, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "password not matching lockout policy not relevant, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordCheckFailedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&user.AuthRequestInfo{
										ID:          "request1",
										UserAgentID: "agent1",
									},
								),
							),
						},
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
				lockoutPolicy: &domain.LockoutPolicy{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "password not matching, max password attempts reached - user locked, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordCheckFailedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&user.AuthRequestInfo{
										ID:          "request1",
										UserAgentID: "agent1",
									},
								),
							),
							eventFromEventPusher(
								user.NewUserLockedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
								),
							),
						},
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
				lockoutPolicy: &domain.LockoutPolicy{
					MaxPasswordAttempts: 1,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "check password, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordCheckSucceededEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&user.AuthRequestInfo{
										ID:          "request1",
										UserAgentID: "agent1",
									},
								),
							),
						},
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
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordCheckSucceededEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&user.AuthRequestInfo{
										ID:          "request1",
										UserAgentID: "agent1",
									},
								),
							),
							eventFromEventPusher(
								user.NewHumanPasswordHashUpdatedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"$plain$x$password",
								),
							),
						},
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
			name: "regression test old version event",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanPasswordCheckSucceededEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									&user.AuthRequestInfo{
										ID:          "request1",
										UserAgentID: "agent1",
									},
								),
							),
							eventFromEventPusher(
								user.NewHumanPasswordHashUpdatedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"$plain$x$password",
								),
							),
						},
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
				eventstore:         tt.fields.eventstore,
				userPasswordHasher: tt.fields.userPasswordHasher,
			}
			err := r.HumanCheckPassword(tt.args.ctx, tt.args.resourceOwner, tt.args.userID, tt.args.password, tt.args.authReq, tt.args.lockoutPolicy)
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
			wantErr: caos_errs.ThrowInvalidArgument(passwap.ErrPasswordMismatch, "COMMAND-3M0fs", "Errors.User.Password.Invalid"),
		},
		{
			name:    "no change",
			args:    args{passwap.ErrPasswordNoChange},
			wantErr: caos_errs.ThrowPreconditionFailed(passwap.ErrPasswordNoChange, "COMMAND-Aesh5", "Errors.User.Password.NotChanged"),
		},
		{
			name:    "other",
			args:    args{io.ErrClosedPipe},
			wantErr: caos_errs.ThrowInternal(io.ErrClosedPipe, "COMMAND-CahN2", "Errors.Internal"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertPasswapErr(tt.args.err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
