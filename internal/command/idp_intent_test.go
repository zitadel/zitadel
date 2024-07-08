package command

import (
	"context"
	"net/url"
	"testing"

	"github.com/crewjam/saml"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	rep_idp "github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateIntent(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id_generator.Generator
	}
	type args struct {
		ctx        context.Context
		idpID      string
		successURL string
		failureURL string
		instanceID string
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
			"error no id generator",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: mock.NewIDGeneratorExpectError(t, zerrors.ThrowInternal(nil, "", "error id")),
			},
			args{
				ctx:        context.Background(),
				idpID:      "",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: zerrors.ThrowInternal(nil, "", "error id"),
			},
		},
		{
			"error no idpID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        context.Background(),
				idpID:      "",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-x8j2bk", "Errors.Intent.IDPMissing"),
			},
		},
		{
			"error no successURL",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        context.Background(),
				idpID:      "idp",
				successURL: ":",
				failureURL: "https://failure.url",
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-x8j3bk", "Errors.Intent.SuccessURLMissing"),
			},
		},
		{
			"error no failureURL",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        context.Background(),
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: ":",
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-x8j4bk", "Errors.Intent.FailureURLMissing"),
			},
		},
		{
			"error idp not existing org",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        context.Background(),
				idpID:      "idp",
				instanceID: "instance",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-39n221fs", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			"error idp not existing instance",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        context.Background(),
				idpID:      "idp",
				instanceID: "instance",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-39n221fs", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			"push, org",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org").Aggregate,
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
								rep_idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							success, _ := url.Parse("https://success.url")
							failure, _ := url.Parse("https://failure.url")
							return idpintent.NewStartedEvent(
								context.Background(),
								&idpintent.NewAggregate("id", "instance").Aggregate,
								success,
								failure,
								"idp",
							)
						}(),
					),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        context.Background(),
				instanceID: "instance",
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				intentID: "id",
				details:  &domain.ObjectDetails{ResourceOwner: "instance"},
			},
		},
		{
			"push, instance",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
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
								rep_idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							success, _ := url.Parse("https://success.url")
							failure, _ := url.Parse("https://failure.url")
							return idpintent.NewStartedEvent(
								context.Background(),
								&idpintent.NewAggregate("id", "instance").Aggregate,
								success,
								failure,
								"idp",
							)
						}(),
					),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        context.Background(),
				instanceID: "instance",
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				intentID: "id",
				details:  &domain.ObjectDetails{ResourceOwner: "instance"},
			},
		},
		{
			"push, instance without org",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
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
								rep_idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							success, _ := url.Parse("https://success.url")
							failure, _ := url.Parse("https://failure.url")
							return idpintent.NewStartedEvent(
								context.Background(),
								&idpintent.NewAggregate("id", "instance").Aggregate,
								success,
								failure,
								"idp",
							)
						}(),
					),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        context.Background(),
				instanceID: "instance",
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				intentID: "id",
				details:  &domain.ObjectDetails{ResourceOwner: "instance"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			intentWriteModel, details, err := c.CreateIntent(tt.args.ctx, tt.args.idpID, tt.args.successURL, tt.args.failureURL, tt.args.instanceID)
			require.ErrorIs(t, err, tt.res.err)
			if intentWriteModel != nil {
				assert.Equal(t, tt.res.intentID, intentWriteModel.AggregateID)
			} else {
				assert.Equal(t, tt.res.intentID, "")
			}
			assert.Equal(t, tt.res.details, details)
		})
	}
}

