package openapi

import (
	"net/http"

	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"

	_ "github.com/caos/zitadel/openapi/statik"
)

const (
	HandlerPrefix = "/openapi/v2/swagger"
)

func Start() (http.Handler, error) {
	statikFS, err := fs.NewWithNamespace("swagger")
	if err != nil {
		return nil, err
	}
	handler := &http.ServeMux{}
	handler.Handle("/", cors.AllowAll().Handler(http.FileServer(statikFS)))
	return handler, nil
}
