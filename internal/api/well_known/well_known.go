package well_known

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zitadel/logging"

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

type appLinkQuerier interface {
	SearchOIDCAppLinkConfigs(ctx context.Context) ([]*query.OIDCAppLinkConfig, error)
}

type Handler struct {
	queries appLinkQuerier
}

func NewHandler(queries appLinkQuerier) *Handler {
	return &Handler{queries: queries}
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
	writeJSON(w, buildAppleAppSiteAssociation(configs))
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
	writeJSON(w, buildAssetLinks(configs))
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
		fps := cfg.AndroidSHA256CertFingerprints
		if fps == nil {
			fps = []string{}
		}
		links = append(links, assetLink{
			Relation: []string{"delegate_permission/common.get_login_creds"},
			Target: assetLinkTarget{
				Namespace:              "android_app",
				PackageName:            cfg.AndroidPackageName,
				SHA256CertFingerprints: fps,
			},
		})
	}
	return links
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logging.WithError(err).Error("unable to encode well-known response")
	}
}
