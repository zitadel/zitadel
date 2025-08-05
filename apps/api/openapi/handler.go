package openapi

import (
	"embed"
	"net/http"

	"github.com/rs/cors"
)

const (
	HandlerPrefix = "/openapi/v2/swagger"
)

//go:embed v2/zitadel/*
var openapi embed.FS

func Start() (http.Handler, error) {
	handler := &http.ServeMux{}
	handler.Handle("/", cors.AllowAll().Handler(http.FileServer(http.FS(openapi))))
	return handler, nil
}
