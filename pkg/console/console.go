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
	EnvironmentPath string
}

type spaHandler struct {
	fileSystem http.FileSystem
}

const (
	envRequestPath = "/assets/environment.json"
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
	envPath := envRequestPath
	if config.EnvironmentPath != "" {
		envPath = config.EnvironmentPath
	}
	http.Handle("/", http.FileServer(&spaHandler{statikFS}))
	http.Handle(envRequestPath, http.FileServer(http.Dir(envPath)))
	return http.ListenAndServe(":"+config.Port, nil)
}
