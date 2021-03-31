package openapi

import (
	"net/http"

	"github.com/rakyll/statik/fs"
)

func Start() (http.Handler, error) {
	statikFS, err := fs.NewWithNamespace("swagger")
	if err != nil {
		return nil, err
	}
	handler := &http.ServeMux{}
	handler.Handle("/", http.FileServer(statikFS))
	return handler, nil
}
