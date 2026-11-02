package console

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testIndex = `<!doctype html>
<html lang="en">
  <head>
    <title>ZITADEL • Management Console</title>
    <meta property="description" content="Console Management Platform for ZITADEL instance"/>
    <meta property="og:title" content="ZITADEL Management Console"/>
    <meta property="og:description" content="Console Management Platform for ZITADEL IAM"/>
    <meta property="og:image" content="https://www.zitadel.com/images/preview.png"/>
    <meta name="twitter:card" content="summary"/>
    <meta name="twitter:site" content="@zitadel"/>
    <meta name="twitter:title" content="ZITADEL Management Console"/>
    <meta name="twitter:description" content="Management Platform for ZITADEL IAM"/>
    <meta name="twitter:image" content="https://www.zitadel.com/images/preview.png"/>
  </head>
  <body></body>
</html>`

func Test_applyLinkPreview(t *testing.T) {
	tests := []struct {
		name     string
		index    string
		preview  *LinkPreview
		contains []string
		excludes []string
	}{
		{
			name: "all values",
			preview: &LinkPreview{
				Title:       "Acme",
				Description: "Sign in to Acme",
				ImageURL:    "https://auth.acme.com/assets/v1/org/preview.png",
			},
			contains: []string{
				`<title>ZITADEL • Management Console</title>`,
				`<meta property="og:title" content="Acme"/>`,
				`<meta name="twitter:title" content="Acme"/>`,
				`<meta name="description" content="Sign in to Acme"/>`,
				`<meta property="og:description" content="Sign in to Acme"/>`,
				`<meta name="twitter:description" content="Sign in to Acme"/>`,
				`<meta property="og:image" content="https://auth.acme.com/assets/v1/org/preview.png"/>`,
				`<meta name="twitter:image" content="https://auth.acme.com/assets/v1/org/preview.png"/>`,
				`<meta name="twitter:card" content="summary_large_image"/>`,
			},
			excludes: []string{"www.zitadel.com/images/preview.png", `content="ZITADEL Management Console"`, "twitter:site"},
		},
		{
			name:    "title only keeps the default description and drops the default image",
			preview: &LinkPreview{Title: "Acme"},
			contains: []string{
				`<meta property="og:title" content="Acme"/>`,
				`<meta property="og:description" content="Console Management Platform for ZITADEL IAM"/>`,
				`<meta name="twitter:card" content="summary"/>`,
			},
			excludes: []string{"og:image", "twitter:image", "www.zitadel.com/images/preview.png", "twitter:site"},
		},
		{
			name:     "escapes html",
			preview:  &LinkPreview{Title: `"><script>alert(1)</script>`},
			contains: []string{`content="&#34;&gt;&lt;script&gt;alert(1)&lt;/script&gt;"`},
			excludes: []string{"<script>"},
		},
		{
			name:    "removes defaults with other attribute order and quotes",
			index:   "<html><head>\n<meta content='ZITADEL Management Console' property='og:title' />\n<META NAME=\"twitter:site\" content=\"@zitadel\">\n</head></html>",
			preview: &LinkPreview{Title: "Acme"},
			contains: []string{
				`<meta property="og:title" content="Acme"/>`,
				`<meta name="twitter:title" content="Acme"/>`,
			},
			excludes: []string{"ZITADEL Management Console", "twitter:site"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := testIndex
			if tt.index != "" {
				index = tt.index
			}
			got := string(applyLinkPreview([]byte(index), tt.preview))
			for _, want := range tt.contains {
				assert.Contains(t, got, want)
			}
			for _, unwanted := range tt.excludes {
				assert.NotContains(t, got, unwanted)
			}
		})
	}
}

