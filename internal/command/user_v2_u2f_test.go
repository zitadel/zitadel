package command

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	webauthn_helper "github.com/zitadel/zitadel/internal/webauthn"
)

func TestCommands_RegisterUserU2F(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	ctx = authz.WithRequestedDomain(ctx, "example.com")

	webauthnConfig := &webauthn_helper.Config{
		DisplayName:    "test",
		ExternalSecure: true,
	}
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
		rpID          string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.WebAuthNRegistrationDetails
		wantErr error
	}{
		{
			name: "wrong user",
			args: args{
				userID:        "foo",
				resourceOwner: "org1",
			},
			wantErr: caos_errs.ThrowPermissionDenied(nil, "AUTH-Bohd2", "Errors.User.UserIDWrong"),
		},
		{
			name: "get human passwordless error",
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
			name: "id generator error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(), // getHumanPasswordlessTokens
					expectFilter(eventFromEventPusher(
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
					)),
					expectFilter(eventFromEventPusher(
						org.NewOrgAddedEvent(ctx,
							&org.NewAggregate("org1").Aggregate,
							"org1",
						),
					)),
					expectFilter(eventFromEventPusher(
						org.NewDomainPolicyAddedEvent(ctx,
							&org.NewAggregate("org1").Aggregate,
							false, false, false,
						),
					)),
				),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:     tt.fields.eventstore,
				idGenerator:    tt.fields.idGenerator,
				webauthnConfig: webauthnConfig,
			}
			_, err := c.RegisterUserU2F(ctx, tt.args.userID, tt.args.resourceOwner, tt.args.rpID)
			require.ErrorIs(t, err, tt.wantErr)
			// successful case can't be tested due to random challenge.
		})
	}
}

func TestCommands_pushUserU2F(t *testing.T) {
	ctx := authz.WithRequestedDomain(context.Background(), "example.com")
	webauthnConfig := &webauthn_helper.Config{
		DisplayName:    "test",
		ExternalSecure: true,
	}
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

	prep := []expect{
		expectFilter(), // getHumanU2FTokens
		expectFilter(eventFromEventPusher(
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
		)),
		expectFilter(eventFromEventPusher(
			org.NewOrgAddedEvent(ctx,
				&org.NewAggregate("org1").Aggregate,
				"org1",
			),
		)),
		expectFilter(eventFromEventPusher(
			org.NewDomainPolicyAddedEvent(ctx,
				&org.NewAggregate("org1").Aggregate,
				false, false, false,
			),
		)),
		expectFilter(eventFromEventPusher(
			user.NewHumanWebAuthNAddedEvent(eventstore.NewBaseEventForPush(
				ctx, &org.NewAggregate("org1").Aggregate, user.HumanPasswordlessTokenAddedType,
			), "111", "challenge", "rpID"),
		)),
	}

	tests := []struct {
		name       string
		expectPush func(challenge string) expect
		wantErr    error
	}{
		{
			name: "push error",
			expectPush: func(challenge string) expect {
				return expectPushFailed(io.ErrClosedPipe, []*repository.Event{eventFromEventPusher(
					user.NewHumanU2FAddedEvent(ctx,
						userAgg, "123", challenge, "rpID",
					),
				)})
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			expectPush: func(challenge string) expect {
				return expectPush([]*repository.Event{eventFromEventPusher(
					user.NewHumanU2FAddedEvent(ctx,
						userAgg, "123", challenge, "rpID",
					),
				)})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:     eventstoreExpect(t, prep...),
				webauthnConfig: webauthnConfig,
				idGenerator:    id_mock.NewIDGeneratorExpectIDs(t, "123"),
			}
			wm, userAgg, webAuthN, err := c.createUserPasskey(ctx, "user1", "org1", "rpID", domain.AuthenticatorAttachmentCrossPlattform)
			require.NoError(t, err)

			c.eventstore = eventstoreExpect(t, tt.expectPush(webAuthN.Challenge))

			got, err := c.pushUserU2F(ctx, wm, userAgg, webAuthN)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.NotEmpty(t, got.PublicKeyCredentialCreationOptions)
				assert.Equal(t, "123", got.ID)
				assert.Equal(t, "org1", got.ObjectDetails.ResourceOwner)
			}
		})
	}
}
