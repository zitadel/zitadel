package command

import (
	"context"
	"testing"
	"time"

	"golang.org/x/text/language"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_RequestPasswordReset(t *testing.T) {
	type fields struct {
		checkPermission domain.PermissionCheck
		eventstore      func(t *testing.T) *eventstore.Eventstore
		userEncryption  crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-SAFdda", "Errors.User.IDMissing"),
		},
		{
			name: "user not existing",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-SAF4f", "Errors.User.NotFound"),
		},
		{
			name: "user not initialized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "keyID", Crypted: []byte("code")}, 10*time.Second),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sfe4g", "Errors.User.NotInitialised"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				checkPermission: tt.fields.checkPermission,
				eventstore:      tt.fields.eventstore(t),
				userEncryption:  tt.fields.userEncryption,
			}
			_, _, err := c.RequestPasswordReset(tt.args.ctx, tt.args.userID)
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_requestPasswordReset
		})
	}
}

func TestCommands_RequestPasswordResetReturnCode(t *testing.T) {
	type fields struct {
		checkPermission domain.PermissionCheck
		eventstore      func(t *testing.T) *eventstore.Eventstore
		userEncryption  crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-SAFdda", "Errors.User.IDMissing"),
		},
		{
			name: "user not existing",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-SAF4f", "Errors.User.NotFound"),
		},
		{
			name: "user not initialized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "keyID", Crypted: []byte("code")}, 10*time.Second),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sfe4g", "Errors.User.NotInitialised"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				checkPermission: tt.fields.checkPermission,
				eventstore:      tt.fields.eventstore(t),
				userEncryption:  tt.fields.userEncryption,
			}
			_, _, err := c.RequestPasswordResetReturnCode(tt.args.ctx, tt.args.userID)
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_requestPasswordReset
		})
	}
}

func TestCommands_RequestPasswordResetURLTemplate(t *testing.T) {
	type fields struct {
		checkPermission domain.PermissionCheck
		eventstore      func(t *testing.T) *eventstore.Eventstore
		userEncryption  crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx              context.Context
		userID           string
		urlTmpl          string
		notificationType domain.NotificationType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "invalid template",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				userID:  "user1",
				urlTmpl: "{{",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},

		{
			name: "missing userID",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-SAFdda", "Errors.User.IDMissing"),
		},
		{
			name: "user not existing",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-SAF4f", "Errors.User.NotFound"),
		},
		{
			name: "user not initialized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "keyID", Crypted: []byte("code")}, 10*time.Second),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sfe4g", "Errors.User.NotInitialised"),
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				checkPermission: tt.fields.checkPermission,
				eventstore:      tt.fields.eventstore(t),
				userEncryption:  tt.fields.userEncryption,
			}
			_, _, err := c.RequestPasswordResetURLTemplate(tt.args.ctx, tt.args.userID, tt.args.urlTmpl, tt.args.notificationType)
			require.ErrorIs(t, err, tt.wantErr)
			// successful cases are tested in TestCommands_requestPasswordReset
		})
	}
}

func TestCommands_requestPasswordReset(t *testing.T) {
	type fields struct {
		checkPermission domain.PermissionCheck
		eventstore      func(t *testing.T) *eventstore.Eventstore
		userEncryption  crypto.EncryptionAlgorithm
		newCode         cryptoCodeFunc
	}
	type args struct {
		ctx              context.Context
		userID           string
		returnCode       bool
		urlTmpl          string
		notificationType domain.NotificationType
	}
	type res struct {
		details *domain.ObjectDetails
		code    *string
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-SAFdda", "Errors.User.IDMissing"),
			},
		},
		{
			name: "user not existing",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			res: res{
				err: zerrors.ThrowNotFound(nil, "COMMAND-SAF4f", "Errors.User.NotFound"),
			},
		},
		{
			name: "user not initialized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{CryptoType: crypto.TypeEncryption, Algorithm: "enc", KeyID: "keyID", Crypted: []byte("code")}, 10*time.Second),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sfe4g", "Errors.User.NotInitialised"),
			},
		},
		{
			name: "missing permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			res: res{
				err: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			name: "code generated",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
					),
					expectPush(
						user.NewHumanPasswordCodeAddedEventV2(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("code"),
							},
							10*time.Minute,
							domain.NotificationTypeEmail,
							"",
							false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockCode("code", 10*time.Minute),
			},
			args: args{
				ctx:    context.Background(),
				userID: "userID",
			},
			res: res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				code: nil,
			},
		},
		{
			name: "code generated template",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
					),
					expectPush(
						user.NewHumanPasswordCodeAddedEventV2(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
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
				newCode:         mockCode("code", 10*time.Minute),
			},
			args: args{
				ctx:     context.Background(),
				userID:  "userID",
				urlTmpl: "https://example.com/password/changey?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
			},
			res: res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				code: nil,
			},
		},
		{
			name: "code generated template sms",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
					),
					expectPush(
						user.NewHumanPasswordCodeAddedEventV2(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("code"),
							},
							10*time.Minute,
							domain.NotificationTypeSms,
							"https://example.com/password/changey?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
							false,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
				newCode:         mockCode("code", 10*time.Minute),
			},
			args: args{
				ctx:              context.Background(),
				userID:           "userID",
				urlTmpl:          "https://example.com/password/changey?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}",
				notificationType: domain.NotificationTypeSms,
			},
			res: res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				code: nil,
			},
		},
		{
			name: "code generated returned",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstname", "lastname", "nickname", "displayname",
								language.English, domain.GenderUnspecified, "email", false),
						),
					),
					expectPush(
						user.NewHumanPasswordCodeAddedEventV2(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
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
				newCode:         mockCode("code", 10*time.Minute),
			},
			args: args{
				ctx:        context.Background(),
				userID:     "userID",
				returnCode: true,
			},
			res: res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				code: gu.Ptr("code"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				checkPermission: tt.fields.checkPermission,
				eventstore:      tt.fields.eventstore(t),
				userEncryption:  tt.fields.userEncryption,
				newCode:         tt.fields.newCode,
			}
			got, gotPlainCode, err := c.requestPasswordReset(tt.args.ctx, tt.args.userID, tt.args.returnCode, tt.args.urlTmpl, tt.args.notificationType)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.details, got)
			assert.Equal(t, tt.res.code, gotPlainCode)
		})
	}
}
