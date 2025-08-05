package integration

import (
	"bytes"
	_ "embed"
	"os/exec"
	"path/filepath"

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
	if err := loadedConfig.Log.SetLogger(); err != nil {
		panic(err)
	}
	SystemToken = createSystemUserToken()
	SystemUserWithNoPermissionsToken = createSystemUserWithNoPermissionsToken()
}
