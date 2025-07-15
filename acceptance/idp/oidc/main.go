package main

import (
	"bytes"
	"context"
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
	"syscall"
	"time"

	"github.com/zitadel/oidc/v3/example/server/exampleop"
	"github.com/zitadel/oidc/v3/example/server/storage"
)

func main() {
	apiURL := os.Getenv("API_URL")
	pat := readPAT(os.Getenv("PAT_FILE"))
	domain := os.Getenv("API_DOMAIN")
	schema := os.Getenv("SCHEMA")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)

	issuer := fmt.Sprintf("%s://%s:%s/", schema, host, port)
	redirectURI := fmt.Sprintf("%s/idps/callback", apiURL)

	clientID := "web"
	clientSecret := "secret"
	storage.RegisterClients(
		storage.WebClient(clientID, clientSecret, redirectURI),
	)

	storage := storage.NewStorage(storage.NewUserStore(issuer))
	router := exampleop.SetupServer(issuer, storage, logger, false)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	createZitadelResources(apiURL, pat, domain, issuer, clientID, clientSecret)

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

func createZitadelResources(apiURL, pat, domain, issuer, clientID, clientSecret string) error {
	idpID, err := CreateIDP(apiURL, pat, domain, issuer, clientID, clientSecret)
	if err != nil {
		return err
	}
	return ActivateIDP(apiURL, pat, domain, idpID)
}

type createIDP struct {
	Name             string          `json:"name"`
	Issuer           string          `json:"issuer"`
	ClientId         string          `json:"clientId"`
	ClientSecret     string          `json:"clientSecret"`
	Scopes           []string        `json:"scopes"`
	ProviderOptions  providerOptions `json:"providerOptions"`
	IsIdTokenMapping bool            `json:"isIdTokenMapping"`
	UsePkce          bool            `json:"usePkce"`
}

type providerOptions struct {
	IsLinkingAllowed  bool   `json:"isLinkingAllowed"`
	IsCreationAllowed bool   `json:"isCreationAllowed"`
	IsAutoCreation    bool   `json:"isAutoCreation"`
	IsAutoUpdate      bool   `json:"isAutoUpdate"`
	AutoLinking       string `json:"autoLinking"`
}

type idp struct {
	ID string `json:"id"`
}

func CreateIDP(apiURL, pat, domain string, issuer, clientID, clientSecret string) (string, error) {
	createIDP := &createIDP{
		Name:         "OIDC",
		Issuer:       issuer,
		ClientId:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "profile", "email"},
		ProviderOptions: providerOptions{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
			AutoLinking:       "AUTO_LINKING_OPTION_USERNAME",
		},
		IsIdTokenMapping: false,
		UsePkce:          false,
	}

	resp, err := doRequestWithHeaders(apiURL+"/admin/v1/idps/generic_oidc", pat, domain, createIDP)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	idp := new(idp)
	if err := json.Unmarshal(data, idp); err != nil {
		return "", err
	}
	return idp.ID, nil
}

type activateIDP struct {
	IdpId string `json:"idpId"`
}

func ActivateIDP(apiURL, pat, domain string, idpID string) error {
	activateIDP := &activateIDP{
		IdpId: idpID,
	}
	_, err := doRequestWithHeaders(apiURL+"/admin/v1/policies/login/idps", pat, domain, activateIDP)
	return err
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
