package command

import (
	"context"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommands_CreateIntent(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx        context.Context
		idpID      string
		successURL string
		failureURL string
	}
	type res struct {
		intentID string
		details  *domain.ObjectDetails
		err      error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"error no idpID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: errors.ThrowInvalidArgument(nil, "COMMAND-x8j2bk", "Errors.Intent.Invalid"),
			},
		},
		{
			"error no successURL",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "idp",
				successURL: ":",
				failureURL: "https://failure.url",
			},
			res{
				err: errors.ThrowInvalidArgument(nil, "COMMAND-x8j3bk", "Errors.Intent.Invalid"),
			},
		},
		{
			"error no failureURL",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: ":",
			},
			res{
				err: errors.ThrowInvalidArgument(nil, "COMMAND-x8j4bk", "Errors.Intent.Invalid"),
			},
		},
		{
			"error idp not existing",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(),
					expectFilter(),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: errors.ThrowPreconditionFailed(nil, "COMMAND-39n221fs", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			"push",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("ro").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								success, _ := url.Parse("https://success.url")
								failure, _ := url.Parse("https://failure.url")
								return idpintent.NewStartedEvent(
									context.Background(),
									&idpintent.NewAggregate("id", "ro").Aggregate,
									success,
									failure,
									"idp",
								)
							}(),
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				intentID: "id",
				details:  &domain.ObjectDetails{ResourceOwner: "ro"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			intentID, details, err := c.CreateIntent(tt.args.ctx, tt.args.idpID, tt.args.successURL, tt.args.failureURL)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.intentID, intentID)
			assert.Equal(t, tt.res.details, details)
		})
	}
}

func TestCommands_AuthURLFromProvider(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx         context.Context
		idpID       string
		state       string
		callbackURL string
	}
	type res struct {
		authURL string
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"idp not existing",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				ctx:         authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:       "idp",
				state:       "state",
				callbackURL: "url",
			},
			res{
				err: errors.ThrowPreconditionFailed(nil, "", ""),
			},
		},
		{
			"push",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								idp.Options{},
							)),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args{
				ctx:         authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:       "idp",
				state:       "state",
				callbackURL: "url",
			},
			res{
				authURL: "auth?client_id=clientID&prompt=select_account&redirect_uri=url&response_type=code&state=state",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			authURL, err := c.AuthURLFromProvider(tt.args.ctx, tt.args.idpID, tt.args.state, tt.args.callbackURL)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.authURL, authURL)
		})
	}
}
