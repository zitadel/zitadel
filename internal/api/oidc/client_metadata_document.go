package oidc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/net/idna"
	"golang.org/x/sync/singleflight"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/connector"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

// Client ID Metadata Document (CIMD) support. A client_id that is an absolute HTTPS URL is
// treated as a pointer to a document of client metadata (the same shape as an RFC 7591
// registration request). The document is fetched, validated and the client is treated as an
// ephemeral public client without any database entry.
//
// The security anchor is the origin match: every redirect_uri must share the exact origin
// (scheme, host, port) of the client_id URL, so a document can only authorize redirects back
// to the host that serves it. CIMD clients are always public (token_endpoint_auth_method
// none) and never run in dev mode, so no client secret is trusted and no redirect wildcards
// are allowed. The document is fetched through the shared SSRF-safe HTTP client, whose
// operator denylist blocks loopback, private, link-local and cloud-metadata addresses at dial
// time.
//
// The fetch happens inline on the public, unauthenticated authorization endpoint, so it is
// guarded against amplification: concurrent resolutions of the same client_id are collapsed
// with a singleflight group, and failures are negatively cached for a short floor so a
// repeated unresolvable client_id does not trigger an outbound request on every call.

const (
	// clientIDMetadataMaxBodyBytes bounds the size of a fetched Client ID Metadata Document.
	clientIDMetadataMaxBodyBytes = 100 * 1024
	// clientIDMetadataFetchTimeout bounds a single document fetch. It is intentionally
	// tighter than the shared HTTP client timeout, as the fetch happens inline on the
	// authorization and token endpoints.
	clientIDMetadataFetchTimeout = 5 * time.Second
	// clientIDMetadataMaxCacheTTL caps how long a fetched document is cached, regardless of
	// the Cache-Control or Expires headers the document is served with. It is the authoritative
	// upper bound; the Caches.ClientIDMetadataDocuments.MaxAge config is a second, independent
	// ceiling and should be kept in sync.
	clientIDMetadataMaxCacheTTL = 15 * time.Minute
	// clientIDMetadataNegativeTTL is how long a failed resolution is cached, so a repeated
	// unresolvable client_id on the public authorize endpoint collapses to one fetch per window.
	clientIDMetadataNegativeTTL = time.Minute
)

// clientIDMetadataDocumentEnabled reports whether the client_id should be resolved as a
// Client ID Metadata Document: the feature must be enabled for the instance and the client_id
// must look like a CIMD URL.
func clientIDMetadataDocumentEnabled(ctx context.Context, clientID string) bool {
	return authz.GetFeatures(ctx).OIDCClientIDMetadataDocument && looksLikeClientIDMetadataURL(clientID)
}

// looksLikeClientIDMetadataURL reports whether the client_id should be resolved as a Client
// ID Metadata Document. A regular ZITADEL client_id is a numeric snowflake (optionally
// suffixed) and is never an absolute HTTPS URL, so this branch is unambiguous.
func looksLikeClientIDMetadataURL(clientID string) bool {
	u, err := url.Parse(clientID)
	return err == nil && u.IsAbs() && strings.EqualFold(u.Scheme, "https") && u.Host != ""
}

type clientIDMetadataCacheIndex int

const (
	clientIDMetadataCacheIndexUnspecified clientIDMetadataCacheIndex = iota
	clientIDMetadataCacheIndexURL
)

// clientIDMetadataCacheEntry caches the outcome of a resolution. A nil Client marks a
// negatively cached (failed) resolution. The per-entry Expiry is honored by the resolver on
// read, on top of the cache's own max age, so the document's own Cache-Control or Expires
// lifetime is respected and capped.
type clientIDMetadataCacheEntry struct {
	Key    string            `json:"key"`
	Client *query.OIDCClient `json:"client,omitempty"`
	Expiry time.Time         `json:"expiry"`
}

// Keys implements cache.Entry.
func (e *clientIDMetadataCacheEntry) Keys(index clientIDMetadataCacheIndex) []string {
	if index == clientIDMetadataCacheIndexURL {
		return []string{e.Key}
	}
	return nil
}

func clientIDMetadataCacheKey(instanceID, clientID string) string {
	return instanceID + "|" + clientID
}

