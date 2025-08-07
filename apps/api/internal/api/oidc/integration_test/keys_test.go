//go:build integration

package oidc_test

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/integration"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

func TestServer_Keys(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ctxLogin := instance.WithAuthorization(CTX, integration.UserTypeLogin)

	clientID, _ := createClient(t, instance)
	authRequestID := createAuthRequest(t, instance, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)

	instance.RegisterUserPasskey(instance.WithAuthorization(CTX, integration.UserTypeOrgOwner), instance.AdminUserID)
	sessionID, sessionToken, _, _ := instance.CreateVerifiedWebAuthNSession(t, ctxLogin, instance.AdminUserID)
	linkResp, err := instance.Client.OIDCv2.CreateCallback(ctxLogin, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// code exchange so we are sure there is 1 legacy key pair.
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	_, err = exchangeTokens(t, instance, clientID, code, redirectURI)
	require.NoError(t, err)

	issuer := http_util.BuildHTTP(instance.Domain, instance.Config.Port, instance.Config.Secure)
	discovery, err := client.Discover(CTX, issuer, http.DefaultClient)
	require.NoError(t, err)

	tests := []struct {
		name    string
		wantLen int
	}{
		{
			name:    "webkeys",
			wantLen: 2, // 2 from instance creation.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EventuallyWithT(t, func(ttt *assert.CollectT) {
				resp, err := http.Get(discovery.JwksURI)
				require.NoError(ttt, err)
				require.Equal(ttt, resp.StatusCode, http.StatusOK)
				defer resp.Body.Close()

				got := new(jose.JSONWebKeySet)
				err = json.NewDecoder(resp.Body).Decode(got)
				require.NoError(ttt, err)

				assert.Len(t, got.Keys, tt.wantLen)
				for _, key := range got.Keys {
					_, ok := key.Key.(*rsa.PublicKey)
					require.True(ttt, ok)
					require.NotEmpty(ttt, key.KeyID)
					require.Equal(ttt, key.Algorithm, string(jose.RS256))
					require.Equal(ttt, key.Use, crypto.KeyUsageSigning.String())
				}

				cacheControl := resp.Header.Get("cache-control")
				require.Equal(ttt, "max-age=300, must-revalidate", cacheControl)

			}, time.Minute, time.Second/10)
		})

	}
}
