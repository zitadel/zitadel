package console

import (
	"bytes"
	"context"
	"html/template"
	"io/fs"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/zitadel/logging"
)

// LinkPreview holds the branding values used for the Open Graph and Twitter tags of the console page.
// Empty values keep the default tags of the console build.
type LinkPreview struct {
	Title       string
	Description string
	// ImageURL must be absolute so link unfurlers can fetch it.
	ImageURL string
}

func (p *LinkPreview) isEmpty() bool {
	return p == nil || (p.Title == "" && p.Description == "" && p.ImageURL == "")
}

// LinkPreviewFunc returns the active link preview branding of the instance in the context,
// or of the organization if orgID is set and the organization has a custom branding.
type LinkPreviewFunc func(ctx context.Context, orgID string) (*LinkPreview, error)

const consoleOrgQueryParam = "org"

var (
	titleMetaKeys       = []string{"og:title", "twitter:title"}
	descriptionMetaKeys = []string{"description", "og:description", "twitter:description"}
	imageMetaKeys       = []string{"og:image", "twitter:image"}
	twitterCardMetaKey  = "twitter:card"
	twitterSiteMetaKey  = "twitter:site"

	metaTagPatterns = compileMetaTagPatterns(titleMetaKeys, descriptionMetaKeys, imageMetaKeys, []string{twitterCardMetaKey, twitterSiteMetaKey})
	headTagPattern  = regexp.MustCompile(`(?i)<head\b[^>]*>`)
)

// compileMetaTagPatterns matches the meta tags of each key independent of attribute order and quoting,
// so a differently formatted console build still has its default tags removed.
func compileMetaTagPatterns(keyGroups ...[]string) map[string]*regexp.Regexp {
	patterns := make(map[string]*regexp.Regexp)
	for _, keys := range keyGroups {
		for _, key := range keys {
			patterns[key] = regexp.MustCompile(`(?i)[ \t]*<meta\b[^>]*\b(?:property|name)\s*=\s*["']` + regexp.QuoteMeta(key) + `["'][^>]*>[ \t]*\r?\n?`)
		}
	}
	return patterns
}

// linkPreviewHandler serves index.html with the link preview tags of the requested instance or organization.
// It delegates to next for every other file and whenever no custom link preview is configured.
func linkPreviewHandler(fSys fs.FS, instanceHandler func(http.Handler) http.Handler, linkPreview LinkPreviewFunc, next http.Handler) http.Handler {
	index, err := fs.ReadFile(fSys, "index.html")
	if linkPreview == nil || err != nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !servesIndex(fSys, r) {
			next.ServeHTTP(w, r)
			return
		}
		instanceRequest := requestWithInstance(instanceHandler, r)
		if instanceRequest == nil {
			next.ServeHTTP(w, r)
			return
		}
		preview, err := linkPreview(instanceRequest.Context(), r.URL.Query().Get(consoleOrgQueryParam))
		if err != nil {
			logging.WithError(err).Warn("unable to get link preview branding for the management console")
			next.ServeHTTP(w, r)
			return
		}
		if preview.isEmpty() {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodHead {
			return
		}
		_, err = w.Write(applyLinkPreview(index, preview))
		logging.OnError(err).Error("error serving management console index")
	})
}

// servesIndex reports whether the file server would answer the request with index.html.
func servesIndex(fSys fs.FS, r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	if r.URL.Path == "/" || r.URL.Path == "" {
		return true
	}
	if !isIndexOrSubPath(r.URL.Path) {
		return false
	}
	_, err := fs.Stat(fSys, strings.TrimPrefix(r.URL.Path, "/"))
	return err != nil
}

// requestWithInstance resolves the instance of the request without writing an error response,
// so the console still loads for hosts the instance handler rejects.
func requestWithInstance(instanceHandler func(http.Handler) http.Handler, r *http.Request) *http.Request {
	var resolved *http.Request
	instanceHandler(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		resolved = req
	})).ServeHTTP(discardResponseWriter{}, r)
	return resolved
}

type discardResponseWriter struct{}

func (discardResponseWriter) Header() http.Header         { return http.Header{} }
func (discardResponseWriter) Write(b []byte) (int, error) { return len(b), nil }
func (discardResponseWriter) WriteHeader(int)             {}

// applyLinkPreview removes the default tags of every configured value and inserts the custom tags
// at the start of the head. Unfurlers use the first matching tag, so the custom values win
// even if a default tag of a future console build is not matched.
func applyLinkPreview(index []byte, preview *LinkPreview) []byte {
	type metaTag struct{ key, value string }
	// the ZITADEL account must not be attributed as the site of a custom branded preview
	removeKeys := []string{twitterSiteMetaKey}
	tags := make([]metaTag, 0, 8)
	add := func(keys []string, value string) {
		for _, key := range keys {
			removeKeys = append(removeKeys, key)
			tags = append(tags, metaTag{key, value})
		}
	}
	if preview.Title != "" {
		add(titleMetaKeys, preview.Title)
	}
	if preview.Description != "" {
		add(descriptionMetaKeys, preview.Description)
	}
	if preview.ImageURL != "" {
		add(imageMetaKeys, preview.ImageURL)
		add([]string{twitterCardMetaKey}, "summary_large_image")
	} else {
		// like the hosted login, a custom preview without an image must not show the ZITADEL image
		removeKeys = append(removeKeys, imageMetaKeys...)
		add([]string{twitterCardMetaKey}, "summary")
	}

	out := index
	for _, key := range removeKeys {
		out = metaTagPatterns[key].ReplaceAll(out, nil)
	}
	head := headTagPattern.FindIndex(out)
	if head == nil {
		return out
	}
	var inserted bytes.Buffer
	for _, tag := range tags {
		attr := "property"
		if strings.HasPrefix(tag.key, "twitter:") || tag.key == "description" {
			attr = "name"
		}
		inserted.WriteString("\n    <meta " + attr + `="` + tag.key + `" content="` + template.HTMLEscapeString(tag.value) + `"/>`)
	}
	return slices.Concat(out[:head[1]], inserted.Bytes(), out[head[1]:])
}