// StartClientIDMetadataDocumentCache starts the cache that holds resolved Client ID Metadata
// Documents. The cached entry type wraps a *query.OIDCClient, so it stays in this package
// rather than a domain sub-package (unlike the federated logout cache) to avoid a domain to
// query dependency; the cache is therefore constructed here instead of exposing its element
// types from the package API.
func StartClientIDMetadataDocumentCache(background context.Context, conf *cache.Config, connectors connector.Connectors) (cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry], error) {
	return connector.StartCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry](
		background,
		[]clientIDMetadataCacheIndex{clientIDMetadataCacheIndexURL},
		cache.PurposeClientIDMetadataDocument,
		conf,
		connectors,
	)
}

// clientIDMetadataResolver fetches and validates Client ID Metadata Documents and turns them
// into ephemeral, in-memory OIDC clients. It reuses the shared SSRF-safe HTTP client and the
// generic cache; it never touches the database.
type clientIDMetadataResolver struct {
	httpClient          *http.Client
	cache               cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]
	group               singleflight.Group
	accessTokenLifetime time.Duration
	idTokenLifetime     time.Duration
}

func newClientIDMetadataResolver(
	httpClient *http.Client,
	documentCache cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry],
	accessTokenLifetime, idTokenLifetime time.Duration,
) *clientIDMetadataResolver {
	return &clientIDMetadataResolver{
		httpClient:          httpClient,
		cache:               documentCache,
		accessTokenLifetime: accessTokenLifetime,
		idTokenLifetime:     idTokenLifetime,
	}
}

// ResolveClient returns the synthetic public OIDC client described by the Client ID Metadata
// Document located at clientID, which must be an absolute HTTPS URL. A valid cache entry
// (positive or negative) is served without a fetch; otherwise the document is fetched,
// validated and cached. Concurrent resolutions of the same client_id share a single fetch.
// Any fetch or validation failure is reported as an invalid_client error.
func (r *clientIDMetadataResolver) ResolveClient(ctx context.Context, instanceID, clientID string) (*query.OIDCClient, error) {
	key := clientIDMetadataCacheKey(instanceID, clientID)
	if client, negative, ok := r.cached(ctx, key); ok {
		if negative {
			return nil, r.invalidClient(ctx, nil, "client id metadata document could not be resolved")
		}
		return client, nil
	}
	resolved, err, _ := r.group.Do(key, func() (any, error) {
		return r.resolveAndCache(ctx, clientID, key)
	})
	if err != nil {
		return nil, err
	}
	return resolved.(*query.OIDCClient), nil
}

func (r *clientIDMetadataResolver) resolveAndCache(ctx context.Context, clientID, key string) (*query.OIDCClient, error) {
	client, ttl, err := r.fetchAndValidate(ctx, clientID)
	if err != nil {
		// Negatively cache the failure for a short floor so a repeated unresolvable client_id
		// on the public authorize endpoint does not trigger an outbound request every time.
		r.cacheSet(ctx, key, nil, clientIDMetadataNegativeTTL)
		return nil, err
	}
	if ttl > 0 {
		r.cacheSet(ctx, key, client, ttl)
	}
	return client, nil
}

// cached returns the cached resolution for key, if any is still valid. negative is true when
// the entry records a failed resolution. The returned client must be treated as immutable: it
// is shared across requests (the in-memory connector hands back the stored pointer).
func (r *clientIDMetadataResolver) cached(ctx context.Context, key string) (client *query.OIDCClient, negative, ok bool) {
	if r.cache == nil {
		return nil, false, false
	}
	entry, found := r.cache.Get(ctx, clientIDMetadataCacheIndexURL, key)
	if !found || entry == nil {
		return nil, false, false
	}
	if time.Now().After(entry.Expiry) {
		_ = r.cache.Invalidate(ctx, clientIDMetadataCacheIndexURL, key)
		return nil, false, false
	}
	return entry.Client, entry.Client == nil, true
}

func (r *clientIDMetadataResolver) cacheSet(ctx context.Context, key string, client *query.OIDCClient, ttl time.Duration) {
	if r.cache == nil {
		return
	}
	r.cache.Set(ctx, &clientIDMetadataCacheEntry{
		Key:    key,
		Client: client,
		Expiry: time.Now().Add(ttl),
	})
}

