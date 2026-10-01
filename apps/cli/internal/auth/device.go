package auth

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

// Login performs the OIDC Device Authorization Grant flow interactively.
// It initiates the device authorization request, displays the user code and
// verification URI, and polls the token endpoint until the user approves or
// the context is cancelled.
func Login(ctx context.Context, instance, clientID string) (*oauth2.Token, error) {
	// Normalise instance to an issuer URL
	issuer := instance
	if !strings.HasPrefix(issuer, "http://") && !strings.HasPrefix(issuer, "https://") {
		issuer = "https://" + issuer
	}

	scopes := []string{
		oidc.ScopeOpenID,
		oidc.ScopeProfile,
		oidc.ScopeEmail,
		oidc.ScopeOfflineAccess,
		"urn:zitadel:iam:org:project:id:zitadel:aud",
	}

	provider, err := rp.NewRelyingPartyOIDC(ctx, issuer, clientID, "", "", scopes)
	if err != nil {
		return nil, fmt.Errorf("creating OIDC relying party: %w", err)
	}

	// 1. Start Device Authorization
	resp, err := rp.DeviceAuthorization(ctx, scopes, provider, nil)
	if err != nil {
		return nil, fmt.Errorf("device authorization request failed: %w", err)
	}

	// 2. Display verification info to user
	fmt.Fprintf(os.Stderr, "\nTo authenticate, visit:\n  %s\n", resp.VerificationURI)
	if resp.VerificationURIComplete != "" {
		fmt.Fprintf(os.Stderr, "Or use the direct link: %s\n", resp.VerificationURIComplete)
	}
	fmt.Fprintf(os.Stderr, "\nAnd enter the code: %s\n", resp.UserCode)
	fmt.Fprintln(os.Stderr, "\nWaiting for authorization...")

	// 3. Poll for the access token
	// rp.DeviceAccessToken handles the interval and polling logic automatically.
	tokenResponse, err := rp.DeviceAccessToken(ctx, resp.DeviceCode, time.Duration(resp.Interval)*time.Second, provider)
	if err != nil {
		return nil, fmt.Errorf("polling for device access token failed: %w", err)
	}

	// Convert *oidc.AccessTokenResponse to *oauth2.Token
	token := &oauth2.Token{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		TokenType:    tokenResponse.TokenType,
		Expiry:       time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
	}

	return token, nil
}
