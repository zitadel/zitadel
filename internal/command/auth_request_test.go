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
		{
			"empty code error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				code: "",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Sfr3s", "Errors.AuthRequest.InvalidCode"),
			},
		},
		{
			"no code added error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
					//eventFromEventPusher(
					//	authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate,
					//		"loginClient",
					//		"clientID",
					//		"redirectURI",
					//		"state",
					//		"nonce",
					//		[]string{"openid"},
					//		[]string{"audience"},
					//		domain.OIDCResponseTypeCode,
					//		&domain.OIDCCodeChallenge{
					//			Challenge: "challenge",
					//			Method:    domain.CodeChallengeMethodS256,
					//		},
					//		[]domain.Prompt{domain.PromptNone},
					//		[]string{"en", "de"},
					//		gu.Ptr(time.Duration(0)),
					//		"loginHint",
					//		"hintUserID",
					//	),
					//),
					),
				),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				code: "authRequestID:codeID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SFwd2", "Errors.AuthRequest.NoCode"),
			},
		},
		{
			"invalid code error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate,
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
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						),
					),
				),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				code: "authRequestID:invalidCodeID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-DBNqz", "Errors.AuthRequest.InvalidCode"),
			},
		},
		{
			"code exchanged",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate,
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

						//TODO: session link event
						//eventFromEventPusher(
						//	authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						//),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
							),
						},
					),
				),
			},
			args{
				ctx:  authz.WithInstanceID(context.Background(), "instanceID"),
				code: "authRequestID:",
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
