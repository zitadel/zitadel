package integration

import (
	"bytes"
	"context"
	_ "embed"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/zitadel/logging"
	"sigs.k8s.io/yaml"
)

type Config struct {
	Log          *logging.Config
	Hostname     string
	Port         uint16
	Secure       bool
	LoginURLV2   string
	LogoutURLV2  string
	WebAuthNName string
}

//go:embed config/client.yaml
var clientYAML []byte

var (
	tmpDir       string
	loadedConfig Config
)

// TmpDir returns the absolute path to the project's temp directory.
func TmpDir() string {
	return tmpDir
}

// ServerAddr returns the physical address (host:port) of the ZITADEL server
// used by the integration test environment, e.g. "localhost:8082".
// Use this when establishing raw TCP/gRPC connections and set the instance
// domain via [grpc.WithAuthority] so that ZITADEL can route by virtual host.
func ServerAddr() string {
	return loadedConfig.Host()
}

func init() {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	tmpDir = filepath.Join(string(bytes.TrimSpace(out)), "tmp")

	if err := yaml.Unmarshal(clientYAML, &loadedConfig); err != nil {
		panic(err)
	}
	if err := loadedConfig.Log.SetLogger(); err != nil {
		panic(err)
	}
	SystemToken = createSystemUserToken()
	SystemUserWithNoPermissionsToken = createSystemUserWithNoPermissionsToken()

	// Override http.DefaultTransport so that HTTP requests to *.integration.localhost
	// are routed to the actual server host (loadedConfig.Hostname).
	// In testcontainers environments these custom instance domains are not DNS-resolvable,
	// while the server always listens on localhost. The Host header is preserved so
	// ZITADEL can still route requests to the correct instance by virtual host.
	http.DefaultTransport = newInstanceRoutingTransport(loadedConfig.Hostname)
}

// newInstanceRoutingTransport returns an http.RoundTripper that dials targetHost
// instead of the host in the URL when the URL host is *.integration.localhost.
func newInstanceRoutingTransport(targetHost string) http.RoundTripper {
	base := http.DefaultTransport.(*http.Transport).Clone()
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	base.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(host, ".integration.localhost") || host == "integration.localhost" {
			addr = net.JoinHostPort(targetHost, port)
		}
		return dialer.DialContext(ctx, network, addr)
	}
	return base
}
