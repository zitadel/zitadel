package denylist

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"syscall"
	"time"

	internal_net "github.com/zitadel/zitadel/internal/net"
)

// lookupIP resolves the hostnames of HTTPS targets sent through a proxy. It is a variable so
// tests can resolve names without DNS.
var lookupIP internal_net.IPLookupFunc = net.LookupIP

// NewHTTPTransport returns a cloned default transport that enforces denylist checks
// right before each TCP dial to avoid DNS rebinding TOCTOU gaps.
//
// When a request goes through a proxy, a plain dial only reaches the proxy, so the dial-time
// check never sees the target. HTTPS requests through an HTTP or HTTPS proxy are therefore
// tunnelled by the transport itself: it resolves the target, checks every address against the
// denylist and asks the proxy to CONNECT to the vetted address rather than to the hostname, so
// the proxy cannot resolve the name differently. Plain HTTP requests and SOCKS proxies receive
// the hostname at the proxy and cannot be pinned; they are checked against the denylist by URL
// before the proxy is contacted.
func NewHTTPTransport(denyList []AddressChecker) http.RoundTripper {
	return newHTTPTransport(denyList, http.DefaultTransport.(*http.Transport).Clone())
}

// newHTTPTransport applies the denylist checks to base.
func newHTTPTransport(denyList []AddressChecker, base *http.Transport) http.RoundTripper {
	if len(denyList) == 0 {
		return base
	}
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
		Control: func(network, address string, c syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return err
			}

			parsedIP := net.ParseIP(host)
			if parsedIP == nil { // at this point, it must be an IP so it should never happen
				return &net.DNSError{Err: "invalid IP address", Name: host}
			}

			return IsHostBlocked(denyList, host, parsedIP)
		},
	}
	base.DialContext = dialer.DialContext
	if base.Proxy == nil {
		return base
	}
	tunnels := &pinnedProxyTunnels{
		denyList:      denyList,
		dialer:        dialer,
		proxy:         base.Proxy,
		connectHeader: base.ProxyConnectHeader,
		tlsConfig:     base.TLSClientConfig,
	}
	base.Proxy = tunnels.selectProxy
	base.DialContext = tunnels.dialContext
	return base
}

// pinnedProxyTunnels sends HTTPS requests through a proxy over a tunnel to the vetted target
// address. selectProxy runs before the connection is dialled: for an HTTPS target it records
// the proxy and lets the transport dial "directly", and dialContext then opens the tunnel to
// the vetted address through the recorded proxy. The transport runs TLS over the returned
// connection with the target hostname, so certificate verification is unchanged.
type pinnedProxyTunnels struct {
	denyList      []AddressChecker
	dialer        *net.Dialer
	proxy         func(*http.Request) (*url.URL, error)
	connectHeader http.Header
	tlsConfig     *tls.Config
	pinned        sync.Map // target "host:port" -> *url.URL of the proxy to tunnel through
}

func (t *pinnedProxyTunnels) selectProxy(req *http.Request) (*url.URL, error) {
	proxyURL, err := t.proxy(req)
	if err != nil || proxyURL == nil {
		return proxyURL, err
	}
	if req.URL.Scheme == "https" && (proxyURL.Scheme == "http" || proxyURL.Scheme == "https") {
		t.pinned.Store(targetAddress(req.URL), proxyURL)
		return nil, nil
	}
	// The proxy receives the hostname and resolves it itself, so the target can only be
	// checked by URL here.
	if err := IsURLBlocked(t.denyList, req.URL, nil); err != nil {
		return nil, err
	}
	return proxyURL, nil
}

func (t *pinnedProxyTunnels) dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	proxyURL, ok := t.pinned.Load(address)
	if !ok {
		return t.dialer.DialContext(ctx, network, address)
	}
	return t.tunnel(ctx, network, address, proxyURL.(*url.URL))
}

// tunnel resolves address, checks every resolved IP against the denylist and opens a CONNECT
// tunnel through proxyURL to the first of them.
func (t *pinnedProxyTunnels) tunnel(ctx context.Context, network, address string, proxyURL *url.URL) (_ net.Conn, err error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	ips, err := internal_net.HostnameToIPList(host, lookupIP)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, &net.DNSError{Err: "no addresses", Name: host}
	}
	if err := IsHostBlocked(t.denyList, host, ips...); err != nil {
		return nil, err
	}
	target := net.JoinHostPort(ips[0].String(), port)

	conn, err := t.dialer.DialContext(ctx, network, proxyAddress(proxyURL))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
		defer func() { _ = conn.SetDeadline(time.Time{}) }()
	}
	if proxyURL.Scheme == "https" {
		tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
		if t.tlsConfig != nil {
			tlsConfig = t.tlsConfig.Clone()
		}
		tlsConfig.ServerName = proxyURL.Hostname()
		tlsConn := tls.Client(conn, tlsConfig)
		if err = tlsConn.HandshakeContext(ctx); err != nil {
			return nil, err
		}
		conn = tlsConn
	}

	connectReq := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Opaque: target},
		Host:   target,
		Header: t.connectHeader.Clone(),
	}
	if connectReq.Header == nil {
		connectReq.Header = make(http.Header)
	}
	if user := proxyURL.User; user != nil {
		password, _ := user.Password()
		connectReq.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user.Username()+":"+password)))
	}
	if err = connectReq.Write(conn); err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, connectReq)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proxy refused CONNECT to %s: %s", target, resp.Status)
	}
	if reader.Buffered() > 0 {
		return nil, fmt.Errorf("proxy sent unexpected data after CONNECT to %s", target)
	}
	return conn, nil
}

// targetAddress returns the "host:port" the transport dials for u.
func targetAddress(u *url.URL) string {
	port := u.Port()
	if port == "" {
		port = "443"
		if u.Scheme == "http" {
			port = "80"
		}
	}
	return net.JoinHostPort(u.Hostname(), port)
}

// proxyAddress returns the "host:port" of proxyURL.
func proxyAddress(proxyURL *url.URL) string {
	port := proxyURL.Port()
	if port == "" {
		port = "80"
		if proxyURL.Scheme == "https" {
			port = "443"
		}
	}
	return net.JoinHostPort(proxyURL.Hostname(), port)
}
