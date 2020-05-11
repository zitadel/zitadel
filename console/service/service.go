package service

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/caos/zitadel/console/config"
)

type spaHandler struct {
	dir       string
	indexFile string
}

func (s *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := filepath.Join(s.dir, filepath.Clean(r.URL.Path))

	if info, err := os.Stat(p); err != nil {
		http.ServeFile(w, r, filepath.Join(s.dir, s.indexFile))
		return
	} else if info.IsDir() {
		http.ServeFile(w, r, filepath.Join(s.dir, s.indexFile))
		return
	}

	http.ServeFile(w, r, p)
}

func Start(ctx context.Context, conf *config.Config) error {
	http.Handle("/", &spaHandler{conf.StaticDir, "index.html"})
	return http.ListenAndServe(":"+conf.Port, nil)
}
