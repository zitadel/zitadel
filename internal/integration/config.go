package integration

import (
	"bytes"
	_ "embed"
	"log/slog"
	"os/exec"
	"path/filepath"

	"github.com/zitadel/logging"
	"sigs.k8s.io/yaml"

	"github.com/zitadel/zitadel/cmd/build"
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

var (
	//go:embed config/client.yaml
	clientYAML []byte
)

var (
	tmpDir       string
	loadedConfig Config
)

// TmpDir returns the absolute path to the projects's temp directory.
func TmpDir() string {
	return tmpDir
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

	loadedConfig.Log.Formatter.Data = map[string]interface{}{
		"service": "zitadel",
		"version": build.Version(),
	}

	slog.SetDefault(loadedConfig.Log.Slog())

	SystemToken = systemUserToken()
}