// Test_applyLinkPreview_consoleIndex guards against changes to the console's index.html
// that would make the tag replacement silently fall back to the defaults.
func Test_applyLinkPreview_consoleIndex(t *testing.T) {
	pages := map[string]bool{
		"../../../../console/src/index.html": true,
		// the built pages only exist after building the console
		"../../../../console/dist/console/index.html": false,
		"static/index.html":                           false,
	}
	for page, required := range pages {
		t.Run(page, func(t *testing.T) {
			index, err := os.ReadFile(page)
			if !required && os.IsNotExist(err) {
				t.Skip("console is not built")
			}
			require.NoError(t, err)

			got := string(applyLinkPreview(index, &LinkPreview{
				Title:       "Acme",
				Description: "Sign in to Acme",
				ImageURL:    "https://auth.acme.com/assets/v1/org/preview.png",
			}))

			for _, key := range append(append(append([]string{}, titleMetaKeys...), descriptionMetaKeys...), imageMetaKeys...) {
				assert.Regexp(t, `(property|name)="`+regexp.QuoteMeta(key)+`" content="(Acme|Sign in to Acme|https://auth\.acme\.com/assets/v1/org/preview\.png)"`, got, key)
			}
			assert.Contains(t, got, `name="twitter:card" content="summary_large_image"`)
			assert.NotContains(t, got, `content="summary"`)
			assert.NotContains(t, got, "twitter:site")
			assert.NotContains(t, got, "www.zitadel.com/images/preview.png")
			assert.NotContains(t, got, `content="ZITADEL Management Console"`)
			assert.NotContains(t, got, "Platform for ZITADEL")
			assert.Contains(t, got, "<title>ZITADEL • Management Console</title>")

			titleOnly := string(applyLinkPreview(index, &LinkPreview{Title: "Acme"}))
			assert.Contains(t, titleOnly, `name="twitter:card" content="summary"`)
			assert.NotContains(t, titleOnly, "og:image")
			assert.NotContains(t, titleOnly, "twitter:image")
			assert.NotContains(t, titleOnly, "summary_large_image")
		})
	}
}

// Test_applyLinkPreview_customTagsFirst covers default tags the patterns don't match:
// the custom tags come first in the head, which is the tag unfurlers use.
func Test_applyLinkPreview_customTagsFirst(t *testing.T) {
	index := "<html><head>\n<meta data-x=\"1\" property=og:title content=Default>\n<title>Console</title>\n</head></html>"

	got := string(applyLinkPreview([]byte(index), &LinkPreview{Title: "Acme"}))

	custom := strings.Index(got, `<meta property="og:title" content="Acme"/>`)
	require.GreaterOrEqual(t, custom, 0)
	assert.Less(t, custom, strings.Index(got, "content=Default"))
	assert.Less(t, strings.Index(got, "<head>"), custom)
}

func Test_linkPreviewHandler(t *testing.T) {
	fSys := fstest.MapFS{
		"index.html":  {Data: []byte(testIndex)},
		"main.js":     {Data: []byte("console.log()")},
		"assets/a.js": {Data: []byte("a")},
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("static"))
	})
	passthrough := func(h http.Handler) http.Handler { return h }
	rejecting := func(http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "instance not found", http.StatusNotFound)
		})
	}
	acme := func(_ context.Context, orgID string) (*LinkPreview, error) {
		if orgID == "org1" {
			return &LinkPreview{Title: "Org One"}, nil
		}
		return &LinkPreview{Title: "Acme"}, nil
	}
	empty := func(context.Context, string) (*LinkPreview, error) { return &LinkPreview{}, nil }
	failing := func(context.Context, string) (*LinkPreview, error) { return nil, errors.New("db down") }

	tests := []struct {
		name            string
		instanceHandler func(http.Handler) http.Handler
		linkPreview     LinkPreviewFunc
		path            string
		wantStatus      int
		wantContains    string
	}{
		{name: "instance branding on root", instanceHandler: passthrough, linkPreview: acme, path: "/", wantStatus: http.StatusOK, wantContains: `content="Acme"`},
		{name: "instance branding on spa route", instanceHandler: passthrough, linkPreview: acme, path: "/users", wantStatus: http.StatusOK, wantContains: `content="Acme"`},
		{name: "organization branding", instanceHandler: passthrough, linkPreview: acme, path: "/?org=org1", wantStatus: http.StatusOK, wantContains: `content="Org One"`},
		{name: "static file untouched", instanceHandler: passthrough, linkPreview: acme, path: "/main.js", wantStatus: http.StatusOK, wantContains: "static"},
		{name: "no custom values", instanceHandler: passthrough, linkPreview: empty, path: "/", wantStatus: http.StatusOK, wantContains: "static"},
		{name: "query error falls back", instanceHandler: passthrough, linkPreview: failing, path: "/", wantStatus: http.StatusOK, wantContains: "static"},
		{name: "unknown instance falls back", instanceHandler: rejecting, linkPreview: acme, path: "/", wantStatus: http.StatusOK, wantContains: "static"},
		{name: "no link preview func", instanceHandler: passthrough, linkPreview: nil, path: "/", wantStatus: http.StatusOK, wantContains: "static"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := linkPreviewHandler(fSys, tt.instanceHandler, tt.linkPreview, next)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantContains)
		})
	}
}