func TestCommands_AuthFromProvider(t *testing.T) {
	type fields struct {
		eventstore   func(t *testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx         context.Context
		idpID       string
		state       string
		callbackURL string
		samlRootURL string
	}
	type res struct {
		content  string
		redirect bool
		err      error
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
				eventstore: expectEventstore(
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
				err: zerrors.ThrowPreconditionFailed(nil, "", ""),
			},
		},
		{
			"idp removed",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: expectEventstore(
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
								rep_idp.Options{},
							)),
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewIDPRemovedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
							),
						),
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
				err: zerrors.ThrowInternal(nil, "COMMAND-xw921211", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			"oauth auth redirect",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: expectEventstore(
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
								rep_idp.Options{},
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
								rep_idp.Options{},
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
				content:  "auth?client_id=clientID&prompt=select_account&redirect_uri=url&response_type=code&state=state",
				redirect: true,
			},
		},
		{
			"migrated and push",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid", "profile", "User.Read"},
								false,
								rep_idp.Options{},
							)),
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOIDCIDPMigratedAzureADEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid", "profile", "User.Read"},
								"tenant",
								true,
								rep_idp.Options{},
							)),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid", "profile", "User.Read"},
								false,
								rep_idp.Options{},
							)),
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOIDCIDPMigratedAzureADEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid", "profile", "User.Read"},
								"tenant",
								true,
								rep_idp.Options{},
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
				content:  "https://login.microsoftonline.com/tenant/oauth2/v2.0/authorize?client_id=clientID&prompt=select_account&redirect_uri=url&response_type=code&scope=openid+profile+User.Read&state=state",
				redirect: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			content, redirect, err := c.AuthFromProvider(tt.args.ctx, tt.args.idpID, tt.args.state, tt.args.callbackURL, tt.args.samlRootURL)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.redirect, redirect)
			assert.Equal(t, tt.res.content, content)
		})
	}
}

