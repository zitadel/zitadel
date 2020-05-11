//go:generate /Users/ffo/go/bin/statik -src=../../console/dist/app

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
	Port string
}

type spaHandler struct {
	fileSystem http.FileSystem
}

func (i *spaHandler) Open(name string) (http.File, error) {
	ret, err := i.fileSystem.Open(name)
	if !os.IsNotExist(err) || path.Ext(name) != "" {
		return ret, err
	}

	return i.fileSystem.Open("/index.html")
}

func Start(ctx context.Context, config Config) error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}
	http.Handle("/", http.FileServer(&spaHandler{statikFS}))
	return http.ListenAndServe(":"+config.Port, nil)
}
