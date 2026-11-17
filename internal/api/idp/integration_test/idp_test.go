//go:build integration

package idp_test

import (
	"context"
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
	"github.com/zitadel/saml/pkg/provider/xml/md"
	"golang.org/x/crypto/bcrypt"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var (
	CTX      context.Context
	Instance *integration.Instance
	Client   user.UserServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		Client = Instance.Client.UserV2
		return m.Run()
	}())
}

func TestServer_SAMLCertificate(t *testing.T) {
	samlRedirectIdpID := Instance.AddSAMLRedirectProvider(CTX, "")
	oauthIdpResp := Instance.AddGenericOAuthProvider(CTX, Instance.DefaultOrg.Id)

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
				idpID: oauthIdpResp.Id,
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
			certificateURL := http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/idps/" + tt.args.idpID + "/saml/certificate"
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
	samlRedirectIdpID := Instance.AddSAMLRedirectProvider(CTX, "")
	oauthIdpResp := Instance.AddGenericOAuthProvider(CTX, Instance.DefaultOrg.Id)

	type args struct {
		ctx        context.Context
		idpID      string
		internalUI bool
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantACS []md.IndexedEndpointType
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
				idpID: oauthIdpResp.Id,
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
			wantACS: []md.IndexedEndpointType{
				{
					XMLName: xml.Name{
						Space: "urn:oasis:names:tc:SAML:2.0:metadata",
						Local: "AssertionConsumerService",
					},
					Index:            "1",
					IsDefault:        "",
					Binding:          saml.HTTPPostBinding,
					Location:         http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/idps/" + samlRedirectIdpID + "/saml/acs",
					ResponseLocation: "",
				},
				{
					XMLName: xml.Name{
						Space: "urn:oasis:names:tc:SAML:2.0:metadata",
						Local: "AssertionConsumerService",
					},
					Index:            "2",
					IsDefault:        "",
					Binding:          saml.HTTPArtifactBinding,
					Location:         http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/idps/" + samlRedirectIdpID + "/saml/acs",
					ResponseLocation: "",
				},
				{
					XMLName: xml.Name{
						Space: "urn:oasis:names:tc:SAML:2.0:metadata",
						Local: "AssertionConsumerService",
					},
					Index:            "3",
					IsDefault:        "",
					Binding:          saml.HTTPPostBinding,
					Location:         http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/ui/login/login/externalidp/saml/acs",
					ResponseLocation: "",
				},
				{
					XMLName: xml.Name{
						Space: "urn:oasis:names:tc:SAML:2.0:metadata",
						Local: "AssertionConsumerService",
					},
					Index:            "4",
					IsDefault:        "",
					Binding:          saml.HTTPArtifactBinding,
					Location:         http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/ui/login/login/externalidp/saml/acs",
					ResponseLocation: "",
				},
			},
		},
		{
			name: "saml metadata, ok (internalUI)",
			args: args{
				ctx:        CTX,
				idpID:      samlRedirectIdpID,
				internalUI: true,
			},
			want: http.StatusOK,
			wantACS: []md.IndexedEndpointType{
				{
					XMLName: xml.Name{
						Space: "urn:oasis:names:tc:SAML:2.0:metadata",
						Local: "AssertionConsumerService",
					},
					Index:            "0",
					IsDefault:        "true",
					Binding:          saml.HTTPPostBinding,
					Location:         http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/ui/login/login/externalidp/saml/acs",
					ResponseLocation: "",
				},
				{
					XMLName: xml.Name{
						Space: "urn:oasis:names:tc:SAML:2.0:metadata",
						Local: "AssertionConsumerService",
					},
					Index:            "1",
					IsDefault:        "",
					Binding:          saml.HTTPArtifactBinding,
					Location:         http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/ui/login/login/externalidp/saml/acs",
					ResponseLocation: "",
				},
				{
					XMLName: xml.Name{
						Space: "urn:oasis:names:tc:SAML:2.0:metadata",
						Local: "AssertionConsumerService",
					},
					Index:            "2",
					IsDefault:        "",
					Binding:          saml.HTTPPostBinding,
					Location:         http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/idps/" + samlRedirectIdpID + "/saml/acs",
					ResponseLocation: "",
				},
				{
					XMLName: xml.Name{
						Space: "urn:oasis:names:tc:SAML:2.0:metadata",
						Local: "AssertionConsumerService",
					},
					Index:            "3",
					IsDefault:        "",
					Binding:          saml.HTTPArtifactBinding,
					Location:         http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/idps/" + samlRedirectIdpID + "/saml/acs",
					ResponseLocation: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadataURL := http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/idps/" + tt.args.idpID + "/saml/metadata"
			if tt.args.internalUI {
				metadataURL = metadataURL + "?internalUI=true"
			}
			resp, err := http.Get(metadataURL)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, resp.StatusCode)
			if tt.want == http.StatusOK {
				b, err := io.ReadAll(resp.Body)
				defer resp.Body.Close()
				assert.NoError(t, err)

				metadata, err := saml_xml.ParseMetadataXmlIntoStruct(b)
				assert.NoError(t, err)

				assert.Equal(t, metadata.SPSSODescriptor.AssertionConsumerService, tt.wantACS)
			}
		})
	}
}

