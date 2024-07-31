//go:build integration

package idp_test

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/beevik/etree"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlidp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	saml_xml "github.com/zitadel/saml/pkg/provider/xml"
	"golang.org/x/crypto/bcrypt"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var (
	CTX    context.Context
	ErrCTX context.Context
	Tester *integration.Tester
	Client user.UserServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(time.Hour)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, ErrCTX = Tester.WithAuthorization(ctx, integration.OrgOwner), errCtx
		Client = Tester.Client.UserV2
		return m.Run()
	}())
}

func TestServer_SAMLCertificate(t *testing.T) {
	samlRedirectIdpID := Tester.AddSAMLRedirectProvider(t, CTX, "")
	oauthIdpID := Tester.AddGenericOAuthProvider(t, CTX)

	type args struct {
		ctx   context.Context
		idpID string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "saml certificate, invalid idp",
			args: args{
				ctx:   CTX,
				idpID: "unknown",
			},
			want: http.StatusNotFound,
		},
		{
			name: "saml certificate, invalid idp type",
			args: args{
				ctx:   CTX,
				idpID: oauthIdpID,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "saml certificate, ok",
			args: args{
				ctx:   CTX,
				idpID: samlRedirectIdpID,
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			certificateURL := http_util.BuildOrigin(Tester.Host(), Tester.Server.Config.ExternalSecure) + "/idps/" + tt.args.idpID + "/saml/certificate"
			resp, err := http.Get(certificateURL)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, resp.StatusCode)
			if tt.want == http.StatusOK {
				b, err := io.ReadAll(resp.Body)
				defer resp.Body.Close()
				assert.NoError(t, err)

				block, _ := pem.Decode(b)
				_, err = x509.ParseCertificate(block.Bytes)
				assert.NoError(t, err)
			}
		})
	}
}

func TestServer_SAMLMetadata(t *testing.T) {
	samlRedirectIdpID := Tester.AddSAMLRedirectProvider(t, CTX, "")
	oauthIdpID := Tester.AddGenericOAuthProvider(t, CTX)

	type args struct {
		ctx   context.Context
		idpID string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "saml metadata, invalid idp",
			args: args{
				ctx:   CTX,
				idpID: "unknown",
			},
			want: http.StatusNotFound,
		},
		{
			name: "saml metadata, invalid idp type",
			args: args{
				ctx:   CTX,
				idpID: oauthIdpID,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "saml metadata, ok",
			args: args{
				ctx:   CTX,
				idpID: samlRedirectIdpID,
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadataURL := http_util.BuildOrigin(Tester.Host(), Tester.Server.Config.ExternalSecure) + "/idps/" + tt.args.idpID + "/saml/metadata"
			resp, err := http.Get(metadataURL)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, resp.StatusCode)
			if tt.want == http.StatusOK {
				b, err := io.ReadAll(resp.Body)
				defer resp.Body.Close()
				assert.NoError(t, err)

				_, err = saml_xml.ParseMetadataXmlIntoStruct(b)
				assert.NoError(t, err)
			}

		})
	}
}