func TestCommands_AuthFromProvider_SAML(t *testing.T) {
	type fields struct {
		eventstore   func(t *testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx         context.Context
		idpID       string
		state       string
		callbackURL string
		samlRootURL string
	}
	type res struct {
		url    string
		values map[string]string
		err    error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"saml auth default redirect",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewSAMLIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								[]byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-08-27T12:40:58.803Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("key"),
								},
								[]byte("certificate"),
								"",
								false,
								gu.Ptr(domain.SAMLNameIDFormatUnspecified),
								"",
								rep_idp.Options{},
							)),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewSAMLIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								[]byte("<EntityDescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" validUntil=\"2023-08-27T12:40:58.803Z\" cacheDuration=\"PT48H\" entityID=\"http://localhost:8000/metadata\">\n  <IDPSSODescriptor xmlns=\"urn:oasis:names:tc:SAML:2.0:metadata\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n    <KeyDescriptor use=\"signing\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n    </KeyDescriptor>\n    <KeyDescriptor use=\"encryption\">\n      <KeyInfo xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n        <X509Data xmlns=\"http://www.w3.org/2000/09/xmldsig#\">\n          <X509Certificate xmlns=\"http://www.w3.org/2000/09/xmldsig#\">MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8Ahs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+aucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWxm+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURNB2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0OBBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uvNONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEfy/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsbGFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTLUzreO96WzlBBMtY=</X509Certificate>\n        </X509Data>\n      </KeyInfo>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes128-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes192-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#aes256-cbc\"></EncryptionMethod>\n      <EncryptionMethod Algorithm=\"http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p\"></EncryptionMethod>\n    </KeyDescriptor>\n    <NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</NameIDFormat>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n    <SingleSignOnService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\" Location=\"http://localhost:8000/sso\"></SingleSignOnService>\n  </IDPSSODescriptor>\n</EntityDescriptor>"),
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAxHd087RoEm9ywVWZ/H+tDWxQsmVvhfRz4jAq/RfU+OWXNH4J\njMMSHdFs0Q+WP98nNXRyc7fgbMb8NdmlB2yD4qLYapN5SDaBc5dh/3EnyFt53oSs\njTlKnQUPAeJr2qh/NY046CfyUyQMM4JR5OiQFo4TssfWnqdcgamGt0AEnk2lvbMZ\nKQdAqNS9lDzYbjMGavEQPTZE35mFXFQXjaooZXq+TIa7hbaq7/idH7cHNbLcPLgj\nfPQA8q+DYvnvhXlmq0LPQZH3Oiixf+SF2vRwrBzT2mqGD2OiOkUmhuPwyqEiiBHt\nfxklRtRU6WfLa1Gcb1PsV0uoBGpV3KybIl/GlwIDAQABAoIBAEQjDduLgOCL6Gem\n0X3hpdnW6/HC/jed/Sa//9jBECq2LYeWAqff64ON40hqOHi0YvvGA/+gEOSI6mWe\nsv5tIxxRz+6+cLybsq+tG96kluCE4TJMHy/nY7orS/YiWbd+4odnEApr+D3fbZ/b\nnZ1fDsHTyn8hkYx6jLmnWsJpIHDp7zxD76y7k2Bbg6DZrCGiVxngiLJk23dvz79W\np03lHLM7XE92aFwXQmhfxHGxrbuoB/9eY4ai5IHp36H4fw0vL6NXdNQAo/bhe0p9\nAYB7y0ZumF8Hg0Z/BmMeEzLy6HrYB+VE8cO93pNjhSyH+p2yDB/BlUyTiRLQAoM0\nVTmOZXECgYEA7NGlzpKNhyQEJihVqt0MW0LhKIO/xbBn+XgYfX6GpqPa/ucnMx5/\nVezpl3gK8IU4wPUhAyXXAHJiqNBcEeyxrw0MXLujDVMJgYaLysCLJdvMVgoY08mS\nK5IQivpbozpf4+0y3mOnA+Sy1kbfxv2X8xiWLODRQW3f3q/xoklwOR8CgYEA1GEe\nfaibOFTQAYcIVj77KXtBfYZsX3EGAyfAN9O7cKHq5oaxVstwnF47WxpuVtoKZxCZ\nbNm9D5WvQ9b+Ztpioe42tzwE7Bff/Osj868GcDdRPK7nFlh9N2yVn/D514dOYVwR\n4MBr1KrJzgRWt4QqS4H+to1GzudDTSNlG7gnK4kCgYBUi6AbOHzoYzZL/RhgcJwp\ntJ23nhmH1Su5h2OO4e3mbhcP66w19sxU+8iFN+kH5zfUw26utgKk+TE5vXExQQRK\nT2k7bg2PAzcgk80ybD0BHhA8I0yrx4m0nmfjhe/TPVLgh10iwgbtP+eM0i6v1vc5\nZWyvxu9N4ZEL6lpkqr0y1wKBgG/NAIQd8jhhTW7Aav8cAJQBsqQl038avJOEpYe+\nCnpsgoAAf/K0/f8TDCQVceh+t+MxtdK7fO9rWOxZjWsPo8Si5mLnUaAHoX4/OpnZ\nlYYVWMqdOEFnK+O1Yb7k2GFBdV2DXlX2dc1qavntBsls5ecB89id3pyk2aUN8Pf6\npYQhAoGAMGtrHFely9wyaxI0RTCyfmJbWZHGVGkv6ELK8wneJjdjl82XOBUGCg5q\naRCrTZ3dPitKwrUa6ibJCIFCIziiriBmjDvTHzkMvoJEap2TVxYNDR6IfINVsQ57\nlOsiC4A2uGq4Lbfld+gjoplJ5GX6qXtTgZ6m7eo0y7U6zm2tkN0=\n-----END RSA PRIVATE KEY-----\n"),
								}, []byte("-----BEGIN CERTIFICATE-----\nMIIC2zCCAcOgAwIBAgIIAy/jm1gAAdEwDQYJKoZIhvcNAQELBQAwEjEQMA4GA1UE\nChMHWklUQURFTDAeFw0yMzA4MzAwNzExMTVaFw0yNDA4MjkwNzExMTVaMBIxEDAO\nBgNVBAoTB1pJVEFERUwwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDE\nd3TztGgSb3LBVZn8f60NbFCyZW+F9HPiMCr9F9T45Zc0fgmMwxId0WzRD5Y/3yc1\ndHJzt+Bsxvw12aUHbIPiothqk3lINoFzl2H/cSfIW3nehKyNOUqdBQ8B4mvaqH81\njTjoJ/JTJAwzglHk6JAWjhOyx9aep1yBqYa3QASeTaW9sxkpB0Co1L2UPNhuMwZq\n8RA9NkTfmYVcVBeNqihler5MhruFtqrv+J0ftwc1stw8uCN89ADyr4Ni+e+FeWar\nQs9Bkfc6KLF/5IXa9HCsHNPaaoYPY6I6RSaG4/DKoSKIEe1/GSVG1FTpZ8trUZxv\nU+xXS6gEalXcrJsiX8aXAgMBAAGjNTAzMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUE\nDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMA0GCSqGSIb3DQEBCwUAA4IBAQCx\n/dRNIj0N/16zJhZR/ahkc2AkvDXYxyr4JRT5wK9GQDNl/oaX3debRuSi/tfaXFIX\naJA6PxM4J49ZaiEpLrKfxMz5kAhjKchCBEMcH3mGt+iNZH7EOyTvHjpGrP2OZrsh\nO17yrvN3HuQxIU6roJlqtZz2iAADsoPtwOO4D7hupm9XTMkSnAmlMWOo/q46Jz89\n1sMxB+dXmH/zV0wgwh0omZfLV0u89mvdq269VhcjNBpBYSnN1ccqYWd5iwziob3I\nvaavGHGfkbvRUn/tKftYuTK30q03R+e9YbmlWZ0v695owh2e/apCzowQsCKfSVC8\nOxVyt5XkHq1tWwVyBmFp\n-----END CERTIFICATE-----\n"),
								"",
								false,
								gu.Ptr(domain.SAMLNameIDFormatUnspecified),
								"",
								rep_idp.Options{},
							)),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							func() eventstore.Command {
								success, _ := url.Parse("https://success.url")
								failure, _ := url.Parse("https://failure.url")
								return idpintent.NewStartedEvent(
									context.Background(),
									&idpintent.NewAggregate("id", "instance").Aggregate,
									success,
									failure,
									"idp",
								)
							}(),
						),
					),
					expectRandomPush(
						[]eventstore.Command{
							idpintent.NewSAMLRequestEvent(
								context.Background(),
								&idpintent.NewAggregate("id", "instance").Aggregate,
								"request",
							),
						},
					),
				),
			},
			args{
				ctx:         authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:       "idp",
				state:       "id",
				callbackURL: "url",
				samlRootURL: "samlurl",
			},
			res{
				url: "http://localhost:8000/sso",
				values: map[string]string{
					"SAMLRequest": "", // generated IDs so not assertable
					"RelayState":  "id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			content, _, err := c.AuthFromProvider(tt.args.ctx, tt.args.idpID, tt.args.state, tt.args.callbackURL, tt.args.samlRootURL)
			require.ErrorIs(t, err, tt.res.err)

			authURL, err := url.Parse(content)
			require.NoError(t, err)

			assert.Equal(t, tt.res.url, authURL.Scheme+"://"+authURL.Host+authURL.Path)
			query := authURL.Query()
			for k, v := range tt.res.values {
				assert.True(t, query.Has(k))
				if v != "" {
					assert.Equal(t, v, query.Get(k))
				}
			}
		})
	}
}

