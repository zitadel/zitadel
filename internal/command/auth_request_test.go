package command

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/oidc/amr"
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
	type res struct {
		details *domain.ObjectDetails
		wantErr error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
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
			res{
				wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sx208nt", "Errors.AuthRequest.AlreadyHandled"),
			},
		},
		{
			"session not existing",
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
					expectFilter(),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:       authz.WithInstanceID(context.Background(), "instanceID"),
				id:        "V2_id",
				sessionID: "session",
			},
			res{
				wantErr: caos_errs.ThrowNotFound(nil, "COMMAND-x0099887", "Errors.Session.NotExisting"),
			},
		},
		{
			"missing permission",
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
				id:        "V2_id",
				sessionID: "sessionID",
			},
			res{
				wantErr: caos_errs.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			"invalid session token",
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
				id:           "V2_id",
				sessionID:    "sessionID",
				sessionToken: "invalid",
			},
			res{
				wantErr: caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
			},
		},
		{
			"linked",
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
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
					expectPush(
						[]*repository.Event{eventFromEventPusherWithInstanceID(
							"instanceID",
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]string{amr.PWD},
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
				id:           "V2_id",
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:           tt.fields.eventstore,
				sessionTokenVerifier: tt.fields.tokenVerifier,
				checkPermission:      tt.fields.checkPermission,
			}
			details, err := c.LinkSessionToAuthRequest(tt.args.ctx, tt.args.id, tt.args.sessionID, tt.args.sessionToken)
			require.ErrorIs(t, err, tt.res.wantErr)
			assert.Equal(t, tt.res.details, details)
		})
	}
}

func TestCommands_AddAuthRequestCode(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx  context.Context
		id   string
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			"empty code error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				id:   "V2_authRequestID",
				code: "",
			},
			caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Ht52d", "Errors.AuthRequest.InvalidCode"),
		},
		{
			"no session linked error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
						),
					),
				),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				id:   "V2_authRequestID",
				code: "code",
			},
			caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.AlreadyHandled"),
		},
		{
			"success",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
						),
						eventFromEventPusher(
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]string{amr.PWD},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate, "code"),
							),
						},
					),
				),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				id:   "V2_authRequestID",
				code: "code",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := c.AddAuthRequestCode(tt.args.ctx, tt.args.id, tt.args.code)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}

func TestCommands_ExchangeAuthCode(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx  context.Context
		code string
	}
	type res struct {
		authRequest *AuthenticatedAuthRequest
		err         error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		//{
		//	"empty code error",
		//	fields{
		//		eventstore: eventstoreExpect(t),
		//	},
		//	args{
		//		ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
		//		code: "",
		//	},
		//	res{
		//		err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sfr3s", "Errors.AuthRequest.InvalidCode"),
		//	},
		//},
		{
			"no code added error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
						),
					),
				),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				code: "V2_authRequestID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.NoCode"),
			},
		},
		//{
		//	"invalid code error",
		//	fields{
		//		eventstore: eventstoreExpect(t,
		//			expectFilter(
		//				eventFromEventPusher(
		//					authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
		//						"loginClient",
		//						"clientID",
		//						"redirectURI",
		//						"state",
		//						"nonce",
		//						[]string{"openid"},
		//						[]string{"audience"},
		//						domain.OIDCResponseTypeCode,
		//						&domain.OIDCCodeChallenge{
		//							Challenge: "challenge",
		//							Method:    domain.CodeChallengeMethodS256,
		//						},
		//						[]domain.Prompt{domain.PromptNone},
		//						[]string{"en", "de"},
		//						gu.Ptr(time.Duration(0)),
		//						gu.Ptr("loginHint"),
		//						gu.Ptr("hintUserID"),
		//					),
		//				),
		//				eventFromEventPusher(
		//					authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate, "code"),
		//				),
		//			),
		//		),
		//	},
		//	args{
		//		ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
		//		code: "invalidCode",
		//	},
		//	res{
		//		err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-DBNqz", "Errors.AuthRequest.InvalidCode"),
		//	},
		//},
		{
			"code exchanged",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
						),
						eventFromEventPusher(
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]string{amr.PWD},
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate, "code"),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
							),
						},
					),
				),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				code: "V2_authRequestID",
			},
			res{
				authRequest: &AuthenticatedAuthRequest{
					AuthRequest: &AuthRequest{
						ID:           "authRequestID",
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
					SessionID: "sessionID",
					UserID:    "userID",
					AMR:       []string{"pwd"},
					AuthTime:  testNow,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.ExchangeAuthCode(tt.args.ctx, tt.args.code)
			assert.Equal(t, tt.res.authRequest, got)
			assert.ErrorIs(t, tt.res.err, err)
		})
	}
}