func TestServer_SAMLACS(t *testing.T) {
	userHuman := Instance.CreateHumanUser(CTX)
	samlRedirectIdpID := Instance.AddSAMLRedirectProvider(CTX, "urn:oid:0.9.2342.19200300.100.1.1") // the username is set in urn:oid:0.9.2342.19200300.100.1.1
	externalUserID := "test1"
	linkedExternalUserID := "test2"
	Instance.CreateUserIDPlink(CTX, userHuman.UserId, linkedExternalUserID, samlRedirectIdpID, linkedExternalUserID)
	idp, err := getIDP(
		http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure),
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
			callbackURL := http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure) + "/idps/" + tt.args.idpID + "/saml/acs"
			response := createResponse(t, idp, samlRequest, tt.args.nameID, tt.args.nameIDFormat, tt.args.username)
			//test purposes, use defined response
			if tt.args.response != "" {
				response = tt.args.response
			}
			location, err := integration.CheckPost(callbackURL, httpPostFormRequest(relayState, response))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, relayState, location.Query().Get("id"))
			if tt.want.successful {
				assert.True(t, strings.HasPrefix(location.String(), tt.args.successURL))
				assert.NotEmpty(t, location.Query().Get("token"))
				assert.Equal(t, tt.want.user, location.Query().Get("user"))
			} else {
				assert.True(t, strings.HasPrefix(location.String(), tt.args.failureURL))
			}

		})
	}
}