func TestCommands_SucceedIDPIntent(t *testing.T) {
	type fields struct {
		eventstore          func(t *testing.T) *eventstore.Eventstore
		idpConfigEncryption crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		writeModel *IDPIntentWriteModel
		idpUser    idp.User
		idpSession idp.Session
		userID     string
	}
	type res struct {
		token string
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"encryption fails",
			fields{
				idpConfigEncryption: func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, zerrors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
				eventstore: expectEventstore(),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
			},
			res{
				err: zerrors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"token encryption fails",
			fields{
				idpConfigEncryption: func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).DoAndReturn(func(value []byte) ([]byte, error) {
						return value, nil
					})
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, zerrors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
				eventstore: expectEventstore(),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
				idpSession: &oauth.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
					},
				},
			},
			res{
				err: zerrors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"push",
			fields{
				idpConfigEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: expectEventstore(
					expectPush(
						func() eventstore.Command {
							event := idpintent.NewSucceededEvent(
								context.Background(),
								&idpintent.NewAggregate("id", "instance").Aggregate,
								[]byte(`{"sub":"id","preferred_username":"username"}`),
								"id",
								"username",
								"",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("accessToken"),
								},
								"idToken",
							)
							return event
						}(),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "instance"),
				idpSession: &openid.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
						IDToken: "idToken",
					},
				},
				idpUser: openid.NewUser(&oidc.UserInfo{
					Subject: "id",
					UserInfoProfile: oidc.UserInfoProfile{
						PreferredUsername: "username",
					},
				}),
			},
			res{
				token: "aWQ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.idpConfigEncryption,
			}
			got, err := c.SucceedIDPIntent(tt.args.ctx, tt.args.writeModel, tt.args.idpUser, tt.args.idpSession, tt.args.userID)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.token, got)
		})
	}
}

