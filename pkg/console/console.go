package console

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	Port      string
	StaticDir string
}

type spaHandler struct {
	dir       string
	indexFile string
}

func Start(ctx context.Context, config Config) error {
	http.Handle("/", &spaHandler{config.StaticDir, "index.html"})
	return http.ListenAndServe(":"+config.Port, nil)
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
