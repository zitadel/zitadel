package well_known

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	AppleAppSiteAssociationPath = "/.well-known/apple-app-site-association"
	AssetLinksPath              = "/.well-known/assetlinks.json"
)

// HandlerPrefixes are registered with [api.API.RegisterHandlerPrefixes].
var HandlerPrefixes = []string{
	AppleAppSiteAssociationPath,
	AssetLinksPath,
}

// Config holds runtime options for the well-known handlers.
type Config struct {
	// AppLinksCacheControlMaxAge sets the Cache-Control max-age for
	// apple-app-site-association and assetlinks.json responses.
	// 0 sets Cache-Control: no-store.
	AppLinksCacheControlMaxAge time.Duration
}

type appLinkQuerier interface {
	SearchOIDCAppLinkConfigs(ctx context.Context) ([]*query.OIDCAppLinkConfig, error)
}

type Handler struct {
	queries            appLinkQuerier
	cacheControlMaxAge time.Duration
}

func NewHandler(queries appLinkQuerier, config Config) *Handler {
	return &Handler{
		queries:            queries,
		cacheControlMaxAge: config.AppLinksCacheControlMaxAge,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case AppleAppSiteAssociationPath:
		h.serveAppleAppSiteAssociation(w, r)
	case AssetLinksPath:
		h.serveAssetLinks(w, r)
	default:
		http.NotFound(w, r)
	}
}

type appleAppSiteAssociation struct {
	WebCredentials appleWebCredentials `json:"webcredentials"`
}

type appleWebCredentials struct {
	Apps []string `json:"apps"`
}

type assetLink struct {
	Relation []string        `json:"relation"`
	Target   assetLinkTarget `json:"target"`
}

type assetLinkTarget struct {
	Namespace              string   `json:"namespace"`
	PackageName            string   `json:"package_name"`
	SHA256CertFingerprints []string `json:"sha256_cert_fingerprints"`
}

func (h *Handler) serveAppleAppSiteAssociation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	configs, err := h.queries.SearchOIDCAppLinkConfigs(r.Context())
	if err != nil {
		logging.Error(r.Context(), "unable to query OIDC app link configs for apple-app-site-association", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, r, buildAppleAppSiteAssociation(configs))
}

func (h *Handler) serveAssetLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	configs, err := h.queries.SearchOIDCAppLinkConfigs(r.Context())
	if err != nil {
		logging.Error(r.Context(), "unable to query OIDC app link configs for assetlinks.json", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, r, buildAssetLinks(r.Context(), configs))
}

func buildAppleAppSiteAssociation(configs []*query.OIDCAppLinkConfig) appleAppSiteAssociation {
	apps := make([]string, 0)
	seen := make(map[string]struct{})
	for _, cfg := range configs {
		if cfg.IOSTeamID == "" || cfg.IOSBundleID == "" {
			continue
		}
		appID := cfg.IOSTeamID + "." + cfg.IOSBundleID
		if _, ok := seen[appID]; ok {
			continue
		}
		seen[appID] = struct{}{}
		apps = append(apps, appID)
	}
	return appleAppSiteAssociation{
		WebCredentials: appleWebCredentials{Apps: apps},
	}
}

func buildAssetLinks(ctx context.Context, configs []*query.OIDCAppLinkConfig) []assetLink {
	links := make([]assetLink, 0)
	for _, cfg := range configs {
		if cfg.AndroidPackageName == "" || len(cfg.AndroidSHA256CertFingerprints) == 0 {
			continue
		}
		links = append(links, assetLink{
			// Google's Credential Manager prerequisites require both relations for
			// passkeys; some devices (e.g. Pixel) reject passkey creation when
			// handle_all_urls is missing.
			// https://developer.android.com/identity/credential-manager/prerequisites
			Relation: []string{
				"delegate_permission/common.handle_all_urls",
				"delegate_permission/common.get_login_creds",
			},
			Target: assetLinkTarget{
				Namespace:              "android_app",
				PackageName:            cfg.AndroidPackageName,
				SHA256CertFingerprints: normalizeSHA256Fingerprints(ctx, cfg.AndroidSHA256CertFingerprints),
			},
		})
	}
	return links
}

// normalizeSHA256Fingerprints returns colon-separated uppercase hex fingerprints
// for assetlinks.json. Normalization is applied at serve time only so stored /
// client state can keep the format clients sent.
func normalizeSHA256Fingerprints(ctx context.Context, fps []string) []string {
	out := make([]string, len(fps))
	for i, fp := range fps {
		out[i] = normalizeSHA256Fingerprint(ctx, fp)
	}
	return out
}

func normalizeSHA256Fingerprint(ctx context.Context, fp string) string {
	if strings.Contains(fp, ":") {
		return strings.ToUpper(fp) // already 32 valid pairs per API validation; only case may vary
	}
	raw, err := hex.DecodeString(fp) // guaranteed 64 hex chars
	if err != nil || len(raw) != sha256.Size {
		logging.Warn(ctx, "stored android sha256 fingerprint is not canonical hex; serving as-is",
			"fingerprint", fp, "err", err)
		return strings.ToUpper(fp)
	}
	parts := make([]string, len(raw))
	for i, b := range raw {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, ":")
}

func (h *Handler) writeJSON(w http.ResponseWriter, r *http.Request, v any) {
	w.Header().Set("Content-Type", "application/json")
	h.setCacheControl(w)
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logging.Error(r.Context(), "unable to encode well-known response", "err", err)
	}
}

// setCacheControl sets the Cache-Control header for the well-known app link files.
//
// The header intentionally omits the "public" directive and mirrors the JWKS
// endpoint (max-age + must-revalidate). In ZITADEL Cloud, the CDN is configured
// to cache only responses that explicitly opt in via "public" or "s-maxage".
// With "public", a fleet of independent edge caches could each pin
// a momentarily inconsistent (e.g. empty) copy and serve diverging content for
// the whole max-age, breaking native passkey verification. Without it, clients
// (the Apple and Google verifiers, browsers) still cache for max-age while the
// origin stays the single source of truth.
func (h *Handler) setCacheControl(w http.ResponseWriter) {
	if h.cacheControlMaxAge <= 0 {
		w.Header().Set(http_util.CacheControl, "no-store")
		return
	}
	w.Header().Set(http_util.CacheControl,
		fmt.Sprintf("max-age=%d, must-revalidate", int(h.cacheControlMaxAge/time.Second)),
	)
}