func TestCommands_SucceedSAMLIDPIntent(t *testing.T) {
	type fields struct {
		eventstore          func(t *testing.T) *eventstore.Eventstore
		idpConfigEncryption crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		writeModel *IDPIntentWriteModel
		idpUser    idp.User
		assertion  *saml.Assertion
		userID     string
	}
	type res struct {
		token string
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"encryption fails",
			fields{
				idpConfigEncryption: func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, zerrors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
				eventstore: expectEventstore(),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
			},
			res{
				err: zerrors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"push",
			fields{
				idpConfigEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: expectEventstore(
					expectPush(
						idpintent.NewSAMLSucceededEvent(
							context.Background(),
							&idpintent.NewAggregate("id", "instance").Aggregate,
							[]byte(`{"sub":"id","preferred_username":"username"}`),
							"id",
							"username",
							"",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("<Assertion xmlns=\"urn:oasis:names:tc:SAML:2.0:assertion\" ID=\"id\" IssueInstant=\"0001-01-01T00:00:00Z\" Version=\"\"><Issuer xmlns=\"urn:oasis:names:tc:SAML:2.0:assertion\" NameQualifier=\"\" SPNameQualifier=\"\" Format=\"\" SPProvidedID=\"\"></Issuer></Assertion>"),
							},
						),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "instance"),
				assertion:  &saml.Assertion{ID: "id"},
				idpUser: openid.NewUser(&oidc.UserInfo{
					Subject: "id",
					UserInfoProfile: oidc.UserInfoProfile{
						PreferredUsername: "username",
					},
				}),
			},
			res{
				token: "aWQ",
			},
		},
		{
			"push with userID",
			fields{
				idpConfigEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: expectEventstore(
					expectPush(
						idpintent.NewSAMLSucceededEvent(
							context.Background(),
							&idpintent.NewAggregate("id", "instance").Aggregate,
							[]byte(`{"sub":"id","preferred_username":"username"}`),
							"id",
							"username",
							"user",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("<Assertion xmlns=\"urn:oasis:names:tc:SAML:2.0:assertion\" ID=\"id\" IssueInstant=\"0001-01-01T00:00:00Z\" Version=\"\"><Issuer xmlns=\"urn:oasis:names:tc:SAML:2.0:assertion\" NameQualifier=\"\" SPNameQualifier=\"\" Format=\"\" SPProvidedID=\"\"></Issuer></Assertion>"),
							},
						),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "instance"),
				assertion:  &saml.Assertion{ID: "id"},
				idpUser: openid.NewUser(&oidc.UserInfo{
					Subject: "id",
					UserInfoProfile: oidc.UserInfoProfile{
						PreferredUsername: "username",
					},
				}),
				userID: "user",
			},
			res{
				token: "aWQ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.idpConfigEncryption,
			}
			got, err := c.SucceedSAMLIDPIntent(tt.args.ctx, tt.args.writeModel, tt.args.idpUser, tt.args.userID, tt.args.assertion)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.token, got)
		})
	}
}

func TestCommands_RequestSAMLIDPIntent(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		writeModel *IDPIntentWriteModel
		request    string
	}
	type res struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"push",
			fields{
				eventstore: expectEventstore(
					expectPush(
						idpintent.NewSAMLRequestEvent(
							context.Background(),
							&idpintent.NewAggregate("id", "instance").Aggregate,
							"request",
						),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "instance"),
				request:    "request",
			},
			res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := c.RequestSAMLIDPIntent(tt.args.ctx, tt.args.writeModel, tt.args.request)
			require.ErrorIs(t, err, tt.res.err)
			require.Equal(t, tt.args.writeModel.RequestID, tt.args.request)
		})
	}
}

