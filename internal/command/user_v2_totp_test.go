package command

import (
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommands_AddUserTOTP(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	type fields struct {
		eventstore *eventstore.Eventstore
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
			name: "wrong user",
			args: args{
				userID:        "foo",
				resourceowner: "org1",
			},
			wantErr: caos_errs.ThrowUnauthenticated(nil, "AUTH-Bohd2", "Errors.User.UserIDWrong"),
		},
		{
			name: "create otp error",
			args: args{
				userID:        "user1",
				resourceowner: "org1",
			},
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-MM9fs", "Errors.User.NotFound"),
		},
		{
			name: "push error",
			args: args{
				userID:        "user1",
				resourceowner: "org1",
			},
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectRandomPushFailed(io.ErrClosedPipe, []*repository.Event{eventFromEventPusher(
						user.NewHumanOTPAddedEvent(ctx, userAgg, nil),
					)}),
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
				eventstore: eventstoreExpect(
					t,
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
					expectRandomPush([]*repository.Event{eventFromEventPusher(
						user.NewHumanOTPAddedEvent(ctx, userAgg, nil),
					)}),
				),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
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
	key, secret, err := domain.NewOTPKey("example.com", "user1", cryptoAlg)
	require.NoError(t, err)
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	code, err := totp.GenerateCode(key.Secret(), time.Now())
	require.NoError(t, err)

	type fields struct {
		eventstore *eventstore.Eventstore
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
			name: "wrong user id",
			args: args{
				userID: "foo",
			},
			wantErr: caos_errs.ThrowUnauthenticated(nil, "AUTH-Bohd2", "Errors.User.UserIDWrong"),
		},
		{
			name: "success",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
					),
					expectPush([]*repository.Event{eventFromEventPusher(
						user.NewHumanOTPVerifiedEvent(ctx, userAgg, ""),
					)}),
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
				eventstore: tt.fields.eventstore,
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
