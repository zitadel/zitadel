package command

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
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
		want    *CurrentAuthRequest
		wantErr error
	}{
		{
			"already exists error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
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
				ctx:     mockCtx,
				request: &AuthRequest{},
			},
			nil,
			caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sf3gt", "Errors.AuthRequest.AlreadyExisting"),
		},
		{
			"added",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
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
						}),
				),
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args{
				ctx: mockCtx,
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
			&CurrentAuthRequest{
				AuthRequest: &AuthRequest{
					ID:           "V2_id",
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
			got, err := c.AddAuthRequest(tt.args.ctx, tt.args.request)
			require.ErrorIs(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_LinkSessionToAuthRequest(t *testing.T) {
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
	type fields struct {
		eventstore      *eventstore.Eventstore
		tokenVerifier   func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx              context.Context
		id               string
		sessionID        string
		sessionToken     string
		checkLoginClient bool
	}
	type res struct {
		details *domain.ObjectDetails
		authReq *CurrentAuthRequest
		wantErr error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"authRequest not found",
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
				ctx:       mockCtx,
				id:        "id",
				sessionID: "sessionID",
			},
			res{
				wantErr: caos_errs.ThrowNotFound(nil, "COMMAND-jae5P", "Errors.AuthRequest.NotExisting"),
			},
		},
		{
			"authRequest not existing",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("id", "instanceID").Aggregate,
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
						eventFromEventPusher(
							authrequest.NewFailedEvent(mockCtx, &authrequest.NewAggregate("id", "instanceID").Aggregate,
								domain.OIDCErrorReasonUnspecified),
						),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:       mockCtx,
				id:        "id",
				sessionID: "sessionID",
			},
			res{
				wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sx208nt", "Errors.AuthRequest.AlreadyHandled"),
			},
		},
		{
			"wrong login client",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("id", "instanceID").Aggregate,
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
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:              authz.NewMockContext("instanceID", "orgID", "wrongLoginClient"),
				id:               "id",
				sessionID:        "sessionID",
				sessionToken:     "token",
				checkLoginClient: true,
			},
			res{
				wantErr: caos_errs.ThrowPermissionDenied(nil, "COMMAND-rai9Y", "Errors.AuthRequest.WrongLoginClient"),
			},
		},
		{
			"session not existing",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
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
				ctx:       mockCtx,
				id:        "V2_id",
				sessionID: "sessionID",
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
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
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
							session.NewAddedEvent(mockCtx, &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld")),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:       mockCtx,
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
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
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
							session.NewAddedEvent(mockCtx, &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld")),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
				},
			},
			args{
				ctx:          mockCtx,
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
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
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
							session.NewAddedEvent(mockCtx, &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld"),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "org1").Aggregate,
								"userID", testNow),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "org1").Aggregate,
								testNow),
						),
					),
					expectPush(
						[]*repository.Event{eventFromEventPusherWithInstanceID(
							"instanceID",
							authrequest.NewSessionLinkedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						)}),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:          mockCtx,
				id:           "V2_id",
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
				authReq: &CurrentAuthRequest{
					AuthRequest: &AuthRequest{
						ID:           "V2_id",
						LoginClient:  "loginClient",
						ClientID:     "clientID",
						RedirectURI:  "redirectURI",
						State:        "state",
						Nonce:        "nonce",
						Scope:        []string{"openid"},
						Audience:     []string{"audience"},
						ResponseType: domain.OIDCResponseTypeCode,
					},
					SessionID:   "sessionID",
					UserID:      "userID",
					AuthMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				},
			},
		},
		{
			"linked with login client check",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
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
							session.NewAddedEvent(mockCtx, &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld"),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "org1").Aggregate,
								"userID", testNow),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "org1").Aggregate,
								testNow),
						),
					),
					expectPush(
						[]*repository.Event{eventFromEventPusherWithInstanceID(
							"instanceID",
							authrequest.NewSessionLinkedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						)}),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:              authz.NewMockContext("instanceID", "orgID", "loginClient"),
				id:               "V2_id",
				sessionID:        "sessionID",
				sessionToken:     "token",
				checkLoginClient: true,
			},
			res{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
				authReq: &CurrentAuthRequest{
					AuthRequest: &AuthRequest{
						ID:           "V2_id",
						LoginClient:  "loginClient",
						ClientID:     "clientID",
						RedirectURI:  "redirectURI",
						State:        "state",
						Nonce:        "nonce",
						Scope:        []string{"openid"},
						Audience:     []string{"audience"},
						ResponseType: domain.OIDCResponseTypeCode,
					},
					SessionID:   "sessionID",
					UserID:      "userID",
					AuthMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				},
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
			details, got, err := c.LinkSessionToAuthRequest(tt.args.ctx, tt.args.id, tt.args.sessionID, tt.args.sessionToken, tt.args.checkLoginClient)
			require.ErrorIs(t, err, tt.res.wantErr)
			assert.Equal(t, tt.res.details, details)
			if err == nil {
				assert.WithinRange(t, got.AuthTime, testNow, testNow)
				got.AuthTime = time.Time{}
			}
			assert.Equal(t, tt.res.authReq, got)
		})
	}
}

