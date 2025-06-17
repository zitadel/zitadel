package command

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/senders/mock"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddHumanTOTP(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		permissionCheck domain.PermissionCheck
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:    authz.NewMockContext("instanceID", "org1", "user1"),
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
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:    authz.NewMockContext("instanceID", "org1", "user1"),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "other user, no permission error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user2", "org1").Aggregate,
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
				permissionCheck: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:    authz.NewMockContext("instanceID", "org1", "user1"),
				orgID:  "org1",
				userID: "user2",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "org not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
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
					expectFilter(),
				),
			},
			args: args{
				ctx:    authz.NewMockContext("instanceID", "org1", "user1"),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "org iam policy not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
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
					expectFilter(
						org.NewOrgAddedEvent(context.Background(),
							&user.NewAggregate("org1", "org1").Aggregate,
							"org",
						),
					),
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:    authz.NewMockContext("instanceID", "org1", "user1"),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "otp already exists, already exists error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
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
					expectFilter(
						org.NewOrgAddedEvent(context.Background(),
							&user.NewAggregate("org1", "org1").Aggregate,
							"org",
						),
					),
					expectFilter(
						org.NewDomainPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							true,
							true,
							true,
						),
					),
					expectFilter(
						user.NewHumanOTPAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("a"),
							},
						),
						user.NewHumanOTPVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"agent1",
						),
					),
				),
			},
			args: args{
				ctx:    authz.NewMockContext("instanceID", "org1", "user1"),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.permissionCheck,
			}
			got, err := r.AddHumanTOTP(tt.args.ctx, tt.args.userID, tt.args.orgID)
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

func TestCommands_createHumanTOTP(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.NewMockContext("instanceID", "org1", "user1"),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SqyJz", "Errors.User.NotFound"),
		},
		{
			name: "other user, no permission error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user2", "org1").Aggregate,
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
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           authz.NewMockContext("instanceID", "org1", "user1"),
				resourceOwner: "org1",
				userID:        "user2",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "org not existing, not found error",
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
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.NewMockContext("instanceID", "org1", "user1"),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-55M9f", "Errors.Org.NotFound"),
		},
		{
			name: "org iam policy not existing, not found error",
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("org1", "org1").Aggregate,
								"org",
							),
						),
					),
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.NewMockContext("instanceID", "org1", "user1"),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-8ugTs", "Errors.Org.DomainPolicy.NotFound"),
		},
		{
			name: "otp already exists, already exists error",
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("org1", "org1").Aggregate,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								}),
						),
						eventFromEventPusher(
							user.NewHumanOTPVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"agent1")),
					),
				),
			},
			args: args{
				ctx:           authz.NewMockContext("instanceID", "org1", "user1"),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: zerrors.ThrowAlreadyExists(nil, "COMMAND-do9se", "Errors.User.MFA.OTP.AlreadyReady"),
		},
		{
			name: "issuer not in context",
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("org1", "org1").Aggregate,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           authz.NewMockContext("instanceID", "org1", "user1"),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: zerrors.ThrowInternal(nil, "TOTP-ieY3o", "Errors.Internal"),
		},
		{
			name: "success",
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("org1", "org1").Aggregate,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:           http_util.WithRequestedHost(authz.NewMockContext("instanceID", "org1", "user1"), "zitadel.com"),
				resourceOwner: "org1",
				userID:        "user1",
			},
			want: true,
		},
		{
			name: "success, other user",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user2", "org1").Aggregate,
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
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("org1", "org1").Aggregate,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           http_util.WithRequestedHost(authz.NewMockContext("instanceID", "org1", "user1"), "zitadel.com"),
				resourceOwner: "org1",
				userID:        "user2",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
				multifactors: domain.MultifactorConfigs{
					OTP: domain.OTPConfig{
						CryptoMFA: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
					},
				},
			}
			got, err := c.createHumanTOTP(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want {
				require.NotNil(t, got)
				assert.NotNil(t, got.wm)
				assert.NotNil(t, got.userAgg)
				require.NotNil(t, got.key)
				assert.NotEmpty(t, got.key.URL())
				assert.NotEmpty(t, got.key.Secret())
				assert.Len(t, got.cmds, 1)
			}
		})
	}
}

