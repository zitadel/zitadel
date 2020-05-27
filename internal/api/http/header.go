package http

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/caos/zitadel/internal/api"
)

type key int

var (
	httpHeaders key
	remoteAddr  key
)

func CopyHeadersToContext(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), httpHeaders, r.Header)
		ctx = context.WithValue(ctx, remoteAddr, r.RemoteAddr)
		r = r.WithContext(ctx)
		h(w, r)
	}
}

func HeadersFromCtx(ctx context.Context) (http.Header, bool) {
	headers, ok := ctx.Value(httpHeaders).(http.Header)
	return headers, ok
}

func RemoteIPFromCtx(ctx context.Context) string {
	ctxHeaders, ok := HeadersFromCtx(ctx)
	if !ok {
		return RemoteAddrFromCtx(ctx)
	}
	forwarded, ok := ForwardedFor(ctxHeaders)
	if ok {
		return forwarded
	}
	return RemoteAddrFromCtx(ctx)
}

func RemoteIPFromRequest(r *http.Request) net.IP {
	return net.ParseIP(RemoteIPStringFromRequest(r))
}

func RemoteIPStringFromRequest(r *http.Request) string {
	ip, ok := ForwardedFor(r.Header)
	if ok {
		return ip
	}
	return r.RemoteAddr
}

func ForwardedFor(headers http.Header) (string, bool) {
	forwarded, ok := headers[api.ForwardedFor]
	if ok {
		ip := strings.Split(forwarded[0], ", ")[0]
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