func TestCommands_FailAuthRequest(t *testing.T) {
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		id     string
		reason domain.OIDCErrorReason
	}
	type res struct {
		details *domain.ObjectDetails
		authReq *CurrentAuthRequest
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
			},
			args{
				ctx:    mockCtx,
				id:     "foo",
				reason: domain.OIDCErrorReasonLoginRequired,
			},
			res{
				wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sx202nt", "Errors.AuthRequest.AlreadyHandled"),
			},
		},
		{
			"failed",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
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
					expectPush(
						[]*repository.Event{eventFromEventPusherWithInstanceID(
							"instanceID",
							authrequest.NewFailedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								domain.OIDCErrorReasonLoginRequired),
						)}),
				),
			},
			args{
				ctx:    mockCtx,
				id:     "V2_id",
				reason: domain.OIDCErrorReasonLoginRequired,
			},
			res{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
				authReq: &CurrentAuthRequest{
					AuthRequest: &AuthRequest{
						ID:           "V2_id",
						LoginClient:  "loginClient",
						ClientID:     "clientID",
						RedirectURI:  "redirectURI",
						State:        "state",
						Nonce:        "nonce",
						Scope:        []string{"openid"},
						Audience:     []string{"audience"},
						ResponseType: domain.OIDCResponseTypeCode,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			details, got, err := c.FailAuthRequest(tt.args.ctx, tt.args.id, tt.args.reason)
			require.ErrorIs(t, err, tt.res.wantErr)
			assert.Equal(t, tt.res.details, details)
			assert.Equal(t, tt.res.authReq, got)
		})
	}
}

func TestCommands_AddAuthRequestCode(t *testing.T) {
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
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
				ctx:  mockCtx,
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
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
				ctx:  mockCtx,
				id:   "V2_authRequestID",
				code: "V2_authRequestID",
			},
			caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.AlreadyHandled"),
		},
		{
			"success",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
							authrequest.NewSessionLinkedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewCodeAddedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
							),
						},
					),
				),
			},
			args{
				ctx:  mockCtx,
				id:   "V2_authRequestID",
				code: "V2_authRequestID",
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
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx  context.Context
		code string
	}
	type res struct {
		authRequest *CurrentAuthRequest
		err         error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"empty code error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:  mockCtx,
				code: "",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sf3g2", "Errors.AuthRequest.InvalidCode"),
			},
		},
		{
			"no code added error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
				ctx:  mockCtx,
				code: "V2_authRequestID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.NoCode"),
			},
		},
		{
			"code exchanged",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
							authrequest.NewSessionLinkedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewCodeExchangedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
							),
						},
					),
				),
			},
			args{
				ctx:  mockCtx,
				code: "V2_authRequestID",
			},
			res{
				authRequest: &CurrentAuthRequest{
					AuthRequest: &AuthRequest{
						ID:           "V2_authRequestID",
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
					SessionID:   "sessionID",
					UserID:      "userID",
					AuthMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
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
			assert.ErrorIs(t, tt.res.err, err)

			if err == nil {
				// equal on time won't work -> test separately and clear it before comparing the rest
				assert.WithinRange(t, got.AuthTime, testNow, testNow)
				got.AuthTime = time.Time{}
			}
			assert.Equal(t, tt.res.authRequest, got)
		})
	}
}
