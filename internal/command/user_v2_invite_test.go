package command

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateInviteCode(t *testing.T) {
	t.Parallel()
	type fields struct {
		checkPermission             domain.PermissionCheck
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		eventstore                  func(*testing.T) *eventstore.Eventstore
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		ctx    context.Context
		invite *CreateUserInvite
	}
	type want struct {
		details    *domain.ObjectDetails
		returnCode *string
		err        error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			"user id missing",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID: "",
				},
			},
			want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-4jio3", "Errors.User.UserIDMissing"),
			},
		},
		{
			"missing permission",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID: "userID",
				},
			},
			want{
				err: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			"user does not exist",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID: "unknown",
				},
			},
			want{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Wgvn4", "Errors.User.NotFound"),
			},
		},
		{
			"create ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID: "userID",
				},
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
				returnCode: nil,
			},
		},
		{
			"return ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								true,
								"",
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID:     "userID",
					ReturnCode: true,
				},
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
				returnCode: gu.Ptr("code"),
			},
		},
		{
			"return ok, with same user requests code",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(authz.SetCtxData(context.Background(), authz.CtxData{UserID: "userID"}),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								true,
								"",
								"",
							),
						),
					),
				),
				// we do not run checkPermission() because the same user is requesting the code as the user to which the code is intended for
				checkPermission:             nil,
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx: authz.SetCtxData(context.Background(), authz.CtxData{UserID: "userID"}),
				invite: &CreateUserInvite{
					UserID:     "userID",
					ReturnCode: true,
				},
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
				returnCode: gu.Ptr("code"),
			},
		},
		{
			"with template and application name ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"https://example.com/invite?userID={{.UserID}}",
								false,
								"applicationName",
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID:          "userID",
					URLTemplate:     "https://example.com/invite?userID={{.UserID}}",
					ReturnCode:      false,
					ApplicationName: "applicationName",
				},
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
				returnCode: nil,
			},
		},
		{
			"create ok after three verification failures",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						// first invite code generated and returned
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code1"),
								},
								time.Hour,
								"",
								true,
								"",
								"",
							),
						),
						// simulate three failed verification attempts
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code2"),
								},
								time.Hour,
								"",
								false,
								"",
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code2", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID: "userID",
				},
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
				returnCode: nil,
			},
		},
		{
			"return ok after three verification failures",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						// first invite code generated and returned
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code1"),
								},
								time.Hour,
								"",
								true,
								"",
								"",
							),
						),
						// simulate three failed verification attempts
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code2"),
								},
								time.Hour,
								"",
								true,
								"",
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code2", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID:     "userID",
					ReturnCode: true,
				},
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
				returnCode: gu.Ptr("code2"),
			},
		},
		{
			"create ok after verification fails due to invite code expiration",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						// first invite code generated and returned
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code1"),
								},
								-5*time.Minute, // expired code
								"",
								true,
								"",
								"",
							),
						),
						// simulate a failed verification attempt due to expiry
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code2"),
								},
								time.Hour,
								"",
								false,
								"",
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code2", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID: "userID",
				},
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
				returnCode: nil,
			},
		},
		{
			"return ok after verification fails due to invite code expiration",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						// first invite code generated and returned
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code1"),
								},
								-5*time.Minute, // expired code
								"",
								true,
								"",
								"",
							),
						),
						// simulate a failed verification attempt due to expiry
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code2"),
								},
								time.Hour,
								"",
								true,
								"",
								"",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code2", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx: context.Background(),
				invite: &CreateUserInvite{
					UserID:     "userID",
					ReturnCode: true,
				},
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
				returnCode: gu.Ptr("code2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{
				checkPermission:             tt.fields.checkPermission,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				eventstore:                  tt.fields.eventstore(t),
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
			}
			gotDetails, gotReturnCode, err := c.CreateInviteCode(tt.args.ctx, tt.args.invite)

			require.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.details, gotDetails)
			assert.Equal(t, tt.want.returnCode, gotReturnCode)
		})
	}
}

func TestCommands_ResendInviteCode(t *testing.T) {
	t.Parallel()
	type fields struct {
		checkPermission             domain.PermissionCheck
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		eventstore                  func(*testing.T) *eventstore.Eventstore
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		ctx           context.Context
		userID        string
		orgID         string
		authRequestID string
	}
	type want struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			"missing user id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:    context.Background(),
				userID: "",
			},
			want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-4jio3", "Errors.User.UserIDMissing"),
			},
		},
		{
			"missing permission",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:    context.Background(),
				userID: "userID",
			},
			want{
				err: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			"user does not exist",
			fields{
				eventstore: expectEventstore(
					// The write model doesn't query any events
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    context.Background(),
				userID: "unknown",
			},
			want{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Wgvn4", "Errors.User.NotFound"),
			},
		},
		{
			"no previous code",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    context.Background(),
				userID: "userID",
			},
			want{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Wr3gq", "Errors.User.Code.NotFound"),
			},
		},
		{
			"previous code returned",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								true,
								"",
								"",
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:    context.Background(),
				userID: "userID",
			},
			want{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Wr3gq", "Errors.User.Code.NotFound"),
			},
		},
		{
			"resend ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx:    context.Background(),
				userID: "userID",
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
			},
		},
		{
			"return ok, with same user requests code",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(authz.SetCtxData(context.Background(), authz.CtxData{UserID: "userID"}),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
				),
				// we do not run checkPermission() because the same user is requesting the code as the user to which the code is intended for
				checkPermission:             nil,
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				// ctx:    context.Background(),
				ctx:    authz.SetCtxData(context.Background(), authz.CtxData{UserID: "userID"}),
				userID: "userID",
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
			},
		},
		{
			"resend with new auth requestID ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID2",
							),
						),
					),
				),
				checkPermission:             newMockPermissionCheckAllowed(),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("code", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args{
				ctx:           context.Background(),
				userID:        "userID",
				authRequestID: "authRequestID2",
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{
				checkPermission:             tt.fields.checkPermission,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				eventstore:                  tt.fields.eventstore(t),
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
			}
			details, err := c.ResendInviteCode(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.authRequestID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.details, details)
		})
	}
}

