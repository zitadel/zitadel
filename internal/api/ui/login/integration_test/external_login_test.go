//go:build integration

// Package login_test contains an end-to-end reproduction of GHSA-738m-7888-jfv8, an
// external-IDP account pre-hijack / takeover in the server-rendered Login V1 UI.
//
// The vulnerability: POST /ui/login/externaluser/option (handleExternalNotFoundOptionCheck)
// is meant to be the *submit* of the "external account not found" page, which only ever appears
// after a genuine, protocol-verified IDP callback has populated the server-side
// authReq.LinkingUsers state. Before the fix, the handler ignored that state and built the new
// user's external identity (IDPConfigID + ExternalUserID) and its "verified" email straight from
// the raw POST body. An UNAUTHENTICATED attacker could therefore POST forged data directly to the
// endpoint - without ever going through the real IDP - and have Zitadel create a real account
// bound to an arbitrary (IDPConfigID, ExternalUserID). By setting ExternalUserID to a victim's
// public external identifier (e.g. a GitHub numeric user id) the attacker pre-creates an account
// that the victim's own, genuine future "Sign in with X" then silently logs into.
//
// This test drives the exact attack over raw HTTP against a running Zitadel instance:
//
//  1. authorize (create a V1 auth request)  ->  2. load login page (grab CSRF)
//     ->  3. select the IDP (sets SelectedIDPConfigID)  ->  4. forge the not-found-option POST
//
// It then asserts, over gRPC, that NO user was created for the attacker's forged identity. Against
// the unpatched code this fails (the forged account is created); against the fixed code it passes
// (the request is rejected with Errors.ExternalIDP.NoExternalUserData). That is the red -> green
// signal for the fix.
//
// Note on setup: no real GitHub account or working OAuth app is needed. The exploit bypasses the
// real OAuth round-trip entirely, so the IDP template only has to EXIST with creation allowed and
// dummy credentials - its client id/secret and endpoints are never exercised on the attack path.
//
// Requires a running backend + database. Run with:
//
//	go test -tags integration ./internal/api/ui/login/integration_test/...
package login_test

import (
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	feature "github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var (
	CTX      context.Context
	Instance *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		CTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		return m.Run()
	}())
}

// csrfTokenRegex extracts the gorilla/csrf hidden form field that the login templates embed via
// {{ .CSRF }} (rendered as <input type="hidden" name="gorilla.csrf.Token" value="...">). The token
// must be echoed back on every POST; the paired cookie travels automatically in the cookie jar.
var csrfTokenRegex = regexp.MustCompile(`name="gorilla\.csrf\.Token"[^>]*value="([^"]*)"`)

