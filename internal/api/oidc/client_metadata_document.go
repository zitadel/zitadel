package oidc

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/net/idna"
	"golang.org/x/sync/singleflight"
	"golang.org/x/time/rate"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/connector"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
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
//
// A client_id URL is only resolved when both the system and the instance allow it. The system
// allows URLs in the runtime configuration (OIDC.ClientIDMetadataDocument), which no instance can
// change, so on a multi-instance system no instance can make the server fetch a URL its operator
// did not allow. Within that, each instance allows URLs in its security settings, which can be
// changed at runtime through the API. A client_id that is not allowed is never fetched or cached;
// it falls through to the regular client lookup and is answered with invalid_client, exactly as
// before CIMD existed. Both layers allow nothing by default.

const (
	// clientIDMetadataMaxBodyBytes bounds the size of a fetched Client ID Metadata Document.
	clientIDMetadataMaxBodyBytes = 100 * 1024
	// clientIDMetadataFetchTimeout bounds a single document fetch. It is intentionally
	// tighter than the shared HTTP client timeout, as the fetch happens inline on the
	// authorization and token endpoints.
	clientIDMetadataFetchTimeout = 5 * time.Second
	// clientIDMetadataResolveTimeout bounds a whole shared resolution: the fetch plus the cache
	// writes that follow it. The resolution is detached from the caller's cancellation, so this
	// is what keeps a blocked cache connector from holding the shared call open.
	clientIDMetadataResolveTimeout = 2 * clientIDMetadataFetchTimeout
	// clientIDMetadataFetchesPerSecond and clientIDMetadataFetchBurst are the per-instance fetch
	// bound when the configuration sets none.
	clientIDMetadataFetchesPerSecond = 5
	clientIDMetadataFetchBurst       = 20
	// clientIDMetadataMaxCacheTTL caps how long a fetched document is cached, regardless of
	// the Cache-Control or Expires headers the document is served with. It is the authoritative
	// upper bound; the Caches.ClientIDMetadataDocuments.MaxAge config is a second, independent
	// ceiling and should be kept in sync.
	clientIDMetadataMaxCacheTTL = 15 * time.Minute
	// clientIDMetadataNegativeTTL is how long a failed resolution is cached, so a repeated
	// unresolvable client_id on the public authorize endpoint collapses to one fetch per window.
	clientIDMetadataNegativeTTL = time.Minute
)

// ClientIDMetadataDocumentConfig is the system configuration of Client ID Metadata Documents.
type ClientIDMetadataDocumentConfig struct {
	// AllowedURLs lists the client_id URLs that may be resolved. An entry ending with a slash
	// allows every client_id it is a prefix of, any other entry allows exactly that client_id.
	// Entries are compared as plain strings, as the specification requires for client_id URLs.
	AllowedURLs []string
	// AllowAnyURL allows every valid client_id URL, regardless of AllowedURLs.
	AllowAnyURL bool
	// FetchesPerSecond and FetchBurst bound the outbound document fetches per instance. The
	// authorize and token endpoints are public, so without a bound a caller could have ZITADEL
	// fetch, and cache, one document per made-up client_id under an allowed prefix.
	FetchesPerSecond float64
	FetchBurst       int
}

// clientIDMetadataAllowlist is the validated system configuration: the client_id URLs the
// system allows, within which every instance chooses its own.
type clientIDMetadataAllowlist struct {
	urls             []string
	allowAny         bool
	fetchesPerSecond float64
	fetchBurst       int
}