func TestServer_SAMLACS(t *testing.T) {
	userHuman := Tester.CreateHumanUser(CTX)
	samlRedirectIdpID := Tester.AddSAMLRedirectProvider(t, CTX, "urn:oid:0.9.2342.19200300.100.1.1") // the username is set in urn:oid:0.9.2342.19200300.100.1.1
	externalUserID := "test1"
	linkedExternalUserID := "test2"
	Tester.CreateUserIDPlink(CTX, userHuman.UserId, linkedExternalUserID, samlRedirectIdpID, linkedExternalUserID)
	idp, err := getIDP(
		http_util.BuildOrigin(Tester.Host(), Tester.Server.Config.ExternalSecure),
		[]string{samlRedirectIdpID},
		externalUserID,
		linkedExternalUserID,
	)
	assert.NoError(t, err)

	type args struct {
		ctx          context.Context
		successURL   string
		failureURL   string
		idpID        string
		username     string
		nameID       string
		nameIDFormat string
		intentID     string
		response     string
	}
	type want struct {
		successful bool
		user       string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "intent invalid",
			args: args{
				ctx:          CTX,
				successURL:   "https://example.com/success",
				failureURL:   "https://example.com/failure",
				idpID:        samlRedirectIdpID,
				username:     externalUserID,
				nameID:       externalUserID,
				nameIDFormat: string(saml.PersistentNameIDFormat),
				intentID:     "notexisting",
			},
			want: want{
				successful: false,
				user:       "",
			},
			wantErr: true,
		},
		{
			name: "response invalid",
			args: args{
				ctx:          CTX,
				successURL:   "https://example.com/success",
				failureURL:   "https://example.com/failure",
				idpID:        samlRedirectIdpID,
				username:     externalUserID,
				nameID:       externalUserID,
				nameIDFormat: string(saml.PersistentNameIDFormat),
				response:     "invalid",
			},
			want: want{
				successful: false,
				user:       "",
			},
		},
		{
			name: "saml flow redirect, ok",
			args: args{
				ctx:          CTX,
				successURL:   "https://example.com/success",
				failureURL:   "https://example.com/failure",
				idpID:        samlRedirectIdpID,
				username:     externalUserID,
				nameID:       externalUserID,
				nameIDFormat: string(saml.PersistentNameIDFormat),
			},
			want: want{
				successful: true,
				user:       "",
			},
		},
		{
			name: "saml flow redirect with link, ok",
			args: args{
				ctx:          CTX,
				successURL:   "https://example.com/success",
				failureURL:   "https://example.com/failure",
				idpID:        samlRedirectIdpID,
				username:     linkedExternalUserID,
				nameID:       linkedExternalUserID,
				nameIDFormat: string(saml.PersistentNameIDFormat),
			},
			want: want{
				successful: true,
				user:       userHuman.UserId,
			},
		},
		{
			name: "saml flow redirect (transient), ok",
			args: args{
				ctx:          CTX,
				successURL:   "https://example.com/success",
				failureURL:   "https://example.com/failure",
				idpID:        samlRedirectIdpID,
				username:     externalUserID,
				nameID:       "genericID",
				nameIDFormat: string(saml.TransientNameIDFormat),
			},
			want: want{
				successful: true,
				user:       "",
			},
		},
		{
			name: "saml flow redirect with link (transient), ok",
			args: args{
				ctx:          CTX,
				successURL:   "https://example.com/success",
				failureURL:   "https://example.com/failure",
				idpID:        samlRedirectIdpID,
				username:     linkedExternalUserID,
				nameID:       "genericID",
				nameIDFormat: string(saml.TransientNameIDFormat),
			},
			want: want{
				successful: true,
				user:       userHuman.UserId,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.StartIdentityProviderIntent(tt.args.ctx,
				&user.StartIdentityProviderIntentRequest{
					IdpId: tt.args.idpID,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: tt.args.successURL,
							FailureUrl: tt.args.failureURL,
						},
					},
				},
			)
			// can't fail as covered in other tests
			require.NoError(t, err)

			//parse returned URL to continue flow to callback with the same intentID==RelayState
			authURL, err := url.Parse(got.GetAuthUrl())
			require.NoError(t, err)
			samlRequest := &http.Request{Method: http.MethodGet, URL: authURL}
			assert.NotEmpty(t, authURL)

			//generate necessary information to create request to callback URL
			relayState := authURL.Query().Get("RelayState")
			//test purposes, use defined intentID
			if tt.args.intentID != "" {
				relayState = tt.args.intentID
			}
			callbackURL := http_util.BuildOrigin(Tester.Host(), Tester.Server.Config.ExternalSecure) + "/idps/" + tt.args.idpID + "/saml/acs"
			response := createResponse(t, idp, samlRequest, tt.args.nameID, tt.args.nameIDFormat, tt.args.username)
			//test purposes, use defined response
			if tt.args.response != "" {
				response = tt.args.response
			}
			location, err := integration.CheckPost(callbackURL, httpPostFormRequest(relayState, response))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, relayState, location.Query().Get("id"))
				if tt.want.successful {
					assert.True(t, strings.HasPrefix(location.String(), tt.args.successURL))
					assert.NotEmpty(t, location.Query().Get("token"))
					assert.Equal(t, tt.want.user, location.Query().Get("user"))
				} else {
					assert.True(t, strings.HasPrefix(location.String(), tt.args.failureURL))
				}
			}
		})
	}
}