// TestServer_SAMLACSIntentIDPMismatch reproduces the SAML IdP-confusion issue: an assertion validly
// signed by and targeted at IdP-A must not be accepted for a login intent that was started for a
// different IdP-B. Before the fix the ACS handler verified the assertion against the IdP named in the
// request path but resolved the user under the intent's IdP, so the forged response below succeeded
// and returned a session as the victim. It must now be rejected.
func TestServer_SAMLACSIntentIDPMismatch(t *testing.T) {
	victim := Instance.CreateHumanUser(CTX)
	externalUserID := "victim-external-id"

	// IdP-B is the victim's IdP (the intent target); IdP-A is the attacker's IdP (whose signing key the
	// attacker controls). Both are backed by the same instance SAML key here, which is irrelevant to the
	// attack — it turns on Destination/Audience/InResponseTo, not on a distinct key.
	idpB := Instance.AddSAMLRedirectProvider(CTX, "")
	idpA := Instance.AddSAMLRedirectProvider(CTX, "")
	Instance.CreateUserIDPlink(CTX, victim.UserId, externalUserID, idpB, externalUserID)

	baseURL := http_util.BuildOrigin(Instance.Host(), Instance.Config.Secure)
	idp, err := getIDP(baseURL, []string{idpA, idpB}, externalUserID, externalUserID+"-2")
	require.NoError(t, err)

	const (
		successURL = "https://example.com/success"
		failureURL = "https://example.com/failure"
	)

	// Start a genuine intent for IdP-B and keep its RelayState; this is the intent the attacker confuses.
	intentB, err := Client.StartIdentityProviderIntent(CTX, &user.StartIdentityProviderIntentRequest{
		IdpId:   idpB,
		Content: &user.StartIdentityProviderIntentRequest_Urls{Urls: &user.RedirectURLs{SuccessUrl: successURL, FailureUrl: failureURL}},
	})
	require.NoError(t, err)
	authURLB, err := url.Parse(intentB.GetAuthUrl())
	require.NoError(t, err)
	relayStateB := authURLB.Query().Get("RelayState")
	reqB := &http.Request{Method: http.MethodGet, URL: authURLB}

	// Start an intent for IdP-A only to obtain a valid AuthnRequest targeting IdP-A's SP.
	intentA, err := Client.StartIdentityProviderIntent(CTX, &user.StartIdentityProviderIntentRequest{
		IdpId:   idpA,
		Content: &user.StartIdentityProviderIntentRequest_Urls{Urls: &user.RedirectURLs{SuccessUrl: successURL, FailureUrl: failureURL}},
	})
	require.NoError(t, err)
	authURLA, err := url.Parse(intentA.GetAuthUrl())
	require.NoError(t, err)
	reqA := &http.Request{Method: http.MethodGet, URL: authURLA}

	// Forge a response validly signed for IdP-A (Destination/Audience = IdP-A's SP), but whose
	// InResponseTo is copied from the IdP-B intent so it satisfies that intent's RequestID check.
	response := createConfusedResponse(t, idp, reqA, reqB, externalUserID, string(saml.PersistentNameIDFormat))

	// Post it to IdP-A's ACS with the IdP-B intent's RelayState: the URL names A, the RelayState names B.
	callbackURL := baseURL + "/idps/" + idpA + "/saml/acs"
	_, err = integration.CheckPost(callbackURL, httpPostFormRequest(relayStateB, response))
	// Must be rejected. Before the fix this returned a 302 to the success URL carrying the victim's
	// session, so CheckPost would return no error; now the mismatch yields a 400 (no Location).
	require.Error(t, err)
}

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
		Key:         Instance.SAMLPrivateKey(),
		Certificate: Instance.SAMLCertificate(),
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

// createConfusedResponse builds a SAML response that is validly signed for and targeted at the SP of
// spReq's IdP, but whose InResponseTo is taken from idReq (a different IdP's intent). This mirrors the
// IdP-confusion attack: the assertion verifies against one IdP while answering another IdP's intent.
func createConfusedResponse(t *testing.T, idp *saml.IdentityProvider, spReq, idReq *http.Request, nameID, nameIDFormat string) string {
	sp, err := saml.NewIdpAuthnRequest(idp, spReq)
	require.NoError(t, err)
	require.NoError(t, sp.Validate())

	id, err := saml.NewIdpAuthnRequest(idp, idReq)
	require.NoError(t, err)
	require.NoError(t, id.Validate())

	// Keep spReq's SP (Destination/Recipient/Audience), but answer idReq's request ID so the response
	// passes the RequestID (InResponseTo) check of the confused intent.
	sp.Request.ID = id.Request.ID

	err = idp.AssertionMaker.MakeAssertion(sp, &saml.Session{
		CreateTime:   time.Now().UTC(),
		NameID:       nameID,
		NameIDFormat: nameIDFormat,
		UserName:     nameID,
	})
	require.NoError(t, err)
	require.NoError(t, sp.MakeResponse())

	doc := etree.NewDocument()
	doc.SetRoot(sp.ResponseEl)
	responseBuf, err := doc.WriteToBytes()
	require.NoError(t, err)
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