// newClientIDMetadataAllowlist validates the configured entries. Every entry must itself be a
// valid client_id URL, so that a prefix entry cannot be widened by a client_id that is not one.
func newClientIDMetadataAllowlist(config ClientIDMetadataDocumentConfig) (clientIDMetadataAllowlist, error) {
	allowlist := clientIDMetadataAllowlist{
		urls:             make([]string, 0, len(config.AllowedURLs)),
		allowAny:         config.AllowAnyURL,
		fetchesPerSecond: config.FetchesPerSecond,
		fetchBurst:       config.FetchBurst,
	}
	for _, entry := range config.AllowedURLs {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if !domain.IsClientIDMetadataDocumentURL(entry) {
			return clientIDMetadataAllowlist{}, zerrors.ThrowInvalidArgumentf(nil, "OIDC-qX3vL", "invalid client id metadata document allowed url %q", entry)
		}
		allowlist.urls = append(allowlist.urls, entry)
	}
	return allowlist, nil
}

// overlaps reports whether at least one client_id URL is allowed by both the system and an
// instance that allows instanceURLs, or every URL the system allows if instanceAllowAny.
func (a clientIDMetadataAllowlist) overlaps(instanceURLs []string, instanceAllowAny bool) bool {
	if a.allowAny {
		return instanceAllowAny || len(instanceURLs) > 0
	}
	if instanceAllowAny {
		return len(a.urls) > 0
	}
	for _, systemURL := range a.urls {
		for _, instanceURL := range instanceURLs {
			if domain.ClientIDMetadataDocumentURLAllowed([]string{systemURL}, instanceURL) ||
				domain.ClientIDMetadataDocumentURLAllowed([]string{instanceURL}, systemURL) {
				return true
			}
		}
	}
	return false
}

// allows reports whether the system allows clientID.
func (a clientIDMetadataAllowlist) allows(clientID string) bool {
	return a.allowAny || domain.ClientIDMetadataDocumentURLAllowed(a.urls, clientID)
}

// instanceAllowsClientIDMetadataDocumentURL reports whether the instance in ctx allows clientID.
func instanceAllowsClientIDMetadataDocumentURL(ctx context.Context, clientID string) bool {
	instance := authz.GetInstance(ctx)
	return instance.ClientIDMetadataDocumentAllowAnyURL() || domain.ClientIDMetadataDocumentURLAllowed(instance.ClientIDMetadataDocumentAllowedURLs(), clientID)
}

// clientIDMetadataDocument is a fetched Client ID Metadata Document: the client metadata of an
// RFC 7591 registration request plus the client_id it describes.
type clientIDMetadataDocument struct {
	ClientID string `json:"client_id"`
	clientRegistrationRequest
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
	allowlist           clientIDMetadataAllowlist
	cache               cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]
	group               singleflight.Group
	accessTokenLifetime time.Duration
	idTokenLifetime     time.Duration
	resolveTimeout      time.Duration
	fetchesPerSecond    rate.Limit
	fetchBurst          int
	fetchLimiters       sync.Map // instance ID -> *rate.Limiter
	logger              *slog.Logger
}

// newClientIDMetadataResolver returns a resolver for the client_id URLs the allowlist allows.
// The document is fetched with a copy of httpClient that does not follow redirects: the
// document must be served from the client_id URL itself, and a redirect could otherwise lead
// the fetch to a URL the allowlist does not allow.
func newClientIDMetadataResolver(
	httpClient *http.Client,
	allowlist clientIDMetadataAllowlist,
	documentCache cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry],
	accessTokenLifetime, idTokenLifetime time.Duration,
	logger *slog.Logger,
) *clientIDMetadataResolver {
	if logger == nil {
		logger = slog.Default()
	}
	documentClient := *httpClient
	documentClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resolver := &clientIDMetadataResolver{
		httpClient:          &documentClient,
		allowlist:           allowlist,
		cache:               documentCache,
		accessTokenLifetime: accessTokenLifetime,
		idTokenLifetime:     idTokenLifetime,
		resolveTimeout:      clientIDMetadataResolveTimeout,
		fetchesPerSecond:    rate.Limit(clientIDMetadataFetchesPerSecond),
		fetchBurst:          clientIDMetadataFetchBurst,
		logger:              logger,
	}
	if allowlist.fetchesPerSecond > 0 {
		resolver.fetchesPerSecond = rate.Limit(allowlist.fetchesPerSecond)
	}
	if allowlist.fetchBurst > 0 {
		resolver.fetchBurst = allowlist.fetchBurst
	}
	return resolver
}