func TestCommands_SucceedLDAPIDPIntent(t *testing.T) {
	type fields struct {
		eventstore          func(t *testing.T) *eventstore.Eventstore
		idpConfigEncryption crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		writeModel *IDPIntentWriteModel
		idpUser    idp.User
		userID     string
		attributes map[string][]string
	}
	type res struct {
		token string
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"encryption fails",
			fields{
				idpConfigEncryption: func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, zerrors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
				eventstore: expectEventstore(),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "instance"),
			},
			res{
				err: zerrors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"push",
			fields{
				idpConfigEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: expectEventstore(
					expectPush(
						idpintent.NewLDAPSucceededEvent(
							context.Background(),
							&idpintent.NewAggregate("id", "instance").Aggregate,
							[]byte(`{"id":"id","preferredUsername":"username","preferredLanguage":"und"}`),
							"id",
							"username",
							"",
							map[string][]string{"id": {"id"}},
						),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "instance"),
				attributes: map[string][]string{"id": {"id"}},
				idpUser: ldap.NewUser(
					"id",
					"",
					"",
					"",
					"",
					"username",
					"",
					false,
					"",
					false,
					language.Tag{},
					"",
					"",
				),
			},
			res{
				token: "aWQ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.idpConfigEncryption,
			}
			got, err := c.SucceedLDAPIDPIntent(tt.args.ctx, tt.args.writeModel, tt.args.idpUser, tt.args.userID, tt.args.attributes)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.token, got)
		})
	}
}

func TestCommands_FailIDPIntent(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		writeModel *IDPIntentWriteModel
		reason     string
	}
	type res struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"push",
			fields{
				eventstore: expectEventstore(
					expectPush(
						idpintent.NewFailedEvent(
							context.Background(),
							&idpintent.NewAggregate("id", "instance").Aggregate,
							"reason",
						),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "instance"),
				reason:     "reason",
			},
			res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := c.FailIDPIntent(tt.args.ctx, tt.args.writeModel, tt.args.reason)
			require.ErrorIs(t, err, tt.res.err)
		})
	}
}

func Test_tokensForSucceededIDPIntent(t *testing.T) {
	type args struct {
		session       idp.Session
		encryptionAlg crypto.EncryptionAlgorithm
	}
	type res struct {
		accessToken *crypto.CryptoValue
		idToken     string
		err         error
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"no tokens",
			args{
				&ldap.Session{},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: nil,
				idToken:     "",
				err:         nil,
			},
		},
		{
			"token encryption fails",
			args{
				&oauth.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
					},
				},
				func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, zerrors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
			},
			res{
				accessToken: nil,
				idToken:     "",
				err:         zerrors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"oauth tokens",
			args{
				&oauth.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
					},
				},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("accessToken"),
				},
				idToken: "",
				err:     nil,
			},
		},
		{
			"oidc tokens",
			args{
				&openid.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
						IDToken: "idToken",
					},
				},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("accessToken"),
				},
				idToken: "idToken",
				err:     nil,
			},
		},
		{
			"jwt tokens",
			args{
				&jwt.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						IDToken: "idToken",
					},
				},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: nil,
				idToken:     "idToken",
				err:         nil,
			},
		},
		{
			"azure tokens",
			args{
				&azuread.Session{
					OAuthSession: &oauth.Session{
						Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
							Token: &oauth2.Token{
								AccessToken: "accessToken",
							},
							IDToken: "idToken",
						},
					},
				},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("accessToken"),
				},
				idToken: "idToken",
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccessToken, gotIDToken, err := tokensForSucceededIDPIntent(tt.args.session, tt.args.encryptionAlg)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.accessToken, gotAccessToken)
			assert.Equal(t, tt.res.idToken, gotIDToken)
		})
	}
}
