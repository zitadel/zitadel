package middleware

import (
	"fmt"
	"strings"
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

var (
	DefaultSCP = CSP{
		DefaultSrc: CSPSourceOptsNone(),
		ScriptSrc:  CSPSourceOptsSelf(),
		ObjectSrc:  CSPSourceOptsNone(),
		StyleSrc:   CSPSourceOptsSelf(),
		ImgSrc:     CSPSourceOptsSelf(),
		MediaSrc:   CSPSourceOptsNone(),
		FrameSrc:   CSPSourceOptsNone(),
		FontSrc:    CSPSourceOptsSelf(),
		ConnectSrc: CSPSourceOptsSelf(),
	}
)

func (csp *CSP) Value(nonce string) string {
	valuesMap := csp.asMap()

	values := make([]string, 0, len(valuesMap))
	for k, v := range valuesMap {
		if v == nil {
			continue
		}

		values = append(values, fmt.Sprintf("%v %v", k, v.String(nonce)))
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

func (srcOpts *CSPSourceOptions) AddSelf() *CSPSourceOptions {
	*srcOpts = append(*srcOpts, "'self'")
	return srcOpts
}

func (srcOpts *CSPSourceOptions) AddInline() *CSPSourceOptions {
	*srcOpts = append(*srcOpts, "'unsafe-inline'")
	return srcOpts
}

func (srcOpts *CSPSourceOptions) AddEval() *CSPSourceOptions {
	*srcOpts = append(*srcOpts, "'unsafe-eval'")
	return srcOpts
}

func (srcOpts *CSPSourceOptions) AddStrictDynamic() *CSPSourceOptions {
	*srcOpts = append(*srcOpts, "'strict-dynamic'")
	return srcOpts
}

func (srcOpts *CSPSourceOptions) AddHost(h ...string) *CSPSourceOptions {
	*srcOpts = append(*srcOpts, h...)
	return srcOpts
}

func (srcOpts *CSPSourceOptions) AddScheme(s ...string) *CSPSourceOptions {
	return srcOpts.add(s, "%v:")
}

func (srcOpts *CSPSourceOptions) AddNonce() *CSPSourceOptions {
	*srcOpts = append(*srcOpts, "'nonce-%v'")
	return srcOpts
}

func (srcOpts *CSPSourceOptions) AddHash(alg, b64v string) *CSPSourceOptions {
	*srcOpts = append(*srcOpts, fmt.Sprintf("'%v-%v'", alg, b64v))
	return srcOpts
}

func (srcOpts *CSPSourceOptions) String(nonce string) string {
	value := strings.Join(*srcOpts, " ")
	if !strings.Contains(value, "%v") {
		return value
	}
	return fmt.Sprintf(value, nonce)
}

func (srcOpts *CSPSourceOptions) add(values []string, format string) *CSPSourceOptions {
	for i, v := range values {
		values[i] = fmt.Sprintf(format, v)
	}

	*srcOpts = append(*srcOpts, values...)
	return srcOpts
}