func TestCommands_HumanCheckMFATOTPSetup(t *testing.T) {
	ctx := authz.NewMockContext("", "org1", "user1")

	cryptoAlg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	key, err := domain.NewTOTPKey("example.com", "user1")
	require.NoError(t, err)
	secret, err := crypto.Encrypt([]byte(key.Secret()), cryptoAlg)
	require.NoError(t, err)
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	userAgg2 := &user.NewAggregate("user2", "org1").Aggregate

	code, err := totp.GenerateCode(key.Secret(), time.Now())
	require.NoError(t, err)

	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID        string
		code          string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "missing user id",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args:    args{},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing"),
		},
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "other user, no permission error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg2, secret),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				resourceOwner: "org1",
				userID:        "user2",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "otp not existing error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
						eventFromEventPusher(
							user.NewHumanOTPRemovedEvent(ctx, userAgg),
						),
					),
				),
			},
			args: args{
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotExisting"),
		},
		{
			name: "otp already ready error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
						eventFromEventPusher(
							user.NewHumanOTPVerifiedEvent(context.Background(),
								userAgg,
								"agent1",
							),
						),
					),
				),
			},
			args: args{
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-qx4ls", "Errors.Users.MFA.OTP.AlreadyReady"),
		},
		{
			name: "wrong code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
					),
				),
			},
			args: args{
				resourceOwner: "org1",
				code:          "wrong",
				userID:        "user1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
					),
					expectPushFailed(io.ErrClosedPipe,
						user.NewHumanOTPVerifiedEvent(ctx,
							userAgg,
							"agent1",
						),
					),
				),
			},
			args: args{
				resourceOwner: "org1",
				code:          code,
				userID:        "user1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
					),
					expectPush(
						user.NewHumanOTPVerifiedEvent(ctx,
							userAgg,
							"agent1",
						),
					),
				),
			},
			args: args{
				resourceOwner: "org1",
				code:          code,
				userID:        "user1",
			},
		},
		{
			name: "success, other user",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg2, secret),
						),
					),
					expectPush(
						user.NewHumanOTPVerifiedEvent(ctx,
							userAgg2,
							"agent1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				resourceOwner: "org1",
				code:          code,
				userID:        "user2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
				multifactors: domain.MultifactorConfigs{
					OTP: domain.OTPConfig{
						CryptoMFA: cryptoAlg,
					},
				},
			}
			got, err := c.HumanCheckMFATOTPSetup(ctx, tt.args.userID, tt.args.code, "agent1", tt.args.resourceOwner)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want {
				require.NotNil(t, got)
				assert.Equal(t, "org1", got.ResourceOwner)
			}
		})
	}
}

func TestCommandSide_RemoveHumanTOTP(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
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
				eventstore: expectEventstore(),
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
			name: "otp not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
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
			name: "otp, no permission error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "otp remove, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil,
							),
						),
					),
					expectPush(
						user.NewHumanOTPRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
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
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.HumanRemoveTOTP(tt.args.ctx, tt.args.userID, tt.args.orgID)
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

func TestCommandSide_AddHumanOTPSMS(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-QSF2s", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "wrong user, permission denied error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user2",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			name: "otp sms already exists, already exists error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowAlreadyExists(nil, "COMMAND-Ad3g2", "Errors.User.MFA.OTP.AlreadyReady"),
			},
		},
		{
			name: "phone not verified, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Q54j2", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "phone removed, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"+4179654321",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneRemovedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Q54j2", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "successful add",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"+4179654321",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
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
			name: "successful add, other user",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
								"+4179654321",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSAddedEvent(ctx,
							&user.NewAggregate("user2", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user2",
				resourceOwner: "org1",
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
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.AddHumanOTPSMS(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
			assertObjectDetails(t, tt.res.want, got)
		})
	}
}

