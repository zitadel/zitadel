package http

import (
	"context"
	"net"
	"net/http"
	"strings"
)

const (
	Authorization   = "authorization"
	Accept          = "accept"
	AcceptLanguage  = "accept-language"
	CacheControl    = "cache-control"
	ContentType     = "content-type"
	ContentLength   = "content-length"
	Expires         = "expires"
	Location        = "location"
	Origin          = "origin"
	Pragma          = "pragma"
	UserAgentHeader = "user-agent"
	ForwardedFor    = "x-forwarded-for"
	XUserAgent      = "x-user-agent"
	XGrpcWeb        = "x-grpc-web"
	XRequestedWith  = "x-requested-with"
	XRobotsTag      = "x-robots-tag"
	IfNoneMatch     = "If-None-Match"
	LastModified    = "Last-Modified"
	Etag            = "Etag"

	ContentSecurityPolicy   = "content-security-policy"
	XXSSProtection          = "x-xss-protection"
	StrictTransportSecurity = "strict-transport-security"
	XFrameOptions           = "x-frame-options"
	XContentTypeOptions     = "x-content-type-options"
	ReferrerPolicy          = "referrer-policy"
	FeaturePolicy           = "feature-policy"
	PermissionsPolicy       = "permissions-policy"

	ZitadelOrgID = "x-zitadel-orgid"
)

type key int

const (
	httpHeaders key = iota
	remoteAddr
	origin
)

func CopyHeadersToContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), httpHeaders, r.Header)
		ctx = context.WithValue(ctx, remoteAddr, r.RemoteAddr)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func HeadersFromCtx(ctx context.Context) (http.Header, bool) {
	headers, ok := ctx.Value(httpHeaders).(http.Header)
	return headers, ok
}

func OriginHeader(ctx context.Context) string {
	headers, ok := ctx.Value(httpHeaders).(http.Header)
	if !ok {
		return ""
	}
	return headers.Get(Origin)
}

func ComposedOrigin(ctx context.Context) string {
	o, ok := ctx.Value(origin).(string)
	if !ok {
		return ""
	}
	return o
}

func WithComposedOrigin(ctx context.Context, composed string) context.Context {
	return context.WithValue(ctx, origin, composed)
}

func RemoteIPFromCtx(ctx context.Context) string {
	ctxHeaders, ok := HeadersFromCtx(ctx)
	if !ok {
		return RemoteAddrFromCtx(ctx)
	}
	forwarded, ok := GetForwardedFor(ctxHeaders)
	if ok {
		return forwarded
	}
	return RemoteAddrFromCtx(ctx)
}

func RemoteIPFromRequest(r *http.Request) net.IP {
	return net.ParseIP(RemoteIPStringFromRequest(r))
}

func RemoteIPStringFromRequest(r *http.Request) string {
	ip, ok := GetForwardedFor(r.Header)
	if ok {
		return ip
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func GetAuthorization(r *http.Request) string {
	return r.Header.Get(Authorization)
}

func GetOrgID(r *http.Request) string {
	return r.Header.Get(ZitadelOrgID)
}

func GetForwardedFor(headers http.Header) (string, bool) {
	forwarded, ok := headers[http.CanonicalHeaderKey(ForwardedFor)]
	if ok {
		ip := strings.TrimSpace(strings.Split(forwarded[0], ",")[0])
		if ip != "" {
			return ip, true
		}
	}
	return "", false
}

func RemoteAddrFromCtx(ctx context.Context) string {
	ctxRemoteAddr, _ := ctx.Value(remoteAddr).(string)
	return ctxRemoteAddr
}