var key = func() crypto.PrivateKey {
	b, _ := pem.Decode([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0OhbMuizgtbFOfwbK7aURuXhZx6VRuAs3nNibiuifwCGz6u9
yy7bOR0P+zqN0YkjxaokqFgra7rXKCdeABmoLqCC0U+cGmLNwPOOA0PaD5q5xKhQ
4Me3rt/R9C4Ca6k3/OnkxnKwnogcsmdgs2l8liT3qVHP04Oc7Uymq2v09bGb6nPu
fOrkXS9F6mSClxHG/q59AGOWsXK1xzIRV1eu8W2SNdyeFVU1JHiQe444xLoPul5t
InWasKayFsPlJfWNc8EoU8COjNhfo/GovFTHVjh9oUR/gwEFVwifIHihRE0Hazn2
EQSLaOr2LM0TsRsQroFjmwSGgI+X2bfbMTqWOQIDAQABAoIBAFWZwDTeESBdrLcT
zHZe++cJLxE4AObn2LrWANEv5AeySYsyzjRBYObIN9IzrgTb8uJ900N/zVr5VkxH
xUa5PKbOcowd2NMfBTw5EEnaNbILLm+coHdanrNzVu59I9TFpAFoPavrNt/e2hNo
NMGPSdOkFi81LLl4xoadz/WR6O/7N2famM+0u7C2uBe+TrVwHyuqboYoidJDhO8M
w4WlY9QgAUhkPyzZqrl+VfF1aDTGVf4LJgaVevfFCas8Ws6DQX5q4QdIoV6/0vXi
B1M+aTnWjHuiIzjBMWhcYW2+I5zfwNWRXaxdlrYXRukGSdnyO+DH/FhHePJgmlkj
NInADDkCgYEA6MEQFOFSCc/ELXYWgStsrtIlJUcsLdLBsy1ocyQa2lkVUw58TouW
RciE6TjW9rp31pfQUnO2l6zOUC6LT9Jvlb9PSsyW+rvjtKB5PjJI6W0hjX41wEO6
fshFELMJd9W+Ezao2AsP2hZJ8McCF8no9e00+G4xTAyxHsNI2AFTCQcCgYEA5cWZ
JwNb4t7YeEajPt9xuYNUOQpjvQn1aGOV7KcwTx5ELP/Hzi723BxHs7GSdrLkkDmi
Gpb+mfL4wxCt0fK0i8GFQsRn5eusyq9hLqP/bmjpHoXe/1uajFbE1fZQR+2LX05N
3ATlKaH2hdfCJedFa4wf43+cl6Yhp6ZA0Yet1r8CgYEAwiu1j8W9G+RRA5/8/DtO
yrUTOfsbFws4fpLGDTA0mq0whf6Soy/96C90+d9qLaC3srUpnG9eB0CpSOjbXXbv
kdxseLkexwOR3bD2FHX8r4dUM2bzznZyEaxfOaQypN8SV5ME3l60Fbr8ajqLO288
wlTmGM5Mn+YCqOg/T7wjGmcCgYBpzNfdl/VafOROVbBbhgXWtzsz3K3aYNiIjbp+
MunStIwN8GUvcn6nEbqOaoiXcX4/TtpuxfJMLw4OvAJdtxUdeSmEee2heCijV6g3
ErrOOy6EqH3rNWHvlxChuP50cFQJuYOueO6QggyCyruSOnDDuc0BM0SGq6+5g5s7
H++S/wKBgQDIkqBtFr9UEf8d6JpkxS0RXDlhSMjkXmkQeKGFzdoJcYVFIwq8jTNB
nJrVIGs3GcBkqGic+i7rTO1YPkquv4dUuiIn+vKZVoO6b54f+oPBXd4S0BnuEqFE
rdKNuCZhiaE2XD9L/O9KP1fh5bfEcKwazQ23EvpJHBMm8BGC+/YZNw==
-----END RSA PRIVATE KEY-----`))
	k, _ := x509.ParsePKCS1PrivateKey(b.Bytes)
	return k
}()

var cert = func() *x509.Certificate {
	b, _ := pem.Decode([]byte(`-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV
BAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5
NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8A
hs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+a
ucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWx
m+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6
D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURN
B2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0O
BBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56
zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5
pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uv
NONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEf
y/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL
/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsb
GFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTL
UzreO96WzlBBMtY=
-----END CERTIFICATE-----`))
	c, _ := x509.ParseCertificate(b.Bytes)
	return c
}()

func getIDP(zitadelBaseURL string, idpIDs []string, user1, user2 string) (*saml.IdentityProvider, error) {
	baseURL, err := url.Parse("http://localhost:8000")
	if err != nil {
		return nil, err
	}

	store := &samlidp.MemoryStore{}
	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	err = store.Put("/users/"+user1, samlidp.User{
		Name:           user1,
		HashedPassword: hashedPassword1,
		Groups:         []string{"Administrators", "Users"},
		Email:          "test@example.com",
		CommonName:     "Test Test",
		Surname:        "Test",
		GivenName:      "Test",
	})
	if err != nil {
		return nil, err
	}
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	err = store.Put("/users/"+user2, samlidp.User{
		Name:           user2,
		HashedPassword: hashedPassword2,
		Groups:         []string{"Administrators", "Users"},
		Email:          "test@example.com",
		CommonName:     "Test Test",
		Surname:        "Test",
		GivenName:      "Test",
	})
	if err != nil {
		return nil, err
	}
	for _, idpID := range idpIDs {
		metadata, err := saml_xml.ReadMetadataFromURL(http.DefaultClient, zitadelBaseURL+"/idps/"+idpID+"/saml/metadata")
		if err != nil {
			return nil, err
		}
		entity := new(saml.EntityDescriptor)
		if err := xml.Unmarshal(metadata, entity); err != nil {
			return nil, err
		}

		if err := store.Put("/services/"+idpID, samlidp.Service{
			Name:     idpID,
			Metadata: *entity,
		}); err != nil {
			return nil, err
		}
	}

	idpServer, err := samlidp.New(samlidp.Options{
		URL:         *baseURL,
		Key:         key,
		Certificate: cert,
		Store:       store,
	})
	if err != nil {
		return nil, err
	}
	if idpServer.IDP.AssertionMaker == nil {
		idpServer.IDP.AssertionMaker = &saml.DefaultAssertionMaker{}
	}
	return &idpServer.IDP, nil
}

func createResponse(t *testing.T, idp *saml.IdentityProvider, req *http.Request, nameID, nameIDFormat, username string) string {
	authnReq, err := saml.NewIdpAuthnRequest(idp, req)
	assert.NoError(t, authnReq.Validate())

	err = idp.AssertionMaker.MakeAssertion(authnReq, &saml.Session{
		CreateTime:   time.Now().UTC(),
		Index:        "",
		NameID:       nameID,
		NameIDFormat: nameIDFormat,
		UserName:     username,
	})
	assert.NoError(t, err)
	err = authnReq.MakeResponse()
	assert.NoError(t, err)

	doc := etree.NewDocument()
	doc.SetRoot(authnReq.ResponseEl)
	responseBuf, err := doc.WriteToBytes()
	assert.NoError(t, err)
	responseBuf = append([]byte("<?xml version=\"1.0\"?>"), responseBuf...)

	return base64.StdEncoding.EncodeToString(responseBuf)
}

func httpGETRequest(t *testing.T, callbackURL string, relayState, response, sig, sigAlg string) *http.Request {
	req, err := http.NewRequest("GET", callbackURL, nil)
	require.NoError(t, err)

	q := req.URL.Query()
	q.Add("RelayState", relayState)
	q.Add("SAMLResponse", response)
	if sig != "" {
		q.Add("Sig", sig)
	}
	if sigAlg != "" {
		q.Add("SigAlg", sigAlg)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

func httpPostFormRequest(relayState, response string) url.Values {
	return url.Values{
		"SAMLResponse": {response},
		"RelayState":   {relayState},
	}
}
