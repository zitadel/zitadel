package well_known

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/zitadel/logging"

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
		logging.WithError(err).Error("unable to query OIDC app link configs for apple-app-site-association")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, buildAppleAppSiteAssociation(configs))
}

func (h *Handler) serveAssetLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	configs, err := h.queries.SearchOIDCAppLinkConfigs(r.Context())
	if err != nil {
		logging.WithError(err).Error("unable to query OIDC app link configs for assetlinks.json")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, buildAssetLinks(configs))
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

func buildAssetLinks(configs []*query.OIDCAppLinkConfig) []assetLink {
	links := make([]assetLink, 0)
	for _, cfg := range configs {
		if cfg.AndroidPackageName == "" || len(cfg.AndroidSHA256CertFingerprints) == 0 {
			continue
		}
		links = append(links, assetLink{
			Relation: []string{"delegate_permission/common.get_login_creds"},
			Target: assetLinkTarget{
				Namespace:              "android_app",
				PackageName:            cfg.AndroidPackageName,
				SHA256CertFingerprints: normalizeSHA256Fingerprints(cfg.AndroidSHA256CertFingerprints),
			},
		})
	}
	return links
}

// normalizeSHA256Fingerprints returns colon-separated uppercase hex fingerprints
// for assetlinks.json. Normalization is applied at serve time only so stored /
// client state can keep the format clients sent.
func normalizeSHA256Fingerprints(fps []string) []string {
	out := make([]string, len(fps))
	for i, fp := range fps {
		out[i] = normalizeSHA256Fingerprint(fp)
	}
	return out
}

func normalizeSHA256Fingerprint(fp string) string {
	cleaned := make([]byte, 0, 64)
	for _, r := range fp {
		switch {
		case r == ':' || unicode.IsSpace(r):
			continue
		case r >= 'a' && r <= 'f':
			cleaned = append(cleaned, byte(r-32))
		case (r >= '0' && r <= '9') || (r >= 'A' && r <= 'F'):
			cleaned = append(cleaned, byte(r))
		default:
			return strings.ToUpper(strings.TrimSpace(fp))
		}
	}
	if len(cleaned) != 64 {
		return strings.ToUpper(strings.TrimSpace(fp))
	}
	var b strings.Builder
	b.Grow(95) // 64 hex chars + 31 colons
	for i := 0; i < 64; i += 2 {
		if i > 0 {
			b.WriteByte(':')
		}
		b.Write(cleaned[i : i+2])
	}
	return b.String()
}

func (h *Handler) writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	h.setCacheControl(w)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logging.WithError(err).Error("unable to encode well-known response")
	}
}

func (h *Handler) setCacheControl(w http.ResponseWriter) {
	if h.cacheControlMaxAge == 0 {
		w.Header().Set(http_util.CacheControl, "no-store")
		return
	}
	w.Header().Set(http_util.CacheControl,
		fmt.Sprintf("public, max-age=%d", int(h.cacheControlMaxAge/time.Second)),
	)
}
