package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

// Login performs the OIDC PKCE authorization code flow interactively.
// It opens the user's browser, waits for the callback on a localhost port,
// exchanges the code, and returns the resulting OAuth2 token.
func Login(ctx context.Context, instance, clientID string) (*oauth2.Token, error) {
	// Normalise instance to an issuer URL
	issuer := instance
	if !strings.HasPrefix(issuer, "http://") && !strings.HasPrefix(issuer, "https://") {
		issuer = "https://" + issuer
	}

	// Listen on a random free port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("starting local server: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", port)

	// PKCE cookie handler (keys don't matter — single-use local server)
	key := []byte("zitadelcli123456") // 16 bytes for AES-128
	cookieHandler := httphelper.NewCookieHandler(key, key, httphelper.WithUnsecure())

	scopes := []string{
		oidc.ScopeOpenID,
		oidc.ScopeProfile,
		oidc.ScopeEmail,
		oidc.ScopeOfflineAccess,
		"urn:zitadel:iam:org:project:id:zitadel:aud",
	}

	provider, err := rp.NewRelyingPartyOIDC(ctx, issuer, clientID, "", redirectURI, scopes,
		rp.WithPKCE(cookieHandler),
		rp.WithCookieHandler(cookieHandler),
	)
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("creating OIDC relying party: %w", err)
	}

	tokenCh := make(chan *oauth2.Token, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()

	// Auth start handler
	mux.Handle("/login", rp.AuthURLHandler(func() string { return "cli-login" }, provider))

	// Callback handler — extracts tokens after code exchange
	mux.Handle("/callback", rp.CodeExchangeHandler(func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, provider rp.RelyingParty) {
		tokenCh <- tokens.Token
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<html><body><h1>Login successful!</h1><p>You can close this window.</p></body></html>")
	}, provider))

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Open the browser at the /login path which triggers the auth URL redirect
	loginURL := fmt.Sprintf("http://localhost:%d/login", port)
	if err := openBrowser(loginURL); err != nil {
		fmt.Fprintf(os.Stderr, "Open this URL in your browser:\n  %s\n", loginURL)
	}

	fmt.Fprintln(os.Stderr, "Waiting for login in browser...")

	// Wait for the token or an error
	select {
	case token := <-tokenCh:
		server.Shutdown(ctx)
		return token, nil
	case err := <-errCh:
		server.Shutdown(ctx)
		return nil, err
	case <-ctx.Done():
		server.Shutdown(context.Background())
		return nil, ctx.Err()
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
