package handler

import (
	"net/http"
	"path"
)

func (l *Login) handleResources(staticDir string) http.Handler {
	return http.StripPrefix(EndpointResources, http.FileServer(http.Dir(path.Join(staticDir, EndpointResources))))
}
