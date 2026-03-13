package command

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_humanVerifyPasswordlessInitCode(t *testing.T) {
	ctx := http_util.WithRequestedHost(context.Background(), "example.com")
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	es := expectEventstore(
		expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
	)(t)
	code, err := newEncryptedCode(ctx, es.Filter, domain.SecretGeneratorTypePasswordlessInitCode, alg) //nolint:staticcheck
	require.NoError(t, err)
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		userID        string
		resourceOwner string
		codeID        string
		code          string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
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
				codeID:        "123",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "code verification error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeAddedEvent(context.Background(),
								userAgg, "123", code.Crypted, time.Minute,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "123"),
						),
					),
					expectPush(
						user.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, "123"),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				codeID:        "123",
				code:          "wrong",
			},
			wantErr: zerrors.ThrowInvalidArgument(err, "COMMAND-Dhz8i", "Errors.User.Code.Invalid"),
		},
		{
			name: "success",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeAddedEvent(context.Background(),
								userAgg, "123", code.Crypted, time.Minute,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "123"),
						),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				codeID:        "123",
				code:          code.Plain,
			},
		},
		{
			name: "expired error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithCreationDate(
							user.NewHumanPasswordlessInitCodeAddedEvent(context.Background(),
								userAgg, "123", code.Crypted, time.Minute,
							),
							time.Now().Add(-2*time.Minute),
						),
						eventFromEventPusherWithCreationDate(
							user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "123"),
							time.Now().Add(-2*time.Minute),
						),
					),
					expectPush(
						user.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, "123"),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				codeID:        "123",
				code:          code.Plain,
			},
			wantErr: zerrors.ThrowInvalidArgument(err, "COMMAND-Dhz8i", "Errors.User.Code.Invalid"),
		},
		{
			// https://github.com/zitadel/zitadel/security/advisories/GHSA-2x66-r53r-9r86
			name: "expired, fail, check again",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithCreationDate(
							user.NewHumanPasswordlessInitCodeAddedEvent(context.Background(),
								userAgg, "123", code.Crypted, time.Minute,
							),
							time.Now().Add(-2*time.Minute),
						),
						eventFromEventPusherWithCreationDate(
							user.NewHumanPasswordlessInitCodeSentEvent(ctx, userAgg, "123"),
							time.Now().Add(-2*time.Minute),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, "123"),
						),
					),
					expectPush(
						user.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, "123"),
					),
				),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				codeID:        "123",
				code:          code.Plain,
			},
			wantErr: zerrors.ThrowInvalidArgument(err, "COMMAND-Dhz8i", "Errors.User.Code.Invalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := c.humanVerifyPasswordlessInitCode(ctx, tt.args.userID, tt.args.resourceOwner, tt.args.codeID, tt.args.code, alg)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
