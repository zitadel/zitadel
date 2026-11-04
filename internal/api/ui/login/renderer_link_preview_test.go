package login

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestLogin_applyLinkPreview(t *testing.T) {
	l := &Login{renderer: &Renderer{pathPrefix: HandlerPrefix}}
	requestWithOrigin := func() *http.Request {
		r := httptest.NewRequest(http.MethodGet, "/ui/login/login", nil)
		return r.WithContext(http_util.WithDomainContext(r.Context(), http_util.NewDomainCtx("auth.acme.com", "", "https")))
	}
	tests := []struct {
		name    string
		request *http.Request
		policy  *domain.LabelPolicy
		want    baseData
	}{
		{
			name:    "no label policy keeps the page title and description",
			request: requestWithOrigin(),
			want: baseData{
				Title:                "Login",
				Description:          "Sign in with ZITADEL",
				OpenGraphTitle:       "Login",
				OpenGraphDescription: "Sign in with ZITADEL",
				TwitterCard:          "summary",
			},
		},
		{
			name:    "empty link preview keeps the defaults",
			request: requestWithOrigin(),
			policy:  &domain.LabelPolicy{},
			want: baseData{
				Title:                "Login",
				Description:          "Sign in with ZITADEL",
				OpenGraphTitle:       "Login",
				OpenGraphDescription: "Sign in with ZITADEL",
				TwitterCard:          "summary",
			},
		},
		{
			name:    "custom values with an absolute image url",
			request: requestWithOrigin(),
			policy: &domain.LabelPolicy{
				LinkPreviewTitle:       "Acme",
				LinkPreviewDescription: "Sign in to Acme",
				LinkPreviewImageURL:    "policy/label/link_preview_image-123",
			},
			want: baseData{
				Title:                "Login",
				Description:          "Sign in to Acme",
				OpenGraphTitle:       "Acme",
				OpenGraphDescription: "Sign in to Acme",
				OpenGraphImage:       "https://auth.acme.com/ui/login/resources/dynamic?orgId=org1&default-policy=false&filename=policy/label/link_preview_image-123",
				TwitterCard:          "summary_large_image",
			},
		},
		{
			name:    "no image without a known origin",
			request: httptest.NewRequest(http.MethodGet, "/ui/login/login", nil),
			policy:  &domain.LabelPolicy{LinkPreviewImageURL: "policy/label/link_preview_image-123"},
			want: baseData{
				Title:                "Login",
				Description:          "Sign in with ZITADEL",
				OpenGraphTitle:       "Login",
				OpenGraphDescription: "Sign in with ZITADEL",
				TwitterCard:          "summary",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := baseData{
				Title:                "Login",
				Description:          "Sign in with ZITADEL",
				LabelPolicy:          tt.policy,
				PrivateLabelingOrgID: "org1",
			}
			got := l.applyLinkPreview(tt.request, data)
			assert.Equal(t, tt.want.Title, got.Title)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.OpenGraphTitle, got.OpenGraphTitle)
			assert.Equal(t, tt.want.OpenGraphDescription, got.OpenGraphDescription)
			assert.Equal(t, tt.want.OpenGraphImage, got.OpenGraphImage)
			assert.Equal(t, tt.want.TwitterCard, got.TwitterCard)
		})
	}
}

func TestLogin_applyLinkPreview_afterPageTextChange(t *testing.T) {
	l := &Login{renderer: &Renderer{pathPrefix: HandlerPrefix}}
	r := httptest.NewRequest(http.MethodGet, "/ui/login/mfa/verify", nil)

	t.Run("defaults follow the final page text", func(t *testing.T) {
		data := l.applyLinkPreview(r, baseData{})
		data.Title = "Verify TOTP"
		data.Description = "Enter the code"
		data = l.applyLinkPreview(r, data)
		assert.Equal(t, "Verify TOTP", data.OpenGraphTitle)
		assert.Equal(t, "Enter the code", data.OpenGraphDescription)
	})

	t.Run("custom values win over the final page text", func(t *testing.T) {
		policy := &domain.LabelPolicy{LinkPreviewTitle: "Acme", LinkPreviewDescription: "Sign in to Acme"}
		data := l.applyLinkPreview(r, baseData{LabelPolicy: policy})
		data.Title = "Verify TOTP"
		data.Description = "Enter the code"
		data = l.applyLinkPreview(r, data)
		assert.Equal(t, "Verify TOTP", data.Title)
		assert.Equal(t, "Sign in to Acme", data.Description)
		assert.Equal(t, "Acme", data.OpenGraphTitle)
		assert.Equal(t, "Sign in to Acme", data.OpenGraphDescription)
	})
}

func TestMainTemplate_linkPreviewTags(t *testing.T) {
	stub := func(...any) string { return "" }
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"resourceThemeUrl":    stub,
		"resourceUrl":         stub,
		"variablesCssFileUrl": stub,
		"customIconResource":  stub,
		"hasCustomPolicy":     func(*domain.LabelPolicy) bool { return false },
	}).ParseFiles("static/templates/main.html"))
	template.Must(tmpl.New("header").Parse(""))

	render := func(data baseData) string {
		var out bytes.Buffer
		require.NoError(t, tmpl.ExecuteTemplate(&out, "main-top", data))
		return out.String()
	}

	t.Run("custom values", func(t *testing.T) {
		got := render(baseData{
			Title:                "Login",
			Description:          "Sign in to Acme",
			OpenGraphTitle:       `Acme "Login"`,
			OpenGraphDescription: "Sign in to Acme",
			OpenGraphImage:       "https://auth.acme.com/ui/login/resources/dynamic?orgId=org1&filename=preview.png",
			TwitterCard:          "summary_large_image",
		})
		assert.Contains(t, got, `<title>Login</title>`)
		assert.Contains(t, got, `<meta property="og:title" content="Acme &#34;Login&#34;"/>`)
		assert.Contains(t, got, `<meta property="og:description" content="Sign in to Acme"/>`)
		assert.Contains(t, got, `<meta property="og:image" content="https://auth.acme.com/ui/login/resources/dynamic?orgId=org1&amp;filename=preview.png"/>`)
		assert.Contains(t, got, `<meta name="twitter:card" content="summary_large_image"/>`)
		assert.Contains(t, got, `<meta name="twitter:title" content="Acme &#34;Login&#34;"/>`)
		assert.Contains(t, got, `<meta name="twitter:image" content="https://auth.acme.com/ui/login/resources/dynamic?orgId=org1&amp;filename=preview.png"/>`)
		assert.Contains(t, got, `<meta name="robots" content="none" />`)
	})

	t.Run("no description and image", func(t *testing.T) {
		got := render(baseData{Title: "Login", OpenGraphTitle: "Login", TwitterCard: "summary"})
		assert.Contains(t, got, `<meta property="og:title" content="Login"/>`)
		assert.Contains(t, got, `<meta name="twitter:card" content="summary"/>`)
		assert.NotContains(t, got, "og:description")
		assert.NotContains(t, got, "og:image")
		assert.NotContains(t, got, "twitter:description")
		assert.NotContains(t, got, "twitter:image")
	})
}