func TestCommandSide_AddHumanOTPSMSWithCheckSucceeded(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
			authRequest   *domain.AuthRequest
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "successful add",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"+4179654321",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
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
			name: "successful add with auth request",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanPhoneChangedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"+4179654321",
							),
						),
						eventFromEventPusher(
							user.NewHumanPhoneVerifiedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
						user.NewHumanOTPSMSCheckSucceededEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.AddHumanOTPSMSWithCheckSucceeded(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.authRequest)
			assert.ErrorIs(t, err, tt.res.err)
			assertObjectDetails(t, tt.res.want, got)
		})
	}
}

func TestCommandSide_RemoveHumanOTPSMS(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-S3br2", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "other user not permission, permission denied error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "other",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			name: "otp sms not added, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowNotFound(nil, "COMMAND-Sr3h3", "Errors.User.MFA.OTP.NotExisting"),
			},
		},
		{
			name: "successful remove",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSRemovedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
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
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.RemoveHumanOTPSMS(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
			assertObjectDetails(t, tt.res.want, got)
		})
	}
}

func TestCommandSide_HumanSendOTPSMS(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
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
		eventstore                  func(*testing.T) *eventstore.Eventstore
		defaultSecretGenerators     *SecretGenerators
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
			authRequest   *domain.AuthRequest
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				eventstore:              expectEventstore(),
				defaultSecretGenerators: defaultGenerators,
			},
			args: args{
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-S3SF1", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp sms not added, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				defaultSecretGenerators: defaultGenerators,
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SFD52", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "successful add",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
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
						user.NewHumanOTPSMSCodeAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
							time.Hour,
							nil,
							"",
						),
					),
				),
				defaultSecretGenerators:     defaultGenerators,
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
			},
			args: args{
				ctx:           ctx,
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
			name: "successful add (external code)",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
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
								"verifyServiceSID",
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
						user.NewHumanOTPSMSCodeAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							nil,
							0,
							nil,
							"id",
						),
					),
				),
				defaultSecretGenerators:     defaultGenerators,
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
			},
			args: args{
				ctx:           ctx,
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
			name: "successful add (without secret config)",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
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
						user.NewHumanOTPSMSCodeAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
							time.Hour,
							nil,
							"",
						),
					),
				),
				defaultSecretGenerators:     defaultGenerators,
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
			},
			args: args{
				ctx:           ctx,
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
			name: "successful add with auth request",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
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
						user.NewHumanOTPSMSCodeAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
							time.Hour,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
							"",
						),
					),
				),
				defaultSecretGenerators:     defaultGenerators,
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
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
				eventstore:                  tt.fields.eventstore(t),
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
			}
			err := r.HumanSendOTPSMS(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.authRequest)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommandSide_HumanOTPSMSCodeSent(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
			generatorInfo *senders.CodeGeneratorInfo
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-AE2h2", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp sms not added, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SD3gh", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "successful add",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSCodeSentEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&senders.CodeGeneratorInfo{},
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				generatorInfo: &senders.CodeGeneratorInfo{},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "successful add (external code)",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSCodeSentEvent(ctx,
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
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				generatorInfo: &senders.CodeGeneratorInfo{
					ID:             "generatorID",
					VerificationID: "verificationID",
				},
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
				eventstore: tt.fields.eventstore(t),
			}
			err := r.HumanOTPSMSCodeSent(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.generatorInfo)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommandSide_HumanCheckOTPSMS(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore        func(*testing.T) *eventstore.Eventstore
		userEncryption    crypto.EncryptionAlgorithm
		phoneCodeVerifier func(ctx context.Context, id string) (senders.CodeGenerator, error)
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			code          string
			resourceOwner string
			authRequest   *domain.AuthRequest
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				ctx:           ctx,
				userID:        "",
				code:          "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-S453v", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "code missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-SJl2g", "Errors.User.Code.Empty"),
			},
		},
		{
			name: "otp sms not added, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-d2r52", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "otp sms code not added, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-S34gh", "Errors.User.Code.NotFound"),
			},
		},
		{
			name: "invalid code, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPSMSCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("other-code"),
								},
								time.Hour,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
								"",
							),
						),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(ctx,
								&org.NewAggregate("orgID").Aggregate,
								3, 3, true,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSCheckFailedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "CODE-woT0xc", "Errors.User.Code.Invalid"),
			},
		},
		{
			name: "invalid code, max attempts reached, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPSMSCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("other-code"),
								},
								time.Hour,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
								"",
							),
						),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(ctx,
								&org.NewAggregate("orgID").Aggregate,
								1, 1, true,
							),
						),
					),
					expectPush(
						user.NewHumanOTPSMSCheckFailedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
						user.NewUserLockedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "CODE-woT0xc", "Errors.User.Code.Invalid"),
			},
		},
		{
			name: "code ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPSMSCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
								"",
							),
						),
					),
					expectFilter(), // recheck
					expectPush(
						user.NewHumanOTPSMSCheckSucceededEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "code ok, locked in the meantime",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPSMSCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
								"",
							),
						),
					),
					expectFilter( // recheck
						user.NewUserLockedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-S6h4R", "Errors.User.Locked"),
			},
		},
		{
			name: "code ok (external)",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPSMSCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								nil,
								0,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
								"id",
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPSMSCodeSentEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&senders.CodeGeneratorInfo{
									ID:             "id",
									VerificationID: "verificationID",
								},
							),
						),
					),
					expectFilter(), // recheck
					expectPush(
						user.NewHumanOTPSMSCheckSucceededEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				phoneCodeVerifier: func(ctx context.Context, id string) (senders.CodeGenerator, error) {
					sender := mock.NewMockCodeGenerator(gomock.NewController(t))
					sender.EXPECT().VerifyCode("verificationID", "code").Return(nil)
					return sender, nil
				},
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
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
				eventstore:        tt.fields.eventstore(t),
				userEncryption:    tt.fields.userEncryption,
				phoneCodeVerifier: tt.fields.phoneCodeVerifier,
			}
			err := r.HumanCheckOTPSMS(tt.args.ctx, tt.args.userID, tt.args.code, tt.args.resourceOwner, tt.args.authRequest)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommandSide_AddHumanOTPEmail(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-Sg1hz", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "wrong user, permission denied error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user2",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			name: "otp email already exists, already exists error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowAlreadyExists(nil, "COMMAND-MKL2s", "Errors.User.MFA.OTP.AlreadyReady"),
			},
		},
		{
			name: "email not verified, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-KLJ2d", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "successful add",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanEmailChangedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"email@test.ch",
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
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
			name: "successful add, other user",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanEmailChangedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
								"email@test.ch",
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailAddedEvent(ctx,
							&user.NewAggregate("user2", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user2",
				resourceOwner: "org1",
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
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.AddHumanOTPEmail(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
			assertObjectDetails(t, tt.res.want, got)
		})
	}
}

