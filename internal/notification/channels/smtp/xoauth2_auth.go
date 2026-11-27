package smtp

import (
	"context"
	"hash/fnv"
	"net/smtp"
	"net/url"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/zitadel/zitadel/internal/denylist"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var tokenSourceCache *lru.TwoQueueCache[string, oauth2.TokenSource]

func init() {
	var err error
	tokenSourceCache, err = lru.New2Q[string, oauth2.TokenSource](100)
	if err != nil {
		panic(err)
	}
}

type XOAuth2AuthConfig struct {
	User                  string
	TokenEndpoint         string
	Scopes                []string
	ClientCredentialsAuth *OAuth2ClientCredentials
}

type OAuth2ClientCredentials struct {
	ClientId     string
	ClientSecret string
}

func (cfg XOAuth2AuthConfig) Hash() string {
	hash := fnv.New128()
	hash.Write([]byte(cfg.User))
	hash.Write([]byte(cfg.TokenEndpoint))
	for _, s := range cfg.Scopes {
		hash.Write([]byte(s))
	}
	if cfg.ClientCredentialsAuth != nil {
		hash.Write([]byte(cfg.ClientCredentialsAuth.ClientId))
		hash.Write([]byte(cfg.ClientCredentialsAuth.ClientSecret))
	}
	return string(hash.Sum(nil))
}

type xoauth2Auth struct {
	host         string
	config       XOAuth2AuthConfig
	restrictions *Restrictions
}

// XOAuth2Auth returns the XOAUTH2 [smtp.Auth]. It fetches the token from the token endpoint of the config.
// If restrictions are set, the token endpoint must use https, and the request is sent through restrictions.HTTPClient.
func XOAuth2Auth(config XOAuth2AuthConfig, host string, restrictions *Restrictions) smtp.Auth {
	return &xoauth2Auth{
		host:         host,
		config:       config,
		restrictions: restrictions,
	}
}

// checkTokenEndpoint refuses a token endpoint that is not allowed for restricted configs.
func (r *Restrictions) checkTokenEndpoint(endpoint string) error {
	uri, err := url.Parse(endpoint)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "SMTP-Wd3Lq8", "Errors.Invalid.Argument")
	}
	if !strings.EqualFold(uri.Scheme, "https") {
		return zerrors.ThrowInvalidArgument(nil, "SMTP-Nx6Tf1", "Errors.SMTPConfig.OrgHTTPSRequired")
	}
	if err = denylist.IsHostBlocked(r.DenyList, uri.Hostname()); err != nil {
		return zerrors.ThrowInvalidArgument(err, "SMTP-Gy2Ke5", "Errors.NotificationProvider.Blocked.Endpoint")
	}
	if r.HTTPClient == nil {
		return zerrors.ThrowInternal(nil, "SMTP-Qb8Zh4", "Errors.Internal")
	}
	return nil
}
func (a *xoauth2Auth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.host {
		return "", nil, zerrors.ThrowInternal(nil, "SMTP-eRJLyi", "wrong host name")
	}

	hash := a.config.Hash()
	tokenCtx := context.Background()
	if a.restrictions != nil {
		if err := a.restrictions.checkTokenEndpoint(a.config.TokenEndpoint); err != nil {
			return "", nil, err
		}
		// restricted token sources must never be shared with unrestricted ones
		hash = "restricted:" + hash
		tokenCtx = context.WithValue(tokenCtx, oauth2.HTTPClient, a.restrictions.HTTPClient)
	}
	tokenSource, ok := tokenSourceCache.Get(hash)
	if !ok {
		if a.config.ClientCredentialsAuth != nil {
			config := &clientcredentials.Config{
				ClientID:     a.config.ClientCredentialsAuth.ClientId,
				ClientSecret: a.config.ClientCredentialsAuth.ClientSecret,
				Scopes:       a.config.Scopes,
				TokenURL:     a.config.TokenEndpoint,
			}
			tokenSource = config.TokenSource(tokenCtx)
			tokenSourceCache.Add(hash, tokenSource)
		}
	}

	token, err := tokenSource.Token()
	if err != nil {
		return "", nil, zerrors.ThrowInternal(err, "SMTP-fjHcJu", "Failed to get token to connect with smtp server")
	}

	resp := []byte("user=" + a.config.User + "\x01" + "auth=Bearer " + token.AccessToken + "\x01\x01")
	return "XOAUTH2", resp, nil
}

func (a *xoauth2Auth) Next(_ []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}
	return nil, zerrors.ThrowInternal(nil, "SMTP-Pqsrj9", "unexpected server challenge for XOAUTH2 auth method")
}