func (r *clientIDMetadataResolver) fetchAndValidate(ctx context.Context, clientID string) (*query.OIDCClient, time.Duration, error) {
	clientURL, err := url.Parse(clientID)
	if err != nil || !clientURL.IsAbs() || !strings.EqualFold(clientURL.Scheme, "https") || clientURL.Host == "" {
		return nil, 0, r.invalidClient(ctx, err, "client_id is not a valid https url")
	}

	fetchCtx, cancel := context.WithTimeout(ctx, clientIDMetadataFetchTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(fetchCtx, http.MethodGet, clientID, nil)
	if err != nil {
		return nil, 0, r.invalidClient(ctx, err, "client id metadata document request could not be built")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		// The SSRF-safe client blocks loopback, private, link-local and cloud-metadata
		// addresses at dial time, so a blocked target surfaces here as a transport error.
		return nil, 0, r.invalidClient(ctx, err, "client id metadata document could not be fetched")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, r.invalidClient(ctx, nil, "client id metadata document could not be fetched")
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, clientIDMetadataMaxBodyBytes+1))
	if err != nil {
		return nil, 0, r.invalidClient(ctx, err, "client id metadata document could not be read")
	}
	if int64(len(body)) > clientIDMetadataMaxBodyBytes {
		return nil, 0, r.invalidClient(ctx, nil, "client id metadata document is too large")
	}

	var doc clientRegistrationRequest
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, 0, r.invalidClient(ctx, err, "client id metadata document could not be parsed")
	}

	client, err := r.documentToClient(ctx, clientID, clientURL, &doc)
	if err != nil {
		return nil, 0, err
	}
	return client, cacheTTLFromResponse(resp.Header, clientIDMetadataMaxCacheTTL), nil
}

// documentToClient maps a validated metadata document to a synthetic public OIDC client.
// CIMD clients are always public and identified by an HTTPS URL, so they are treated as web
// clients without a project, secret or dev mode. The reused RFC 7591 mappers reject
// unsupported grant, response and application types; the origin match restricts the redirect
// URIs to the client_id origin.
func (r *clientIDMetadataResolver) documentToClient(ctx context.Context, clientID string, clientURL *url.URL, doc *clientRegistrationRequest) (*query.OIDCClient, error) {
	if doc.JWKsURI != "" || len(doc.JWKs) > 0 {
		return nil, r.invalidClient(ctx, nil, "jwks and jwks_uri are not supported for client id metadata documents")
	}
	// A document that declares a confidential auth method is rejected: a client_id that is a
	// public URL cannot safely be associated with a client secret. An empty value defaults to
	// the only supported method, none.
	if doc.TokenEndpointAuthMethod != "" && doc.TokenEndpointAuthMethod != "none" {
		return nil, r.invalidClient(ctx, nil, "only token_endpoint_auth_method none is supported for client id metadata documents")
	}
	// application_type is validated but not used: a CIMD client is identified by an HTTPS URL
	// and its redirect URIs are origin matched against it, so it is always treated as a web
	// client (a native document would in any case fail the origin match for its loopback or
	// custom-scheme redirect URIs).
	if _, regErr := registrationApplicationTypeToDomain(doc.ApplicationType); regErr != nil {
		return nil, r.invalidClient(ctx, regErr, regErr.ErrorDescription)
	}
	grantTypes, regErr := registrationGrantTypesToDomain(doc.GrantTypes)
	if regErr != nil {
		return nil, r.invalidClient(ctx, regErr, regErr.ErrorDescription)
	}
	responseTypes, regErr := registrationResponseTypesToDomain(doc.ResponseTypes)
	if regErr != nil {
		return nil, r.invalidClient(ctx, regErr, regErr.ErrorDescription)
	}

	redirectURIs := originMatchedURIs(clientURL, doc.RedirectURIs)
	if len(redirectURIs) == 0 {
		return nil, r.invalidClient(ctx, nil, "no redirect_uri matches the client_id origin")
	}
	postLogoutRedirectURIs := originMatchedURIs(clientURL, doc.PostLogoutRedirectURIs)

	applicationType := domain.OIDCApplicationTypeWeb
	authMethod := domain.OIDCAuthMethodTypeNone
	if compliance := domain.GetOIDCV1Compliance(&applicationType, grantTypes, &authMethod, redirectURIs); compliance.NoneCompliant {
		return nil, r.invalidClient(ctx, nil, "the client id metadata document is not compliant")
	}

	return &query.OIDCClient{
		InstanceID:      authz.GetInstance(ctx).InstanceID(),
		ClientID:        clientID,
		State:           domain.AppStateActive,
		RedirectURIs:    redirectURIs,
		ResponseTypes:   responseTypes,
		GrantTypes:      grantTypes,
		ApplicationType: applicationType,
		AuthMethodType:  authMethod,
		// PostLogoutRedirectURIs are origin matched as well so they cannot point off-origin.
		PostLogoutRedirectURIs: postLogoutRedirectURIs,
		IsDevMode:              false,
		AccessTokenType:        domain.OIDCTokenTypeBearer,
		// No project: CIMD clients are ephemeral and carry no project roles. The token flow
		// already tolerates an empty project id.
		ProjectID: "",
		Settings: &query.OIDCSettings{
			AccessTokenLifetime: r.accessTokenLifetime,
			IdTokenLifetime:     r.idTokenLifetime,
		},
	}, nil
}

