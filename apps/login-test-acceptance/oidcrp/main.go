package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

var (
	callbackPath = "/auth/callback"
	key          = []byte("test1234test1234")
)

func main() {
	apiURL := os.Getenv("API_URL")
	pat := readPAT(os.Getenv("PAT_FILE"))
	domain := os.Getenv("API_DOMAIN")
	loginURL := os.Getenv("LOGIN_URL")
	issuer := os.Getenv("ISSUER")
	port := os.Getenv("PORT")
	scopeList := strings.Split(os.Getenv("SCOPES"), " ")

	redirectURI := fmt.Sprintf("%s%s", issuer, callbackPath)
	cookieHandler := httphelper.NewCookieHandler(key, key, httphelper.WithUnsecure())

	clientID, clientSecret, err := createZitadelResources(apiURL, pat, domain, redirectURI, loginURL)
	if err != nil {
		panic(err)
	}

	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)
	client := &http.Client{
		Timeout: time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	// enable outgoing request logging
	logging.EnableHTTPClient(client,
		logging.WithClientGroup("client"),
	)

	options := []rp.Option{
		rp.WithCookieHandler(cookieHandler),
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
		rp.WithHTTPClient(client),
		rp.WithLogger(logger),
		rp.WithSigningAlgsFromDiscovery(),
		rp.WithCustomDiscoveryUrl(issuer + "/.well-known/openid-configuration"),
	}
	if clientSecret == "" {
		options = append(options, rp.WithPKCE(cookieHandler))
	}

	// One can add a logger to the context,
	// pre-defining log attributes as required.
	ctx := logging.ToContext(context.TODO(), logger)
	provider, err := rp.NewRelyingPartyOIDC(ctx, issuer, clientID, clientSecret, redirectURI, scopeList, options...)
	if err != nil {
		logrus.Fatalf("error creating provider %s", err.Error())
	}

	// generate some state (representing the state of the user in your application,
	// e.g. the page where he was before sending him to login
	state := func() string {
		return uuid.New().String()
	}

	urlOptions := []rp.URLParamOpt{
		rp.WithPromptURLParam("Welcome back!"),
	}

	// register the AuthURLHandler at your preferred path.
	// the AuthURLHandler creates the auth request and redirects the user to the auth server.
	// including state handling with secure cookie and the possibility to use PKCE.
	// Prompts can optionally be set to inform the server of
	// any messages that need to be prompted back to the user.
	http.Handle("/login", rp.AuthURLHandler(
		state,
		provider,
		urlOptions...,
	))

	// for demonstration purposes the returned userinfo response is written as JSON object onto response
	marshalUserinfo := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens[*oidc.IDTokenClaims], state string, rp rp.RelyingParty, info *oidc.UserInfo) {
		fmt.Println("access token", tokens.AccessToken)
		fmt.Println("refresh token", tokens.RefreshToken)
		fmt.Println("id token", tokens.IDToken)

		data, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(data)
	}

	// register the CodeExchangeHandler at the callbackPath
	// the CodeExchangeHandler handles the auth response, creates the token request and calls the callback function
	// with the returned tokens from the token endpoint
	// in this example the callback function itself is wrapped by the UserinfoCallback which
	// will call the Userinfo endpoint, check the sub and pass the info into the callback function
	http.Handle(callbackPath, rp.CodeExchangeHandler(rp.UserinfoCallback(marshalUserinfo), provider))

	// if you would use the callback without calling the userinfo endpoint, simply switch the callback handler for:
	//
	// http.Handle(callbackPath, rp.CodeExchangeHandler(marshalToken, provider))

	// simple counter for request IDs
	var counter atomic.Int64
	// enable incomming request logging
	mw := logging.Middleware(
		logging.WithLogger(logger),
		logging.WithGroup("server"),
		logging.WithIDFunc(func() slog.Attr {
			return slog.Int64("id", counter.Add(1))
		}),
	)

	http.Handle("/healthy", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { return }))
	fmt.Println("/healthy returns 200 OK")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mw(http.DefaultServeMux),
	}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}

func readPAT(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	pat, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(pat), "\n")
}

