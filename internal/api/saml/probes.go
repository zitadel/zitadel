package saml

import (
	"context"
	"errors"
	"net/http"

	httphelper "github.com/caos/oidc/pkg/http"
)

type ProbesFn func(context.Context) error

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ok(w)
}

func readyHandler(probes []ProbesFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Readiness(w, r, probes...)
	}
}

func Readiness(w http.ResponseWriter, r *http.Request, probes ...ProbesFn) {
	ctx := r.Context()
	for _, probe := range probes {
		if err := probe(ctx); err != nil {
			http.Error(w, "not ready", http.StatusInternalServerError)
			return
		}
	}
	ok(w)
}

func ReadyStorage(s Storage) ProbesFn {
	return func(ctx context.Context) error {
		if s == nil {
			return errors.New("no storage")
		}
		return s.Health(ctx)
	}
}

func ok(w http.ResponseWriter) {
	httphelper.MarshalJSON(w, status{"ok"})
}

type status struct {
	Status string `json:"status,omitempty"`
}
