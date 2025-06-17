package command

import (
	"io"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
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

func TestCommands_AddUserTOTP(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	userAgg2 := &user.NewAggregate("user2", "org1").Aggregate

	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		userID        string
		resourceowner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "other user, permission error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								userAgg,
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
				userID:        "foo",
				resourceowner: "org1",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "create otp error",
			args: args{
				userID:        "user1",
				resourceowner: "org1",
			},
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SqyJz", "Errors.User.NotFound"),
		},
		{
			name: "push error",
			args: args{
				userID:        "user1",
				resourceowner: "org1",
			},
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								userAgg,
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
							org.NewOrgAddedEvent(ctx,
								userAgg,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(ctx,
								userAgg,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(),
					expectRandomPushFailed(io.ErrClosedPipe, []eventstore.Command{
						user.NewHumanOTPAddedEvent(ctx, userAgg, nil),
					}),
				),
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			args: args{
				userID:        "user1",
				resourceowner: "org1",
			},
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								userAgg,
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
							org.NewOrgAddedEvent(ctx,
								userAgg,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(ctx,
								userAgg,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(),
					expectRandomPush([]eventstore.Command{
						user.NewHumanOTPAddedEvent(ctx, userAgg, nil),
					}),
				),
			},
			want: true,
		},
		{
			name: "success, other user",
			args: args{
				userID:        "user2",
				resourceowner: "org1",
			},
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								userAgg2,
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
							org.NewOrgAddedEvent(ctx,
								userAgg2,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(ctx,
								userAgg2,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(),
					expectRandomPush([]eventstore.Command{
						user.NewHumanOTPAddedEvent(ctx, userAgg2, nil),
					}),
				),
				checkPermission: newMockPermissionCheckAllowed(),
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
						Issuer:    "zitadel.com",
						CryptoMFA: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
					},
				},
			}
			got, err := c.AddUserTOTP(ctx, tt.args.userID, tt.args.resourceowner)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want {
				require.NotNil(t, got)
				assert.Equal(t, "org1", got.ResourceOwner)
				assert.NotEmpty(t, got.Secret)
				assert.NotEmpty(t, got.URI)
			}
		})
	}
}

func TestCommands_CheckUserTOTP(t *testing.T) {
	ctx := authz.NewMockContext("", "org1", "user1")

	cryptoAlg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	key, err := domain.NewTOTPKey("example.com", "user1")
	require.NoError(t, err)
	secret, err := crypto.Encrypt([]byte(key.Secret()), cryptoAlg)
	require.NoError(t, err)

	userAgg := &user.NewAggregate("user1", "org1").Aggregate

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
			name: "other user, no permission, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				userID: "foo",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "user id, with permission, success",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, &user.NewAggregate("foo", "org1").Aggregate, secret),
						),
					),
					expectPush(
						user.NewHumanOTPVerifiedEvent(ctx, &user.NewAggregate("foo", "org1").Aggregate, ""),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				userID: "foo",
				code:   code,
			},
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
						user.NewHumanOTPVerifiedEvent(ctx, userAgg, ""),
					),
				),
			},
			args: args{
				resourceOwner: "org1",
				code:          code,
				userID:        "user1",
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
			got, err := c.CheckUserTOTP(ctx, tt.args.userID, tt.args.code, tt.args.resourceOwner)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want {
				require.NotNil(t, got)
				assert.Equal(t, "org1", got.ResourceOwner)
			}
		})
	}
}