func createZitadelResources(apiURL, pat, domain, redirectURI, loginURL string) (string, string, error) {
	projectID, err := CreateProject(apiURL, pat, domain)
	if err != nil {
		return "", "", err
	}
	return CreateApp(apiURL, pat, domain, projectID, redirectURI, loginURL)
}

type project struct {
	ID string `json:"id"`
}
type createProject struct {
	Name                   string `json:"name"`
	ProjectRoleAssertion   bool   `json:"projectRoleAssertion"`
	ProjectRoleCheck       bool   `json:"projectRoleCheck"`
	HasProjectCheck        bool   `json:"hasProjectCheck"`
	PrivateLabelingSetting string `json:"privateLabelingSetting"`
}

func CreateProject(apiURL, pat, domain string) (string, error) {
	createProject := &createProject{
		Name:                   "OIDC",
		ProjectRoleAssertion:   false,
		ProjectRoleCheck:       false,
		HasProjectCheck:        false,
		PrivateLabelingSetting: "PRIVATE_LABELING_SETTING_UNSPECIFIED",
	}
	resp, err := doRequestWithHeaders(apiURL+"/management/v1/projects", pat, domain, createProject)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	p := new(project)
	if err := json.Unmarshal(data, p); err != nil {
		return "", err
	}
	fmt.Printf("projectID: %+v\n", p.ID)
	return p.ID, nil
}

type createApp struct {
	Name                     string   `json:"name"`
	RedirectUris             []string `json:"redirectUris"`
	ResponseTypes            []string `json:"responseTypes"`
	GrantTypes               []string `json:"grantTypes"`
	AppType                  string   `json:"appType"`
	AuthMethodType           string   `json:"authMethodType"`
	PostLogoutRedirectUris   []string `json:"postLogoutRedirectUris"`
	Version                  string   `json:"version"`
	DevMode                  bool     `json:"devMode"`
	AccessTokenType          string   `json:"accessTokenType"`
	AccessTokenRoleAssertion bool     `json:"accessTokenRoleAssertion"`
	IdTokenRoleAssertion     bool     `json:"idTokenRoleAssertion"`
	IdTokenUserinfoAssertion bool     `json:"idTokenUserinfoAssertion"`
	ClockSkew                string   `json:"clockSkew"`
	AdditionalOrigins        []string `json:"additionalOrigins"`
	SkipNativeAppSuccessPage bool     `json:"skipNativeAppSuccessPage"`
	BackChannelLogoutUri     []string `json:"backChannelLogoutUri"`
	LoginVersion             version  `json:"loginVersion"`
}

type version struct {
	LoginV2 loginV2 `json:"loginV2"`
}
type loginV2 struct {
	BaseUri string `json:"baseUri"`
}

type app struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func CreateApp(apiURL, pat, domain, projectID string, redirectURI, loginURL string) (string, string, error) {
	createApp := &createApp{
		Name:                     "OIDC",
		RedirectUris:             []string{redirectURI},
		ResponseTypes:            []string{"OIDC_RESPONSE_TYPE_CODE"},
		GrantTypes:               []string{"OIDC_GRANT_TYPE_AUTHORIZATION_CODE"},
		AppType:                  "OIDC_APP_TYPE_WEB",
		AuthMethodType:           "OIDC_AUTH_METHOD_TYPE_BASIC",
		Version:                  "OIDC_VERSION_1_0",
		DevMode:                  true,
		AccessTokenType:          "OIDC_TOKEN_TYPE_BEARER",
		AccessTokenRoleAssertion: true,
		IdTokenRoleAssertion:     true,
		IdTokenUserinfoAssertion: true,
		ClockSkew:                "1s",
		SkipNativeAppSuccessPage: true,
		LoginVersion: version{
			LoginV2: loginV2{
				BaseUri: loginURL,
			},
		},
	}

	resp, err := doRequestWithHeaders(apiURL+"/management/v1/projects/"+projectID+"/apps/oidc", pat, domain, createApp)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	a := new(app)
	if err := json.Unmarshal(data, a); err != nil {
		return "", "", err
	}
	return a.ClientID, a.ClientSecret, err
}

func doRequestWithHeaders(apiURL, pat, domain string, body any) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, err
	}
	values := http.Header{}
	values.Add("Authorization", "Bearer "+pat)
	values.Add("x-forwarded-host", domain)
	values.Add("Content-Type", "application/json")
	req.Header = values

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
