package command

import (
	"context"
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

func TestCommandSide_AddHumanOTP(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "org not existing, not found error",
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
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "org iam policy not existing, not found error",
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "otp already exists, already exists error",
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddHumanOTP(tt.args.ctx, tt.args.userID, tt.args.orgID)
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

func TestCommands_createHumanOTP(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-MM9fs", "Errors.User.NotFound"),
		},
		{
			name: "org not existing, not found error",
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
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-55M9f", "Errors.Org.NotFound"),
		},
		{
			name: "org iam policy not existing, not found error",
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-8ugTs", "Errors.Org.DomainPolicy.NotFound"),
		},
		{
			name: "otp already exists, already exists error",
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: caos_errs.ThrowAlreadyExists(nil, "COMMAND-do9se", "Errors.User.MFA.OTP.AlreadyReady"),
		},
		{
			name: "issuer not in context",
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				userID:        "user1",
			},
			wantErr: caos_errs.ThrowInternal(nil, "TOTP-ieY3o", "Errors.Internal"),
		},
		{
			name: "success",
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
				ctx:           authz.WithRequestedDomain(context.Background(), "zitadel.com"),
				resourceOwner: "org1",
				userID:        "user1",
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

func TestCommands_HumanCheckMFAOTPSetup(t *testing.T) {
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
			name:    "missing user id",
			args:    args{},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing"),
		},
		{
			name: "filter error",
			fields: fields{
				eventstore: eventstoreExpect(t,
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
			name: "otp not existing error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			wantErr: caos_errs.ThrowNotFound(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotExisting"),
		},
		{
			name: "otp already ready error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-qx4ls", "Errors.Users.MFA.OTP.AlreadyReady"),
		},
		{
			name: "wrong code",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
			wantErr: caos_errs.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
		},
		{
			name: "push error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
					),
					expectPushFailed(io.ErrClosedPipe,
						[]*repository.Event{eventFromEventPusher(
							user.NewHumanOTPVerifiedEvent(ctx,
								userAgg,
								"agent1",
							),
						)},
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
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
					),
					expectPush([]*repository.Event{eventFromEventPusher(
						user.NewHumanOTPVerifiedEvent(ctx,
							userAgg,
							"agent1",
						),
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
			got, err := c.HumanCheckMFAOTPSetup(ctx, tt.args.userID, tt.args.code, "agent1", tt.args.resourceOwner)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want {
				require.NotNil(t, got)
				assert.Equal(t, "org1", got.ResourceOwner)
			}
		})
	}
}

func TestCommandSide_RemoveHumanOTP(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "otp not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "otp not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewHumanOTPRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
								),
							),
						},
					),
				),
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
				eventstore: tt.fields.eventstore,
			}
			got, err := r.HumanRemoveOTP(tt.args.ctx, tt.args.userID, tt.args.orgID)
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