// fetchAllowed reports whether the instance may fetch another document now. Each instance has
// its own bucket, so one instance cannot use up another's fetches.
func (r *clientIDMetadataResolver) fetchAllowed(instanceID string) bool {
	limiter, _ := r.fetchLimiters.LoadOrStore(instanceID, rate.NewLimiter(r.fetchesPerSecond, r.fetchBurst))
	return limiter.(*rate.Limiter).Allow()
}

// forRequest returns a copy of the cached client with the instance's current login settings
// applied. The cache holds only what the document says, so a change to the required Login V2
// base URI applies to the next request rather than when the cached document expires.
func forRequest(ctx context.Context, cached *query.OIDCClient) *query.OIDCClient {
	client := *cached
	client.LoginVersion = domain.LoginVersion2
	client.LoginBaseURI = nil
	if loginV2 := authz.GetFeatures(ctx).LoginV2; loginV2.Required {
		client.LoginBaseURI = (*query.URL)(loginV2.BaseURI)
	}
	return &client
}

// ResolveClient returns the synthetic public OIDC client described by the Client ID Metadata
// Document located at clientID, which must be an absolute HTTPS URL. A valid cache entry
// (positive or negative) is served without a fetch; otherwise the document is fetched,
// validated and cached. Concurrent resolutions of the same client_id share a single fetch.
// Any fetch or validation failure is reported as an invalid_client error.
//
// The shared resolution deliberately does not inherit the caller's cancellation. It is started
// by whichever request happens to miss the cache first, but its result is shared, so binding it
// to that one request would let a single client disconnect abort the fetch for every waiter and,
// worse, store the resulting failure as a negative cache entry that blocks the client_id for
// everyone until it expires. context.WithoutCancel keeps the values (instance, features,
// tracing) the resolution needs; the whole resolution, fetch and cache writes, stays bounded by
// the resolve timeout, so detaching cannot leak a call that runs forever. Each caller still waits on its own
// context, so a disconnected caller returns immediately instead of blocking on the shared work.
// Supported reports whether Client ID Metadata Documents can be resolved for the instance in
// ctx: the instance security settings must enable them, and at least one client_id URL must be
// allowed by both the system and the instance. Discovery advertises support only then, so a
// client is not steered towards a mechanism that cannot resolve it.
func (r *clientIDMetadataResolver) Supported(ctx context.Context) bool {
	if r == nil {
		return false
	}
	instance := authz.GetInstance(ctx)
	return instance.EnableClientIDMetadataDocument() &&
		r.allowlist.overlaps(instance.ClientIDMetadataDocumentAllowedURLs(), instance.ClientIDMetadataDocumentAllowAnyURL())
}

// Handles reports whether clientID is resolved as a Client ID Metadata Document: support must be
// enabled for the instance in ctx and clientID must be a valid client_id URL that both the system
// and the instance allow. Any other client_id goes through the regular client lookup.
func (r *clientIDMetadataResolver) Handles(ctx context.Context, clientID string) bool {
	if r == nil {
		return false
	}
	return authz.GetInstance(ctx).EnableClientIDMetadataDocument() && r.allowed(ctx, clientID)
}

// allowed reports whether clientID is a valid client_id URL that both the system and the
// instance in ctx allow.
func (r *clientIDMetadataResolver) allowed(ctx context.Context, clientID string) bool {
	return domain.IsClientIDMetadataDocumentURL(clientID) && r.allowlist.allows(clientID) && instanceAllowsClientIDMetadataDocumentURL(ctx, clientID)
}

