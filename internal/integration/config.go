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

var (
	//go:embed config/client.yaml
	clientYAML []byte
)

var (
	loadedConfig Config
)

func init() {
	if err := yaml.Unmarshal(clientYAML, &loadedConfig); err != nil {
		panic(err)
	}
	if err := loadedConfig.Log.SetLogger(); err != nil {
		panic(err)
	}
	SystemToken = systemUserToken()
}
