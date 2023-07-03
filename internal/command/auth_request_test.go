package command

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/session"
)

func TestCommands_AddAuthRequest(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx     context.Context
		request *AuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			"already exists error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
							),
						),
					),
				),
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args{
				ctx:     authz.WithInstanceID(context.Background(), "instanceID"),
				request: &AuthRequest{},
			},
			caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sf3gt", "Errors.AuthRequest.AlreadyExisting"),
		},
		{
			"added",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{eventFromEventPusherWithInstanceID(
							"instanceID",
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
							),
						)}),
				),
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				request: &AuthRequest{
					LoginClient:  "loginClient",
					ClientID:     "clientID",
					RedirectURI:  "redirectURI",
					State:        "state",
					Nonce:        "nonce",
					Scope:        []string{"openid"},
					Audience:     []string{"audience"},
					ResponseType: domain.OIDCResponseTypeCode,
					CodeChallenge: &domain.OIDCCodeChallenge{
						Challenge: "challenge",
						Method:    domain.CodeChallengeMethodS256,
					},
					Prompt:     []domain.Prompt{domain.PromptNone},
					UILocales:  []string{"en", "de"},
					MaxAge:     gu.Ptr(time.Duration(0)),
					LoginHint:  gu.Ptr("loginHint"),
					HintUserID: gu.Ptr("hintUserID"),
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			err := c.AddAuthRequest(tt.args.ctx, tt.args.request)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}

func TestCommands_LinkSessionToAuthRequest(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		tokenVerifier   func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx          context.Context
		id           string
		sessionID    string
		sessionToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			"authRequest not existing",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:       authz.WithInstanceID(context.Background(), "instanceID"),
				id:        "id",
				sessionID: "session",
			},
			caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sx208nt", "Errors.AuthRequest.AlreadyHandled"),
		},
		{
			"session not existing",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("id", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								nil,
								nil,
								nil,
								nil,
								"",
								"",
							),
						),
					),
					expectFilter(),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:       authz.WithInstanceID(context.Background(), "instanceID"),
				id:        "id",
				sessionID: "session",
			},
			caos_errs.ThrowNotFound(nil, "COMMAND-x0099887", "Errors.Session.NotExisting"),
		},
		{
			"missing permission",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("id", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								nil,
								nil,
								nil,
								nil,
								"",
								"",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:       authz.WithInstanceID(context.Background(), "instanceID"),
				id:        "id",
				sessionID: "session",
			},
			caos_errs.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			"invalid session token",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("id", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								nil,
								nil,
								nil,
								nil,
								"",
								"",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
				},
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				id:           "id",
				sessionID:    "session",
				sessionToken: "invalid",
			},
			caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
		},
		{
			"linked",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("id", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								nil,
								nil,
								nil,
								nil,
								"",
								"",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
					expectPush(
						[]*repository.Event{eventFromEventPusherWithInstanceID(
							"instanceID",
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("id", "instanceID").Aggregate,
								"session",
							),
						)}),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				id:           "id",
				sessionID:    "session",
				sessionToken: "token",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:           tt.fields.eventstore,
				sessionTokenVerifier: tt.fields.tokenVerifier,
				checkPermission:      tt.fields.checkPermission,
			}
			err := c.LinkSessionToAuthRequest(tt.args.ctx, tt.args.id, tt.args.sessionID, tt.args.sessionToken)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}
