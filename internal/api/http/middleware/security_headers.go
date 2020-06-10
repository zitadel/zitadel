package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caos/zitadel/internal/api"
)

type CSP struct {
	DefaultSrc CSPSourceOptions
	ScriptSrc  CSPSourceOptions
	ObjectSrc  CSPSourceOptions
	StyleSrc   CSPSourceOptions
	ImgSrc     CSPSourceOptions
	MediaSrc   CSPSourceOptions
	FrameSrc   CSPSourceOptions
	FontSrc    CSPSourceOptions
	ConnectSrc CSPSourceOptions
	FormAction CSPSourceOptions
}

func (csp *CSP) Value() string {
	valuesMap := csp.asMap()

	values := make([]string, 0, len(valuesMap))
	for k, v := range valuesMap {
		if v == nil {
			continue
		}

		values = append(values, fmt.Sprintf("%v %v", k, v.String()))
	}

	return strings.Join(values, ";")
}

func (csp *CSP) asMap() map[string]CSPSourceOptions {
	return map[string]CSPSourceOptions{
		"default-src": csp.DefaultSrc,
		"script-src":  csp.ScriptSrc,
		"object-src":  csp.ObjectSrc,
		"style-src":   csp.StyleSrc,
		"img-src":     csp.ImgSrc,
		"media-src":   csp.MediaSrc,
		"frame-src":   csp.FrameSrc,
		"font-src":    csp.FontSrc,
		"connect-src": csp.ConnectSrc,
		"form-action": csp.FormAction,
	}
}

type CSPSourceOptions []string

func CSPSourceOpts() CSPSourceOptions {
	return CSPSourceOptions{}
}

func CSPSourceOptsNone() CSPSourceOptions {
	return []string{"'none'"}
}

func CSPSourceOptsSelf() CSPSourceOptions {
	return []string{"'self'"}
}

func (srcOpts CSPSourceOptions) AddSelf() CSPSourceOptions {
	return append(srcOpts, "'self'")
}

func (srcOpts CSPSourceOptions) AddInline() CSPSourceOptions {
	return append(srcOpts, "'unsafe-inline'")
}

func (srcOpts CSPSourceOptions) AddEval() CSPSourceOptions {
	return append(srcOpts, "'unsafe-eval'")
}

func (srcOpts CSPSourceOptions) AddStrictDynamic() CSPSourceOptions {
	return append(srcOpts, "'strict-dynamic'")
}

func (srcOpts CSPSourceOptions) AddHost(h ...string) CSPSourceOptions {
	return append(srcOpts, h...)
}

func (srcOpts CSPSourceOptions) AddScheme(s ...string) CSPSourceOptions {
	return srcOpts.add(s, "%v:")
}

func (srcOpts CSPSourceOptions) AddNonce(b64n ...string) CSPSourceOptions {
	return srcOpts.add(b64n, "'nonce-%v'")
}

func (srcOpts CSPSourceOptions) AddHash(alg, b64v string) CSPSourceOptions {
	return append(srcOpts, fmt.Sprintf("'%v-%v'", alg, b64v))
}

func (srcOpts CSPSourceOptions) String() string {
	return strings.Join(srcOpts, " ")
}

func (srcOpts CSPSourceOptions) add(values []string, format string) CSPSourceOptions {
	for i, v := range values {
		values[i] = fmt.Sprintf(format, v)
	}

	return append(srcOpts, values...)
}

var (
	DefaultSCP = &CSP{
		DefaultSrc: CSPSourceOptsSelf(),
		ScriptSrc:  CSPSourceOptsSelf(),
		ObjectSrc:  CSPSourceOptsNone(),
		StyleSrc:   CSPSourceOptsSelf().AddScheme("data"),
		ImgSrc:     CSPSourceOptsSelf(),
		MediaSrc:   CSPSourceOptsNone(),
		FrameSrc:   CSPSourceOptsNone(),
		FontSrc:    CSPSourceOptsSelf(),
		ConnectSrc: CSPSourceOptsSelf(),
		FormAction: CSPSourceOptsSelf(),
	}
)

func SecurityHeaders(csp *CSP) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			headers := w.Header()
			if csp == nil {
				csp = DefaultSCP
			}
			headers.Set(api.ContentSecurityPolicy, csp.Value())
			headers.Set(api.XXSSProtection, "1; mode=block")
			headers.Set(api.StrictTransportSecurity, "max-age=31536000; includeSubDomains")
			headers.Set(api.XFrameOptions, "DENY")
			headers.Set(api.XContentTypeOptions, "nosniff")
			headers.Set(api.ReferrerPolicy, "same-origin")
			headers.Set(api.FeaturePolicy, "payment 'none'")
			//PLANNED: add expect-ct

			handler.ServeHTTP(w, req)
		})
	}
}