func (r *clientIDMetadataResolver) invalidClient(ctx context.Context, parent error, description string) error {
	// description is passed as an argument rather than as the format string: some
	// descriptions are derived from attacker-controlled document fields and could contain
	// formatting verbs.
	return oidc.ErrInvalidClient().
		WithParent(parent).
		WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError).
		WithDescription("%s", description)
}

// originMatchedURIs returns the subset of uris whose origin (scheme, host, port) exactly
// matches the client_id URL. This is the CIMD trust anchor: a document can only authorize
// redirects back to the origin that serves it.
func originMatchedURIs(clientURL *url.URL, uris []string) []string {
	matched := make([]string, 0, len(uris))
	for _, raw := range trimSpaceSlice(uris) {
		if sameOrigin(clientURL, raw) {
			matched = append(matched, raw)
		}
	}
	if len(matched) == 0 {
		return nil
	}
	return matched
}

// sameOrigin reports whether rawURI is absolute and shares the exact origin of clientURL.
// The scheme comparison is case-insensitive, the host is compared in lower-case ASCII
// (punycode) form, and the port is normalized to the scheme default. Only the origin is
// compared; any userinfo, path or query on the redirect URI is ignored, which is safe because
// it does not change the destination origin.
func sameOrigin(clientURL *url.URL, rawURI string) bool {
	u, err := url.Parse(rawURI)
	if err != nil || !u.IsAbs() {
		return false
	}
	// A redirect URI with embedded credentials is rejected: it is non-standard and only
	// invites confusion, while the origin it points to is unchanged.
	if u.User != nil {
		return false
	}
	if !strings.EqualFold(u.Scheme, clientURL.Scheme) {
		return false
	}
	return normalizedHostPort(u) == normalizedHostPort(clientURL)
}

func normalizedHostPort(u *url.URL) string {
	host := strings.ToLower(strings.TrimSuffix(u.Hostname(), "."))
	if ascii, err := idna.Lookup.ToASCII(host); err == nil {
		host = ascii
	}
	port := u.Port()
	if port == "" {
		port = defaultPortForScheme(u.Scheme)
	}
	return host + ":" + port
}

func defaultPortForScheme(scheme string) string {
	switch strings.ToLower(scheme) {
	case "https":
		return "443"
	case "http":
		return "80"
	default:
		return ""
	}
}

// cacheTTLFromResponse derives a cache lifetime from the response Cache-Control and Expires
// headers, capped at max. A no-store or no-cache directive disables caching (returns 0). When
// no caching information is present the lifetime defaults to max.
func cacheTTLFromResponse(header http.Header, max time.Duration) time.Duration {
	cacheControl := strings.ToLower(header.Get("Cache-Control"))
	if strings.Contains(cacheControl, "no-store") || strings.Contains(cacheControl, "no-cache") {
		return 0
	}
	if maxAge, ok := maxAgeFromCacheControl(cacheControl); ok {
		if maxAge <= 0 {
			return 0
		}
		return capDuration(maxAge, max)
	}
	if expires := header.Get("Expires"); expires != "" {
		if t, err := http.ParseTime(expires); err == nil {
			ttl := time.Until(t)
			if ttl <= 0 {
				return 0
			}
			return capDuration(ttl, max)
		}
	}
	return max
}

func maxAgeFromCacheControl(cacheControl string) (time.Duration, bool) {
	for _, directive := range strings.Split(cacheControl, ",") {
		directive = strings.TrimSpace(directive)
		value, ok := strings.CutPrefix(directive, "max-age=")
		if !ok {
			continue
		}
		seconds, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			return 0, false
		}
		return time.Duration(seconds) * time.Second, true
	}
	return 0, false
}

func capDuration(d, max time.Duration) time.Duration {
	if d > max {
		return max
	}
	return d
}
