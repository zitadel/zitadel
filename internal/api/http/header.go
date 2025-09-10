package http

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	Authorization          = "authorization"
	Accept                 = "accept"
	AcceptLanguage         = "accept-language"
	CacheControl           = "cache-control"
	ContentType            = "content-type"
	ContentLength          = "content-length"
	ContentLocation        = "content-location"
	Expires                = "expires"
	Location               = "location"
	Origin                 = "origin"
	Pragma                 = "pragma"
	UserAgentHeader        = "user-agent"
	ForwardedFor           = "x-forwarded-for"
	ForwardedHost          = "x-forwarded-host"
	ForwardedProto         = "x-forwarded-proto"
	Forwarded              = "forwarded"
	Host                   = "host"
	ZitadelForwarded       = "x-zitadel-forwarded"
	XUserAgent             = "x-user-agent"
	XGrpcWeb               = "x-grpc-web"
	XRequestedWith         = "x-requested-with"
	XRobotsTag             = "x-robots-tag"
	IfNoneMatch            = "if-none-match"
	LastModified           = "last-modified"
	Etag                   = "etag"
	GRPCTimeout            = "grpc-timeout"
	ConnectProtocolVersion = "connect-protocol-version"
	ConnectTimeoutMS       = "connect-timeout-ms"
	GrpcStatus             = "grpc-status"
	GrpcMessage            = "grpc-message"
	GrpcStatusDetailsBin   = "grpc-status-details-bin"

	ContentSecurityPolicy   = "content-security-policy"
	XXSSProtection          = "x-xss-protection"
	StrictTransportSecurity = "strict-transport-security"
	XFrameOptions           = "x-frame-options"
	XContentTypeOptions     = "x-content-type-options"
	ReferrerPolicy          = "referrer-policy"
	FeaturePolicy           = "feature-policy"
	PermissionsPolicy       = "permissions-policy"

	ZitadelOrgID = "x-zitadel-orgid"

	OrgIdInPathVariableName = "orgId"
	OrgIdInPathVariable     = "{" + OrgIdInPathVariableName + "}"
)

type key int

const (
	httpHeaders key = iota
	remoteAddr
	domainCtx
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
	// path variable takes precedence over header
	orgID, ok := mux.Vars(r)[OrgIdInPathVariableName]
	if ok {
		return orgID
	}

	return r.Header.Get(ZitadelOrgID)
}

func GetForwardedFor(headers http.Header) (string, bool) {
	forwarded := strings.Split(headers.Get(ForwardedFor), ",")[0]
	return forwarded, forwarded != ""
}

func RemoteAddrFromCtx(ctx context.Context) string {
	ctxRemoteAddr, _ := ctx.Value(remoteAddr).(string)
	return ctxRemoteAddr
}
