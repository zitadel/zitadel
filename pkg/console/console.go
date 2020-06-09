package console

import (
	"context"
	"net/http"
	"os"
	"path"

	"github.com/rakyll/statik/fs"

	_ "github.com/caos/zitadel/pkg/console/statik"
)

type Config struct {
	Port            string
	EnvOverwriteDir string
}

type spaHandler struct {
	fileSystem http.FileSystem
}

const (
	envRequestPath = "/assets/environment.json"
	envDefaultDir  = "/console"
)

func (i *spaHandler) Open(name string) (http.File, error) {
	ret, err := i.fileSystem.Open(name)
	if !os.IsNotExist(err) || path.Ext(name) != "" {
		return ret, err
	}

	return i.fileSystem.Open("/index.html")
}

func Start(ctx context.Context, config Config) error {
	statikFS, err := fs.NewWithNamespace("console")
	if err != nil {
		return err
	}
	envDir := envDefaultDir
	if config.EnvOverwriteDir != "" {
		envDir = config.EnvOverwriteDir
	}
	http.Handle("/", http.FileServer(&spaHandler{statikFS}))
	http.Handle(envRequestPath, http.StripPrefix("/assets", http.FileServer(http.Dir(envDir))))
	return http.ListenAndServe(":"+config.Port, nil)
}
