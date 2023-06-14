package command

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"golang.org/x/text/language"
)

func TestCommands_AddUserOTP(t *testing.T) {
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
							user.NewHumanAddedEvent(context.Background(),
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
							org.NewOrgAddedEvent(context.Background(),
								userAgg,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
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
							user.NewHumanAddedEvent(context.Background(),
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
							org.NewOrgAddedEvent(context.Background(),
								userAgg,
								"org",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
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
			got, err := c.AddUserOTP(ctx, tt.args.userID, tt.args.resourceowner)
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
