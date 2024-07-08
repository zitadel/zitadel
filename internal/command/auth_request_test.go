package command

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddAuthRequest(t *testing.T) {
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id_generator.Generator
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								false,
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
			zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sf3gt", "Errors.AuthRequest.AlreadyExisting"),
		},
		{
			"added",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						authrequest.NewAddedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							"loginClient",
							"clientID",
							"redirectURI",
							"state",
							"nonce",
							[]string{"openid"},
							[]string{"audience"},
							domain.OIDCResponseTypeCode,
							domain.OIDCResponseModeQuery,
							&domain.OIDCCodeChallenge{
								Challenge: "challenge",
								Method:    domain.CodeChallengeMethodS256,
							},
							[]domain.Prompt{domain.PromptNone},
							[]string{"en", "de"},
							gu.Ptr(time.Duration(0)),
							gu.Ptr("loginHint"),
							gu.Ptr("hintUserID"),
							false,
						),
					),
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
					ResponseMode: domain.OIDCResponseModeQuery,
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
					ResponseMode: domain.OIDCResponseModeQuery,
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
				eventstore: tt.fields.eventstore,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.AddAuthRequest(tt.args.ctx, tt.args.request)
			require.ErrorIs(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_LinkSessionToAuthRequest(t *testing.T) {
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
	type fields struct {
		eventstore    *eventstore.Eventstore
		tokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
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
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:       mockCtx,
				id:        "id",
				sessionID: "sessionID",
			},
			res{
				wantErr: zerrors.ThrowNotFound(nil, "COMMAND-jae5P", "Errors.AuthRequest.NotExisting"),
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								true,
							),
						),
						eventFromEventPusher(
							authrequest.NewFailedEvent(mockCtx, &authrequest.NewAggregate("id", "instanceID").Aggregate,
								domain.OIDCErrorReasonUnspecified),
						),
					),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:       mockCtx,
				id:        "id",
				sessionID: "sessionID",
			},
			res{
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sx208nt", "Errors.AuthRequest.AlreadyHandled"),
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								true,
							),
						),
					),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:              authz.NewMockContext("instanceID", "orgID", "wrongLoginClient"),
				id:               "id",
				sessionID:        "sessionID",
				sessionToken:     "token",
				checkLoginClient: true,
			},
			res{
				wantErr: zerrors.ThrowPermissionDenied(nil, "COMMAND-rai9Y", "Errors.AuthRequest.WrongLoginClient"),
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								true,
							),
						),
					),
					expectFilter(),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:       mockCtx,
				id:        "V2_id",
				sessionID: "sessionID",
			},
			res{
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Flk38", "Errors.Session.NotExisting"),
			},
		},
		{
			"session expired",
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(mockCtx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewUserCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "org1", testNow.Add(-5*time.Minute), &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								testNow.Add(-5*time.Minute)),
						),
						eventFromEventPusher(
							session.NewLifetimeSetEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								2*time.Minute),
						),
					),
				),
			},
			args{
				ctx:          mockCtx,
				id:           "V2_id",
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Hkl3d", "Errors.Session.Expired"),
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(mockCtx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
					),
				),
				tokenVerifier: newMockTokenVerifierInvalid(),
			},
			args{
				ctx:          mockCtx,
				id:           "V2_id",
				sessionID:    "sessionID",
				sessionToken: "invalid",
			},
			res{
				wantErr: zerrors.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(mockCtx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewUserCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							session.NewLifetimeSetEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								2*time.Minute),
						),
					),
					expectPush(
						authrequest.NewSessionLinkedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							"sessionID",
							"userID",
							testNow,
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
						),
					),
				),
				tokenVerifier: newMockTokenVerifierValid(),
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
						ResponseMode: domain.OIDCResponseModeQuery,
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(mockCtx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewUserCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							session.NewLifetimeSetEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								2*time.Minute),
						),
					),
					expectPush(
						authrequest.NewSessionLinkedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							"sessionID",
							"userID",
							testNow,
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
						),
					),
				),
				tokenVerifier: newMockTokenVerifierValid(),
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
						ResponseMode: domain.OIDCResponseModeQuery,
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
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sx202nt", "Errors.AuthRequest.AlreadyHandled"),
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
								domain.OIDCResponseModeQuery,
								nil,
								nil,
								nil,
								nil,
								nil,
								nil,
								true,
							),
						),
					),
					expectPush(
						authrequest.NewFailedEvent(mockCtx, &authrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							domain.OIDCErrorReasonLoginRequired),
					),
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
						ResponseMode: domain.OIDCResponseModeQuery,
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
			zerrors.ThrowPreconditionFailed(nil, "COMMAND-Ht52d", "Errors.AuthRequest.InvalidCode"),
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
								domain.OIDCResponseModeQuery,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
								true,
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
			zerrors.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.AlreadyHandled"),
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
								domain.OIDCResponseModeQuery,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
								true,
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
						authrequest.NewCodeAddedEvent(mockCtx, &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
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