// TestExternalNotFoundOption_ForgedRegistration_IsRejected reproduces the core exploit and asserts
// the fix rejects it. It is the end-to-end, behavioral red -> green test for GHSA-738m-7888-jfv8.
func TestExternalNotFoundOption_ForgedRegistration_IsRejected(t *testing.T) {
	// --- Server setup (all via gRPC as the IAM owner) -------------------------------------------

	// Force Login V1. Three things select V1 at auth-request creation (internal/api/oidc/auth_request.go):
	// (a) the instance feature LoginV2.Required must be false, (b) the OIDC app's login version must be
	// unspecified/V1, and (c) the authorize request must NOT carry the x-zitadel-login-client header.
	// We pin (a) here explicitly so the test does not depend on the instance default.
	integration.EnsureInstanceFeature(t, CTX, Instance,
		&feature.SetInstanceFeaturesRequest{LoginV2: &feature.LoginV2{Required: false}},
		func(tt *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.False(tt, got.GetLoginV2().GetRequired())
		},
	)

	// A project + OIDC app. loginVersion is left nil (unspecified => V1). The redirect URI only has to
	// be registered; the attack never actually redeems a code there.
	redirectURI := "http://localhost:9999/callback"
	project := Instance.CreateProject(CTX, t, Instance.DefaultOrg.GetId(), integration.RandString(10), false, false)
	oidcApp, err := Instance.CreateOIDCClient(CTX, redirectURI, "", project.GetId(),
		app.OIDCAppType_OIDC_APP_TYPE_WEB, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, true)
	require.NoError(t, err)
	clientID := oidcApp.GetClientId()

	// An external OAuth IDP with creation ALLOWED but auto-creation DISABLED - the configuration under
	// which the genuine flow renders the not-found-option page. Its clientID/secret + example.com
	// endpoints are dummy and are never contacted on the attack path. auto-linking is left unspecified
	// so it plays no role here.
	idpResp := Instance.AddGenericOAuthProviderWithOptions(CTX, integration.RandString(10),
		true /*isLinkingAllowed*/, true /*isCreationAllowed*/, false /*isAutoCreation*/, idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_UNSPECIFIED)
	idpID := idpResp.GetId()
	// The IDP must be enabled on the login policy to be selectable in the login UI.
	Instance.AddProviderToDefaultLoginPolicy(CTX, idpID)

	// --- Attacker HTTP client -------------------------------------------------------------------

	// One shared cookie jar for the WHOLE flow. This is mandatory: the server issues a user-agent
	// cookie (zitadel.useragent) on the authorize response and binds the created auth request to that
	// user-agent id; getAuthRequest later looks it up by (authRequestID, userAgentID). A different jar
	// (or no jar) would make the auth request unresolvable. The csrf cookie also lives here.
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	httpClient := &http.Client{
		Jar: jar,
		// Never auto-follow redirects: we want to inspect each hop (and must NOT chase the IDP redirect
		// out to the dummy example.com authorize URL in step 3).
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	issuer := Instance.OIDCIssuer() // http://<instance-domain>:<port>

	// Step 1: create a Login V1 auth request by hitting the OIDC authorize endpoint WITHOUT the
	// login-client header. The response is a 302 to /ui/login/login?authRequestID=<id>; we extract id.
	authURL := issuer + "/oauth/v2/authorize?" + url.Values{
		"client_id":             {clientID},
		"redirect_uri":          {redirectURI},
		"response_type":         {"code"},
		"scope":                 {"openid"},
		"state":                 {"attacker-state"},
		"code_challenge":        {oidc.NewSHACodeChallenge(integration.CodeVerifier)},
		"code_challenge_method": {"S256"},
	}.Encode()
	loc := getLocation(t, httpClient, authURL)
	authRequestID := loc.Query().Get(login.QueryAuthRequestID)
	require.NotEmpty(t, authRequestID, "expected a V1 auth request id in the authorize redirect %q", loc.String())

	// Step 2: load the login page. This renders the login form and lets us scrape the CSRF token
	// (valid for any subsequent POST in this session, since it is bound to the csrf cookie in the jar).
	loginPageURL := issuer + login.HandlerPrefix + login.EndpointLogin + "?" + login.QueryAuthRequestID + "=" + authRequestID
	csrfToken := scrapeCSRF(t, httpClient, loginPageURL)

	// Step 3: select the external IDP. handleExternalLogin -> handleIDP persists SelectedIDPConfigID on
	// the auth request (SelectExternalIDP) and then 302-redirects toward the (dummy) IDP authorize URL,
	// which we deliberately do not follow. This is the ONLY legitimate step the attacker performs; it
	// requires no interaction with the real IDP.
	selectIDPURL := issuer + login.HandlerPrefix + login.EndpointExternalLogin + "?" +
		login.QueryAuthRequestID + "=" + authRequestID + "&idpConfigID=" + idpID
	resp := doGet(t, httpClient, selectIDPURL)
	resp.Body.Close()

	// Step 4: the attack. POST forged "external account not found" data. Every external-* field below is
	// attacker-controlled and, before the fix, was trusted verbatim to build the created user's identity.
	// The fix must ignore them (authReq.LinkingUsers is empty here because no genuine callback happened)
	// and reject the request.
	forgedEmail := "attacker-" + integration.RandString(8) + "@evil.test"
	forgedVictimID := "999999999" // e.g. a victim's public GitHub numeric user id
	form := url.Values{
		login.QueryAuthRequestID:   {authRequestID},
		"gorilla.csrf.Token":       {csrfToken},
		"external-idp-config-id":   {idpID},          // attacker-controlled: which IDP the link is bound to
		"external-idp-ext-user-id": {forgedVictimID}, // attacker-controlled: the identity being claimed
		"external-email":           {forgedEmail},    // attacker-controlled: echoed "IDP" email
		"external-email-verified":  {"true"},         // attacker-controlled: forged "verified" flag
		"email":                    {forgedEmail},    // attacker-controlled
		"username":                 {"attacker-" + integration.RandString(8)},
		"firstname":                {"At"},
		"lastname":                 {"Tacker"},
		"language":                 {"en"},
		"terms-confirm":            {"true"},
	}
	postURL := issuer + login.HandlerPrefix + login.EndpointExternalNotFoundOption + "?autoregisterbutton=true"
	body := doPostForm(t, httpClient, postURL, form)

	// The fixed handler rejects with the LinkingUsers guard (zerrors id LOGIN-Dju3f,
	// Errors.ExternalIDP.NoExternalUserData). We match the error code, which the error page renders
	// verbatim regardless of the negotiated UI locale (the translated message text is locale-dependent
	// and may fall back to the raw i18n key). This is an immediate, projection-lag-free signal.
	assert.Contains(t, body, "LOGIN-Dju3f",
		"forged registration should be rejected by the LinkingUsers guard; got a page without that rejection (attack likely succeeded)")

	// The load-bearing assertion: no account may exist for the attacker's forged identity. We use
	// assert.Never (not Eventually) so that if the account IS created but its projection lags, the user
	// still surfaces within the window and fails the test - which is exactly the unpatched behavior.
	assert.Never(t, func() bool {
		return userCountByEmail(t, forgedEmail) > 0
	}, 5*time.Second, 250*time.Millisecond,
		"a user was created for the attacker-forged external identity %q - account pre-hijack succeeded", forgedEmail)
}

// --- helpers ------------------------------------------------------------------------------------

func doGet(t *testing.T, c *http.Client, rawURL string) *http.Response {
	t.Helper()
	req, err := http.NewRequestWithContext(CTX, http.MethodGet, rawURL, nil)
	require.NoError(t, err)
	resp, err := c.Do(req)
	require.NoError(t, err)
	return resp
}

// getLocation issues a GET and returns the parsed Location of the expected redirect response.
func getLocation(t *testing.T, c *http.Client, rawURL string) *url.URL {
	t.Helper()
	resp := doGet(t, c, rawURL)
	defer resp.Body.Close()
	require.GreaterOrEqual(t, resp.StatusCode, 300)
	require.Less(t, resp.StatusCode, 400, "expected a redirect from %q", rawURL)
	loc, err := resp.Location()
	require.NoError(t, err)
	return loc
}

// scrapeCSRF GETs a rendered login page and returns its gorilla/csrf token.
func scrapeCSRF(t *testing.T, c *http.Client, rawURL string) string {
	t.Helper()
	resp := doGet(t, c, rawURL)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	m := csrfTokenRegex.FindStringSubmatch(string(b))
	require.Len(t, m, 2, "could not find gorilla.csrf.Token on %q", rawURL)
	return m[1]
}

// doPostForm submits a urlencoded form and returns the response body as a string.
func doPostForm(t *testing.T, c *http.Client, rawURL string, form url.Values) string {
	t.Helper()
	req, err := http.NewRequestWithContext(CTX, http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(b)
}

// userCountByEmail returns how many users match the given email address (via the User v2 API).
func userCountByEmail(t *testing.T, email string) int {
	t.Helper()
	resp, err := Instance.Client.UserV2.ListUsers(CTX, &user.ListUsersRequest{
		Queries: []*user.SearchQuery{
			{Query: &user.SearchQuery_EmailQuery{EmailQuery: &user.EmailQuery{EmailAddress: email}}},
		},
	})
	require.NoError(t, err)
	return len(resp.GetResult())
}