func (r *clientIDMetadataResolver) ResolveClient(ctx context.Context, instanceID, clientID string) (*query.OIDCClient, error) {
	if !r.allowed(ctx, clientID) {
		return nil, r.invalidClient(ctx, nil, "client_id is not an allowed client id metadata document url")
	}
	key := clientIDMetadataCacheKey(instanceID, clientID)
	if client, negative, ok := r.cached(ctx, key); ok {
		if negative {
			return nil, r.invalidClient(ctx, nil, "client id metadata document could not be resolved")
		}
		return forRequest(ctx, client), nil
	}
	results := r.group.DoChan(key, func() (any, error) {
		resolveCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), r.resolveTimeout)
		defer cancel()
		return r.resolveAndCache(resolveCtx, instanceID, clientID, key)
	})
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-results:
		if result.Err != nil {
			return nil, result.Err
		}
		return forRequest(ctx, result.Val.(*query.OIDCClient)), nil
	}
}

func (r *clientIDMetadataResolver) resolveAndCache(ctx context.Context, instanceID, clientID, key string) (*query.OIDCClient, error) {
	// A fetch that is not allowed now is not a property of the document, so it is not cached:
	// a later request resolves normally once the bucket refills.
	if !r.fetchAllowed(instanceID) {
		return nil, r.invalidClient(ctx, nil, "too many client id metadata document fetches, retry later")
	}
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
	if err != nil {
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

	var doc clientIDMetadataDocument
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, 0, r.invalidClient(ctx, err, "client id metadata document could not be parsed")
	}
	// The document must name the URL it was fetched from as its client_id, compared as plain
	// strings, so a document cannot be served for a client_id other than its own.
	if doc.ClientID != clientID {
		return nil, 0, r.invalidClient(ctx, nil, "client id metadata document client_id does not match the client_id url")
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
func (r *clientIDMetadataResolver) documentToClient(ctx context.Context, clientID string, clientURL *url.URL, doc *clientIDMetadataDocument) (*query.OIDCClient, error) {
	if doc.JWKsURI != "" || len(doc.JWKs) > 0 {
		return nil, r.invalidClient(ctx, nil, "jwks and jwks_uri are not supported for client id metadata documents")
	}
	// A document that declares a confidential auth method is rejected: a client_id that is a
	// public URL cannot safely be associated with a client secret. An empty value defaults to
	// the only supported method, none.
	if doc.TokenEndpointAuthMethod != "" && doc.TokenEndpointAuthMethod != "none" {
		return nil, r.invalidClient(ctx, nil, "only token_endpoint_auth_method none is supported for client id metadata documents")
	}
	// application_type is validated but its result is not used: a CIMD client is identified by
	// an HTTPS URL and its redirect URIs are origin matched against it, so it is always treated
	// as a web client (a native document would in any case fail the origin match for its
	// loopback or custom-scheme redirect URIs).
	if _, regErr := registrationApplicationTypeToDomain(doc.ApplicationType, doc.RedirectURIs); regErr != nil {
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

	client := &query.OIDCClient{
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
		// CIMD clients always use the v2 login; the base URI is applied per request.
		LoginVersion: domain.LoginVersion2,
	}
	return client, nil
}

// invalidClient is the single funnel for every resolution failure. They are all reported to the
// caller as invalid_client, because at this layer an unreachable or malformed document is not
// distinguishable from an unknown client and both mean the same thing to the client: this
// client_id does not resolve.
//
// The underlying cause is logged so an operator can tell a client-side problem (bad document,
// origin mismatch) from an instance-side one (egress broken, denylist too wide), which the
// invalid_client answer alone hides. It is logged at debug level: the resolution happens on the
// public, unauthenticated authorize endpoint, so the client_id is attacker controlled and a
// higher level would let anyone fill the logs.
func (r *clientIDMetadataResolver) invalidClient(ctx context.Context, parent error, description string) error {
	r.logger.DebugContext(ctx, "client id metadata document not resolved",
		"instanceID", authz.GetInstance(ctx).InstanceID(),
		"reason", description,
		"err", parent,
	)
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