func TestCommandSide_AddHumanOTPEmailWithCheckSucceeded(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
			authRequest   *domain.AuthRequest
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "successful add",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanEmailChangedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"email@test.ch",
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
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
			name: "successful add with auth request",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanEmailChangedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"email@test.ch",
							),
						),
						eventFromEventPusher(
							user.NewHumanEmailVerifiedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
						user.NewHumanOTPEmailCheckSucceededEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.AddHumanOTPEmailWithCheckSucceeded(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.authRequest)
			assert.ErrorIs(t, err, tt.res.err)
			assertObjectDetails(t, tt.res.want, got)
		})
	}
}

func TestCommandSide_RemoveHumanOTPEmail(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-S2h11", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "other user not permission, permission denied error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "other",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			name: "otp email not added, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowNotFound(nil, "COMMAND-b312D", "Errors.User.MFA.OTP.NotExisting"),
			},
		},
		{
			name: "successful remove",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailRemovedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
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
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := r.RemoveHumanOTPEmail(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
			assertObjectDetails(t, tt.res.want, got)
		})
	}
}

func TestCommandSide_HumanSendOTPEmail(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	defaultGenerators := &SecretGenerators{
		OTPEmail: &crypto.GeneratorConfig{
			Length:              8,
			Expiry:              time.Hour,
			IncludeLowerLetters: true,
			IncludeUpperLetters: true,
			IncludeDigits:       true,
			IncludeSymbols:      true,
		},
	}
	type fields struct {
		eventstore                  func(*testing.T) *eventstore.Eventstore
		defaultSecretGenerators     *SecretGenerators
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
			authRequest   *domain.AuthRequest
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				eventstore:              expectEventstore(),
				defaultSecretGenerators: defaultGenerators,
			},
			args: args{
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-S3SF1", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp email not added, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				defaultSecretGenerators: defaultGenerators,
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SFD52", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "successful add",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailCodeAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
							time.Hour,
							nil,
						),
					),
				),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				ctx:           ctx,
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
			name: "successful add with auth request",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailCodeAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
							time.Hour,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("12345678", time.Hour),
				defaultSecretGenerators:     defaultGenerators,
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
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
				eventstore:                  tt.fields.eventstore(t),
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
			}
			err := r.HumanSendOTPEmail(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.authRequest)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommandSide_HumanOTPEmailCodeSent(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			resourceOwner string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-AE2h2", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp email not added, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SD3gh", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "successful add",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailCodeSentEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
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
				eventstore: tt.fields.eventstore(t),
			}
			err := r.HumanOTPEmailCodeSent(tt.args.ctx, tt.args.userID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommandSide_HumanCheckOTPEmail(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore     func(*testing.T) *eventstore.Eventstore
		userEncryption crypto.EncryptionAlgorithm
	}
	type (
		args struct {
			ctx           context.Context
			userID        string
			code          string
			resourceOwner string
			authRequest   *domain.AuthRequest
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  error
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
				ctx:           ctx,
				userID:        "",
				code:          "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-S453v", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "code missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-SJl2g", "Errors.User.Code.Empty"),
			},
		},
		{
			name: "otp email not added, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-d2r52", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "otp email code not added, precondition failed error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-S34gh", "Errors.User.Code.NotFound"),
			},
		},
		{
			name: "invalid code, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPEmailCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("other-code"),
								},
								time.Hour,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
							),
						),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(ctx,
								&org.NewAggregate("orgID").Aggregate,
								3, 3, true,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailCheckFailedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "CODE-woT0xc", "Errors.User.Code.Invalid"),
			},
		},
		{
			name: "invalid code, max attempts reached, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPEmailCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("other-code"),
								},
								time.Hour,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
							),
						),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(ctx,
								&org.NewAggregate("orgID").Aggregate,
								1, 1, true,
							),
						),
					),
					expectPush(
						user.NewHumanOTPEmailCheckFailedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
						user.NewUserLockedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "CODE-woT0xc", "Errors.User.Code.Invalid"),
			},
		},
		{
			name: "code ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPEmailCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
							),
						),
					),
					expectFilter(), // recheck
					expectPush(
						user.NewHumanOTPEmailCheckSucceededEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "code ok, locked in the meantime",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanOTPEmailCodeAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour,
								&user.AuthRequestInfo{
									ID:          "authRequestID",
									UserAgentID: "userAgentID",
									BrowserInfo: &user.BrowserInfo{
										UserAgent:      "user-agent",
										AcceptLanguage: "en",
										RemoteIP:       net.IP{192, 0, 2, 1},
									},
								},
							),
						),
					),
					expectFilter( // recheck
						user.NewUserLockedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
				userEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "code",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-S6h4R", "Errors.User.Locked"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore(t),
				userEncryption: tt.fields.userEncryption,
			}
			err := r.HumanCheckOTPEmail(tt.args.ctx, tt.args.userID, tt.args.code, tt.args.resourceOwner, tt.args.authRequest)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}