func TestCommands_InviteCodeSent(t *testing.T) {
	t.Parallel()
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		userID string
		orgID  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			"missing user id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:    context.Background(),
				userID: "",
			},
			zerrors.ThrowInvalidArgument(nil, "COMMAND-Sgf31", "Errors.User.UserIDMissing"),
		},
		{
			"user does not exist",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx:    context.Background(),
				userID: "unknown",
			},
			zerrors.ThrowPreconditionFailed(nil, "COMMAND-HN34a", "Errors.User.NotFound"),
		},
		{
			"code does not exist",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
					),
				),
			},
			args{
				ctx:    context.Background(),
				userID: "unknown",
			},
			zerrors.ThrowPreconditionFailed(nil, "COMMAND-Wr3gq", "Errors.User.Code.NotFound"),
		},
		{
			"sent ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusher(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCodeSentEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
					),
				),
			},
			args{
				ctx:    context.Background(),
				userID: "userID",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := c.InviteCodeSent(tt.args.ctx, tt.args.userID, tt.args.orgID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_VerifyInviteCode(t *testing.T) {
	t.Parallel()
	type fields struct {
		eventstore     func(*testing.T) *eventstore.Eventstore
		userEncryption crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx    context.Context
		userID string
		code   string
	}
	type want struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			"code ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCheckSucceededEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:    context.Background(),
				userID: "userID",
				code:   "code",
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
			},
		},
		// all other cases are tested in TestCommands_VerifyInviteCodeSetPassword
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{
				eventstore:     tt.fields.eventstore(t),
				userEncryption: tt.fields.userEncryption,
			}
			gotDetails, err := c.VerifyInviteCode(tt.args.ctx, tt.args.userID, tt.args.code)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.details, gotDetails)
		})
	}
}

func TestCommands_VerifyInviteCodeSetPassword(t *testing.T) {
	t.Parallel()
	type fields struct {
		eventstore         func(*testing.T) *eventstore.Eventstore
		userEncryption     crypto.EncryptionAlgorithm
		userPasswordHasher *crypto.Hasher
	}
	type args struct {
		ctx         context.Context
		userID      string
		code        string
		password    string
		userAgentID string
	}
	type want struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			"missing user id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:    context.Background(),
				userID: "",
			},
			want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-Gk3f2", "Errors.User.UserIDMissing"),
			},
		},
		{
			"user does not exist",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx:    context.Background(),
				userID: "unknown",
			},
			want{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-F5g2h", "Errors.User.NotFound"),
			},
		},
		{
			"invalid code",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCheckFailedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:    context.Background(),
				userID: "userID",
				code:   "invalid",
			},
			want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-Wgn4q", "Errors.User.Code.Invalid"),
			},
		},
		{
			"code ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCheckSucceededEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:    context.Background(),
				userID: "userID",
				code:   "code",
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
			},
		},
		{
			"code ok, with password and user agent",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								6,
								true,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						eventFromEventPusher(
							user.NewHumanInviteCheckSucceededEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanPasswordChangedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"$plain$x$Password1!",
								false,
								"userAgentID",
							),
						),
					),
				),
				userEncryption:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx:         context.Background(),
				userID:      "userID",
				code:        "code",
				password:    "Password1!",
				userAgentID: "userAgentID",
			},
			want{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
					ID:            "userID",
				},
			},
		},
		{
			"code ok, with non compliant password",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName",
								"lastName",
								"nickName",
								"displayName",
								language.Afrikaans,
								domain.GenderUnspecified,
								"email",
								false,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanInviteCodeAddedEvent(context.Background(),
								&user.NewAggregate("userID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								"",
								false,
								"",
								"authRequestID",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								6,
								true,
								true,
								true,
								true,
							),
						),
					),
				),
				userEncryption:     crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args{
				ctx:         context.Background(),
				userID:      "userID",
				code:        "code",
				password:    "pw",
				userAgentID: "userAgentID",
			},
			want{
				err: zerrors.ThrowInvalidArgument(nil, "DOMAIN-HuJf6", "Errors.User.PasswordComplexityPolicy.MinLength"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := &Commands{
				eventstore:         tt.fields.eventstore(t),
				userEncryption:     tt.fields.userEncryption,
				userPasswordHasher: tt.fields.userPasswordHasher,
			}
			gotDetails, err := c.VerifyInviteCodeSetPassword(tt.args.ctx, tt.args.userID, tt.args.code, tt.args.password, tt.args.userAgentID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.details, gotDetails)
		})
	}
}
