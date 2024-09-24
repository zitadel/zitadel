package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/passwap"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user/authenticator"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func filterSchemaUserPasswordExisting() expect {
	return expectFilter(
		eventFromEventPusher(
			authenticator.NewPasswordCreatedEvent(
				context.Background(),
				&authenticator.NewAggregate("user1", "org1").Aggregate,
				"user1",
				"$plain$x$password",
				false,
			),
		),
	)
}

func filterPasswordComplexityPolicyExisting() expect {
	return expectFilter(
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
	)
}

func TestCommands_SetSchemaUserPassword(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		userPasswordHasher *crypto.Hasher
		checkPermission    domain.PermissionCheck
		codeAlg            crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx  context.Context
		user *SetSchemaUserPassword
	}
	type res struct {
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no userID, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-aS3Vz5t6BS", "Errors.IDMissing"))
				},
			},
		},
		{
			"no password, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty"))
				},
			},
		},
		{
			"user not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:   "notexisting",
					Password: "password",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-TODO", "Errors.User.Password.NotFound"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:   "user1",
					Password: "password",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"password added, ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					filterSchemaUserExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password",
							false,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:         "user1",
					Password:       "password",
					ChangeRequired: false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, complexity failed",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password",
							false,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:         "user1",
					Password:       "password",
					ChangeRequired: false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password",
							false,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:         "user1",
					Password:       "password",
					ChangeRequired: false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, changeRequired, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password",
							true,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:         "user1",
					Password:       "password",
					ChangeRequired: true,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, encoded, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password2",
							false,
						),
					),
				),
				checkPermission:    newMockPermissionCheckAllowed(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:              "user1",
					Password:            "passwordnotused",
					EncodedPasswordHash: "$plain$x$password2",
					ChangeRequired:      false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, current password, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password2",
							false,
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:          "user1",
					Password:        "password2",
					CurrentPassword: "password",
					ChangeRequired:  false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, current password, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password2",
							false,
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:          "user1",
					Password:        "password2",
					CurrentPassword: "password",
					ChangeRequired:  false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		}, {
			"password set, current password, failed",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:          "user1",
					Password:        "password2",
					CurrentPassword: "notreally",
					ChangeRequired:  false,
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(passwap.ErrPasswordMismatch, "COMMAND-3M0fs", "Errors.User.Password.Invalid"))
				},
			},
		},
		{
			"password set, code, ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							authenticator.NewPasswordCreatedEvent(
								context.Background(),
								&authenticator.NewAggregate("user1", "org1").Aggregate,
								"user1",
								"$plain$x$password",
								false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							authenticator.NewPasswordCodeAddedEvent(context.Background(),
								&authenticator.NewAggregate("user1", "org1").Aggregate,
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
							),
						),
					),
					filterPasswordComplexityPolicyExisting(),
					expectPush(
						authenticator.NewPasswordCreatedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							"user1",
							"$plain$x$password2",
							false,
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:           "user1",
					Password:         "password2",
					VerificationCode: "code",
					ChangeRequired:   false,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password set, code, failed",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							authenticator.NewPasswordCreatedEvent(
								context.Background(),
								&authenticator.NewAggregate("user1", "org1").Aggregate,
								"user1",
								"$plain$x$password",
								false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							authenticator.NewPasswordCodeAddedEvent(context.Background(),
								&authenticator.NewAggregate("user1", "org1").Aggregate,
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
							),
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:           "user1",
					Password:         "password2",
					VerificationCode: "notreally",
					ChangeRequired:   false,
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "CODE-woT0xc", "Errors.User.Code.Invalid"))
				},
			},
		},
		{
			"password set, code, no code",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				codeAlg:            crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &SetSchemaUserPassword{
					UserID:           "user1",
					Password:         "password2",
					VerificationCode: "notreally",
					ChangeRequired:   false,
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-TODO", "Errors.User.Code.NotFound"))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:         tt.fields.eventstore(t),
				checkPermission:    tt.fields.checkPermission,
				userPasswordHasher: tt.fields.userPasswordHasher,
				userEncryption:     tt.fields.codeAlg,
			}
			details, err := c.SetSchemaUserPassword(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_RequestSchemaUserPasswordReset(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
		newCode         encrypedCodeFunc
	}
	type args struct {
		ctx  context.Context
		user *RequestSchemaUserPasswordReset
	}
	type res struct {
		details   *domain.ObjectDetails
		plainCode string
		err       func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no userID, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:  authz.NewMockContext("instanceID", "", ""),
				user: &RequestSchemaUserPasswordReset{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing"))
				},
			},
		},
		{
			"password not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &RequestSchemaUserPasswordReset{
					UserID: "notexisting",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-TODO", "Errors.User.Password.NotFound"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &RequestSchemaUserPasswordReset{
					UserID: "user1",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"password reset, email, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					expectPush(
						authenticator.NewPasswordCodeAddedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("code"),
							},
							10*time.Minute,
							domain.NotificationTypeEmail,
							"https://example.com/password/changey?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
							false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("code", 10*time.Minute),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &RequestSchemaUserPasswordReset{
					UserID:           "user1",
					NotificationType: domain.NotificationTypeEmail,
					URLTemplate:      "https://example.com/password/changey?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password reset, sms, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					expectPush(
						authenticator.NewPasswordCodeAddedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("code"),
							},
							10*time.Minute,
							domain.NotificationTypeSms,
							"",
							false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("code", 10*time.Minute),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &RequestSchemaUserPasswordReset{
					UserID:           "user1",
					NotificationType: domain.NotificationTypeSms,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"password reset, returned, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					expectPush(
						authenticator.NewPasswordCodeAddedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("code"),
							},
							10*time.Minute,
							domain.NotificationTypeEmail,
							"",
							true,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockEncryptedCode("code", 10*time.Minute),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				user: &RequestSchemaUserPasswordReset{
					UserID:     "user1",
					ReturnCode: true,
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				plainCode: "code",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:       tt.fields.eventstore(t),
				checkPermission:  tt.fields.checkPermission,
				newEncryptedCode: tt.fields.newCode,
			}
			details, err := c.RequestSchemaUserPasswordReset(tt.args.ctx, tt.args.user)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}
			if tt.res.plainCode != "" {
				assert.Equal(t, tt.res.plainCode, tt.args.user.PlainCode)
			}
		})
	}
}

func TestCommands_DeleteSchemaUserPassword(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
	}
	type res struct {
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no ID, error",
			fields{
				eventstore:      expectEventstore(),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-PoSU5BOZCi", "Errors.IDMissing"))
				},
			},
		},
		{
			"password not existing, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "notexisting",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
				},
			},
		},
		{
			"password already removed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							authenticator.NewPasswordCreatedEvent(
								context.Background(),
								&authenticator.NewAggregate("user1", "org1").Aggregate,
								"id1",
								"hash",
								false,
							),
						),
						eventFromEventPusher(
							authenticator.NewPasswordDeletedEvent(
								context.Background(),
								&authenticator.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "user1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "TODO", "TODO"))
				},
			},
		},
		{
			"no permission, error",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "user1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			"password removed, ok",
			fields{
				eventstore: expectEventstore(
					filterSchemaUserPasswordExisting(),
					expectPush(
						authenticator.NewPasswordDeletedEvent(
							context.Background(),
							&authenticator.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "user1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			details, err := c.DeleteSchemaUserPassword(tt.args.ctx, tt.args.resourceOwner, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.details, details)
			}
		})
	}
}
