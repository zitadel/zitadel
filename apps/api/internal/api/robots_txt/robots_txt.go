package robots_txt

import (
	"fmt"
	"net/http"
)

const (
	HandlerPrefix = "/robots.txt"
)

func Start() (http.Handler, error) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		fmt.Fprintf(w, "User-agent: *\nDisallow: /\n")
	})
	return handler, nil
}
