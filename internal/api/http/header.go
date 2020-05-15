package http

import (
	"context"
	"net/http"
	"strings"
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
	forwarded, ok := ctxHeaders["X-Forwarded-For"]
	if ok {
		ip := strings.Split(forwarded[0], ", ")[0]
		if ip != "" {
			return ip
		}
	}
	return RemoteAddrFromCtx(ctx)
}

func RemoteAddrFromCtx(ctx context.Context) string {
	ctxRemoteAddr, _ := ctx.Value(remoteAddr).(string)
	return ctxRemoteAddr
}
